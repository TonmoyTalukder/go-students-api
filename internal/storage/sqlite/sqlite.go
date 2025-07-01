package sqlite

import (
	"database/sql"

	"github.com/TonmoyTalukder/go-students-api/internal/config"
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

	stmt, err  := s.Db.Prepare("INSERT INTO students (name, email, age) VALUES (?, ?, ?)")
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