package services

import (
	"crypto/rand"
	"errors"
	"fmt"
	"log"
	"math/big"
	"time"

	"github.com/axellelanca/urlshortener/internal/models"
	"github.com/axellelanca/urlshortener/internal/repository"
	"gorm.io/gorm"
)


const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"


type LinkService struct {
	linkRepo  repository.LinkRepository
	clickRepo repository.ClickRepository
}

func NewLinkService(linkRepo repository.LinkRepository, clickRepo repository.ClickRepository) *LinkService {
	return &LinkService{
		linkRepo:  linkRepo,
		clickRepo: clickRepo,
	}
}

func (s *LinkService) GenerateShortCode(length int) (string, error) {
	result := make([]byte, length)
	for i := range result {
		num, err := rand.Int(rand.Reader, big.NewInt(int64(len(charset))))
		if err != nil {
			return "", fmt.Errorf("failed to generate random number: %w", err)
		}
		result[i] = charset[num.Int64()]
	}
	return string(result), nil
}

func (s *LinkService) CreateLink(longURL string) (*models.Link, error) {
	var shortCode string
	const maxRetries = 5

	for i := 0; i < maxRetries; i++ {
		code, err := s.GenerateShortCode(6)
		if err != nil {
			return nil, fmt.Errorf("failed to generate short code: %w", err)
		}
		_, err = s.linkRepo.GetLinkByShortCode(code)

		if err != nil {
	
			if errors.Is(err, gorm.ErrRecordNotFound) {
				shortCode = code 
				break            
			}

			return nil, fmt.Errorf("database error checking short code uniqueness: %w", err)
		}
		log.Printf("Short code '%s' already exists, retrying generation (%d/%d)...", code, i+1, maxRetries)
	}

	if shortCode == "" {
		return nil, models.ErrShortCodeGenerationFailed
	}

	link := &models.Link{
		ShortCode: shortCode,
		LongURL:   longURL,
		CreatedAt: time.Now(),
	}

	if err := s.linkRepo.CreateLink(link); err != nil {
		return nil, fmt.Errorf("failed to create link in database: %w", err)
	}

	return link, nil
}


func (s *LinkService) GetLinkByShortCode(shortCode string) (*models.Link, error) {
	link, err := s.linkRepo.GetLinkByShortCode(shortCode)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, models.ErrLinkNotFound
		}
		return nil, fmt.Errorf("database error retrieving link: %w", err)
	}

	return link, nil
}


func (s *LinkService) GetLinkStats(shortCode string) (*models.Link, int, error) {
	link, err := s.GetLinkByShortCode(shortCode)
	if err != nil {
		return nil, 0, err
	}

	totalClicks, err := s.clickRepo.CountClicksByLinkID(link.ID)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count clicks: %w", err)
	}

	return link, totalClicks, nil
}

