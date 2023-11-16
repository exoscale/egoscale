## Project

PACKAGE := github.com/exoscale/egoscale
PROJECT_URL := https://$(PACKAGE)

# Dependencies

# Requires: https://github.com/exoscale/go.mk
# - install: git submodule update --init --recursive go.mk
# - update:  git submodule update --remote
include go.mk/init.mk
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
	wget -q --show-progress --progress=dot https://openapi-v2.exoscale.com/source.json -O- > v2/oapi/source.json
	@echo
	cd v2/oapi/; go generate
	@rm v2/oapi/source.json
	ls -l v2/oapi/oapi.gen.go

.PHONY: generate
generate:
	@go install go.uber.org/mock/mockgen@latest
	@wget -q --show-progress --progress=dot https://openapi-v2.exoscale.com/source.yaml -O- > v3/generator/source.yaml
	@echo
	@cd v3/generator/; go generate
	@go mod tidy && go mod vendor
	@rm v3/generator/source.yaml
	@ls -l v3/*.go
