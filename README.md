# Flotta device worker


This is an _unofficial_ implementation of the device-worker for [Project Flotta](https://project-flotta.io/).

## Motivation

The current implementation of device-worker does not take into consideration the resources available on the device _while_ running workloads.
This implementation aims to have a high resilience by ensuring that _workloads_ do not deplete the device of resources.
Another motivation is to make the device-worker agnostic about the work to do. The device-worker should only manage and run workloads without any other additional work like monitoring, logging, or data collection. 
All these additional tasks could be run like any other workloads.

Last but not least, this implementation *does not* use _yggdrasil_ as broker. It has a simple implementation of _yggdrasil_ API but it is a standalone in this regard.

## TO DO

- [x] Edge HTTP client (enrolling, registration, heartbeat and configuration)
- [ ] Generating *Podman* pod specification to be able to run workloads
- [ ] Execute workloads
- [ ] Integration tests

## Prerequisites

You must have flotta operator and API server running and client certificates already generated.
For more information, please read the [documentation](https://project-flotta.io/documentation/latest/intro/overview.html).

## Usage

Create a configuration yaml like:
```
LOG_LEVEL: debug
CA_ROOT: /home/cosmin/projects/device-worker-ng/resources/certificates/ca.pem
CERT: /home/cosmin/projects/device-worker-ng/resources/certificates/cert.pem
KEY: /home/cosmin/projects/device-worker-ng/resources/certificates/key.pem
SERVER: https://127.0.0.1:8043
DEVICE_ID: toto
```

Run the device-worker:

```
device-worker-ng --config config.yaml 
```

If you prefer, you can use environment variable with the prefix `EDGE_DEVICE`:
```
EDGE_DEVICE_CA_ROOT=path_to_ca EDGE_DEVICE_CERT=path_to_cert EDGE_DEVICE_KEY=path_to_key EDGE_DEVICE_SERVER=https://127.0.0.1:8043 device-worker-ng
```


