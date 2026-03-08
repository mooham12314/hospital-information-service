package auth

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestJWT_GenerateAndParse_Success(t *testing.T) {
	manager := NewManager("test-secret-key", 60)

	token, err := manager.GenerateToken(1, 10, "alice", "HOSPITAL_A")
	assert.NoError(t, err)
	assert.NotEmpty(t, token)

	claims, err := manager.ParseToken(token)
	assert.NoError(t, err)
	assert.Equal(t, int64(1), claims.StaffID)
	assert.Equal(t, int64(10), claims.HospitalID)
	assert.Equal(t, "alice", claims.Username)
	assert.Equal(t, "HOSPITAL_A", claims.Hospital)
}

func TestJWT_ParseToken_InvalidToken(t *testing.T) {
	manager := NewManager("test-secret-key", 60)

	_, err := manager.ParseToken("invalid.token.here")
	assert.Error(t, err)
}

func TestJWT_ParseToken_WrongSecret(t *testing.T) {
	manager1 := NewManager("secret1", 60)
	manager2 := NewManager("secret2", 60)

	token, err := manager1.GenerateToken(1, 10, "alice", "HOSPITAL_A")
	assert.NoError(t, err)

	_, err = manager2.ParseToken(token)
	assert.Error(t, err)
}

func TestJWT_TokenExpiration(t *testing.T) {
	manager := NewManager("test-secret-key", -1)

	token, err := manager.GenerateToken(1, 10, "alice", "HOSPITAL_A")
	assert.NoError(t, err)

	time.Sleep(2 * time.Second)

	_, err = manager.ParseToken(token)
	assert.Error(t, err)
}

func TestJWT_Claims_Content(t *testing.T) {
	manager := NewManager("test-secret-key", 120)

	token, err := manager.GenerateToken(123, 456, "testuser", "TEST_HOSPITAL")
	assert.NoError(t, err)

	claims, err := manager.ParseToken(token)
	assert.NoError(t, err)

	assert.Equal(t, int64(123), claims.StaffID)
	assert.Equal(t, int64(456), claims.HospitalID)
	assert.Equal(t, "testuser", claims.Username)
	assert.Equal(t, "TEST_HOSPITAL", claims.Hospital)
	assert.NotNil(t, claims.IssuedAt)
	assert.NotNil(t, claims.ExpiresAt)
}
