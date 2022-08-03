# Service Expose Support
Authors:
- [nullday](https://github.com/aseaday)

## Summary

This proposal is designed for exposing service from the application running in the envd container such as Juypter.We introduce a built-in function `expose` for achieving this goal.

## Motivation

There a lot of AI/ML SDK can be accessible from Web or GUI Client. Envd users need to configure their container to use these tools better. For example, a user want to write some code on Juypter in their Chrome.

In the future, we should also keep the possibility for remote or cloud access.

## Goals
- http/tcp 
- OS support:
    - Linux/Mac
    - wsl access from the windows

## API
```
expose(local_port=<LPORT>, host_port=<RPORT>, svc=<SERVICE_NAME>)
```
It will bind the host port to the local port.
