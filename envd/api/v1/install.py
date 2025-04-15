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
"""

from typing import Optional, Sequence


def python(version: str = "3.11"):
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


def pixi(use_pixi_mirror: bool = False, pypi_index: Optional[str] = None):
    """Install Pixi (https://github.com/prefix-dev/pixi).

    `pixi` is an alternative to `conda` that is written in Rust and provides faster
    dependency resolution and installation. It also simplify the project management.

    This doesn't support installing Python packages through `install.python_packages`
    because that part should be managed by `pixi`. You can run `pixi shell` in the
    `envd` environment to sync all the dependencies.

    Args:
        use_pixi_mirror (bool): use pixi mirror
        pypi_index (Optional[str]): customize pypi index url
    """


def uv(python_version: str = "3.11"):
    """Install UV (an extremely fast Python package and project manager).

    `uv` is much faster than `conda`. Choose this one instead of `conda` if you don't
    need any machine learning packages.

    This doesn't support installing Python packages through `install.python_packages`
    because that part should be managed by `uv`. You can run `uv sync` in the `envd`
    environment to install all the dependencies.

    Args:
        python_version (str): install this Python version through UV
    """


def r_lang():
    """Install R Lang."""


def julia():
    """Install Julia."""


def apt_packages(name: Sequence[str] = ()):
    """Install package using the system package manager (apt on Ubuntu).

    Args:
        name (Sequence[str]): apt package name list
    """


def python_packages(
    name: Sequence[str] = (), requirements: str = "", local_wheels: Sequence[str] = ()
):
    """Install python package by pip.

    Args:
        name (Sequence[str]): package name list
        requirements (str): requirements file path
        local_wheels (Sequence[str]): local wheels
            (wheel files should be placed under the current directory)
    """


def conda_packages(
    name: Sequence[str] = (), channel: Sequence[str] = (), env_file: str = ""
):
    """Install python package by Conda

    Args:
        name (Sequence[str]): List of package names with optional version assignment,
            such as ['pytorch', 'tensorflow==1.13.0']
        channel (Sequence[str]): additional channels
        env_file (str): conda env file path
    """


def r_packages(name: Sequence[str]):
    """Install R packages by R package manager.

    Args:
        name (Sequence[str]): package name list
    """


def julia_packages(name: Sequence[str]):
    """Install Julia packages.

    Args:
        name (Sequence[str]): List of Julia packages
    """


def vscode_extensions(name: Sequence[str]):
    """Install VS Code extensions

    Args:
        name (Sequence[str]): extension names, such as ['ms-python.python']
    """


def cuda(version: str, cudnn: Optional[str] = "8"):
    """Replace the base image with a `nvidia/cuda` image.

    This will replace the default base image to an `nvidia/cuda` image. You can
    also use a CUDA base image directly like
    `base(image="nvidia/cuda:12.2.0-devel-ubuntu22.04", dev=True)`.

    Args:
        version (str): CUDA version, such as '11.6.2'
        cudnn (optional, str): CUDNN version, such as '8'

    Example usage:
    ```python
    install.cuda(version="11.6.2", cudnn="8")
    ```
    """
