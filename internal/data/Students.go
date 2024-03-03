package data

import (
	"context"
	"database/sql"
	"time"
)

type StudentModel struct {
	DB *sql.DB
}

type Student struct {
	StudentID   int64  `json:"student_id"`
	FirstName   string `json:"first_name"`
	LastName    string `json:"last_name"`
	Gender      string `json:"gender"`
	DateOfBirth Date   `json:"date_of_birth"`
}

func (m StudentModel) Insert(student *Student) error {
	query := `
		INSERT INTO students (first_name, last_name, GENDER, date_of_birth)
		VALUES ($1, $2, $3, $4)
		RETURNING student_id
		`
	args := []any{student.FirstName, student.LastName, student.Gender, time.Time(student.DateOfBirth)}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	return m.DB.QueryRowContext(ctx, query, args...).Scan(&student.StudentID)
}
