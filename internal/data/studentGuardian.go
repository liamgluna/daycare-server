package data

import "database/sql"

type StudentGuardianModel struct {
	DB *sql.DB
}

type StudentGuardian struct {
	StudentID  int64 `json:"student_id"`
	GuardianID int64 `json:"guardian_id"`
}
