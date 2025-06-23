package mocks

import (
	"errors"

	"github.com/axellelanca/urlshortener/internal/models"
	"gorm.io/gorm"
)

type MockLinkRepository struct {
	links   map[string]*models.Link
	nextID  uint
	shouldFail bool
}

func NewMockLinkRepository() *MockLinkRepository {
	return &MockLinkRepository{
		links:  make(map[string]*models.Link),
		nextID: 1,
	}
}

func (m *MockLinkRepository) SetShouldFail(shouldFail bool) {
	m.shouldFail = shouldFail
}

func (m *MockLinkRepository) CreateLink(link *models.Link) error {
	if m.shouldFail {
		return errors.New("mock database error")
	}
	
	if link.ID == 0 {
		link.ID = m.nextID
		m.nextID++
	}
	
	m.links[link.ShortCode] = link
	return nil
}

func (m *MockLinkRepository) GetLinkByShortCode(shortCode string) (*models.Link, error) {
	if m.shouldFail {
		return nil, errors.New("mock database error")
	}
	
	link, exists := m.links[shortCode]
	if !exists {
		return nil, gorm.ErrRecordNotFound
	}
	
	return link, nil
}

func (m *MockLinkRepository) GetAllLinks() ([]models.Link, error) {
	if m.shouldFail {
		return nil, errors.New("mock database error")
	}
	
	var allLinks []models.Link
	for _, link := range m.links {
		allLinks = append(allLinks, *link)
	}
	
	return allLinks, nil
}

func (m *MockLinkRepository) CountClicksByLinkID(linkID uint) (int, error) {
	if m.shouldFail {
		return 0, errors.New("mock database error")
	}
	
	return 5, nil
}

type MockClickRepository struct {
	clicks     map[uint][]models.Click
	shouldFail bool
}

func NewMockClickRepository() *MockClickRepository {
	return &MockClickRepository{
		clicks: make(map[uint][]models.Click),
	}
}

func (m *MockClickRepository) SetShouldFail(shouldFail bool) {
	m.shouldFail = shouldFail
}

func (m *MockClickRepository) CreateClick(click *models.Click) error {
	if m.shouldFail {
		return errors.New("mock database error")
	}
	
	m.clicks[click.LinkID] = append(m.clicks[click.LinkID], *click)
	return nil
}

func (m *MockClickRepository) GetClicksByLinkID(linkID uint) ([]models.Click, error) {
	if m.shouldFail {
		return nil, errors.New("mock database error")
	}
	
	clicks, exists := m.clicks[linkID]
	if !exists {
		return []models.Click{}, nil
	}
	
	return clicks, nil
}

func (m *MockClickRepository) CountClicksByLinkID(linkID uint) (int, error) {
	if m.shouldFail {
		return 0, errors.New("mock database error")
	}
	
	clicks, exists := m.clicks[linkID]
	if !exists {
		return 0, nil
	}
	
	return len(clicks), nil
} 