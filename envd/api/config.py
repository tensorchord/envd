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

from typing import Optional


def apt_source(mode: Optional[str], source: Optional[str]):
    r"""Configure apt sources

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
    '''
    )
    ```

    Args:
        mode (str, optional): This argument is not supported currently
        source (str, optional): The apt source configuration
    """
    pass


def jupyter(password: str, port: int):
    """
    Configure jupyter notebook configuration

    Args:
        password (str): Password for access authenticatioin
        port (int): Port to serve jupyter notebook
    """
    pass

def pip_index(url):
    """
    Configure pypi index mirror

    Args:
        url (str): Pypi index url (i.e. https://mirror.sjtu.edu.cn/pypi/web/simple)
    """
    pass

def conda_channel(channel:str):
    """
    Configure conda channel mirror

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
    '''
    )
    ```

    Args:
        channel (str): Basically the same with file content of an usual .condarc
    """
