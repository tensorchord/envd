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
import logging

with open("README.md", "r", encoding="utf-8") as f:
    readme = f.read()


def build_envd_if_not_found():
    if not os.path.exists("bin/envd"):
        logging.info("envd not found. Build from scratch")
        errno = subprocess.call(["make", "build"])
        assert errno == 0, "Failed to build envd"


class EnvdExtension(Extension):
    """Extension for `envd`"""


class EnvdBuildExt(build_ext):
    def build_extension(self, ext: Extension) -> None:
        if not isinstance(ext, EnvdExtension):
            return super().build_extension(ext)

        bin_path = os.path.join(self.build_lib, "envd", "bin")
        build_envd_if_not_found()
        os.makedirs(bin_path, exist_ok=True)
        shutil.copy("bin/envd", bin_path)


def get_version():
    # Remove prefix v in versioning
    build_envd_if_not_found()
    version = subprocess.check_output(
        ["./bin/envd", "version", "--short"], universal_newlines=True
    ).strip()
    ver = version.rsplit(" ", 1)[-1][1:]
    return ver

classifiers = [
    'Development Status :: 3 - Alpha',
    'Topic :: Software Development :: Build Tools',
    'Intended Audience :: Science/Research',
    'Intended Audience :: Developers',
    'License :: OSI Approved :: Apache Software License',
    'Programming Language :: Python :: 3.6',
    'Programming Language :: Python :: 3.7',
    'Programming Language :: Python :: 3.8',
    'Programming Language :: Python :: 3.9',
    'Programming Language :: Python :: 3.10',
]


setup(
    name="envd",
    version=get_version(),
    use_scm_version=True,
    setup_requires=["setuptools_scm"],
    description="A development environment management tool for data scientists.",
    long_description=readme,
    long_description_content_type="text/markdown",
    url="https://github.com/tensorchord/envd",
    license="Apache License 2.0",
    author="TensorChord",
    author_email="envd-maintainers@tensorchord.ai",
    packages=find_packages(),
    include_package_data=True,
    entry_points={
        "console_scripts": [
            "envd=envd.cmd:envd",
        ],
    },
    classifiers=classifiers,
    zip_safe=False,
    ext_modules=[
        EnvdExtension(name="envd", sources=["cmd/*"]),
    ],
    cmdclass=dict(build_ext=EnvdBuildExt),
)
