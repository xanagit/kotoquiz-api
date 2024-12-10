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
	WordService                services.WordService
	LabelService               services.LabelService
	LevelService               services.LevelService
	WordLearningHistoryService services.WordLearningHistoryService
	WordDtoService             services.WordDtoService

	// Controllers
	WordController                controllers.WordController
	LevelController               controllers.LevelController
	TagController                 controllers.TagController
	WordLearningHistoryController controllers.WordLearningHistoryController
	WordDtoController             controllers.WordDtoController
}

type MiddlewareComponents struct {
	// Middlewares
	AuthMiddleware middlewares.AuthMiddleware
}

func InitializeAppComponents(db *gorm.DB) *AppComponents {
	// Repositories
	wordRepo := &repositories.WordRepositoryImpl{DB: db}
	labelRepo := &repositories.LabelRepositoryImpl{DB: db}
	levelRepo := &repositories.LevelRepositoryImpl{DB: db}
	wordLearningHistoryRepo := &repositories.WordLearningHistoryRepositoryImpl{DB: db}

	// Services
	wordService := &services.WordServiceImpl{Repo: wordRepo}
	labelService := &services.LabelServiceImpl{Repo: labelRepo}
	levelService := &services.LevelServiceImpl{Repo: levelRepo}
	wordLearningHistoryService := &services.WordLearningHistoryServiceImpl{Repo: wordLearningHistoryRepo}
	wordDtoService := &services.WordDtoServiceImpl{
		WordRepo:            wordRepo,
		LearningHistoryRepo: wordLearningHistoryRepo,
	}

	// Controllers
	wordController := &controllers.WordControllerImpl{Service: wordService}
	levelController := &controllers.LevelControllerImpl{Service: levelService}
	tagController := &controllers.TagControllerImpl{Service: labelService}
	wordLearningHistoryController := &controllers.WordLearningHistoryControllerImpl{Service: wordLearningHistoryService}
	wordDtoController := &controllers.WordDtoControllerImpl{WordDtoService: wordDtoService}

	// Return an instance of AppComponents with interfaces
	return &AppComponents{
		// Repositories
		WordRepository:                wordRepo,
		LabelRepository:               labelRepo,
		LevelRepository:               levelRepo,
		WordLearningHistoryRepository: wordLearningHistoryRepo,

		// Services
		WordService:                wordService,
		LabelService:               labelService,
		LevelService:               levelService,
		WordLearningHistoryService: wordLearningHistoryService,
		WordDtoService:             wordDtoService,

		// Controllers
		WordController:                wordController,
		LevelController:               levelController,
		TagController:                 tagController,
		WordLearningHistoryController: wordLearningHistoryController,
		WordDtoController:             wordDtoController,
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
