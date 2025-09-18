package player

import "time"

type Player struct {
	ID        int32     `json:"id"`
	Username  string    `json:"username"`
	Class     string    `json:"class"`
	Level     int32     `json:"level"`
	Gold      int32     `json:"gold"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
