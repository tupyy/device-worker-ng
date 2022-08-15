package scheduler

import (
	"context"
	"time"

	"github.com/tupyy/device-worker-ng/internal/entity"
	"github.com/tupyy/device-worker-ng/internal/scheduler/containers"
	"go.uber.org/zap"
)

type actionType int

const (
	defaultHeartbeatPeriod = 2 * time.Second
	gracefullShutdown      = 5 * time.Second

	// action type
	runAction actionType = iota
	stopAction
)

//go:generate mockgen -package=scheduler -destination=mock_executor.go --build_flags=--mod=mod . Executor
type Executor interface {
	Run(ctx context.Context, w entity.Workload) *Future[TaskState]
	Stop(ctx context.Context, w entity.Workload)
}

type Mutator interface {
	Mutate(t *Task) (mutated bool)
}

type Scheduler struct {
	// tasks holds all the current tasks
	tasks *containers.Store[*Task]
	// futures holds the futures of executed tasks
	// the hash of the task is the key
	futures map[string]*Future[TaskState]
	// executor
	executor Executor
	// mutator is responsible with mutating task
	mutator Mutator
	// runCancel is the cancel function of the run goroutine
	runCancel context.CancelFunc
	// executionQueue holds the tasks which must be executed by executor
	executionQueue *containers.ExecutionQueue[actionType, *Task]
}

// New creates a new scheduler with the default heartbeat period of 2 seconds.
func New(executor Executor) *Scheduler {
	return newExecutor(executor, defaultHeartbeatPeriod)
}

// New creates a new scheduler with the hearbeat period provided by the user.
func NewWitHeartbeatPeriod(executor Executor, heartbeatPeriod time.Duration) *Scheduler {
	return newExecutor(executor, heartbeatPeriod)
}

func newExecutor(executor Executor, heartbeatPeriod time.Duration) *Scheduler {
	return &Scheduler{
		tasks:          containers.NewStore[*Task](),
		futures:        make(map[string]*Future[TaskState]),
		executionQueue: containers.NewExecutionQueue[actionType, *Task](),
		executor:       executor,
		mutator:        NewMutator(), // mutator with standard RestartGuard
	}
}

func (s *Scheduler) Start(ctx context.Context, input chan entity.Message, profileUpdateCh chan entity.Message) {
	runCtx, cancel := context.WithCancel(ctx)
	s.runCancel = cancel

	taskCh := make(chan entity.Option[[]entity.Workload])
	go func(ctx context.Context) {
		for {
			select {
			case message := <-input:
				switch message.Kind {
				case entity.WorkloadConfigurationMessage:
					val, ok := message.Payload.(entity.Option[[]entity.Workload])
					if !ok {
						zap.S().Errorf("mismatch message payload type. expected workload. got %v", message)
					}
					taskCh <- val
				}
			case <-ctx.Done():
				return
			}
		}
	}(runCtx)
	go s.run(runCtx, taskCh, profileUpdateCh)
}

func (s *Scheduler) Stop(ctx context.Context) {
	zap.S().Info("shutting down scheduler")

	// shutdown goroutines
	s.runCancel()

	zap.S().Info("scheduler shutdown")
}

