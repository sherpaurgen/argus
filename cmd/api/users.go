package main

import (
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/sherpaurgen/argus/internal/validator"

	"github.com/sherpaurgen/argus/internal/data"
)

func (app *application) createUserHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
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

	//err := json.NewDecoder(r.Body).Decode(&input)
	err := app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		//app.errorResponse(w, r, http.StatusBadRequest, err.Error())
		return
	}

	v := validator.New()
	v.Check(input.Email != "", "email", "must be provided")
	v.Check(input.PasswordHash != "", "secret", "must be provided")
	v.Check(input.FirstName != "", "firstname", "must be provided")
	v.Check(input.LastName != "", "lastname", "must be provided")
	if !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
	}
	usr := &data.Users{
		UserID:       input.UserID,
		UserGroup:    input.UserGroup,
		Email:        input.Email,
		FirstName:    input.FirstName,
		LastName:     input.LastName,
		PasswordHash: input.PasswordHash,
		CreatedAt:    input.CreatedAt,
		UpdatedAt:    input.CreatedAt,
		LastLogin:    input.LastLogin,
	}
	err = app.models.UserModel.Insert(usr)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
	// %+v will also print key and values
	//fmt.Fprintf(w, "%+v\n", input)
	headers := make(http.Header)
	headers.Set("Location", fmt.Sprintf("/v1/users/%s", usr.UserID))
	err = app.writeJSON(w, http.StatusCreated, envelope{"users": usr}, headers)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}

}

func (app *application) getUserHandler(w http.ResponseWriter, r *http.Request) {
	user_id, err := app.readIDParam(r)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}
	usr, err := app.models.UserModel.Get(user_id)

	if err != nil {
		app.notFoundResponse(w, r)
		return
	}
	err = app.writeJSON(w, http.StatusOK, envelope{"user": usr}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}

}
func (app *application) deleteUserHandler(w http.ResponseWriter, r *http.Request) {
	user_id, err := app.readIDParam(r)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}
	err = app.models.UserModel.Delete(user_id)

	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}
	// Return a 200 OK status code along with a success message.
	err = app.writeJSON(w, http.StatusOK, envelope{"message": "movie successfully deleted"}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}

}

func (app *application) updateUserHandler(w http.ResponseWriter, r *http.Request) {
	user_id, err := app.readIDParam(r) //get the user id from url
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}
	usr, err := app.models.UserModel.Get(user_id) //fetching 1 matched user from db
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}
	var input struct {
		Email     *string `json:"email" db:"email"`
		FirstName *string `json:"fname" db:"fname"`
		LastName  *string `json:"lname" db:"lname"`
	}
	err = app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}
	usr.UserID = user_id
	if input.Email != nil {
		usr.Email = *input.Email
	}
	if input.FirstName != nil {
		usr.LastName = *input.LastName
	}

	if input.LastName != nil {
		usr.LastName = *input.LastName
	}

	v := validator.New()

	v.Check(usr.Email != "", "email", "must be provided")
	v.Check(usr.FirstName != "", "firstname", "must be provided")
	v.Check(usr.LastName != "", "lastname", "must be provided")
	if !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}
	err = app.models.UserModel.Update(usr)
	if err != nil {
		app.logger.Printf("data.update: throws error %v", err)
		app.serverErrorResponse(w, r, err)
	}
	err = app.writeJSON(w, http.StatusOK, envelope{"user": "updateduser"}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}
