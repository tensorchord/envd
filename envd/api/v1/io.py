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

"""IO functions

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

from typing import Optional


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
