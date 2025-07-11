package sqlite

import (
	"database/sql"
	"fmt"

	"github.com/TonmoyTalukder/go-students-api/internal/config"
	"github.com/TonmoyTalukder/go-students-api/internal/types"

	// _ "github.com/mattn/go-sqlite3"
	_ "modernc.org/sqlite"
)

type Sqlite struct {
	Db *sql.DB
}

func New(cfg *config.Config) (*Sqlite, error) {
	db, err := sql.Open("sqlite", cfg.StoragePath)
	if err != nil {
		return nil, err
	}

	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS students (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			name TEXT,
			email TEXT,
			age INTEGER
	)`)

	if err != nil {
		return nil, err
	}

	return &Sqlite{
		Db: db,
	}, nil
}

func (s *Sqlite) CreateStudent(name string, email string, age int) (int64, error) {

	stmt, err := s.Db.Prepare("INSERT INTO students (name, email, age) VALUES (?, ?, ?)")
	if err != nil {
		return 0, err
	}

	defer stmt.Close()

	result, errExe := stmt.Exec(name, email, age)
	if errExe != nil {
		return 0, errExe
	}

	lastId, errLII := result.LastInsertId()
	if errLII != nil {
		return 0, errLII
	}

	return lastId, nil
}

func (s *Sqlite) GetStudentById(id int64) (types.Student, error) {
	stmt, err := s.Db.Prepare("SELECT id, name, email, age FROM students WHERE id = ? LIMIT 1")
	if err != nil {
		return types.Student{}, err
	}

	defer stmt.Close()

	var student types.Student

	err = stmt.QueryRow(id).Scan(&student.Id, &student.Name, &student.Email, &student.Age)
	if err != nil {
		if err == sql.ErrNoRows {
			return types.Student{}, fmt.Errorf("No student fount with id %s", fmt.Sprint(id))
		}
		return types.Student{}, fmt.Errorf("Query error: %w", err)
	}

	return student, nil
}

func (s *Sqlite) GetStudents() ([]types.Student, error) {
	stmt, err := s.Db.Prepare("SELECT id, name, email, age FROM students")
	if err != nil {
		return nil, err
	}

	defer stmt.Close()

	rows, err := stmt.Query()
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var students []types.Student

	for rows.Next() {
		var student types.Student

		err := rows.Scan(&student.Id, &student.Name, &student.Email, &student.Age)
		if err != nil {
			return nil, err
		}

		students = append(students, student)
	}

	return students, nil
}

func (s *Sqlite) GetStudentsWithPagination(limit, offset int) ([]types.Student, error) {
	stmt, err := s.Db.Prepare("SELECT id, name, email, age FROM students LIMIT ? OFFSET ?")
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	rows, err := stmt.Query(limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var students []types.Student
	for rows.Next() {
		var student types.Student
		if err := rows.Scan(&student.Id, &student.Name, &student.Email, &student.Age); err != nil {
			return nil, err
		}
		students = append(students, student)
	}

	return students, nil
}

func (s *Sqlite) CountStudents() (int, error) {
	stmt, err := s.Db.Prepare("SELECT COUNT(*) FROM students")
	if err != nil {
		return 0, err
	}
	defer stmt.Close()

	var count int
	if err := stmt.QueryRow().Scan(&count); err != nil {
		return 0, err
	}

	return count, nil
}

func (s *Sqlite) UpdateStudentById(id int64, name string, email string, age int) error {
	stmt, err := s.Db.Prepare("UPDATE students SET name = ?, email = ?, age = ? WHERE id = ?")
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(name, email, age, id)
	return err
}

func (s *Sqlite) DeleteStudentById(id int64) error {
	stmt, err := s.Db.Prepare("DELETE FROM students WHERE id = ?")
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(id)
	return err
}