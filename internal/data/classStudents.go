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
