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
