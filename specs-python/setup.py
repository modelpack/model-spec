#     Copyright 2025 The CNCF ModelPack Authors
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

from setuptools import setup, find_packages

setup(
    name="modelpack",
    version="0.1.0",
    description="Python SDK for the CNCF ModelPack specification",
    packages=find_packages(),
    package_data={"modelpack.v1": ["config-schema.json"]},
    python_requires=">=3.10",
    install_requires=[
        "jsonschema[format]>=4.20.0",
    ],
    extras_require={
        "dev": [
            "pytest>=7.0",
            "ruff>=0.4.0",
        ],
    },
)
