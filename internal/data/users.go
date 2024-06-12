package data

import "time"

type Users struct {
	UserID    string    `json:"user_id" db:"user_id"`
	UserGroup string    `json:"group_id" db:"group_id"`
	Email     string    `json:"email" db:"email"`
	FirstName string    `json:"fname" db:"fname"`
	LastName  string    `json:"lname" db:"lname"`
	Secret    string    `json:"secret" db:"secret"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
}
