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
	Students          StudentModel
	Guardians         GuardianModel
	StudentGuardian   StudentGuardianModel
	Faculty           FacultyModel
	Classes           ClassesModel
	ClassStudents     ClassStudentsModel
	StudentAttendance StudentAttendanceModel
}

func NewModels(db *sql.DB) Models {
	return Models{
		Students:          StudentModel{DB: db},
		Guardians:         GuardianModel{DB: db},
		StudentGuardian:   StudentGuardianModel{DB: db},
		Faculty:           FacultyModel{DB: db},
		Classes:           ClassesModel{DB: db},
		ClassStudents:     ClassStudentsModel{DB: db},
		StudentAttendance: StudentAttendanceModel{DB: db},
	}
}
