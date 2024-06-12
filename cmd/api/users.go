package main

import (
	"fmt"
	"github.com/sherpaurgen/argus/internal/validator"
	"net/http"
	"time"

	"github.com/sherpaurgen/argus/internal/data"
)

func (app *application) createUserHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		UserID    string    `json:"user_id" db:"user_id"`
		UserGroup string    `json:"group_id" db:"group_id"`
		Email     string    `json:"email" db:"email"`
		FirstName string    `json:"fname" db:"fname"`
		LastName  string    `json:"lname" db:"lname"`
		Secret    string    `json:"secret" db:"secret"`
		CreatedAt time.Time `json:"created_at" db:"created_at"`
	}

	//err := json.NewDecoder(r.Body).Decode(&input)
	err := app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		//app.errorResponse(w, r, http.StatusBadRequest, err.Error())
		return
	}
	v := validator.New()
	v.Check(input.Email == "", "email", "must be provided")
	v.Check(input.Secret == "", "secret", "must be provided")
	v.Check(input.FirstName == "", "firstname", "must be provided")
	v.Check(input.LastName == "", "lastname", "must be provided")
	if !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
	}
	// %+v will also print key and values
	fmt.Fprintf(w, "%+v\n", input)
}

func (app *application) getUnitHandler(w http.ResponseWriter, r *http.Request) {
	userid, err := app.readIDParam(r)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}
	user := data.Users{
		UserID:    userid,
		UserGroup: "usr100",
		Email:     "usr100@example.com",
		FirstName: "bob",
		LastName:  "rolling",
		Secret:    "swolxo@los23",
		CreatedAt: time.Now(),
	}
	err = app.writeJSON(w, http.StatusOK, envelope{"user": user}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}
