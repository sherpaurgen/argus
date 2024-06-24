package data

import (
	"database/sql"
	"errors"
	"time"
)

// Define a custom ErrRecordNotFound error. We'll return this from our Get() method when
// looking up a movie that doesn't exist in our database.
var (
	ErrRecordNotFound = errors.New("record not found")
)

// Create a Models struct which wraps other models,
// like a UserModel and PermissionMode.
type Models struct {
	ChildrenModel ChildModel
	UserModel     UserModel
	Users         interface {
		Insert(movie *Users) error
		Get(user_id string) (*Users, error)
		Update(movie *Users) error
		Delete(user_id string) error
	}
	SleepDataModel SleepDataModel
	SleepData      interface {
		Insert(sd *SleepData) error
		GetSleepData(child_id string, user_id string, StartSleep time.Time, EndSleep time.Time, limit int, page int)
	}
}

// For ease of use,this New() method which returns a Models struct containing
// the initialized ChildrenModel.
func NewModels(db *sql.DB) Models {
	return Models{
		ChildrenModel: ChildModel{
			DB: db},
		UserModel: UserModel{
			DB: db,
		},
		SleepDataModel: SleepDataModel{DB: db},
	}
}

func NewMockModels() Models {
	return Models{
		Users: MockUserModel{},
	}
}
