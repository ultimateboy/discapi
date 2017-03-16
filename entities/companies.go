package entities

import "github.com/rs/rest-layer/schema"

var (
	Companies = schema.Schema{
		Description: `Represents a company entity`,
		Fields: schema.Fields{
			// Custom ID field
			"id": schema.Field{
				Description: "The item's id",
				Required:    true,
				Filterable:  true,
				Sortable:    true,
				Validator: &schema.String{
					Regexp: "^[0-9a-z]*$",
				},
			},
			"created": schema.CreatedField,
			"updated": schema.UpdatedField,
			"website": {
				Validator: &schema.URL{},
			},
		},
	}
)
