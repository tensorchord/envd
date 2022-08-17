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
    """Execute commands

    Args:
        commands (Dict[str, str]): map name to command, similar to Makefile
    """


def expose(envd_port: str, host_port: Optional[str], service: Optional[str]):
    """Expose port to host
    Proposal: https://github.com/tensorchord/envd/pull/780

    Args:
        envd_port (str): port in `envd` container
        host_port (Optional[str]): port in the host, if not provided, `envd` will
            randomly choose a free port
        service (Optional[str]): service name
    """


def daemon(commands: List[List[str]]):
    """Run daemon processes in the container
    Proposal: https://github.com/tensorchord/envd/pull/769

    It's better to redirect the logs to local files for debug purposes.

    Args:
        commands (List[List[str]]): run multiple commands in the background

    Example usage:
    ```
    runtime.daemon([
        ["jupyter-lab", "--port", "8080"],
        ["python3", "serving.py", ">>serving.log", "2>&1"],
    ])
    ```
    """
