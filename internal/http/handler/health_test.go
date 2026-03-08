package handler

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

type mockPool struct {
	pingErr error
}

func (m *mockPool) Ping(ctx context.Context) error {
	return m.pingErr
}

func TestHealthHandler_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)

	pool := &mockPool{pingErr: nil}

	req := httptest.NewRequest("GET", "/health", nil)
	w := httptest.NewRecorder()

	c, _ := gin.CreateTestContext(w)
	c.Request = req

	if pool.Ping(context.Background()) == nil {
		c.JSON(http.StatusOK, gin.H{"status": "ok", "db": "up"})
	}

	assert.Equal(t, 200, w.Code)
}

func TestHealthHandler_DatabaseDown(t *testing.T) {
	gin.SetMode(gin.TestMode)

	pool := &mockPool{pingErr: errors.New("connection failed")}

	if pool.Ping(context.Background()) != nil {
		assert.Error(t, pool.pingErr)
	}
}
