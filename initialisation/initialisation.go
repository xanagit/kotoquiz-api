package initialisation

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/xanagit/kotoquiz-api/config"
	"github.com/xanagit/kotoquiz-api/middlewares"
	"github.com/xanagit/kotoquiz-api/models"
	"go.uber.org/zap"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func ConfigureRoutes(r *gin.Engine, components *AppComponents, middlewareComponents *MiddlewareComponents, log *zap.Logger) {
	r.Use(middlewareComponents.CORSMiddleware.HandleCORS())
	r.GET("/health", components.HealthController.HealthCheck)

	public := r.Group("/api/v1/public")
	public.Use(middlewareComponents.CORSMiddleware.HandleCORS())
	{
		public.POST("/register", components.RegistrationController.RegisterUser)
	}

	appUserGroup := r.Group("/api/v1/app")
	appUserGroup.Use(middlewareComponents.CORSMiddleware.HandleCORS())
	appUserGroup.Use(middlewareComponents.AuthMiddleware.AuthRequired())
	appUserGroup.Use(middlewareComponents.AuthMiddleware.RequireRoles(string(middlewares.UserRole)))
	{
		appUserGroup.GET("/words/q", components.WordDtoController.ListWordsIDs)
		appUserGroup.GET("/words", components.WordDtoController.ListDtoWords)    // query param: ids, lang
		appUserGroup.GET("/words/:id", components.WordDtoController.ReadDtoWord) // query param: lang
		appUserGroup.GET("/tags", components.TagController.ListTags)
		appUserGroup.GET("/levels", components.LevelController.ListLevels)
		appUserGroup.POST("/quiz/results", components.WordLearningHistoryController.ProcessQuizResults)
	}

	techGroup := r.Group("/api/v1/tech")
	techGroup.Use(middlewareComponents.CORSMiddleware.HandleCORS())
	techGroup.Use(middlewareComponents.AuthMiddleware.AuthRequired())
	techGroup.Use(middlewareComponents.AuthMiddleware.RequireRoles(string(middlewares.AdminRole)))
	{
		techGroup.GET("/words/:id", components.WordController.ReadWord)
		techGroup.POST("/words", components.WordController.CreateWord)
		techGroup.PUT("/words/:id", components.WordController.UpdateWord)
		techGroup.DELETE("/words/:id", components.WordController.DeleteWord)

		techGroup.GET("/tags/:id", components.TagController.ReadTag)
		techGroup.POST("/tags", components.TagController.CreateTag)
		techGroup.PUT("/tags/:id", components.TagController.UpdateTag)
		techGroup.DELETE("/tags/:id", components.TagController.DeleteTag)

		techGroup.GET("/levels/:id", components.LevelController.ReadLevel)
		techGroup.POST("/levels", components.LevelController.CreateLevel)
		techGroup.PUT("/levels/:id", components.LevelController.UpdateLevel)
		techGroup.DELETE("/levels/:id", components.LevelController.DeleteLevel)
	}

	log.Info("Routes configured successfully")
}

func DatabaseConnectionFromConfig(cfg *config.Config, log *zap.Logger) (*gorm.DB, error) {
	// Use the loaded configuration values
	dbConfig := cfg.Database
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d",
		dbConfig.Host, dbConfig.User, dbConfig.Password, dbConfig.Name, dbConfig.Port)

	log.Info("Connecting to database",
		zap.String("host", dbConfig.Host),
		zap.String("dbname", dbConfig.Name),
		zap.Int("port", dbConfig.Port))

	db, err := DatabaseConnection(dsn, log)
	return db, err
}

func DatabaseConnection(dsn string, log *zap.Logger) (*gorm.DB, error) {
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Error("Failed to connect to database", zap.Error(err))
		return nil, err
	}

	// Auto-migrate the models to keep the schema in sync
	log.Info("Running database migrations")
	err = db.AutoMigrate(
		&models.Label{},
		&models.Word{},
		&models.Level{},
		&models.WordTag{},
		&models.WordLevel{},
		&models.WordLearningHistory{})
	if err != nil {
		log.Error("Failed to migrate database", zap.Error(err))
		return nil, err
	}

	log.Info("Database connected and migrations complete")
	return db, nil
}

// TODO :Endpoints à implémenter
//POST /levels/{id}/words {"words": ["uuid1", "uuid2"]} Ajouter des mots à un level
//DELETE /levels/{id}/words {"words": ["uuid1", "uuid2"]} Retirer des mots à un level

//Validation via middleware github.com/go-playground/validator
