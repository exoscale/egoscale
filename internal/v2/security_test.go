package v2

import (
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestNewSecurityProviderExoscaleV2(t *testing.T) {
	var (
		provider *SecurityProviderExoscaleV2
		err      error
	)

	provider, err = NewSecurityProviderExoscaleV2("key", "")
	require.NotNil(t, err, "expected an error")
	require.Nil(t, provider)

	provider, err = NewSecurityProviderExoscaleV2("", "secret")
	require.NotNil(t, err, "expected an error")
	require.Nil(t, provider)

	provider, err = NewSecurityProviderExoscaleV2("key", "secret")
	require.Nil(t, err)
	require.NotNil(t, provider)
	require.Equal(t, "key", provider.apiKey)
	require.Equal(t, "secret", provider.apiSecret)
}

func TestSecurityProviderExoscaleV2SignRequest(t *testing.T) {
	// In order to test the signing process validation, we have to compute expected signatures using an external
	// (verified) implementation with the same properties and compare them to the output of the signRequest()
	// method, e.g. https://github.com/exoscale/requests-exoscale-auth

	var (
		testAPIKey     = "EXOxxxxxxxxxxxxxxxxxxxxxxxx"
		testAPISecret  = "XXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXX"
		testExpireDate = time.Date(2077, 1, 1, 0, 0, 0, 0, time.UTC)
	)

	provider := &SecurityProviderExoscaleV2{
		apiKey:    testAPIKey,
		apiSecret: testAPISecret,
	}

	// Request without URL parameters
	req, err := http.NewRequest("GET", "https://api.exoscale.com/v2/zone", nil)
	require.NoError(t, err)
	require.NoError(t, provider.signRequest(req, testExpireDate))
	require.Equal(t,
		"EXO2-HMAC-SHA256 "+
			"credential="+testAPIKey+
			",expires="+fmt.Sprint(testExpireDate.Unix())+
			",signature=Ntbq/p0HVmA3Zg1HHY+Lq1vjFGi7HeMrrgXDS5jRNlY=",
		req.Header.Get("Authorization"))

	// Request with URL parameters
	req, err = http.NewRequest("GET", "https://api.exoscale.com/v2/zone?k1=v1&k2=v2", nil)
	require.NoError(t, err)
	require.NoError(t, provider.signRequest(req, testExpireDate))
	require.Equal(t,
		"EXO2-HMAC-SHA256 "+
			"credential="+testAPIKey+
			",signed-query-args=k1;k2"+
			",expires="+fmt.Sprint(testExpireDate.Unix())+
			",signature=iqOBz13+44L5j0uJclE8hmUhQQcvtCSoPEOXYK6liqY=",
		req.Header.Get("Authorization"))
}
