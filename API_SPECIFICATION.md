# Hospital Information Service - API Specification

## Overview
REST API for Hospital Information Middleware system that enables staff to search and manage patient information across multiple Hospital Information Systems (HIS).

## Base URL
```
http://localhost:8088/api/v1
```

## Authentication
All protected endpoints require Bearer token authentication in the Authorization header:
```
Authorization: Bearer <jwt_token>
```

---

## API Endpoints

### 1. Health Check

#### GET /health
Check service and database health status.

**Authentication:** None required

**Response (200 OK):**
```json
{
  "status": "ok",
  "db": "up"
}
```

**Response (503 Service Unavailable):**
```json
{
  "status": "unhealthy",
  "db": "down"
}
```

---

### 2. Staff Management

#### POST /staff/create
Create a new staff account and receive authentication token.

**Authentication:** None required

**Request Body:**
```json
{
  "username": "alice",
  "password": "password123",
  "hospital": "HOSPITAL_A"
}
```

**Validation Rules:**
- `username`: Required, string
- `password`: Required, minimum 8 characters
- `hospital`: Required, hospital code string

**Response (201 Created):**
```json
{
  "staff_id": 1,
  "username": "alice",
  "hospital_id": 1,
  "hospital": "HOSPITAL_A",
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
}
```

**Error Responses:**

- **400 Bad Request** - Invalid input
```json
{
  "error": "username, password(min 8), and hospital are required"
}
```

- **409 Conflict** - Staff already exists
```json
{
  "error": "staff already exists in this hospital"
}
```

- **500 Internal Server Error**
```json
{
  "error": "failed to create staff"
}
```

---

#### POST /staff/login
Authenticate staff and receive JWT token.

**Authentication:** None required

**Request Body:**
```json
{
  "username": "alice",
  "password": "password123",
  "hospital": "HOSPITAL_A"
}
```

**Validation Rules:**
- `username`: Required, string
- `password`: Required, string
- `hospital`: Required, hospital code string

**Response (200 OK):**
```json
{
  "staff_id": 1,
  "username": "alice",
  "hospital_id": 1,
  "hospital": "HOSPITAL_A",
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
}
```

**Error Responses:**

- **400 Bad Request** - Invalid input
```json
{
  "error": "username, password(min 8), and hospital are required"
}
```

- **401 Unauthorized** - Invalid credentials
```json
{
  "error": "invalid username/password/hospital"
}
```

- **500 Internal Server Error**
```json
{
  "error": "failed to login"
}
```

---

### 3. Patient Search

#### POST /patient/search
Search for patients in the same hospital as the authenticated staff member.

**Authentication:** Required (Bearer token)

**Request Body (all fields optional):**
```json
{
  "national_id": "1234567890123",
  "passport_id": "AB1234567",
  "first_name": "John",
  "middle_name": "Michael",
  "last_name": "Doe",
  "date_of_birth": "1990-01-15",
  "phone_number": "+66812345678",
  "email": "john.doe@example.com"
}
```

**Search Behavior:**
- Empty request `{}` returns all patients in the staff's hospital (up to 100 results)
- Combines criteria with AND logic
- Name fields support partial, case-insensitive matching (ILIKE)
- If no results found and ID provided (national_id or passport_id), attempts to fetch from Hospital A API
- Returns only patients belonging to the staff's hospital (hospital-scoped access)

**Response (200 OK):**
```json
{
  "patients": [
    {
      "id": 1,
      "hospital_id": 1,
      "first_name_th": "สมชาย",
      "middle_name_th": null,
      "last_name_th": "ใจดี",
      "first_name_en": "Somchai",
      "middle_name_en": null,
      "last_name_en": "Jaidee",
      "date_of_birth": "1990-01-15T00:00:00Z",
      "patient_hn": "HN001",
      "national_id": "1234567890123",
      "passport_id": null,
      "phone_number": "+66812345678",
      "email": "somchai@example.com",
      "gender": "M",
      "created_at": "2026-03-08T10:00:00Z",
      "updated_at": "2026-03-08T10:00:00Z"
    }
  ],
  "count": 1
}
```

**Error Responses:**

- **400 Bad Request** - Invalid JSON
```json
{
  "error": "invalid request body"
}
```

- **401 Unauthorized** - Missing or invalid token
```json
{
  "error": "unauthorized"
}
```
OR
```json
{
  "error": "missing or invalid bearer token"
}
```

- **500 Internal Server Error**
```json
{
  "error": "failed to search patients"
}
```

---

## Data Models

