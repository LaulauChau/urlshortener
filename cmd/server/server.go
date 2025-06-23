package server

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/axellelanca/urlshortener/cmd"
	"github.com/axellelanca/urlshortener/internal/api"
	"github.com/axellelanca/urlshortener/internal/models"
	"github.com/axellelanca/urlshortener/internal/monitor"
	"github.com/axellelanca/urlshortener/internal/repository"
	"github.com/axellelanca/urlshortener/internal/services"
	"github.com/axellelanca/urlshortener/internal/workers"
	"github.com/gin-gonic/gin"
	"github.com/spf13/cobra"
	"gorm.io/driver/sqlite" 
	"gorm.io/gorm"
)


var RunServerCmd = &cobra.Command{
	Use:   "run-server",
	Short: "Lance le serveur API de raccourcissement d'URLs et les processus de fond.",
	Long: `Cette commande initialise la base de données, configure les APIs,
démarre les workers asynchrones pour les clics et le moniteur d'URLs,
puis lance le serveur HTTP.`,
	Run: func(cobraCmd *cobra.Command, args []string) {

		cfg := cmd.Cfg
		if cfg == nil {
			log.Fatalf("FATAL: Configuration non chargée")
		}


		db, err := gorm.Open(sqlite.Open(cfg.Database.Name), &gorm.Config{})
		if err != nil {
			log.Fatalf("FATAL: Échec de la connexion à la base de données: %v", err)
		}

		linkRepo := repository.NewLinkRepository(db)
		clickRepo := repository.NewClickRepository(db)

		log.Println("Repositories initialisés.")

		linkService := services.NewLinkService(linkRepo, clickRepo)
		_ = services.NewClickService(clickRepo)

	
		log.Println("Services métiers initialisés.")


		clickEventsChannel := make(chan models.ClickEvent, cfg.Analytics.BufferSize)
		workers.StartClickWorkers(cfg.Analytics.WorkerCount, clickEventsChannel, clickRepo)

		log.Printf("Channel d'événements de clic initialisé avec un buffer de %d. %d worker(s) de clics démarré(s).",
			cfg.Analytics.BufferSize, cfg.Analytics.WorkerCount)

		monitorInterval := time.Duration(cfg.Monitor.IntervalMinutes) * time.Minute
		urlMonitor := monitor.NewUrlMonitor(linkRepo, monitorInterval)
		go urlMonitor.Start()
		log.Printf("Moniteur d'URLs démarré avec un intervalle de %v.", monitorInterval)


		router := gin.Default()
		api.SetupRoutes(router, linkService, cfg.Analytics.BufferSize, cfg.Server.BaseURL)


		api.ClickEventsChannel = clickEventsChannel


		log.Println("Routes API configurées.")

	
		serverAddr := fmt.Sprintf(":%d", cfg.Server.Port)
		srv := &http.Server{
			Addr:    serverAddr,
			Handler: router,
		}

		go func() {
			log.Printf("Serveur HTTP démarré sur le port %d", cfg.Server.Port)
			log.Printf("URL de base: %s", cfg.Server.BaseURL)
			if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
				log.Fatalf("FATAL: Échec du démarrage du serveur HTTP: %v", err)
			}
		}()
	

		quit := make(chan os.Signal, 1)
		signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM) 

		<-quit
		log.Println("Signal d'arrêt reçu. Arrêt du serveur...")
		
		log.Println("Arrêt en cours... Donnez un peu de temps aux workers pour finir.")
		time.Sleep(5 * time.Second)

		log.Println("Serveur arrêté proprement.")
	},
}

func init() {
	cmd.RootCmd.AddCommand(RunServerCmd)
}
