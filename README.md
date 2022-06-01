<div align="center">
<h1>envd</h1>
<p>Development Environment for Data Scientists</p>
</div>

<p align=center>
<a href="https://discord.gg/KqswhpVgdU"><img alt="discord invitation link" src="https://img.shields.io/discord/974584200327991326?label=discord&style=social"></a>
<a href="https://github.com/tensorchord/envd/actions/workflows/CI.yml"><img alt="continuous integration" src="https://github.com/tensorchord/envd/actions/workflows/CI.yml/badge.svg"></a>
<a href="https://trackgit.com"><img src="https://us-central1-trackgit-analytics.cloudfunctions.net/token/ping/l3ldvdaswvnjpty9u7l3" alt="trackgit-views" /></a>
</p>

envd is a development environment management tool for data scientists.

🐍 **No docker, only python** - Write python code to build the development environment, we help you take care of Docker.

🖨️ **Built-in jupyter/vscode** - Jupyter and VSCode remote extension are the first-class support.

⏱️ **Save time** - Better cache management to save your time, keep the focus on the model, instead of dependencies

☁️ **Local & cloud** - Run the environment locally or in the cloud, without any code change

🐳 **Container native** - Leverage container technologies but no need to learn how to use them, we optimize it for you

🤟 **Infrastructure as code** - Describe your project in a declarative way, 100% reproducible

## Why use envd?

It is still too difficult to configure development environments and reproduce results for data scientists and AI/ML researchers.

They have to play with Docker, conda, CUDA, GPU Drivers, and even Kubernetes if the training jobs are running in the cloud, to make things happen.

Thus, researchers have to find infra guys to help them. But the infra guys also struggle to build environments for machine learning. Infra guys love immutable infrastructure. But researchers optimize AI/ML models by trial and error. The environment will be updated, modified, or rebuilt again, and again, in place. Researchers do not have the bandwidth to be the expert on Dockerfile. They prefer `docker commit`, then the image is error-prone and hard to maintain, or debug.

envd provides another way to solve the problem. As the infra guys, we accept the reality of the differences between AI/ML and traditional workloads. We do not expect researchers to learn the basics of infrastructure, instead, we build tools to help researchers manage their development environments easily, and in a cloud-native way.

envd provides build language similar to Python and has first-class support for jupyter, vscode, and python dependencies in container technologies.

## How does envd work?

## Install

### From binary

TODO

### From source code

```bash
git clone https://github.com/tensorchord/envd
go mod tidy
make
./bin/envd --version
```

## Quickstart

Checkout the [examples](./examples/mnist), and configure envd with the manifest `build.envd`:

```python
vscode(plugins=[
    "ms-python.python",
])

base(os="ubuntu20.04", language="python3")
pip_package(name=[
    "tensorflow",
    "numpy",
])
cuda(version="11.6", cudnn="8")
shell("zsh")
jupyter(password="", port=8888)
```

Then you can run `envd up` to create the development environment.

<a href="https://asciinema.org/a/498012" target="_blank"><img src="https://asciinema.org/a/498012.svg" /></a>

TODO: illustrate that the cache will be persistent.

```
$ envd up
[+] ⌚ parse build.envd and download/cache dependencies 0.0s ✅ (finished)        
 => 💽 (cached) download oh-my-zsh                                            0.0s
 => 💽 (cached) download ms-python.python                                     0.0s
[+] 🐋 build envd environment 7.7s (24/25)                                        
 => 💽 (cached) (built-in packages) apt-get install curl openssh-client g     0.0s
 => 💽 (cached) create user group envd                                        0.0s
 => 💽 (cached) create user envd                                              0.0s
 => 💽 (cached) add user envd to sudoers                                      0.0s
 => 💽 (cached) (user-defined packages) apt-get install screenfetch           0.0s
 => 💽 (cached) install system packages                                       0.0s
 => 💽 (cached) pip install jupyter                                           0.0s
 => 💽 (cached) install PyPI packages                                         0.0s
 => 💽 (cached) install envd-ssh                                              0.0s
 => 💽 (cached) install vscode plugin ms-python.python                        0.0s
 => 💽 (cached) copy /oh-my-zsh /home/envd/.oh-my-zsh                         0.0s
 => 💽 (cached) mkfile /home/envd/install.sh                                  0.0s
 => 💽 (cached) install oh-my-zsh                                             0.0s
...
# You are in the docker container for dev
(envd 🐳)  ➜  mnist-dev git:(master) python3 ./main.py
...
```

Jupyter notebook service and sshd server are running inside the container. You can use jupyter or vscode remote-ssh extension to develop AI/ML models.

```
$ envd ls
NAME            JUPYTER                 SSH TARGET      GPU     STATUS  CONTAINER ID 
mnist           http://localhost:8888   mnist.envd      true    running 253f656b8c40
```

## Features

### Configure mirrors

envd supports PyPI mirror and apt source configuration. You can configure them in `build.env` or `$HOME/.config/envd/config.envd` to set up in all environments.

```text
cat ~/.config/envd/config.envd
ubuntu_apt(source="""
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
pip_mirror(mirror = "https://mirror.sjtu.edu.cn/pypi/web/simple")
vscode(plugins = [
    "ms-python.python",
    "github.copilot"
])
```

## Join Us

envd is backed by [TensorChord](https://github.com/tensorchord) and licensed under Apache-2.0. We are actively hiring engineers to build developer tools for machine learning practitioners in open source.

## Contribute

We welcome all kinds of contributions from the open-source community, individuals, and partners.

- Join our [discord community](https://discord.gg/KqswhpVgdU)! 
- To build from the source, check the [contributing page](./CONTRIBUTING.md).
