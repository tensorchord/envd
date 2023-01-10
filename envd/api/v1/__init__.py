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

"""# v1 syntax

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


def run(commands: List[str], mount_host: bool = False):
    """Execute command

    Args:
        commands (List[str]): command to run during the building process
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
