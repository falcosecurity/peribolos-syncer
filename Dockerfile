FROM scratch
COPY peribolos-owners-syncer /peribolos-owners-syncer
ENTRYPOINT ["/peribolos-owners-syncer"]

