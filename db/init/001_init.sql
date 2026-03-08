CREATE TABLE IF NOT EXISTS hospitals (
    id BIGSERIAL PRIMARY KEY,
    code VARCHAR(50) UNIQUE NOT NULL,
    name VARCHAR(255) NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS staff (
    id BIGSERIAL PRIMARY KEY,
    username VARCHAR(100) NOT NULL,
    password_hash TEXT NOT NULL,
    hospital_id BIGINT NOT NULL REFERENCES hospitals(id),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    UNIQUE (username, hospital_id)
);

CREATE TABLE IF NOT EXISTS patients (
    id BIGSERIAL PRIMARY KEY,
    hospital_id BIGINT NOT NULL REFERENCES hospitals(id),
    patient_hn VARCHAR(100),
    national_id VARCHAR(20),
    passport_id VARCHAR(50),
    first_name_th VARCHAR(255),
    middle_name_th VARCHAR(255),
    last_name_th VARCHAR(255),
    first_name_en VARCHAR(255),
    middle_name_en VARCHAR(255),
    last_name_en VARCHAR(255),
    date_of_birth DATE,
    phone_number VARCHAR(30),
    email VARCHAR(255),
    gender CHAR(1),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT chk_gender CHECK (gender IN ('M', 'F') OR gender IS NULL)
);

CREATE INDEX IF NOT EXISTS idx_staff_hospital_id ON staff (hospital_id);
CREATE INDEX IF NOT EXISTS idx_patients_hospital_id ON patients (hospital_id);
CREATE INDEX IF NOT EXISTS idx_patients_national_id ON patients (national_id);
CREATE INDEX IF NOT EXISTS idx_patients_passport_id ON patients (passport_id);
