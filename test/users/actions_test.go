package apod

import (
	"bytes"
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/kkwon1/APODViewerService/api/users"
)

// Mock the token verifier by implementing the interface
type mockTokenVerifier struct{}

func (m *mockTokenVerifier) VerifyToken(ctx context.Context, idToken string) (bool, error) {
	return true, nil
}

var userAction users.UserAction = users.NewUserAction(&mockTokenVerifier{})

/*
	Currently tests are not atomic. Depends on order of tests being run.
	If it crashes while running a test, records will persist in the DB and not be deleted.

	TODO: Maybe add a before/after function to do some setup/teardown? Check if record exists
	delete if it does. Delete after everything. Probably need to instantiate or bring in mongodb client
	to manually create/delete users
*/

const testUserID = "DB_TEST_ID"
const testApodURL = "https://apod.nasa.gov/apod/image/2005/c2020_f8_2020_05_02dp_1024.jpg"
const testApodName = "Long Tailed Comet SWAN"
const testApodDate = "2020-05-08"
const testMediaType = "image"
const testDescription = "Test description"
const testActionDate = "2020-05-11T18:19:58.747Z"

const testToken = "Bearer mock-token"

func getTestJson(action string) string {
	testJSON := fmt.Sprintf(`{
		"UserID": "%s",
		"Action": "%s",
		"ApodURL": "%s",
		"ApodName": "%s",
		"ApodDate": "%s",
		"MediaType": "%s",
		"Description": "%s",
		"ActionDate": "%s"
	}`, testUserID, action, testApodURL, testApodName, testApodDate, testMediaType, testDescription, testActionDate)

	return testJSON
}

func ApodActionTestHelper(t *testing.T, payload []byte, expected string, expectedStatus int) {
	req, err := http.NewRequest("POST", "/users/action", bytes.NewBuffer(payload))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", testToken)

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(userAction.ApplyAction)
	handler.ServeHTTP(rr, req)
	if status := rr.Code; status != expectedStatus {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, expectedStatus)
	}

	if rr.Body.String() != expected {
		t.Errorf("handler returned unexpected body: got %v want %v",
			rr.Body.String(), expected)
	}
}

func TestSaveApodHappyPath(t *testing.T) {
	var payload = []byte(getTestJson("save"))

	ApodActionTestHelper(t, payload, "Successfully saved APOD", http.StatusOK)
}

func TestUnsaveApodHappyPath(t *testing.T) {
	var payload = []byte(getTestJson("unsave"))

	ApodActionTestHelper(t, payload, "Successfully unsaved APOD", http.StatusOK)
}

func TestLikeApodHappyPath(t *testing.T) {
	var payload = []byte(getTestJson("like"))

	ApodActionTestHelper(t, payload, "Successfully liked APOD", http.StatusOK)
}

func TestUnlikeApodHappyPath(t *testing.T) {
	var payload = []byte(getTestJson("unlike"))

	ApodActionTestHelper(t, payload, "Successfully unliked APOD", http.StatusOK)
}

func TestInvalidAction(t *testing.T) {
	var payload = []byte(getTestJson("invalidAction"))

	ApodActionTestHelper(t, payload, "Invalid action in request body", http.StatusBadRequest)
}
