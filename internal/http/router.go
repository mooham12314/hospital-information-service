package http

import (
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/mooham12314/hospital-information-service/internal/auth"
	"github.com/mooham12314/hospital-information-service/internal/client"
	"github.com/mooham12314/hospital-information-service/internal/config"
	"github.com/mooham12314/hospital-information-service/internal/http/handler"
	"github.com/mooham12314/hospital-information-service/internal/http/middleware"
	"github.com/mooham12314/hospital-information-service/internal/repository"
	"github.com/mooham12314/hospital-information-service/internal/service"
)

func NewRouter(db *pgxpool.Pool, cfg config.Config) *gin.Engine {
	r := gin.New()
	r.Use(gin.Logger(), gin.Recovery())

	healthHandler := handler.NewHealthHandler(db)
	r.GET("/health", healthHandler.Health)

	jwtManager := auth.NewManager(cfg.JWTSecret, cfg.JWTTTLMin)

	staffRepo := repository.NewStaffRepository(db)
	staffService := service.NewStaffService(staffRepo, jwtManager)
	staffHandler := handler.NewStaffHandler(staffService)

	patientRepo := repository.NewPatientRepository(db)
	hospitalAClient := client.NewHospitalAClient(cfg.HospitalAAPIURL)
	patientService := service.NewPatientService(patientRepo, hospitalAClient)
	patientHandler := handler.NewPatientHandler(patientService)

	api := r.Group("/api/v1")
	{
		api.GET("/health", healthHandler.Health)
		api.POST("/staff/create", staffHandler.Create)
		api.POST("/staff/login", staffHandler.Login)

		secured := api.Group("")
		secured.Use(middleware.RequireAuth(jwtManager))
		secured.POST("/patient/search", patientHandler.Search)
	}

	r.POST("/staff/create", staffHandler.Create)
	r.POST("/staff/login", staffHandler.Login)

	securedRoot := r.Group("")
	securedRoot.Use(middleware.RequireAuth(jwtManager))
	securedRoot.POST("/patient/search", patientHandler.Search)

	return r
}
