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


def copy(src: str, dest: str):
    """Copy from host `src` to container `dest` (build time)

    Args:
        src (str): source path
        dest (str): destination path
    """


def mount(src: str, dest: str):
    """Mount from host `src` to container `dest` (runtime)

    Args:
        src (str): source path
        dest (str): destination path
    """
