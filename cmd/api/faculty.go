package main

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/liamgluna/daycare-server/internal/data"
)

func (app *application) createFacultyHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		FirstName string `json:"first_name"`
		LastName  string `json:"last_name"`
		Email     string `json:"email"`
		Contact   string  `json:"contact"`
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

	err = app.writeJSON(w, http.StatusCreated, envelope{"faculty": faculty}, headers)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) updateFacultyHandler(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil || id < 1 {
		app.notFoundResponse(w, r)
		return
	}

	faculty, err := app.models.Faculty.Get(id)
	if err != nil {
		if errors.Is(err, data.ErrRecordNotFound) {
			app.notFoundResponse(w, r)
		} else {
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	var input struct {
		FirstName *string `json:"first_name"`
		LastName  *string `json:"last_name"`
		Email     *string `json:"email"`
		Contact   *string  `json:"contact"`
		Position  *string `json:"position"`
	}

	err = app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	if input.FirstName != nil {
		faculty.FirstName = *input.FirstName
	}

	if input.LastName != nil {
		faculty.LastName = *input.LastName
	}

	if input.Email != nil {
		faculty.Email = *input.Email
	}

	if input.Contact != nil {
		faculty.Contact = *input.Contact
	}

	if input.Position != nil {
		faculty.Position = *input.Position
	}

	err = app.models.Faculty.Update(faculty)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"faculty": faculty}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) showFacultyHandler(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil || id < 1 {
		app.notFoundResponse(w, r)
		return
	}

	faculty, err := app.models.Faculty.Get(id)
	if err != nil {
		if errors.Is(err, data.ErrRecordNotFound) {
			app.notFoundResponse(w, r)
		} else {
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"faculty": faculty}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}