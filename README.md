<div align="center">
<img src="https://user-images.githubusercontent.com/12974685/200007223-cd94fe9a-266d-4bbd-ac23-e71043d0c3dc.svg#gh-light-mode-only" alt="envd cat wink"/>
<img src="https://user-images.githubusercontent.com/12974685/200007265-4e47ff2c-c2a0-4e77-baaa-760ee8728388.svg#gh-dark-mode-only" alt="envd cat wink"/>

<p>Development environment for AI/ML</p>
</div>

<p align=center>
<a href="https://discord.gg/KqswhpVgdU"><img alt="discord invitation link" src="https://dcbadge.vercel.app/api/server/KqswhpVgdU?style=flat"></a>
<a href="https://twitter.com/TensorChord"><img src="https://img.shields.io/twitter/follow/tensorchord?style=social" alt="trackgit-views" /></a>
<a href="https://pypi.org/project/envd"><img src="https://img.shields.io/pypi/pyversions/envd" alt="Python Version" /></a>
<a href="https://github.com/tensorchord/envd#contributors-"><img alt="all-contributors" src="https://img.shields.io/github/all-contributors/tensorchord/envd/main"></a>
<a href="https://pypi.org/project/envd/"><img alt="envd package donwloads" src="https://static.pepy.tech/personalized-badge/envd?period=month&units=international_system&left_color=grey&right_color=brightgreen&left_text=downloads/month"</a>
<a href="https://github.com/tensorchord/envd/actions/workflows/CI.yml"><img alt="continuous integration" src="https://github.com/tensorchord/envd/actions/workflows/CI.yml/badge.svg"></a>
<a href='https://coveralls.io/github/tensorchord/envd?branch=main'><img src='https://coveralls.io/repos/github/tensorchord/envd/badge.svg?branch=main' alt='Coverage Status' /></a>
</p>

## What is envd?

envd (`ɪnˈvdɪ`) is a command-line tool that helps you create the container-based development environment for AI/ML.

Creating development environments is not easy, especially with today's complex systems and dependencies. With everything from Python to CUDA, BASH scripts, and Dockerfiles constantly breaking, it can feel like a nightmare - until now!

Instantly get your environment running exactly as you need with a simple declaration of the packages you seek in build.envd and just one command: `envd up`!

<p align="center">
  <img src="https://user-images.githubusercontent.com/5100735/207217321-34c30dde-4b55-4871-b6fe-f9fc6ec19986.svg" width="75%"/>
</p>

## Why use `envd`?

Environments built with `envd` provide the following features out-of-the-box:

**Simple CLI and language**

`envd` enables you to quickly and seamlessly integrate powerful CLI tools into your existing Python workflow to provision your programming environment without learning a new language or DSL.

```python
def build():
    install.python_packages(name = [
        "numpy",
    ])
    shell("zsh")
    config.jupyter()
```

**Isolation, compatible with OCI image**

