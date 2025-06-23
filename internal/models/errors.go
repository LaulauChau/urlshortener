package models

import "errors"

// Custom errors for the URL shortener application
var (
	// ErrLinkNotFound is returned when a short link is not found
	ErrLinkNotFound = errors.New("link not found")
	
	// ErrInvalidURL is returned when the provided URL is invalid
	ErrInvalidURL = errors.New("invalid URL format")
	
	// ErrDuplicateShortCode is returned when trying to create a link with an existing short code
	ErrDuplicateShortCode = errors.New("short code already exists")
	
	// ErrShortCodeGenerationFailed is returned when unable to generate a unique short code after retries
	ErrShortCodeGenerationFailed = errors.New("failed to generate unique short code after maximum retries")
	
	// ErrDatabaseConnection is returned when there's a database connection issue
	ErrDatabaseConnection = errors.New("database connection error")
	
	// ErrConfigurationLoad is returned when configuration cannot be loaded
	ErrConfigurationLoad = errors.New("failed to load configuration")
) 