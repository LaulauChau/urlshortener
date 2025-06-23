package models

import "errors"

var (
	ErrLinkNotFound = errors.New("link not found")
	
	ErrInvalidURL = errors.New("invalid URL format")
	
	ErrDuplicateShortCode = errors.New("short code already exists")

	ErrShortCodeGenerationFailed = errors.New("failed to generate unique short code after maximum retries")	

	ErrDatabaseConnection = errors.New("database connection error")
	
	ErrConfigurationLoad = errors.New("failed to load configuration")
) 