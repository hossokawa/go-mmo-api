package item

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

type ItemRepository interface {
	CreateItem(ctx context.Context, args CreateItemParams) (*Item, error)
	GetAllItems(ctx context.Context) ([]*Item, error)
	GetItemByID(ctx context.Context, id uuid.UUID) (*Item, error)
	GetItemByName(ctx context.Context, name string) (*Item, error)
	UpdateItemValue(ctx context.Context, args UpdateItemValueParams) error
	DeleteItemByID(ctx context.Context, id uuid.UUID) error
}

type pgRepository struct {
	db *pgx.Conn
}

func NewPostgresRepository(db *pgx.Conn) ItemRepository {
	return &pgRepository{db: db}
}

const createItem = `
INSERT INTO item (name, value)
VALUES ($1, $2)
RETURNING id, name, value
`

type CreateItemParams struct {
	Name  string `json:"name"`
	Value int32  `json:"value"`
}

func (r *pgRepository) CreateItem(ctx context.Context, args CreateItemParams) (*Item, error) {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return nil, fmt.Errorf("beginning transaction: %w", err)
	}
	defer tx.Rollback(ctx)

	var i Item

	row := tx.QueryRow(ctx, createItem, args.Name, args.Value)
	err = row.Scan(
		&i.ID,
		&i.Name,
		&i.Value,
	)
	if err != nil {
		return nil, fmt.Errorf("scanning row into item struct: %w", err)
	}

	if err = tx.Commit(ctx); err != nil {
		return nil, fmt.Errorf("commiting transaction: %w", err)
	}

	return &i, nil
}

const getAllItems = `
SELECT id, name, value FROM item
`

func (r *pgRepository) GetAllItems(ctx context.Context) ([]*Item, error) {
	rows, err := r.db.Query(ctx, getAllItems)
	if err != nil {
		return nil, fmt.Errorf("querying for all items: %w", err)
	}
	defer rows.Close()

	var items []*Item

	for rows.Next() {
		var i Item

		err = rows.Scan(
			&i.ID,
			&i.Name,
			&i.Value,
		)
		if err != nil {
			return nil, fmt.Errorf("scanning rows into item struct: %w", err)
		}

		items = append(items, &i)
	}

	return items, nil
}

const getItemByID = `
SELECT id, name, value FROM item WHERE id = $1
`

func (r *pgRepository) GetItemByID(ctx context.Context, id uuid.UUID) (*Item, error) {
	var i Item

	row := r.db.QueryRow(ctx, getItemByID, id)
	err := row.Scan(
		&i.ID,
		&i.Name,
		&i.Value,
	)
	if err != nil {
		return nil, fmt.Errorf("scanning row into item struct: %w", err)
	}

	return &i, nil
}

const getItemByName = `
SELECT id, name, value FROM item WHERE name = $1
`

func (r *pgRepository) GetItemByName(ctx context.Context, name string) (*Item, error) {
	var i Item

	row := r.db.QueryRow(ctx, getItemByName, name)
	err := row.Scan(
		&i.ID,
		&i.Name,
		&i.Value,
	)
	if err != nil {
		return nil, fmt.Errorf("scanning row into item struct: %w", err)
	}

	return &i, nil
}

const updateItemValue = `
UPDATE item SET value = $2 WHERE id = $1
`

type UpdateItemValueParams struct {
	ID       uuid.UUID `json:"id"`
	NewValue int32     `json:"new_value"`
}

func (r *pgRepository) UpdateItemValue(ctx context.Context, args UpdateItemValueParams) error {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return fmt.Errorf("beginning transaction: %w", err)
	}
	defer tx.Rollback(ctx)

	_, err = tx.Exec(ctx, updateItemValue, args.ID, args.NewValue)
	if err != nil {
		return fmt.Errorf("updating item value: %w", err)
	}

	if err = tx.Commit(ctx); err != nil {
		return fmt.Errorf("commiting transaction: %w", err)
	}

	return nil
}

const deleteItemByID = `
DELETE FROM item WHERE id = $1
`

func (r *pgRepository) DeleteItemByID(ctx context.Context, id uuid.UUID) error {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return fmt.Errorf("beginning transaction: %w", err)
	}
	defer tx.Rollback(ctx)

	_, err = tx.Exec(ctx, deleteItemByID, id)
	if err != nil {
		return fmt.Errorf("deleting item: %w", err)
	}

	if err = tx.Commit(ctx); err != nil {
		return fmt.Errorf("commiting transaction: %w", err)
	}

	return nil
}
