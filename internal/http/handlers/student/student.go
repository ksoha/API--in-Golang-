package student

import (
	"encoding/json"
	"errors"
	"io"
	"log/slog"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/ksoha/API-in-Golang/internal/storage"
	"github.com/ksoha/API-in-Golang/internal/types"
	"github.com/ksoha/API-in-Golang/internal/utils/response"
)

// receiving the storage interface as a dependency injection
func New(storage storage.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		slog.Info("creating a new student")

		//to serealize the data sent by the client we use the json packge

		var student types.Student
		//encoding json
		err := json.NewDecoder(r.Body).Decode(&student)
		//first check the explicit error type

		//io.EOF checks if the request body is empty
		if errors.Is(err, io.EOF) {
			response.WriteJSON(w, http.StatusBadRequest, response.GeneralError(err))
			return //returning from the function so that the rest of code is not executed
		}

		//check is it is a general error
		if err != nil {
			response.WriteJSON(w, http.StatusBadRequest, response.GeneralError(err))
			return
		}

		//request validation
		//for production application we always validate the request
		if err := validator.New().Struct(student); err != nil {
			//typecast the error to validator.ValidationErrors type before passing it
			validationErrs := err.(validator.ValidationErrors)

			response.WriteJSON(w, http.StatusBadRequest, response.ValidationError(validationErrs))
			return
		}

		//create the student in the database
		lastid, err := storage.CreateStudent(
			student.Name,
			student.Email,
			student.Age,
		)

		slog.Info("user created successfully", slog.Int64("id", lastid))

		if err != nil {
			response.WriteJSON(w, http.StatusInternalServerError, err)
			return
		}

		response.WriteJSON(w, http.StatusAccepted, map[string]int64{"id": lastid})
	}
}
