package initialisation

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

func GinHandlers(db *gorm.DB) *gin.Engine {
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

	// Tag controller
	tagController := &controllers.TagControllerImpl{Service: labelService}

	// Category controller
	categoryController := &controllers.CategoryControllerImpl{Service: labelService}

	// LevelName controller
	levelNameController := &controllers.LevelNameControllerImpl{Service: labelService}

	// Configuration de l'application Gin
	r := gin.Default()
	appUserGroup := r.Group("/api/v1/app")
	{
		appUserGroup.GET("/words", wordDtoController.ListDtoWords)    // query param: ids, lang
		appUserGroup.GET("/words/:id", wordDtoController.ReadDtoWord) // query param: lang
		appUserGroup.GET("/tags", tagController.ListTags)
		appUserGroup.GET("/categories", categoryController.ListCategories)
		appUserGroup.GET("/categories/:id/levelNames", levelNameController.ListLevelNames)
	}

	techGroup := r.Group("/api/v1/tech")
	{
		techGroup.GET("/words/:id", wordController.ReadWord)
		techGroup.POST("/words", wordController.CreateWord)
		techGroup.PUT("/words/:id", wordController.UpdateWord)
		techGroup.DELETE("/words/:id", wordController.DeleteWord)

		techGroup.GET("/labels/:id", labelController.ReadLabel)
		techGroup.POST("/labels", labelController.CreateLabel)
		techGroup.PUT("/labels/:id", labelController.UpdateLabel)
		techGroup.DELETE("/labels/:id", labelController.DeleteLabel)
	}
	return r
}

func DatabaseConnectionFromConfig(cfg *config.Config) (*gorm.DB, error) {
	// Use the loaded configuration values
	dbConfig := cfg.Database
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d",
		dbConfig.Host, dbConfig.User, dbConfig.Password, dbConfig.Name, dbConfig.Port)

	db, err := DatabaseConnection(dsn)
	return db, err
}

func DatabaseConnection(dsn string) (*gorm.DB, error) {
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	// Auto-migrate the models to keep the schema in sync
	err = db.AutoMigrate(&models.Label{}, &models.Word{}, &models.Level{})
	if err != nil {
		log.Fatal("failed to migrate database:", err)
	}
	return db, err
}

// Endpoints à implémenter
//CRUD /labels
//CRUD /levels
// GET /level-names (liste levelNames)
//CRUD /words
//POST /words/{id}/tags {"tags": ["uuid1", "uuid2"]} Ajouter des tags à un mot
//DELETE /words/{id}/tags {"tags": ["uuid1", "uuid2"]} Retirer des tags à un mot
//POST /levels/{id}/words {"words": ["uuid1", "uuid2"]} Ajouter des mots à un level
//DELETE /levels/{id}/words {"words": ["uuid1", "uuid2"]} Retirer des mots à un level
//
//Validation via middleware github.com/go-playground/validator
