package data

import (
	"database/sql"
	"errors"
)

// Define a custom ErrRecordNotFound error. We'll return this from our Get() method when
// looking up a movie that doesn't exist in our database.
var (
	ErrRecordNotFound = errors.New("record not found")
)

// Create a Models struct which wraps other models,
// like a UserModel and PermissionMode.
type Models struct {
	ChildrenOb ChildModel
	UserOb     UserModel
	Users      interface {
		Insert(movie *Users) error
		Get(user_id string) (*Users, error)
		Update(movie *Users) error
		Delete(user_id string) error
	}
}

// For ease of use,this New() method which returns a Models struct containing
// the initialized ChildrenModel.
func NewModels(db *sql.DB) Models {
	return Models{
		ChildrenOb: ChildModel{
			DB: db},
		UserOb: UserModel{
			DB: db,
		},
	}
}

func NewMockModels() Models {
	return Models{
		Users: MockUserModel{},
	}
}
