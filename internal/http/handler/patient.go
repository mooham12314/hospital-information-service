package handler

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/mooham12314/hospital-information-service/internal/http/middleware"
	"github.com/mooham12314/hospital-information-service/internal/service"
)

type PatientServiceInterface interface {
	Search(ctx context.Context, hospitalID int64, req service.PatientSearchRequest) (service.PatientSearchResponse, error)
}

type PatientHandler struct {
	service PatientServiceInterface
}

func NewPatientHandler(service PatientServiceInterface) *PatientHandler {
	return &PatientHandler{service: service}
}

func (h *PatientHandler) Search(c *gin.Context) {
	hospitalID, exists := c.Get(middleware.ContextHospitalIDKey)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	var req service.PatientSearchRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}

	result, err := h.service.Search(c.Request.Context(), hospitalID.(int64), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to search patients"})
		return
	}

	c.JSON(http.StatusOK, result)
}
