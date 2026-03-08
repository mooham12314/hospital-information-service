package repository

import (
	"context"
	"errors"
	"strings"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

var ErrDuplicateStaff = errors.New("staff already exists in hospital")

type Staff struct {
	ID           int64
	Username     string
	PasswordHash string
	HospitalID   int64
	HospitalCode string
}

type StaffRepository struct {
	db *pgxpool.Pool
}

func NewStaffRepository(db *pgxpool.Pool) *StaffRepository {
	return &StaffRepository{db: db}
}

func (r *StaffRepository) GetOrCreateHospital(ctx context.Context, hospitalCode string) (int64, error) {
	code := strings.TrimSpace(hospitalCode)
	var hospitalID int64
	err := r.db.QueryRow(ctx, `
		INSERT INTO hospitals (code, name)
		VALUES ($1, $1)
		ON CONFLICT (code)
		DO UPDATE SET updated_at = NOW()
		RETURNING id
	`, code).Scan(&hospitalID)
	if err != nil {
		return 0, err
	}
	return hospitalID, nil
}

func (r *StaffRepository) CreateStaff(ctx context.Context, username, passwordHash string, hospitalID int64) (int64, error) {
	var staffID int64
	err := r.db.QueryRow(ctx, `
		INSERT INTO staff (username, password_hash, hospital_id)
		VALUES ($1, $2, $3)
		RETURNING id
	`, username, passwordHash, hospitalID).Scan(&staffID)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return 0, ErrDuplicateStaff
		}
		return 0, err
	}
	return staffID, nil
}

func (r *StaffRepository) FindStaffByUsernameAndHospital(ctx context.Context, username, hospitalCode string) (Staff, error) {
	var staff Staff
	err := r.db.QueryRow(ctx, `
		SELECT s.id, s.username, s.password_hash, s.hospital_id, h.code
		FROM staff s
		JOIN hospitals h ON h.id = s.hospital_id
		WHERE s.username = $1 AND h.code = $2
		LIMIT 1
	`, username, strings.TrimSpace(hospitalCode)).Scan(
		&staff.ID,
		&staff.Username,
		&staff.PasswordHash,
		&staff.HospitalID,
		&staff.HospitalCode,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return Staff{}, pgx.ErrNoRows
		}
		return Staff{}, err
	}
	return staff, nil
}
