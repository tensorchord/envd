<div align="center">
<h1>envd</h1>
<p>Development Environment for Data Scientists</p>
</div>

<p align=center>
<a href="https://discord.gg/KqswhpVgdU"><img alt="discord invitation link" src="https://img.shields.io/discord/974584200327991326?label=discord&style=social"></a>
<a href="https://github.com/tensorchord/envd/actions/workflows/CI.yml"><img alt="continuous integration" src="https://github.com/tensorchord/envd/actions/workflows/CI.yml/badge.svg"></a>
<a href="https://trackgit.com"><img src="https://us-central1-trackgit-analytics.cloudfunctions.net/token/ping/l3ldvdaswvnjpty9u7l3" alt="trackgit-views" /></a>
</p>

> **⚠️ envd is still under heavy development, and subject to change. it is not feature-complete or production-ready. Please contact us in [discord](https://discord.gg/KqswhpVgdU) if there is any problem.**

envd is a development environment management tool for data scientists.

🐍 **No docker, only python** - Write python code, we help you take care of Docker, CUDA, SSH keys and so on.

🖨️ **Built-in jupyter/vscode** - Jupyter and VSCode remote extension are first-class support.

⏱️ **Save time** - Better cache management to save your time, keep the focus on the model, instead of dependencies.

☁️ **Local & cloud** - Run the environment locally or in the cloud, without any code change.

🐳 **Container native** - Leverage container technologies but no need to learn how to use them, we optimize it for you.

🤟 **Infrastructure as code** - Describe your project in a declarative way, 100% reproducible.

## Why use envd?

It is still too difficult to configure development environments and reproduce results for data scientists and AI/ML researchers.

They have to play with Docker, conda, CUDA, GPU Drivers, and even Kubernetes if the training jobs are running in the cloud, to make things happen.

Thus, researchers have to find infra guys to help them. But the infra guys also struggle to build environments for machine learning. Infra guys love immutable infrastructure. But researchers optimize AI/ML models by trial and error. The environment will be updated, modified, or rebuilt again, and again, in place. Researchers do not have the bandwidth to be the expert on Dockerfile. They prefer `docker commit`, then the image is error-prone and hard to maintain, or debug.

envd provides another way to solve the problem. As the infra guys, we accept the reality of the differences between AI/ML and traditional workloads. We do not expect researchers to learn the basics of infrastructure, instead, we build tools to help researchers manage their development environments easily, and in a cloud-native way.

envd provides build language similar to Python and has first-class support for jupyter, vscode, and python dependencies in container technologies.

## Documentation

See [envd documentation](https://envd.tensorchord.ai/docs/intro).

## Getting Started

Get started by **creating a new envd environment**.

### What you'll need

- Docker (20.10.0 or above)

### Install envd

You can download the binary from the [latest release page](https://github.com/tensorchord/envd/releases/latest), and add it in `$PATH`.

After the download, please run `envd bootstrap` to bootstrap.

> You can add `--dockerhub-mirror` or `-m` flag when running `envd boostrap`, to configure the mirror for docker.io registry:
>
>```bash title="Set docker mirror"
>envd bootstrap --dockerhub-mirror https://docker.mirrors.sjtug.sjtu.edu.cn
>```

### Create an envd environment

Please clone the [`envd-quick-start`](https://github.com/tensorchord/envd-quick-start):

```
git clone https://github.com/tensorchord/envd-quick-start.git
```

The build manifest `build.envd` looks like:

```python title=build.envd
def build():
    base(os="ubuntu20.04", language="python3")
    install.python_packages(name = [
        "numpy",
    ])
    shell("zsh")
```

Then please run the command below to setup a new environment:

```
cd envd-quick-start && envd up
```

```
$ cd envd-quick-start && envd up
[+] ⌚ parse build.envd and download/cache dependencies 2.8s ✅ (finished)     
 => download oh-my-zsh                                                    2.8s 
[+] 🐋 build envd environment 18.3s (25/25) ✅ (finished)                      
 => create apt source dir                                                 0.0s 
 => local://cache-dir                                                     0.1s 
 => => transferring cache-dir: 5.12MB                                     0.1s 
...
 => pip install numpy                                                    13.0s 
 => copy /oh-my-zsh /home/envd/.oh-my-zsh                                 0.1s 
 => mkfile /home/envd/install.sh                                          0.0s 
 => install oh-my-zsh                                                     0.1s 
 => mkfile /home/envd/.zshrc                                              0.0s 
 => install shell                                                         0.0s
 => install PyPI packages                                                 0.0s
 => merging all components into one                                       0.3s
 => => merging                                                            0.3s
 => mkfile /home/envd/.gitconfig                                          0.0s 
 => exporting to oci image format                                         2.4s 
 => => exporting layers                                                   2.0s 
 => => exporting manifest sha256:7dbe9494d2a7a39af16d514b997a5a8f08b637f  0.0s
 => => exporting config sha256:1da06b907d53cf8a7312c138c3221e590dedc2717  0.0s
 => => sending tarball                                                    0.4s
(envd) ➜  demo git:(master) ✗ # You are in the container-based environment!
```

### Play with the environment

You can run `ssh envd-quick-start.envd` to reconnect if you exit from the environment. Or you can execute `git` or `python` commands inside.

```bash
$ python demo.py
[2 3 4]
$ git fetch
$
```

### Setup jupyter notebook

Please edit the `build.envd` to enable jupyter notebook:

```python title=build.envd
def build():
    base(os="ubuntu20.04", language="python3")
    install.python_packages(name = [
        "numpy",
    ])
    shell("zsh")
    config.jupyter(password="", port=8888)
```

You can get the endpoint of jupyter notebook via `envd get envs`.

```bash
$ envd up --detach
$ envd get env
NAME                    JUPYTER                 SSH TARGET              CONTEXT                                 IMAGE                   GPU     CUDA    CUDNN   STATUS          CONTAINER ID 
envd-quick-start        http://localhost:8888   envd-quick-start.envd   /home/gaocegege/code/envd-quick-start   envd-quick-start:dev    false   <none>  <none>  Up 54 seconds   bd3f6a729e94
```

## Features

### Pause and resume

```
$ envd pause --env mnist
mnist
$ env get envs
NAME         JUPYTER                 SSH TARGET   CONTEXT  IMAGE      GPU  CUDA  CUDNN  STATUS              CONTAINER ID 
mnist        http://localhost:9999   mnist.envd   /mnist   mnist:dev  true 11.6  8      Up 23 hours(Paused) 74a9f1007004
$ envd resume --env mnist
$ ssh mnist.envd
(envd 🐳) $ # The environment is resumed!
```

### Configure mirrors

envd supports PyPI mirror and apt source configuration. You can configure them in `build.env` or `$HOME/.config/envd/config.envd` to set up in all environments.

```text
cat ~/.config/envd/config.envd
config.apt_source(source="""
deb https://mirror.sjtu.edu.cn/ubuntu focal main restricted
deb https://mirror.sjtu.edu.cn/ubuntu focal-updates main restricted
deb https://mirror.sjtu.edu.cn/ubuntu focal universe
deb https://mirror.sjtu.edu.cn/ubuntu focal-updates universe
deb https://mirror.sjtu.edu.cn/ubuntu focal multiverse
deb https://mirror.sjtu.edu.cn/ubuntu focal-updates multiverse
deb https://mirror.sjtu.edu.cn/ubuntu focal-backports main restricted universe multiverse
deb http://archive.canonical.com/ubuntu focal partner
deb https://mirror.sjtu.edu.cn/ubuntu focal-security main restricted universe multiverse
""")
config.pip_index(url = "https://mirror.sjtu.edu.cn/pypi/web/simple")
install.vscode_extensions([
    "ms-python.python",
    "github.copilot"
])
```

## Contribute

We welcome all kinds of contributions from the open-source community, individuals, and partners.

- Join our [discord community](https://discord.gg/KqswhpVgdU)! 
- To build from the source, check the [contributing page](./CONTRIBUTING.md).
