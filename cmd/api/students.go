package main

import (
	"net/http"
	"strconv"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/liamgluna/daycare-server/internal/data"
)

// func (app *application) createStudentHandler(w http.ResponseWriter, r *http.Request) {
// 	var input struct {
// 		FirstName   string    `json:"first_name"`
// 		LastName    string    `json:"last_name"`
// 		Gender      string    `json:"gender"`
// 		DateOfBirth data.Date `json:"date_of_birth"`
// 	}

// 	err := app.readJSON(w, r, &input)
// 	if err != nil {
// 		app.badRequestResponse(w, r, err)
// 		return
// 	}

// 	student := &data.Student{
// 		FirstName:   input.FirstName,
// 		LastName:    input.LastName,
// 		Gender:      input.Gender,
// 		DateOfBirth: input.DateOfBirth,
// 	}

// 	err = app.models.Students.Insert(student)
// 	if err != nil {
// 		app.serverErrorResponse(w, r, err)
// 		return
// 	}
// 	app.logger.Info("sheesh")
// 	err = app.writeJSON(w, http.StatusCreated, envelope{"student": student}, nil)
// 	if err != nil {
// 		app.serverErrorResponse(w, r, err)
// 	}
// }

/*
The post request:

	{
	  "student": {
	    "first_name": "John",
	    "last_name": "Wick",
	    "gender": "Male",
	    "date_of_birth": "2020-Mar-01"
	  },
	  "guardians": [
	    {
	      "first_name": "John",
	      "last_name": "Doe",
	      "gender": "Male",
	      "relationship": "Father",
	      "ocupation": "IT Specialist",
	      "contact": 1234567890
	    },
	    {
	      "first_name": "Jane",
	      "last_name": "Doe",
	      "gender": "Male",
	      "relationship": "Mother",
	      "ocupation": "attorney",
	      "contact": 1234567890
	    }
	  ]
	}
*/

func (app *application) createStudentHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Student struct {
			FirstName   string    `json:"first_name"`
			LastName    string    `json:"last_name"`
			Gender      string    `json:"gender"`
			DateOfBirth data.Date `json:"date_of_birth"`
		} `json:"student"`
		Guardians []struct {
			FirstName    string `json:"first_name"`
			LastName     string `json:"last_name"`
			Gender       string `json:"gender"`
			Relationship string `json:"relationship"`
			Occupation   string `json:"ocupation"`
			Contact      int64  `json:"contact"`
		} `json:"guardians"`
	}
	err := app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	student := &data.Student{
		FirstName:   input.Student.FirstName,
		LastName:    input.Student.LastName,
		Gender:      input.Student.Gender,
		DateOfBirth: input.Student.DateOfBirth,
	}

	guardians := make([]*data.Guardian, len(input.Guardians))
	for i, g := range input.Guardians {
		guardians[i] = &data.Guardian{
			FirstName:    g.FirstName,
			LastName:     g.LastName,
			Gender:       g.Gender,
			Relationship: g.Relationship,
			Occupation:   g.Occupation,
			Contact:      g.Contact,
		}
	}

	err = app.models.Students.InsertWithGuardians(student, guardians)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"message": "sheesh"}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) showStudentHandler(w http.ResponseWriter, r *http.Request) {
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

func (app *application) updateStudentHandler(w http.ResponseWriter, r *http.Request) {
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

func (app *application) deleteStudentHandler(w http.ResponseWriter, r *http.Request) {
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
