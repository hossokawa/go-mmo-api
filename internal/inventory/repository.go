package inventory

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/hossokawa/go-nethttp-example/internal/item"
	"github.com/jackc/pgx/v5"
)

type InventoryRepository interface {
	AddItem(ctx context.Context, args AddItemParams) error
	ListPlayerItems(ctx context.Context, playerID int32) ([]item.Item, error)
	RemoveItem(ctx context.Context, args RemoveItemParams) error
}

type pgRepository struct {
	db *pgx.Conn
}

func NewPostgresRepository(db *pgx.Conn) InventoryRepository {
	return &pgRepository{db: db}
}

const addItem = `
INSERT INTO inventory (player_id, item_id)
VALUES ($1, $2)
`

type AddItemParams struct {
	PlayerID int32     `json:"player_id"`
	ItemID   uuid.UUID `json:"item_id"`
}

func (r *pgRepository) AddItem(ctx context.Context, args AddItemParams) error {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return fmt.Errorf("beginning transaction: %w", err)
	}
	defer tx.Rollback(ctx)

	_, err = tx.Exec(ctx, addItem, args.PlayerID, args.ItemID)
	if err != nil {
		return fmt.Errorf("adding item to player's inventory: %w", err)
	}

	if err = tx.Commit(ctx); err != nil {
		return fmt.Errorf("commiting transaction: %w", err)
	}

	return nil
}

const listPlayerItems = `
SELECT item.id, item.name, item.value
FROM inventory
JOIN item ON item.id = item_id
WHERE player_id = $1
`

func (r *pgRepository) ListPlayerItems(ctx context.Context, playerID int32) ([]item.Item, error) {
	rows, err := r.db.Query(ctx, listPlayerItems, playerID)
	if err != nil {
		return nil, fmt.Errorf("getting all items for player: %w", err)
	}
	defer rows.Close()

	var items []item.Item

	for rows.Next() {
		var i item.Item

		if err := rows.Scan(&i.ID, &i.Name, &i.Value); err != nil {
			return nil, fmt.Errorf("scanning rows from inventory into item struct: %w", err)
		}

		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return items, nil
}

const removeItem = `
DELETE FROM inventory
WHERE player_id = $1
AND item_id = $2
`

type RemoveItemParams struct {
	PlayerID int32     `json:"player_id"`
	ItemID   uuid.UUID `json:"item_id"`
}

func (r *pgRepository) RemoveItem(ctx context.Context, args RemoveItemParams) error {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return fmt.Errorf("beginning transaction: %w", err)
	}
	defer tx.Rollback(ctx)

	_, err = tx.Exec(ctx, removeItem, args.PlayerID, args.ItemID)
	if err != nil {
		return fmt.Errorf("removing item from player's inventory: %w", err)
	}

	if err = tx.Commit(ctx); err != nil {
		return fmt.Errorf("commiting transaction: %w", err)
	}

	return nil
}
