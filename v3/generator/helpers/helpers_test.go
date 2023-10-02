package helpers

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestToCamel(t *testing.T) {
	require.Equal(t, "InstancePoolID", ToCamel("instance-pool-id"))
	require.Equal(t, "InstancePoolID", ToCamel("-instance-pool-id-"))
	require.Equal(t, "ID", ToCamel("id"))
	require.Equal(t, "Test", ToCamel("test"))
	require.Equal(t, "", ToCamel("."))
	require.Equal(t, "TestFoo", ToCamel("test...foo"))
}

func TestToLowerCamel(t *testing.T) {
	require.Equal(t, "instancePoolID", ToLowerCamel("instance-pool-id"))
	require.Equal(t, "id", ToLowerCamel("id"))
}
