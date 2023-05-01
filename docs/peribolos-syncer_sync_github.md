---
title: peribolos-syncer sync github
---	

## peribolos-syncer sync github

Synchronize Peribolos config on remote GitHub repositories via Pull Request

```
peribolos-syncer sync github [flags]
```

### Examples

```

peribolos-syncer sync github --org=acme --team=app-maintainers
--peribolos-config-path=config/org.yaml --peribolos-config-repository=community --peribolos-config-git-ref=main
--owners-repository=app --owners-git-ref=main --owners-path=OWNERS
--github-username=bot --github-token-path=./bot_token
--git-author-name=bot --git-author-email="bot@acme.org"
--gpg-public-key=./bot.pub --gpg-private-key=./bot.asc

```

### Options

```
      --dry-run                                  Dry run for testing. Uses API tokens but does not mutate.
      --git-author-email string                  The Git author email with which write commits for the update of the Peribolos config
      --git-author-name string                   The Git author name with which write commits for the update of the Peribolos config
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
      --gpg-private-key string                   The path to the private GPG key for signing git commits
      --gpg-public-key string                    The path to the public GPG key for signing git commits
  -h, --help                                     help for github
      --org string                               The name of the GitHub organization to update configuration for
  -r, --owners-git-ref string                    The base Git reference at which parse the OWNERS hierarchy (default "master")
  -o, --owners-path string                       The path to the OWNERS file from the root of the Git repository. Ignored with sync-github.
      --owners-repository string                 The name of the github repository from which parse OWNERS file
      --peribolos-config-git-ref string          The base Git reference at which pull the peribolos config repository (default "master")
  -c, --peribolos-config-path string             The path to the peribolos organization config file from the root of the Git repository (default "org.yaml")
      --peribolos-config-repository string       The name of the github repository that contains the peribolos organization config file
      --team string                              The name of the GitHub team to update configuration for
```

### SEE ALSO

* [peribolos-syncer sync](peribolos-syncer_sync.md)	 - Synchronize Peribolos config with external GitHub people source of truth

