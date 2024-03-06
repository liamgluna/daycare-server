package data

import (
	"context"
	"database/sql"
	"time"
)

type FacultyModel struct {
	DB *sql.DB
}

type Faculty struct {
	FacultyID int64  `json:"faculty_id"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Email     string `json:"email"`
	Contact   int64  `json:"contact"`
	Position  string `json:"position"`
}

func (m FacultyModel) Insert(faculty *Faculty) error {
	query := `
		INSERT INTO faculty (first_name, last_name, email, contact, position) 
		VALUES ($1, $2, $3, $4, $5)
		RETURNING faculty_id
		`
	args := []any{faculty.FirstName, faculty.LastName, faculty.Email, faculty.Contact, faculty.Position}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	return m.DB.QueryRowContext(ctx, query, args...).Scan(&faculty.FacultyID)
}
