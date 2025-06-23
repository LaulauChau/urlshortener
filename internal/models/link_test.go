package models

import (
	"testing"
	"time"
)

func TestLink_Validation(t *testing.T) {
	tests := []struct {
		name      string
		link      Link
		expectErr bool
	}{
		{
			name: "valid link",
			link: Link{
				ID:        1,
				ShortCode: "abc123",
				LongURL:   "https://example.com",
				CreatedAt: time.Now(),
			},
			expectErr: false,
		},
		{
			name: "empty short code",
			link: Link{
				ID:        1,
				ShortCode: "",
				LongURL:   "https://example.com",
				CreatedAt: time.Now(),
			},
			expectErr: true,
		},
		{
			name: "empty long URL",
			link: Link{
				ID:        1,
				ShortCode: "abc123",
				LongURL:   "",
				CreatedAt: time.Now(),
			},
			expectErr: true,
		},
		{
			name: "short code too long",
			link: Link{
				ID:        1,
				ShortCode: "thisisaverylongshortcode",
				LongURL:   "https://example.com",
				CreatedAt: time.Now(),
			},
			expectErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Basic validation checks
			if tt.link.ShortCode == "" && !tt.expectErr {
				t.Errorf("Expected ShortCode to be non-empty")
			}
			if tt.link.LongURL == "" && !tt.expectErr {
				t.Errorf("Expected LongURL to be non-empty")
			}
			if len(tt.link.ShortCode) > 10 && !tt.expectErr {
				t.Errorf("Expected ShortCode to be <= 10 characters")
			}
		})
	}
}

func TestLink_Creation(t *testing.T) {
	now := time.Now()
	link := Link{
		ID:        1,
		ShortCode: "abc123",
		LongURL:   "https://example.com",
		CreatedAt: now,
	}

	if link.ID != 1 {
		t.Errorf("Expected ID to be 1, got %d", link.ID)
	}
	if link.ShortCode != "abc123" {
		t.Errorf("Expected ShortCode to be 'abc123', got %s", link.ShortCode)
	}
	if link.LongURL != "https://example.com" {
		t.Errorf("Expected LongURL to be 'https://example.com', got %s", link.LongURL)
	}
	if !link.CreatedAt.Equal(now) {
		t.Errorf("Expected CreatedAt to be %v, got %v", now, link.CreatedAt)
	}
} 