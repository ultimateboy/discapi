package entities

import "github.com/rs/rest-layer/schema"

var (
	Discs = schema.Schema{
		Description: `Represents a disc entity`,
		Fields: schema.Fields{
			"id":      schema.IDField,
			"created": schema.CreatedField,
			"updated": schema.UpdatedField,
			"company": {
				Validator: &schema.Reference{
					Path: "companies",
				},
			},
			"name": {
				Required:  true,
				Validator: &schema.String{},
			},
			"speed": {
				Validator:  &schema.Float{},
				Filterable: true,
			},
			"glide": {
				Validator:  &schema.Float{},
				Filterable: true,
			},
			"turn": {
				Validator:  &schema.Float{},
				Filterable: true,
			},
			"fade": {
				Validator:  &schema.Float{},
				Filterable: true,
			},
		},
	}
)
