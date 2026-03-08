package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/mooham12314/hospital-information-service/internal/auth"
	"github.com/stretchr/testify/assert"
)

func TestRequireAuth_ValidToken(t *testing.T) {
	gin.SetMode(gin.TestMode)

	manager := auth.NewManager("test-secret", 60)
	token, _ := manager.GenerateToken(1, 10, "alice", "HOSPITAL_A")

	req := httptest.NewRequest(http.MethodGet, "/protected", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()

	router := gin.New()
	router.Use(RequireAuth(manager))
	router.GET("/protected", func(c *gin.Context) {
		staffID, _ := c.Get(ContextStaffIDKey)
		hospitalID, _ := c.Get(ContextHospitalIDKey)
		username, _ := c.Get(ContextUsernameKey)
		hospital, _ := c.Get(ContextHospitalKey)

		c.JSON(http.StatusOK, gin.H{
			"staff_id":    staffID,
			"hospital_id": hospitalID,
			"username":    username,
			"hospital":    hospital,
		})
	})
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), `"staff_id":1`)
	assert.Contains(t, w.Body.String(), `"hospital_id":10`)
	assert.Contains(t, w.Body.String(), `"username":"alice"`)
	assert.Contains(t, w.Body.String(), `"hospital":"HOSPITAL_A"`)
}

func TestRequireAuth_MissingToken(t *testing.T) {
	gin.SetMode(gin.TestMode)

	manager := auth.NewManager("test-secret", 60)

	req := httptest.NewRequest(http.MethodGet, "/protected", nil)
	w := httptest.NewRecorder()

	router := gin.New()
	router.Use(RequireAuth(manager))
	router.GET("/protected", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "success"})
	})
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
	assert.Contains(t, w.Body.String(), "missing or invalid bearer token")
}

func TestRequireAuth_InvalidTokenFormat(t *testing.T) {
	gin.SetMode(gin.TestMode)

	manager := auth.NewManager("test-secret", 60)

	req := httptest.NewRequest(http.MethodGet, "/protected", nil)
	req.Header.Set("Authorization", "InvalidFormat token123")
	w := httptest.NewRecorder()

	router := gin.New()
	router.Use(RequireAuth(manager))
	router.GET("/protected", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "success"})
	})
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
	assert.Contains(t, w.Body.String(), "missing or invalid bearer token")
}

func TestRequireAuth_InvalidToken(t *testing.T) {
	gin.SetMode(gin.TestMode)

	manager := auth.NewManager("test-secret", 60)

	req := httptest.NewRequest(http.MethodGet, "/protected", nil)
	req.Header.Set("Authorization", "Bearer invalid.token.here")
	w := httptest.NewRecorder()

	router := gin.New()
	router.Use(RequireAuth(manager))
	router.GET("/protected", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "success"})
	})
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
	assert.Contains(t, w.Body.String(), "invalid or expired token")
}

func TestRequireAuth_ExpiredToken(t *testing.T) {
	gin.SetMode(gin.TestMode)

	manager := auth.NewManager("test-secret", -1)
	token, _ := manager.GenerateToken(1, 10, "alice", "HOSPITAL_A")

	req := httptest.NewRequest(http.MethodGet, "/protected", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()

	router := gin.New()
	router.Use(RequireAuth(manager))
	router.GET("/protected", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "success"})
	})
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
	assert.Contains(t, w.Body.String(), "invalid or expired token")
}
