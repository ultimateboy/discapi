package entities

import "github.com/rs/rest-layer/schema"

var (
	// EmailField provides the user's email address for reuse on the auth endpoint
	EmailField = schema.Field{
		Required:   true,
		Filterable: true,
		Validator: &schema.String{
			// anything@anything
			Regexp: "^.{1,}@.{1,}$",
		},
	}

	Users = schema.Schema{
		Description: `The user object`,
		Fields: schema.Fields{
			"id":       schema.IDField,
			"email":    EmailField,
			"password": schema.PasswordField,
			"bag": {
				Validator: &schema.Array{
					ValuesValidator: &schema.Reference{
						Path: "discs",
					},
				},
			},
		},
	}
)
