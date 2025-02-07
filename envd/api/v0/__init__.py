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

"""Global functions

::: tip
Note that the documentation is automatically generated from [envd/api](https://github.com/tensorchord/envd/tree/main/envd/api) folder
in [tensorchord/envd](https://github.com/tensorchord/envd/tree/main/envd/api) repo.
Please update the python file there instead of directly editing file inside envd-docs repo.
:::

::: warning
v0 syntax is no longer supported from `envd>=v1.0`. Check the [upgrade guide](https://envd.tensorchord.ai/guide/v1.html).

If you want to use v0 syntax, try `pip install 'envd<1'`.
:::
"""

from typing import List, Optional


def base(os: str, language: str, image: Optional[str]):
    """Set base image

    Args:
        os (str): The operating system (i.e. `ubuntu22.04`)
        language (str): The programming language dependency (i.e. `python3.8`)
        image (Optional[str]): Custom image (i.e. `python:3.11-slim`)
    """


def shell(name: str):
    """Interactive shell

    Args:
        name (str): shell name (i.e. `zsh`, `bash`)
    """


def run(commands: List[str], mount_host: bool = False):
    """Execute command

    Args:
        commands (List[str]): command to run during the building process
        mount_host (bool): mount the host directory. Default is False.
            Enabling this will disable the build cache for this operation.

    Example:
    ```python
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
        name (optional, str): User name
        email (optional, str): User email
        editor (optional, str): Editor for git operations

    Example usage:
    ```python
    git_config(name="My Name", email="my@email.com", editor="vim")
    ```
    """


def include(git: str):
    """Import from another git repo

    This will pull the git repo and execute all the `envd` files. The return value will be a module
    contains all the variables/functions defined (except the ones with `_` prefix).

    Args:
        git (str): git URL

    Example usage:
    ```python
    envd = include("https://github.com/tensorchord/envdlib")

    def build():
        base(os="ubuntu22.04", language="python")
        envd.tensorboard(8000)
    ```
    """
