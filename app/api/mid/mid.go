// Package mid provides app level middleware support.
package mid

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/gsemer/ardanlabs-service/business/api/auth"
)

// Handler represents the handler function that needs to be called.
type Handler func(ctx context.Context) error

type ctxKey int

const (
	claimKey ctxKey = iota + 1
	userIDKey
)

func setClaims(ctx context.Context, claims auth.Claims) context.Context {
	return context.WithValue(ctx, claimKey, claims)
}

func GetClaims(ctx context.Context) auth.Claims {
	v, ok := ctx.Value(claimKey).(auth.Claims)
	if !ok {
		return auth.Claims{}
	}
	return v
}

func setUserID(ctx context.Context, userID uuid.UUID) context.Context {
	return context.WithValue(ctx, userIDKey, userID)
}

func GetUserID(ctx context.Context) (uuid.UUID, error) {
	v, ok := ctx.Value(userIDKey).(uuid.UUID)
	if !ok {
		return uuid.UUID{}, errors.New("user not found in context")
	}
	return v, nil
}
