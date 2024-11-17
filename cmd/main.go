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
	err = db.AutoMigrate(&models.Label{}, &models.Word{}, &models.Level{})
	if err != nil {
		log.Fatal("failed to migrate database:", err)
	}

	// Word repository, service and controller
	wordRepository := &repositories.WordRepositoryImpl{DB: db}
	wordService := &services.WordServiceImpl{Repo: wordRepository}
	wordController := &controllers.WordControllerImpl{Service: wordService}

	// WordDto service and controller
	wordDtoService := &services.WordDtoServiceImpl{Repo: wordRepository}
	wordDtoController := &controllers.WordDtoControllerImpl{WordDtoService: wordDtoService}

	// Label repository, service and controller
	labelRepository := &repositories.LabelRepositoryImpl{DB: db}
	labelService := &services.LabelServiceImpl{Repo: labelRepository}
	labelController := &controllers.LabelControllerImpl{Service: labelService}

	// Configuration de l'application Gin
	r := gin.Default()
	appUserGroup := r.Group("/api/v1/app")
	{
		appUserGroup.GET("/words", wordDtoController.ListDtoWords)    // query param: ids, lang
		appUserGroup.GET("/words/:id", wordDtoController.ReadDtoWord) // query param: lang
	}

	labelsGroup := r.Group("/api/v1")
	{
		labelsGroup.GET("/tags", labelController.ListTags)
		labelsGroup.GET("/categories", labelController.ListCategories)
		labelsGroup.GET("/levelNames", labelController.ListLevelNames)
	}

	techGroup := r.Group("/api/v1/tech")
	{
		techGroup.GET("/words/:id", wordController.ReadWord)
		techGroup.POST("/words", wordController.CreateWord)
		techGroup.PUT("/words/:id", wordController.UpdateWord)
		techGroup.DELETE("/words/:id", wordController.DeleteWord)
	}

	runError := r.Run()
	if runError != nil {
		return
	}
}
