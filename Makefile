GO_MK_REF := v1.0.0

# make go.mk a dependency for all targets
.EXTRA_PREREQS = go.mk

ifndef MAKE_RESTARTS
# This section will be processed the first time that make reads this file.

# This causes make to re-read the Makefile and all included
# makefiles after go.mk has been cloned.
Makefile:
	@touch Makefile
endif

.PHONY: go.mk
.ONESHELL:
go.mk:
	@if [ ! -d "go.mk" ]; then
		git clone https://github.com/exoscale/go.mk.git
	fi
	@cd go.mk
	@if ! git show-ref --quiet --verify "refs/heads/${GO_MK_REF}"; then
		git fetch
	fi
	@if ! git show-ref --quiet --verify "refs/tags/${GO_MK_REF}"; then
		git fetch --tags
	fi
	git checkout --quiet ${GO_MK_REF}
## Project

PACKAGE := github.com/exoscale/egoscale
PROJECT_URL := https://$(PACKAGE)

# Dependencies

go.mk/init.mk:
include go.mk/init.mk
go.mk/public.mk:
include go.mk/public.mk

GO_TEST_EXTRA_ARGS := -mod=readonly -v
GOLANGCI_EXTRA_ARGS := --modules-download-mode=readonly

# OpenAPI code generator
# REF: https://github.com/deepmap/oapi-codegen

OAPI_CODEGEN_VERSION := v1.9.1

OAPI_CODEGEN_MOD_VERSION := $(shell sed -nE 's|^\s*github.com/deepmap/oapi-codegen\s+(v[.0-9]+)$$|\1|p' go.mod)
ifneq ($(OAPI_CODEGEN_VERSION), $(OAPI_CODEGEN_MOD_VERSION))
$(warning OpenAPI code generator (oapi-codegen) versions mismatch (Makefile: $(OAPI_CODEGEN_VERSION); go.mod: $(OAPI_CODEGEN_MOD_VERSION)))
endif


## Targets

# Dependencies
.PHONY: install-oapi-codegen
install-oapi-codegen:
	go install github.com/deepmap/oapi-codegen/cmd/oapi-codegen@$(OAPI_CODEGEN_VERSION)

# OpenAPI specifications (JSON)
.PHONY: oapigen
oapigen: install-oapi-codegen
	@wget -q --show-progress --progress=dot https://openapi-v2.exoscale.com/source.json -O- > v2/oapi/source.json
	@cd v2/oapi/
	@go generate
	@rm source.json
	@ls -l oapi.gen.go

.PHONY: generate
generate:
	@cd v3/generator/
	@go generate
	@go mod tidy
	@go mod vendor
	@cd ..
	@ls -l *.go
