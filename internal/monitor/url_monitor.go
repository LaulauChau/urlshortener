package monitor

import (
	"log"
	"net/http"
	"sync"
	"time"

	_ "github.com/axellelanca/urlshortener/internal/models"
	"github.com/axellelanca/urlshortener/internal/repository"
)

type UrlMonitor struct {
	linkRepo    repository.LinkRepository
	interval    time.Duration
	knownStates map[uint]bool
	mu          sync.Mutex
}

func NewUrlMonitor(linkRepo repository.LinkRepository, interval time.Duration) *UrlMonitor {
	return &UrlMonitor{
		linkRepo:    linkRepo,
		interval:    interval,
		knownStates: make(map[uint]bool),
	}
}

func (m *UrlMonitor) Start() {
	log.Printf("[MONITOR] Démarrage du moniteur d'URLs avec un intervalle de %v...", m.interval)
	ticker := time.NewTicker(m.interval)
	defer ticker.Stop()

	m.checkUrls()

	for range ticker.C {
		m.checkUrls()
	}
}

func (m *UrlMonitor) checkUrls() {
	log.Println("[MONITOR] Lancement de la vérification de l'état des URLs...")

	links, err := m.linkRepo.GetAllLinks()
	if err != nil {
		log.Printf("[MONITOR] ERREUR lors de la récupération des liens pour la surveillance : %v", err)
		return
	}

	for _, link := range links {
		currentState := m.isUrlAccessible(link.LongURL)

		m.mu.Lock()
		previousState, exists := m.knownStates[link.ID]
		m.knownStates[link.ID] = currentState
		m.mu.Unlock()

		if !exists {
			log.Printf("[MONITOR] État initial pour le lien %s (%s) : %s",
				link.ShortCode, link.LongURL, formatState(currentState))
			continue
		}

		if previousState != currentState {
			log.Printf("[NOTIFICATION] Le lien %s (%s) est passé de %s à %s !",
				link.ShortCode, link.LongURL, formatState(previousState), formatState(currentState))
		}

	}
	log.Println("[MONITOR] Vérification de l'état des URLs terminée.")
}

func (m *UrlMonitor) isUrlAccessible(url string) bool {
	client := &http.Client{
		Timeout: 5 * time.Second,
	}

	resp, err := client.Head(url)
	if err != nil {
		log.Printf("[MONITOR] Erreur d'accès à l'URL '%s': %v", url, err)
		return false
	}

	defer resp.Body.Close()

	return resp.StatusCode >= 200 && resp.StatusCode < 400
}

func formatState(accessible bool) string {
	if accessible {
		return "ACCESSIBLE"
	}
	return "INACCESSIBLE"
}
