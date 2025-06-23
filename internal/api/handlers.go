package api

import (
	"errors"
	"log"
	"net/http"
	"time"

	"github.com/axellelanca/urlshortener/internal/models"
	"github.com/axellelanca/urlshortener/internal/services"
	"github.com/gin-gonic/gin"
)

var ClickEventsChannel chan models.ClickEvent

func SetupRoutes(router *gin.Engine, linkService *services.LinkService, bufferSize int, baseURL string) {
	if ClickEventsChannel == nil {
		ClickEventsChannel = make(chan models.ClickEvent, bufferSize)
	}

	router.GET("/health", HealthCheckHandler)

	apiV1 := router.Group("/api/v1")
	{
		apiV1.POST("/links", CreateShortLinkHandler(linkService, baseURL))
		apiV1.GET("/links/:shortCode/stats", GetLinkStatsHandler(linkService))
	}

	router.GET("/:shortCode", RedirectHandler(linkService))
}

func HealthCheckHandler(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}

type CreateLinkRequest struct {
	LongURL string `json:"long_url" binding:"required,url"`
}

func CreateShortLinkHandler(linkService *services.LinkService, baseURL string) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req CreateLinkRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		link, err := linkService.CreateLink(req.LongURL)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusCreated, gin.H{
			"short_code":     link.ShortCode,
			"long_url":       link.LongURL,
			"full_short_url": baseURL + "/" + link.ShortCode,
		})
	}
}

func RedirectHandler(linkService *services.LinkService) gin.HandlerFunc {
	return func(c *gin.Context) {
		shortCode := c.Param("shortCode")

		link, err := linkService.GetLinkByShortCode(shortCode)
		if err != nil {
			if errors.Is(err, models.ErrLinkNotFound) {
				c.JSON(http.StatusNotFound, gin.H{"error": "Link not found"})
				return
			}
			log.Printf("Error retrieving link for %s: %v", shortCode, err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
			return
		}

		clickEvent := models.ClickEvent{
			LinkID:    link.ID,
			Timestamp: time.Now(),
			UserAgent: c.GetHeader("User-Agent"),
			IPAddress: c.ClientIP(),
		}

		select {
		case ClickEventsChannel <- clickEvent:
		default:
			log.Printf("Warning: ClickEventsChannel is full, dropping click event for %s.", shortCode)
		}

		c.Redirect(http.StatusFound, link.LongURL)
	}
}

func GetLinkStatsHandler(linkService *services.LinkService) gin.HandlerFunc {
	return func(c *gin.Context) {
		shortCode := c.Param("shortCode")

		link, totalClicks, err := linkService.GetLinkStats(shortCode)
		if err != nil {
			if errors.Is(err, models.ErrLinkNotFound) {
				c.JSON(http.StatusNotFound, gin.H{"error": "Link not found"})
				return
			}
			log.Printf("Error retrieving stats for %s: %v", shortCode, err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"short_code":   link.ShortCode,
			"long_url":     link.LongURL,
			"total_clicks": totalClicks,
		})
	}
}
