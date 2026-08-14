package response

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

func getTestContext() (*gin.Context, *httptest.ResponseRecorder) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	return c, w
}

func TestSuccess(t *testing.T) {
	c, w := getTestContext()
	
	Success(c, http.StatusOK, gin.H{"message": "test"})
	
	assert.Equal(t, http.StatusOK, w.Code)
	
	var res map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &res)
	assert.NoError(t, err)
	
	assert.True(t, res["success"].(bool))
	data := res["data"].(map[string]interface{})
	assert.Equal(t, "test", data["message"])
}

func TestError(t *testing.T) {
	c, w := getTestContext()
	
	Error(c, http.StatusBadRequest, "bad request error")
	
	assert.Equal(t, http.StatusBadRequest, w.Code)
	
	var res map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &res)
	assert.NoError(t, err)
	
	assert.False(t, res["success"].(bool))
	assert.Equal(t, "bad request error", res["error"])
}

func TestOK(t *testing.T) {
	c, w := getTestContext()
	OK(c, "test")
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestCreated(t *testing.T) {
	c, w := getTestContext()
	Created(c, "test")
	assert.Equal(t, http.StatusCreated, w.Code)
}

func TestBadRequest(t *testing.T) {
	c, w := getTestContext()
	BadRequest(c, "error")
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestUnauthorized(t *testing.T) {
	c, w := getTestContext()
	Unauthorized(c, "error")
	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

func TestNotFound(t *testing.T) {
	c, w := getTestContext()
	NotFound(c, "error")
	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestInternalError(t *testing.T) {
	c, w := getTestContext()
	InternalError(c, "error")
	assert.Equal(t, http.StatusInternalServerError, w.Code)
}
