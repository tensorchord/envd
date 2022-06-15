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

from typing import List, Optional


def system_packages(name: List[str]):
    """Install package by system-level package manager(apt on Ubuntu)
    """
    pass


def python_packages(name: List[str]):
    """Install python package by pip
    """
    pass

def conda_packages(name: List[str]):
    """Install python package by Conda

    Args:
        name (List[str]): List of package names with optional version assignment, 
            such as ['pytorch', 'tensorflow==1.13.0']
    """

def r_packages(name: List[str]):
    """Install R packages by R package manager
    """
    pass

def cuda(version: str, cudnn: Optional[str] = None):
    """Install CUDA dependency
    Args:
      version (str): CUDA version, such as '11.6'
      cudnn (optional, str): CUDNN version, such as '6'
    """

def vscode_extensions(name: List[str]):
    """Install VS Code extensions
    Args:
      name (list of str): extension names, such as ['ms-python.python']
    """
