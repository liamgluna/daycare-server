package data

import (
	"context"
	"database/sql"
	"fmt"
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
	Schedule  string `json:"schedule"`
}

func (m ClassModel) Insert(class *Class) error {
	query := `
		INSERT INTO classes (faculty_id, class_name, term, schedule) 
		VALUES ($1, $2, $3, $4)
		RETURNING class_id
		`
	args := []any{class.FacultyID, class.ClassName, class.Term, class.Schedule}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	return m.DB.QueryRowContext(ctx, query, args...).Scan(&class.ClassID)
}

func (m ClassModel) Get(id int64) (*Class, error) {
	if id < 1 {
		return nil, ErrRecordNotFound
	}

	query := `SELECT class_id, faculty_id, class_name, term, schedule
	 	FROM classes WHERE class_id = $1`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	row := m.DB.QueryRowContext(ctx, query, id)

	class := &Class{}
	err := row.Scan(&class.ClassID, &class.FacultyID, &class.ClassName, &class.Term, &class.Schedule)
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
		SET class_name = $1, term = $2, schedule = $3
		WHERE class_id = $4
		`
	args := []any{class.ClassName, class.Term, class.Schedule, class.ClassID}

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

func (m ClassModel) GetAll(name string, filters Filters) ([]*Class, Metadata, error) {
	query := fmt.Sprintf(`
		SELECT count(*) OVER(), class_id, faculty_id, class_name, term
		FROM classes
		WHERE (to_tsvector('simple', class_name) @@ plainto_tsquery('simple', $1))
		OR $1 = ''
		ORDER BY %s %s, class_id ASC
		LIMIT $2 OFFSET $3`, filters.sortColumn(), filters.sortDirection())

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	rows, err := m.DB.QueryContext(ctx, query, name, filters.limit(), filters.offset())
	if err != nil {
		return nil, Metadata{}, err
	}
	defer rows.Close()

	var classes []*Class
	totalRecords := 0

	for rows.Next() {
		var class Class
		err := rows.Scan(
			&totalRecords,
			&class.ClassID,
			&class.FacultyID,
			&class.ClassName,
			&class.Term,
		)
		if err != nil {
			return nil, Metadata{}, err
		}

		classes = append(classes, &class)
	}

	if err := rows.Err(); err != nil {
		return nil, Metadata{}, err
	}

	metaData := calculateMetadata(totalRecords, filters.Page, filters.PageSize)

	return classes, metaData, nil
}

func (m ClassModel) GetAllByFacultyID(faculty_id int64) ([]*Class, error) {
	query := `SELECT class_id, class_name, term, schedule
			FROM classes
			WHERE faculty_id = $1
			ORDER BY class_id ASC`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	rows, err := m.DB.QueryContext(ctx, query, faculty_id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var classes []*Class
	for rows.Next() {
		var class Class
		err := rows.Scan(
			&class.ClassID,
			&class.ClassName,
			&class.Term,
			&class.Schedule,
		)
		if err != nil {
			return nil, err
		}
		classes = append(classes, &class)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return classes, nil
}
