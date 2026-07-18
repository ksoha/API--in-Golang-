package types

// creating a student struct to hold the student data
type Student struct {
	Id    int64
	Name  string `validate:"required`
	Email string `validate:"required,email"`
	Age   int    `validate:"required"`
}
