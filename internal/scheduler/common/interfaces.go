package common

import (
	"context"

	"github.com/tupyy/device-worker-ng/internal/entity"
)

type Reconciler interface {
	Reconcile(ctx context.Context, job *entity.Job, executor Executor) *entity.Future[entity.Result[entity.JobState]]
}

//go:generate mockgen -package=scheduler -destination=mock_executor.go --build_flags=--mod=mod . Executor
type Executor interface {
	Remove(ctx context.Context, w entity.Workload) error
	Run(ctx context.Context, w entity.Workload) error
	Stop(ctx context.Context, w entity.Workload) error
	GetState(ctx context.Context, w entity.Workload) (entity.JobState, error)
	Exists(ctx context.Context, w entity.Workload) (bool, error)
}
