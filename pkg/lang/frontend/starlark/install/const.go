// Copyright 2022 The envd Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package install

const (
	// language
	rulePython = "install.python"
	ruleConda  = "install.conda"
	ruleRLang  = "install.r_lang"
	ruleJulia  = "install.julia"

	// packages
	ruleSystemPackage = "install.apt_packages"
	rulePyPIPackage   = "install.python_packages"
	ruleCondaPackages = "install.conda_packages"
	ruleRPackage      = "install.r_packages"
	ruleJuliaPackages = "install.julia_packages"

	// others
	ruleCUDA   = "install.cuda"
	ruleVSCode = "install.vscode_extensions"
)
