package main

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/liamgluna/daycare-server/internal/data"
)

func (app *application) createClassStudentHandler(w http.ResponseWriter, r *http.Request) {
	classID, err := strconv.ParseInt(chi.URLParam(r, "classID"), 10, 64)
	if err != nil || classID < 1 {
		app.notFoundResponse(w, r)
		return
	}

	var input struct {
		StudentID int64 `json:"student_id"`
	}

	err = app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	classStudent := &data.ClassStudents{
		ClassID:   classID,
		StudentID: input.StudentID,
	}

	err = app.models.ClassStudents.Insert(classStudent)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	headers := make(http.Header)
	headers.Set("Location", fmt.Sprintf("/classes/%d/students", classStudent.ClassID))

	err = app.writeEnvelopedJSON(w, http.StatusCreated, envelope{"classStudent": classStudent}, headers)
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

	err = app.writeEnvelopedJSON(w, http.StatusOK, envelope{"message": "class student deleted successfully"}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) listClassStudentsHandler(w http.ResponseWriter, r *http.Request) {
	classID, err := strconv.ParseInt(chi.URLParam(r, "classID"), 10, 64)
	if err != nil || classID < 1 {
		app.notFoundResponse(w, r)
		return
	}

	classStudents, err := app.models.ClassStudents.GetStudentsByClassIDWithGuardian(classID)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	err = app.writeJSON(w, http.StatusOK, classStudents, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) addStudentAttendance(w http.ResponseWriter, r *http.Request) {
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

	var input struct {
		Present bool `json:"present"`
	}

	err = app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	studentAttendance := &data.StudentAttendance{
		ClassID:   classID,
		StudentID: studentID,
		Present:   input.Present,
		ClassDate: data.Date(time.Now()),
	}

	err = app.models.StudentAttendance.Insert(studentAttendance)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	headers := make(http.Header)
	headers.Set("Location", fmt.Sprintf("/classes/%d/students/%d/attendance", studentAttendance.ClassID, studentAttendance.StudentID))

	err = app.writeEnvelopedJSON(w, http.StatusCreated, envelope{"studentAttendance": studentAttendance}, headers)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}
