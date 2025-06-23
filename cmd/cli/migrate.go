package cli

import (
	"fmt"
	"log"

	"github.com/axellelanca/urlshortener/cmd"
	"github.com/axellelanca/urlshortener/internal/models"
	"github.com/spf13/cobra"
	"gorm.io/driver/sqlite" // Driver SQLite pour GORM
	"gorm.io/gorm"
)

var MigrateCmd = &cobra.Command{
	Use:   "migrate",
	Short: "Exécute les migrations de la base de données pour créer ou mettre à jour les tables.",
	Long: `Cette commande se connecte à la base de données configurée (SQLite)
et exécute les migrations automatiques de GORM pour créer les tables 'links' et 'clicks'
basées sur les modèles Go.`,
	Run: func(cobraCmd *cobra.Command, args []string) {
		cfg := cmd.Cfg
		if cfg == nil {
			log.Fatalf("FATAL: Configuration non chargée")
		}

		db, err := gorm.Open(sqlite.Open(cfg.Database.Name), &gorm.Config{})
		if err != nil {
			log.Fatalf("FATAL: Échec de la connexion à la base de données: %v", err)
		}

		sqlDB, err := db.DB()
		if err != nil {
			log.Fatalf("FATAL: Échec de l'obtention de la base de données SQL sous-jacente: %v", err)
		}
		defer sqlDB.Close()

		if err := db.AutoMigrate(&models.Link{}, &models.Click{}); err != nil {
			log.Fatalf("FATAL: Échec des migrations: %v", err)
		}

		fmt.Println("Migrations de la base de données exécutées avec succès.")
	},
}

func init() {
	cmd.RootCmd.AddCommand(MigrateCmd)
}
