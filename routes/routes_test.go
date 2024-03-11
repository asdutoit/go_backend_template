package routes_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/asdutoit/gotraining/section11/models"
	"github.com/asdutoit/gotraining/section11/routes"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

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

func TestRegisterUser(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := routes.SetupRouter()

	// Register a new user

	var user models.User
	user.Username = "testuser"
	user.Email = "testuser@gmail.com"
	user.Password = "password"

	userJSON, err := json.Marshal(user)
	if err != nil {
		t.Fatalf("Could not convert user data to JSON: %v", err)
	}

	req, err := http.NewRequest(http.MethodPost, "/signup", bytes.NewBuffer(userJSON))
	if err != nil {
		t.Fatalf("Could not create request: %v", err)
	}

	req.Header.Set("Content-Type", "application/json")

	resp := httptest.NewRecorder()
	fmt.Println("Request: ", req)
	router.ServeHTTP(resp, req)

	fmt.Println("Response: ", resp.Body.String())

	assert.Equal(t, http.StatusOK, resp.Code, "Registration should succeed")
}
