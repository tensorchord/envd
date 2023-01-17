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

"""Install functions

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
    """Install R Lang."""


def julia():
    """Install Julia."""


def apt_packages(name: List[str] = []):
    """Install package by system-level package manager (apt on Ubuntu).

    Args:
        name (List[str]): apt package name list
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


def conda_packages(name: List[str] = [], channel: List[str] = [], env_file: str = ""):
    """Install python package by Conda

    Args:
        name (List[str]): List of package names with optional version assignment,
            such as ['pytorch', 'tensorflow==1.13.0']
        channel (List[str]): additional channels
        env_file (str): conda env file path
    """


def r_packages(name: List[str]):
    """Install R packages by R package manager.

    Args:
        name (List[str]): package name list
    """


def julia_packages(name: List[str]):
    """Install Julia packages.

    Args:
        name (List[str]): List of Julia packages
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
