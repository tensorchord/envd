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

import os
import shutil
from setuptools import setup, Extension, find_packages
from setuptools.command.build_ext import build_ext
import subprocess


with open("README.md", "r", encoding="utf-8") as f:
    readme = f.read()


class EnvdExtension(Extension):
    """Extension for `envd`"""


class EnvdBuildExt(build_ext):
    def build_extension(self, ext: Extension) -> None:
        if not isinstance(ext, EnvdExtension):
            return super().build_extension(ext)

        bin_path = os.path.join(self.build_lib, "envd", "bin")
        errno = subprocess.call(["make", "build"])
        assert errno == 0, "Failed to build envd"
        os.makedirs(bin_path, exist_ok=True)
        shutil.copy("bin/envd", bin_path)


def get_version():
    subprocess.call(["make", "build"])
    version = subprocess.check_output(
        ["./bin/envd", "-v"], universal_newlines=True
    ).strip()
    return version.rsplit(" ", 1)[-1]


setup(
    name="envd",
    version=get_version(),
    description="A development environment management tool for data scientists.",
    long_description=readme,
    url="https://github.com/tensorchord/envd",
    license="Apache License 2.0",
    author="TensorChord",
    author_email="envd-maintainers@tensorchord.ai",
    packages=find_packages(),
    include_package_data=True,
    python_requires=">=3.6",
    entry_points={
        "console_scripts": [
            "envd=envd.cmd:envd",
        ],
    },
    zip_safe=False,
    ext_modules=[
        EnvdExtension(name="envd", sources=["cmd/*"]),
    ],
    cmdclass=dict(build_ext=EnvdBuildExt),
)
