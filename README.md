# Peribolos config syncer

Tool to synchronize [Peribolos](https://docs.prow.k8s.io/docs/components/cli-tools/peribolos/) configuration from external sources.

It synchronizes GitHub Team configurations with external GitHub people source of truth, like [OWNERS](https://docs.prow.k8s.io/docs/components/plugins/approve/approvers/#overview) files.

## Supported source of truth

The currently supported GitHub people source of truth are:

* [OWNERS](https://docs.prow.k8s.io/docs/components/plugins/approve/approvers/#overview).

## Usage

### Local files

The `sync local` retrieves both the GitHub people source of truth and Peribolos config files from a **local filesystem**.

It updates in-place the specified GitHub team according to the approvers of the specified [OWNERS](https://docs.prow.k8s.io/docs/components/plugins/approve/approvers/#overview) file.

#### Documentation

Please refer to the [`sync local`](./docs/peribolos-syncer_sync_local.md) command documentation.

### Remote files

The `sync github` synchronizes Peribolos config on remote **GitHub repositories** via Pull Request.

It updates the specified GitHub team according to the approvers of the specified [OWNERS](https://docs.prow.k8s.io/docs/components/plugins/approve/approvers/#overview) file.

#### Documentation

Please refer to the [`sync github`](./docs/peribolos-syncer_sync_github.md) command documentation.

## Goals

- Synchronize Github teams in a Peribolos configuration.

## Non-goals

- Synchronize Github organization and settings.
