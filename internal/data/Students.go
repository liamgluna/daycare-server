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
	      "contact": 1234567890
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
	var studentID int
	err = tx.QueryRow(query, student.FirstName, student.LastName, student.Gender, time.Time(student.DateOfBirth)).Scan(&studentID)
	if err != nil {
		return err
	}

	// Insert guardians and associate them with the student
	for _, guardian := range guardians {
		// Insert guardian
		query = `
			INSERT INTO guardians (first_name, last_name, gender, relationship, occupation, contact)
			VALUES ($1, $2, $3, $4, $5, $6)
			RETURNING guardian_id
		`
		var guardianID int
		err := tx.QueryRow(query, guardian.FirstName, guardian.LastName, guardian.Gender, guardian.Relationship, guardian.Occupation, guardian.Contact).Scan(&guardianID)
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
