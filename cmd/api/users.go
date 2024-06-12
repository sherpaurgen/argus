package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/sherpaurgen/Tilicho_v1/internal/data"
)

func (app *application) createUserHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		UnitID      int       `json:"unit_id" db:"unit_id"`
		BuildingID  string    `json:"building_id" db:"building_id"`
		OwnerID     int       `json:"owner_id" db:"owner_id"`
		Title       string    `json:"title" db:"title"`
		Description string    `json:"description" db:"description"`
		FloorNumber int       `json:"floor_number" db:"floor_number"`
		PricePerDay float64   `json:"price_per_day,omitempty" db:"price_per_day"`
		CreatedAt   time.Time `json:"created_at" db:"created_at"`
		Features    []string  `json:"generes"`
	}

	//err := json.NewDecoder(r.Body).Decode(&input)
	err := app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		//app.errorResponse(w, r, http.StatusBadRequest, err.Error())
		return
	}
	fmt.Fprintf(w, "%+v\n", input)
}

func (app *application) getUnitHandler(w http.ResponseWriter, r *http.Request) {
	unitid, err := app.readIDParam(r)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}
	unit := data.Unit{
		UnitID:      int(unitid),
		BuildingID:  "B17",
		OwnerID:     int(unitid),
		Title:       "titelsdf",
		Description: "asdf",
		FloorNumber: 123,
		PricePerDay: 232,
		CreatedAt:   time.Now(),
	}
	err = app.writeJSON(w, http.StatusOK, envelope{"unit": unit}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}
