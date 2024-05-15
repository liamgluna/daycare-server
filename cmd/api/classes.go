package main

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/golang-jwt/jwt/v5"
	"github.com/liamgluna/daycare-server/internal/data"
)

func (app *application) createClassHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		FacultyID int64  `json:"faculty_id"`
		ClassName string `json:"class_name"`
		Term      string `json:"term"`
	}

	err := app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	class := &data.Class{
		FacultyID: input.FacultyID,
		ClassName: input.ClassName,
		Term:      input.Term,
	}

	err = app.models.Classes.Insert(class)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	headers := make(http.Header)
	headers.Set("Location", fmt.Sprintf("/classes/%d", class.ClassID))

	err = app.writeEnvelopedJSON(w, http.StatusCreated, envelope{"class": class}, headers)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) updateClassHandler(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil || id < 1 {
		app.notFoundResponse(w, r)
		return
	}

	class, err := app.models.Classes.Get(id)
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
		ClassName *string `json:"class_name"`
		Term      *string `json:"term"`
	}

	err = app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	if input.ClassName != nil {
		class.ClassName = *input.ClassName
	}

	if input.Term != nil {
		class.Term = *input.Term
	}

	err = app.models.Classes.Update(class)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	err = app.writeEnvelopedJSON(w, http.StatusOK, envelope{"class": class}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) deleteClassHandler(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil || id < 1 {
		app.notFoundResponse(w, r)
		return
	}

	err = app.models.Classes.Delete(id)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	err = app.writeEnvelopedJSON(w, http.StatusOK, envelope{"message": "class deleted successfully"}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) showClassHandler(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil || id < 1 {
		app.notFoundResponse(w, r)
		return
	}

	class, err := app.models.Classes.Get(id)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	err = app.writeJSON(w, http.StatusOK, class, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) listClassesHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Name string
		data.Filters
	}

	qs := r.URL.Query()

	input.Name = app.readString(qs, "name", "")
	input.Filters.Page = app.readInt(qs, "page", 1)
	input.Filters.PageSize = app.readInt(qs, "page_size", 20)
	input.Filters.Sort = app.readString(qs, "sort", "class_id")
	input.Filters.SortSafelist = []string{"class_id", "class_name", "term", "-class_id", "-class_name", "-term"}

	classes, metadata, err := app.models.Classes.GetAll(input.Name, input.Filters)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	err = app.writeEnvelopedJSON(w, http.StatusOK, envelope{"classes": classes, "metadata": metadata}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}
}

func (app *application) listClassesByFacultyIDHandler(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("jwt")
	if err != nil {
		switch {
		case errors.Is(err, http.ErrNoCookie):
			app.badRequestResponse(w, r, err)

		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	token, err := jwt.ParseWithClaims(cookie.Value, &jwt.RegisteredClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(app.cfg.jwtSecret), nil
	})
	if err != nil {
		app.invalidCredentialsResponse(w, r)
		return
	}

	claims := token.Claims.(*jwt.RegisteredClaims)

	id, err := strconv.ParseInt(claims.Issuer, 10, 64)
	if err != nil {
		app.invalidCredentialsResponse(w, r)
		return
	}

	classes, err := app.models.Classes.GetAllByFacultyID(id)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	err = app.writeJSON(w, http.StatusOK, classes, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}
}
