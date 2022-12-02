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

<p><img src="https://user-images.githubusercontent.com/52693877/193411608-01864ed9-84e2-4ea8-b978-8cfe09d394ca.png#gh-light-mode-only" alt="explain"/></p>
<p><img src="https://user-images.githubusercontent.com/16186646/200111649-771d3214-dd23-49fa-a3d8-761a88f805bf.svg#gh-dark-mode-only" alt="explain"/></p>

## What is envd?

envd (`ɪnˈvdɪ`) is a command-line tool that helps you create the container-based development environment for AI/ML.

Development environments are full of python and system dependencies, CUDA, BASH scripts, Dockerfiles, SSH configurations, Kubernetes YAMLs, and many other clunky things that are always breaking. envd is to solve the problem:

1. Declare the list of dependencies (CUDA, python packages, your favorite IDE, and so on) in `build.envd`
1. Simply run `envd up`.
1. Develop in the isolated environment.

<p align="center">
  <img src="https://user-images.githubusercontent.com/5100735/189058399-3865a039-9459-4e74-83dd-3ee2ecadfef5.svg" width="75%"/>
</p>

## Why use `envd`?

Environments built with `envd` provide the following features out-of-the-box:

❤️ **Knowledge reuse in your team**

`envd` build functions can be reused. Use `include` function to import any git repositories. No more copy/paste Dockerfile instructions, let's reuse them.

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

⏱️ **Builtkit native, build up to 6x faster**

