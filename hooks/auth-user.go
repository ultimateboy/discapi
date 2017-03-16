// This file provides rest-layer hooks which responds to all requests, inspects the
// context for a User ID, and reject access if not found.

package hooks

import (
	"context"

	"github.com/ultimateboy/discapi/utils"

	"github.com/rs/rest-layer/resource"
)

// AuthUser is a rest-layer resource event handler that protect the resource from unauthorized users
type AuthUser struct{}

// OnFind implements resource.FindEventHandler interface
func (a AuthUser) OnFind(ctx context.Context, lookup *resource.Lookup, offset, limit int) error {
	// Reject unauthorized users
	if err := rejectUnauthorizedUser(ctx); err != nil {
		return err
	}

	return nil
}

// OnGot implements resource.GotEventHandler interface
func (a AuthUser) OnGot(ctx context.Context, item **resource.Item, err *error) {
	// Do not override existing errors
	if err != nil {
		return
	}

	if e := rejectUnauthorizedUser(ctx); e != nil {
		*err = e
		return
	}

	return
}

// OnInsert implements resource.InsertEventHandler interface
func (a AuthUser) OnInsert(ctx context.Context, items []*resource.Item) error {
	if err := rejectUnauthorizedUser(ctx); err != nil {
		return err
	}

	return nil
}

// OnUpdate implements resource.UpdateEventHandler interface
func (a AuthUser) OnUpdate(ctx context.Context, item *resource.Item, original *resource.Item) error {
	if err := rejectUnauthorizedUser(ctx); err != nil {
		return err
	}

	return nil
}

// OnDelete implements resource.DeleteEventHandler interface
func (a AuthUser) OnDelete(ctx context.Context, item *resource.Item) error {
	if err := rejectUnauthorizedUser(ctx); err != nil {
		return err
	}

	return nil
}

// OnClear implements resource.ClearEventHandler interface
func (a AuthUser) OnClear(ctx context.Context, lookup *resource.Lookup) error {
	if err := rejectUnauthorizedUser(ctx); err != nil {
		return err
	}

	return nil
}

// rejectUnauthorizedUser looks for the User ID in the context, provided by the
// JWT middleware. If not found, returns an unauthored error, otherwise nil.
func rejectUnauthorizedUser(ctx context.Context) error {
	_, found := utils.UserFromContext(ctx)
	if !found {
		return resource.ErrUnauthorized
	}

	return nil
}
