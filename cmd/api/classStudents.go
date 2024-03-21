package main

import (
	"fmt"
	"net/http"

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
