// This file provides rest-layer hooks which respond to /auth requests,
// verifies the email and password combination, and attaches the signed JWT to
// the response.

package hooks

import (
	"context"
	"errors"
	"log"

	"github.com/ultimateboy/discapi/utils"

	"github.com/rs/rest-layer/resource"
	"github.com/rs/rest-layer/rest"
	"github.com/rs/rest-layer/schema"
)

// ErrNoInsert is a fake error to trick the storage handler to not store the entity
var ErrNoInsert = errors.New("No insert error")

// AuthLogin is a rest-layer resource event handler which generates a jwt and adds
// to the payload if a valid email/password is provided
type AuthLogin struct {
	SigningKey   []byte
	UserResource *resource.Resource
}

// OnInsert implements rest-layer InsertEventHandler
func (h AuthLogin) OnInsert(ctx context.Context, items []*resource.Item) error {
	for _, i := range items {

		// Check the provided email and password
		email := i.GetField("email").(string)
		password := i.GetField("password").(string)
		if email == "" || password == "" {
			// We dont allow users with blank email or passwords, deny.
			return rest.ErrUnauthorized
		}

		// Find the user with the provided email address
		userList, err := utils.FindUser(ctx, h.UserResource, email)
		if err != nil {
			log.Printf("Error finding user: %s\n", err)
			return err
		}
		// There will only ever be a single user returned
		if len(userList.Items) != 1 {
			return rest.ErrUnauthorized
		}

		user := userList.Items[0]

		// If the password verifies, include the JWT in the response
		if schema.VerifyPassword(user.Payload["password"], []byte(password)) {
			// @todo how expensive is this? should we store it?
			userJWT, err := utils.GenerateJWT(h.SigningKey, user.ID.(string))
			if err != nil {
				log.Printf("Error generating jwt: %s", err)
				return err
			}
			i.Payload["jwt"] = userJWT
		} else {
			// If the provided password failed to validate, deny.
			return rest.ErrUnauthorized
		}
	}

	// Success. Return a predefined error to trick the storage system to not process the post
	return ErrNoInsert
}

// OnInserted implements InsertedEventHandler
func (h AuthLogin) OnInserted(ctx context.Context, items []*resource.Item, err *error) {
	// Catch the fake successful authentication, no-insert error.
	// Any other error will fall through.
	if *err == ErrNoInsert {
		*err = nil
	}
}
