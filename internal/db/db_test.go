package db

import (
	"htmx_demo/internal/components"
	"io"
	"log"
	"os"
	"path"
	"runtime"
	"testing"
)

func BenchmarkDBCalls(t *testing.B) {
	log.SetOutput(io.Discard)
	db, err := NewSqliteDB()
	if err != nil {
		t.Fatal(err)
	}
	for i := 0; i < t.N; i++ {
		db.QueryOptions()
	}
}

func init() {
	_, filename, _, _ := runtime.Caller(0)
	dir := path.Join(path.Dir(filename), "../../")
	err := os.Chdir(dir)
	if err != nil {
		log.Fatal(err)
	}
}

// There should be 3 species and 2 types in the database
func TestQueryOptions(t *testing.T) {
	db, err := NewSqliteDB()
	if err != nil {
		t.Fatal(err)
	}
	options, err := db.QueryOptions()
	if err != nil {
		t.Fatal(err)
	}
	if len(options.Species) != 3 {
		t.Log(options.Species)
		t.Fatalf("Expected 3 species, got %d", len(options.Species))
	}
	if len(options.Type) != 2 {
		t.Log(options.Type)
		t.Fatalf("Expected 2 types, got %d", len(options.Type))
	}
}

// There should be no dogs of type plant in the database
func TestQueryFiltered(t *testing.T) {
	db, err := NewSqliteDB()
	if err != nil {
		t.Fatal(err)
	}
	filter := components.Filter{
		Species: "Dog",
		Type:    "Plant",
		Name:    "",
	}
	results, err := db.QueryFiltered(filter)
	if err != nil {
		t.Fatal(err)
	}
	if len(results) != 0 {
		log.Println(results)
		t.Fatalf("Expected no results, got %d", len(results))
	}
}

// Testing invalid searches
func TestQueryFilteredInvalid(t *testing.T) {
	db, err := NewSqliteDB()
	if err != nil {
		t.Fatal(err)
	}
	filter := components.Filter{
		Species: "",
		Type:    "🐻",
		Name:    "Lorem ipsum dolor sit amet, consetetur sadipscing elitr, sed diam nonumy eirmod tempor invidunt ut labore et dolore magna aliquyam erat, sed diam voluptua. At vero eos et accusam et justo duo dolores et ea rebum. Stet clita kasd gubergren, no sea takimata sanctus est Lorem ipsum dolor sit amet. Lorem ipsum dolor sit amet, consetetur sadipscing elitr, sed diam nonumy eirmod tempor invidunt ut labore et dolore magna aliquyam erat, sed diam voluptua. At vero eos et accusam et justo duo dolores et ea rebum. Stet clita kasd gubergren, no sea takimata sanctus est Lorem ipsum dolor sit amet.",
	}
	results, err := db.QueryFiltered(filter)
	if err != nil {
		t.Fatal(err)
	}
	if len(results) != 0 {
		log.Println(results)
		t.Fatalf("Expected no results, got %d", len(results))
	}
}
