package types

import "time"

// UserResultAddress ...
type UserResultAddress struct {
	Address  string    `json:"address"`
	Asset    string    `json:"asset"`
	CreateAt time.Time `json:"createAt"`
	CybName  string    `json:"cybName"`
}
