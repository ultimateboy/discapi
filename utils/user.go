package utils

import (
	"context"

	"github.com/rs/rest-layer/resource"
	"github.com/rs/rest-layer/schema"
)

// userKey is the key for user.User values in Contexts.
// @see context.Value()
const userKey key = "user"

// NewContextWithUser stores user into context
func NewContextWithUser(ctx context.Context, user *resource.Item) context.Context {
	return context.WithValue(ctx, userKey, user)
}

// UserFromContext retrieves user from context
func UserFromContext(ctx context.Context) (*resource.Item, bool) {
	user, ok := ctx.Value(userKey).(*resource.Item)
	return user, ok
}

// FindUser finds a user by a given email address
func FindUser(ctx context.Context, userResource *resource.Resource, email string) (*resource.ItemList, error) {
	query, err := schema.NewQuery(map[string]interface{}{}, userResource.Validator())
	if err != nil {
		return nil, err
	}
	userLookup := resource.NewLookupWithQuery(query)
	userLookup.AddQuery(schema.Query{
		schema.Equal{
			Field: "email",
			Value: email,
		},
	})

	// Limit one. There can only ever be one.
	return userResource.Find(ctx, userLookup, 0, 1)
}
