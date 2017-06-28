package mysql

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
)

type MySQLRepository struct {
	db *sql.DB
}

func (r MySQLRepository) Write(key, value string) error {
	_, err := r.db.Exec("REPLACE INTO stuff (id, payload) VALUES(?, ?)", key, value)
	return err
}

func (r MySQLRepository) Read(key string) (string, error) {
	var payload string
	err := r.db.QueryRow("SELECT payload FROM stuff WHERE id = ?", key).Scan(&payload)

	if err == sql.ErrNoRows {
		return "", nil
	}

	return payload, err
}

func NewMySQLRepository(url string) (MySQLRepository, error) {
	db, err := sql.Open("mysql", url)
	if err != nil {
		return MySQLRepository{}, err
	}

	_, err = db.Exec("CREATE TABLE IF NOT EXISTS stuff (id VARCHAR(255) NOT NULL PRIMARY KEY, payload VARCHAR(255))")
	if err != nil {
		return MySQLRepository{}, err
	}

	return MySQLRepository{db}, nil
}
