package db

import (
	"database/sql"
	"htmx_demo/internal/components"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

var DB *sql.DB

func InitDB() error {
	var err error
	DB, err = sql.Open("sqlite3", "demo.db")
	if err != nil {
		return err
	}
	if err := DB.Ping(); err != nil {
		return err
	}
	log.Println("Connected to the database")
	// Copy database into memory
	if _, err := DB.Exec(`ATTACH DATABASE ':memory:' AS memdb;`); err != nil {
		return err
	}
	if _, err := DB.Exec(`CREATE TABLE memdb.subjects AS SELECT * FROM subjects;`); err != nil {
		return err
	}
	return nil
}

func QueryOptions() (components.AllOptions, error) {
	var allOptions components.AllOptions
	species_rows, err := DB.Query(`SELECT DISTINCT species FROM memdb.subjects;`)
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
	type_rows, err := DB.Query(`SELECT DISTINCT type FROM memdb.subjects;`)
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

func QueryFiltered(filter components.Filter) ([]components.Profile, error) {
	var profiles []components.Profile
	rows, err := DB.Query(`SELECT name, picture_url, age FROM memdb.subjects WHERE species = ? AND type = ? AND UPPER(name) LIKE '%' || UPPER(?) || '%'`, filter.Species, filter.Type, filter.Name)
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
