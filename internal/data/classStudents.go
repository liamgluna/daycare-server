package data

import (
	"context"
	"database/sql"
	"time"
)

type ClassStudentsModel struct {
	DB *sql.DB
}

type ClassStudents struct {
	ClassID   int64 `json:"class_id"`
	StudentID int64 `json:"student_id"`
}

func (m ClassStudentsModel) Insert(classStudent *ClassStudents) error {
	query := `
		INSERT INTO class_students (class_id, student_id)
		VALUES ($1, $2)
		`
	args := []any{classStudent.ClassID, classStudent.StudentID}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	_, err := m.DB.ExecContext(ctx, query, args...)
	return err
}

func (m ClassStudentsModel) Delete(classID, studentID int64) error {
	if classID < 1 || studentID < 1 {
		return ErrRecordNotFound
	}

	query := `DELETE FROM class_students WHERE class_id = $1 AND student_id = $2`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	result, err := m.DB.ExecContext(ctx, query, classID, studentID)
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
