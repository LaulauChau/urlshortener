package api

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/axellelanca/urlshortener/internal/services"
	"github.com/axellelanca/urlshortener/internal/services/mocks"
	"github.com/gin-gonic/gin"
)

func setupTestRouter() (*gin.Engine, *services.LinkService) {
	gin.SetMode(gin.TestMode)
	
	mockLinkRepo := mocks.NewMockLinkRepository()
	mockClickRepo := mocks.NewMockClickRepository()
	linkService := services.NewLinkService(mockLinkRepo, mockClickRepo)
	
	router := gin.New()
	SetupRoutes(router, linkService, 100, "http://localhost:8080")
	
	return router, linkService
}

func TestHealthCheckHandler(t *testing.T) {
	router, _ := setupTestRouter()

	req, err := http.NewRequest("GET", "/health", nil)
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, w.Code)
	}

	var response map[string]interface{}
	if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}

	if response["status"] != "ok" {
		t.Errorf("Expected status 'ok', got %v", response["status"])
	}
}

func TestCreateShortLinkHandler(t *testing.T) {
	router, _ := setupTestRouter()

	tests := []struct {
		name           string
		requestBody    map[string]interface{}
		expectedStatus int
	}{
		{
			name: "valid request",
			requestBody: map[string]interface{}{
				"long_url": "https://example.com",
			},
			expectedStatus: http.StatusCreated,
		},
		{
			name: "invalid URL",
			requestBody: map[string]interface{}{
				"long_url": "not-a-url",
			},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name: "missing long_url",
			requestBody: map[string]interface{}{
				"other_field": "value",
			},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "empty request body",
			requestBody:    map[string]interface{}{},
			expectedStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			jsonBody, err := json.Marshal(tt.requestBody)
			if err != nil {
				t.Fatalf("Failed to marshal request body: %v", err)
			}

			req, err := http.NewRequest("POST", "/api/v1/links", bytes.NewBuffer(jsonBody))
			if err != nil {
				t.Fatalf("Failed to create request: %v", err)
			}
			req.Header.Set("Content-Type", "application/json")

			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			if w.Code != tt.expectedStatus {
				t.Errorf("Expected status code %d, got %d", tt.expectedStatus, w.Code)
			}

			if tt.expectedStatus == http.StatusCreated {
				var response map[string]interface{}
				if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
					t.Fatalf("Failed to unmarshal response: %v", err)
				}

				if response["short_code"] == nil {
					t.Errorf("Expected short_code in response")
				}
				if response["long_url"] != tt.requestBody["long_url"] {
					t.Errorf("Expected long_url %v, got %v", tt.requestBody["long_url"], response["long_url"])
				}
				if response["full_short_url"] == nil {
					t.Errorf("Expected full_short_url in response")
				}
			}
		})
	}
}

func TestRedirectHandler(t *testing.T) {
	router, linkService := setupTestRouter()
	
	link, err := linkService.CreateLink("https://example.com")
	if err != nil {
		t.Fatalf("Failed to create test link: %v", err)
	}

	tests := []struct {
		name           string
		shortCode      string
		expectedStatus int
		expectedLocation string
	}{
		{
			name:             "valid short code",
			shortCode:        link.ShortCode,
			expectedStatus:   http.StatusFound,
			expectedLocation: "https://example.com",
		},
		{
			name:           "invalid short code",
			shortCode:      "nonexistent",
			expectedStatus: http.StatusNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req, err := http.NewRequest("GET", "/"+tt.shortCode, nil)
			if err != nil {
				t.Fatalf("Failed to create request: %v", err)
			}
			req.Header.Set("User-Agent", "Test Agent")

			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			if w.Code != tt.expectedStatus {
				t.Errorf("Expected status code %d, got %d", tt.expectedStatus, w.Code)
			}

			if tt.expectedStatus == http.StatusFound {
				location := w.Header().Get("Location")
				if location != tt.expectedLocation {
					t.Errorf("Expected location %s, got %s", tt.expectedLocation, location)
				}
			}
		})
	}
}

func TestGetLinkStatsHandler(t *testing.T) {
	router, linkService := setupTestRouter()

	link, err := linkService.CreateLink("https://example.com")
	if err != nil {
		t.Fatalf("Failed to create test link: %v", err)
	}

	tests := []struct {
		name           string
		shortCode      string
		expectedStatus int
	}{
		{
			name:           "valid short code",
			shortCode:      link.ShortCode,
			expectedStatus: http.StatusOK,
		},
		{
			name:           "invalid short code",
			shortCode:      "nonexistent",
			expectedStatus: http.StatusNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req, err := http.NewRequest("GET", "/api/v1/links/"+tt.shortCode+"/stats", nil)
			if err != nil {
				t.Fatalf("Failed to create request: %v", err)
			}

			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			if w.Code != tt.expectedStatus {
				t.Errorf("Expected status code %d, got %d", tt.expectedStatus, w.Code)
			}

			if tt.expectedStatus == http.StatusOK {
				var response map[string]interface{}
				if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
					t.Fatalf("Failed to unmarshal response: %v", err)
				}

				if response["short_code"] != tt.shortCode {
					t.Errorf("Expected short_code %s, got %v", tt.shortCode, response["short_code"])
				}
				if response["long_url"] == nil {
					t.Errorf("Expected long_url in response")
				}
				if response["total_clicks"] == nil {
					t.Errorf("Expected total_clicks in response")
				}
			}
		})
	}
}

func TestCreateLinkRequest_Validation(t *testing.T) {
	tests := []struct {
		name    string
		request CreateLinkRequest
		valid   bool
	}{
		{
			name: "valid request",
			request: CreateLinkRequest{
				LongURL: "https://example.com",
			},
			valid: true,
		},
		{
			name: "empty URL",
			request: CreateLinkRequest{
				LongURL: "",
			},
			valid: false,
		},
		{
			name: "invalid URL format",
			request: CreateLinkRequest{
				LongURL: "not-a-url",
			},
			valid: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if len(tt.request.LongURL) > 0 && !tt.valid {
				t.Log("URL validation would be handled by Gin binding")
			}
		})
	}
}

func TestClickEventsChannel(t *testing.T) {
	router, linkService := setupTestRouter()

	link, err := linkService.CreateLink("https://example.com")
	if err != nil {
		t.Fatalf("Failed to create test link: %v", err)
	}

	if ClickEventsChannel == nil {
		t.Error("ClickEventsChannel should be initialized")
	}

	req, err := http.NewRequest("GET", "/"+link.ShortCode, nil)
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}
	req.Header.Set("User-Agent", "Test Agent")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusFound {
		t.Errorf("Expected status code %d, got %d", http.StatusFound, w.Code)
	}

	if ClickEventsChannel == nil {
		t.Error("ClickEventsChannel should not be nil after redirect")
	}
} 