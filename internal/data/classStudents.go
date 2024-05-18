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

// func (m ClassStudentsModel) GetByClassID(classID int64) ([]int64, error) {
// 	if classID < 1 {
// 		return nil, ErrRecordNotFound
// 	}

// 	query := `SELECT student_id FROM class_students WHERE class_id = $1`

// 	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
// 	defer cancel()

// 	rows, err := m.DB.QueryContext(ctx, query, classID)
// 	if err != nil {
// 		return nil, err
// 	}
// 	defer rows.Close()

// 	var studentIDs []int64
// 	for rows.Next() {
// 		var studentID int64
// 		err := rows.Scan(&studentID)
// 		if err != nil {
// 			return nil, err
// 		}
// 		studentIDs = append(studentIDs, studentID)
// 	}

// 	if err = rows.Err(); err != nil {
// 		return nil, err
// 	}

// 	return studentIDs, nil
// }

func (m ClassStudentsModel) GetStudentsByClassID(classID int64) ([]*Student, error) {
	if classID < 1 {
		return nil, ErrRecordNotFound
	}

	query := `
		SELECT s.student_id, s.first_name, s.last_name, s.gender, s.date_of_birth
		FROM students s
		INNER JOIN class_students cs ON s.student_id = cs.student_id
		WHERE cs.class_id = $1
		`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	rows, err := m.DB.QueryContext(ctx, query, classID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var students []*Student

	for rows.Next() {
		var student Student
		err := rows.Scan(
			&student.StudentID,
			&student.FirstName,
			&student.LastName,
			&student.Gender,
			&student.DateOfBirth,
		)
		if err != nil {
			return nil, err
		}
		students = append(students, &student)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return students, nil
}

type StudentWithGuardian struct {
	StudentID       int64  `json:"student_id"`
	FirstName       string `json:"first_name"`
	LastName        string `json:"last_name"`
	Gender          string `json:"gender"`
	DateOfBirth     Date   `json:"date_of_birth"`
	Guardian        string `json:"guardian_name"`
	GuardianContact string `json:"guardian_contact"`
}

func (m ClassStudentsModel) GetStudentsByClassIDWithGuardian(classID int64) ([]*StudentWithGuardian, error) {
	if classID < 1 {
		return nil, ErrRecordNotFound
	}

	query := `
		SELECT s.student_id, s.first_name, s.last_name, s.gender, s.date_of_birth, concat(g.first_name, ' ', g.last_name) as guardian_name, g.contact
		FROM students s
		INNER JOIN class_students cs ON s.student_id = cs.student_id
		INNER JOIN student_guardian sg ON s.student_id = sg.student_id
		INNER JOIN guardians g ON sg.guardian_id = g.guardian_id
		WHERE cs.class_id = $1
		`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	rows, err := m.DB.QueryContext(ctx, query, classID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var students []*StudentWithGuardian
	for rows.Next() {
		var student StudentWithGuardian
		err := rows.Scan(
			&student.StudentID,
			&student.FirstName,
			&student.LastName,
			&student.Gender,
			&student.DateOfBirth,
			&student.Guardian,
			&student.GuardianContact,
		)
		if err != nil {
			return nil, err
		}
		students = append(students, &student)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return students, nil
}
