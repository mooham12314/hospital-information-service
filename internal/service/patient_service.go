package service

import (
	"context"
	"time"

	"github.com/mooham12314/hospital-information-service/internal/client"
	"github.com/mooham12314/hospital-information-service/internal/repository"
)

type PatientSearchRequest struct {
	NationalID  *string `json:"national_id,omitempty"`
	PassportID  *string `json:"passport_id,omitempty"`
	FirstName   *string `json:"first_name,omitempty"`
	MiddleName  *string `json:"middle_name,omitempty"`
	LastName    *string `json:"last_name,omitempty"`
	DateOfBirth *string `json:"date_of_birth,omitempty"`
	PhoneNumber *string `json:"phone_number,omitempty"`
	Email       *string `json:"email,omitempty"`
}

type PatientSearchResponse struct {
	Patients []repository.Patient `json:"patients"`
	Count    int                  `json:"count"`
}

type PatientService struct {
	repo            *repository.PatientRepository
	hospitalAClient *client.HospitalAClient
}

func NewPatientService(repo *repository.PatientRepository, hospitalAClient *client.HospitalAClient) *PatientService {
	return &PatientService{
		repo:            repo,
		hospitalAClient: hospitalAClient,
	}
}

func (s *PatientService) Search(ctx context.Context, hospitalID int64, req PatientSearchRequest) (PatientSearchResponse, error) {
	criteria := repository.PatientSearchCriteria{
		HospitalID:  hospitalID,
		NationalID:  req.NationalID,
		PassportID:  req.PassportID,
		FirstName:   req.FirstName,
		MiddleName:  req.MiddleName,
		LastName:    req.LastName,
		DateOfBirth: req.DateOfBirth,
		PhoneNumber: req.PhoneNumber,
		Email:       req.Email,
	}

	patients, err := s.repo.Search(ctx, criteria)
	if err != nil {
		return PatientSearchResponse{}, err
	}

	if len(patients) == 0 && s.hospitalAClient != nil {
		var searchID string
		if req.NationalID != nil && *req.NationalID != "" {
			searchID = *req.NationalID
		} else if req.PassportID != nil && *req.PassportID != "" {
			searchID = *req.PassportID
		}

		if searchID != "" {
			externalPatient, err := s.hospitalAClient.SearchPatient(ctx, searchID)
			if err == nil && externalPatient != nil {
				patient := s.mapHospitalAToPatient(hospitalID, externalPatient)
				_, _ = s.repo.CreateOrUpdate(ctx, patient)
				patients, _ = s.repo.Search(ctx, criteria)
			}
		}
	}

	return PatientSearchResponse{
		Patients: patients,
		Count:    len(patients),
	}, nil
}

func (s *PatientService) mapHospitalAToPatient(hospitalID int64, ext *client.HospitalAPatient) repository.Patient {
	patient := repository.Patient{
		HospitalID: hospitalID,
	}

	if ext.PatientHN != "" {
		patient.PatientHN = &ext.PatientHN
	}
	if ext.NationalID != "" {
		patient.NationalID = &ext.NationalID
	}
	if ext.PassportID != "" {
		patient.PassportID = &ext.PassportID
	}
	if ext.FirstNameTH != "" {
		patient.FirstNameTH = &ext.FirstNameTH
	}
	if ext.MiddleNameTH != "" {
		patient.MiddleNameTH = &ext.MiddleNameTH
	}
	if ext.LastNameTH != "" {
		patient.LastNameTH = &ext.LastNameTH
	}
	if ext.FirstNameEN != "" {
		patient.FirstNameEN = &ext.FirstNameEN
	}
	if ext.MiddleNameEN != "" {
		patient.MiddleNameEN = &ext.MiddleNameEN
	}
	if ext.LastNameEN != "" {
		patient.LastNameEN = &ext.LastNameEN
	}
	if ext.DateOfBirth != "" {
		dob, err := time.Parse("2006-01-02", ext.DateOfBirth)
		if err == nil {
			patient.DateOfBirth = &dob
		}
	}
	if ext.PhoneNumber != "" {
		patient.PhoneNumber = &ext.PhoneNumber
	}
	if ext.Email != "" {
		patient.Email = &ext.Email
	}
	if ext.Gender != "" {
		patient.Gender = &ext.Gender
	}

	return patient
}
