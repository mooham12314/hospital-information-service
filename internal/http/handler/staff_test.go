package handler

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/mooham12314/hospital-information-service/internal/service"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockStaffService struct {
	mock.Mock
}

func (m *MockStaffService) CreateStaff(ctx context.Context, req service.CreateStaffRequest) (service.AuthResult, error) {
	args := m.Called(ctx, req)
	return args.Get(0).(service.AuthResult), args.Error(1)
}

func (m *MockStaffService) Login(ctx context.Context, req service.LoginStaffRequest) (service.AuthResult, error) {
	args := m.Called(ctx, req)
	return args.Get(0).(service.AuthResult), args.Error(1)
}

func TestStaffHandler_Create(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockSvc := new(MockStaffService)
	handler := NewStaffHandler(mockSvc)

	expected := service.AuthResult{Token: "test-token"}
	mockSvc.On("CreateStaff", mock.Anything, mock.Anything).Return(expected, nil)

	body := `{"username":"alice","password":"password123","hospital":"HOSPITAL_A"}`
	req := httptest.NewRequest("POST", "/staff/create", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router := gin.New()
	router.POST("/staff/create", handler.Create)
	router.ServeHTTP(w, req)

	assert.Equal(t, 201, w.Code)
	var result service.AuthResult
	json.Unmarshal(w.Body.Bytes(), &result)
	assert.Equal(t, "test-token", result.Token)
}

func TestStaffHandler_Login(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockSvc := new(MockStaffService)
	handler := NewStaffHandler(mockSvc)

	expected := service.AuthResult{Token: "login-token"}
	mockSvc.On("Login", mock.Anything, mock.Anything).Return(expected, nil)

	body := `{"username":"alice","password":"password123","hospital":"HOSPITAL_A"}`
	req := httptest.NewRequest("POST", "/staff/login", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router := gin.New()
	router.POST("/staff/login", handler.Login)
	router.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)
}
