// Copyright 2023 The Falco Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//	http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package local

const (
	commandName             = "local"
	commandShortDescription = "Synchronize Peribolos config on local filesystem"
	commandExample          = `
peribolos-syncer sync local --owners-file OWNERS --peribolos-config org.yaml --org acme --team app-maintainers
`

	flagOwnersFilePath             = "owners-file"
	flagPeribolosConfigFilepath    = "orgs-config"
	defaultOwnersFilepath          = "OWNERS"
	defaultPeribolosConfigFilepath = "org.yaml"
	FilePerm                       = 0o644
)
