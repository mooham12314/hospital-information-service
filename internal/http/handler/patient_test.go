package handler

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/mooham12314/hospital-information-service/internal/http/middleware"
	"github.com/mooham12314/hospital-information-service/internal/repository"
	"github.com/mooham12314/hospital-information-service/internal/service"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockPatientService struct {
	mock.Mock
}

func (m *MockPatientService) Search(ctx context.Context, hospitalID int64, req service.PatientSearchRequest) (service.PatientSearchResponse, error) {
	args := m.Called(ctx, hospitalID, req)
	return args.Get(0).(service.PatientSearchResponse), args.Error(1)
}

func TestPatientHandler_Search(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockSvc := new(MockPatientService)
	handler := NewPatientHandler(mockSvc)

	expected := service.PatientSearchResponse{Patients: []repository.Patient{}, Count: 0}
	mockSvc.On("Search", mock.Anything, mock.Anything, mock.Anything).Return(expected, nil)

	body := `{}`
	req := httptest.NewRequest("POST", "/patient/search", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router := gin.New()
	router.POST("/patient/search", func(c *gin.Context) {
		c.Set(middleware.ContextHospitalIDKey, int64(1))
		handler.Search(c)
	})
	router.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)
}

func TestPatientHandler_Search_Unauthorized(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockSvc := new(MockPatientService)
	handler := NewPatientHandler(mockSvc)

	body := `{}`
	req := httptest.NewRequest("POST", "/patient/search", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router := gin.New()
	router.POST("/patient/search", handler.Search)
	router.ServeHTTP(w, req)

	assert.Equal(t, 401, w.Code)
}

func TestPatientHandler_Search_InvalidJSON(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockSvc := new(MockPatientService)
	handler := NewPatientHandler(mockSvc)

	body := `{invalid json}`
	req := httptest.NewRequest("POST", "/patient/search", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router := gin.New()
	router.POST("/patient/search", func(c *gin.Context) {
		c.Set(middleware.ContextHospitalIDKey, int64(1))
		handler.Search(c)
	})
	router.ServeHTTP(w, req)

	assert.Equal(t, 400, w.Code)
	assert.Contains(t, w.Body.String(), "invalid request body")
}

func TestPatientHandler_Search_ServiceError(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockSvc := new(MockPatientService)
	handler := NewPatientHandler(mockSvc)

	mockSvc.On("Search", mock.Anything, mock.Anything, mock.Anything).Return(service.PatientSearchResponse{}, assert.AnError)

	body := `{}`
	req := httptest.NewRequest("POST", "/patient/search", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router := gin.New()
	router.POST("/patient/search", func(c *gin.Context) {
		c.Set(middleware.ContextHospitalIDKey, int64(1))
		handler.Search(c)
	})
	router.ServeHTTP(w, req)

	assert.Equal(t, 500, w.Code)
	assert.Contains(t, w.Body.String(), "failed to search patients")
}

func TestPatientHandler_Search_WithCriteria(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockSvc := new(MockPatientService)
	handler := NewPatientHandler(mockSvc)

	firstName := "John"
	nationalID := "1234567890123"
	patient := repository.Patient{
		ID:          1,
		HospitalID:  1,
		FirstNameEN: &firstName,
		NationalID:  &nationalID,
	}

	expected := service.PatientSearchResponse{
		Patients: []repository.Patient{patient},
		Count:    1,
	}
	mockSvc.On("Search", mock.Anything, int64(1), mock.Anything).Return(expected, nil)

	body := `{"national_id":"1234567890123","first_name":"John"}`
	req := httptest.NewRequest("POST", "/patient/search", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router := gin.New()
	router.POST("/patient/search", func(c *gin.Context) {
		c.Set(middleware.ContextHospitalIDKey, int64(1))
		handler.Search(c)
	})
	router.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)

	var response service.PatientSearchResponse
	json.Unmarshal(w.Body.Bytes(), &response)
	assert.Equal(t, 1, response.Count)
	assert.Equal(t, nationalID, *response.Patients[0].NationalID)
}
