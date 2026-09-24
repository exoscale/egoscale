package naming

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCamel(t *testing.T) {
	n := New(nil)
	require.Equal(t, "InstancePoolID", n.Camel("instance-pool-id"))
	require.Equal(t, "InstancePoolID", n.Camel("-instance-pool-id-"))
	require.Equal(t, "ID", n.Camel("id"))
	require.Equal(t, "Test", n.Camel("test"))
	require.Equal(t, "", n.Camel("."))
	require.Equal(t, "TestFoo", n.Camel("test...foo"))

	// trimming
	require.Equal(t, "Test", n.Camel("---test"))
	require.Equal(t, "Test", n.Camel("test___"))
	require.Equal(t, "Test", n.Camel(".-_test_-."))
	require.Equal(t, "Test", n.Camel("  test  "))
	require.Equal(t, "Test", n.Camel(" .-_test_-. "))

	// splitting
	require.Equal(t, "TestFooBarBaz", n.Camel("test-foo_bar.baz"))
	require.Equal(t, "TestFooBar", n.Camel("test--_foo..bar"))
	require.Equal(t, "TESTFooBAR", n.Camel("TEST-Foo-BAR"))
	require.Equal(t, "Test123Foo", n.Camel("test-123-foo"))
	require.Equal(t, "123", n.Camel("1-2-3"))

	// trimming and splitting
	require.Equal(t, "TrimMixedSeparators", n.Camel("-_. trim-mixed/separators ._-"))

	// Acronym handling
	require.Equal(t, "HandleXMLHTTPRequest", n.Camel("handle-xml-http-request"))
	require.Equal(t, "ParseJSONAPI", n.Camel("parse_json_api"))
	require.Equal(t, "NewURLForID", n.Camel("new-url-for-id"))
	require.Equal(t, "SupportSQLQueries", n.Camel("support-sql-queries"))

}

func TestLowerCamel(t *testing.T) {
	n := New(nil)
	require.Equal(t, "instancePoolID", n.LowerCamel("instance-pool-id"))
	require.Equal(t, "id", n.LowerCamel("id"))
}

func TestCamelAcronyms(t *testing.T) {
	n := New(map[string]string{"dbaas": "DBAAS", "ai": "AI"})
	// From github.com/BluntSporks/abbreviation.
	require.Equal(t, "ListXMLFoo", n.Camel("list-xml-foo"))
	// From the configured acronyms.
	require.Equal(t, "CreateDBAASService", n.Camel("create-dbaas-service"))
	require.Equal(t, "ListAIAPIKeys", n.Camel("list-ai-api-keys"))
	require.Equal(t, "dbaasService", n.LowerCamel("dbaas-service"))
	// Not configured: kept as a regular word.
	require.Equal(t, "Dbaas", New(nil).Camel("dbaas"))
}

func TestDoc(t *testing.T) {
	require.Equal(t, "", Doc(""))
	require.Equal(t, "", Doc("null"))
	require.Equal(t, "// a", Doc("a"))
	require.Equal(t, "// a\n// b", Doc(" a \n\nb\n"))
	require.Equal(t, "// a\n// b\n// c", Doc("a\n\n\nb\nc"))
}
