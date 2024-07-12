package main

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/sherpaurgen/argus/internal/data"
	"github.com/sherpaurgen/argus/internal/validator"
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
	fmt.Printf("sdata obj: %+v \n", sdata)
	fmt.Println(sdata.StartSleep.Equal(time.Time{}))
	v := validator.New()
	v.IsTimeEmptyCheck(sdata.StartSleep.Equal(time.Time{}), "start_sleep", "must be provided")
	v.IsTimeEmptyCheck(sdata.EndSleep.Equal(time.Time{}), "end_sleep", "must be provided")
	v.Check(sdata.StartSleep.Before(sdata.EndSleep), "end_sleep", "must be after start_sleep")
	v.Check(sdata.DeviceID != "", "device_id", "must be provided")
	v.Check(sdata.UserID != "", "user_id", "must be provided")
	v.Check(sdata.ChildId != "", "child_id", "must be provided")
	if !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
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

func (app *application) getSleepDataHandler(w http.ResponseWriter, r *http.Request) {
	//input example.com/v1/sleepdata/ch112?start_sleep=2024-06-22T01:20:00Z&end_sleep=2024-06-22T05:20:00Z&limit=5&page=1

	child_id, err := app.readIDParam(r)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	queryparams := r.URL.Query()
	raw_start_sleep := queryparams.Get("start_sleep")
	raw_end_sleep := queryparams.Get("end_sleep")
	limit_string := queryparams.Get("limit")
	page_string := queryparams.Get("page")

	page := 1
	limit := 360
	//considering metric is collected every 1 minute default limit to show data for last 6 hours 6*60=360 metric points

	if page_string != "" || limit_string != "" {
		var err1, err2 error
		page, err1 = strconv.Atoi(page_string)
		limit, err2 = strconv.Atoi(limit_string)

		if err1 != nil || err2 != nil {
			app.badRequestResponse(w, r, fmt.Errorf("invalid page or limit"))
			return
		}
	}

	user_id := "2iFGw10P8Ne8jsFjx2wgjCuKJOI"
	start_sleep, err3 := time.Parse(time.RFC3339, raw_start_sleep)
	end_sleep, err4 := time.Parse(time.RFC3339, raw_end_sleep)
	if err3 != nil || err4 != nil {
		app.badRequestResponse(w, r, fmt.Errorf("invalid start_sleep or end_sleep"))
		return
	}

	sleepdata_response, err := app.models.SleepDataModel.GetSleepData(child_id, user_id, start_sleep, end_sleep, page, limit)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	headers := make(http.Header)
	headers.Set("Location", fmt.Sprintf("/v1/sleepdata/%s", child_id))
	err = app.writeJSON(w, http.StatusCreated, envelope{"sleepdata": sleepdata_response}, headers)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}
