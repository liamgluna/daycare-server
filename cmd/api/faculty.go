package main

import (
	"fmt"
	"net/http"

	"github.com/liamgluna/daycare-server/internal/data"
)

func (app *application) createFacultyHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		FirstName string `json:"first_name"`
		LastName  string `json:"last_name"`
		Email     string `json:"email"`
		Contact   int64  `json:"contact"`
		Position  string `json:"position"`
	}

	err := app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	faculty := &data.Faculty{
		FirstName: input.FirstName,
		LastName:  input.LastName,
		Email:     input.Email,
		Contact:   input.Contact,
		Position:  input.Position,
	}

	err = app.models.Faculty.Insert(faculty)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	headers := make(http.Header)
	headers.Set("Location", fmt.Sprintf("/faculty/%d", faculty.FacultyID))

	err = app.writeJSON(w, http.StatusCreated, envelope{"faculty": faculty}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}
