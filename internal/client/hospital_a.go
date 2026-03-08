package client

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type HospitalAPatient struct {
	FirstNameTH  string `json:"first_name_th"`
	MiddleNameTH string `json:"middle_name_th"`
	LastNameTH   string `json:"last_name_th"`
	FirstNameEN  string `json:"first_name_en"`
	MiddleNameEN string `json:"middle_name_en"`
	LastNameEN   string `json:"last_name_en"`
	DateOfBirth  string `json:"date_of_birth"`
	PatientHN    string `json:"patient_hn"`
	NationalID   string `json:"national_id"`
	PassportID   string `json:"passport_id"`
	PhoneNumber  string `json:"phone_number"`
	Email        string `json:"email"`
	Gender       string `json:"gender"`
}

type HospitalAClient struct {
	baseURL    string
	httpClient *http.Client
}

func NewHospitalAClient(baseURL string) *HospitalAClient {
	return &HospitalAClient{
		baseURL: baseURL,
		httpClient: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

func (c *HospitalAClient) SearchPatient(ctx context.Context, idValue string) (*HospitalAPatient, error) {
	url := fmt.Sprintf("%s/patient/search/%s", c.baseURL, idValue)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return nil, nil
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("hospital A API returned status %d", resp.StatusCode)
	}

	var patient HospitalAPatient
	if err := json.NewDecoder(resp.Body).Decode(&patient); err != nil {
		return nil, err
	}

	return &patient, nil
}
