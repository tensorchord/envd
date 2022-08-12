# Daemon Process Support

Author:
- [kemingy](https://github.com/kemingy/)


## Summary

This proposal is designed to support run daemon processes in the `envd` container.

## Motivation

There can be a general solution for several use cases:

1. run a Jupyter notebook service
2. run a TensorBoard service
3. run a demo machine learning model serving service
4. run a metrics exporter

This is related to the following:

* https://github.com/tensorchord/envd/issues/527
* https://github.com/tensorchord/envd/pull/568
* https://github.com/tensorchord/envd/pull/708
* https://github.com/tensorchord/envd/pull/497

## Goals

* able to run multiple daemon processes controlled by `tini`

## API

```python
runtime.daemon(commands=[
    "python3 serving.py --port 8080",
    "watch date > /dev/null"
])
```
