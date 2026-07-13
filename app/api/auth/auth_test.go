package auth_test

import (
	"bytes"
	"context"
	"fmt"
	"runtime/debug"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
	"github.com/gsemer/ardanlabs-service/app/api/auth"
	"github.com/gsemer/ardanlabs-service/foundation/logger"
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
		Issuer:    "service project",
	}

	a, err := auth.New(cfg)
	if err != nil {
		t.Fatalf("Should be able to create an authenticator: %s", err)
	}

	claims := auth.Claims{
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    "service project",
			Subject:   "57af22c4-c715-4a6b-aef9-5bb69dfdba79",
			ExpiresAt: jwt.NewNumericDate(time.Now().UTC().Add(time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now().UTC()),
		},
		Roles: []string{"ADMIN"},
	}

	token, err := a.GenerateToken(kid, claims)
	if err != nil {
		t.Fatalf("Should be able to generate a JWT: %s", err)
	}

	parsedClaims, err := a.Authenticate(context.Background(), "Bearer "+token)
	if err != nil {
		t.Fatalf("Should be able to authenticate the claims: %s", err)
	}

	userID := uuid.MustParse(claims.Subject)

	err = a.Authorize(context.Background(), parsedClaims, userID, auth.RuleAdminOnly)
	if err != nil {
		t.Errorf("Should be able to authorize the RoleAdmin claims: %s", err)
	}

	err = a.Authorize(context.Background(), parsedClaims, userID, auth.RuleUserOnly)
	if err == nil {
		t.Error("Should NOT be able to authorize the RoleUser claims")
	}

	err = a.Authorize(context.Background(), parsedClaims, userID, auth.RuleAdminOrSubject)
	if err != nil {
		t.Errorf("Should be able to authorize the RuleAdminOrSubject claim with RoleAdmin only: %s", err)
	}

	// -------------------------------------------------------------------------

	claims = auth.Claims{
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    "service project",
			Subject:   "57af22c4-c715-4a6b-aef9-5bb69dfdba79",
			ExpiresAt: jwt.NewNumericDate(time.Now().UTC().Add(time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now().UTC()),
		},
		Roles: []string{"USER"},
	}

	token, err = a.GenerateToken(kid, claims)
	if err != nil {
		t.Fatalf("Should be able to generate a JWT: %s", err)
	}

	parsedClaims, err = a.Authenticate(context.Background(), "Bearer "+token)
	if err != nil {
		t.Fatalf("Should be able to authenticate the claims: %s", err)
	}

	userID = uuid.MustParse(claims.Subject)

	err = a.Authorize(context.Background(), parsedClaims, userID, auth.RuleUserOnly)
	if err != nil {
		t.Errorf("Should be able to authorize the RoleUserOnly claim with RoleUser only: %s", err)
	}

	err = a.Authorize(context.Background(), parsedClaims, userID, auth.RuleAdminOnly)
	if err == nil {
		t.Error("Should NOT be able to authorize the RuleAdminOnly claim with RoleUser only")
	}

	err = a.Authorize(context.Background(), parsedClaims, userID, auth.RuleAdminOrSubject)
	if err != nil {
		t.Errorf("Should be able to authorize the RuleAdminOrSubject claim with RoleUser only: %s", err)
	}

	err = a.Authorize(context.Background(), parsedClaims, userID, auth.RuleAny)
	if err != nil {
		t.Errorf("Should be able to authorize the RuleAny claim with RoleUser only: %s", err)
	}

	// -------------------------------------------------------------------------

	claims = auth.Claims{
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    "service project",
			Subject:   "57af22c4-c715-4a6b-aef9-5bb69dfdba79",
			ExpiresAt: jwt.NewNumericDate(time.Now().UTC().Add(time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now().UTC()),
		},
		Roles: []string{"USER"},
	}

	token, err = a.GenerateToken(kid, claims)
	if err != nil {
		t.Fatalf("Should be able to generate a JWT: %s", err)
	}

	parsedClaims, err = a.Authenticate(context.Background(), "Bearer "+token)
	if err != nil {
		t.Fatalf("Should be able to authenticate the claims: %s", err)
	}

	userID = uuid.MustParse("9e9722c4-c715-4a6b-aef9-5bb69dfdba79")

	err = a.Authorize(context.Background(), parsedClaims, userID, auth.RuleAdminOrSubject)
	if err == nil {
		t.Error("Should NOT be able to authorize the RuleAdminOrSubject claim with RoleUser only and different userID")
	}
}

