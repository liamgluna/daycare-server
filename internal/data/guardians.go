package data

import "database/sql"

type GuardianModel struct {
	DB *sql.DB
}

type Guardian struct {
	GuardianID   int64  `json:"guardian_id"`
	FirstName    string `json:"first_name"`
	LastName     string `json:"last_name"`
	Gender       string `json:"gender"`
	Relationship string `json:"relationship"`
	Occupation    string `json:"ocupation"`
	Contact      int64  `json:"contact"`
}
