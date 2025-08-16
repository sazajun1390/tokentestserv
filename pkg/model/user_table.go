package model

import (
	"time"

	"github.com/uptrace/bun"
)

type User struct {
	bun.BaseModel `bun:"table:users,alias:u"`
	ID            int64     `bun:"id,pk,autoincrement"`
	CreatedAt     time.Time `bun:"created_at,notnull"`
}

type UserProfile struct {
	bun.BaseModel `bun:"table:user_profiles,alias:up"`
	ID            int64        `bun:"id,pk,autoincrement"`
	UserID        int64        `bun:"user_id,notnull,unique"`
	CreatedAt     time.Time    `bun:"created_at,notnull"`
	UpdatedAt     time.Time    `bun:"updated_at,notnull"`
	DeletedAt     bun.NullTime `bun:"deleted_at"`
	Email         string       `bun:"email,unique,notnull"`
	Password      string       `bun:"password,notnull"`
	Tel           string       `bun:"tel"`
}