[Buildkit](https://github.com/moby/buildkit) supports parallel builds and software cache (e.g. pip index cache and apt cache). You can enjoy the benefits without knowledge of it.

For example, the PyPI cache is shared across builds and thus the package will be cached if it has been downloaded before.

<p align=center>
  <img src="https://user-images.githubusercontent.com/5100735/189928628-543f4851-87b7-462b-b811-372cbf46ff25.svg#gh-light-mode-only" width="65%"/>
</p>

<p align=center>
  <img src="https://user-images.githubusercontent.com/16186646/197944452-4a5dcd5f-68d0-4505-b419-e95c298978d7.svg#gh-dark-mode-only" width="65%"/>
</p>

🐍 **One configuration to rule them all**

Development environments are full of Dockerfiles, bash scripts, Kubernetes YAML manifests, and many other clunky files that are always breaking. You just need one configuration file `build.envd`[^1], it works both for local Docker and Kubernetes clusters in the cloud.

![envd](https://user-images.githubusercontent.com/5100735/188821980-dcbd9069-b504-436a-9ffd-05ac5543a6d1.png)

[^1]: The build language is [starlark](https://docs.bazel.build/versions/main/skylark/language.html), which is a dialect of Python.

✍️ **Don't sacrifice your developer experience**

SSH is configured for the created environment. You can use vscode-remote, jupyter, pycharm or other IDEs that you love. Besides this, declare the IDE extensions you want, let `envd` take care of them.

```python
def build():
    install.vscode_extensions([
        "ms-python.python",
    ])
```

☁️ **No polluted environment**

Are you working on multiple projects, all of which need different versions of CUDA? `envd` helps you create isolated and clean environments.

## Who should use envd?

We're focused on helping data scientists and teams that develop AI/ML models. And they may suffer from:

- building the development environments with Python/R/Julia, CUDA, Docker, SSH, and so on. Do you have a complicated Dockerfile or build script that sets up all your dev environments, but is always breaking?
- Updating the environment. Do you always need to ask infrastructure engineers how to add a new Python/R/Julia package in the Dockerfile?
- Managing environments and machines. Do you always forget which machines are used for the specific project, because you handle multiple projects concurrently?

---

**Talk with us**

💬 Interested in talking with us about your experience building or managing AI/ML applications?

[**Set up a time to chat!**](https://forms.gle/9HDBHX5Y3fzuDCDAA)

## Getting Started 🚀

### Requirements

- Docker (20.10.0 or above)

### Install and bootstrap `envd`

`envd` can be installed with `pip` (only support Python3). After the installation, please run `envd bootstrap` to bootstrap.

```bash
pip3 install --pre --upgrade envd
envd bootstrap
```

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
- To build from the source, please read our [contributing documentation](https://envd.tensorchord.ai/community/contributing.html) and [development tutorial](https://envd.tensorchord.ai/community/development.html).

[![Open in Gitpod](https://gitpod.io/button/open-in-gitpod.svg)](https://gitpod.io/#https://github.com/tensorchord/envd)

## Contributors ✨

Thanks goes to these wonderful people ([emoji key](https://allcontributors.org/docs/en/emoji-key)):

<!-- ALL-CONTRIBUTORS-LIST:START - Do not remove or modify this section -->
<!-- prettier-ignore-start -->
<!-- markdownlint-disable -->
<table>
  <tbody>
    <tr>
      <td align="center"><a href="http://blog.duanfei.org"><img src="https://avatars.githubusercontent.com/u/16186646?v=4?s=70" width="70px;" alt=" Friends A."/><br /><sub><b> Friends A.</b></sub></a><br /><a href="https://github.com/tensorchord/envd/commits?author=shaonianche" title="Documentation">📖</a> <a href="#design-shaonianche" title="Design">🎨</a></td>
      <td align="center"><a href="https://github.com/aaronzs"><img src="https://avatars.githubusercontent.com/u/1827365?v=4?s=70" width="70px;" alt="Aaron Sun"/><br /><sub><b>Aaron Sun</b></sub></a><br /><a href="#userTesting-aaronzs" title="User Testing">📓</a> <a href="https://github.com/tensorchord/envd/commits?author=aaronzs" title="Code">💻</a></td>
      <td align="center"><a href="https://github.com/popfido"><img src="https://avatars.githubusercontent.com/u/3928409?v=4?s=70" width="70px;" alt="Aka.Fido"/><br /><sub><b>Aka.Fido</b></sub></a><br /><a href="#platform-popfido" title="Packaging/porting to new platform">📦</a> <a href="https://github.com/tensorchord/envd/commits?author=popfido" title="Documentation">📖</a> <a href="https://github.com/tensorchord/envd/commits?author=popfido" title="Code">💻</a></td>
      <td align="center"><a href="http://alexhxi.com"><img src="https://avatars.githubusercontent.com/u/68758451?v=4?s=70" width="70px;" alt="Alex Xi"/><br /><sub><b>Alex Xi</b></sub></a><br /><a href="https://github.com/tensorchord/envd/commits?author=AlexXi19" title="Code">💻</a></td>
      <td align="center"><a href="https://github.com/sunby"><img src="https://avatars.githubusercontent.com/u/9817127?v=4?s=70" width="70px;" alt="Bingyi Sun"/><br /><sub><b>Bingyi Sun</b></sub></a><br /><a href="https://github.com/tensorchord/envd/commits?author=sunby" title="Code">💻</a></td>
      <td align="center"><a href="http://gaocegege.com/Blog"><img src="https://avatars.githubusercontent.com/u/5100735?v=4?s=70" width="70px;" alt="Ce Gao"/><br /><sub><b>Ce Gao</b></sub></a><br /><a href="https://github.com/tensorchord/envd/commits?author=gaocegege" title="Code">💻</a> <a href="https://github.com/tensorchord/envd/commits?author=gaocegege" title="Documentation">📖</a> <a href="#design-gaocegege" title="Design">🎨</a> <a href="#projectManagement-gaocegege" title="Project Management">📆</a></td>
      <td align="center"><a href="https://frostming.com"><img src="https://avatars.githubusercontent.com/u/16336606?v=4?s=70" width="70px;" alt="Frost Ming"/><br /><sub><b>Frost Ming</b></sub></a><br /><a href="https://github.com/tensorchord/envd/commits?author=frostming" title="Code">💻</a></td>
    </tr>
    <tr>
      <td align="center"><a href="https://GuangyangLi.com"><img src="https://avatars.githubusercontent.com/u/2060045?v=4?s=70" width="70px;" alt="Guangyang Li"/><br /><sub><b>Guangyang Li</b></sub></a><br /><a href="https://github.com/tensorchord/envd/commits?author=gyli" title="Code">💻</a></td>
      <td align="center"><a href="https://github.com/Gui-Yue"><img src="https://avatars.githubusercontent.com/u/78520005?v=4?s=70" width="70px;" alt="Gui-Yue"/><br /><sub><b>Gui-Yue</b></sub></a><br /><a href="https://github.com/tensorchord/envd/commits?author=Gui-Yue" title="Code">💻</a></td>
      <td align="center"><a href="https://github.com/haiker2011"><img src="https://avatars.githubusercontent.com/u/8073429?v=4?s=70" width="70px;" alt="Haiker Sun"/><br /><sub><b>Haiker Sun</b></sub></a><br /><a href="https://github.com/tensorchord/envd/commits?author=haiker2011" title="Code">💻</a></td>
      <td align="center"><a href="https://bandism.net/"><img src="https://avatars.githubusercontent.com/u/22633385?v=4?s=70" width="70px;" alt="Ikko Ashimine"/><br /><sub><b>Ikko Ashimine</b></sub></a><br /><a href="https://github.com/tensorchord/envd/commits?author=eltociear" title="Code">💻</a></td>
      <td align="center"><a href="https://github.com/nasnoisaac"><img src="https://avatars.githubusercontent.com/u/11145462?v=4?s=70" width="70px;" alt="Isaac "/><br /><sub><b>Isaac </b></sub></a><br /><a href="https://github.com/tensorchord/envd/commits?author=nasnoisaac" title="Code">💻</a></td>
      <td align="center"><a href="https://belyenochi.github.io/"><img src="https://avatars.githubusercontent.com/u/26409132?v=4?s=70" width="70px;" alt="JasonZhu"/><br /><sub><b>JasonZhu</b></sub></a><br /><a href="https://github.com/tensorchord/envd/commits?author=Belyenochi" title="Code">💻</a></td>
      <td align="center"><a href="https://github.com/knight42"><img src="https://avatars.githubusercontent.com/u/4237254?v=4?s=70" width="70px;" alt="Jian Zeng"/><br /><sub><b>Jian Zeng</b></sub></a><br /><a href="#design-knight42" title="Design">🎨</a> <a href="#ideas-knight42" title="Ideas, Planning, & Feedback">🤔</a> <a href="#research-knight42" title="Research">🔬</a></td>
    </tr>
    <tr>
      <td align="center"><a href="https://github.com/VoVAllen"><img src="https://avatars.githubusercontent.com/u/8686776?v=4?s=70" width="70px;" alt="Jinjing Zhou"/><br /><sub><b>Jinjing Zhou</b></sub></a><br /><a href="https://github.com/tensorchord/envd/issues?q=author%3AVoVAllen" title="Bug reports">🐛</a> <a href="https://github.com/tensorchord/envd/commits?author=VoVAllen" title="Code">💻</a> <a href="#design-VoVAllen" title="Design">🎨</a> <a href="https://github.com/tensorchord/envd/commits?author=VoVAllen" title="Documentation">📖</a></td>
      <td align="center"><a href="http://jun.dev/blog/issues"><img src="https://avatars.githubusercontent.com/u/8097526?v=4?s=70" width="70px;" alt="Jun"/><br /><sub><b>Jun</b></sub></a><br /><a href="#platform-junnplus" title="Packaging/porting to new platform">📦</a> <a href="https://github.com/tensorchord/envd/commits?author=junnplus" title="Code">💻</a></td>
      <td align="center"><a href="https://kemingy.github.io/"><img src="https://avatars.githubusercontent.com/u/12974685?v=4?s=70" width="70px;" alt="Keming"/><br /><sub><b>Keming</b></sub></a><br /><a href="https://github.com/tensorchord/envd/commits?author=kemingy" title="Code">💻</a> <a href="https://github.com/tensorchord/envd/commits?author=kemingy" title="Documentation">📖</a> <a href="#ideas-kemingy" title="Ideas, Planning, & Feedback">🤔</a> <a href="#infra-kemingy" title="Infrastructure (Hosting, Build-Tools, etc)">🚇</a></td>
      <td align="center"><a href="https://github.com/pingsutw"><img src="https://avatars.githubusercontent.com/u/37936015?v=4?s=70" width="70px;" alt="Kevin Su"/><br /><sub><b>Kevin Su</b></sub></a><br /><a href="https://github.com/tensorchord/envd/commits?author=pingsutw" title="Code">💻</a></td>
      <td align="center"><a href="https://github.com/3AceShowHand"><img src="https://avatars.githubusercontent.com/u/7138436?v=4?s=70" width="70px;" alt="Ling Jin"/><br /><sub><b>Ling Jin</b></sub></a><br /><a href="https://github.com/tensorchord/envd/issues?q=author%3A3AceShowHand" title="Bug reports">🐛</a> <a href="#infra-3AceShowHand" title="Infrastructure (Hosting, Build-Tools, etc)">🚇</a></td>
      <td align="center"><a href="http://manjusaka.itscoder.com"><img src="https://avatars.githubusercontent.com/u/7054676?v=4?s=70" width="70px;" alt="Manjusaka"/><br /><sub><b>Manjusaka</b></sub></a><br /><a href="https://github.com/tensorchord/envd/commits?author=Zheaoli" title="Code">💻</a></td>
      <td align="center"><a href="https://github.com/lilylee1874"><img src="https://avatars.githubusercontent.com/u/52693877?v=4?s=70" width="70px;" alt="Nino"/><br /><sub><b>Nino</b></sub></a><br /><a href="#design-lilylee1874" title="Design">🎨</a> <a href="https://github.com/tensorchord/envd/commits?author=lilylee1874" title="Code">💻</a></td>
    </tr>
    <tr>
      <td align="center"><a href="http://phillipw.info"><img src="https://avatars.githubusercontent.com/u/34707116?v=4?s=70" width="70px;" alt="Pengyu Wang"/><br /><sub><b>Pengyu Wang</b></sub></a><br /><a href="https://github.com/tensorchord/envd/commits?author=cswpy" title="Documentation">📖</a></td>
      <td align="center"><a href="https://github.com/Sepush"><img src="https://avatars.githubusercontent.com/u/39197136?v=4?s=70" width="70px;" alt="Sepush"/><br /><sub><b>Sepush</b></sub></a><br /><a href="https://github.com/tensorchord/envd/commits?author=sepush" title="Documentation">📖</a></td>
      <td align="center"><a href="https://blog.thrimbda.com/"><img src="https://avatars.githubusercontent.com/u/15231162?v=4?s=70" width="70px;" alt="Siyuan Wang"/><br /><sub><b>Siyuan Wang</b></sub></a><br /><a href="https://github.com/tensorchord/envd/commits?author=Thrimbda" title="Code">💻</a> <a href="#infra-Thrimbda" title="Infrastructure (Hosting, Build-Tools, etc)">🚇</a> <a href="#maintenance-Thrimbda" title="Maintenance">🚧</a></td>
      <td align="center"><a href="http://suyan.moe"><img src="https://avatars.githubusercontent.com/u/24221472?v=4?s=70" width="70px;" alt="Suyan"/><br /><sub><b>Suyan</b></sub></a><br /><a href="https://github.com/tensorchord/envd/commits?author=suyanhanx" title="Documentation">📖</a></td>
      <td align="center"><a href="http://myfra.vercel.app"><img src="https://avatars.githubusercontent.com/u/60420319?v=4?s=70" width="70px;" alt="To My"/><br /><sub><b>To My</b></sub></a><br /><a href="https://github.com/tensorchord/envd/commits?author=MyFRA" title="Documentation">📖</a></td>
      <td align="center"><a href="https://www.iam.rw"><img src="https://avatars.githubusercontent.com/u/29533246?v=4?s=70" width="70px;" alt="Tumushimire Yves"/><br /><sub><b>Tumushimire Yves</b></sub></a><br /><a href="https://github.com/tensorchord/envd/commits?author=yvestumushimire" title="Code">💻</a></td>
      <td align="center"><a href="https://page.codespaper.com"><img src="https://avatars.githubusercontent.com/u/3764335?v=4?s=70" width="70px;" alt="Wei Zhang"/><br /><sub><b>Wei Zhang</b></sub></a><br /><a href="https://github.com/tensorchord/envd/commits?author=zwpaper" title="Code">💻</a></td>
    </tr>
    <tr>
      <td align="center"><a href="https://www.hawkingrei.com/"><img src="https://avatars.githubusercontent.com/u/3427324?v=4?s=70" width="70px;" alt="Weizhen Wang"/><br /><sub><b>Weizhen Wang</b></sub></a><br /><a href="https://github.com/tensorchord/envd/commits?author=hawkingrei" title="Code">💻</a></td>
      <td align="center"><a href="http://blog.xuruowei.com"><img src="https://avatars.githubusercontent.com/u/18398013?v=4?s=70" width="70px;" alt="XRW"/><br /><sub><b>XRW</b></sub></a><br /><a href="https://github.com/tensorchord/envd/commits?author=Xuruowei" title="Code">💻</a></td>
      <td align="center"><a href="https://github.com/jiayouxujin"><img src="https://avatars.githubusercontent.com/u/29749249?v=4?s=70" width="70px;" alt="Xu Jin"/><br /><sub><b>Xu Jin</b></sub></a><br /><a href="https://github.com/tensorchord/envd/commits?author=jiayouxujin" title="Code">💻</a></td>
      <td align="center"><a href="https://xuanwo.io/"><img src="https://avatars.githubusercontent.com/u/5351546?v=4?s=70" width="70px;" alt="Xuanwo"/><br /><sub><b>Xuanwo</b></sub></a><br /><a href="#question-Xuanwo" title="Answering Questions">💬</a> <a href="#design-Xuanwo" title="Design">🎨</a> <a href="#ideas-Xuanwo" title="Ideas, Planning, & Feedback">🤔</a> <a href="https://github.com/tensorchord/envd/pulls?q=is%3Apr+reviewed-by%3AXuanwo" title="Reviewed Pull Requests">👀</a></td>
      <td align="center"><a href="https://github.com/enjoyliu"><img src="https://avatars.githubusercontent.com/u/13460894?v=4?s=70" width="70px;" alt="Yijiang Liu"/><br /><sub><b>Yijiang Liu</b></sub></a><br /><a href="https://github.com/tensorchord/envd/commits?author=enjoyliu" title="Code">💻</a></td>
      <td align="center"><a href="https://elon.fun/"><img src="https://avatars.githubusercontent.com/u/8433465?v=4?s=70" width="70px;" alt="Yilong Li"/><br /><sub><b>Yilong Li</b></sub></a><br /><a href="https://github.com/tensorchord/envd/commits?author=dragonly" title="Documentation">📖</a> <a href="https://github.com/tensorchord/envd/issues?q=author%3Adragonly" title="Bug reports">🐛</a> <a href="https://github.com/tensorchord/envd/commits?author=dragonly" title="Code">💻</a></td>
      <td align="center"><a href="https://terrytangyuan.github.io/about/"><img src="https://avatars.githubusercontent.com/u/4269898?v=4?s=70" width="70px;" alt="Yuan Tang"/><br /><sub><b>Yuan Tang</b></sub></a><br /><a href="https://github.com/tensorchord/envd/commits?author=terrytangyuan" title="Code">💻</a> <a href="#design-terrytangyuan" title="Design">🎨</a> <a href="https://github.com/tensorchord/envd/commits?author=terrytangyuan" title="Documentation">📖</a> <a href="#ideas-terrytangyuan" title="Ideas, Planning, & Feedback">🤔</a></td>
    </tr>
    <tr>
      <td align="center"><a href="https://rudeigerc.dev/"><img src="https://avatars.githubusercontent.com/u/18243819?v=4?s=70" width="70px;" alt="Yuchen Cheng"/><br /><sub><b>Yuchen Cheng</b></sub></a><br /><a href="https://github.com/tensorchord/envd/issues?q=author%3Arudeigerc" title="Bug reports">🐛</a> <a href="#infra-rudeigerc" title="Infrastructure (Hosting, Build-Tools, etc)">🚇</a> <a href="#maintenance-rudeigerc" title="Maintenance">🚧</a> <a href="#tool-rudeigerc" title="Tools">🔧</a></td>
      <td align="center"><a href="http://lunarwhite.github.io"><img src="https://avatars.githubusercontent.com/u/57584831?v=4?s=70" width="70px;" alt="Yuedong Wu"/><br /><sub><b>Yuedong Wu</b></sub></a><br /><a href="https://github.com/tensorchord/envd/commits?author=lunarwhite" title="Code">💻</a></td>
      <td align="center"><a href="https://github.com/yczheng0"><img src="https://avatars.githubusercontent.com/u/21327543?v=4?s=70" width="70px;" alt="Yunchuan Zheng"/><br /><sub><b>Yunchuan Zheng</b></sub></a><br /><a href="https://github.com/tensorchord/envd/commits?author=yczheng0" title="Code">💻</a></td>
      <td align="center"><a href="http://lizheming.top"><img src="https://avatars.githubusercontent.com/u/9639449?v=4?s=70" width="70px;" alt="Zheming Li"/><br /><sub><b>Zheming Li</b></sub></a><br /><a href="https://github.com/tensorchord/envd/commits?author=lizhemingi" title="Code">💻</a></td>
      <td align="center"><a href="https://github.com/Xiaoaier-Z-L"><img src="https://avatars.githubusercontent.com/u/96805673?v=4?s=70" width="70px;" alt="Zhenguo.Li"/><br /><sub><b>Zhenguo.Li</b></sub></a><br /><a href="https://github.com/tensorchord/envd/commits?author=Xiaoaier-Z-L" title="Code">💻</a> <a href="https://github.com/tensorchord/envd/commits?author=Xiaoaier-Z-L" title="Documentation">📖</a></td>
      <td align="center"><a href="https://blog.triplez.cn/"><img src="https://avatars.githubusercontent.com/u/16285716?v=4?s=70" width="70px;" alt="Zhenzhen Zhao"/><br /><sub><b>Zhenzhen Zhao</b></sub></a><br /><a href="#infra-Triple-Z" title="Infrastructure (Hosting, Build-Tools, etc)">🚇</a> <a href="#userTesting-Triple-Z" title="User Testing">📓</a> <a href="https://github.com/tensorchord/envd/commits?author=Triple-Z" title="Code">💻</a></td>
      <td align="center"><a href="https://t.me/littlepoint"><img src="https://avatars.githubusercontent.com/u/7611700?v=4?s=70" width="70px;" alt="Zhizhen He"/><br /><sub><b>Zhizhen He</b></sub></a><br /><a href="https://github.com/tensorchord/envd/commits?author=hezhizhen" title="Code">💻</a> <a href="https://github.com/tensorchord/envd/commits?author=hezhizhen" title="Documentation">📖</a></td>
    </tr>
    <tr>
      <td align="center"><a href="https://github.com/cutecutecat"><img src="https://avatars.githubusercontent.com/u/19801166?v=4?s=70" width="70px;" alt="cutecutecat"/><br /><sub><b>cutecutecat</b></sub></a><br /><a href="https://github.com/tensorchord/envd/commits?author=cutecutecat" title="Code">💻</a></td>
      <td align="center"><a href="https://github.com/dqhl76"><img src="https://avatars.githubusercontent.com/u/69136320?v=4?s=70" width="70px;" alt="dqhl76"/><br /><sub><b>dqhl76</b></sub></a><br /><a href="https://github.com/tensorchord/envd/commits?author=dqhl76" title="Documentation">📖</a> <a href="https://github.com/tensorchord/envd/commits?author=dqhl76" title="Code">💻</a></td>
      <td align="center"><a href="https://github.com/jimoosciuc"><img src="https://avatars.githubusercontent.com/u/33337387?v=4?s=70" width="70px;" alt="jimoosciuc"/><br /><sub><b>jimoosciuc</b></sub></a><br /><a href="#userTesting-jimoosciuc" title="User Testing">📓</a></td>
      <td align="center"><a href="https://kenwoodjw.github.io"><img src="https://avatars.githubusercontent.com/u/10386710?v=4?s=70" width="70px;" alt="kenwoodjw"/><br /><sub><b>kenwoodjw</b></sub></a><br /><a href="https://github.com/tensorchord/envd/commits?author=kenwoodjw" title="Code">💻</a></td>
      <td align="center"><a href="https://github.com/aseaday"><img src="https://avatars.githubusercontent.com/u/3927355?v=4?s=70" width="70px;" alt="nullday"/><br /><sub><b>nullday</b></sub></a><br /><a href="#ideas-aseaday" title="Ideas, Planning, & Feedback">🤔</a> <a href="https://github.com/tensorchord/envd/commits?author=aseaday" title="Code">💻</a></td>
      <td align="center"><a href="https://tisonkun.org/"><img src="https://avatars.githubusercontent.com/u/18818196?v=4?s=70" width="70px;" alt="tison"/><br /><sub><b>tison</b></sub></a><br /><a href="https://github.com/tensorchord/envd/commits?author=tisonkun" title="Code">💻</a></td>
      <td align="center"><a href="http://fatelei.github.io"><img src="https://avatars.githubusercontent.com/u/961094?v=4?s=70" width="70px;" alt="wangxiaolei"/><br /><sub><b>wangxiaolei</b></sub></a><br /><a href="https://github.com/tensorchord/envd/commits?author=fatelei" title="Code">💻</a></td>
    </tr>
    <tr>
      <td align="center"><a href="https://github.com/sea-wyq"><img src="https://avatars.githubusercontent.com/u/22475606?v=4?s=70" width="70px;" alt="wyq"/><br /><sub><b>wyq</b></sub></a><br /><a href="https://github.com/tensorchord/envd/issues?q=author%3Asea-wyq" title="Bug reports">🐛</a> <a href="#design-sea-wyq" title="Design">🎨</a> <a href="https://github.com/tensorchord/envd/commits?author=sea-wyq" title="Code">💻</a></td>
      <td align="center"><a href="https://github.com/xiangtianyu"><img src="https://avatars.githubusercontent.com/u/10825900?v=4?s=70" width="70px;" alt="xiangtianyu"/><br /><sub><b>xiangtianyu</b></sub></a><br /><a href="https://github.com/tensorchord/envd/commits?author=xiangtianyu" title="Documentation">📖</a></td>
      <td align="center"><a href="https://github.com/xieydd"><img src="https://avatars.githubusercontent.com/u/20329697?v=4?s=70" width="70px;" alt="xieydd"/><br /><sub><b>xieydd</b></sub></a><br /><a href="https://github.com/tensorchord/envd/commits?author=xieydd" title="Code">💻</a></td>
      <td align="center"><a href="https://github.com/xing0821"><img src="https://avatars.githubusercontent.com/u/54933318?v=4?s=70" width="70px;" alt="xing0821"/><br /><sub><b>xing0821</b></sub></a><br /><a href="#ideas-xing0821" title="Ideas, Planning, & Feedback">🤔</a> <a href="#userTesting-xing0821" title="User Testing">📓</a> <a href="https://github.com/tensorchord/envd/commits?author=xing0821" title="Code">💻</a></td>
      <td align="center"><a href="https://github.com/zhyon404"><img src="https://avatars.githubusercontent.com/u/32242529?v=4?s=70" width="70px;" alt="zhyon404"/><br /><sub><b>zhyon404</b></sub></a><br /><a href="https://github.com/tensorchord/envd/commits?author=zhyon404" title="Code">💻</a></td>
      <td align="center"><a href="https://www.homeboyc.cn/"><img src="https://avatars.githubusercontent.com/u/22193008?v=4?s=70" width="70px;" alt="杨成锴"/><br /><sub><b>杨成锴</b></sub></a><br /><a href="https://github.com/tensorchord/envd/commits?author=asjdf" title="Code">💻</a></td>
    </tr>
  </tbody>
</table>

<!-- markdownlint-restore -->
<!-- prettier-ignore-end -->

<!-- ALL-CONTRIBUTORS-LIST:END -->

This project follows the [all-contributors](https://github.com/all-contributors/all-contributors) specification. Contributions of any kind welcome!

## Star History

[![Star History Chart](https://api.star-history.com/svg?repos=tensorchord/envd&type=Date)](https://star-history.com/#tensorchord/envd&Date)

## License 📋

[Apache 2.0](./LICENSE)

<a href="https://trackgit.com"><img src="https://us-central1-trackgit-analytics.cloudfunctions.net/token/ping/l3ldvdaswvnjpty9u7l3" alt="trackgit-views" /></a>
