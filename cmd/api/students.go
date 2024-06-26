package main

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/liamgluna/daycare-server/internal/data"
)

func (app *application) createStudentHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		FirstName   string    `json:"first_name"`
		LastName    string    `json:"last_name"`
		Gender      string    `json:"gender"`
		DateOfBirth data.Date `json:"date_of_birth"`
	}

	err := app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	student := &data.Student{
		FirstName:   input.FirstName,
		LastName:    input.LastName,
		Gender:      input.Gender,
		DateOfBirth: input.DateOfBirth,
	}

	err = app.models.Students.Insert(student)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	headers := make(http.Header)
	headers.Set("Location", fmt.Sprintf("/students/%d", student.StudentID))

	err = app.writeEnvelopedJSON(w, http.StatusCreated, envelope{"student": student}, headers)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) createStudentWithGuardiansHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Student struct {
			FirstName   string    `json:"first_name"`
			LastName    string    `json:"last_name"`
			Gender      string    `json:"gender"`
			DateOfBirth data.Date `json:"date_of_birth"`
		} `json:"student"`
		Guardian struct {
			FirstName    string `json:"first_name"`
			LastName     string `json:"last_name"`
			Gender       string `json:"gender"`
			Relationship string `json:"relationship"`
			Occupation   string `json:"occupation"`
			Contact      string `json:"contact"`
		} `json:"guardian"`
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

	guardian := &data.Guardian{
		FirstName:    input.Guardian.FirstName,
		LastName:     input.Guardian.LastName,
		Gender:       input.Guardian.Gender,
		Relationship: input.Guardian.Relationship,
		Occupation:   input.Guardian.Occupation,
		Contact:      input.Guardian.Contact,
	}

	err = app.models.Students.InsertWithGuardian(student, guardian)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	headers := make(http.Header)
	headers.Set("Location", fmt.Sprintf("/students/%d", student.StudentID))

	err = app.writeJSON(w, http.StatusCreated, student, headers)
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

	err = app.models.Students.Delete(id)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	err = app.writeEnvelopedJSON(w, http.StatusOK, envelope{"message": "student deleted successfully"}, nil)
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

	student, err := app.models.Students.Get(id)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	var input struct {
		FirstName   *string    `json:"first_name"`
		LastName    *string    `json:"last_name"`
		Gender      *string    `json:"gender"`
		DateOfBirth *data.Date `json:"date_of_birth"`
	}

	err = app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	if input.FirstName != nil {
		student.FirstName = *input.FirstName
	}

	if input.LastName != nil {
		student.LastName = *input.LastName
	}

	if input.Gender != nil {
		student.Gender = *input.Gender
	}

	if input.DateOfBirth != nil {
		student.DateOfBirth = *input.DateOfBirth
	}

	err = app.models.Students.Update(student)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	err = app.writeEnvelopedJSON(w, http.StatusOK, envelope{"student": student}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}

}

func (app *application) updateStudentAndGuardianHandler(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil || id < 1 {
		app.notFoundResponse(w, r)
		return
	}

	student, err := app.models.Students.Get(id)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	guardian, err := app.models.Guardians.GetByStudentID(id)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	var input struct {
		Student struct {
			FirstName   *string    `json:"first_name"`
			LastName    *string    `json:"last_name"`
			Gender      *string    `json:"gender"`
			DateOfBirth *data.Date `json:"date_of_birth"`
		} `json:"student"`
		Guardian struct {
			FirstName    *string `json:"first_name"`
			LastName     *string `json:"last_name"`
			Gender       *string `json:"gender"`
			Relationship *string `json:"relationship"`
			Occupation   *string `json:"occupation"`
			Contact      *string `json:"contact"`
		} `json:"guardian"`
	}

	err = app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}
	if input.Student.FirstName != nil {
		student.FirstName = *input.Student.FirstName
	}
	if input.Student.LastName != nil {
		student.LastName = *input.Student.LastName
	}
	if input.Student.Gender != nil {
		student.Gender = *input.Student.Gender
	}
	if input.Student.DateOfBirth != nil {
		student.DateOfBirth = *input.Student.DateOfBirth
	}
	if input.Guardian.FirstName != nil {
		guardian.FirstName = *input.Guardian.FirstName
	}
	if input.Guardian.LastName != nil {
		guardian.LastName = *input.Guardian.LastName
	}
	if input.Guardian.Gender != nil {
		guardian.Gender = *input.Guardian.Gender
	}
	if input.Guardian.Relationship != nil {
		guardian.Relationship = *input.Guardian.Relationship
	}
	if input.Guardian.Occupation != nil {
		guardian.Occupation = *input.Guardian.Occupation
	}
	if input.Guardian.Contact != nil {
		guardian.Contact = *input.Guardian.Contact
	}

	err = app.models.Students.UpdateWithGuardian(student, guardian)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	err = app.writeEnvelopedJSON(w, http.StatusOK, envelope{"student": student, "guardian": guardian}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) showStudentGuardiansHandler(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil || id < 1 {
		app.notFoundResponse(w, r)
		return
	}

	guardians, err := app.models.Guardians.GetByStudentID(id)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	err = app.writeJSON(w, http.StatusOK, guardians, nil)
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

	student, err := app.models.Students.Get(id)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	err = app.writeEnvelopedJSON(w, http.StatusOK, envelope{"student": student}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) listStudentsHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Name string
		data.Filters
	}

	qs := r.URL.Query()

	input.Name = app.readString(qs, "name", "")
	input.Filters.Page = app.readInt(qs, "page", 1)
	input.Filters.PageSize = app.readInt(qs, "page_size", 20)
	input.Filters.Sort = app.readString(qs, "sort", "last_name")
	input.Filters.SortSafelist = []string{"student_id", "first_name", "last_name", "-student_id", "-first_name", "-last_name"}

	students, metadata, err := app.models.Students.GetAll(input.Name, input.Filters)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	err = app.writeEnvelopedJSON(w, http.StatusOK, envelope{"students": students, "metadata": metadata}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}
