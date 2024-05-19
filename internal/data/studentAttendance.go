package data

import (
	"context"
	"database/sql"
	"time"
)

type StudentAttendanceModel struct {
	DB *sql.DB
}

type StudentAttendance struct {
	StudentID int64 `json:"student_id"`
	ClassID   int64 `json:"class_id"`
	ClassDate Date  `json:"class_date"`
	Present   bool  `json:"present"`
}

func (m StudentAttendanceModel) Insert(studentAttendance *StudentAttendance) error {
	query := `
		INSERT INTO student_attendance (student_id, class_id, class_date, present) 
		VALUES ($1, $2, $3, $4)
		RETURNING student_id
		`
	args := []any{studentAttendance.StudentID, studentAttendance.ClassID, time.Time(studentAttendance.ClassDate), studentAttendance.Present}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := m.DB.QueryRowContext(ctx, query, args...).Scan(&studentAttendance.StudentID)
	if err != nil {
		return err
	}

	return nil
}

func (m StudentAttendanceModel) GetAttendance(date time.Time, classID int64) ([]*StudentAttendance, error) {
	query := `
		SELECT student_id, class_id, class_date, present
		FROM student_attendance
		WHERE class_date = $1 AND class_id = $2
		`
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	rows, err := m.DB.QueryContext(ctx, query, date, classID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	studentAttendances := []*StudentAttendance{}
	for rows.Next() {
		studentAttendance := &StudentAttendance{}
		err := rows.Scan(&studentAttendance.StudentID, &studentAttendance.ClassID, &studentAttendance.ClassDate, &studentAttendance.Present)
		if err != nil {
			return nil, err
		}
		studentAttendances = append(studentAttendances, studentAttendance)
	}

	return studentAttendances, nil
}	