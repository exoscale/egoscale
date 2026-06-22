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

func TestEnumValueName(t *testing.T) {
	require.Equal(t, "", EnumValueName(""))

	require.Equal(t, "Standby", EnumValueName("standby"))
	require.Equal(t, "Master", EnumValueName("master"))

	require.Equal(t, "ReadMinusReplica", EnumValueName("read-replica"))
	require.Equal(t, "WritePlusReplica", EnumValueName("write+replica"))

	require.Equal(t, "1gDot24gb", EnumValueName("1g.24gb"))
	require.Equal(t, "1gDot24gbMinusMe", EnumValueName("1g.24gb-me"))
	require.Equal(t, "1gDot24gbPlusMe", EnumValueName("1g.24gb+me"))
	require.Equal(t, "2gDot48gbPlusMeDotAll", EnumValueName("2g.48gb+me.all"))
	require.Equal(t, "4gDot96gbPlusGfx", EnumValueName("4g.96gb+gfx"))
	require.Equal(t, "1gDot24gbPlusMeDotAll", EnumValueName("1g.24gb+me.all"))

	seen := map[string]string{
		EnumValueName("1g.24gb-me"):     "1g.24gb-me",
		EnumValueName("1g.24gb+me"):     "1g.24gb+me",
		EnumValueName("2g.48gb-me"):     "2g.48gb-me",
		EnumValueName("2g.48gb"):        "2g.48gb",
		EnumValueName("4g.96gb+gfx"):    "4g.96gb+gfx",
		EnumValueName("1g.24gb"):        "1g.24gb",
		EnumValueName("2g.48gb+me.all"): "2g.48gb+me.all",
		EnumValueName("1g.24gb+gfx"):    "1g.24gb+gfx",
		EnumValueName("1g.24gb+me.all"): "1g.24gb+me.all",
		EnumValueName("4g.96gb"):        "4g.96gb",
		EnumValueName("2g.48gb+gfx"):    "2g.48gb+gfx",
	}
	require.Len(t, seen, 11, "every enum value must produce a distinct identifier")

	require.Equal(t, "HandleMinusXMLMinusHTTPMinusRequest", EnumValueName("handle-xml-http-request"))
	require.Equal(t, "ParseJSONAPI", EnumValueName("parse_json_api"))
	require.Equal(t, "SupportMinusSQLMinusQueries", EnumValueName("support-sql-queries"))

	require.Equal(t, "Trim", EnumValueName(" ---trim--- "))
	require.Equal(t, "Trim", EnumValueName("---trim---"))
	require.Equal(t, "", EnumValueName("./"))
}
