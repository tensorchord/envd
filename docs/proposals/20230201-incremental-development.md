## Summary

This proposal is intended to support incremental development in the `envd` environment.

## Motivation

When users start working on new projects, they don't really know all the Python and system packages they need. Typically, they will need to try different libraries and choose the best one. Then they update the `build.envd` file. This allows them to reproduce the environment. We could help improve this process.

There are several requirements for incremental development:

1. pip
2. apt
3. expose port (can use SSH tunneling, may require re-login)
    - This is useful if users want to expose some services like tensorboard, demo service, etc.
4. mount (optional, cannot be done without restarting the container)
    - useful for caching, like huggingface's models

## Goals

* Track manually installed packages in the environment
* Help users update the `build.envd` file.

A hidden benefit of this proposal is that users can quickly start from a base image (as long as that image has been downloaded).

This is related to the following issue:

* https://github.com/tensorchord/envd/issues/9

## Solution

1. provide the CLI to record the installation commands

```bash
envd add pip --packages numpy "jax[cuda] -f https://storage.googleapis.com/jax-releases/jax_cuda_releases.html"
envd add apt --packages htop screenfetch
envd add expose --envd 8000 --host 9000 --service demo --addr 0.0.0.0
```

Many ML packages have their own index services. This way, some pip arguments can be passed in a string. So it's also possible to reproduce it like this:

```python
install.python_packages(name=[
    "numpy",
    "jax[cuda] -f https://storage.googleapis.com/jax-releases/jax_cuda_releases.html",
])
```

2. add to the `build.envd` file

This part can be tricky. Because we need to parse the `build.envd` file and get the function position to insert these new configurations. Fortunately, starlark already supports this: [Function.Position](https://pkg.go.dev/go.starlark.net/starlark#Function.Position).

Back to our goals, we want to help users to update the `build.envd`. This means we don't have to insert it in the exact right place in the `build.envd`. Other options:

- Provide a lock file
  - JSON file should be easy to add/delete/update
  - need another CLI like `envd sync` to generate a list of `envd` function calls for users to add to their `build.envd` file

## Other solutions

- See https://github.com/tensorchord/envd/issues/9#issuecomment-1383694132
  - the dependency tree may lack of some information like where to find these packages
- We should be able to get the command history and extract the relevant lines for users to review.
  - it's hard to generate the `envd` function calls, so users may need to take more effort to find out which one they need

## Others

- For the docker runner, since it will mount the working directory, these changes will be written to the disk directly.
- For the envd-server runner, users may need to sync changes by using git.

## TBD

* Should we use a lock file or not? If so, should we use a JSON file?
