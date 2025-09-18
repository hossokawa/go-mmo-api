package inventory

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/hossokawa/go-nethttp-example/internal/item"
)

type InventoryService struct {
	repo InventoryRepository
}

func NewInventoryService(repo InventoryRepository) *InventoryService {
	return &InventoryService{repo: repo}
}

func (s *InventoryService) AddItem(ctx context.Context, playerID int32, itemID uuid.UUID) error {
	err := s.repo.AddItem(ctx, AddItemParams{PlayerID: playerID, ItemID: itemID})
	if err != nil {
		return fmt.Errorf("adding item with id %v to inventory of player with id %v: %w", itemID, playerID, err)
	}

	return nil
}

func (s *InventoryService) ListPlayerItems(ctx context.Context, playerID int32) ([]item.Item, error) {
	items, err := s.repo.ListPlayerItems(ctx, playerID)
	if err != nil {
		return nil, fmt.Errorf("getting items for player with id %v: %w", playerID, err)
	}

	return items, nil
}

func (s *InventoryService) RemoveItem(ctx context.Context, playerID int32, itemID uuid.UUID) error {
	err := s.repo.RemoveItem(ctx, RemoveItemParams{PlayerID: playerID, ItemID: itemID})
	if err != nil {
		return fmt.Errorf("removing item with id %v from inventory of player with id %v: %w", itemID, playerID, err)
	}

	return nil
}
