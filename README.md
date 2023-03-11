# Peribolos config syncer

The tool synchronizes the org.yaml Peribolos config file parsing the OWNERS standard file.

## Quickstart

```shell
peribolos-syncer sync \
  --peribolos-config=/path/to/org.yaml \
  --owners-file=/path/to/OWNERS \
  --organization=<name> \
  --team=<name>
```

## Goals

- Synchronize `org.yaml` Peribolos desired state file, parsing `OWNERS` file.

## Non-goals

- Synchronize Github organization and settings.

## The supported source of the `org.yaml` file

- Local filesystem.
- (next) Remote Github repository.

## The supported target of the `org.yaml` file

- Local filesystem.
- (next) Pull request on the `org.yaml` repository.

