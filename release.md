# Release Process

Our release process is fully automated using [GitHub Actions](https://github.com/falcosecurity/peribolos-syncer/actions) and [goreleaser](https://github.com/goreleaser/goreleaser) tool for artifacts.

When we release we do the following process:

1. We decide together (usually in the #falco channel in [slack](https://kubernetes.slack.com/messages/falco)) what's the next version to tag
2. A person with repository rights does the tag
3. The same person runs commands in their machine following the "Release commands" section below
4. Once the CI has done its job, the tag is live on [GitHub](https://github.com/falcosecurity/peribolos-syncer/releases) with the artifacts, and the container image is live on [GitHub Container Registry](https://github.com/falcosecurity/peribolos-syncer/pkgs/container/peribolos-syncer) with proper tags.

## Release commands

Tag the version

```bash
git tag -a v0.1.0-rc.0 -m "v0.1.0-rc.0"
git push origin v0.1.0-rc.0
```

