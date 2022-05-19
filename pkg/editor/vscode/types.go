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

package vscode

import "fmt"

const (
	vendorVSCodeTemplate  = "https://%s.gallery.vsassets.io/_apis/public/gallery/publisher/%s/extension/%s/%s/assetbyname/Microsoft.VisualStudio.Services.VSIXPackage"
	vendorOpenVSXTemplate = "https://open-vsx.org/api/%s/%s/latest"
)

type MarketplaceVendor string

const (
	MarketplaceVendorVSCode  MarketplaceVendor = "vscode"
	MarketplaceVendorOpenVSX MarketplaceVendor = "openvsx"
)

type Plugin struct {
	Publisher string
	Extension string
	Version   *string
}

func (p Plugin) String() string {
	if p.Version != nil {
		return fmt.Sprintf("%s.%s-%s", p.Publisher, p.Extension, *p.Version)
	}
	return fmt.Sprintf("%s.%s", p.Publisher, p.Extension)
}
