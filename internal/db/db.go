package db

import (
	"database/sql"
	"htmx_demo/internal/components"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

type DB interface {
	// initDB initializes the database
	initDB() error
	// QueryOptions returns all the options available in the database
	QueryOptions() (components.AllOptions, error)
	// QueryFiltered returns a list of profiles that match the filter
	QueryFiltered(filter components.Filter) ([]components.Profile, error)
}

type SqliteDB struct {
	Db *sql.DB
}

// Create a new SqliteDB instance and initialize the database
func NewSqliteDB() (*SqliteDB, error) {
	db := &SqliteDB{}
	if err := db.initDB(); err != nil {
		return nil, err
	}
	return db, nil
}

func (s *SqliteDB) initDB() error {
	var err error
	s.Db, err = sql.Open("sqlite3", "demo.db")
	if err != nil {
		return err
	}
	if err := s.Db.Ping(); err != nil {
		return err
	}
	log.Println("Connected to the database")
	// Copy database into memory
	if _, err := s.Db.Exec(`ATTACH DATABASE ':memory:' AS memdb;`); err != nil {
		return err
	}
	if _, err := s.Db.Exec(`CREATE TABLE memdb.subjects AS SELECT * FROM subjects;`); err != nil {
		return err
	}
	return nil
}

func (s *SqliteDB) QueryOptions() (components.AllOptions, error) {
	var allOptions components.AllOptions
	species_rows, err := s.Db.Query(`SELECT DISTINCT species FROM memdb.subjects;`)
	if err != nil {
		log.Println(err)
		return allOptions, err
	}

	for species_rows.Next() {
		var species string
		if err := species_rows.Scan(&species); err != nil {
			return allOptions, err
		}
		allOptions.Species = append(allOptions.Species, species)
	}
	type_rows, err := s.Db.Query(`SELECT DISTINCT type FROM memdb.subjects;`)
	if err != nil {
		return allOptions, err
	}
	for type_rows.Next() {
		var kind string //type is reserved
		if err := type_rows.Scan(&kind); err != nil {
			return allOptions, err
		}
		allOptions.Type = append(allOptions.Type, kind)
	}
	return allOptions, nil
}

func (s *SqliteDB) QueryFiltered(filter components.Filter) ([]components.Profile, error) {
	var profiles []components.Profile
	rows, err := s.Db.Query(`SELECT name, picture_url, age FROM memdb.subjects WHERE species = ? AND type = ? AND UPPER(name) LIKE '%' || UPPER(?) || '%'`, filter.Species, filter.Type, filter.Name)
	if err != nil {
		return profiles, err
	}
	for rows.Next() {
		var profile components.Profile
		if err := rows.Scan(&profile.Name, &profile.PictureUrl, &profile.Age); err != nil {
			return profiles, err
		}
		profiles = append(profiles, profile)
	}
	return profiles, nil
}
