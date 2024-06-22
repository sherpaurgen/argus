package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/sherpaurgen/argus/internal/data"
)

type SleepData struct {
	StartSleep time.Time `json:"start_sleep" db:"start_sleep"`
	EndSleep   time.Time `json:"end_sleep" db:"end_sleep"`
	DeviceID   string    `json:"device_id" db:"device_id"`
	UserID     string    `json:"user_id" db:"user_id"`
	ChildId    string    `json:"child_id" db:"child_id"`
}

func (app *application) addsleepdataHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		StartSleep   time.Time `json:"start_sleep" db:"start_sleep"`
		EndSleep     time.Time `json:"end_sleep" db:"end_sleep"`
		DeviceID     string    `json:"device_id" db:"device_id"`
		UserID       string    `json:"user_id" db:"user_id"`
		ChildId      string    `json:"child_id" db:"child_id"`
		SleepQuality int       `json:"sleep_quality" db:"sleep_quality"`
	}
	err := app.readJSON(w, r, &input)
	input.UserID = "2iFGw10P8Ne8jsFjx2wgjCuKJOI"
	input.DeviceID = "Nanit Pro1"
	input.ChildId = "e79d0004947a20fa8f0cb0a8326"
	fmt.Println(err)
	if err != nil {
		app.badRequestResponse(w, r, err)
		//app.errorResponse(w, r, http.StatusBadRequest, err.Error())
		return
	}
	sdata := &data.SleepData{
		StartSleep: input.StartSleep,
		EndSleep:   input.EndSleep,
		DeviceID:   input.DeviceID,
		UserID:     input.UserID,
		ChildId:    input.ChildId,
	}
	err = app.models.SleepDataModel.Insert(sdata)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
	headers := make(http.Header)
	headers.Set("Location", fmt.Sprintf("/v1/sleepdata/%s", sdata.ChildId))
	err = app.writeJSON(w, http.StatusCreated, envelope{"sleepdata": "data added"}, headers)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}
