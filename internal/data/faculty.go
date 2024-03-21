package data

import (
	"context"
	"database/sql"
	"errors"
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
	Contact   string  `json:"contact"`
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

func (m FacultyModel) Update(faculty *Faculty) error {
	query := `
		UPDATE faculty 
		SET first_name = $1, last_name = $2, email = $3, contact = $4, position = $5
		WHERE faculty_id = $6
		`
	args := []any{faculty.FirstName, faculty.LastName, faculty.Email, faculty.Contact, faculty.Position, faculty.FacultyID}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	result, err := m.DB.ExecContext(ctx, query, args...)
	if err != nil {
		return err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rows == 0 {
		return ErrRecordNotFound
	}

	return nil
}

func (m FacultyModel) Get(id int64) (*Faculty, error) {
	if id < 1 {
		return nil, ErrRecordNotFound
	}

	query := `
		SELECT faculty_id, first_name, last_name, email, contact, position
		FROM faculty
		WHERE faculty_id = $1
		`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var faculty Faculty

	row := m.DB.QueryRowContext(ctx, query, id)
	err := row.Scan(
		&faculty.FacultyID,
		&faculty.FirstName,
		&faculty.LastName,
		&faculty.Email,
		&faculty.Contact,
		&faculty.Position,
	)

	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrRecordNotFound
		default:
			return nil, err
		}
	}

	return &faculty, nil
}

func (m FacultyModel) Delete(id int64) error {
	if id < 1 {
		return ErrRecordNotFound
	}

	query := `DELETE FROM faculty WHERE faculty_id = $1`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	result, err := m.DB.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return ErrRecordNotFound
	}

	return nil
}