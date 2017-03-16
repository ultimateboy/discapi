package hooks

import (
	"context"
	"log"

	"github.com/rs/rest-layer/resource"
	"github.com/rs/rest-layer/rest"

	"github.com/ultimateboy/discapi/utils"
)

// UserHook provides a rest-layer hook to react to changes to users
type UserHook struct {
	UserResource *resource.Resource
}

// OnInsert implements resource.InsertEventHandler interface
func (h UserHook) OnInsert(ctx context.Context, items []*resource.Item) error {
	for _, i := range items {
		email := i.GetField("email")
		password := i.GetField("password")
		if email == nil || password == nil {
			// We dont allow users with blank email or passwords, deny if blank.
			return rest.ErrUnauthorized
		}

		// Determine if a user with this email address already exists
		userList, err := utils.FindUser(ctx, h.UserResource, email.(string))
		if err != nil {
			log.Printf("Error finding user: %s\n", err)
			return err
		}
		if len(userList.Items) > 0 {
			// @todo add details for deny to body... user already exists
			return rest.ErrConflict
		}
	}
	return nil
}

// OnUpdate implements UpdateEventHandler
func (h UserHook) OnUpdate(ctx context.Context, item *resource.Item, original *resource.Item) error {
	actingUser, ok := utils.UserFromContext(ctx)
	if !ok {
		return rest.ErrUnauthorized
	}

	// If the user being updated is not the same as the actor
	if actingUser.ID != item.ID {
		return rest.ErrUnauthorized
	}

	return nil
}

func (h UserHook) OnDelete(ctx context.Context, item *resource.Item) error {
	actingUser, ok := utils.UserFromContext(ctx)
	if !ok {
		return rest.ErrUnauthorized
	}

	// If the user being deleted is not the same as the actor
	if actingUser.ID != item.ID {
		return rest.ErrUnauthorized
	}
	return nil
}
