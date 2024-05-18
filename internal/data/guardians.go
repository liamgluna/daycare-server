package data

import (
	"context"
	"database/sql"
	"time"
)

type GuardianModel struct {
	DB *sql.DB
}

type Guardian struct {
	GuardianID   int64  `json:"guardian_id"`
	FirstName    string `json:"first_name"`
	LastName     string `json:"last_name"`
	Gender       string `json:"gender"`
	Relationship string `json:"relationship"`
	Occupation   string `json:"occupation"`
	Contact      string `json:"contact"`
}

// CREATE TABLE IF NOT EXISTS student_guardian (
//     student_id integer REFERENCES students(student_id) ON DELETE CASCADE,
//     guardian_id integer REFERENCES guardians(guardian_id) ON DELETE CASCADE,
//     PRIMARY KEY (student_id, guardian_id)
// );

func (m *GuardianModel) GetByStudentID(studentID int64) (*Guardian, error) {
	if studentID < 1 {
		return nil, ErrRecordNotFound
	}

	query := `
		SELECT g.guardian_id, g.first_name, g.last_name, gender, relationship, occupation, contact
		FROM guardians g
		INNER JOIN student_guardian sg ON g.guardian_id = sg.guardian_id
		WHERE sg.student_id = $1
	`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	g := &Guardian{}
	err := m.DB.QueryRowContext(ctx, query, studentID).Scan(
		&g.GuardianID,
		&g.FirstName,
		&g.LastName,
		&g.Gender,
		&g.Relationship,
		&g.Occupation,
		&g.Contact,
	)
	if err != nil {
		return nil, err
	}

	return g, nil
}
