package main

import (
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
		Email:        input.UserGroup,
		FirstName:    input.FirstName,
		LastName:     input.LastName,
		PasswordHash: input.PasswordHash,
		CreatedAt:    input.CreatedAt,
		UpdatedAt:    input.CreatedAt,
		LastLogin:    input.LastLogin,
	}
	err = app.models.UserOb.Insert(usr)
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

func (app *application) getUnitHandler(w http.ResponseWriter, r *http.Request) {
	userid, err := app.readIDParam(r)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}
	user := data.Users{
		UserID:       userid,
		UserGroup:    "usr100",
		Email:        "usr100@example.com",
		FirstName:    "bob",
		LastName:     "rolling",
		PasswordHash: "swolxo@los23",
		CreatedAt:    time.Now(),
	}
	err = app.writeJSON(w, http.StatusOK, envelope{"user": user}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}
