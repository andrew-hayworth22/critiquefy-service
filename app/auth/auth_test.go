package auth_test

import (
	"bytes"
	"context"
	"fmt"
	"runtime/debug"
	"testing"
	"time"

	"github.com/andrew-hayworth22/critiquefy-service/app/auth"
	"github.com/andrew-hayworth22/critiquefy-service/foundation/logger"
	"github.com/golang-jwt/jwt/v5"
)

func Test_Auth(t *testing.T) {
	log, teardown := newUnit(t)
	defer func() {
		if r := recover(); r != nil {
			t.Log(r)
			t.Error(string(debug.Stack()))
		}
		teardown()
	}()

	cfg := auth.Config{
		Log:       log,
		KeyLookup: &keyStore{},
		Issuer:    "critiquefy",
	}

	a, err := auth.New(cfg)
	if err != nil {
		t.Fatalf("Should be able to create an authenticator: %s", err)
	}

	cases := []struct {
		name                   string
		claims                 auth.Claims
		requiredRole           string
		tokenOverride          string
		expectedJWTGeneration  bool
		expectedAuthentication bool
		expectedAuthorization  bool
	}{
		{
			name: "Success_SingleRole",
			claims: auth.Claims{
				RegisteredClaims: jwt.RegisteredClaims{
					Issuer:    "critiquefy",
					Subject:   "c11eabcc-8492-4dfa-a586-97d9f1694a8a",
					ExpiresAt: jwt.NewNumericDate(time.Now().UTC().Add(time.Hour)),
					IssuedAt:  jwt.NewNumericDate(time.Now().UTC()),
				},
				Roles: []string{"admin"},
			},
			requiredRole:           "admin",
			expectedJWTGeneration:  true,
			expectedAuthentication: true,
			expectedAuthorization:  true,
		},
		{
			name: "Success_MultiRole",
			claims: auth.Claims{
				RegisteredClaims: jwt.RegisteredClaims{
					Issuer:    "critiquefy",
					Subject:   "c11eabcc-8492-4dfa-a586-97d9f1694a8a",
					ExpiresAt: jwt.NewNumericDate(time.Now().UTC().Add(time.Hour)),
					IssuedAt:  jwt.NewNumericDate(time.Now().UTC()),
				},
				Roles: []string{"admin", "user", "chad"},
			},
			requiredRole:           "admin",
			expectedJWTGeneration:  true,
			expectedAuthentication: true,
			expectedAuthorization:  true,
		},
		{
			name: "Fail_NoRoles",
			claims: auth.Claims{
				RegisteredClaims: jwt.RegisteredClaims{
					Issuer:    "critiquefy",
					Subject:   "c11eabcc-8492-4dfa-a586-97d9f1694a8a",
					ExpiresAt: jwt.NewNumericDate(time.Now().UTC().Add(time.Hour)),
					IssuedAt:  jwt.NewNumericDate(time.Now().UTC()),
				},
				Roles: []string{},
			},
			requiredRole:           "user",
			expectedJWTGeneration:  true,
			expectedAuthentication: true,
			expectedAuthorization:  false,
		},
		{
			name: "Fail_SingleRole",
			claims: auth.Claims{
				RegisteredClaims: jwt.RegisteredClaims{
					Issuer:    "critiquefy",
					Subject:   "c11eabcc-8492-4dfa-a586-97d9f1694a8a",
					ExpiresAt: jwt.NewNumericDate(time.Now().UTC().Add(time.Hour)),
					IssuedAt:  jwt.NewNumericDate(time.Now().UTC()),
				},
				Roles: []string{"admin"},
			},
			requiredRole:           "user",
			expectedJWTGeneration:  true,
			expectedAuthentication: true,
			expectedAuthorization:  false,
		},
		{
			name: "Fail_IssuerMismatch",
			claims: auth.Claims{
				RegisteredClaims: jwt.RegisteredClaims{
					Issuer:    "critics",
					Subject:   "c11eabcc-8492-4dfa-a586-97d9f1694a8a",
					ExpiresAt: jwt.NewNumericDate(time.Now().UTC().Add(time.Hour)),
					IssuedAt:  jwt.NewNumericDate(time.Now().UTC()),
				},
			},
			expectedJWTGeneration:  true,
			expectedAuthentication: false,
			expectedAuthorization:  false,
		},
		{
			name: "Fail_ExpiredToken",
			claims: auth.Claims{
				RegisteredClaims: jwt.RegisteredClaims{
					Issuer:    "critics",
					Subject:   "c11eabcc-8492-4dfa-a586-97d9f1694a8a",
					ExpiresAt: jwt.NewNumericDate(time.Now().UTC().Add(-time.Minute)),
					IssuedAt:  jwt.NewNumericDate(time.Now().UTC().Add(-time.Hour)),
				},
			},
			expectedJWTGeneration:  true,
			expectedAuthentication: false,
			expectedAuthorization:  false,
		},
		{
			name: "Fail_TamperedToken",
			claims: auth.Claims{
				RegisteredClaims: jwt.RegisteredClaims{
					Issuer:    "critics",
					Subject:   "c11eabcc-8492-4dfa-a586-97d9f1694a8a",
					ExpiresAt: jwt.NewNumericDate(time.Now().UTC().Add(-time.Minute)),
					IssuedAt:  jwt.NewNumericDate(time.Now().UTC().Add(-time.Hour)),
				},
				Roles: []string{},
			},
			tokenOverride:          "eyJhbGciOiJSUzI1NiIsImtpZCI6InM0c0tJakQ5a0lSanhzMnR1bFBxR0xkeFnmZ1BFclJOMU11M0hkOWs5TlEiLCJ0eXAiOiJKV1QifQ.eyJpc3MiOiJjcml0aWNzIiwic3ViIjoiYzExZWFiY2MtODQ5Mi00ZGZhLWE1ODYtOTdkOWYxNjk0YThhIiwiZXhwIjoxNzQ3NDY0NDA1LCJpYXQiOjE3NDc0NjA4NjUsInJvbGVzIjpbXX0.GUF9giEm9kPo-KKMZE9K_o52DfiGNj2j_ocnfVMMJEWjfgAj4YKpS8AMQjFwAuyAwbJckmBd8n2wpWqxVBH0g4i9E_6p2hsmjPO0Olf6hmOyKeQ5QX5CvOdNyv3x2Wt9LfyqCfKeaCWB4YQV5H8N8gPuMUXeQGqG_raDrDx2q4ZZgPORU8KsD852uhgGkTIEeTnjUYxecW-RfTWPG-FjrSgSBM7HB6xE1mlIYm7glhv0Bkh-Bz6drG2YJmRCgHHmWmPl4ElYfQ9w5xSIUXMR_Z6QNsnZIhoFRVJ5Qo7HwQd3qPAcho1mQMr_MzspzQLHP5nWFLhtTmp79sqy_EO-Sw",
			expectedJWTGeneration:  true,
			expectedAuthentication: false,
			expectedAuthorization:  false,
		},
		{
			name: "Fail_MissingKid",
			claims: auth.Claims{
				RegisteredClaims: jwt.RegisteredClaims{
					Issuer:    "critics",
					Subject:   "c11eabcc-8492-4dfa-a586-97d9f1694a8a",
					ExpiresAt: jwt.NewNumericDate(time.Now().UTC().Add(-time.Minute)),
					IssuedAt:  jwt.NewNumericDate(time.Now().UTC().Add(-time.Hour)),
				},
				Roles: []string{},
			},
			tokenOverride:          "eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9.eyJpc3MiOiJjcml0aXF1ZWZ5Iiwic3ViIjoiYzExZWFiY2MtODQ5Mi00ZGZhLWE1ODYtOTdkOWYxNjk0YThhIiwiZXhwIjoxNzQ3NDY4ODk2LCJpYXQiOjE3NDc0NjUyOTYsInJvbGVzIjpbImFkbWluIl19.mdZn-n4x0-3pW2E1v0V1mZqqrlciWDYoRwWbc1oQAM7bgxj902PRP5HeSRsCOV7hDWYMNPH8Z7At4o1sE5NV-Zx9OHvxsojMbpo5ZhSuQLS4PzxzKT8MsLX-ddyw16j7MBO_pSIKR-lM66IN9xB2tehKFR9s9o7GjqSg8CHxh-Hs2snthDIIDPKHTcL38ER46njOmFGso17FqZfM_3Fb_wUAuqyoq9mWbPV_CQHuCCWZTjli7dR_FB62QgYa1exFafLXIV9JGrBXjTDA-s7jQFls7kDfhR5-UVlvTVR1rNLD8PyfyepqSflZB48-FhIz7fujtVG46WlB8Bqv3Tzk8g",
			expectedJWTGeneration:  true,
			expectedAuthentication: false,
			expectedAuthorization:  false,
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			token, err := a.GenerateToken(kid, c.claims)
			if err != nil && c.expectedJWTGeneration {
				t.Fatalf("Should be able to generate a JWT: %s", err)
			}

			if c.tokenOverride != "" {
				token = c.tokenOverride
			}

			fmt.Println(token)

			parsedClaims, err := a.Authenticate(context.Background(), "Bearer "+token)
			if err != nil && c.expectedAuthentication {
				t.Fatalf("Should be able to authenticate the claims: %s", err)
			}

			if !parsedClaims.HasRole(c.requiredRole) && c.expectedAuthorization {
				t.Fatalf("Should be able to check claim roles: %s", err)
			}

			err = a.Authorize(context.Background(), parsedClaims, "admin")
			if err != nil && c.expectedAuthorization {
				t.Errorf("Should be able to authorize admin role: %s", err)
			}
		})
	}
}

