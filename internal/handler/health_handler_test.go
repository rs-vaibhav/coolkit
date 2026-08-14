package handler

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func init() {
	gin.SetMode(gin.TestMode)
}

func TestHealth(t *testing.T) {
	h := NewHealthHandler(nil) // nil db is fine — Health doesn't use it

	router := gin.New()
	router.GET("/api/v1/health", h.Health)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/v1/health", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestHealthResponse(t *testing.T) {
	h := NewHealthHandler(nil)

	router := gin.New()
	router.GET("/api/v1/health", h.Health)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/v1/health", nil)
	router.ServeHTTP(w, req)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)

	assert.True(t, response["success"].(bool))

	data := response["data"].(map[string]interface{})
	assert.Equal(t, "ok", data["status"])
	assert.Equal(t, "0.1.0", data["version"])
}
