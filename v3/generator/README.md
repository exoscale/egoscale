# Exoscale Consumer API code generator (V3)

This is a CLI tool that generates API on top of [oapi-codegen](https://github.com/deepmap/oapi-codegen) generated code.

## Commands

Tool supports 2 commands: `generate` and `list-unimplemented`.

### `generate`

Main command that generates the API according to the mapping from `mapping.go` file.
Mapping defines following parameters:

- group (one of: compute, dbaas, dns, iam and global);
- struct name (resource) for a group of API calls (eg. `AccessKey`), also used as a file name (in snake-case, eg. `iam_access_key.gen.go`);
- function name (eg. `List`);
- function name in `client.gen.go` (common part, eg. `ListAccessKeys`).

Optional overrides (per function) are:

- `ResDefOverride`: overrides function response definition;
- `ResPassthroughOverride:`: overrides function response return argument;

### `list-unimplemented`

Helper command `list-unimplemented` can be used to print all functions from `client.gen.go` that are not yet implemented in API.

## Adding new functions

Simple workflow for adding new functions in Consumer API looks like this:

- inside generator folder run `go run . list-unimplemented` to find names of functions you want to add;
- update `mapping.go` to include the new functions (use name from previous command as `OAPIName:`);
- run `go run . generate`;

If new entity (with `RootName`) was added, update main group file (eg. `iam.go`) to attach the new resource.

For complex nested structures in response some additional work is needed:

- add nested struct definition in `types_nested_structs.go` (see existing implementations);
- override template response using `ResDefOverride` and `ResPassthroughOverride:` (check for example `dbaas/Integrations/ListSettings` in `mapping.go`).
