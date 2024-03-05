package routes_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

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
