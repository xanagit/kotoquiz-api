package initialisation

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/xanagit/kotoquiz-api/config"
	"github.com/xanagit/kotoquiz-api/controllers"
	"github.com/xanagit/kotoquiz-api/middlewares"
	"github.com/xanagit/kotoquiz-api/models"
	"github.com/xanagit/kotoquiz-api/repositories"
	"github.com/xanagit/kotoquiz-api/services"
	"golang.org/x/time/rate"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
	"time"
)

func GinHandlers(cfg *config.Config, db *gorm.DB) *gin.Engine {
	// Word repository, service and controller
	wordRepository := &repositories.WordRepositoryImpl{DB: db}
	wordService := &services.WordServiceImpl{Repo: wordRepository}
	wordController := &controllers.WordControllerImpl{Service: wordService}

	// Label repository, service and controller
	labelRepository := &repositories.LabelRepositoryImpl{DB: db}
	labelService := &services.LabelServiceImpl{Repo: labelRepository}

	// Level repository, service and controller
	levelRepository := &repositories.LevelRepositoryImpl{DB: db}
	levelService := &services.LevelServiceImpl{Repo: levelRepository}
	levelController := &controllers.LevelControllerImpl{Service: levelService}

	// Tag controller
	tagController := &controllers.TagControllerImpl{Service: labelService}

	// WordLearningHistory repository, service and controller
	wordLearningHistoryRepository := &repositories.WordLearningHistoryRepositoryImpl{DB: db}
	wordLearningHistoryService := &services.WordLearningHistoryServiceImpl{Repo: wordLearningHistoryRepository}
	wordLearningHistoryController := &controllers.WordLearningHistoryControllerImpl{Service: wordLearningHistoryService}

	// User repository, service and controller
	userRepository := &repositories.UserRepositoryImpl{DB: db}
	userService := &services.UserServiceImpl{Repo: userRepository}
	userController := &controllers.UserControllerImpl{Service: userService}

	// WordDto service and controller
	wordDtoService := &services.WordDtoServiceImpl{WordRepo: wordRepository, LearningHistoryRepo: wordLearningHistoryRepository}
	wordDtoController := &controllers.WordDtoControllerImpl{WordDtoService: wordDtoService}

	// Auth service and controller
	authService := &services.AuthServiceImpl{JwtSecret: getJwtSecret(cfg), UserService: userService}
	authController := &controllers.AuthControllerImpl{AuthService: authService, UserService: userService}

	// Headers and security controller
	headersAndSecurityController := &controllers.HeadersAndSecurityControllerImpl{Config: *cfg}

	// Auth middleware
	authMiddleware := &middlewares.AuthenticationMiddlewareImpl{AuthService: authService}
	// Rate limiters
	authLimiter := middlewares.NewRateLimiterImpl(rate.Every(time.Minute), 30) // 30 req/min
	apiLimiter := middlewares.NewRateLimiterImpl(rate.Every(time.Second), 10)

	r := gin.Default()
	r.Use(headersAndSecurityController.AddHeaders)

	r.POST("/csp-report", headersAndSecurityController.CspReport) // Content Security Policy violation report

	// Auth routes
	auth := r.Group("/api/v1/auth")
	auth.Use(authLimiter.LimitRate())
	{
		auth.POST("/login", authController.Login)
		auth.POST("/refresh", authController.RefreshToken)
	}

	// App user routes
	appUserGroup := r.Group("/api/v1/app")
	appUserGroup.Use(authMiddleware.AuthRequired(models.RoleAppUser))
	appUserGroup.Use(apiLimiter.LimitRate())
	{
		appUserGroup.GET("/words/q", wordDtoController.ListWordsIDs)
		appUserGroup.GET("/words", wordDtoController.ListDtoWords)    // query param: ids, lang
		appUserGroup.GET("/words/:id", wordDtoController.ReadDtoWord) // query param: lang
		appUserGroup.GET("/tags", tagController.ListTags)
		appUserGroup.GET("/levels", levelController.ListLevels)
		appUserGroup.POST("/quiz/results", wordLearningHistoryController.ProcessQuizResults)
	}

	// Tech routes
	techGroup := r.Group("/api/v1/tech")
	techGroup.Use(authMiddleware.AuthRequired(models.RoleTech))
	techGroup.Use(apiLimiter.LimitRate())
	{
		techGroup.GET("/words/:id", wordController.ReadWord)
		techGroup.POST("/words", wordController.CreateWord)
		techGroup.PUT("/words/:id", wordController.UpdateWord)
		techGroup.DELETE("/words/:id", wordController.DeleteWord)

		techGroup.GET("/tags/:id", tagController.ReadTag)
		techGroup.POST("/tags", tagController.CreateTag)
		techGroup.PUT("/tags/:id", tagController.UpdateTag)
		techGroup.DELETE("/tags/:id", tagController.DeleteTag)

		techGroup.GET("/levels/:id", levelController.ReadLevel)
		techGroup.POST("/levels", levelController.CreateLevel)
		techGroup.PUT("/levels/:id", levelController.UpdateLevel)
		techGroup.DELETE("/levels/:id", levelController.DeleteLevel)

		techGroup.POST("/users", userController.CreateUser)
		techGroup.GET("/users/:id", userController.ReadUser)
		techGroup.PUT("/users/:id", userController.UpdateUser)
		techGroup.DELETE("/users/:id", userController.DeleteUser)
	}
	return r
}

func getJwtSecret(cfg *config.Config) []byte {
	jwtConfig := cfg.Jwt
	jwtSecret := []byte(jwtConfig.SecretKey)
	return jwtSecret
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

// 2. Validation des tokens
// 3. Gestion du refresh token
