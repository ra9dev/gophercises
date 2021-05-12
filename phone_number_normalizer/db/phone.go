package db

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq" //...
)

// Phone represents the phone_numbers table in the DB
type Phone struct {
	ID     int
	Number string
}

// Store ...
type Store struct {
	db *sql.DB
}

// Open ...
func Open(driverName, dataSource string) (*Store, error) {
	db, err := sql.Open(driverName, dataSource)
	if err != nil {
		return nil, err
	}
	return &Store{db}, nil
}

// Close ...
func (s *Store) Close() error {
	return s.db.Close()
}

// Seed ...
func (s *Store) Seed() error {
	data := []string{
		"1234567890",
		"123 456 7891",
		"(123) 456 7892",
		"(123) 456-7893",
		"123-456-7894",
		"(123)456-7892",
	}

	for _, number := range data {
		if _, err := insertPhone(s.db, number); err != nil {
			return err
		}
	}
	return nil
}

func insertPhone(db *sql.DB, phone string) (int, error) {
	stmt, err := db.Prepare(`INSERT INTO phone_numbers(value) VALUES($1) RETURNING id`)
	if err != nil {
		return -1, err
	}
	defer stmt.Close()

	var id int
	err = stmt.QueryRow(phone).Scan(&id)
	if err != nil {
		return -1, err
	}

	return id, nil
}

// AllPhones ...
func (s *Store) AllPhones() ([]Phone, error) {
	rows, err := s.db.Query("SELECT * from phone_numbers")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var phones []Phone
	for rows.Next() {
		var p Phone
		if err := rows.Scan(&p.ID, &p.Number); err != nil {
			return nil, err
		}
		phones = append(phones, p)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return phones, nil
}

// FindPhone ...
func (s *Store) FindPhone(number string) (*Phone, error) {
	var p Phone
	stmt, err := s.db.Prepare("SELECT * FROM phone_numbers WHERE value=$1")
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	row := stmt.QueryRow(number)
	err = row.Scan(&p.ID, &p.Number)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return &p, nil
}

// UpdatePhone ...
func (s *Store) UpdatePhone(p *Phone) error {
	stmt, err := s.db.Prepare(`UPDATE phone_numbers SET value=$2 WHERE id=$1`)
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(p.ID, p.Number)
	return err
}

// DeletePhone ...
func (s *Store) DeletePhone(id int) error {
	stmt, err := s.db.Prepare(`DELETE FROM phone_numbers WHERE id=$1`)
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(id)
	return err
}

// Migrate ...
func Migrate(driverName, dataSource string) error {
	db, err := sql.Open(driverName, dataSource)
	if err != nil {
		return err
	}
	err = createPhoneNumbersTable(db)
	if err != nil {
		return err
	}
	return db.Close()
}

func createPhoneNumbersTable(db *sql.DB) error {
	statement := `
    CREATE TABLE IF NOT EXISTS phone_numbers (
      id SERIAL,
      value VARCHAR(255)
	)`

	_, err := db.Exec(statement)
	return err
}

// Reset ...
func Reset(driverName, dataSource, dbName string) error {
	source := fmt.Sprintf("%s dbname=%s", dataSource, "template1")
	db, err := sql.Open(driverName, source)
	if err != nil {
		return err
	}
	err = resetDB(db, dbName)
	if err != nil {
		return err
	}
	return db.Close()
}

func resetDB(db *sql.DB, name string) error {
	_, err := db.Exec("DROP DATABASE IF EXISTS " + name)
	if err != nil {
		return err
	}
	return createDB(db, name)
}

func createDB(db *sql.DB, name string) error {
	_, err := db.Exec("CREATE DATABASE " + name)
	if err != nil {
		return err
	}
	return nil
}