With `envd`, users can create an isolated space to train, fine-tune, or serve. By utilizing sophisticated virtualization technology as well as other features like [buildkit](https://github.com/moby/buildkit), it's an ideal solution for environment setup.

`envd` environment image is compatible with [OCI image specification](https://github.com/opencontainers/image-spec). By leveraging the power of an OCI image, you can make your environment available to anyone and everyone! Make it happen with a container registry like Harbor or Docker Hub.

**Local, and cloud**

`envd` can now be used on a hybrid platform, ranging from local machines to clusters hosted by Kubernetes. Any of these options offers an efficient and versatile way for developers to create their projects!

```sh
$ envd context use local
# Run envd environments locally
$ envd up
...
$ envd context use cluster
# Run envd environments in the cluster with the same experience
$ envd up
```

Check out the [doc](https://envd.tensorchord.ai/teams/kubernetes.html) for more details.

**Build anywhere, faster**

`envd` offers a wealth of advantages, such as remote build and software caching capabilities like pip index caches or apt cache, with the help of [buildkit](https://github.com/moby/buildkit) - all designed to make your life easier without ever having to step foot in the code itself!

Reusing previously downloaded packages from the PyPI/APT cache saves time and energy, making builds more efficient. No need to redownload what was already acquired before – a single download is enough for repeat usage! 

With Dockerfile v1, users are unable to take advantage of PyPI caching for faster installation speeds - but `envd` offers this support and more!

<p align=center>
  <img src="https://user-images.githubusercontent.com/5100735/189928628-543f4851-87b7-462b-b811-372cbf46ff25.svg#gh-light-mode-only" width="65%"/>
</p>

<p align=center>
  <img src="https://user-images.githubusercontent.com/16186646/197944452-4a5dcd5f-68d0-4505-b419-e95c298978d7.svg#gh-dark-mode-only" width="65%"/>
</p>

Besides, `envd` also supports remote build, which means you can build your environment on a remote machine, such as a cloud server, and then push it to the registry. This is especially useful when you are working on a machine with limited resources, or when you expect a build machine with higher performance.

**Knowledge reuse in your team**

Forget copy-pasting Dockerfile instructions - use envd to easily build functions and reuse them by importing any Git repositories with the `include` function! Craft powerful custom solutions quickly.

```python
envdlib = include("https://github.com/tensorchord/envdlib")

def build():
    base(os="ubuntu20.04", language="python")
    envdlib.tensorboard(host_port=8888)
```

<details>
  <summary><code>envdlib.tensorboard</code> is defined in <a href="https://github.com/tensorchord/envdlib/blob/main/src/monitoring.envd">github.com/tensorchord/envdlib</a></summary>

```python
def tensorboard(
    envd_port=6006,
    envd_dir="/home/envd/logs",
    host_port=0,
    host_dir="/tmp",
):
    """Configure TensorBoard.

    Make sure you have permission for `host_dir`

    Args:
        envd_port (Optional[int]): port used by envd container
        envd_dir (Optional[str]): log storage mount path in the envd container
        host_port (Optional[int]): port used by the host, if not specified or equals to 0,
            envd will randomly choose a free port
        host_dir (Optional[str]): log storage mount path in the host
    """
    install.python_packages(["tensorboard"])
    runtime.mount(host_path=host_dir, envd_path=envd_dir)
    runtime.daemon(
        commands=[
            [
                "tensorboard",
                "--logdir",
                envd_dir,
                "--port",
                str(envd_port),
                "--host",
                "0.0.0.0",
            ],
        ]
    )
    runtime.expose(envd_port=envd_port, host_port=host_port, service="tensorboard")
```
</details>

## Getting Started 🚀

### Requirements

- Docker (20.10.0 or above)

### Install and bootstrap `envd`

`envd` can be installed with `pip`, or you can download the binary [release](https://github.com/tensorchord/envd/releases) directly. After the installation, please run `envd bootstrap` to bootstrap.

```bash
pip3 install --upgrade envd
```

After the installation, please run `envd bootstrap` to bootstrap:

```bash
envd bootstrap
```

Read the [documentation](https://envd.tensorchord.ai/guide/getting-started.html#install-and-bootstrap-envd) for more alternative installation methods.

> You can add `--dockerhub-mirror` or `-m` flag when running `envd bootstrap`, to configure the mirror for docker.io registry:
>
>```bash title="Set docker mirror"
>envd bootstrap --dockerhub-mirror https://docker.mirrors.sjtug.sjtu.edu.cn
>```

### Create an `envd` environment

Please clone the [`envd-quick-start`](https://github.com/tensorchord/envd-quick-start):

```bash
git clone https://github.com/tensorchord/envd-quick-start.git
```

The build manifest `build.envd` looks like:

```python title=build.envd
def build():
    base(os="ubuntu20.04", language="python3")
    # Configure the pip index if needed.
    # config.pip_index(url = "https://pypi.tuna.tsinghua.edu.cn/simple")
    install.python_packages(name = [
        "numpy",
    ])
    shell("zsh")
```

*Note that we use Python here as an example but please check out examples for other languages such as R and Julia [here](https://github.com/tensorchord/envd/tree/main/examples).*

Then please run the command below to set up a new environment:

```bash
cd envd-quick-start && envd up
```

```bash
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
envd-quick-start via Py v3.9.13 via 🅒 envd
⬢ [envd]❯ # You are in the container-based environment!
```

### Set up Jupyter notebook

Please edit the `build.envd` to enable jupyter notebook:

```python title=build.envd
def build():
    base(os="ubuntu20.04", language="python3")
    # Configure the pip index if needed.
    # config.pip_index(url = "https://pypi.tuna.tsinghua.edu.cn/simple")
    install.python_packages(name = [
        "numpy",
    ])
    shell("zsh")
    config.jupyter()
```

You can get the endpoint of the running Jupyter notebook via `envd envs ls`.

```bash
$ envd up --detach
$ envd envs ls
NAME                    JUPYTER                 SSH TARGET              CONTEXT                                 IMAGE                   GPU     CUDA    CUDNN   STATUS          CONTAINER ID
envd-quick-start        http://localhost:42779   envd-quick-start.envd   /home/gaocegege/code/envd-quick-start   envd-quick-start:dev    false   <none>  <none>  Up 54 seconds   bd3f6a729e94
```

## More on documentation 📝

See [envd documentation](https://envd.tensorchord.ai/guide/getting-started.html).

## Roadmap 🗂️

Please checkout [ROADMAP](https://envd.tensorchord.ai/community/roadmap.html).

## Contribute 😊

We welcome all kinds of contributions from the open-source community, individuals, and partners.

- Join our [discord community](https://discord.gg/KqswhpVgdU)!
- To build from the source, please read our [contributing documentation](https://envd.tensorchord.ai/community/contributing.html) and [development tutorial](https://envd.tensorchord.ai/developers/development.html).

[![Open in Gitpod](https://gitpod.io/button/open-in-gitpod.svg)](https://gitpod.io/#https://github.com/tensorchord/envd)

## Contributors ✨

Thanks goes to these wonderful people ([emoji key](https://allcontributors.org/docs/en/emoji-key)):

<!-- ALL-CONTRIBUTORS-LIST:START - Do not remove or modify this section -->
<!-- prettier-ignore-start -->
<!-- markdownlint-disable -->
<table>
  <tbody>
    <tr>
      <td align="center" valign="top" width="14.28%"><a href="http://blog.duanfei.org"><img src="https://avatars.githubusercontent.com/u/16186646?v=4?s=70" width="70px;" alt=" Friends A."/><br /><sub><b> Friends A.</b></sub></a><br /><a href="https://github.com/tensorchord/envd/commits?author=shaonianche" title="Documentation">📖</a> <a href="#design-shaonianche" title="Design">🎨</a></td>
      <td align="center" valign="top" width="14.28%"><a href="https://github.com/aaronzs"><img src="https://avatars.githubusercontent.com/u/1827365?v=4?s=70" width="70px;" alt="Aaron Sun"/><br /><sub><b>Aaron Sun</b></sub></a><br /><a href="#userTesting-aaronzs" title="User Testing">📓</a> <a href="https://github.com/tensorchord/envd/commits?author=aaronzs" title="Code">💻</a></td>
      <td align="center" valign="top" width="14.28%"><a href="https://github.com/popfido"><img src="https://avatars.githubusercontent.com/u/3928409?v=4?s=70" width="70px;" alt="Aka.Fido"/><br /><sub><b>Aka.Fido</b></sub></a><br /><a href="#platform-popfido" title="Packaging/porting to new platform">📦</a> <a href="https://github.com/tensorchord/envd/commits?author=popfido" title="Documentation">📖</a> <a href="https://github.com/tensorchord/envd/commits?author=popfido" title="Code">💻</a></td>
      <td align="center" valign="top" width="14.28%"><a href="http://alexhxi.com"><img src="https://avatars.githubusercontent.com/u/68758451?v=4?s=70" width="70px;" alt="Alex Xi"/><br /><sub><b>Alex Xi</b></sub></a><br /><a href="https://github.com/tensorchord/envd/commits?author=AlexXi19" title="Code">💻</a></td>
      <td align="center" valign="top" width="14.28%"><a href="https://github.com/LuBingtan"><img src="https://avatars.githubusercontent.com/u/30698342?v=4?s=70" width="70px;" alt="Bingtan Lu"/><br /><sub><b>Bingtan Lu</b></sub></a><br /><a href="https://github.com/tensorchord/envd/commits?author=LuBingtan" title="Code">💻</a></td>
      <td align="center" valign="top" width="14.28%"><a href="https://github.com/sunby"><img src="https://avatars.githubusercontent.com/u/9817127?v=4?s=70" width="70px;" alt="Bingyi Sun"/><br /><sub><b>Bingyi Sun</b></sub></a><br /><a href="https://github.com/tensorchord/envd/commits?author=sunby" title="Code">💻</a></td>
      <td align="center" valign="top" width="14.28%"><a href="http://gaocegege.com/Blog"><img src="https://avatars.githubusercontent.com/u/5100735?v=4?s=70" width="70px;" alt="Ce Gao"/><br /><sub><b>Ce Gao</b></sub></a><br /><a href="https://github.com/tensorchord/envd/commits?author=gaocegege" title="Code">💻</a> <a href="https://github.com/tensorchord/envd/commits?author=gaocegege" title="Documentation">📖</a> <a href="#design-gaocegege" title="Design">🎨</a> <a href="#projectManagement-gaocegege" title="Project Management">📆</a></td>
    </tr>
    <tr>
      <td align="center" valign="top" width="14.28%"><a href="https://frostming.com"><img src="https://avatars.githubusercontent.com/u/16336606?v=4?s=70" width="70px;" alt="Frost Ming"/><br /><sub><b>Frost Ming</b></sub></a><br /><a href="https://github.com/tensorchord/envd/commits?author=frostming" title="Code">💻</a> <a href="https://github.com/tensorchord/envd/commits?author=frostming" title="Documentation">📖</a></td>
      <td align="center" valign="top" width="14.28%"><a href="https://GuangyangLi.com"><img src="https://avatars.githubusercontent.com/u/2060045?v=4?s=70" width="70px;" alt="Guangyang Li"/><br /><sub><b>Guangyang Li</b></sub></a><br /><a href="https://github.com/tensorchord/envd/commits?author=gyli" title="Code">💻</a></td>
      <td align="center" valign="top" width="14.28%"><a href="https://github.com/Gui-Yue"><img src="https://avatars.githubusercontent.com/u/78520005?v=4?s=70" width="70px;" alt="Gui-Yue"/><br /><sub><b>Gui-Yue</b></sub></a><br /><a href="https://github.com/tensorchord/envd/commits?author=Gui-Yue" title="Code">💻</a></td>
      <td align="center" valign="top" width="14.28%"><a href="https://github.com/haiker2011"><img src="https://avatars.githubusercontent.com/u/8073429?v=4?s=70" width="70px;" alt="Haiker Sun"/><br /><sub><b>Haiker Sun</b></sub></a><br /><a href="https://github.com/tensorchord/envd/commits?author=haiker2011" title="Code">💻</a></td>
      <td align="center" valign="top" width="14.28%"><a href="https://bandism.net/"><img src="https://avatars.githubusercontent.com/u/22633385?v=4?s=70" width="70px;" alt="Ikko Ashimine"/><br /><sub><b>Ikko Ashimine</b></sub></a><br /><a href="https://github.com/tensorchord/envd/commits?author=eltociear" title="Code">💻</a></td>
      <td align="center" valign="top" width="14.28%"><a href="https://github.com/nasnoisaac"><img src="https://avatars.githubusercontent.com/u/11145462?v=4?s=70" width="70px;" alt="Isaac "/><br /><sub><b>Isaac </b></sub></a><br /><a href="https://github.com/tensorchord/envd/commits?author=nasnoisaac" title="Code">💻</a></td>
      <td align="center" valign="top" width="14.28%"><a href="https://belyenochi.github.io/"><img src="https://avatars.githubusercontent.com/u/26409132?v=4?s=70" width="70px;" alt="JasonZhu"/><br /><sub><b>JasonZhu</b></sub></a><br /><a href="https://github.com/tensorchord/envd/commits?author=Belyenochi" title="Code">💻</a></td>
    </tr>
    <tr>
      <td align="center" valign="top" width="14.28%"><a href="https://github.com/knight42"><img src="https://avatars.githubusercontent.com/u/4237254?v=4?s=70" width="70px;" alt="Jian Zeng"/><br /><sub><b>Jian Zeng</b></sub></a><br /><a href="#design-knight42" title="Design">🎨</a> <a href="#ideas-knight42" title="Ideas, Planning, & Feedback">🤔</a> <a href="#research-knight42" title="Research">🔬</a></td>
      <td align="center" valign="top" width="14.28%"><a href="https://github.com/VoVAllen"><img src="https://avatars.githubusercontent.com/u/8686776?v=4?s=70" width="70px;" alt="Jinjing Zhou"/><br /><sub><b>Jinjing Zhou</b></sub></a><br /><a href="https://github.com/tensorchord/envd/issues?q=author%3AVoVAllen" title="Bug reports">🐛</a> <a href="https://github.com/tensorchord/envd/commits?author=VoVAllen" title="Code">💻</a> <a href="#design-VoVAllen" title="Design">🎨</a> <a href="https://github.com/tensorchord/envd/commits?author=VoVAllen" title="Documentation">📖</a></td>
      <td align="center" valign="top" width="14.28%"><a href="http://jun.dev/blog/issues"><img src="https://avatars.githubusercontent.com/u/8097526?v=4?s=70" width="70px;" alt="Jun"/><br /><sub><b>Jun</b></sub></a><br /><a href="#platform-junnplus" title="Packaging/porting to new platform">📦</a> <a href="https://github.com/tensorchord/envd/commits?author=junnplus" title="Code">💻</a></td>
      <td align="center" valign="top" width="14.28%"><a href="https://github.com/Kaiyang-Chen"><img src="https://avatars.githubusercontent.com/u/48289729?v=4?s=70" width="70px;" alt="Kaiyang Chen"/><br /><sub><b>Kaiyang Chen</b></sub></a><br /><a href="https://github.com/tensorchord/envd/commits?author=Kaiyang-Chen" title="Code">💻</a></td>
      <td align="center" valign="top" width="14.28%"><a href="https://kemingy.github.io/"><img src="https://avatars.githubusercontent.com/u/12974685?v=4?s=70" width="70px;" alt="Keming"/><br /><sub><b>Keming</b></sub></a><br /><a href="https://github.com/tensorchord/envd/commits?author=kemingy" title="Code">💻</a> <a href="https://github.com/tensorchord/envd/commits?author=kemingy" title="Documentation">📖</a> <a href="#ideas-kemingy" title="Ideas, Planning, & Feedback">🤔</a> <a href="#infra-kemingy" title="Infrastructure (Hosting, Build-Tools, etc)">🚇</a></td>
      <td align="center" valign="top" width="14.28%"><a href="https://github.com/pingsutw"><img src="https://avatars.githubusercontent.com/u/37936015?v=4?s=70" width="70px;" alt="Kevin Su"/><br /><sub><b>Kevin Su</b></sub></a><br /><a href="https://github.com/tensorchord/envd/commits?author=pingsutw" title="Code">💻</a></td>
      <td align="center" valign="top" width="14.28%"><a href="https://github.com/3AceShowHand"><img src="https://avatars.githubusercontent.com/u/7138436?v=4?s=70" width="70px;" alt="Ling Jin"/><br /><sub><b>Ling Jin</b></sub></a><br /><a href="https://github.com/tensorchord/envd/issues?q=author%3A3AceShowHand" title="Bug reports">🐛</a> <a href="#infra-3AceShowHand" title="Infrastructure (Hosting, Build-Tools, etc)">🚇</a></td>
    </tr>
    <tr>
      <td align="center" valign="top" width="14.28%"><a href="http://manjusaka.itscoder.com"><img src="https://avatars.githubusercontent.com/u/7054676?v=4?s=70" width="70px;" alt="Manjusaka"/><br /><sub><b>Manjusaka</b></sub></a><br /><a href="https://github.com/tensorchord/envd/commits?author=Zheaoli" title="Code">💻</a></td>
      <td align="center" valign="top" width="14.28%"><a href="https://github.com/lilylee1874"><img src="https://avatars.githubusercontent.com/u/52693877?v=4?s=70" width="70px;" alt="Nino"/><br /><sub><b>Nino</b></sub></a><br /><a href="#design-lilylee1874" title="Design">🎨</a> <a href="https://github.com/tensorchord/envd/commits?author=lilylee1874" title="Code">💻</a></td>
      <td align="center" valign="top" width="14.28%"><a href="http://phillipw.info"><img src="https://avatars.githubusercontent.com/u/34707116?v=4?s=70" width="70px;" alt="Pengyu Wang"/><br /><sub><b>Pengyu Wang</b></sub></a><br /><a href="https://github.com/tensorchord/envd/commits?author=cswpy" title="Documentation">📖</a></td>
      <td align="center" valign="top" width="14.28%"><a href="https://github.com/Sepush"><img src="https://avatars.githubusercontent.com/u/39197136?v=4?s=70" width="70px;" alt="Sepush"/><br /><sub><b>Sepush</b></sub></a><br /><a href="https://github.com/tensorchord/envd/commits?author=sepush" title="Documentation">📖</a></td>
      <td align="center" valign="top" width="14.28%"><a href="https://blog.thrimbda.com/"><img src="https://avatars.githubusercontent.com/u/15231162?v=4?s=70" width="70px;" alt="Siyuan Wang"/><br /><sub><b>Siyuan Wang</b></sub></a><br /><a href="https://github.com/tensorchord/envd/commits?author=Thrimbda" title="Code">💻</a> <a href="#infra-Thrimbda" title="Infrastructure (Hosting, Build-Tools, etc)">🚇</a> <a href="#maintenance-Thrimbda" title="Maintenance">🚧</a></td>
      <td align="center" valign="top" width="14.28%"><a href="http://suyan.moe"><img src="https://avatars.githubusercontent.com/u/24221472?v=4?s=70" width="70px;" alt="Suyan"/><br /><sub><b>Suyan</b></sub></a><br /><a href="https://github.com/tensorchord/envd/commits?author=suyanhanx" title="Documentation">📖</a></td>
      <td align="center" valign="top" width="14.28%"><a href="http://myfra.vercel.app"><img src="https://avatars.githubusercontent.com/u/60420319?v=4?s=70" width="70px;" alt="To My"/><br /><sub><b>To My</b></sub></a><br /><a href="https://github.com/tensorchord/envd/commits?author=MyFRA" title="Documentation">📖</a></td>
    </tr>
    <tr>
      <td align="center" valign="top" width="14.28%"><a href="https://www.iam.rw"><img src="https://avatars.githubusercontent.com/u/29533246?v=4?s=70" width="70px;" alt="Tumushimire Yves"/><br /><sub><b>Tumushimire Yves</b></sub></a><br /><a href="https://github.com/tensorchord/envd/commits?author=yvestumushimire" title="Code">💻</a></td>
      <td align="center" valign="top" width="14.28%"><a href="https://page.codespaper.com"><img src="https://avatars.githubusercontent.com/u/3764335?v=4?s=70" width="70px;" alt="Wei Zhang"/><br /><sub><b>Wei Zhang</b></sub></a><br /><a href="https://github.com/tensorchord/envd/commits?author=zwpaper" title="Code">💻</a></td>
      <td align="center" valign="top" width="14.28%"><a href="http://weixiao-huang.github.io"><img src="https://avatars.githubusercontent.com/u/8834733?v=4?s=70" width="70px;" alt="Weixiao Huang"/><br /><sub><b>Weixiao Huang</b></sub></a><br /><a href="https://github.com/tensorchord/envd/commits?author=weixiao-huang" title="Code">💻</a></td>
      <td align="center" valign="top" width="14.28%"><a href="https://www.hawkingrei.com/"><img src="https://avatars.githubusercontent.com/u/3427324?v=4?s=70" width="70px;" alt="Weizhen Wang"/><br /><sub><b>Weizhen Wang</b></sub></a><br /><a href="https://github.com/tensorchord/envd/commits?author=hawkingrei" title="Code">💻</a></td>
      <td align="center" valign="top" width="14.28%"><a href="http://blog.xuruowei.com"><img src="https://avatars.githubusercontent.com/u/18398013?v=4?s=70" width="70px;" alt="XRW"/><br /><sub><b>XRW</b></sub></a><br /><a href="https://github.com/tensorchord/envd/commits?author=Xuruowei" title="Code">💻</a></td>
      <td align="center" valign="top" width="14.28%"><a href="https://github.com/jiayouxujin"><img src="https://avatars.githubusercontent.com/u/29749249?v=4?s=70" width="70px;" alt="Xu Jin"/><br /><sub><b>Xu Jin</b></sub></a><br /><a href="https://github.com/tensorchord/envd/commits?author=jiayouxujin" title="Code">💻</a></td>
      <td align="center" valign="top" width="14.28%"><a href="https://xuanwo.io/"><img src="https://avatars.githubusercontent.com/u/5351546?v=4?s=70" width="70px;" alt="Xuanwo"/><br /><sub><b>Xuanwo</b></sub></a><br /><a href="#question-Xuanwo" title="Answering Questions">💬</a> <a href="#design-Xuanwo" title="Design">🎨</a> <a href="#ideas-Xuanwo" title="Ideas, Planning, & Feedback">🤔</a> <a href="https://github.com/tensorchord/envd/pulls?q=is%3Apr+reviewed-by%3AXuanwo" title="Reviewed Pull Requests">👀</a></td>
    </tr>
    <tr>
      <td align="center" valign="top" width="14.28%"><a href="https://github.com/enjoyliu"><img src="https://avatars.githubusercontent.com/u/13460894?v=4?s=70" width="70px;" alt="Yijiang Liu"/><br /><sub><b>Yijiang Liu</b></sub></a><br /><a href="https://github.com/tensorchord/envd/commits?author=enjoyliu" title="Code">💻</a></td>
      <td align="center" valign="top" width="14.28%"><a href="https://elon.fun/"><img src="https://avatars.githubusercontent.com/u/8433465?v=4?s=70" width="70px;" alt="Yilong Li"/><br /><sub><b>Yilong Li</b></sub></a><br /><a href="https://github.com/tensorchord/envd/commits?author=dragonly" title="Documentation">📖</a> <a href="https://github.com/tensorchord/envd/issues?q=author%3Adragonly" title="Bug reports">🐛</a> <a href="https://github.com/tensorchord/envd/commits?author=dragonly" title="Code">💻</a></td>
      <td align="center" valign="top" width="14.28%"><a href="https://terrytangyuan.github.io/about/"><img src="https://avatars.githubusercontent.com/u/4269898?v=4?s=70" width="70px;" alt="Yuan Tang"/><br /><sub><b>Yuan Tang</b></sub></a><br /><a href="https://github.com/tensorchord/envd/commits?author=terrytangyuan" title="Code">💻</a> <a href="#design-terrytangyuan" title="Design">🎨</a> <a href="https://github.com/tensorchord/envd/commits?author=terrytangyuan" title="Documentation">📖</a> <a href="#ideas-terrytangyuan" title="Ideas, Planning, & Feedback">🤔</a></td>
      <td align="center" valign="top" width="14.28%"><a href="https://rudeigerc.dev/"><img src="https://avatars.githubusercontent.com/u/18243819?v=4?s=70" width="70px;" alt="Yuchen Cheng"/><br /><sub><b>Yuchen Cheng</b></sub></a><br /><a href="https://github.com/tensorchord/envd/issues?q=author%3Arudeigerc" title="Bug reports">🐛</a> <a href="#infra-rudeigerc" title="Infrastructure (Hosting, Build-Tools, etc)">🚇</a> <a href="#maintenance-rudeigerc" title="Maintenance">🚧</a> <a href="#tool-rudeigerc" title="Tools">🔧</a></td>
      <td align="center" valign="top" width="14.28%"><a href="http://lunarwhite.github.io"><img src="https://avatars.githubusercontent.com/u/57584831?v=4?s=70" width="70px;" alt="Yuedong Wu"/><br /><sub><b>Yuedong Wu</b></sub></a><br /><a href="https://github.com/tensorchord/envd/commits?author=lunarwhite" title="Code">💻</a></td>
      <td align="center" valign="top" width="14.28%"><a href="https://github.com/yczheng0"><img src="https://avatars.githubusercontent.com/u/21327543?v=4?s=70" width="70px;" alt="Yunchuan Zheng"/><br /><sub><b>Yunchuan Zheng</b></sub></a><br /><a href="https://github.com/tensorchord/envd/commits?author=yczheng0" title="Code">💻</a></td>
      <td align="center" valign="top" width="14.28%"><a href="http://lizheming.top"><img src="https://avatars.githubusercontent.com/u/9639449?v=4?s=70" width="70px;" alt="Zheming Li"/><br /><sub><b>Zheming Li</b></sub></a><br /><a href="https://github.com/tensorchord/envd/commits?author=lizhemingi" title="Code">💻</a></td>
    </tr>
    <tr>
      <td align="center" valign="top" width="14.28%"><a href="https://github.com/Xiaoaier-Z-L"><img src="https://avatars.githubusercontent.com/u/96805673?v=4?s=70" width="70px;" alt="Zhenguo.Li"/><br /><sub><b>Zhenguo.Li</b></sub></a><br /><a href="https://github.com/tensorchord/envd/commits?author=Xiaoaier-Z-L" title="Code">💻</a> <a href="https://github.com/tensorchord/envd/commits?author=Xiaoaier-Z-L" title="Documentation">📖</a></td>
      <td align="center" valign="top" width="14.28%"><a href="https://blog.triplez.cn/"><img src="https://avatars.githubusercontent.com/u/16285716?v=4?s=70" width="70px;" alt="Zhenzhen Zhao"/><br /><sub><b>Zhenzhen Zhao</b></sub></a><br /><a href="#infra-Triple-Z" title="Infrastructure (Hosting, Build-Tools, etc)">🚇</a> <a href="#userTesting-Triple-Z" title="User Testing">📓</a> <a href="https://github.com/tensorchord/envd/commits?author=Triple-Z" title="Code">💻</a></td>
      <td align="center" valign="top" width="14.28%"><a href="https://t.me/littlepoint"><img src="https://avatars.githubusercontent.com/u/7611700?v=4?s=70" width="70px;" alt="Zhizhen He"/><br /><sub><b>Zhizhen He</b></sub></a><br /><a href="https://github.com/tensorchord/envd/commits?author=hezhizhen" title="Code">💻</a> <a href="https://github.com/tensorchord/envd/commits?author=hezhizhen" title="Documentation">📖</a></td>
      <td align="center" valign="top" width="14.28%"><a href="https://github.com/cutecutecat"><img src="https://avatars.githubusercontent.com/u/19801166?v=4?s=70" width="70px;" alt="cutecutecat"/><br /><sub><b>cutecutecat</b></sub></a><br /><a href="https://github.com/tensorchord/envd/commits?author=cutecutecat" title="Code">💻</a></td>
      <td align="center" valign="top" width="14.28%"><a href="https://github.com/dqhl76"><img src="https://avatars.githubusercontent.com/u/69136320?v=4?s=70" width="70px;" alt="dqhl76"/><br /><sub><b>dqhl76</b></sub></a><br /><a href="https://github.com/tensorchord/envd/commits?author=dqhl76" title="Documentation">📖</a> <a href="https://github.com/tensorchord/envd/commits?author=dqhl76" title="Code">💻</a></td>
      <td align="center" valign="top" width="14.28%"><a href="https://github.com/jimoosciuc"><img src="https://avatars.githubusercontent.com/u/33337387?v=4?s=70" width="70px;" alt="jimoosciuc"/><br /><sub><b>jimoosciuc</b></sub></a><br /><a href="#userTesting-jimoosciuc" title="User Testing">📓</a></td>
      <td align="center" valign="top" width="14.28%"><a href="https://kenwoodjw.github.io"><img src="https://avatars.githubusercontent.com/u/10386710?v=4?s=70" width="70px;" alt="kenwoodjw"/><br /><sub><b>kenwoodjw</b></sub></a><br /><a href="https://github.com/tensorchord/envd/commits?author=kenwoodjw" title="Code">💻</a></td>
    </tr>
    <tr>
      <td align="center" valign="top" width="14.28%"><a href="http://www.hwdef.org"><img src="https://avatars.githubusercontent.com/u/13084946?v=4?s=70" width="70px;" alt="li mengyang"/><br /><sub><b>li mengyang</b></sub></a><br /><a href="https://github.com/tensorchord/envd/commits?author=hwdef" title="Code">💻</a></td>
      <td align="center" valign="top" width="14.28%"><a href="https://github.com/aseaday"><img src="https://avatars.githubusercontent.com/u/3927355?v=4?s=70" width="70px;" alt="nullday"/><br /><sub><b>nullday</b></sub></a><br /><a href="#ideas-aseaday" title="Ideas, Planning, & Feedback">🤔</a> <a href="https://github.com/tensorchord/envd/commits?author=aseaday" title="Code">💻</a></td>
      <td align="center" valign="top" width="14.28%"><a href="https://github.com/rrain7"><img src="https://avatars.githubusercontent.com/u/49144127?v=4?s=70" width="70px;" alt="rrain7"/><br /><sub><b>rrain7</b></sub></a><br /><a href="https://github.com/tensorchord/envd/commits?author=rrain7" title="Code">💻</a></td>
      <td align="center" valign="top" width="14.28%"><a href="https://tisonkun.org/"><img src="https://avatars.githubusercontent.com/u/18818196?v=4?s=70" width="70px;" alt="tison"/><br /><sub><b>tison</b></sub></a><br /><a href="https://github.com/tensorchord/envd/commits?author=tisonkun" title="Code">💻</a></td>
      <td align="center" valign="top" width="14.28%"><a href="http://fatelei.github.io"><img src="https://avatars.githubusercontent.com/u/961094?v=4?s=70" width="70px;" alt="wangxiaolei"/><br /><sub><b>wangxiaolei</b></sub></a><br /><a href="https://github.com/tensorchord/envd/commits?author=fatelei" title="Code">💻</a></td>
      <td align="center" valign="top" width="14.28%"><a href="https://github.com/sea-wyq"><img src="https://avatars.githubusercontent.com/u/22475606?v=4?s=70" width="70px;" alt="wyq"/><br /><sub><b>wyq</b></sub></a><br /><a href="https://github.com/tensorchord/envd/issues?q=author%3Asea-wyq" title="Bug reports">🐛</a> <a href="#design-sea-wyq" title="Design">🎨</a> <a href="https://github.com/tensorchord/envd/commits?author=sea-wyq" title="Code">💻</a></td>
      <td align="center" valign="top" width="14.28%"><a href="https://oubotong.github.io/johan"><img src="https://avatars.githubusercontent.com/u/26356127?v=4?s=70" width="70px;" alt="x0oo0x"/><br /><sub><b>x0oo0x</b></sub></a><br /><a href="https://github.com/tensorchord/envd/commits?author=oubotong" title="Code">💻</a></td>
    </tr>
    <tr>
      <td align="center" valign="top" width="14.28%"><a href="https://github.com/xiangtianyu"><img src="https://avatars.githubusercontent.com/u/10825900?v=4?s=70" width="70px;" alt="xiangtianyu"/><br /><sub><b>xiangtianyu</b></sub></a><br /><a href="https://github.com/tensorchord/envd/commits?author=xiangtianyu" title="Documentation">📖</a></td>
      <td align="center" valign="top" width="14.28%"><a href="https://github.com/xieydd"><img src="https://avatars.githubusercontent.com/u/20329697?v=4?s=70" width="70px;" alt="xieydd"/><br /><sub><b>xieydd</b></sub></a><br /><a href="https://github.com/tensorchord/envd/commits?author=xieydd" title="Code">💻</a></td>
      <td align="center" valign="top" width="14.28%"><a href="https://github.com/xing0821"><img src="https://avatars.githubusercontent.com/u/54933318?v=4?s=70" width="70px;" alt="xing0821"/><br /><sub><b>xing0821</b></sub></a><br /><a href="#ideas-xing0821" title="Ideas, Planning, & Feedback">🤔</a> <a href="#userTesting-xing0821" title="User Testing">📓</a> <a href="https://github.com/tensorchord/envd/commits?author=xing0821" title="Code">💻</a></td>
      <td align="center" valign="top" width="14.28%"><a href="https://xxchan.github.io"><img src="https://avatars.githubusercontent.com/u/37948597?v=4?s=70" width="70px;" alt="xxchan"/><br /><sub><b>xxchan</b></sub></a><br /><a href="https://github.com/tensorchord/envd/commits?author=xxchan" title="Documentation">📖</a></td>
      <td align="center" valign="top" width="14.28%"><a href="https://github.com/zhyon404"><img src="https://avatars.githubusercontent.com/u/32242529?v=4?s=70" width="70px;" alt="zhyon404"/><br /><sub><b>zhyon404</b></sub></a><br /><a href="https://github.com/tensorchord/envd/commits?author=zhyon404" title="Code">💻</a></td>
      <td align="center" valign="top" width="14.28%"><a href="https://www.homeboyc.cn/"><img src="https://avatars.githubusercontent.com/u/22193008?v=4?s=70" width="70px;" alt="杨成锴"/><br /><sub><b>杨成锴</b></sub></a><br /><a href="https://github.com/tensorchord/envd/commits?author=asjdf" title="Code">💻</a></td>
    </tr>
  </tbody>
</table>

<!-- markdownlint-restore -->
<!-- prettier-ignore-end -->

<!-- ALL-CONTRIBUTORS-LIST:END -->

This project follows the [all-contributors](https://github.com/all-contributors/all-contributors) specification. Contributions of any kind welcome!

## License 📋

[Apache 2.0](./LICENSE)

<a href="https://trackgit.com"><img src="https://us-central1-trackgit-analytics.cloudfunctions.net/token/ping/l3ldvdaswvnjpty9u7l3" alt="trackgit-views" /></a>
