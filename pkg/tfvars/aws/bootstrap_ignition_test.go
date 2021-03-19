package aws

import (
	"fmt"
	"testing"

	igntypes "github.com/coreos/ignition/v2/config/v3_2/types"
)

func Test_parseCertificateBundle(t *testing.T) {
	type checkFunc func([]igntypes.Resource, error) error
	check := func(fns ...checkFunc) []checkFunc { return fns }

	hasNilError := func(_ []igntypes.Resource, have error) error {
		if have != nil {
			return fmt.Errorf("expected nil error, found %q", have)
		}
		return nil
	}

	hasSomeError := func(_ []igntypes.Resource, have error) error {
		if have == nil {
			return fmt.Errorf("expected error, found nil")
		}
		return nil
	}

	chainHasLength := func(want int) checkFunc {
		return func(caref []igntypes.Resource, _ error) error {
			if have := len(caref); want != have {
				return fmt.Errorf("expected bundle length %d, found %d", want, have)
			}
			return nil
		}
	}

	for _, tc := range [...]struct {
		name   string
		bundle string
		checks []checkFunc
	}{{
		"no certificate",
		``,
		check(
			hasNilError,
			chainHasLength(0),
		),
	}, {
		"no certificate, EOLs",
		`




			`,
		check(
			hasNilError,
			chainHasLength(0),
		),
	}, {
		"one certificate",
		`-----BEGIN CERTIFICATE-----
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
11Bq4Ry36CAT3PoIXyhO78iAL6Aly+gheyFT
-----END CERTIFICATE-----`,
		check(
			hasNilError,
			chainHasLength(1),
		),
	}, {
		"one certificate with additional EOL",
		`-----BEGIN CERTIFICATE-----
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
11Bq4Ry36CAT3PoIXyhO78iAL6Aly+gheyFT
-----END CERTIFICATE-----

`,
		check(
			hasNilError,
			chainHasLength(1),
		),
	}, {
		"two certificates",
		`-----BEGIN CERTIFICATE-----
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
11Bq4Ry36CAT3PoIXyhO78iAL6Aly+gheyFT
-----END CERTIFICATE-----
-----BEGIN CERTIFICATE-----
MIIF5zCCA8+gAwIBAgIJAJGoemeMXvTiMA0GCSqGSIb3DQEBCwUAMIGIMQswCQYD
VQQGEwJVUzEXMBUGA1UECAwOTm9ydGggQ2Fyb2xpbmExDTALBgNVBAcMBEFwZXgx
EDAOBgNVBAoMB1JlZCBIYXQxEzARBgNVBAsMClRpZ2VyIFRlYW0xFDASBgNVBAMM
CzE3Mi4zMS44LjEwMRQwEgYDVQQDDAsxNzIuMzEuMC4xMDAgFw0xNzA4MTExNTQy
NDlaGA8yMjE3MDYyNDE1NDI0OVowgYgxCzAJBgNVBAYTAlVTMRcwFQYDVQQIDA5O
b3J0aCBDYXJvbGluYTENMAsGA1UEBwwEQXBleDEQMA4GA1UECgwHUmVkIEhhdDET
MBEGA1UECwwKVGlnZXIgVGVhbTEUMBIGA1UEAwwLMTcyLjMxLjguMTAxFDASBgNV
BAMMCzE3Mi4zMS4wLjEwMIICIjANBgkqhkiG9w0BAQEFAAOCAg8AMIICCgKCAgEA
wt4NNVOMvDd/Tr4HgGZb/43ZWvbrR3TDpSmlWc+MKk6UkBZ4Zy4lWie5xmlpb8Kp
JAbCFGcfJ/zpf433BvYnl0Em6EDwL6I7YCxwSytBr3TeyYT/JMbeHl6z/QKrxG1a
JLLrDcSBQjzCkaA5TMdMoUKQyICRWydgVGo39sAqYKhPII0Z6OoM6olrFwZNuN+T
4CA8uKSFcgDRsNf7iFjklSIsy6DfmAcOM0CYNwVtWYKG/SYVMStkmZvyZeZgeKzB
gy0SAqHxJ16tZghhLdS4edl8mNaVtpbBMcZU3Z9WqmP2IwR1MoDpebD2ZFnrY1k+
efgTyFZklIUjwmvm9YO1nX9epRsVwozH6YV4lJeRpZBVb7XzPeWm58xSvVFHiPJx
VZhFQb9PFv8fhihjAaaYf3AqhILiBlUH1Qb2pCdKpvQGetTEDD9jLAJHkpF72nfy
t3xOoCJGqhNbiUMvEKvpUlO6rTosITEkmtaRdGC1IwaBETTmBQDwS4nBg7hOlGO2
t7XrUXxrGbuGXyCY7Zz1HVVObH/m4I0yPpGf0BCZl1PbxVXeN+jOg29XX01mhG6x
7bEWsUTCfuokvJHRWUrkwZUpcgwBMFNHqOeq78NcUGlEpkbdoptWHpjxe+w6xZsz
AGJaFcvUU1EohoFdzsnF7ra7l4TEZZyFEdwlYXZnj8UCAwEAAaNQME4wHQYDVR0O
BBYEFBpbryhpHb6rSdaYsn1x7iE/KNy6MB8GA1UdIwQYMBaAFBpbryhpHb6rSdaY
sn1x7iE/KNy6MAwGA1UdEwQFMAMBAf8wDQYJKoZIhvcNAQELBQADggIBAECgEP0w
AP9BG3b7xXJdPV5oHkOIDmDcZtlSChLBQtyPwHez+8VcBxJqGF39FrYTqDaVKE+W
u1EV8oiC/JKjT9StTWNUKHrzrzNfVP378D8PShJOnlasTwIjbp4XycTV6uju8NDb
XKq8Lsd5XkpX+syhfvSw8KyuRP8nzrmUh7usokqXoMuohk94fp0nqFA2lo+16Xj5
zmZX+Ukyyho5Uh7Rpcx0W/ergsxsyfGt3yBVQO/ohsau9ezyBJwucY8vCLQ8Vd1i
OQao4dnOg4eIV8g5jPhll4nEgZqOp9NLmnoXWqOVPZ4FQ6V19nd21pRTvHmTXDFG
HO33EZdRAsSFr6f5cnHqqrv5PoTZQtyiXe01FUE9dpFS+gQf4cd98tr/Z39AYNZ7
5cllfZkM0xVuhR5fOzGlbDfVQD5eSqwM0Qp5wwMwJhGkxCLe/6l6xV0TDV5Rcb7i
BQ3zPGw7uJtOrxeFnGHWNWYdl24gxA1h7G3uKeKo/MNiJHvOLNul9wuvlWCCEJI0
bxr7iKAGSun8qkAWkYQy6ClARIGSiy+ah6pwOgObbge6eR2o/QOYFhIE1yleTNJU
JqYys+7rMa2qN2mfsnkn/IN2G80+Qs+qTLv2LDucO2PweBaslWh2ejc6axzLBJza
11Bq4Ry36CAT3PoIXyhO78iAL6Aly+gheyFT
-----END CERTIFICATE-----`,
		check(
			hasNilError,
			chainHasLength(2),
		),
	}, {
		"two certificates with many EOLs",
		`

-----BEGIN CERTIFICATE-----
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
11Bq4Ry36CAT3PoIXyhO78iAL6Aly+gheyFT
-----END CERTIFICATE-----




-----BEGIN CERTIFICATE-----
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
11Bq4Ry36CAT3PoIXyhO78iAL6Aly+gheyFT
-----END CERTIFICATE-----

`,
		check(
			hasNilError,
			chainHasLength(2),
		),
	}, {
		"invalid certificate",
		`-----BEGIN CERTIFICATE-----
invalid
-----END CERTIFICATE-----`,
		check(
			hasSomeError,
		),
	}} {
		t.Run(tc.name, func(t *testing.T) {
			caref, err := parseCertificateBundle([]byte(tc.bundle))
			for _, check := range tc.checks {
				if e := check(caref, err); e != nil {
					t.Error(e)
				}
			}
		})
	}
}
