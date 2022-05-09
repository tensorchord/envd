# MIDI

Development Environment for Data Scientists

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

Checkout [examples](./examples/mnist) and run

```
midi up
```

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
