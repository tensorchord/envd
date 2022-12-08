# Copyright 2022 The envd Authors
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#      http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

"""#v1 syntax

::: tip
Note that the documentation is automatically generated from [envd/api](https://github.com/tensorchord/envd/tree/main/envd/api) folder
in [tensorchord/envd](https://github.com/tensorchord/envd/tree/main/envd/api) repo.
Please update the python file there instead of directly editing file inside envd-docs repo.
:::

::: warning
Enable v1 by adding `# syntax=v1` to the 1st line of your envd file.

v1 is experimental and may change in the future. Make sure to freeze the envd version for online CI/CD.
:::
"""

from typing import Dict, List, Optional

"""## Universe
"""


def base(image: str = "ubuntu:20.04", dev: bool = False):
    """Set up the base env.

    Args:
        image (str): docker image, can be any Debian-based images
        dev (bool): enabling the dev env will add lots of development related libraries like
            envd-sshd, vim, git, shell prompt, vscode extensions, etc.
    """


def shell(name: str = "base"):
    """Interactive shell

    Args:
        name (str): shell name (i.e. `zsh`, `bash`)
    """


def run(commands: str, mount_host: bool = False):
    """Execute command

    Args:
        commands (str): command to run during the building process
        mount_host (bool): mount the host directory. Default is False.
            Enabling this will disable the build cache for this operation.

    Example:
    ```
    run(commands=["conda install -y -c conda-forge exa"])
    ```
    """


def git_config(
    name: Optional[str] = None,
    email: Optional[str] = None,
    editor: Optional[str] = None,
):
    """Setup git config

    Args:
        name (str): User name
        email (str): User email
        editor (str): Editor for git operations

    Example usage:
    ```
    git_config(name="My Name", email="my@email.com", editor="vim")
    ```
    """


def include(git: str):
    """Import from another git repo

    This will pull the git repo and execute all the `envd` files. The return value will be a module
    contains all the variables/functions defined (expect those has `_` prefix).

    Args:
        git (str): git URL

    Example usage:
    ```
    envd = include("https://github.com/tensorchord/envdlib")

    def build():
        base(os="ubuntu20.04", language="python")
        envd.tensorboard(host_port=8000)
    ```
    """


class install:
    """## Install

    All the following modules should be called with `install.` prefix.
    """

    def python(version: str = "3.9"):
        """Install python.

        If `install.conda` is not used, this will create a solo Python environment. Otherwise, it
        will be a conda environment.

        Args:
            version (str): Python version
        """

    def conda(use_mamba: bool = False):
        """Install MiniConda or MicroMamba.

        Args:
            use_mamba (bool): use mamba instead of conda
        """

    def r_lang():
        """Install R Lang.

        Not implemented yet. Please use v0 if you need R.
        """

    def julia():
        """Install Julia.

        Not implemented yet. Please use v0 if you need Julia.
        """

    def apt_packages(name: List[str] = []):
        """Install package by system-level package manager (apt on Ubuntu).

        Args:
            name (str): apt package name list
        """

    def python_packages(
        name: List[str] = [], requirements: str = "", local_wheels: List[str] = []
    ):
        """Install python package by pip.

        Args:
            name (List[str]): package name list
            requirements (str): requirements file path
            local_wheels (List[str]): local wheels
                (wheel files should be placed under the current directory)
        """

    def conda_packages(
        name: List[str] = [], channel: List[str] = [], env_file: str = ""
    ):
        """Install python package by Conda

        Args:
            name (List[str]): List of package names with optional version assignment,
                such as ['pytorch', 'tensorflow==1.13.0']
            channel (List[str]): additional channels
            env_file (str): conda env file path
        """

    def r_packages(name: List[str]):
        """Install R packages by R package manager.

        Not implemented yet. Please use v0 if you need R.

        Args:
            name (List[str]): package name list
        """

    def julia_packages(name: List[str]):
        """Install Julia packages.

        Not implemented yet. Please use v0 if you need Julia.

        Args:
            name (List(str)): List of Julia packages
        """

    def vscode_extensions(name: List[str]):
        """Install VS Code extensions

        Args:
            name (List[str]): extension names, such as ['ms-python.python']
        """

    def cuda(version: str, cudnn: Optional[str] = "8"):
        """Install CUDA dependency

        Args:
            version (str): CUDA version, such as '11.6.2'
            cudnn (optional, str): CUDNN version, such as '8'
        """


