package handlers

import (
	"htmx_demo/internal/components"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

var handler *Handlers

type mockDB struct {
	queryVal any
}

func (m *mockDB) QueryOptions() (components.AllOptions, error) {
	return components.AllOptions{}, nil
}
func (m *mockDB) QueryFiltered(filter components.Filter) ([]components.Profile, error) {
	m.queryVal = &filter
	return nil, nil
}

func init() {
	handler = &Handlers{Router: gin.Default(), Db: &mockDB{}}
	handler.SetupRoutes()
}

// Test that the index page renders
func TestIndex(t *testing.T) {
	req, err := http.NewRequest("GET", "/", nil)
	if err != nil {
		t.Fatal(err)
	}
	resp := httptest.NewRecorder()
	handler.Router.ServeHTTP(resp, req)
	if resp.Code != http.StatusOK {
		t.Fatalf("Expected status 200, got %d", resp.Code)
	}
}

// Test that the filter page marshals the query parameters correctly
func TestFilter(t *testing.T) {
	req, err := http.NewRequest("GET", "/filter?species=dog&type=2345&name=", nil)
	if err != nil {
		t.Fatal(err)
	}
	resp := httptest.NewRecorder()
	handler.Router.ServeHTTP(resp, req)
	if resp.Code != http.StatusOK {
		t.Fatalf("Expected status 200, got %d", resp.Code)
	}
	if filter, ok := handler.Db.(*mockDB).queryVal.(*components.Filter); !ok {
		t.Fatalf("Expected filter to be of type components.Filter, got %T", handler.Db.(*mockDB).queryVal)
	} else {
		if filter.Species != "dog" {
			t.Fatalf("Expected species to be dog, got %s", filter.Species)
		}
		if filter.Type != "2345" {
			t.Fatalf("Expected type to be plant, got %s", filter.Type)
		}
		if filter.Name != "" {
			t.Fatalf("Expected name to be empty string, got %s", filter.Name)
		}
	}
	req, err = http.NewRequest("GET", `/filter?type=🐻&name=Lorem ipsum dolor sit amet, consetetur sadipscing elitr, sed diam nonumy eirmod tempor invidunt ut labore et dolore magna aliquyam erat, sed diam voluptua. At vero eos et accusam et justo duo dolores et ea rebum. Stet clita kasd gubergren, no sea takimata sanctus est Lorem ipsum dolor sit amet. Lorem ipsum dolor sit amet, consetetur sadipscing elitr, sed diam nonumy eirmod tempor invidunt ut labore et dolore magna aliquyam erat, sed diam voluptua. At vero eos et accusam et justo duo dolores et ea rebum. Stet clita kasd gubergren, no sea takimata sanctus est Lorem ipsum dolor sit amet.`, nil)
	if err != nil {
		t.Fatal(err)
	}
	resp = httptest.NewRecorder()
	handler.Router.ServeHTTP(resp, req)
	if resp.Code != http.StatusOK {
		t.Fatalf("Expected status 200, got %d", resp.Code)
	}
	if filter, ok := handler.Db.(*mockDB).queryVal.(*components.Filter); !ok {
		t.Fatalf("Expected filter to be of type components.Filter, got %T", handler.Db.(*mockDB).queryVal)
	} else {
		if filter.Species != "" {
			t.Fatalf("Expected species to be empty string, got %s", filter.Species)
		}
		if filter.Type != "🐻" {
			t.Fatalf("Expected type to be empty string, got %s", filter.Type)
		}
		if filter.Name != `Lorem ipsum dolor sit amet, consetetur sadipscing elitr, sed diam nonumy eirmod tempor invidunt ut labore et dolore magna aliquyam erat, sed diam voluptua. At vero eos et accusam et justo duo dolores et ea rebum. Stet clita kasd gubergren, no sea takimata sanctus est Lorem ipsum dolor sit amet. Lorem ipsum dolor sit amet, consetetur sadipscing elitr, sed diam nonumy eirmod tempor invidunt ut labore et dolore magna aliquyam erat, sed diam voluptua. At vero eos et accusam et justo duo dolores et ea rebum. Stet clita kasd gubergren, no sea takimata sanctus est Lorem ipsum dolor sit amet.` {
			t.Fatalf("Expected name to be empty string, got %s", filter.Name)
		}
	}
}
