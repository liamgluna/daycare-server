package main

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/golang-jwt/jwt/v5"
	"github.com/liamgluna/daycare-server/internal/data"
	"golang.org/x/crypto/bcrypt"
)

func (app *application) createFacultyHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		FirstName string `json:"first_name"`
		LastName  string `json:"last_name"`
		Email     string `json:"email"`
		Password  string `json:"password"`
		Contact   string `json:"contact"`
		Position  string `json:"position"`
	}

	err := app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	password, err := bcrypt.GenerateFromPassword([]byte(input.Password), 12)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	faculty := &data.Faculty{
		FirstName: input.FirstName,
		LastName:  input.LastName,
		Email:     input.Email,
		Password:  password,
		Contact:   input.Contact,
		Position:  input.Position,
	}

	err = app.models.Faculty.Insert(faculty)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrDuplicateEmail):
			app.userAlreadyExistResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	headers := make(http.Header)
	headers.Set("Location", fmt.Sprintf("/faculty/%d", faculty.FacultyID))

	claims := &jwt.RegisteredClaims{
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 72)),
		Issuer:    strconv.Itoa(int(faculty.FacultyID)),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString([]byte(app.cfg.jwtSecret))
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "jwt",
		Value:    tokenString,
		Expires:  time.Now().Add(time.Hour * 72),
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteStrictMode,
	})

	if err := app.writeJSON(w, http.StatusOK, faculty, nil); err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) loginFacultyHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	err := app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	faculty, err := app.models.Faculty.GetByEmail(input.Email)
	if err != nil {
		if errors.Is(err, data.ErrRecordNotFound) {
			app.invalidCredentialsResponse(w, r)
		} else {
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	if err = bcrypt.CompareHashAndPassword(faculty.Password, []byte(input.Password)); err != nil {
		if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
			app.invalidCredentialsResponse(w, r)
		} else {
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	claims := &jwt.RegisteredClaims{
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 72)),
		Issuer:    strconv.Itoa(int(faculty.FacultyID)),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString([]byte(app.cfg.jwtSecret))
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "jwt",
		Value:    tokenString,
		Expires:  time.Now().Add(time.Hour * 72),
		HttpOnly: true,
		Secure:   true,
	})

	err = app.writeJSON(w, http.StatusOK, faculty, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) logoutFacultyHandler(w http.ResponseWriter, r *http.Request) {
	_, err := r.Cookie("jwt")
	if err != nil {
		switch {
		case errors.Is(err, http.ErrNoCookie):
			app.badRequestResponse(w, r, err)

		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "jwt",
		Value:    "",
		Expires:  time.Now().Add(-time.Hour),
		HttpOnly: true,
		Secure:   true,
	})

	err = app.writeEnvelopedJSON(w, http.StatusOK, envelope{"message": "You have been logged out"}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) getUserWithTokenHandler(w http.ResponseWriter, r *http.Request) {
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

	faculty, err := app.models.Faculty.Get(id)
	if err != nil {
		if errors.Is(err, data.ErrRecordNotFound) {
			app.notFoundResponse(w, r)
		} else {
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	err = app.writeEnvelopedJSON(w, http.StatusOK, envelope{"faculty": faculty}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) updateFacultyHandler(w http.ResponseWriter, r *http.Request) {
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
		Contact   *string `json:"contact"`
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

	err = app.writeJSON(w, http.StatusOK, faculty, nil)
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

	err = app.writeEnvelopedJSON(w, http.StatusOK, envelope{"faculty": faculty}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}
