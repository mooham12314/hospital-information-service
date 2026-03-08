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

func TestStaffHandler_Create_InvalidJSON(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockSvc := new(MockStaffService)
	handler := NewStaffHandler(mockSvc)

	body := `{"username":"alice","password":invalid}`
	req := httptest.NewRequest("POST", "/staff/create", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router := gin.New()
	router.POST("/staff/create", handler.Create)
	router.ServeHTTP(w, req)

	assert.Equal(t, 400, w.Code)
	assert.Contains(t, w.Body.String(), "invalid request body")
}

func TestStaffHandler_Create_InvalidInput(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockSvc := new(MockStaffService)
	handler := NewStaffHandler(mockSvc)

	mockSvc.On("CreateStaff", mock.Anything, mock.Anything).Return(service.AuthResult{}, service.ErrInvalidInput)

	body := `{"username":"alice","password":"short","hospital":"HOSPITAL_A"}`
	req := httptest.NewRequest("POST", "/staff/create", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router := gin.New()
	router.POST("/staff/create", handler.Create)
	router.ServeHTTP(w, req)

	assert.Equal(t, 400, w.Code)
	assert.Contains(t, w.Body.String(), "password(min 8)")
}

func TestStaffHandler_Create_DuplicateStaff(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockSvc := new(MockStaffService)
	handler := NewStaffHandler(mockSvc)

	mockSvc.On("CreateStaff", mock.Anything, mock.Anything).Return(service.AuthResult{}, service.ErrDuplicateStaff)

	body := `{"username":"alice","password":"password123","hospital":"HOSPITAL_A"}`
	req := httptest.NewRequest("POST", "/staff/create", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router := gin.New()
	router.POST("/staff/create", handler.Create)
	router.ServeHTTP(w, req)

	assert.Equal(t, 409, w.Code)
	assert.Contains(t, w.Body.String(), "already exists")
}

func TestStaffHandler_Create_InternalError(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockSvc := new(MockStaffService)
	handler := NewStaffHandler(mockSvc)

	mockSvc.On("CreateStaff", mock.Anything, mock.Anything).Return(service.AuthResult{}, assert.AnError)

	body := `{"username":"alice","password":"password123","hospital":"HOSPITAL_A"}`
	req := httptest.NewRequest("POST", "/staff/create", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router := gin.New()
	router.POST("/staff/create", handler.Create)
	router.ServeHTTP(w, req)

	assert.Equal(t, 500, w.Code)
	assert.Contains(t, w.Body.String(), "failed to create staff")
}

func TestStaffHandler_Login_InvalidJSON(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockSvc := new(MockStaffService)
	handler := NewStaffHandler(mockSvc)

	body := `{invalid json`
	req := httptest.NewRequest("POST", "/staff/login", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router := gin.New()
	router.POST("/staff/login", handler.Login)
	router.ServeHTTP(w, req)

	assert.Equal(t, 400, w.Code)
	assert.Contains(t, w.Body.String(), "invalid request body")
}

func TestStaffHandler_Login_InvalidCredentials(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockSvc := new(MockStaffService)
	handler := NewStaffHandler(mockSvc)

	mockSvc.On("Login", mock.Anything, mock.Anything).Return(service.AuthResult{}, service.ErrInvalidCredentials)

	body := `{"username":"alice","password":"wrongpass","hospital":"HOSPITAL_A"}`
	req := httptest.NewRequest("POST", "/staff/login", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router := gin.New()
	router.POST("/staff/login", handler.Login)
	router.ServeHTTP(w, req)

	assert.Equal(t, 401, w.Code)
	assert.Contains(t, w.Body.String(), "invalid username/password/hospital")
}

func TestStaffHandler_Login_InvalidInput(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockSvc := new(MockStaffService)
	handler := NewStaffHandler(mockSvc)

	mockSvc.On("Login", mock.Anything, mock.Anything).Return(service.AuthResult{}, service.ErrInvalidInput)

	body := `{"username":"alice","password":"short","hospital":"HOSPITAL_A"}`
	req := httptest.NewRequest("POST", "/staff/login", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router := gin.New()
	router.POST("/staff/login", handler.Login)
	router.ServeHTTP(w, req)

	assert.Equal(t, 400, w.Code)
	assert.Contains(t, w.Body.String(), "password(min 8)")
}

func TestStaffHandler_Login_InternalError(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockSvc := new(MockStaffService)
	handler := NewStaffHandler(mockSvc)

	mockSvc.On("Login", mock.Anything, mock.Anything).Return(service.AuthResult{}, assert.AnError)

	body := `{"username":"alice","password":"password123","hospital":"HOSPITAL_A"}`
	req := httptest.NewRequest("POST", "/staff/login", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router := gin.New()
	router.POST("/staff/login", handler.Login)
	router.ServeHTTP(w, req)

	assert.Equal(t, 500, w.Code)
	assert.Contains(t, w.Body.String(), "failed to login")
}
