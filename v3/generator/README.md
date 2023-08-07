# Exoscale Consumer API code generator (V3)

This is a CLI tool that generates Consumer API on top of [oapi-codegen](https://github.com/deepmap/oapi-codegen) generated code.
Main motivation is grouping API functions into namespaces (thus shortening names) and implementing boilerplate response parsing.

## Commands

Tool supports 2 commands: `generate` and `list-unimplemented`.

### `generate`

Main command that generates the Consumer API according to the mapping from `mapping.go` file.
Mapping defines following parameters:

- subfolder (one of: compute, dbaas, dns, iam and global);
- struct name (entity) for a group of API calls (eg. `AccessKey`), also used as a file name (in snake-case, eg. `api/iam/access_key.gen.go`);
- function name in Consumer API (eg. `List`);
- function name in `oapi/oapi.gen.go` (common part from `ClientWithResponsesInterface` interface, eg. `ListAccessKeys`).

Optional overrides (per function) are:

- `ResDefOverride`: overrides function response definition;
- `ResPassthroughOverride:`: overrides function response return argument;

### `list-unimplemented`

Helper command `list-unimplemented` can be used to print all functions from `oapi/oapi.gen.go` that are not yet implemented in Consumer API.

## Adding new functions

Simple workflow for adding new functions in Consumer API looks like this:

- inside generator folder run `go run . list-unimplemented` to find names of functions you want to add;
- update `mapping.go` to include the new functions (use name from previous command as `OAPIName:`);
- run `go run . generate`;

If new entity (with `RootName`) was added, update main subfolder file (eg. `api/iam/iam.go`) to attach the new entity.

For complex nested structures in response some additional work is needed:

- add nested struct definition in `oapi/oapi_nested_structs.go` (see existing implementations);
- override template response using `ResDefOverride` and `ResPassthroughOverride:` (check for example `dbaas/Integrations/ListSettings` in `mapping.go`).
