package inventory

import "github.com/google/uuid"

type Inventory struct {
	PlayerID int32     `json:"player_id"`
	ItemID   uuid.UUID `json:"item_id"`
}
