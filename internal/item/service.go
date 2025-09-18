package item

import (
	"context"
	"fmt"

	"github.com/google/uuid"
)

type ItemService struct {
	repo ItemRepository
}

func NewItemService(repo ItemRepository) *ItemService {
	return &ItemService{repo: repo}
}

func (s *ItemService) CreateItem(ctx context.Context, name string, value int32) (*Item, error) {
	_, err := s.repo.GetItemByName(ctx, name)
	if err == nil {
		return nil, fmt.Errorf("item with name '%v' already exists", name)
	}

	newItem, err := s.repo.CreateItem(ctx, CreateItemParams{Name: name, Value: value})
	if err != nil {
		return nil, fmt.Errorf("creating new item: %w", err)
	}

	return newItem, nil
}

func (s *ItemService) GetAllItems(ctx context.Context) ([]*Item, error) {
	return s.repo.GetAllItems(ctx)
}

func (s *ItemService) GetItemByID(ctx context.Context, id uuid.UUID) (*Item, error) {
	item, err := s.repo.GetItemByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("getting item with id %v: %w", id, err)
	}

	return item, nil
}

func (s *ItemService) GetItemByName(ctx context.Context, name string) (*Item, error) {
	item, err := s.repo.GetItemByName(ctx, name)
	if err != nil {
		return nil, fmt.Errorf("getting item with name '%v': %w", name, err)
	}

	return item, nil
}

func (s *ItemService) UpdateItemValue(ctx context.Context, id uuid.UUID, newValue int32) error {
	err := s.repo.UpdateItemValue(ctx, UpdateItemValueParams{ID: id, NewValue: newValue})
	if err != nil {
		return fmt.Errorf("updating value for item with id %v: %w", id, err)
	}

	return nil
}

func (s *ItemService) DeleteItemByID(ctx context.Context, id uuid.UUID) error {
	err := s.repo.DeleteItemByID(ctx, id)
	if err != nil {
		return fmt.Errorf("deleting item with id %v: %w", id, err)
	}

	return nil
}
