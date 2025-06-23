package services

import (
	"fmt"

	"github.com/axellelanca/urlshortener/internal/models"
	"github.com/axellelanca/urlshortener/internal/repository"
)

// ClickService fournit des méthodes pour la logique métier des clics.
type ClickService struct {
	clickRepo repository.ClickRepository
}

// NewClickService crée et retourne une nouvelle instance de ClickService.
// C'est la fonction recommandée pour obtenir un service, assurant que toutes ses dépendances sont injectées.
func NewClickService(clickRepo repository.ClickRepository) *ClickService {
	return &ClickService{
		clickRepo: clickRepo,
	}
}

// ProcessClickEvent traite un événement de clic et le persiste en base de données.
func (s *ClickService) ProcessClickEvent(event models.ClickEvent) error {
	click := &models.Click{
		LinkID:    event.LinkID,
		Timestamp: event.Timestamp,
		UserAgent: event.UserAgent,
		IPAddress: event.IPAddress,
	}

	if err := s.clickRepo.CreateClick(click); err != nil {
		return fmt.Errorf("failed to create click: %w", err)
	}

	return nil
}

// GetClicksByLinkID récupère tous les clics pour un lien donné.
func (s *ClickService) GetClicksByLinkID(linkID uint) ([]models.Click, error) {
	clicks, err := s.clickRepo.GetClicksByLinkID(linkID)
	if err != nil {
		return nil, fmt.Errorf("failed to get clicks: %w", err)
	}
	return clicks, nil
}

// CountClicksByLinkID compte le nombre de clics pour un lien donné.
func (s *ClickService) CountClicksByLinkID(linkID uint) (int, error) {
	count, err := s.clickRepo.CountClicksByLinkID(linkID)
	if err != nil {
		return 0, fmt.Errorf("failed to count clicks: %w", err)
	}
	return count, nil
}
