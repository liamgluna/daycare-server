package data

import (
	"database/sql"
	"errors"
)

var (
	ErrRecordNotFound = errors.New("record not found")
	ErrEditConflict   = errors.New("edit conflict")
)

type Models struct {
	Students  StudentModel
	Guardians GuardianModel
}

func NewModels(db *sql.DB) Models {
	return Models{
		Students: StudentModel{DB: db},
		Guardians: GuardianModel{DB: db},
	}
}
