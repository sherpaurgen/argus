package data

import "time"

type Users struct {
	UserID      int       `json:"user_id" db:"user_id"`
	UserGroup   int       `json:"group_id" db:"group_id"`
	Email       string    `json:"email" db:"email"`
	Description string    `json:"description" db:"description"`
	Secret      string    `json:"secret" db:"secret"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
}
