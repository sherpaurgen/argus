package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"

	"github.com/go-chi/chi/v5"
)

type envelope map[string]any

func (app *application) readIDParam(r *http.Request) (string, error) {
	id := chi.URLParam(r, "id")
	//id, err := strconv.ParseInt(id_parameter, 10, 64)
	//if err != nil {
	//	return 0, err
	//}
	//app.logger.Printf("Error parsing id from url: %v", user_id)
	return id, nil
}

// here "any" is equivalent to interface{} is used to accept any type of data
// the status is the status code written to the response
// the headers is map of key value pairs
func (app *application) writeJSON(w http.ResponseWriter, status int, data envelope, headers http.Header) error {
	js, err := json.Marshal(data)
	if err != nil {
		return err
	}
	for key, value := range headers {
		w.Header()[key] = value
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	w.Write(js)
	return nil
}

func (app *application) readJSON(w http.ResponseWriter, r *http.Request, dest interface{}) error {
	//dest destination
	const maxBytes = 1048576
	r.Body = http.MaxBytesReader(w, r.Body, int64(maxBytes))
	//err := json.NewDecoder(r.Body).Decode(dest)
	dcoder := json.NewDecoder(r.Body)
	dcoder.DisallowUnknownFields() //allows decoder to return error when input doesnt match target struct
	err := dcoder.Decode(dest)
	if err != nil {
		fmt.Println(err)
		// while reading json from reqest body there is likely to happen errors below
		var syntaxError *json.SyntaxError
		var unmarshallTypeError *json.UnmarshalTypeError
		var invalidUnmarshalError *json.InvalidUnmarshalError
		var maxBytesError *http.MaxBytesError
		switch {
		case errors.As(err, &syntaxError):
			return fmt.Errorf("readJSON: bad json body invalid character %d", syntaxError.Offset)
		case errors.Is(err, io.ErrUnexpectedEOF):
			return errors.New("readJSON: body contains badly formed json")
		case errors.As(err, &unmarshallTypeError):
			if unmarshallTypeError.Field != "" {
				return fmt.Errorf("readJSON: body contains incorrect json type for field %q", unmarshallTypeError.Field)
			}
			return fmt.Errorf("readJSON: body contains incorrect type at character %d", unmarshallTypeError.Offset)
		case errors.Is(err, io.EOF):
			return errors.New("readJSON: body must not be empty")
		case errors.As(err, &invalidUnmarshalError):
			panic(err)
		case errors.As(err, &maxBytesError):
			return fmt.Errorf("readJSON: body must not be larger than 1MB")
		default:
			return fmt.Errorf("readJSON: Got error %w", err)
		}
	}
	return nil
}

func (app *application) badRequestResponse(w http.ResponseWriter, r *http.Request, err error) {
	app.errorResponse(w, r, http.StatusBadRequest, err.Error())
}
