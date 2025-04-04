// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.28.0

package database

import (
	"database/sql"
	"time"

	"github.com/google/uuid"
)

type RefreshToken struct {
	Token     string
	CreatedAt time.Time
	UpdatedAt time.Time
	RevokedAt sql.NullTime
	UserID    uuid.UUID
	ExpiresAt time.Time
}

type User struct {
	ID             uuid.UUID
	CreatedAt      time.Time
	UpdatedAt      time.Time
	Email          string
	HashedPassword string
}

type Video struct {
	ID           uuid.UUID
	CreatedAt    time.Time
	UpdatedAt    time.Time
	Title        string
	Description  string
	ThumbnailUrl string
	VideoUrl     string
	UserID       uuid.UUID
}
