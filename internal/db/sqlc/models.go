// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.24.0

package db

import (
	"time"
)

type Profile struct {
	WalletAddress string    `json:"wallet_address"`
	GamerTag      string    `json:"gamer_tag"`
	CreatedAt     time.Time `json:"created_at"`
}