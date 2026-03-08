package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Patient struct {
	ID           int64      `json:"id"`
	HospitalID   int64      `json:"hospital_id"`
	PatientHN    *string    `json:"patient_hn,omitempty"`
	NationalID   *string    `json:"national_id,omitempty"`
	PassportID   *string    `json:"passport_id,omitempty"`
	FirstNameTH  *string    `json:"first_name_th,omitempty"`
	MiddleNameTH *string    `json:"middle_name_th,omitempty"`
	LastNameTH   *string    `json:"last_name_th,omitempty"`
	FirstNameEN  *string    `json:"first_name_en,omitempty"`
	MiddleNameEN *string    `json:"middle_name_en,omitempty"`
	LastNameEN   *string    `json:"last_name_en,omitempty"`
	DateOfBirth  *time.Time `json:"date_of_birth,omitempty"`
	PhoneNumber  *string    `json:"phone_number,omitempty"`
	Email        *string    `json:"email,omitempty"`
	Gender       *string    `json:"gender,omitempty"`
	CreatedAt    time.Time  `json:"created_at"`
	UpdatedAt    time.Time  `json:"updated_at"`
}

type PatientSearchCriteria struct {
	HospitalID  int64
	NationalID  *string
	PassportID  *string
	FirstName   *string
	MiddleName  *string
	LastName    *string
	DateOfBirth *string
	PhoneNumber *string
	Email       *string
}

type PatientRepository struct {
	db *pgxpool.Pool
}

func NewPatientRepository(db *pgxpool.Pool) *PatientRepository {
	return &PatientRepository{db: db}
}

func (r *PatientRepository) Search(ctx context.Context, criteria PatientSearchCriteria) ([]Patient, error) {
	query := `
		SELECT id, hospital_id, patient_hn, national_id, passport_id,
		       first_name_th, middle_name_th, last_name_th,
		       first_name_en, middle_name_en, last_name_en,
		       date_of_birth, phone_number, email, gender,
		       created_at, updated_at
		FROM patients
		WHERE hospital_id = $1
	`
	args := []interface{}{criteria.HospitalID}
	argIdx := 2

	if criteria.NationalID != nil && *criteria.NationalID != "" {
		query += fmt.Sprintf(" AND national_id = $%d", argIdx)
		args = append(args, *criteria.NationalID)
		argIdx++
	}

	if criteria.PassportID != nil && *criteria.PassportID != "" {
		query += fmt.Sprintf(" AND passport_id = $%d", argIdx)
		args = append(args, *criteria.PassportID)
		argIdx++
	}

	if criteria.FirstName != nil && *criteria.FirstName != "" {
		query += fmt.Sprintf(" AND (LOWER(first_name_th) LIKE LOWER($%d) OR LOWER(first_name_en) LIKE LOWER($%d))", argIdx, argIdx)
		args = append(args, "%"+*criteria.FirstName+"%")
		argIdx++
	}

	if criteria.MiddleName != nil && *criteria.MiddleName != "" {
		query += fmt.Sprintf(" AND (LOWER(middle_name_th) LIKE LOWER($%d) OR LOWER(middle_name_en) LIKE LOWER($%d))", argIdx, argIdx)
		args = append(args, "%"+*criteria.MiddleName+"%")
		argIdx++
	}

	if criteria.LastName != nil && *criteria.LastName != "" {
		query += fmt.Sprintf(" AND (LOWER(last_name_th) LIKE LOWER($%d) OR LOWER(last_name_en) LIKE LOWER($%d))", argIdx, argIdx)
		args = append(args, "%"+*criteria.LastName+"%")
		argIdx++
	}

	if criteria.DateOfBirth != nil && *criteria.DateOfBirth != "" {
		query += fmt.Sprintf(" AND date_of_birth = $%d", argIdx)
		args = append(args, *criteria.DateOfBirth)
		argIdx++
	}

	if criteria.PhoneNumber != nil && *criteria.PhoneNumber != "" {
		query += fmt.Sprintf(" AND phone_number = $%d", argIdx)
		args = append(args, *criteria.PhoneNumber)
		argIdx++
	}

	if criteria.Email != nil && *criteria.Email != "" {
		query += fmt.Sprintf(" AND LOWER(email) = LOWER($%d)", argIdx)
		args = append(args, *criteria.Email)
		argIdx++
	}

	query += " ORDER BY created_at DESC LIMIT 100"

	rows, err := r.db.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var patients []Patient
	for rows.Next() {
		var p Patient
		err := rows.Scan(
			&p.ID, &p.HospitalID, &p.PatientHN, &p.NationalID, &p.PassportID,
			&p.FirstNameTH, &p.MiddleNameTH, &p.LastNameTH,
			&p.FirstNameEN, &p.MiddleNameEN, &p.LastNameEN,
			&p.DateOfBirth, &p.PhoneNumber, &p.Email, &p.Gender,
			&p.CreatedAt, &p.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		patients = append(patients, p)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return patients, nil
}

func (r *PatientRepository) CreateOrUpdate(ctx context.Context, patient Patient) (int64, error) {
	var patientID int64
	err := r.db.QueryRow(ctx, `
		INSERT INTO patients (
			hospital_id, patient_hn, national_id, passport_id,
			first_name_th, middle_name_th, last_name_th,
			first_name_en, middle_name_en, last_name_en,
			date_of_birth, phone_number, email, gender
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14)
		ON CONFLICT (hospital_id, national_id)
		WHERE national_id IS NOT NULL
		DO UPDATE SET
			patient_hn = EXCLUDED.patient_hn,
			passport_id = EXCLUDED.passport_id,
			first_name_th = EXCLUDED.first_name_th,
			middle_name_th = EXCLUDED.middle_name_th,
			last_name_th = EXCLUDED.last_name_th,
			first_name_en = EXCLUDED.first_name_en,
			middle_name_en = EXCLUDED.middle_name_en,
			last_name_en = EXCLUDED.last_name_en,
			date_of_birth = EXCLUDED.date_of_birth,
			phone_number = EXCLUDED.phone_number,
			email = EXCLUDED.email,
			gender = EXCLUDED.gender,
			updated_at = NOW()
		RETURNING id
	`,
		patient.HospitalID, patient.PatientHN, patient.NationalID, patient.PassportID,
		patient.FirstNameTH, patient.MiddleNameTH, patient.LastNameTH,
		patient.FirstNameEN, patient.MiddleNameEN, patient.LastNameEN,
		patient.DateOfBirth, patient.PhoneNumber, patient.Email, patient.Gender,
	).Scan(&patientID)

	if err != nil {
		return 0, err
	}
	return patientID, nil
}
