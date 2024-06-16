package data

import (
	"database/sql"
	"time"

	"github.com/segmentio/ksuid"
)

type UserModel struct {
	DB *sql.DB
}
type Users struct {
	UserID       string    `json:"user_id" db:"user_id"`
	UserGroup    string    `json:"group_id" db:"group_id"`
	Email        string    `json:"email" db:"email"`
	FirstName    string    `json:"fname" db:"fname"`
	LastName     string    `json:"lname" db:"lname"`
	PasswordHash string    `json:"password_hash" db:"password_hash"`
	CreatedAt    time.Time `json:"created_at" db:"created_at"`
	UpdatedAt    time.Time `json:"updated_at" db:"updated_at"`
	LastLogin    time.Time `json:"last_login" db:"last_login"`
}

func (m UserModel) Insert(u *Users) error {
	user_id := ksuid.New()
	query := `INSERT INTO users ( user_id,fname,lname,email,password_hash,created_at,updated_at,last_login) VALUES 
	( $1,$2,$3,$4,$5,$6,$7,$8 )`
	_, err := m.DB.Exec(query, user_id.String(), u.FirstName, u.LastName, u.Email, u.PasswordHash, time.Now(), time.Now(), time.Now())
	if err != nil {
		return err
	}
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

type MockUserModel struct{}

func (m MockUserModel) Insert(u *Users) error {
	return nil
}
func (m MockUserModel) Get(user_id string) (*Users, error) {
	return nil, nil
}
func (m MockUserModel) Update(u *Users) error {
	return nil
}
func (m MockUserModel) Delete(user_id string) error {
	return nil
}
