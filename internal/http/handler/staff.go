package handler

import (
	"context"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/mooham12314/hospital-information-service/internal/service"
)

type StaffServiceInterface interface {
	CreateStaff(ctx context.Context, req service.CreateStaffRequest) (service.AuthResult, error)
	Login(ctx context.Context, req service.LoginStaffRequest) (service.AuthResult, error)
}

type StaffHandler struct {
	service StaffServiceInterface
}

func NewStaffHandler(service StaffServiceInterface) *StaffHandler {
	return &StaffHandler{service: service}
}

func (h *StaffHandler) Create(c *gin.Context) {
	var req service.CreateStaffRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}

	result, err := h.service.CreateStaff(c.Request.Context(), req)
	if err != nil {
		switch {
		case errors.Is(err, service.ErrInvalidInput):
			c.JSON(http.StatusBadRequest, gin.H{"error": "username, password(min 8), and hospital are required"})
		case errors.Is(err, service.ErrDuplicateStaff):
			c.JSON(http.StatusConflict, gin.H{"error": "staff already exists in this hospital"})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create staff"})
		}
		return
	}

	c.JSON(http.StatusCreated, result)
}

func (h *StaffHandler) Login(c *gin.Context) {
	var req service.LoginStaffRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}

	result, err := h.service.Login(c.Request.Context(), req)
	if err != nil {
		switch {
		case errors.Is(err, service.ErrInvalidInput):
			c.JSON(http.StatusBadRequest, gin.H{"error": "username, password(min 8), and hospital are required"})
		case errors.Is(err, service.ErrInvalidCredentials):
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid username/password/hospital"})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to login"})
		}
		return
	}

	c.JSON(http.StatusOK, result)
}
