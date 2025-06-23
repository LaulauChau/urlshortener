package models

import (
	"testing"
	"time"
)

func TestClick_Creation(t *testing.T) {
	now := time.Now()
	link := Link{
		ID:        1,
		ShortCode: "abc123",
		LongURL:   "https://example.com",
		CreatedAt: now,
	}

	click := Click{
		ID:        1,
		LinkID:    1,
		Link:      link,
		Timestamp: now,
		UserAgent: "Mozilla/5.0 (Test Browser)",
		IPAddress: "127.0.0.1",
	}

	if click.ID != 1 {
		t.Errorf("Expected Click ID to be 1, got %d", click.ID)
	}
	if click.LinkID != 1 {
		t.Errorf("Expected LinkID to be 1, got %d", click.LinkID)
	}
	if click.Link.ID != link.ID {
		t.Errorf("Expected Link ID to be %d, got %d", link.ID, click.Link.ID)
	}
	if !click.Timestamp.Equal(now) {
		t.Errorf("Expected Timestamp to be %v, got %v", now, click.Timestamp)
	}
	if click.UserAgent != "Mozilla/5.0 (Test Browser)" {
		t.Errorf("Expected UserAgent to be 'Mozilla/5.0 (Test Browser)', got %s", click.UserAgent)
	}
	if click.IPAddress != "127.0.0.1" {
		t.Errorf("Expected IPAddress to be '127.0.0.1', got %s", click.IPAddress)
	}
}

func TestClickEvent_Creation(t *testing.T) {
	now := time.Now()
	event := ClickEvent{
		LinkID:    1,
		Timestamp: now,
		UserAgent: "Mozilla/5.0 (Test Browser)",
		IPAddress: "192.168.1.1",
	}

	if event.LinkID != 1 {
		t.Errorf("Expected LinkID to be 1, got %d", event.LinkID)
	}
	if !event.Timestamp.Equal(now) {
		t.Errorf("Expected Timestamp to be %v, got %v", now, event.Timestamp)
	}
	if event.UserAgent != "Mozilla/5.0 (Test Browser)" {
		t.Errorf("Expected UserAgent to be 'Mozilla/5.0 (Test Browser)', got %s", event.UserAgent)
	}
	if event.IPAddress != "192.168.1.1" {
		t.Errorf("Expected IPAddress to be '192.168.1.1', got %s", event.IPAddress)
	}
}

func TestClick_Validation(t *testing.T) {
	tests := []struct {
		name   string
		click  Click
		valid  bool
	}{
		{
			name: "valid click",
			click: Click{
				ID:        1,
				LinkID:    1,
				Timestamp: time.Now(),
				UserAgent: "Mozilla/5.0",
				IPAddress: "127.0.0.1",
			},
			valid: true,
		},
		{
			name: "missing LinkID",
			click: Click{
				ID:        1,
				LinkID:    0,
				Timestamp: time.Now(),
				UserAgent: "Mozilla/5.0",
				IPAddress: "127.0.0.1",
			},
			valid: false,
		},
		{
			name: "empty IP address",
			click: Click{
				ID:        1,
				LinkID:    1,
				Timestamp: time.Now(),
				UserAgent: "Mozilla/5.0",
				IPAddress: "",
			},
			valid: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Basic validation checks
			hasValidLinkID := tt.click.LinkID > 0
			hasValidIP := tt.click.IPAddress != ""
			
			isValid := hasValidLinkID && hasValidIP
			
			if isValid != tt.valid {
				t.Errorf("Expected validity to be %v, got %v", tt.valid, isValid)
			}
		})
	}
} 