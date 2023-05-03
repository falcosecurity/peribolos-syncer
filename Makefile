APP := peribolos-syncer

user := maxgio92
bins := git go gofumpt golangci-lint ginkgo

define declare_binpaths
$(1) = $(shell command -v 2>/dev/null $(1))
endef

$(foreach bin,$(bins),\
	$(eval $(call declare_binpaths,$(bin)))\
)

.PHONY: ginkgo
ginkgo:
	@hash ginkgo || \
		$(go) install github.com/onsi/ginkgo/v2/ginkgo

.PHONY: docs
docs:
	@go run docs/docs.go

.PHONY: build
build:
	@$(go) build .

.PHONY: run
run:
	@$(go) run .

.PHONY: test
test: ginkgo
	@$(ginkgo) ./...

.PHONY: lint
lint: golangci-lint
	@$(golangci-lint) run ./...

.PHONY: golangci-lint
golangci-lint:
	@$(go) install github.com/golangci/golangci-lint/cmd/golangci-lint@v1.50.0

.PHONY: gofumpt
gofumpt:
	@$(go) install mvdan.cc/gofumpt@v0.3.1

.PHONY: clean
clean:
	@rm -f $(APP)

.PHONY: help
help: list

.PHONY: list
list:
	@LC_ALL=C $(MAKE) -pRrq -f $(lastword $(MAKEFILE_LIST)) : 2>/dev/null | awk -v RS= -F: '/^# File/,/^# Finished Make data base/ {if ($$1 !~ "^[#.]") {print $$1}}' | sort | egrep -v -e '^[^[:alnum:]]' -e '^$@$$'
