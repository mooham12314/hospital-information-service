package handler

import (
	"bytes"
	"context"
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