func (s *Scheduler) run(ctx context.Context, input chan entity.Option[[]entity.Workload], profileCh chan entity.Message) {
	execution := make(chan struct{}, 1)
	doneExecutionCh := make(chan struct{})
	mutate := make(chan struct{}, 1)
	mark := make(chan struct{}, 1)

	heartbeat := time.NewTicker(defaultHeartbeatPeriod)

	for {
		select {
		case o := <-input:
			if o.None {
				// stop tasks
				iter := s.tasks.Iter()
				for iter.HasNext() {
					task, _ := iter.Next()
					if !s.isMarked(task, inactiveMark) && (task.CurrentState() == TaskStateRunning || task.CurrentState() == TaskStateDeploying) {
						s.mark(task, stopMark)
						s.mark(task, deletionMark)
					}
				}
				break
			}
			// add tasks
			for _, w := range o.Value {
				task := NewTask(w.ID(), w)
				if oldTask, found := s.tasks.FindByName(task.Name()); found {
					if oldTask.Hash() == task.Hash() {
						continue
					}
					// something changed in the workload. Stop the old one and start the new one
					zap.S().Infow("workload changed", "name", oldTask.Name)
					// stop the old one and remove it from store
					s.mark(oldTask, stopMark)
					s.mark(oldTask, deletionMark)
				}
				s.tasks.Add(task)
			}
		case <-mark:
			iter := s.tasks.Iter()
			for iter.HasNext() {
				task, _ := iter.Next()
				// poll his future if any
				if future, found := s.futures[task.ID()]; found {
					result, _ := future.Poll()
					if result.IsReady() {
						zap.S().Debugw("poll future", "id", task.Hash(), "result", result)
						s.markWithValue(task, mutateMark, result.Value)
					}

					// future is resolved when task has either been stopped or exited.
					if future.Resolved() {
						zap.S().Debugw("future resolved", "id", task.Hash())
						delete(s.futures, task.ID())
					}
					continue
				}
				// no future yet meaning the task has not been deployed yet or it exited.
				// first evaluate task. if true than deploy it.
				// evaluate the task
				if s.evaluate(task) {
					s.markWithValue(task, mutateMark, TaskStateDeploying)
				}
			}
			mutate <- struct{}{}
		case <-mutate:
			taskIter := s.tasks.Iter()
			for taskIter.HasNext() {
				task, _ := taskIter.Next()

				if !s.mutator.Mutate(task) {
					continue
				}

				// resolve the mutations
				switch task.NextState() {
				case TaskStateStopping:
					zap.S().Debugw("stop task", "id", task.Name())
					s.executionQueue.Push(stopAction, task)
				case TaskStateDeploying:
					zap.S().Debugw("deploy task", "id", task.Name())
					s.executionQueue.Push(runAction, task)
				default:
					task.MutateTo(task.NextState())
				}
			}
			if s.tasks.Len() > 0 {
				execution <- struct{}{}
			}
		case <-execution:
			// execute every task in the execution queue
			go s.execute(context.Background(), doneExecutionCh)
			// clean task marked for deletion
			s.clean()
			// stop heartbeat while we are consuming the execution queue.
			// Once is done, reset the timer.
			heartbeat.Stop()
		case <-doneExecutionCh:
			heartbeat.Reset(defaultHeartbeatPeriod)
		case <-heartbeat.C:
			mark <- struct{}{}
		case <-ctx.Done():
			return
		}
	}
}

func (s *Scheduler) execute(ctx context.Context, doneCh chan struct{}) {
	ticker := time.NewTicker(1 * time.Second)
	for s.executionQueue.Size() > 0 {
		select {
		case <-ticker.C:
			// stopping task has higher priority
			s.executionQueue.Sort(stopAction)
			action, task, err := s.executionQueue.Pop()
			if err != nil {
				zap.S().Errorw("failed to pop task from queue", "error", err)
				break
			}
			switch action {
			case stopAction:
				task.MutateTo(TaskStateStopping)
				s.executor.Stop(context.Background(), task.Workload)
			case runAction:
				task.MutateTo(TaskStateDeploying)
				future := s.executor.Run(context.Background(), task.Workload)
				s.futures[task.Hash()] = future
			}
		case <-ctx.Done():
			doneCh <- struct{}{}
			return
		}
	}
	doneCh <- struct{}{}
}

func (s *Scheduler) evaluate(t *Task) bool {
	return true
}

// remove task marked for deletion and which are stopped, exited or unknown
func (s *Scheduler) clean() {
	for {
		dirty := false
		for i := 0; i < s.tasks.Len(); i++ {
			if t, ok := s.tasks.Get(i); ok && s.isMarked(t, deletionMark) && (t.CurrentState().OneOf(TaskStateExited, TaskStateUnknown)) {
				zap.S().Debugw("task removed", "task_id", t.ID())
				s.tasks.Delete(t)
				dirty = true
				break
			}
		}
		if !dirty {
			break
		}
	}
}

func (s *Scheduler) mark(t *Task, mark string) {
	t.SetMark(mark, mark)
}

func (s *Scheduler) markWithValue(t *Task, mark string, value interface{}) {
	t.SetMark(mark, value)
}

func (s *Scheduler) isMarked(t *Task, mark string) bool {
	_, marked := t.GetMark(mark)
	return marked
}
