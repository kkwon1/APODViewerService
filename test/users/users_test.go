package apod

import (
	"bytes"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/kkwon1/APODViewerService/api/users"
)

/*
	Currently tests are not atomic. If it crashes after creating a user,
	it will persist in the DB and not be deleted.

	TODO: Maybe add a before/after function to do some setup/teardown? Check if user exists
	delete if it does. Delete after everything.
*/

const testUserName = "DB_TEST_USER"
const testEmail = "DB_TEST_USER@test.com"
const testPassword = "TEST_PASSWORD"

func TestCreateUserHappyPath(t *testing.T) {
	testJSON := fmt.Sprintf(`{"Username": "%s", "Email": "%s", "Password": "%s"}`, testUserName, testEmail, testPassword)
	var jsonStr = []byte(testJSON)

	req, err := http.NewRequest("POST", "/users", bytes.NewBuffer(jsonStr))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(users.CreateUser)
	handler.ServeHTTP(rr, req)
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}
	expected := fmt.Sprintf("Successfully created user: %s", testUserName)
	if rr.Body.String() != expected {
		t.Errorf("handler returned unexpected body: got %v want %v",
			rr.Body.String(), expected)
	}
}

func TestDeleteUserHappyPath(t *testing.T) {
	testJSON := fmt.Sprintf(`{"Username": "%s", "Password": "%s"}`, testUserName, testPassword)
	var jsonStr = []byte(testJSON)

	req, err := http.NewRequest("POST", "/users", bytes.NewBuffer(jsonStr))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(users.DeleteUser)
	handler.ServeHTTP(rr, req)
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}
	expected := fmt.Sprintf("Successfully deleted user: %s", testUserName)
	if rr.Body.String() != expected {
		t.Errorf("handler returned unexpected body: got %v want %v",
			rr.Body.String(), expected)
	}
}

func TestCreateUserNoUsername(t *testing.T) {
	testJSON := fmt.Sprintf(`"Email": "%s", "Password": "%s"}`, testEmail, testPassword)
	var jsonStr = []byte(testJSON)

	req, err := http.NewRequest("POST", "/users", bytes.NewBuffer(jsonStr))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(users.CreateUser)
	handler.ServeHTTP(rr, req)
	if status := rr.Code; status != http.StatusBadRequest {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}
	expected := "Unexpected request payload"
	if rr.Body.String() != expected {
		t.Errorf("handler returned unexpected body: got %v want %v",
			rr.Body.String(), expected)
	}
}

func TestCreateUserNoEmail(t *testing.T) {
	testJSON := fmt.Sprintf(`"Username": "%s", "Password": "%s"}`, testUserName, testPassword)
	var jsonStr = []byte(testJSON)

	req, err := http.NewRequest("POST", "/users", bytes.NewBuffer(jsonStr))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(users.CreateUser)
	handler.ServeHTTP(rr, req)
	if status := rr.Code; status != http.StatusBadRequest {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}
	expected := "Unexpected request payload"
	if rr.Body.String() != expected {
		t.Errorf("handler returned unexpected body: got %v want %v",
			rr.Body.String(), expected)
	}
}

func TestCreateUserNoPassword(t *testing.T) {
	testJSON := fmt.Sprintf(`"Username": "%s", "Email": "%s"}`, testUserName, testEmail)
	var jsonStr = []byte(testJSON)

	req, err := http.NewRequest("POST", "/users", bytes.NewBuffer(jsonStr))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(users.CreateUser)
	handler.ServeHTTP(rr, req)
	if status := rr.Code; status != http.StatusBadRequest {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}
	expected := "Unexpected request payload"
	if rr.Body.String() != expected {
		t.Errorf("handler returned unexpected body: got %v want %v",
			rr.Body.String(), expected)
	}
}
