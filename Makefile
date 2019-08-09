GOTEST ?= go test
PKGS ?= . \
	./compute \
	./dns \
	./storage \
	./runstatus
TESTOPTS = -v -timeout 120m -parallel 3 -count=1

help:
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-10s\033[0m %s\n", $$1, $$2}'

lint:
	@golangci-lint run --build-tags=testacc --skip-dirs internal/ $(LINTARGS)

test: ## Run unit tests
	@$(GOTEST) $(TESTOPTS) $(TESTARGS) $(PKGS)

testacc: ## Run acceptance tests (require Exoscale API credentials)
	@$(GOTEST) $(TESTOPTS) --tags=testacc $(TESTARGS) $(PKGS)

testall: test testacc ## Run all tests (unit + acceptance)

#TODO: docs

.PHONY: test testacc testall help
