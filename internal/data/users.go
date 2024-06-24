package data

import (
	"context"
	"database/sql"
	"errors"
	"log"
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
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	query := `INSERT INTO users ( user_id,fname,lname,email,password_hash,created_at,updated_at,last_login) VALUES 
	( $1,$2,$3,$4,$5,$6,$7,$8 )`
	_, err := m.DB.ExecContext(ctx, query, user_id.String(), u.FirstName, u.LastName, u.Email, u.PasswordHash, time.Now(), time.Now(), time.Now())
	if err != nil {
		return err
	}
	return nil
}

func (m UserModel) Get(user_id string) (*Users, error) {
	query := `
	SELECT user_id,fname,lname,email,password_hash,created_at,updated_at,last_login 
	FROM users 
	WHERE user_id = $1`
	var usr Users
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	err := m.DB.QueryRowContext(ctx, query, user_id).Scan(
		&usr.UserID,
		&usr.FirstName,
		&usr.LastName,
		&usr.Email,
		&usr.PasswordHash,
		&usr.CreatedAt,
		&usr.UpdatedAt,
		&usr.LastLogin,
	)
	log.Println(usr)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrRecordNotFound
		default:
			return nil, err
		}
	}
	return &usr, nil
}

func (m UserModel) Update(u *Users) error {

	query := `
		UPDATE users
		SET fname=$1, lname=$2, email=$3
		WHERE user_id=$4
		RETURNING user_id, fname, lname, email
	`

	var updatedUser Users
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	err := m.DB.QueryRowContext(ctx, query, u.FirstName, u.LastName, u.Email, u.UserID).Scan(&updatedUser.UserID, &updatedUser.FirstName, &updatedUser.LastName, &updatedUser.Email)

	if err != nil {
		if err == sql.ErrNoRows {
			return err
		}
		return err
	}

	//fmt.Printf("data.Update: User updated successfully: %+v\n", updatedUser)
	return nil

}

func (m UserModel) Delete(user_id string) error {
	query := `
	DELETE FROM users WHERE user_id = $1
	`
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	result, err := m.DB.ExecContext(ctx, query, user_id)
	if err != nil {
		return err
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return ErrRecordNotFound
	}
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
