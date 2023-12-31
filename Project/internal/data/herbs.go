package data

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/lib/pq"
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

type HerbModel struct {
	DB *sql.DB
}

func (h HerbModel) Insert(herb *Herb) error {

	query := `INSERT INTO herbs (name, description, price, culinary_uses) 
		VALUES ($1, $2, $3, $4)
		RETURNING id, created_at, version`

	args := []interface{}{
		herb.Name,
		herb.Description,
		herb.Price,
		pq.Array(herb.CulinaryUses),
	}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	return h.DB.QueryRowContext(ctx, query, args...).Scan(
		&herb.ID,
		&herb.CreatedAt,
		&herb.Version,
	)
}

func (h HerbModel) Get(id int64) (*Herb, error) {

	if id < 1 {
		return nil, ErrRecordNotFound
	}

	query := `SELECT id, created_at, name, description, price, culinary_uses, version
		FROM herbs
		WHERE id = $1`

	var herb Herb

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := h.DB.QueryRowContext(ctx, query, id).Scan(
		&herb.ID,
		&herb.CreatedAt,
		&herb.Name,
		&herb.Description,
		&herb.Price,
		pq.Array(&herb.CulinaryUses),
		&herb.Version,
	)

	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrRecordNotFound
		default:
			return nil, err
		}

	}

	return &herb, nil
}

func (h HerbModel) GetAll(name string, culinaryUses []string, filters Filters) ([]*Herb, Metadata, error) {

	query := fmt.Sprintf(`SELECT count(*) OVER(), id, created_at, name, description, price, culinary_uses, version
		FROM herbs
		WHERE (to_tsvector('simple', name) @@ plainto_tsquery('simple', $1) OR $1 = '')
		AND (culinary_uses @> $2 OR $2 = '{}')
		ORDER BY %s %s, id ASC
		LIMIT $3 OFFSET $4`, filters.sortColumn(), filters.sortDirection())

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	args := []interface{}{name, pq.Array(culinaryUses), filters.limit(), filters.offset()}

	rows, err := h.DB.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, Metadata{}, err
	}
	defer rows.Close()

	totalRecords := 0
	herbs := []*Herb{}

	for rows.Next() {

		var herb Herb

		err := rows.Scan(
			&totalRecords,
			&herb.ID,
			&herb.CreatedAt,
			&herb.Name,
			&herb.Description,
			&herb.Price,
			pq.Array(&herb.CulinaryUses),
			&herb.Version,
		)
		if err != nil {
			return nil, Metadata{}, err
		}

		herbs = append(herbs, &herb)
	}

	if err = rows.Err(); err != nil {
		return nil, Metadata{}, err
	}

	metadata := calculateMetadata(totalRecords, filters.Page, filters.PageSize)

	return herbs, metadata, nil
}

func (h HerbModel) Update(herb *Herb) error {

	query := `UPDATE herbs
		SET name = $1, description = $2, price = $3, culinary_uses = $4, version = version + 1
		WHERE id = $5 AND version = $6
		RETURNING version`

	args := []interface{}{
		herb.Name,
		herb.Description,
		herb.Price,
		pq.Array(herb.CulinaryUses),
		herb.ID,
		herb.Version,
	}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := h.DB.QueryRowContext(ctx, query, args...).Scan(&herb.Version)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return ErrEditConflict
		default:
			return err
		}
	}

	return nil
}

func (h HerbModel) Delete(id int64) error {

	if id < 1 {
		return ErrRecordNotFound
	}

	query := `DELETE FROM herbs
		WHERE id = $1`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	result, err := h.DB.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return ErrRecordNotFound
	}

	return nil
}
