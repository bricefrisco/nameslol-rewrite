package shared

import "testing"

var origins = "http://localhost:3000,http://localhost:3001"
var methods = "GET,POST,PUT,DELETE"
var expectedHeaders = map[string]string{
	"Access-Control-Allow-Origin":  origins,
	"Access-Control-Allow-Methods": methods,
	"Access-Control-Allow-Headers": "*",
	"Content-Type":                 "application/json",
}

func TestNewHttpResponses(t *testing.T) {
	responses := NewHttpResponses(origins, methods)
	if responses == nil {
		t.Errorf("Expected headers to not be nil")
	}
}

func TestSuccess_HasCorrectStatusCode(t *testing.T) {
	responses := NewHttpResponses(origins, methods)

	if responses == nil {
		t.Errorf("Expected headers to not be nil")
	}

	response := responses.Success(map[string]string{"message": "success"})

	if response.StatusCode != 200 {
		t.Errorf("Expected status code to be 200, got %v", response.StatusCode)
	}
}

func TestSuccess_HasCorrectJsonBody(t *testing.T) {
	responses := NewHttpResponses(origins, methods)

	if responses == nil {
		t.Errorf("Expected headers to not be nil")
	}

	response := responses.Success(map[string]string{"message": "success"})

	if response.Body != "{\"message\":\"success\"}" {
		t.Errorf("Expected body to be {\"message\":\"success\"}, got %v", response.Body)
	}
}

func TestSuccess_HasCorrectHeaders(t *testing.T) {
	responses := NewHttpResponses(origins, methods)

	if responses == nil {
		t.Errorf("Expected headers to not be nil")
	}

	response := responses.Success(map[string]string{"message": "success"})

	for key, value := range expectedHeaders {
		if response.Headers[key] != value {
			t.Errorf("Expected header %v to be %v, got %v", key, value, response.Headers[key])
		}
	}
}

func TestError_HasCorrectStatusCode(t *testing.T) {
	responses := NewHttpResponses(origins, methods)

	if responses == nil {
		t.Errorf("Expected headers to not be nil")
	}

	response := responses.Error(400, "error")

	if response.StatusCode != 400 {
		t.Errorf("Expected status code to be 400, got %v", response.StatusCode)
	}
}

func TestError_HasCorrectJsonBody(t *testing.T) {
	responses := NewHttpResponses(origins, methods)

	if responses == nil {
		t.Errorf("Expected headers to not be nil")
	}

	response := responses.Error(400, "error")

	if response.Body != "{\"message\":\"error\"}" {
		t.Errorf("Expected body to be {\"Message\":\"error\"}, got %v", response.Body)
	}
}

func TestError_HasCorrectHeaders(t *testing.T) {
	responses := NewHttpResponses(origins, methods)

	if responses == nil {
		t.Errorf("Expected headers to not be nil")
	}

	response := responses.Error(400, "error")

	for key, value := range expectedHeaders {
		if response.Headers[key] != value {
			t.Errorf("Expected header %v to be %v, got %v", key, value, response.Headers[key])
		}
	}
}
