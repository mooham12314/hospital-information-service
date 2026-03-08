package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/mooham12314/hospital-information-service/internal/auth"
)

const (
	ContextStaffIDKey    = "staff_id"
	ContextHospitalIDKey = "hospital_id"
	ContextUsernameKey   = "username"
	ContextHospitalKey   = "hospital"
)

func RequireAuth(jwtManager *auth.Manager) gin.HandlerFunc {
	return func(c *gin.Context) {
		authorization := c.GetHeader("Authorization")
		parts := strings.SplitN(authorization, " ", 2)
		if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "missing or invalid bearer token"})
			return
		}

		claims, err := jwtManager.ParseToken(parts[1])
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid or expired token"})
			return
		}

		c.Set(ContextStaffIDKey, claims.StaffID)
		c.Set(ContextHospitalIDKey, claims.HospitalID)
		c.Set(ContextUsernameKey, claims.Username)
		c.Set(ContextHospitalKey, claims.Hospital)
		c.Next()
	}
}
