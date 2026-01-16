package diagnocat

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

const BaseURL = "https://api.diagnocat.com"

type Client struct {
	APIKey     string
	HTTPClient *http.Client
}

func NewClient(apiKey string) *Client {
	return &Client{
		APIKey: apiKey,
		HTTPClient: &http.Client{
			Timeout: 60 * time.Second,
		},
	}
}

// Request structures
type OpenSessionRequest struct {
	StudyUID   string `json:"study_uid"`
	PatientUID string `json:"patient_uid,omitempty"`
}

type OpenSessionResponse struct {
	OK        bool   `json:"ok"`
	SessionID string `json:"session_id"`
	Hostname  string `json:"hostname"`
	Error     string `json:"error,omitempty"`
}

type RequestUploadURLsRequest struct {
	SessionID string   `json:"session_id"`
	Keys      []string `json:"keys"`
}

type UploadURL struct {
	Key string `json:"key"`
	URL string `json:"url"`
}

type RequestUploadURLsResponse struct {
	OK         bool        `json:"ok"`
	UploadURLs []UploadURL `json:"upload_urls"`
	Error      string      `json:"error,omitempty"`
}

type CloseSessionRequest struct {
	SessionID string `json:"session_id"`
}

type CloseSessionResponse struct {
	OK    bool   `json:"ok"`
	Error string `json:"error,omitempty"`
}

type Analysis struct {
	UID          string    `json:"uid"`
	StudyUID     string    `json:"study_uid"`
	PatientUID   string    `json:"patient_uid"`
	AnalysisType string    `json:"analysis_type"`
	Complete     bool      `json:"complete"`
	Started      bool      `json:"started"`
	Error        string    `json:"error"`
	PDFUrl       string    `json:"pdf_url"`
	PreviewUrl   string    `json:"preview_url"`
	WebpageUrl   string    `json:"webpage_url"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

// OpenSession creates a new upload session
func (c *Client) OpenSession(studyUID, patientUID string) (*OpenSessionResponse, error) {
	reqBody := OpenSessionRequest{
		StudyUID:   studyUID,
		PatientUID: patientUID,
	}

	body, err := json.Marshal(reqBody)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", BaseURL+"/v1/upload/open-session", bytes.NewBuffer(body))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", c.APIKey)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var result OpenSessionResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	if !result.OK {
		return nil, fmt.Errorf("diagnocat error: %s", result.Error)
	}

	return &result, nil
}

// RequestUploadURLs gets presigned S3 URLs for uploading files
func (c *Client) RequestUploadURLs(sessionID string, fileKeys []string) (*RequestUploadURLsResponse, error) {
	reqBody := RequestUploadURLsRequest{
		SessionID: sessionID,
		Keys:      fileKeys,
	}

	body, err := json.Marshal(reqBody)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", BaseURL+"/v1/upload/request-upload-urls", bytes.NewBuffer(body))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", c.APIKey)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var result RequestUploadURLsResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	if !result.OK {
		return nil, fmt.Errorf("diagnocat error: %s", result.Error)
	}

	return &result, nil
}

// UploadFileToS3 uploads file to presigned S3 URL
func (c *Client) UploadFileToS3(uploadURL string, fileContent []byte) error {
	req, err := http.NewRequest("PUT", uploadURL, bytes.NewReader(fileContent))
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/dicom")
	req.ContentLength = int64(len(fileContent))

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("S3 upload failed: %s", string(body))
	}

	return nil
}

// CloseSession closes the upload session and triggers analysis
func (c *Client) CloseSession(sessionID string) error {
	reqBody := CloseSessionRequest{
		SessionID: sessionID,
	}

	body, err := json.Marshal(reqBody)
	if err != nil {
		return err
	}

	req, err := http.NewRequest("POST", BaseURL+"/v1/upload/start-session-close", bytes.NewBuffer(body))
	if err != nil {
		return err
	}

	req.Header.Set("Authorization", c.APIKey)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	var result CloseSessionResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return err
	}

	if !result.OK {
		return fmt.Errorf("diagnocat error: %s", result.Error)
	}

	return nil
}

// GetAnalyses lists all analyses for a patient
func (c *Client) GetAnalyses(patientUID string) ([]Analysis, error) {
	url := fmt.Sprintf("%s/v2/analyses?patient_uid=%s", BaseURL, patientUID)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", c.APIKey)
	req.Header.Set("Accept", "application/json")

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var analyses []Analysis
	if err := json.NewDecoder(resp.Body).Decode(&analyses); err != nil {
		return nil, err
	}

	return analyses, nil
}

// GetAnalysis gets a specific analysis by ID
func (c *Client) GetAnalysis(analysisUID string) (*Analysis, error) {
	url := fmt.Sprintf("%s/v2/analyses/%s", BaseURL, analysisUID)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", c.APIKey)
	req.Header.Set("Accept", "application/json")

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var analysis Analysis
	if err := json.NewDecoder(resp.Body).Decode(&analysis); err != nil {
		return nil, err
	}

	return &analysis, nil
}
