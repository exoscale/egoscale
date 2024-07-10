module github.com/exoscale/egoscale

require (
	github.com/deepmap/oapi-codegen v1.9.1
	github.com/gofrs/uuid v4.4.0+incompatible
	github.com/hashicorp/go-retryablehttp v0.7.7
	github.com/stretchr/testify v1.8.4
)

require (
	github.com/kr/pretty v0.3.1 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)

require (
	github.com/davecgh/go-spew v1.1.2-0.20180830191138-d8f796af33cc // indirect
	github.com/hashicorp/go-cleanhttp v0.5.2 // indirect
	github.com/pmezard/go-difflib v1.0.1-0.20181226105442-5d4384ee4fb2 // indirect
	github.com/stretchr/objx v0.5.0 // indirect
)

go 1.22

retract (
	v1.19.1 // Retracts the previous version
	v1.19.0 // Published accidentally.
)
