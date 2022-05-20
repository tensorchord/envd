<div align="center">
<h1>envd</h1>
<p>Development Environment for Data Scientists</p>
</div>

<p align=center>
<a href="https://discord.gg/KqswhpVgdU"><img alt="discord invitation link" src="https://img.shields.io/discord/974584200327991326?label=discord&style=social"></a>
<a href="https://github.com/tensorchord/envd/actions/workflows/CI.yml"><img alt="continuous integration" src="https://github.com/tensorchord/envd/actions/workflows/CI.yml/badge.svg"></a>
</p>

envd is a development environment management tool for data scientists.

:snake:  **No docker, only python** - Write python code to build the development environment, we help you take care of Docker.

:pager: **Built-in jupyter/vscode** - Provision jupyter notebooks and vscode remote in the image, remote development is possible.

:man_technologist: **Manage code and data** - Help you manage the source code and dataset in the environment

:stopwatch: **Save time** - Better cache management to save your time, keep the focus on the model, instead of dependencies

☁️ **Local & cloud** - Run the environment locally or in the cloud, without any code change

:whale: **Container native** - Leverage container technologies but no need to learn how to use them, we optimize it for you

🤟  **Infrastructure as code** - Describe your project in a declarative way, 100% reproducible

## Install

### From binary

```bash
sudo /bin/sh -c 'wget https://github.com/tensorchord/envd/releases/download/0.0.1-alpha.5/envd_0.0.1-alpha.5_Linux_x86_64 -O /usr/local/bin/envd && chmod +x /usr/local/bin/envd && /usr/local/bin/envd bootstrap'
```

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

Then you can run `envd up` and open jupyter notebook at [`http://localhost:8888`](http://localhost:8888), or open vscode remote to attach to the container.

```
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
envd > 
```

## Features

### Configure mirrors

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

envd configures Ubuntu APT source, PyPI mirror, and others in the development environment.

## Join Us

envd is backed by [TensorChord](https://github.com/tensorchord) and licensed under Apache-2.0. We are actively hiring engineers to build developer tools for machine learning practitioners in open source.

## Contribute

We welcome all kinds of contributions from the open-source community, individuals, and partners.

- Join our [discord community](https://discord.gg/KqswhpVgdU)! 
- To build from the source, check the [contributing page](./CONTRIBUTING.md).
