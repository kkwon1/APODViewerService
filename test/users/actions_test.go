package users

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/kkwon1/APODViewerService/api/users"
	"github.com/kkwon1/APODViewerService/models"
	"go.mongodb.org/mongo-driver/mongo"
)

/*
	Mock dependencies
*/
type mockTokenVerifier struct{}

func (m *mockTokenVerifier) VerifyToken(ctx context.Context, idToken string) (bool, error) {
	return true, nil
}

type mockUserActionDAO struct {
	savesCollection *mongo.Collection
	likesCollection *mongo.Collection
}

//TODO: Make a before/after to mock the functions so we can test error paths as well
func (m *mockUserActionDAO) SaveApod(ctx context.Context, userAction *models.UserAction) error {
	return nil
}
func (m *mockUserActionDAO) UnsaveApod(ctx context.Context, userAction *models.UserAction) error {
	return nil
}
func (m *mockUserActionDAO) LikeApod(ctx context.Context, userAction *models.UserAction) error {
	return nil
}
func (m *mockUserActionDAO) UnlikeApod(ctx context.Context, userAction *models.UserAction) error {
	return nil
}

var userAction users.UserAction = users.NewUserAction(&mockTokenVerifier{}, &mockUserActionDAO{})

const testUserID = "DB_TEST_ID"
const testApodURL = "https://apod.nasa.gov/apod/image/2005/c2020_f8_2020_05_02dp_1024.jpg"
const testApodName = "Long Tailed Comet SWAN"
const testApodDate = "2020-05-08"
const testMediaType = "image"
const testDescription = "Test description"
const testActionDate = "2020-05-11T18:19:58.747Z"

const testToken = "Bearer mock-token"

func getTestJsonString(action string) string {
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

func getTestOutput() *models.ApodObject {
	testJSON := fmt.Sprintf(`{
		"UserID": "%s",
		"ApodURL": "%s",
		"ApodName": "%s",
		"ApodDate": "%s",
		"MediaType": "%s",
		"Description": "%s",
		"ActionDate": "%s"
	}`, testUserID, testApodURL, testApodName, testApodDate, testMediaType, testDescription, testActionDate)

	testOutput := models.ApodObject{}
	json.Unmarshal([]byte(testJSON), &testOutput)
	return &testOutput
}

func ApodActionTestHelper(t *testing.T, payload []byte, expected *models.ApodObject, expectedStatus int) {
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

	var apodObject *models.ApodObject
	_ = json.NewDecoder(rr.Body).Decode(&apodObject)

	if !reflect.DeepEqual(apodObject, expected) {
		t.Errorf("handler returned unexpected body: got %v want %v",
			apodObject, expected)
	}
}

func TestSaveApodHappyPath(t *testing.T) {
	var payload = []byte(getTestJsonString("save"))

	ApodActionTestHelper(t, payload, getTestOutput(), http.StatusOK)
}

func TestUnsaveApodHappyPath(t *testing.T) {
	var payload = []byte(getTestJsonString("unsave"))

	ApodActionTestHelper(t, payload, getTestOutput(), http.StatusOK)
}

func TestLikeApodHappyPath(t *testing.T) {
	var payload = []byte(getTestJsonString("like"))

	ApodActionTestHelper(t, payload, getTestOutput(), http.StatusOK)
}

func TestUnlikeApodHappyPath(t *testing.T) {
	var payload = []byte(getTestJsonString("unlike"))

	ApodActionTestHelper(t, payload, getTestOutput(), http.StatusOK)
}

func TestInvalidAction(t *testing.T) {
	var payload = []byte(getTestJsonString("invalidAction"))

	ApodActionTestHelper(t, payload, nil, http.StatusBadRequest)
}
