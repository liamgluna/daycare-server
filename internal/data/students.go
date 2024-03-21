package data

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
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

/*
The post request:

	{
	  "student": {
	    "first_name": "John",
	    "last_name": "Wick",
	    "gender": "Male",
	    "date_of_birth": "2020-Mar-01"
	  },
	  "guardians": [
	    {
	      "first_name": "John",
	      "last_name": "Doe",
	      "gender": "Male",
	      "relationship": "Father",
	      "ocupation": "IT Specialist",
	      "contact": 1234567890
	    },
	    {
	      "first_name": "Jane",
	      "last_name": "Doe",
	      "gender": "Male",
	      "relationship": "Mother",
	      "ocupation": "attorney",
	      "contact": "1234567890"
	    }
	  ]
	}
*/
func (m StudentModel) InsertWithGuardians(student *Student, guardians []*Guardian) error {
	tx, err := m.DB.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// Insert student
	query := `
		INSERT INTO students (first_name, last_name, gender, date_of_birth)
		VALUES ($1, $2, $3, $4)
		RETURNING student_id
	`
	var studentID int64
	err = tx.QueryRow(query, student.FirstName, student.LastName, student.Gender, time.Time(student.DateOfBirth)).Scan(&studentID)
	if err != nil {
		return err
	}

	student.StudentID = studentID

	// Insert guardians and associate them with the student
	for _, guardian := range guardians {
		// Insert guardian
		query = `
			INSERT INTO guardians (first_name, last_name, gender, relationship, occupation, contact, student_id)
			VALUES ($1, $2, $3, $4, $5, $6, $7)
			RETURNING guardian_id
		`
		var guardianID int
		err := tx.QueryRow(query, guardian.FirstName, guardian.LastName, guardian.Gender, guardian.Relationship, guardian.Occupation, guardian.Contact, studentID).Scan(&guardianID)
		if err != nil {
			return err
		}

		// Associate guardian with the student
		query = `
			INSERT INTO student_guardian (student_id, guardian_id)
			VALUES ($1, $2)
		`
		_, err = tx.Exec(query, studentID, guardianID)
		if err != nil {
			return err
		}
	}

	// Commit transaction
	err = tx.Commit()
	if err != nil {
		return err
	}

	return nil
}

func (m StudentModel) Get(id int64) (*Student, error) {
	if id < 1 {
		return nil, ErrRecordNotFound
	}

	query := `
		SELECT student_id, first_name, last_name, gender, date_of_birth
		FROM students
		WHERE student_id = $1`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var student Student

	row := m.DB.QueryRowContext(ctx, query, id)
	err := row.Scan(
		&student.StudentID,
		&student.FirstName,
		&student.LastName,
		&student.Gender,
		&student.DateOfBirth,
	)

	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrRecordNotFound
		default:
			return nil, err
		}
	}

	return &student, nil
}

func (m StudentModel) Delete(id int64) error {
	if id < 1 {
		return ErrRecordNotFound
	}

	query := `DELETE FROM students WHERE student_id = $1`

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

// delete also the associated guardians
// func (m StudentModel) Delete(id int64) error {
// 	if id < 1 {
// 		return ErrRecordNotFound
// 	}

// 	tx, err := m.DB.Begin()
// 	if err != nil {
// 		return err
// 	}
// 	defer tx.Rollback()

// 	query := `
// 		DELETE FROM students
// 		WHERE student_id = $1
// 	`

// 	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
// 	defer cancel()

// 	result, err := tx.ExecContext(ctx, query, id)
// 	if err != nil {
// 		return err
// 	}

// 	rowsAffected, err := result.RowsAffected()
// 	if err != nil {
// 		return err
// 	}

// 	if rowsAffected == 0 {
// 		return ErrRecordNotFound
// 	}

// 	// Delete associated guardians
// 	query = `
// 		DELETE FROM student_guardian
// 		WHERE student_id = $1
// 		`

// 	_, err = tx.ExecContext(ctx, query, id)
// 	if err != nil {
// 		return err
// 	}

// 	// Commit transaction
// 	err = tx.Commit()
// 	if err != nil {
// 		return err
// 	}

// 	return nil
// }

func (m StudentModel) Update(student *Student) error {
	query := `
		UPDATE students
		SET first_name = $1, last_name = $2, gender = $3, date_of_birth = $4
		WHERE student_id = $5
		`

	args := []any{student.FirstName, student.LastName, student.Gender, time.Time(student.DateOfBirth), student.StudentID}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	result, err := m.DB.ExecContext(ctx, query, args...)
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

func (m StudentModel) GetAll(name string, filters Filters) ([]*Student, Metadata, error) {
	query := fmt.Sprintf(`
		SELECT count(*) OVER(), student_id, first_name, last_name, gender, date_of_birth
		FROM students
		WHERE (to_tsvector('simple', first_name || ' ' || last_name) @@ plainto_tsquery('simple', $1))
		OR $1 = ''
		ORDER BY %s %s, last_name ASC
		LIMIT $2 OFFSET $3`, filters.sortColumn(), filters.sortDirection())

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	rows, err := m.DB.QueryContext(ctx, query, name, filters.limit(), filters.offset())
	if err != nil {
		return nil, Metadata{}, err
	}
	defer rows.Close()

	var students []*Student
	totalRecords := 0

	for rows.Next() {
		var student Student
		err := rows.Scan(
			&totalRecords,
			&student.StudentID,
			&student.FirstName,
			&student.LastName,
			&student.Gender,
			&student.DateOfBirth,
		)
		if err != nil {
			return nil, Metadata{}, err
		}

		students = append(students, &student)
	}

	if err := rows.Err(); err != nil {
		return nil, Metadata{}, err
	}

	metadata := calculateMetadata(totalRecords, filters.Page, filters.PageSize)

	return students, metadata, nil
}
