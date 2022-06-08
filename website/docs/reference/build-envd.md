---
sidebar_position: 1
---

# build.envd Reference

envd can build images automatically by reading the instructions from a `build.envd`. A `build.envd` is a text document that contains all the commands a user could call on the command line to assemble an image.

## Usage

The `envd build` command builds an image from a `build.envd`. `envd up` command builds an image and runs immediately. Traditionally, the `build.envd` and located in the root of the context.  You use the `--path` or `-p` flag with `envd build`/`envd up` to point to a directory anywhere in your file system which contains a `build.envd`.

```bash
$ ls
build.envd ...
$ envd build
```

```bash
# Or you can specify the path.
$ tree .
./examples
└── mnist
    ├── build.envd
    ├── main.py
    ├── mnist.ipynb
    └── README.md
$ envd build --path examples/mnist
```

## Format

The syntax of `build.envd` is [starlark](https://docs.bazel.build/versions/main/skylark/language.html), which is inspired by Python3. If you know Python, then you can write `build.envd` without an issue.

Here is an example of `build.envd`:

```python
def build():
    # Here is a comment.
    base(os="ubuntu20.04", language="python3")
    install.python_packages(name = [
        "numpy",
    ])
    shell("zsh")
```

`build` is the default function that envd executes. In this function, `base()`, `install.python_packages()` and `shell()` are invoked. `base()` configures the operating system and programming language in the environment. `install.python_packages` declares the python packages that will be used in the environment. `shell` specifies which shell will be used.

### build

`build` is the default function name in `build.envd`. envd only invokes `build` function if you run `envd build` or `envd up`.

:::caution

**A `build.envd` must have a `build` function**.

:::

```python title=build.envd
def build():
    ...
```

### base



```python title=build.envd
def build():
    base(os="ubuntu20.04", language="python3")
```
