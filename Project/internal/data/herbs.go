package data

import (
	"time"
)

type Herb struct {
	ID           int64     `json:"id"`                      // Unique integer ID for the movie
	CreatedAt    time.Time `json:"-"`                       // Timestamp for when the movie is added to our database
	Name         string    `json:"name"`                    // Herb name
	Description  string    `json:"description"`             // Herb description
	Price        Price     `json:"price"`                   // Herb price
	CulinaryUses []string  `json:"culinary_uses,omitempty"` // Culinary uses of Herb
	Version      int32     `json:"version"`                 // The version number starts at 1 and will be incremented
	// each time the herb information is updated
}
