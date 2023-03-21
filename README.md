# Peribolos config syncer

The tool synchronizes [Peribolos](https://docs.prow.k8s.io/docs/components/cli-tools/peribolos/) GitHub organization config with GitHub repository's [OWNERS](https://docs.prow.k8s.io/docs/components/plugins/approve/approvers/#overview) files.

## Usage

### Local files

The `sync local` retrieves both `OWNERS` and `orgs.yaml` files from **local filesystem**.

It updates the specified GitHub team according to the leaf approvers, in-place.

```shell
syncer sync local [flags]

Flags:
  -h, --help                 help for local
      --org string           The name of the GitHub organization to update
  -c, --orgs-config string   The path to the Peribolos org.yaml file (default "org.yaml")
  -o, --owners-file string   The path to the OWNERS file (default "OWNERS")
      --team string          The name of the GitHub organization to update
```

### Remote files

The `sync github` retrieves both `OWNERS` and `orgs.yaml` files from **GitHub repositories**.

It updates the specified GitHub team according to the leaf approvers in-place, via a **Pull Request** against the `orgs.yaml` repository base reference.

```shell
syncer sync github [flags]

Flags:
      --author-email string                      The Git author email with which write commits for the update of the Peribolos config
      --author-name string                       The Git author name with which write commits for the update of the Peribolos config
      --dry-run                                  Dry run for testing. Uses API tokens but does not mutate. (default true)
      --github-allowed-burst int                 Size of token consumption bursts. If set, --github-hourly-tokens must be positive too and set to a higher or equal number.
      --github-app-id string                     ID of the GitHub app. If set, requires --github-app-private-key-path to be set and --github-token-path to be unset.
      --github-app-private-key-path string       Path to the private key of the github app. If set, requires --github-app-id to bet set and --github-token-path to be unset
      --github-client.backoff-timeout duration   Largest allowable Retry-After time for requests to the GitHub API. (default 2m0s)
      --github-client.initial-delay duration     Initial delay before retries begin for requests to the GitHub API. (default 2s)
      --github-client.max-404-retries int        Maximum number of retries that will be used for a 404-ing request to the GitHub API. (default 2)
      --github-client.max-retries int            Maximum number of retries that will be used for a failing request to the GitHub API. (default 8)
      --github-client.request-timeout duration   Timeout for any single request to the GitHub API. (default 2m0s)
      --github-endpoint Strings                  GitHub's API endpoint (may differ for enterprise). (default https://api.github.com)
      --github-graphql-endpoint string           GitHub GraphQL API endpoint (may differ for enterprise). (default "https://api.github.com/graphql")
      --github-host string                       GitHub's default host (may differ for enterprise) (default "github.com")
      --github-hourly-tokens int                 If set to a value larger than zero, enable client-side throttling to limit hourly token consumption. If set, --github-allowed-burst must be positive too.
      --github-throttle-org Strings              Throttler settings for a specific org in org:hourlyTokens:burst format. Can be passed multiple times. Only valid when using github apps auth.
      --github-token-path string                 Path to the file containing the GitHub OAuth secret.
      --github-username string                   The GitHub username
  -h, --help                                     help for github
      --org string                               The name of the GitHub organization to update configuration for
  -c, --orgs-config string                       The path to the Peribolos organization config file from the root of the Git repository (default "/org.yaml")
      --orgs-config-base-ref string              The base Git reference at which pull the Peribolos config repository (default "master")
      --orgs-config-repository string            The name of the github repository that contains the Peribolos organization config file
  -o, --owners-file string                       The path to the OWNERS file from the root of the Git repository
  -r, --owners-reference string                  The base Git reference at which parse the OWNERS hierarchy (default "master")
      --owners-repository string                 The name of the github repository from which parse OWNERS file
      --team string                              The name of the GitHub team to update configuration for
```

#### Example

```shell
syncer sync github \
  --org maxgio92 --team foo-maintainers --owners-repository "foo" \
  --orgs-config org.yaml --orgs-config-repository ".github" \
  --github-username=mybot --github-token-path=/path/to/token \
  --git-author-name="My Bot" --git-author-email="bot@example.org" \
  --dry-run=false
```

## Goals

- Synchronize Github teams in a Peribolos configuration, from `OWNERS` leaf approvers.

## Non-goals

- Synchronize Github organization and settings.
