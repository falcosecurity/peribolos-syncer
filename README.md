# Peribolos config syncer

The tool synchronizes [Peribolos](https://docs.prow.k8s.io/docs/components/cli-tools/peribolos/) GitHub organization config with GitHub repository's [OWNERS](https://docs.prow.k8s.io/docs/components/plugins/approve/approvers/#overview) files.

## Usage

### Synchronize local files

```shell
peribolos-syncer sync local \
  --peribolos-config=/path/to/org.yaml \
  --owners-file=/path/to/OWNERS \
  --organization=<name> \
  --team=<name>
```

### Synchronize remote files with Pull Request

```shell
peribolos-syncer sync github \
--org <GithHub organization> \
--team <GitHub team> \
--org-config </path/to/org.yaml> \
--org-config-repository <Peribolos org config GitHub repository> \
--owners-repository <GitHub OWNERS repository> \
--github-token-path=</path/to/github_token>
```

## Goals

- Synchronize `org.yaml` Peribolos desired state file, parsing `OWNERS` file.

## Non-goals

- Synchronize Github organization and settings.
