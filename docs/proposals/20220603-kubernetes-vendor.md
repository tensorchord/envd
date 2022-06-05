# Multi-runtime Support

Authors:
- [Ce Gao](https://github.com/gaocegege)

## Summary

envd builds and runs the environments with Docker. This proposal is to support Kubernetes as a runtime option. After this feature, envd can build the environment with Docker and run on Kubernetes.

## Motivation

Docker runtime only works for individual developers. There are many data scientists developing models on Kubernetes. In the current design, they may need to run `envd build`, then push the image to a OCI image registry manually. The image can be used after a deployment or a pod is created in Kubernetes.

We can provide native support for ease of use.

## Goals

- Run environments on Kubernetes
  - Compose Kubernetes resources (e.g. deployments, services, pods) when run `envd up`.
  - Push the image to a registry.
  - Port forward for sshd and jupyter notebook service.
  - Sync the pod file system with the local host file system.
- Manage environments on Kubernetes
  - List environments running on Kubernetes, add a new field `vendor` when run `envd get envs`
  - List dependencies of the environment on Kubernetes

## Non-Goals

- Monitoring on Kubernetes
- [Support context in envd CLI to switch runtimes between Docker and Kubernetes](https://github.com/tensorchord/envd/issues/92)
- Pause and unpause environments on Kubernetes (Kubernetes does not support the primitives)
- Build the envd environments with the buildkitd pod on Kubernetes

## Proposal

### User Stories

### Implementation Details/Notes/Constraints

## Design Details

### Test Plan

### Graduation Criteria
