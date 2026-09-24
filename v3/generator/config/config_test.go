package config

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestParse(t *testing.T) {
	c, err := Parse([]byte(`
acronyms:
  sks: SKS
schema-overrides:
  Foo:
    props:
      a-target: a
    refs:
      "#/components/schemas/a-ref": ATarget
aliases:
  - name: ATarget
    target: ARef
`))
	require.NoError(t, err)
	require.Equal(t, "SKS", c.Acronyms["sks"])
	require.Equal(t, "a", c.PropName("Foo", "a-target"))
	require.Equal(t, "b", c.PropName("Foo", "b"))
	require.Equal(t, "b", c.PropName("Bar", "b"))
	typ, ok := c.RefType("Foo", "#/components/schemas/a-ref")
	require.True(t, ok)
	require.Equal(t, "ATarget", typ)
	_, ok = c.RefType("Bar", "#/components/schemas/a-ref")
	require.False(t, ok)
	require.Equal(t, []Alias{{Name: "ATarget", Target: "ARef"}}, c.Aliases)
}

func TestParseRejectsUnknownFields(t *testing.T) {
	_, err := Parse([]byte("acronym:\n  sks: SKS\n"))
	require.Error(t, err)
}
