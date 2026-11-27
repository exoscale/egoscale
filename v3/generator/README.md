# Egoscale v3 generator

The generator writes the Egoscale v3 SDK from the Exoscale OpenAPI spec (`source.yaml`). It is its own
Go module: it isn't part of the SDK users import.

It generates four files of the `v3` package:

- `schemas.go`: a Go type per component schema of the spec, and the backward compatibility aliases;
- `operations.go`: a `Client` method per API operation, with its inline request, response and
  parameter types, and the `Find` methods of list responses;
- `constants.go`: the zone API endpoints;
- `problems.go`: an error value per service error type of the `x-problems` of the spec, e.g.
  `ErrKmsKeyNotFound`, which `errors.Is` matches against the errors of the client.

The rest of the `v3` package is hand-written: the client (`client.go`, `api.go`), the errors
(`errors.go`), `version.go`, `doc.go` and the `credentials`, `metadata` and `utils` packages.
Never edit the generated files: change the spec, `config.yaml` or the generator, and regenerate.

## Regenerate the SDK

From the root of the repository:

```shell
make pull-oapi-spec     # optional: fetch the latest spec into source.yaml (make pull-oapi-spec-pp for preprod)
make generate           # go generate, then go mod tidy and go mod vendor
```

`go generate` runs `go run . ./source.yaml ../ v3` in this directory: the arguments are the spec,
the output directory and the package name.

The `check_and_regenerate_v3` workflow does the same every day, and opens a PR when the SDK
changes. `regenerate_v3_preprod`, run by hand, pushes the SDK generated from the preprod spec to the
`pp` branch.

## Layout

Generation has two stages, joined by an intermediate representation:

- `build/` reads the OpenAPI model and builds the Go declarations of each file. All the
  interpretation of the spec lives here.
- `ir/` holds those declarations (`ir.File`): Go types, struct fields, enums, operations… It doesn't
  depend on libopenapi.
- `render/` turns an `ir.File` into formatted Go source, with the templates of the package.
- `naming/` turns spec names into Go names, with acronyms (`sks-nodepool` → `SKSNodepool`).
- `config/` and `config.yaml` hold the settings that are not part of the spec.
- `main.go` builds and renders the three files, and writes them only if all of them succeeded.

The spec is OpenAPI 3.0: `build.Load` turns off the OpenAPI 3.1 handling of properties next to a
`$ref` in libopenapi, which would otherwise turn hundreds of references into `any`. Always load a
spec with `build.Load`, tests included.

Generation fails, instead of writing an SDK that doesn't compile or silently changes, when:

- two declarations have the same Go type name (e.g. a component schema and an inline schema);
- a schema has a `format` or a type the generator doesn't know;
- a component enum has values that can't be Go constant names;
- a reference override of `config.yaml` doesn't match any reference of the spec.

## Configuration: `config.yaml`

- `acronyms`: Exoscale specific acronyms (`sks: SKS`). They take precedence over the generic list of
  [abbreviation](https://github.com/BluntSporks/abbreviation).
- `ignored-schemas`: component schemas not generated, e.g. when one would get the same Go name as
  another type.
- `skip-auth-operations`: operation IDs whose requests are sent without credentials.
- `ref-overrides`: per Go type name, the Go type to use for a reference of the spec instead of the
  type of the referenced schema. They keep the SDK backward compatible when the spec changes a
  reference, e.g. `CreateInstanceRequest` keeps a `*Template` field when the spec now refers to
  `template-ref`.
- `aliases`: type aliases kept for backward compatibility (`type InstanceTarget = InstanceRef`).

## OpenAPI extensions

- `x-go-type` gives the Go type of a simple schema, instead of the one of its type and format. The type
  must be declared in the package: `Endpoint` is declared in `constants.go`.

  ```yaml
  api-endpoint:
    type: string
    x-go-type: Endpoint
  ```

  generates the field `APIEndpoint Endpoint`.

- `x-go-enum-varnames` gives the suffixes of the enum constant names, in the order of the enum values,
  when the values don't make distinct Go names (`1g.24gb+me` and `1g.24gb-me` both give `1g24gbMe`).
  It is ignored if it doesn't have as many names as the enum has values.

  ```yaml
  nvidia-mig-profile-rtxpro6000.96gb:
    type: string
    enum: [1g.24gb, 1g.24gb+me, …, 1g.24gb-me, …]
    x-go-enum-varnames: [1G24Gb, 1G24GbMe, …, 1G24GbNoMe, …]
  ```

  generates `NvidiaMigProfileRtxpro600096gb1G24GbMe NvidiaMigProfileRtxpro600096gb = "1g.24gb+me"`.

- `x-go-findable` (`"1"` or `"2"`) picks the fields a `Find` method matches, see below.

## `Find` methods

The response of a list operation (its operation ID starts with `list`) gets a `Find<Item>` method.
The response must be inline, or a `$ref` to a schema whose name ends with `-response`. The method
looks in the first array property whose items are objects, and matches the `name` and `id` fields of
the items. A property with `x-go-findable: "1"` replaces `name`, and one with `x-go-findable: "2"`
replaces `id`. For example, `zone` has `x-go-findable: "2"` on `api-endpoint`:

```go
func (l ListZonesResponse) FindZone(nameOrAPIEndpoint string) (Zone, error)
```

It returns an error wrapping `ErrNotFound` when no item matches, and `ErrConflict` when several do.

## Tests

`make test` runs the tests of this module. `golden_test.go` generates the SDK from `source.yaml` and
compares it with the committed files: a change of the generator that changes the generated SDK must
regenerate it in the same commit. The tests of `build/` check the declarations built from small
specs, without rendering Go.

To see the generated source of a file before `gofmt`, e.g. when it doesn't compile, set
`GENERATOR_DEBUG` to its name: `schemas`, `operations`, `constants` or `problems`.

```shell
GENERATOR_DEBUG=schemas go generate > /tmp/schemas.go
```
