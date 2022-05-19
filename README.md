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
sudo /bin/sh -c 'wget https://github.com/tensorchord/envd/releases/download/0.0.1-alpha.3/envd_0.0.1-alpha.3_Linux_x86_64 -O /usr/local/bin/envd && chmod +x /usr/local/bin/envd && /usr/local/bin/envd bootstrap'
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
    "ms-python.python-2021.12.1559732655",
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
[+] ⌚ parse build.envd and download/cache dependencies 23.7s ✅ (finished)        
 => download oh-my-zsh                                                       13.5s 
 => download ms-vscode.cpptools-1.7.1                                         2.1s 
 => download github.copilot-1.12.5517                                         1.3s 
 => download ms-python.python-2021.12.1559732655                              6.7s 
[+] 🐋 build envd environment 4.4s (23/24)                                         
 => 💽 (cached) copy /ms-python.python-2021.12.1559732655/extension /root/     0.0s
 => 💽 (cached) merge (copy /ms-vscode.cpptools-1.7.1/extension /root/.vsc     0.0s
 => 💽 (cached) mkfile /etc/apt/sources.list                                   0.0s
 => 💽 (cached) merge (docker-image://docker.io/nvidia/cuda:11.6.0-cudnn8-     0.0s
 => 💽 (cached) mkfile /etc/pip.conf                                           0.0s
 => 💽 (cached) merge (merge (docker-image://docker.io/nvidia/cuda:11.6.0-     0.0s
 => 💽 (cached) sh -c apt-get update && apt-get install -y --no-install-re     0.0s
 => 💽 (cached) pip install tensorflow numpy jupyter                           0.0s
 => 💽 (cached) diff (sh -c apt-get update && apt-get install -y --no-inst     0.0s
...
# You are in the docker container for dev
envd > 
```

## Features

### Configure mirrors

```
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
    "ms-vscode.cpptools-1.7.1",
    "github.copilot-1.12.5517"
])
```

envd configures Ubuntu APT source, PyPI mirror, and others in the development environment.

## Join Us

envd is backed by [TensorChord](https://github.com/tensorchord) and licensed under Apache-2.0. We are actively hiring engineers to build developer tools for machine learning practitioners in open source.

## Contribute

We welcome all kinds of contributions from the open-source community, individuals, and partners.

- Join our [discord community](https://discord.gg/KqswhpVgdU)! 
- To build from source, check the [contributing page](./CONTRIBUTING.md).
