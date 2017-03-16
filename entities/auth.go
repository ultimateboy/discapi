package entities

import "github.com/rs/rest-layer/schema"

var (
	// IDFieldNoRead is a hidden schema.IDField
	IDFieldNoRead = schema.Field{
		Description: "The item's id",
		Required:    true,
		Hidden:      true,
		OnInit:      schema.NewID,
		Filterable:  true,
		Sortable:    true,
		Validator: &schema.String{
			// This regexp matches a base32 id
			Regexp: "^[0-9a-v]{20}$",
		},
	}

	// Auth is a resource schema used to retrieve JWT tokens.
	// The entire entity skips storage and only alters the response of a POST
	// if a valid email/password is provded.
	//
	// It's similar to the user resource, but has a special no-read ID Field
	// (doesn't really make sense in this context given we're skipping storage)
	// and a plain-text password field to do the validation.
	Auth = schema.Schema{
		Description: `The auth endpoint`,
		Fields: schema.Fields{
			// rest-layer requires an ID field although it's not used on this endpoint
			// @todo determine if we can avoid having this field all together
			"id":    IDFieldNoRead,
			"email": EmailField,
			// This field cannot be shared with the user entity because the raw
			// password is needed to perform the authentication. This entire entity
			// skips storage, so the plain-text aspect is nothing to worry about.
			"password": schema.Field{
				Required:   true,
				Filterable: false,
				Sortable:   false,
				Hidden:     true,
				Validator:  &schema.String{},
			},
		},
	}
)
