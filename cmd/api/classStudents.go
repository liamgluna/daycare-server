package main

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/liamgluna/daycare-server/internal/data"
)

func (app *application) createClassStudentHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		ClassID   int64 `json:"class_id"`
		StudentID int64 `json:"student_id"`
	}

	err := app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	classStudent := &data.ClassStudents{
		ClassID:   input.ClassID,
		StudentID: input.StudentID,
	}

	err = app.models.ClassStudents.Insert(classStudent)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	headers := make(http.Header)
	headers.Set("Location", fmt.Sprintf("/class-students/%d", classStudent.ClassID))

	err = app.writeJSON(w, http.StatusCreated, envelope{"classStudent": classStudent}, headers)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) deleteClassStudentHandler(w http.ResponseWriter, r *http.Request) {
	classID, err := strconv.ParseInt(chi.URLParam(r, "classID"), 10, 64)
	if err != nil || classID < 1 {
		app.notFoundResponse(w, r)
		return
	}

	studentID, err := strconv.ParseInt(chi.URLParam(r, "studentID"), 10, 64)
	if err != nil || studentID < 1 {
		app.notFoundResponse(w, r)
		return
	}

	err = app.models.ClassStudents.Delete(classID, studentID)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"message": "class student deleted successfully"}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}
