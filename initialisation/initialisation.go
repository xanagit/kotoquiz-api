package initialisation

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/xanagit/kotoquiz-api/config"
	middleware "github.com/xanagit/kotoquiz-api/middlewares"
	"github.com/xanagit/kotoquiz-api/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
)

func GinHandlers(cfg *config.Config, db *gorm.DB) (*gin.Engine, error) {
	// Initialize the application components
	components := InitializeAppComponents(db)
	var middlewares *MiddlewareComponents
	var err error
	if gin.Mode() == gin.ReleaseMode {
		middlewares, err = InitializeMiddlewareComponents(cfg)
		if err != nil {
			log.Fatalf("Failed to initialize app components: %v", err)
		}
	}

	// Gin application configuration
	r := gin.Default()

	configureRoutes(r, components, middlewares)

	return r, nil
}

func configureRoutes(r *gin.Engine, components *AppComponents, middlewares *MiddlewareComponents) {
	appUserGroup := r.Group("/api/v1/app")
	if middlewares != nil {
		appUserGroup.Use(middlewares.AuthMiddleware.AuthRequired())
		appUserGroup.Use(middlewares.AuthMiddleware.RequireRoles(string(middleware.UserRole)))
	}
	{
		appUserGroup.GET("/words/q", components.WordDtoController.ListWordsIDs)
		appUserGroup.GET("/words", components.WordDtoController.ListDtoWords)    // query param: ids, lang
		appUserGroup.GET("/words/:id", components.WordDtoController.ReadDtoWord) // query param: lang
		appUserGroup.GET("/tags", components.TagController.ListTags)
		appUserGroup.GET("/levels", components.LevelController.ListLevels)
		appUserGroup.POST("/quiz/results", components.WordLearningHistoryController.ProcessQuizResults)
	}

	techGroup := r.Group("/api/v1/tech")
	if middlewares != nil {
		techGroup.Use(middlewares.AuthMiddleware.AuthRequired())
		techGroup.Use(middlewares.AuthMiddleware.RequireRoles(string(middleware.AdminRole)))
	}
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

		techGroup.POST("/users", components.UserController.CreateUser)
		techGroup.GET("/users/:id", components.UserController.ReadUser)
		techGroup.PUT("/users/:id", components.UserController.UpdateUser)
		techGroup.DELETE("/users/:id", components.UserController.DeleteUser)
	}
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
	err = db.AutoMigrate(
		&models.Label{},
		&models.Word{},
		&models.Level{},
		&models.WordTag{},
		&models.WordLevel{},
		&models.User{},
		&models.WordLearningHistory{})
	if err != nil {
		log.Fatal("failed to migrate database:", err)
	}
	return db, err
}

// TODO :Endpoints à implémenter
//POST /levels/{id}/words {"words": ["uuid1", "uuid2"]} Ajouter des mots à un level
//DELETE /levels/{id}/words {"words": ["uuid1", "uuid2"]} Retirer des mots à un level

//Validation via middleware github.com/go-playground/validator
