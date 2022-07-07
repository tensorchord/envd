# Service Expose Support
Authors:
- [nullday](https://github.com/aseaday)

## Summary

This proposal is designed for exposing service from the environment running application such as Juypter for external use. It introduce a first-level function `serve` and its domain related packages:

- host
- cloud

And maybe a utils package for the support of define the default service port or socket predefined:

- presets

## Motivation

There a lot of AI/ML develop toolkits could be accessible from Web or GUI Client.  Users need a interface to configure the service which enviroment runtime could expose. The juypter is a apparent case need to be supported.

There is also such the kind of demand when users run their env in the cloud which accessing the service could be a little complex depends on the cloud conditions. For example, ops team allows exposing the service or ingress or only permit a port-forward.

## Goals
- http/tcp support
- local support:
    - linux/mac local access
    - wsl access from the windows
- remote support:
    - k8s port-foward support
    - k8s service/node-port support for http

Note: There is no need for support non-tcp protocl for now.

## Proposal

Add `def serve` as the first level function as `def build` but could be omitted.

```python
def build():
    base(os="ubuntu20.04", language="python")
    # Use `config.jupyter(password="")` 
    # if you do not need to set up password.
    config.jupyter(password="password")

def server():
    pass
```

import two package `host` and `cloud` in the use of `server` function. The follow example give a user story user want to export juypter port to the host port 1234:

```python
def build():
    base(os="ubuntu20.04", language="python")
    # Use `config.jupyter(password="")` 
    # if you do not need to set up password.
    config.jupyter(password="password")
def serve():
    host.bind(env_port=presets.juypter.port,host_port=1234)
```

for the cloud user scenario:


```python
def build():
    base(os="ubuntu20.04", language="python")
    # Use `config.jupyter(password="")` 
    # if you do not need to set up password.
    config.jupyter(password="password")
def serve():
    cloud.k8s.bind(env_port=presets.juypter.port,host_port=1234) # It would use port forward
    cloud.k8s.bind(env_port=presets.juypter.port, type=cloud.k8s.clusterIP, service_port=1234)
```
 It would expose service with type clusterIP, we could support clusterIP, NodePort and LoadBalancer. And for those k8s cluster auto config ingress with clusterIP and HTTP, It will be very useful.