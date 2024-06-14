package data

import (
	"database/sql"
	"time"
)

type UserModel struct {
	DB *sql.DB
}
type Users struct {
	UserID    string    `json:"user_id" db:"user_id"`
	UserGroup string    `json:"group_id" db:"group_id"`
	Email     string    `json:"email" db:"email"`
	FirstName string    `json:"fname" db:"fname"`
	LastName  string    `json:"lname" db:"lname"`
	Secret    string    `json:"secret" db:"secret"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
}

func (m UserModel) Insert(u *Users) error {
	return nil
}

func (m UserModel) Get(user_id string) (*Users, error) {
	return nil, nil
}

func (m UserModel) Update(u *Users) error {
	return nil
}

func (m UserModel) Delete(user_id string) error {
	return nil
}
