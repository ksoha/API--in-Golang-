package sqlite

import (
	"database/sql"

	"github.com/ksoha/API-in-Golang/internal/config"
	_ "github.com/mattn/go-sqlite3" //importing the sqlite3 driver
)

// implementing the creatStudent method of storage interface
type Sqlite struct {
	Db *sql.DB //field name = Db, type = pointer to sql.DB
}

// intializing the database usig config
// just take the reference of the config struct and return the sqlite and error

func New(cfg *config.Config) (*Sqlite, error) {
	//open the database using the storage path from confi
	db, err := sql.Open("sqlite3", cfg.StoragePath)
	if err != nil {
		return nil, err
	}

	//creating a table(only if it deos not exist)
	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS students(
      id INTEGER PRIMARY KEY AUTOINCREMENT, 
	  name TEXT NOT NULL, 
	  age INTEGER NOT NULL, 
	  email TEXT NOT NULL UNIQUE 
   )`)

	//if error occurs while creating the table, return nil and error
	if err != nil {
		return nil, err
	}

	//returning the sqlite struct with the db reference and nil as error

	return &Sqlite{
		Db: db,
	}, nil
}

// implmenting the storage interface in this fucntion so the student.go can implemwnt the storage interface
// to attach this function to the sqlite struct we use the reciever
func (s *Sqlite) CreateStudent(name string, email string, age int) (int64, error) {

	//creating record in the database
	stmt, err := s.Db.Prepare("INSERT INTO students(name, email, age) VALUES(?, ?, ?)")
	if err != nil {
		return 0, err
	}
	defer stmt.Close() //closing the statement after the function is executed

	//executing the statement with the values passed to the function
	result, err := stmt.Exec(name, email, age)
	if err != nil {
		return 0, err
	}

	lastid, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	return lastid, nil

}
