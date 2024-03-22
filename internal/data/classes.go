package data

import (
	"context"
	"database/sql"
	"time"
)

type ClassModel struct {
	DB *sql.DB
}

type Class struct {
	ClassID   int64  `json:"class_id"`
	FacultyID int64  `json:"faculty_id"`
	ClassName string `json:"class_name"`
	Term      string `json:"term"`
}

func (m ClassModel) Insert(class *Class) error {
	query := `
		INSERT INTO classes (faculty_id, class_name, term) 
		VALUES ($1, $2, $3)
		RETURNING class_id
		`
	args := []any{class.FacultyID, class.ClassName, class.Term}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	return m.DB.QueryRowContext(ctx, query, args...).Scan(&class.ClassID)
}

func (m ClassModel) Get(id int64) (*Class, error) {
	if id < 1 {
		return nil, ErrRecordNotFound
	}

	query := `SELECT * FROM classes WHERE class_id = $1`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	row := m.DB.QueryRowContext(ctx, query, id)

	class := &Class{}
	err := row.Scan(&class.ClassID, &class.FacultyID, &class.ClassName, &class.Term)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrRecordNotFound
		}
		return nil, err
	}

	return class, nil
}

func (m ClassModel) Update(class *Class) error {
	query := `
		UPDATE classes 
		SET class_name = $1, term = $2
		WHERE class_id = $3
		`
	args := []any{class.ClassName, class.Term, class.ClassID}

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

func (m ClassModel) Delete(id int64) error {
	if id < 1 {
		return ErrRecordNotFound
	}

	query := `DELETE FROM classes WHERE class_id = $1`

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
