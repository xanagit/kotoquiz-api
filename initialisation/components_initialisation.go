package initialisation

import (
	"fmt"
	"github.com/xanagit/kotoquiz-api/config"
	"github.com/xanagit/kotoquiz-api/controllers"
	"github.com/xanagit/kotoquiz-api/middlewares"
	"github.com/xanagit/kotoquiz-api/repositories"
	"github.com/xanagit/kotoquiz-api/services"
	"gorm.io/gorm"
)

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

type MiddlewareComponents struct {
	// Middlewares
	AuthMiddleware middlewares.AuthMiddleware
}

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

func InitializeMiddlewareComponents(cfg *config.Config) (*MiddlewareComponents, error) {
	authMiddleware, err := middlewares.NewAuthMiddleware(&cfg.Auth.Keycloak)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize auth middleware: %v", err)
	}

	return &MiddlewareComponents{
		AuthMiddleware: authMiddleware,
	}, nil
}
