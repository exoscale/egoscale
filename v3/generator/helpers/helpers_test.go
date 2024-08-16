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

	// trimming
	require.Equal(t, "Test", ToCamel("---test"))
    require.Equal(t, "Test", ToCamel("test___"))
    require.Equal(t, "Test", ToCamel(".-_test_-."))
    require.Equal(t, "Test", ToCamel("  test  "))
    require.Equal(t, "Test", ToCamel(" .-_test_-. "))

	// splitting
    require.Equal(t, "TestFooBarBaz", ToCamel("test-foo_bar.baz"))
    require.Equal(t, "TestFooBar", ToCamel("test--_foo..bar"))
    require.Equal(t, "TESTFooBAR", ToCamel("TEST-Foo-BAR"))
    require.Equal(t, "Test123Foo", ToCamel("test-123-foo"))
    require.Equal(t, "123", ToCamel("1-2-3"))

	// trimming and splitting
	require.Equal(t, "TrimMixedSeparators", ToCamel("-_. trim-mixed/separators ._-"))

	// Acronym handling
    require.Equal(t, "HandleXMLHTTPRequest", ToCamel("handle-xml-http-request"))
    require.Equal(t, "ParseJSONAPI", ToCamel("parse_json_api"))
    require.Equal(t, "NewURLForID", ToCamel("new-url-for-id"))
    require.Equal(t, "SupportSQLQueries", ToCamel("support-sql-queries"))

}

func TestToLowerCamel(t *testing.T) {
	require.Equal(t, "instancePoolID", ToLowerCamel("instance-pool-id"))
	require.Equal(t, "id", ToLowerCamel("id"))
}
