package shared

import "testing"

func TestNewHttpHeaders(t *testing.T) {
	origins := "http://localhost:3000,http://localhost:3001"
	methods := "GET,POST,PUT,DELETE"
	headers := NewHttpHeaders(origins, methods)

	if headers == nil {
		t.Errorf("Expected headers to not be nil")
	}
}

func TestCreateHeaders(t *testing.T) {
	origins := "http://localhost:3000,http://localhost:3001"
	methods := "GET,POST,PUT,DELETE"
	headers := NewHttpHeaders(origins, methods)

	if headers == nil {
		t.Errorf("Expected headers to not be nil")
	}

	expected := map[string]string{
		"Access-Control-Allow-Origin":  origins,
		"Access-Control-Allow-Methods": methods,
		"Access-Control-Allow-Headers": "*",
		"Content-Type":                 "application/json",
	}

	actual := headers.CreateHeaders()

	if len(actual) != len(expected) {
		t.Errorf("Expected headers to have %d items, got %d", len(expected), len(actual))
	}

	for k, v := range expected {
		if actual[k] != v {
			t.Errorf("Expected %s to be %s, got %s", k, v, actual[k])
		}
	}
}
