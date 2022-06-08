---
sidebar_position: 1
---

# Tutorial

envd is a container-based development environment management tool for data scientists.

🐍 **No docker, only python** - Write python code to build the development environment, we help you take care of Docker.

🖨️ **Built-in jupyter/vscode** - Jupyter and VSCode remote extension are the first-class support.

⏱️ **Save time** - Better cache management to save your time, keep the focus on the model, instead of dependencies

☁️ **Local & cloud** - Run the environment locally or in the cloud, without any code change

🐳 **Container native** - Leverage container technologies but no need to learn how to use them, we optimize it for you

🤟 **Infrastructure as code** - Describe your project in a declarative way, 100% reproducible

Let's discover **envd in less than 5 minutes**.

## Getting Started

Get started by **creating a new envd environment**.

### What you'll need

- Docker (20.10.0 or above)

### Install envd

You can download the binary from the [latest release page](https://github.com/tensorchord/envd/releases/latest), and add it in `$PATH`.

After the download, please run `envd bootstrap` to bootstrap.

:::tip Tip

You can add `--dockerhub-mirror` or `-m` flag when running `envd boostrap`, to configure the mirror for docker.io registry:

```bash title="Set docker mirror"
envd bootstrap --dockerhub-mirror https://docker.mirrors.sjtug.sjtu.edu.cn
```

:::

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
