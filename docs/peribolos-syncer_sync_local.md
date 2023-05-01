---
title: peribolos-syncer sync local
---	

## peribolos-syncer sync local

Synchronize Peribolos config on local filesystem

```
peribolos-syncer sync local [flags]
```

### Examples

```

peribolos-syncer sync local --owners-file OWNERS --peribolos-config org.yaml --org acme --team app-maintainers

```

### Options

```
  -h, --help                 help for local
      --org string           The name of the GitHub organization to update
  -c, --orgs-config string   The path to the Peribolos org.yaml file (default "org.yaml")
  -o, --owners-file string   The path to the OWNERS file (default "OWNERS")
      --team string          The name of the GitHub organization to update
```

### SEE ALSO

* [peribolos-syncer sync](peribolos-syncer_sync.md)	 - Synchronize Peribolos config with external GitHub people source of truth

