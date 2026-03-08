-- Add unique constraint for hospital_id + national_id (when national_id is not null)
CREATE UNIQUE INDEX IF NOT EXISTS idx_patients_hospital_national 
ON patients (hospital_id, national_id) 
WHERE national_id IS NOT NULL;

-- Add unique constraint for hospital_id + passport_id (when passport_id is not null)
CREATE UNIQUE INDEX IF NOT EXISTS idx_patients_hospital_passport 
ON patients (hospital_id, passport_id) 
WHERE passport_id IS NOT NULL;

-- Insert sample hospital for testing
INSERT INTO hospitals (code, name) VALUES ('HOSPITAL_A', 'Hospital A') ON CONFLICT (code) DO NOTHING;

-- Insert sample patients for testing
INSERT INTO patients (
    hospital_id, patient_hn, national_id, passport_id,
    first_name_th, middle_name_th, last_name_th,
    first_name_en, middle_name_en, last_name_en,
    date_of_birth, phone_number, email, gender
) VALUES 
(
    (SELECT id FROM hospitals WHERE code = 'HOSPITAL_A' LIMIT 1),
    'HN001', '1234567890123', NULL,
    'สมชาย', NULL, 'ใจดี',
    'Somchai', NULL, 'Jaidee',
    '1990-01-15', '0812345678', 'somchai@example.com', 'M'
),
(
    (SELECT id FROM hospitals WHERE code = 'HOSPITAL_A' LIMIT 1),
    'HN002', '9876543210987', NULL,
    'สมหญิง', NULL, 'รักดี',
    'Somying', NULL, 'Rakdee',
    '1995-05-20', '0898765432', 'somying@example.com', 'F'
)
ON CONFLICT DO NOTHING;
