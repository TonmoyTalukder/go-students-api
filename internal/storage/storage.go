package storage

import "github.com/TonmoyTalukder/go-students-api/internal/types"

type Storage interface {
	CreateStudent(name string, email string, age int) (int64, error)
	GetStudentById(id int64) (types.Student, error)
	GetStudents() ([]types.Student, error)
	GetStudentsWithPagination(limit int, offset int) ([]types.Student, error)
	CountStudents() (int, error)
	UpdateStudentById(id int64, name string, email string, age int) error
	DeleteStudentById(id int64) error
}