class config:
    """## Config

    All the following modules should be called with `config.` prefix.
    """

    def apt_source(source: Optional[str]):
        """Configure apt sources

        Example usage:
        ```
        apt_source(source='''
            deb https://mirror.sjtu.edu.cn/ubuntu focal main restricted
            deb https://mirror.sjtu.edu.cn/ubuntu focal-updates main restricted
            deb https://mirror.sjtu.edu.cn/ubuntu focal universe
            deb https://mirror.sjtu.edu.cn/ubuntu focal-updates universe
            deb https://mirror.sjtu.edu.cn/ubuntu focal multiverse
            deb https://mirror.sjtu.edu.cn/ubuntu focal-updates multiverse
            deb https://mirror.sjtu.edu.cn/ubuntu focal-backports main restricted universe multiverse
            deb http://archive.canonical.com/ubuntu focal partner
            deb https://mirror.sjtu.edu.cn/ubuntu focal-security main restricted universe multiverse
        ''')
        ```

        Args:
            source (str, optional): The apt source configuration
        """

    def jupyter(token: str, port: int):
        """Configure jupyter notebook configuration

        Args:
            token (str): Token for access authentication
            port (int): Port to serve jupyter notebook
        """

    def pip_index(url: str, extra_url: str):
        """Configure pypi index mirror

        Args:
            url (str): PyPI index URL (i.e. https://mirror.sjtu.edu.cn/pypi/web/simple)
            extra_url (str): PyPI extra index URL. `url` and `extra_url` will be
                treated equally, see https://github.com/pypa/pip/issues/8606
        """

    def conda_channel(channel: str):
        """Configure conda channel mirror

        Example usage:
        ```
        config.conda_channel(channel='''
        channels:
            - defaults
        show_channel_urls: true
        default_channels:
            - https://mirrors.tuna.tsinghua.edu.cn/anaconda/pkgs/main
            - https://mirrors.tuna.tsinghua.edu.cn/anaconda/pkgs/r
            - https://mirrors.tuna.tsinghua.edu.cn/anaconda/pkgs/msys2
        custom_channels:
            conda-forge: https://mirrors.tuna.tsinghua.edu.cn/anaconda/cloud
        ''')
        ```

        Args:
            channel (str): Basically the same with file content of an usual .condarc
        """

    def entrypoint(args: List[str]):
        """Configure entrypoint for custom base image

        Example usage:
        ```
        config.entrypoint(["date", "-u"])
        ```

        Args:
            args (List[str]): list of arguments to run
        """

    def gpu(count: int):
        """Configure the number of GPUs required

        Example usage:
        ```
        config.gpu(count=2)
        ```

        Args:
            count (int): number of GPUs
        """

    def cran_mirror(url: str):
        """Configure the mirror URL, default is https://cran.rstudio.com

        Args:
            url (str): mirror URL
        """

    def julia_pkg_server(url: str):
        """Configure the package server for Julia.
        Since Julia 1.5, https://pkg.julialang.org is the default pkg server.

        Args:
            url (str): Julia pkg server URL
        """

    def rstudio_server():
        """
        Enable the RStudio Server (only work for `base(os="ubuntu20.04", language="r")`)
        """

    def repo(url: str, description: str):
        """Setup repo related information. Will save to the image labels.

        Args:
            url (str): repo URL
            description (str): repo description
        """


class runtime:
    """# Runtime

    All the following modules should be called with `runtime.` prefix.
    """

    def command(commands: Dict[str, str]):
        """Execute commands during runtime

        Args:
            commands (Dict[str, str]): map name to command, similar to Makefile

        Example usage:
        ```
        runtime.command(commands={
            "train": "python train.py --epoch 20 --notify me@tensorchord.ai",
            "run": "python server.py --batch 1 --host 0.0.0.0 --port 8000",
        })
        ```

        You can run `envd run --command train` to train the model.
        """

    def expose(
        envd_port: str,
        host_port: Optional[str],
        service: Optional[str],
        listen_addr: Optional[str],
    ):
        """Expose port to host
        Proposal: https://github.com/tensorchord/envd/pull/780

        Args:
            envd_port (str): port in `envd` container
            host_port (Optional[str]): port in the host, if not provided or
                `host_port=0`, `envd` will randomly choose a free port
            service (Optional[str]): service name
            listen_addr (Optional[str]): address to listen on
        """

    def daemon(commands: List[List[str]]):
        """Run daemon processes in the container
        Proposal: https://github.com/tensorchord/envd/pull/769

        It's better to redirect the logs to local files for debug purposes.

        Args:
            commands (List[List[str]]): run multiple commands in the background

        Example usage:
        ```
        runtime.daemon(commands=[
            ["jupyter-lab", "--port", "8080"],
            ["python3", "serving.py", ">>serving.log", "2>&1"],
        ])
        ```
        """

    def environ(env: Dict[str, str], extra_path: List[str]):
        """Add runtime environments

        Args:
            env (Dict[str, str]): environment name to value
            extra_path (List[str]): additional PATH

        Example usage:
        ```
        runtime.environ(env={"ENVD_MODE": "DEV"}, extra_path=["/usr/bin/go/bin"])
        ```
        """

    def mount(host_path: str, envd_path: str):
        """Mount from host path to container path (runtime)

        Args:
            host_path (str): source path in the host machine
            envd_path (str): destination path in the envd container
        """

    def init(commands: List[str]):
        """Commands to be executed when start the container

        Args:
            commands (List[str]): list of commands
        """


class io:
    """IO

    All the following modules should be called with `io.` prefix.
    """

    def copy(host_path: str, envd_path: str):
        """Copy from host path to container path (build time)

        Args:
            host_path (str): source path in the host machine
            envd_path (str): destination path in the envd container
        """

    def http(url: str, checksum: Optional[str], filename: Optional[str]):
        """Download file with HTTP to `/home/envd/extra_source`

        Args:
            url (str): URL
            checksum (Optional[str]): checksum for the downloaded file
            filename (Optional[str]): rewrite the filename
        """
