package player

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
)

type PlayerRepository interface {
	CreatePlayer(ctx context.Context, args CreatePlayerParams) (*Player, error)
	GetAllPlayers(ctx context.Context) ([]*Player, error)
	GetPlayerByID(ctx context.Context, id int32) (*Player, error)
	GetPlayerByUsername(ctx context.Context, username string) (*Player, error)
	UpdatePlayerLevel(ctx context.Context, args UpdatePlayerLevelParams) error
	IncreasePlayerGold(ctx context.Context, args UpdatePlayerGoldParams) error
	DecreasePlayerGold(ctx context.Context, args UpdatePlayerGoldParams) error
	DeletePlayerByID(ctx context.Context, id int32) error
}

type pgRepository struct {
	db *pgx.Conn
}

func NewPostgresRepository(db *pgx.Conn) PlayerRepository {
	return &pgRepository{db: db}
}

const createPlayer = `
INSERT INTO player (username, class, level, gold, created_at, updated_at)
VALUES ($1, $2, 1, 0, now(), now())
RETURNING id, username, class, level, gold, created_at, updated_at
`

type CreatePlayerParams struct {
	Username string `json:"username"`
	Class    string `json:"class"`
}

func (r *pgRepository) CreatePlayer(ctx context.Context, args CreatePlayerParams) (*Player, error) {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return nil, fmt.Errorf("beginning transaction: %w", err)
	}
	defer tx.Rollback(ctx)

	var p Player

	row := tx.QueryRow(ctx, createPlayer, args.Username, args.Class)
	err = row.Scan(
		&p.ID,
		&p.Username,
		&p.Class,
		&p.Level,
		&p.Gold,
		&p.CreatedAt,
		&p.UpdatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("scanning row into player struct: %w", err)
	}

	if err = tx.Commit(ctx); err != nil {
		return nil, fmt.Errorf("commiting transaction: %w", err)
	}

	return &p, nil
}

const getAllPlayers = `
SELECT id, username, class, level, gold, created_at, updated_at FROM player
`

func (r *pgRepository) GetAllPlayers(ctx context.Context) ([]*Player, error) {
	rows, err := r.db.Query(ctx, getAllPlayers)
	if err != nil {
		return nil, fmt.Errorf("querying for all players: %w", err)
	}
	defer rows.Close()

	var ps []*Player

	for rows.Next() {
		var p Player

		err = rows.Scan(
			&p.ID,
			&p.Username,
			&p.Class,
			&p.Level,
			&p.Gold,
			&p.CreatedAt,
			&p.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("scanning rows into player struct: %w", err)
		}

		ps = append(ps, &p)
	}

	return ps, nil
}

const getPlayerByID = `
SELECT id, username, class, level, gold, created_at, updated_at FROM player WHERE id = $1
`

func (r *pgRepository) GetPlayerByID(ctx context.Context, id int32) (*Player, error) {
	var p Player

	row := r.db.QueryRow(ctx, getPlayerByID, id)
	err := row.Scan(
		&p.ID,
		&p.Username,
		&p.Class,
		&p.Level,
		&p.Gold,
		&p.CreatedAt,
		&p.UpdatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("scanning row into player struct: %w", err)
	}

	return &p, nil
}

const getPlayerByUsername = `
SELECT id, username, class, level, gold, created_at, updated_at FROM player WHERE username = $1
`

func (r *pgRepository) GetPlayerByUsername(ctx context.Context, username string) (*Player, error) {
	var p Player

	row := r.db.QueryRow(ctx, getPlayerByUsername, username)
	err := row.Scan(
		&p.ID,
		&p.Username,
		&p.Class,
		&p.Level,
		&p.Gold,
		&p.CreatedAt,
		&p.UpdatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("scanning row into player struct: %w", err)
	}

	return &p, nil
}

const updatePlayerLevel = `
UPDATE player SET level = $2 WHERE id = $1
`

type UpdatePlayerLevelParams struct {
	ID    int32 `json:"id"`
	Level int32 `json:"level"`
}

func (r *pgRepository) UpdatePlayerLevel(ctx context.Context, args UpdatePlayerLevelParams) error {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return fmt.Errorf("beginning transaction: %w", err)
	}
	defer tx.Rollback(ctx)

	_, err = tx.Exec(ctx, updatePlayerLevel, args.ID, args.Level)
	if err != nil {
		return fmt.Errorf("updating player level: %w", err)
	}

	if err = tx.Commit(ctx); err != nil {
		return fmt.Errorf("commiting transaction: %w", err)
	}

	return nil
}

type UpdatePlayerGoldParams struct {
	ID     int32 `json:"id"`
	Amount int32 `json:"amount"`
}

const increasePlayerGold = `
UPDATE player SET gold = gold + $2 WHERE id = $1
`

func (r *pgRepository) IncreasePlayerGold(ctx context.Context, args UpdatePlayerGoldParams) error {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return fmt.Errorf("beginning transaction: %w", err)
	}
	defer tx.Rollback(ctx)

	_, err = r.db.Exec(ctx, increasePlayerGold, args.ID, args.Amount)
	if err != nil {
		return fmt.Errorf("increasing player gold: %w", err)
	}

	if err = tx.Commit(ctx); err != nil {
		return fmt.Errorf("commiting transaction: %w", err)
	}

	return nil
}

const decreasePlayerGold = `
UPDATE player SET gold = gold - $2 WHERE id = $1
`

func (r *pgRepository) DecreasePlayerGold(ctx context.Context, args UpdatePlayerGoldParams) error {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return fmt.Errorf("beginning transaction: %w", err)
	}
	defer tx.Rollback(ctx)

	_, err = tx.Exec(ctx, decreasePlayerGold, args.ID, args.Amount)
	if err != nil {
		return fmt.Errorf("decreasing player gold: %w", err)
	}

	if err = tx.Commit(ctx); err != nil {
		return fmt.Errorf("commiting transaction: %w", err)
	}

	return nil
}

const deletePlayerByID = `
DELETE FROM player WHERE id = $1
`

func (r *pgRepository) DeletePlayerByID(ctx context.Context, id int32) error {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return fmt.Errorf("beginning transaction: %w", err)
	}
	defer tx.Rollback(ctx)

	_, err = tx.Exec(ctx, deletePlayerByID, id)
	if err != nil {
		return fmt.Errorf("deleting player: %w", err)
	}

	if err = tx.Commit(ctx); err != nil {
		return fmt.Errorf("commiting transaction: %w", err)
	}

	return nil
}
