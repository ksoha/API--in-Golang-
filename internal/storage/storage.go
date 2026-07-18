package storage

import "github.com/ksoha/API-in-Golang/internal/types"

// creating an interface
type Storage interface {
	//createStudent method will return an id and an error
	CreateStudent(name string, email string, age int) (int64, error)
	GetStudentById(id int64) (types.Student, error)
}
