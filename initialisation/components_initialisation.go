// Package initialisation handles application component initialization and wiring
package initialisation

import (
	"fmt"
	"github.com/xanagit/kotoquiz-api/config"
	"github.com/xanagit/kotoquiz-api/controllers"
	"github.com/xanagit/kotoquiz-api/middlewares"
	"github.com/xanagit/kotoquiz-api/repositories"
	"github.com/xanagit/kotoquiz-api/services"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

// AppComponents holds all the application components (repositories, services, controllers)
// used for dependency injection throughout the application
type AppComponents struct {
	// Repositories
	WordRepository                repositories.WordRepository
	LabelRepository               repositories.LabelRepository
	LevelRepository               repositories.LevelRepository
	WordLearningHistoryRepository repositories.WordLearningHistoryRepository

	// Services
	HealthService              services.ApiHealthService
	WordService                services.WordService
	LabelService               services.LabelService
	LevelService               services.LevelService
	WordLearningHistoryService services.WordLearningHistoryService
	WordDtoService             services.WordDtoService
	RegistrationService        services.RegistrationService

	// Controllers
	HealthController              controllers.HealthController
	WordController                controllers.WordController
	LevelController               controllers.LevelController
	TagController                 controllers.TagController
	WordLearningHistoryController controllers.WordLearningHistoryController
	WordDtoController             controllers.WordDtoController
	RegistrationController        controllers.RegistrationController
}

// MiddlewareComponents holds all middleware components used across the application
type MiddlewareComponents struct {
	// Middlewares
	CORSMiddleware middlewares.CORSMiddleware
	AuthMiddleware middlewares.AuthMiddleware
}

// InitializeAppComponents creates and wires together all application components
// This function implements the dependency injection pattern, creating repositories,
// services, and controllers and connecting them appropriately.
//
// Parameters:
//   - db: *gorm.DB - The database connection to use for repositories
//   - cfg: *config.Config - The application configuration
//
// Returns:
//   - *AppComponents - The initialized application components
func InitializeAppComponents(db *gorm.DB, cfg *config.Config) *AppComponents {
	// Repositories
	wordRepo := &repositories.WordRepositoryImpl{DB: db}
	labelRepo := &repositories.LabelRepositoryImpl{DB: db}
	levelRepo := &repositories.LevelRepositoryImpl{DB: db}
	wordLearningHistoryRepo := &repositories.WordLearningHistoryRepositoryImpl{DB: db}

	// Services
	healthService := &services.ApiHealthServiceImpl{DB: db}
	wordService := &services.WordServiceImpl{Repo: wordRepo}
	labelService := &services.LabelServiceImpl{Repo: labelRepo}
	levelService := &services.LevelServiceImpl{Repo: levelRepo}
	wordLearningHistoryService := &services.WordLearningHistoryServiceImpl{Repo: wordLearningHistoryRepo}
	wordDtoService := &services.WordDtoServiceImpl{
		WordRepo:            wordRepo,
		LearningHistoryRepo: wordLearningHistoryRepo,
	}
	registrationService := &services.RegistrationServiceImpl{KeycloakConfig: &cfg.Auth.Keycloak}

	// Controllers
	healthController := &controllers.HealthControllerImpl{Service: healthService}
	wordController := &controllers.WordControllerImpl{Service: wordService}
	levelController := &controllers.LevelControllerImpl{Service: levelService}
	tagController := &controllers.TagControllerImpl{Service: labelService}
	wordLearningHistoryController := &controllers.WordLearningHistoryControllerImpl{Service: wordLearningHistoryService}
	wordDtoController := &controllers.WordDtoControllerImpl{WordDtoService: wordDtoService}
	registrationController := &controllers.RegistrationControllerImpl{Service: registrationService}

	// Return an instance of AppComponents with interfaces
	return &AppComponents{
		// Repositories
		WordRepository:                wordRepo,
		LabelRepository:               labelRepo,
		LevelRepository:               levelRepo,
		WordLearningHistoryRepository: wordLearningHistoryRepo,

		// Services
		HealthService:              healthService,
		WordService:                wordService,
		LabelService:               labelService,
		LevelService:               levelService,
		WordLearningHistoryService: wordLearningHistoryService,
		WordDtoService:             wordDtoService,
		RegistrationService:        registrationService,

		// Controllers
		HealthController:              healthController,
		WordController:                wordController,
		LevelController:               levelController,
		TagController:                 tagController,
		WordLearningHistoryController: wordLearningHistoryController,
		WordDtoController:             wordDtoController,
		RegistrationController:        registrationController,
	}
}

// InitializeMiddlewareComponents creates and initializes middleware components
//
// Parameters:
//   - cfg: *config.Config - The application configuration containing middleware settings
//   - log: *zap.Logger - Logger for middleware initialization
//
// Returns:
//   - *MiddlewareComponents - The initialized middleware components
//   - error - Any error that occurred during initialization
func InitializeMiddlewareComponents(cfg *config.Config, log *zap.Logger) (*MiddlewareComponents, error) {
	corsMiddleware, corsErr := middlewares.NewCORSMiddleware(&cfg.Auth.ApiConfig)
	if corsErr != nil {
		log.Error("Failed to initialize CORS middleware", zap.Error(corsErr))
		return nil, fmt.Errorf("failed to initialize cors middleware: %v", corsErr)
	}

	authMiddleware, authErr := middlewares.NewAuthMiddleware(&cfg.Auth.Keycloak, log)
	if authErr != nil {
		log.Error("Failed to initialize auth middleware", zap.Error(authErr))
		return nil, fmt.Errorf("failed to initialize auth middleware: %v", authErr)
	}

	log.Info("Middleware components initialized successfully")

	return &MiddlewareComponents{
		CORSMiddleware: corsMiddleware,
		AuthMiddleware: authMiddleware,
	}, nil
}
