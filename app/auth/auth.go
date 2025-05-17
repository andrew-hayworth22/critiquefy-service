package auth

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/andrew-hayworth22/critiquefy-service/foundation/logger"
	"github.com/golang-jwt/jwt/v5"
)

// ErrForbidden represents the error thrown when a user is unauthorized to do something
var ErrForbidden = errors.New("attempted action is not allowed")

// Claims are all of the claims in a user's JWT token
type Claims struct {
	jwt.RegisteredClaims
	Roles []string `json:"roles"`
}

// HasRole checks if a user has a given role
func (c Claims) HasRole(r string) bool {
	for _, role := range c.Roles {
		if role == r {
			return true
		}
	}
	return false
}

// KeyLookup defines the functionality needed to fetch RSA keys
type KeyLookup interface {
	PrivateKey(kid string) (key string, err error)
	PublicKey(kid string) (key string, err error)
}

// Config represents the configuration needed for auth logic
type Config struct {
	Log       *logger.Logger
	KeyLookup KeyLookup
	Issuer    string
}

// Auth represents a package providing auth logic
type Auth struct {
	keyLookup KeyLookup
	method    jwt.SigningMethod
	parser    *jwt.Parser
	issuer    string
}

// New constructs a new Auth package
func New(cfg Config) (*Auth, error) {
	a := Auth{
		keyLookup: cfg.KeyLookup,
		method:    jwt.GetSigningMethod(jwt.SigningMethodRS256.Name),
		parser:    jwt.NewParser(jwt.WithValidMethods([]string{jwt.SigningMethodRS256.Name})),
		issuer:    cfg.Issuer,
	}
	return &a, nil
}

// GenerateToken produces an auth token with the provided claims
func (a *Auth) GenerateToken(kid string, claims Claims) (string, error) {
	token := jwt.NewWithClaims(a.method, claims)
	token.Header["kid"] = kid

	privateKeyPEM, err := a.keyLookup.PrivateKey(kid)
	if err != nil {
		return "", fmt.Errorf("private key lookup: %w", err)
	}

	privateKey, err := jwt.ParseRSAPrivateKeyFromPEM([]byte(privateKeyPEM))
	if err != nil {
		return "", fmt.Errorf("parsing private key: %w", err)
	}

	str, err := token.SignedString(privateKey)
	if err != nil {
		return "", fmt.Errorf("parsing private key: %w", err)
	}

	return str, nil
}

// Authenticate validates a bearer token and returns its claims
func (a *Auth) Authenticate(ctx context.Context, bearerToken string) (Claims, error) {
	parts := strings.Split(bearerToken, " ")
	if len(parts) != 2 || parts[0] != "Bearer" {
		return Claims{}, errors.New("expected authorization header format: Bearer <token>")
	}

	var claims Claims
	_, err := a.parser.ParseWithClaims(parts[1], &claims, func(t *jwt.Token) (interface{}, error) {
		kidRaw, ok := t.Header["kid"]
		if !ok {
			return []byte{}, errors.New("kid missing from header")
		}

		issuer, err := t.Claims.GetIssuer()
		if err != nil {
			return Claims{}, fmt.Errorf("fetching issuer: %w", err)
		}

		if issuer != a.issuer {
			return Claims{}, errors.New("issuer mismatch")
		}

		kid, ok := kidRaw.(string)
		if !ok {
			return Claims{}, errors.New("kid malformed")
		}

		publicKeyPEM, err := a.keyLookup.PublicKey(kid)
		if err != nil {
			return Claims{}, fmt.Errorf("fetching public key: %w", err)
		}

		publicKey, err := jwt.ParseRSAPublicKeyFromPEM([]byte(publicKeyPEM))
		if err != nil {
			return Claims{}, fmt.Errorf("parsing public key")
		}

		return publicKey, nil
	})
	if err != nil {
		return Claims{}, fmt.Errorf("error validating token: %w", err)
	}

	return claims, nil
}

// Authorize validates that a user has a role
func (a *Auth) Authorize(ctx context.Context, claims Claims, role string) error {
	for _, r := range claims.Roles {
		if r == role {
			return nil
		}
	}
	return errors.New("user does not have role")
}
