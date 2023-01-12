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

"""Config functions

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

from typing import List, Optional


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


def pip_index(url: str, extra_url: str = "", trust: bool = False):
    """Configure pypi index mirror

    Args:
        url (str): PyPI index URL (i.e. https://mirror.sjtu.edu.cn/pypi/web/simple)
        extra_url (str): PyPI extra index URL. `url` and `extra_url` will be
            treated equally, see https://github.com/pypa/pip/issues/8606
        trust (bool): trust the provided index
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


def owner(uid: int, gid: int):
    """Configure uid:gid as the environmen owner.
    This can also be achieved by using flag `envd --owner uid:gid build` or environment
        variable `ENVD_BUILD_OWNER=uid:gid envd build`

    Args:
        uid (int): UID
        gid (int): GID
    """
