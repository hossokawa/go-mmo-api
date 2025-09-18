package item

import "github.com/google/uuid"

type Item struct {
	ID    uuid.UUID `json:"id"`
	Name  string    `json:"name"`
	Value int32     `json:"value"`
}
