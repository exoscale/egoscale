# Egoscale v3

Exoscale API Golang wrapper

**Egoscale v3** is based on a generator written from scratch with [libopenapi](https://github.com/pb33f/libopenapi).

The core base of the generator is using libopenapi to parse and read the [Exoscale OpenAPI spec](https://api-ch-gva-2.exoscale.com/v2/openapi.json) and then generate the code from it.

## Installation

Install the following dependencies:

```shell
go get "github.com/exoscale/egoscale/v3"
```

Add the following import:

```golang
import "github.com/exoscale/egoscale/v3"
```
## Examples

```Golang
package main

import (
	"context"
	"log"

	"github.com/davecgh/go-spew/spew"
	v3 "github.com/exoscale/egoscale/v3"
	"github.com/exoscale/egoscale/v3/credentials"
)

func main() {
	creds := credentials.NewEnvCredentials()
	// OR
	creds = credentials.NewStaticCredentials("EXOxxx..", "...")

	client, err := v3.NewClient(creds)
	if err != nil {
		log.Fatal(err)
	}

	ctx := context.Background()
	op, err := client.CreateInstance(ctx, v3.CreateInstanceRequest{
		Name:     "egoscale-v3",
		DiskSize: 50,
		// Ubuntu 24.04 LTS
		Template: &v3.Template{ID: v3.UUID("cbd89eb1-c66c-4637-9483-904d7e36c318")},
		// Medium type
		InstanceType: &v3.InstanceType{ID: v3.UUID("b6e9d1e8-89fc-4db3-aaa4-9b4c5b1d0844")},
	})
	if err != nil {
		log.Fatal(err)
	}

	op, err = client.Wait(ctx, op, v3.OperationStateSuccess)
	if err != nil {
		log.Fatal(err)
	}

	instance, err := client.GetInstance(ctx, op.Reference.ID)
	if err != nil {
		log.Fatal(err)
	}

	spew.Dump(instance)
}
```

### Findable

Most of the list request `ListX()` return a type containing the list of the resource requested and a method `FindX()` to be able to retrieve a resource by its `name` or `id` most of the time.

```Golang
pools, err := client.ListInstancePools(ctx)
if err != nil {
	log.Fatal(err)
}
pool, err := pools.FindInstancePool("my-pool-example")
if err != nil {
	log.Fatal(err)
}

fmt.Println(pool.Name)
```

## Development

### Generate Egoscale v3

From the root repo
```Bash
make pull-oapi-spec # Optional(to pull latest Exoscale Open-API spec)
make generate
```

### Debug generator output

```Bash
mkdir test
GENERATOR_DEBUG=client make generate > test/client.go
GENERATOR_DEBUG=schemas make generate > test/schemas.go
GENERATOR_DEBUG=operations make generate > test/operations.go
```

### OpenAPI Extensions

The generator support two types of extension:
- `x-go-type` to specify a type definition in Golang.

	OpenAPI Spec
	```yaml
	api-endpoint:
	  type: string
	  x-go-type: Endpoint
	  description: Zone API endpoint
	```
	Generated code
	```Golang
	type Endpoint string

	type Zone struct {
		APIEndpoint Endpoint // Here is the generated type definition.
		...
	}
	```
- `x-go-findable` to specify which fields in the findable resource to fetch
	OpenAPI Spec
	```yaml
	elastic-ip:
      type: object
      properties:
        id:
          type: string
		  x-go-findable: "1"
          description: Elastic IP ID
        ip:
          type: string
		  x-go-findable: "2"
          description: Elastic IP address
	```
	Generated code
	```Golang
	// FindElasticIP attempts to find an ElasticIP by idOrIP.
	func (l ListElasticIPSResponse) FindElasticIP(idOrIP string) (ElasticIP, error) {
		for i, elem := range l.ElasticIPS {
			if string(elem.ID) == idOrIP || string(elem.IP) == idOrIP {
				return l.ElasticIPS[i], nil
			}
		}

		return ElasticIP{}, fmt.Errorf("%q not found in ListElasticIPSResponse: %w", idOrIP, ErrNotFound)
	}
	```

## Generator Overrides System

The Egoscale v3 generator incorporates an overrides system to preserve backwards compatibility in the Go API when the OpenAPI specification changes, such as renaming schemas or references. This prevents breaking existing code that relies on specific type names, field names, or JSON tags.

The system consists of several components working together during code generation.

### Reference overrides

Defined in `generator/helpers/helpers.go`: they intercept OpenAPI `$ref` paths and redirect them to custom Go type names, bypassing the default camel-case conversion. For example, a schema referencing `"#/components/schemas/ssh-key-ref"` will generate code using the type `SSHKey` instead of `SSHKeyRef`.

### Property overrides

Also in `generator/helpers/helpers.go`: they adjust property names within schemas to maintain historical field names and JSON tags in the generated Go structs. This ensures that changes to property names in the OpenAPI spec do not alter the API surface. For instance, a property named `"block-storage-volume-ref"` can be overridden to generate a field named `BlockStorageVolume` with a JSON tag `"block-storage-volume"`.

### Backwards compatibility aliases

Configured in `generator/schemas/schemas.go`: they map updated schema names back to their original names, producing `type Old = New` declarations in the output. This allows legacy type names to remain valid even after spec updates. An alias like `"ssh-key": "ssh-key-ref"` generates `type SSHKey = SSHKeyRef`, letting code continue using `SSHKey` while the spec defines `ssh-key-ref`.

For cases requiring manual intervention, special aliases are hardcoded in the `Generate` function of `generator/schemas/schemas.go`. These add specific type equivalences, such as `type InstanceTarget = Instance` or `type BlockStorageSnapshotTarget = BlockStorageSnapshot`, to support operations needing distinct field names.

During generation, the system applies these overrides sequentially: first resolving references, then adjusting properties for struct fields, and finally appending aliases. To modify overrides, edit the relevant files and regenerate with `make generate`, followed by `go build ./...` to verify. Debugging output can be obtained using `GENERATOR_DEBUG=schemas make generate`.

This approach enables the SDK to track OpenAPI spec evolution without disrupting user-facing APIs.
