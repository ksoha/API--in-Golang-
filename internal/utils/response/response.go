package response

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/go-playground/validator/v10"
)

// creating the response struct
type Response struct {
	Status string
	Error  string
}

// creating teh constants for the response status
const (
	StatusOk    = " OK"
	StatusError = "error"
)

// to send a Json response back to the user
func WriteJSON(w http.ResponseWriter, status int, data interface{}) error {

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	return json.NewEncoder(w).Encode(data)
}

// to handle a general error and return a json message to the user
func GeneralError(err error) Response {
	return Response{
		Status: StatusError,
		Error:  err.Error(),
	}
}

// fucntion to validate the req body
func ValidationError(errs validator.ValidationErrors) Response {
	var errMsgs []string //slice

	//loop through the validation erros and append the error messages to the slices
	for _, err := range errs {
		switch err.ActualTag() { //check the tag of the error
		case "required":
			errMsgs = append(errMsgs, fmt.Sprintf("field %s is required", err.Field()))
		default:
			errMsgs = append(errMsgs, fmt.Sprintf("field %s is invalid", err.Field()))
		}
	}

	return Response{
		Status: StatusError,
		//convert the slice into an error message of type string
		Error: strings.Join(errMsgs, ","),
	}
}
