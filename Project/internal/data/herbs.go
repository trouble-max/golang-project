package data

import (
	"time"

	"gourmetspices.yerassyl.net/internal/validator"
)

type Herb struct {
	ID           int64     `json:"id"`                      // Unique integer ID for the movie
	CreatedAt    time.Time `json:"-"`                       // Timestamp for when the movie is added to our database
	Name         string    `json:"name"`                    // Herb name
	Description  string    `json:"description,omitempty"`   // Herb description
	Price        Price     `json:"price"`                   // Herb price
	CulinaryUses []string  `json:"culinary_uses,omitempty"` // Culinary uses of Herb
	Version      int32     `json:"version"`                 // The version number starts at 1 and will be incremented
	// each time the herb information is updated
}

func ValidateHerb(v *validator.Validator, herb *Herb) {
	v.Check(herb.Name != "", "name", "must be provided")
	v.Check(len(herb.Name) <= 500, "name", "must not be more than 500 bytes long")

	v.Check(herb.Description != "", "description", "must be provided")
	v.Check(len(herb.Description) <= 500, "description", "must not be more than 500 bytes long")

	v.Check(herb.Price != 0, "price", "must be provided")
	v.Check(herb.Price >= 0, "price", "must be a positive integer")

	v.Check(herb.CulinaryUses != nil, "culinary_uses", "must be provided")
	v.Check(len(herb.CulinaryUses) >= 1, "culinary_uses", "must contain at least 1 use")
	v.Check(len(herb.CulinaryUses) <= 5, "culinary_uses", "must not contain more than 5 uses")
	v.Check(validator.Unique(herb.CulinaryUses), "culinary_uses", "must not contain duplicate values")
}
