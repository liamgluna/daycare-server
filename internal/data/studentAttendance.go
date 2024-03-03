package data

import "database/sql"

type StudentAttendanceModel struct {
	DB *sql.DB
}

type StudentAttendance struct {
	StudentID int64 `json:"student_id"`
	ClassID   int64 `json:"class_id"`
	ClassDate Date  `json:"class_date"`
	Present   bool  `json:"present"`
}
