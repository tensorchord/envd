# MIDI

Development Environment for Data Scientists

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
sudo /bin/sh -c 'wget https://github.com/tensorchord/midi/releases/download/0.0.1-alpha.3/midi_0.0.1-alpha.3_Linux_x86_64 -O /usr/local/bin/midi && chmod +x /usr/local/bin/midi && /usr/local/bin/midi bootstrap'
```

### From source code

```bash
git clone https://github.com/tensorchord/midi
go mod tidy
make
./bin/midi --version
```

## Quickstart

Checkout the [examples](./examples/mnist), and configure MIDI with the manifest `build.MIDI`:

```
vscode(plugins = [
    "ms-python.python-2021.12.1559732655",
])

base(os="ubuntu20.04", language="python3")
pip_package(name = [
    "tensorflow",
    "numpy",
])
cuda(version="11.6", cudnn="8")
shell("zsh")
jupyter(password="", port=8888)
```

Then you can run `midi up` and open jupyter notebook at [`http://localhost:8888`](http://localhost:8888), or open vscode remote to attach to the container.

## Features

### Configure mirrors

```
cat ~/.config/midi/config.MIDI
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

MIDI configures Ubuntu APT source, PyPI mirror, and others in the development environment.
