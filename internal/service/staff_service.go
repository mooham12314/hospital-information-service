package service

import (
	"context"
	"errors"
	"strings"

	"github.com/jackc/pgx/v5"
	"github.com/mooham12314/hospital-information-service/internal/auth"
	"github.com/mooham12314/hospital-information-service/internal/repository"
	"golang.org/x/crypto/bcrypt"
)

var (
	ErrInvalidInput       = errors.New("invalid input")
	ErrDuplicateStaff     = errors.New("staff already exists")
	ErrInvalidCredentials = errors.New("invalid credentials")
)

type CreateStaffRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Hospital string `json:"hospital"`
}

type LoginStaffRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Hospital string `json:"hospital"`
}

type AuthResult struct {
	StaffID    int64  `json:"staff_id"`
	Username   string `json:"username"`
	HospitalID int64  `json:"hospital_id"`
	Hospital   string `json:"hospital"`
	Token      string `json:"token"`
}

type StaffService struct {
	repo       *repository.StaffRepository
	jwtManager *auth.Manager
}

func NewStaffService(repo *repository.StaffRepository, jwtManager *auth.Manager) *StaffService {
	return &StaffService{repo: repo, jwtManager: jwtManager}
}

func (s *StaffService) CreateStaff(ctx context.Context, req CreateStaffRequest) (AuthResult, error) {
	if err := validateCredentials(req.Username, req.Password, req.Hospital); err != nil {
		return AuthResult{}, err
	}

	hospitalID, err := s.repo.GetOrCreateHospital(ctx, req.Hospital)
	if err != nil {
		return AuthResult{}, err
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return AuthResult{}, err
	}

	staffID, err := s.repo.CreateStaff(ctx, strings.TrimSpace(req.Username), string(hash), hospitalID)
	if err != nil {
		if errors.Is(err, repository.ErrDuplicateStaff) {
			return AuthResult{}, ErrDuplicateStaff
		}
		return AuthResult{}, err
	}

	token, err := s.jwtManager.GenerateToken(staffID, hospitalID, strings.TrimSpace(req.Username), strings.TrimSpace(req.Hospital))
	if err != nil {
		return AuthResult{}, err
	}

	return AuthResult{
		StaffID:    staffID,
		Username:   strings.TrimSpace(req.Username),
		HospitalID: hospitalID,
		Hospital:   strings.TrimSpace(req.Hospital),
		Token:      token,
	}, nil
}

func (s *StaffService) Login(ctx context.Context, req LoginStaffRequest) (AuthResult, error) {
	if err := validateCredentials(req.Username, req.Password, req.Hospital); err != nil {
		return AuthResult{}, err
	}

	staff, err := s.repo.FindStaffByUsernameAndHospital(ctx, strings.TrimSpace(req.Username), strings.TrimSpace(req.Hospital))
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return AuthResult{}, ErrInvalidCredentials
		}
		return AuthResult{}, err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(staff.PasswordHash), []byte(req.Password)); err != nil {
		return AuthResult{}, ErrInvalidCredentials
	}

	token, err := s.jwtManager.GenerateToken(staff.ID, staff.HospitalID, staff.Username, staff.HospitalCode)
	if err != nil {
		return AuthResult{}, err
	}

	return AuthResult{
		StaffID:    staff.ID,
		Username:   staff.Username,
		HospitalID: staff.HospitalID,
		Hospital:   staff.HospitalCode,
		Token:      token,
	}, nil
}

func validateCredentials(username, password, hospital string) error {
	if strings.TrimSpace(username) == "" || strings.TrimSpace(password) == "" || strings.TrimSpace(hospital) == "" {
		return ErrInvalidInput
	}
	if len(strings.TrimSpace(password)) < 8 {
		return ErrInvalidInput
	}
	return nil
}
