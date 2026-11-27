package build

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/exoscale/egoscale/v3/generator/config"
	"github.com/exoscale/egoscale/v3/generator/ir"
)

const problemsSpec = `openapi: 3.0.0
info: {title: t, version: "1"}
paths: {}
x-problems:
  types:
    not-found: {status: 404, title: The resource is not found.}
    kms/key-not-found: {parent: not-found, status: 404, title: Key not found, description: The key does not exist.}
`

func TestProblems(t *testing.T) {
	model, err := Load([]byte(problemsSpec))
	require.NoError(t, err)

	f, err := New(model, config.Config{}).Problems("v3")
	require.NoError(t, err)
	// The shared types are left out
	require.Equal(t, []ir.Decl{ir.ProblemTypes{Types: []ir.ProblemType{{
		Doc:    "// ErrKmsKeyNotFound (kms/key-not-found, 404): The key does not exist.",
		Name:   "ErrKmsKeyNotFound",
		ID:     "#/problems/kms/key-not-found",
		Status: 404,
		Title:  "Key not found",
	}}}}, f.Decls)
}

func TestProblemsRefusesTheSameGoName(t *testing.T) {
	model, err := Load([]byte(`openapi: 3.0.0
info: {title: t, version: "1"}
paths: {}
x-problems:
  types:
    kms/key-gone: {parent: not-found, status: 404, title: Gone}
    kms/key.gone: {parent: not-found, status: 404, title: Gone}
`))
	require.NoError(t, err)

	_, err = New(model, config.Config{}).Problems("v3")
	require.ErrorContains(t, err, "ErrKmsKeyGone")
}
