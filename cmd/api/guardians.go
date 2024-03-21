package main

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/liamgluna/daycare-server/internal/data"
)

func (app *application) createGuardianHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		FirstName    string `json:"first_name"`
		LastName     string `json:"last_name"`
		Gender       string `json:"gender"`
		Relationship string `json:"relationship"`
		Ocupation    string `json:"ocupation"`
		Contact      string  `json:"contact"`
	}

	err := app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	student := &data.Student{
		FirstName: input.FirstName,
		LastName:  input.LastName,
		Gender:    input.Gender,
	}

	headers := make(http.Header)
	headers.Set("Location", fmt.Sprintf("/guardians/%d", student.StudentID))

	err = app.writeJSON(w, http.StatusCreated, envelope{"student": student}, headers)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) showGuardianHandler(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil || id < 1 {
		app.notFoundResponse(w, r)
		return
	}

	bday := time.Date(1995, time.January, 1, 0, 0, 0, 0, time.UTC)
	student := &data.Student{
		StudentID:   id,
		FirstName:   "Liam",
		LastName:    "Luna",
		Gender:      "male",
		DateOfBirth: data.Date(bday),
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"student": student}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}

}

func (app *application) updateGuardianHandler(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil || id < 1 {
		app.notFoundResponse(w, r)
		return
	}

	var input struct {
		FirstName   string    `json:"first_name"`
		LastName    string    `json:"last_name"`
		Gender      string    `json:"gender"`
		DateOfBirth data.Date `json:"date_of_birth"`
	}

	err = app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	student := &data.Student{
		StudentID:   id,
		FirstName:   input.FirstName,
		LastName:    input.LastName,
		Gender:      input.Gender,
		DateOfBirth: input.DateOfBirth,
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"student": student}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) deleteGuardianHandler(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil || id < 1 {
		app.notFoundResponse(w, r)
		return
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"message": "student deleted successfully"}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}
