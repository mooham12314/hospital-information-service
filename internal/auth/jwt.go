package auth

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type Manager struct {
	secretKey []byte
	ttl       time.Duration
}

type Claims struct {
	StaffID    int64  `json:"staff_id"`
	HospitalID int64  `json:"hospital_id"`
	Username   string `json:"username"`
	Hospital   string `json:"hospital"`
	jwt.RegisteredClaims
}

func NewManager(secret string, ttlMinutes int) *Manager {
	return &Manager{
		secretKey: []byte(secret),
		ttl:       time.Duration(ttlMinutes) * time.Minute,
	}
}

func (m *Manager) GenerateToken(staffID, hospitalID int64, username, hospital string) (string, error) {
	now := time.Now()
	claims := Claims{
		StaffID:    staffID,
		HospitalID: hospitalID,
		Username:   username,
		Hospital:   hospital,
		RegisteredClaims: jwt.RegisteredClaims{
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(now.Add(m.ttl)),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(m.secretKey)
}

func (m *Manager) ParseToken(tokenString string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (any, error) {
		if token.Method != jwt.SigningMethodHS256 {
			return nil, errors.New("unexpected signing method")
		}
		return m.secretKey, nil
	})
	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(*Claims)
	if !ok || !token.Valid {
		return nil, errors.New("invalid token claims")
	}
	return claims, nil
}