func newUnit(t *testing.T) (*logger.Logger, func()) {
	var buf bytes.Buffer
	log := logger.New(&buf, logger.LevelInfo, "TEST", func(context.Context) string { return "00000000-0000-0000-0000-000000" })

	teardown := func() {
		t.Helper()

		fmt.Println("******************** LOGS ******************** ")
		fmt.Print(buf.String())
		fmt.Println("******************** LOGS ******************** ")
	}

	return log, teardown
}

type keyStore struct{}

func (ks *keyStore) PrivateKey(kid string) (string, error) {
	return privateKeyPEM, nil
}

func (ks *keyStore) PublicKey(kid string) (string, error) {
	return publicKeyPEM, nil
}

const (
	kid = "s4sKIjD9kIRjxs2tulPqGLdxSfgPErRN1Mu3Hd9k9NQ"

	privateKeyPEM = `-----BEGIN PRIVATE KEY-----
MIIEvQIBADANBgkqhkiG9w0BAQEFAASCBKcwggSjAgEAAoIBAQDUc8weNl/Ghyd5
PzaDjkHIPzsdR5y1R4xNCoLLoNcvlB8L5gK495QwcEa5rToBqU9hniaRQqWF71Xv
ejtCRr87kiOS2fhHI94kQzLEfCMTGV6q7U1k5c7VpqeR6+9sdLOGjd5K3xWgYX47
zWdttJAr6BkjsG6RvRUvYf5lAgZZbuC9l7ZmiYrr+wKuUC50/ooZDBBtByqmJmkU
s/jSGJjfxlpnIYcsLykNna0DIHmVJD6N+S93QSKq8bWobMD0u8yilkc1Nex6VlVD
5EWQpOhX0vKa61RtFTMbAMEEeODrulU+J2lV1JlrV1I4ueR2vuTD3sJTGi1LJ9Dc
hqASvHVFAgMBAAECggEAaLKcT/NJ5cNrT5Q4YEK15mJK1pYZAzk8Sic45/LeuQLM
/gcfJlpUPD7Ii+5zXKg8h4XxybpHaibVecwJ8hJ9YXUWdONYOG7TpZk8JppqipoB
Dkkdz/B0qtOTVxUni5JDerblao5f0Qbat8v1AZpvRkP+R5lGFCpTi2NGhC6oRF+7
ES41s4UPbob1dOIMCX6QDi3e+hj7/PaRhcmc8Nmmy7pi90mb633eLC1VEBh0wbkc
NutW8JGp0m3v031YH3U/Pi7wEbzUQ4dtX2s6KpvRIwTh2zVYgOhC7oasCCRr193H
8rFaCqLhc1a1WFbvgpyaJzgAbDGxRoKX5BEgk0dt7QKBgQDvGaPQFagL+pQDyKbG
QzwtYw2fIv12R21uAe8SjFlCVfLHM32WDnO2GNT6TF4t/5vFc1cl4P7qBLzLSyy4
9y81Lu3q66dT4E/c5bTp9UEkPxuJUqrKYhp6p4q0L+jReJ+WNHQKEMt0rnjNx1gC
39lfYEypfixw5zeU05gjEKld2wKBgQDjd/uMs6t8TH9qywBMFTwz8oNHx1rpAoac
wwUtVlEfsUBH0OVxusnKZnX73xmzgZNNEo0xKhLt2MV1uzl8r1gU1CUKkYGMlxxd
IEBCodF0OSKqMOpDJZkjtIxXvphR3vBLMn7iqo7eonkRx1uvFfbdrtKxPVm+5L52
FTWhgZAzXwKBgGnR0jNM8mPi0dFe45jJtv9rYGL27HCFqkPOrU1rOjHmsh1Bh6p6
2PFVyiTA2cnH39wicQZ9rrRJxni+25s9IvKJw5h+FT9E/nOIYmpNNjhhicFcCeSq
SIfSUMvwjDzxAshKjLTLvA/3C9YfDK+w/JZ+m09EXUzWuD2w7BtQy3STAoGBANRM
MRx6u/xAsVMMr/RShWO+XcRqTXDXiKdaZMSRoRlBJ0tfriVdPeSHiGpRKP2eW8o9
HEXcjNorzO86lEbIqB6YeRHKB+0dQ72u0greWExu3umUya9tseXfJnTmT+dpeT/V
mxMWOE2VugVb2Tgp+cOg3MfLCK3fc9tlpC5ebCVlAoGATh60RuO+6k5yvIseXHGO
zaGZTZJwU0IsKvLrYTXbap9g8p0e5bLbcdHrQbU8KooSS3dE9L5MDavkyvZuaJ5X
hhxnf440I8OyV65IlU4uTlDhqFm1eWTy4zZoixdEEN7d46oUPEh7DrY1a/D65Feq
q50/OohIJ6CTKrKt8qwpEn4=
-----END PRIVATE KEY-----
`
	publicKeyPEM = `-----BEGIN PUBLIC KEY-----
MIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEA1HPMHjZfxocneT82g45B
yD87HUectUeMTQqCy6DXL5QfC+YCuPeUMHBGua06AalPYZ4mkUKlhe9V73o7Qka/
O5Ijktn4RyPeJEMyxHwjExlequ1NZOXO1aankevvbHSzho3eSt8VoGF+O81nbbSQ
K+gZI7Bukb0VL2H+ZQIGWW7gvZe2ZomK6/sCrlAudP6KGQwQbQcqpiZpFLP40hiY
38ZaZyGHLC8pDZ2tAyB5lSQ+jfkvd0EiqvG1qGzA9LvMopZHNTXselZVQ+RFkKTo
V9LymutUbRUzGwDBBHjg67pVPidpVdSZa1dSOLnkdr7kw97CUxotSyfQ3IagErx1
RQIDAQAB
-----END PUBLIC KEY-----
`
)
