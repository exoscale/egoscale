package oapi

import (
	"encoding/json"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestDbaasServiceUser_UnmarshalJSON(t *testing.T) {
	var (
		testAccessCert                     = testRandomString(10)
		testAccessCertNotValidAfterTime, _ = time.Parse(iso8601Format, "2020-08-12T11:12:36Z")
		testAccessKey                      = testRandomString(10)
		testAuthentication                 = DbaasServiceUserAuthentication(testRandomString(10))
		testPassword                       = testRandomString(10)
		testType                           = testRandomString(10)
		testUsername                       = testRandomString(10)

		expected = DbaasServiceUser{
			AccessCert:                  &testAccessCert,
			AccessCertNotValidAfterTime: &testAccessCertNotValidAfterTime,
			AccessControl: &DbaasServiceUserAccessControl{
				RedisAclCategories: &[]string{testRandomString(10)},
				RedisAclCommands:   &[]string{testRandomString(10)},
				RedisAclKeys:       &[]string{testRandomString(10)},
			},
			AccessKey:      &testAccessKey,
			Authentication: &testAuthentication,
			Password:       &testPassword,
			Type:           testType,
			Username:       testUsername,
		}

		actual DbaasServiceUser

		jsonDbaasService = `{
  "access-cert": "` + testAccessCert + `",
  "access-cert-not-valid-after-time": "` + testAccessCertNotValidAfterTime.Format(iso8601Format) + `",
  "access-control": {
	"redis-acl-categories": ["` + (*expected.AccessControl.RedisAclCategories)[0] + `"],
	"redis-acl-commands": ["` + (*expected.AccessControl.RedisAclCommands)[0] + `"],
	"redis-acl-keys": ["` + (*expected.AccessControl.RedisAclKeys)[0] + `"]
  },
  "access-key": "` + testAccessKey + `",
  "authentication": "` + fmt.Sprint(testAuthentication) + `",
  "password": "` + testPassword + `",
  "type": "` + testType + `",
  "username": "` + testUsername + `"
}`
	)

	require.NoError(t, json.Unmarshal([]byte(jsonDbaasService), &actual))
	require.Equal(t, expected, actual)
}

func TestDbaasServiceUser_MarshalJSON(t *testing.T) {
	var (
		testAccessCert                     = testRandomString(10)
		testAccessCertNotValidAfterTime, _ = time.Parse(iso8601Format, "2020-08-12T11:12:36Z")
		testAccessKey                      = testRandomString(10)
		testAuthentication                 = DbaasServiceUserAuthentication(testRandomString(10))
		testPassword                       = testRandomString(10)
		testType                           = testRandomString(10)
		testUsername                       = testRandomString(10)

		dbaasServiceUser = DbaasServiceUser{
			AccessCert:                  &testAccessCert,
			AccessCertNotValidAfterTime: &testAccessCertNotValidAfterTime,
			AccessControl: &DbaasServiceUserAccessControl{
				RedisAclCategories: &[]string{testRandomString(10)},
				RedisAclCommands:   &[]string{testRandomString(10)},
				RedisAclKeys:       &[]string{testRandomString(10)},
			},
			AccessKey:      &testAccessKey,
			Authentication: &testAuthentication,
			Password:       &testPassword,
			Type:           testType,
			Username:       testUsername,
		}

		expected = []byte(`{` +
			`"access-cert":"` + testAccessCert + `",` +
			`"access-cert-not-valid-after-time":"` + testAccessCertNotValidAfterTime.Format(iso8601Format) + `",` +
			`"access-control":{` +
			`"redis-acl-categories":["` + (*dbaasServiceUser.AccessControl.RedisAclCategories)[0] + `"],` +
			`"redis-acl-commands":["` + (*dbaasServiceUser.AccessControl.RedisAclCommands)[0] + `"],` +
			`"redis-acl-keys":["` + (*dbaasServiceUser.AccessControl.RedisAclKeys)[0] + `"]` +
			`},` +
			`"access-key":"` + testAccessKey + `",` +
			`"authentication":"` + fmt.Sprint(testAuthentication) + `",` +
			`"password":"` + testPassword + `",` +
			`"type":"` + testType + `",` +
			`"username":"` + testUsername + `"` +
			`}`)
	)

	actual, err := json.Marshal(dbaasServiceUser)
	require.NoError(t, err)
	require.Equal(t, expected, actual)
}
