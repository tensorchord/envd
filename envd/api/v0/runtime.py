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

"""Runtime functions

::: tip
Note that the documentation is automatically generated from [envd/api](https://github.com/tensorchord/envd/tree/main/envd/api) folder
in [tensorchord/envd](https://github.com/tensorchord/envd/tree/main/envd/api) repo.
Please update the python file there instead of directly editing file inside envd-docs repo.
:::
"""

from typing import Dict, Optional, List


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

    You can run `envd exec --command train` to train the model.
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
