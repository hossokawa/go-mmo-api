package player

import (
	"context"
	"errors"
	"fmt"
)

type PlayerService struct {
	repo PlayerRepository
}

func NewPlayerService(repo PlayerRepository) *PlayerService {
	return &PlayerService{repo: repo}
}

type NotFoundErr struct {
	resource  string
	attribute string
	value     any
}

func (e *NotFoundErr) Error() string {
	return fmt.Sprintf("%s with %s '%v' not found", e.resource, e.attribute, e.value)
}

func (s *PlayerService) CreatePlayer(ctx context.Context, username, class string) (*Player, error) {
	_, err := s.repo.GetPlayerByUsername(ctx, username)
	if err == nil {
		return nil, errors.New("username already in use")
	}

	newPlayer, err := s.repo.CreatePlayer(ctx, CreatePlayerParams{Username: username, Class: class})
	if err != nil {
		return nil, fmt.Errorf("creating new player: %w", err)
	}

	return newPlayer, nil
}

func (s *PlayerService) GetAllPlayers(ctx context.Context) ([]*Player, error) {
	return s.repo.GetAllPlayers(ctx)
}

func (s *PlayerService) GetPlayerByID(ctx context.Context, id int32) (*Player, error) {
	player, err := s.repo.GetPlayerByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("getting player with id %v: %w", id, err)
	}
	if player == nil {
		return nil, &NotFoundErr{resource: "player", attribute: "id", value: id}
	}

	return player, nil
}

func (s *PlayerService) GetPlayerByUsername(ctx context.Context, username string) (*Player, error) {
	player, err := s.repo.GetPlayerByUsername(ctx, username)
	if err != nil {
		return nil, fmt.Errorf("getting player with username '%v': %w", username, err)
	}
	if player == nil {
		return nil, &NotFoundErr{resource: "player", attribute: "username", value: username}
	}

	return player, nil
}

func (s *PlayerService) UpdatePlayerLevel(ctx context.Context, id, level int32) error {
	err := s.repo.UpdatePlayerLevel(ctx, UpdatePlayerLevelParams{ID: id, Level: level})
	if err != nil {
		return fmt.Errorf("updating level for player with id %v: %w", id, err)
	}
	return nil
}

func (s *PlayerService) IncreasePlayerGold(ctx context.Context, id, amount int32) error {
	err := s.repo.IncreasePlayerGold(ctx, UpdatePlayerGoldParams{ID: id, Amount: amount})
	if err != nil {
		return fmt.Errorf("increasing gold for player with id %v: %w", id, err)
	}

	return nil
}

func (s *PlayerService) DecreasePlayerGold(ctx context.Context, id, amount int32) error {
	err := s.repo.DecreasePlayerGold(ctx, UpdatePlayerGoldParams{ID: id, Amount: amount})
	if err != nil {
		return fmt.Errorf("decreasing gold for player with id %v: %w", id, err)
	}

	return nil
}

func (s *PlayerService) DeletePlayerByID(ctx context.Context, id int32) error {
	err := s.repo.DeletePlayerByID(ctx, id)
	if err != nil {
		return fmt.Errorf("deleting player with id %v: %w", id, err)
	}

	return nil
}
