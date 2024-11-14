package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/xanagit/kotoquiz-api/config"
	"github.com/xanagit/kotoquiz-api/controllers"
	"github.com/xanagit/kotoquiz-api/models"
	"github.com/xanagit/kotoquiz-api/repositories"
	"github.com/xanagit/kotoquiz-api/services"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
)

func main() {
	// Load configuration
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// Use the loaded configuration values
	dbConfig := cfg.Database
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d",
		dbConfig.Host, dbConfig.User, dbConfig.Password, dbConfig.Name, dbConfig.Port)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	// Auto-migrate the models to keep the schema in sync
	err = db.AutoMigrate(&models.Label{}, &models.Word{})
	if err != nil {
		log.Fatal("failed to migrate database:", err)
	}

	// Initialisation du repository et du service
	wordRepository := &repositories.WordRepositoryImpl{DB: db}
	wordService := &services.WordServiceImpl{Repo: wordRepository}
	wordController := &controllers.WordControllerImpl{Service: wordService}
	// Configuration de l'application Gin
	r := gin.Default()
	apiGroup := r.Group("/api/v1/words")
	{
		// Utilisez les méthodes du service
		apiGroup.GET("", wordController.GetWords)
		apiGroup.GET("/:id", wordController.GetWordByID)
		apiGroup.POST("", wordController.CreateWord)
		apiGroup.PUT("/:id", wordController.UpdateWord)
		apiGroup.DELETE("/:id", wordController.DeleteWord)
	}

	runError := r.Run()
	if runError != nil {
		return
	}
}
