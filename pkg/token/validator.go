package token

import (
	"context"
	"fmt"
	"time"

	"github.com/lestrrat-go/httprc/v3"
	"github.com/lestrrat-go/jwx/v3/jwk"
	"github.com/lestrrat-go/jwx/v3/jwt"
)

type Validator struct {
	cache *jwk.Cache
	url   string
}

func NewValidator(jwksURL string) *Validator {
	ctx := context.Background()
	c, _ := jwk.NewCache(ctx, httprc.NewClient())
	// Register URL with a 15-minute refresh timer
	if err := c.Register(ctx, jwksURL, jwk.WithConstantInterval(15*time.Minute)); err != nil {
		fmt.Printf("Warning: Failed to register JWKS URL %s: %v\n", jwksURL, err)
	}

	return &Validator{
		cache: c,
		url:   jwksURL,
	}
}

func (v *Validator) Validate(ctx context.Context, tokenString string) (jwt.Token, error) {
	keySet, err := v.cache.Lookup(ctx, v.url)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch public keys: %w", err)
	}

	token, err := jwt.Parse(
		[]byte(tokenString),
		jwt.WithKeySet(keySet),
		jwt.WithValidate(true),
		// FIX: Allow a 1-minute difference between Auth and Account service clocks
		jwt.WithAcceptableSkew(1*time.Minute),
	)
	if err != nil {
		// FIX: Return the REAL error so we can debug it (e.g. "iat is in the future")
		return nil, fmt.Errorf("token validation failed: %w", err)
	}

	return token, nil
}
