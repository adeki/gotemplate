package app

import (
	"context"
	"errors"
)

// sample
type User struct {
	ID   int64
	Name string
}

type ctxKey string

const (
	userCtxKey ctxKey = "user"
)

func setCtxUser(ctx context.Context, u *User) context.Context {
	return context.WithValue(ctx, userCtxKey, u)
}

func getCtxUser(ctx context.Context) (*User, error) {
	v := ctx.Value(userCtxKey)
	u, ok := v.(*User)
	if !ok {
		return nil, errors.New("user not found")
	}
	return u, nil
}