func newUnit(t *testing.T) (*logger.Logger, func()) {
	var buf bytes.Buffer
	log := logger.New(&buf, logger.LevelInfo, "TEST", func(context.Context) string { return "00000000-0000-0000-0000-000000000000" })

	teardown := func() {
		t.Helper()

		fmt.Println("************** LOGS **************")
		fmt.Print(buf.String())
		fmt.Println("************** LOGS **************")
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
	kid = "57af22c4-c715-4a6b-aef9-5bb69dfdba79"

	privateKeyPEM = `-----BEGIN PRIVATE KEY-----
MIIEpAIBAAKCAQEAvsKxsiVVYyKKObFUp2f/LwQ8iMfc/jYOGZzh/GuJDmZcAi8k
eS7W/ZNnFZRvlM7E+nYdvYjx3huKns0x5efOBx2dT/fRXOU+3XMR89aAIVcU0e9b
6nM9HSPrX98RscPzq4XeKjcmVmHVOfKKOo1ENnUyMm5JgADP79kcsL5GfqFpoEn1
3mXZdwnkfGfBNCWxw6Rlv+Vqb1DtT7kon0FXUeRuMapKQYudgEkF61d4Fry4wt4K
A1eLfRzbTlbRibMo7K4+v8q+PVW/Kf/piY8fDIvPCCtk6crNYSy3nXGQdjw3zzH2
EUW0jZzSl93c2KInc0Ad8qpGYIbjPAaI9eXxpQIDAQABAoIBAA6wAuqKgVaOtEHY
64GwOi+ujdKiQNu54cALGkNLLFRVgUQRyScjeh4wGUHKGgVFHlmCequ7PZQyXqv3
dJ4VCQH3P8OGezJB3GNElt9FZrwqbknzugoFMXFq8JaDIGOliL9uITry4BrKkZZS
nF4BvnzK7UCAyVv3tArtlo3tOJRLS0A8Q/elMptPIMa1Ztq64WOBzAKyqkDEs3DW
KPwGrAaOG1TnxLJH2H+5fGvDkZrsm5oOgCxy4+t2D2gb93aI7KIDzJsC87XI2vFf
xxf3Gkh7zV1XN6wc07V7Bj1PxABGgDgig172I03m6MMPVuK6MoWYKsUQgDq3m/yH
UL6h2UECgYEA1cdsXGDZa7VTjrZ6dPKVJFNW6RnMkqQuzSP1XxPOZLHFp0h5T0zr
oPRFsvyhviooh3n6XikHEv3j0V+lQpqmcD/B0jopjjpmUO0Snf1TKVgrKxCrDLUX
7p3KEdCMMBWHBHwTKJ3MuraTo69e6PG2w5bDE8Xm/DwQMpkDVPbN5T0CgYEA5G93
aZpOUYX2QS8JhxmZra90x2ePNMP7dNFG5B6Xxwi5uNtwAt4SZrPj3Kkat6q3X+z0
msIcgErEAP6n1tkc/42MZ7ZBKeLX2n26O9VhsVZQT+UBQdDUkOw9OID/sPCpz6b+
95YY+aKRHVsl3i9Px08cUt0mDWdWl0NdcZbblIkCgYBPzp0jd3xze0PwWSsqEY2f
/ATMDLeUXvqh1rS5g9lfOgaBxsqS0jJ86fRDN5DiPzbWLLFNCZ/8dQ/hkAVP8hAE
g6jF5LSyxhaAS6DRnkq3epTTBOv2WHzQtdNEB0jugnrfL7qvRQmzAonnZ4bVC7eh
GN1GunDa5UleukTGKUUAOQKBgQDZNjiY3NYt5LDlCIIJycj8g4MKfSmJ5fU7/idn
kMOHyX84DMi0oU9kAxffYZj7HkSh3SI16e/J+c3omD0mKWrOgV0J3R6XYpEXvEeS
z3LGeqmBXuNUHuuRJmGMUfVP3XfK8SMub7Yt4WwVOu+GFvzIKyxmisy9IA8RZEf9
U5JV6QKBgQCU45PIqORDaXQSJgxHQWOxOHGGqx99Zz6k9LC+NLHgSI5XkmPJWKqb
Jj3IcWl5/lQSW8HNPy4wL/IQ7XNZX4mlzSUvq9uEhpaKKYE01Ed4sXyZSh/MPWy/
y5sRAu/2bdCx0aLCAURttzIN6gmGqrfYyf6DhGYEuxMZ4hGuzAw5zA==
-----END PRIVATE KEY-----
`

	publicKeyPEM = `-----BEGIN PUBLIC KEY-----
MIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEAvsKxsiVVYyKKObFUp2f/
LwQ8iMfc/jYOGZzh/GuJDmZcAi8keS7W/ZNnFZRvlM7E+nYdvYjx3huKns0x5efO
Bx2dT/fRXOU+3XMR89aAIVcU0e9b6nM9HSPrX98RscPzq4XeKjcmVmHVOfKKOo1E
NnUyMm5JgADP79kcsL5GfqFpoEn13mXZdwnkfGfBNCWxw6Rlv+Vqb1DtT7kon0FX
UeRuMapKQYudgEkF61d4Fry4wt4KA1eLfRzbTlbRibMo7K4+v8q+PVW/Kf/piY8f
DIvPCCtk6crNYSy3nXGQdjw3zzH2EUW0jZzSl93c2KInc0Ad8qpGYIbjPAaI9eXx
pQIDAQAB
-----END PUBLIC KEY-----`
)