### AuthResult
```json
{
  "staff_id": integer,
  "username": string,
  "hospital_id": integer,
  "hospital": string,
  "token": string
}
```

### Patient
```json
{
  "id": integer,
  "hospital_id": integer,
  "first_name_th": string | null,
  "middle_name_th": string | null,
  "last_name_th": string | null,
  "first_name_en": string | null,
  "middle_name_en": string | null,
  "last_name_en": string | null,
  "date_of_birth": string (ISO 8601) | null,
  "patient_hn": string | null,
  "national_id": string | null,
  "passport_id": string | null,
  "phone_number": string | null,
  "email": string | null,
  "gender": string ("M" | "F") | null,
  "created_at": string (ISO 8601),
  "updated_at": string (ISO 8601)
}
```

### PatientSearchResponse
```json
{
  "patients": Patient[],
  "count": integer
}
```

---

## JWT Token

### Token Structure
- **Algorithm:** HS256 (HMAC with SHA-256)
- **Expiration:** 60 minutes (configurable via JWT_TTL_MINUTES)
- **Format:** `eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.payload.signature`

### Claims
```json
{
  "staff_id": integer,
  "hospital_id": integer,
  "username": string,
  "hospital": string,
  "iat": integer (issued at timestamp),
  "exp": integer (expiration timestamp)
}
```

### Usage
Include in Authorization header:
```
Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...
```

---

## External Integration

### Hospital A API

#### GET https://hospital-a.api.co.th/patient/search/{id}
External API for fetching patient data from Hospital A when not found locally.

**Parameters:**
- `id`: national_id or passport_id

**Response:**
```json
{
  "first_name_th": "สมชาย",
  "middle_name_th": null,
  "last_name_th": "ใจดี",
  "first_name_en": "Somchai",
  "middle_name_en": null,
  "last_name_en": "Jaidee",
  "date_of_birth": "1990-01-15",
  "patient_hn": "HN001",
  "national_id": "1234567890123",
  "passport_id": null,
  "phone_number": "+66812345678",
  "email": "somchai@example.com",
  "gender": "M"
}
```

**Notes:**
- Called automatically when patient search returns empty and ID provided
- Fetched data is saved to local database
- Re-search performed after save to return standardized response

---

## Security Features

### Password Security
- Passwords are hashed using bcrypt (cost factor: 10)
- Never stored or transmitted in plain text
- Minimum 8 characters required

### Hospital Isolation
- Staff can only access patients from their own hospital
- Hospital ID injected via JWT claims
- Enforced at service layer

### Authentication Flow
1. Staff creates account or logs in
2. Server validates credentials
3. Server generates JWT with staff/hospital info
4. Client stores token
5. Client includes token in subsequent requests
6. Server validates token and extracts context
7. Server enforces hospital-scoped access

---

## Error Handling

### HTTP Status Codes
- **200 OK**: Successful request
- **201 Created**: Resource created successfully
- **400 Bad Request**: Invalid request format or validation error
- **401 Unauthorized**: Missing or invalid authentication
- **409 Conflict**: Resource already exists
- **500 Internal Server Error**: Server-side error
- **503 Service Unavailable**: Service or dependency unavailable

### Error Response Format
```json
{
  "error": "human-readable error message"
}
```

---

## Rate Limiting & Constraints

### Search Constraints
- Maximum 100 patients returned per search
- All searches are hospital-scoped

### Database Constraints
- Unique: (username, hospital_id) for staff
- Unique (where not null): national_id, passport_id for patients
- Gender must be 'M' or 'F'

---

## Environment Configuration

### Required Variables
```
APP_PORT=8080                    # Application port
DATABASE_URL=postgres://...      # PostgreSQL connection string
JWT_SECRET=your-secret-key       # JWT signing secret
JWT_TTL_MINUTES=60              # Token expiration time
HOSPITAL_A_API_URL=https://...  # Hospital A API endpoint
```

---

## Testing

### Example: Create Staff and Search Patients

```bash
# 1. Create staff account
curl -X POST http://localhost:8088/api/v1/staff/create \
  -H "Content-Type: application/json" \
  -d '{
    "username": "alice",
    "password": "password123",
    "hospital": "HOSPITAL_A"
  }'

# Response: Save the token
# {"staff_id":1,"username":"alice","hospital_id":1,"hospital":"HOSPITAL_A","token":"eyJ..."}

# 2. Search patients
curl -X POST http://localhost:8088/api/v1/patient/search \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer eyJ..." \
  -d '{
    "first_name": "som"
  }'
```

---

## Version History

- **v1 (Current)**: Initial release with staff authentication and patient search
