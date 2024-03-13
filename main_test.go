package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/asdutoit/gotraining/section11/db"
	"github.com/asdutoit/gotraining/section11/models"
	"github.com/asdutoit/gotraining/section11/routes"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

type ResponseToken struct {
	Token string `json:"token"`
}

func TestMain(m *testing.M) {
	checkEnv()
	os.Setenv("GO_ENV", "test")
	db.InitDB()

	// Run the tests
	code := m.Run()

	// Close the database connection
	db.DB.Close()

	// Exit with the code returned from running the tests
	os.Exit(code)
}

func TestHealthCheck(t *testing.T) {
	// Switch to test mode so you don't get such noisy output
	gin.SetMode(gin.TestMode)

	// Setup your router, just like you did in your main function
	router := routes.SetupRouter()

	// Perform a GET request with that handler.
	req, _ := http.NewRequest(http.MethodGet, "/healthcheck", nil)
	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	// Check the status code is what we expect.
	assert.Equal(t, http.StatusOK, resp.Code)

	// You can check the response body here
	// assert.Equal(t, "{\"status\":\"ok\"}", resp.Body.String())
	var response map[string]string
	json.Unmarshal(resp.Body.Bytes(), &response)
	assert.Equal(t, "ok", response["status"])
}

func TestDBInit(t *testing.T) {
	assert.NotNil(t, db.DB, "Database should be initialized")
	err := db.DB.Ping()
	assert.NoError(t, err, "Database should be reachable")
}

func TestRegisterUser(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := routes.SetupRouter()

	// Register a new user

	user := models.User{
		Username: "testuser",
		Email:    "testuser@gmail.com",
		Password: "password",
	}

	userJSON, err := json.Marshal(&user)
	if err != nil {
		t.Fatalf("Could not convert user data to JSON: %v", err)
	}

	req, err := http.NewRequest(http.MethodPost, "/signup", bytes.NewBuffer(userJSON))
	if err != nil {
		t.Fatalf("Could not create request: %v", err)
	}

	req.Header.Set("Content-Type", "application/json")

	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusCreated, resp.Code, "Registration should succeed")
}

func TestLoginUser(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := routes.SetupRouter()

	user := models.User{
		Email:    "testuser@gmail.com",
		Password: "password",
	}

	userJSON, err := json.Marshal(&user)
	if err != nil {
		t.Fatalf("Could not convert user data to JSON: %v", err)
	}

	req, err := http.NewRequest(http.MethodPost, "/login", bytes.NewBuffer(userJSON))
	if err != nil {
		t.Fatalf("Could not process login request: %v", err)
	}

	req.Header.Set("Content-Type", "application/json")

	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusOK, resp.Code, "User Login should succeed")
}

func TestDeleteUser(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := routes.SetupRouter()

	user := models.User{
		Username: "testuser",
		Email:    "testuser@gmail.com",
		Password: "password",
	}

	userJSON, err := json.Marshal(&user)
	if err != nil {
		t.Fatalf("Could not convert user data to JSON: %v", err)
	}

	req, err := http.NewRequest(http.MethodPost, "/login", bytes.NewBuffer(userJSON))
	if err != nil {
		t.Fatalf("Could not process login request: %v", err)
	}

	req.Header.Set("Content-Type", "application/json")

	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	var respBody ResponseToken
	err = json.Unmarshal(resp.Body.Bytes(), &respBody)
	if err != nil {
		t.Fatalf("Could not parse response body: %v", err)
	}

	token := respBody.Token

	req, err = http.NewRequest(http.MethodPost, "/deleteUser", bytes.NewBuffer(userJSON))
	if err != nil {
		t.Fatalf("Could not process delete request: %v", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", token)

	resp = httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusOK, resp.Code, "User delete should succeed")

}
