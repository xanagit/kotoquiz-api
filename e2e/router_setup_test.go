package main

import (
	"errors"
	"github.com/xanagit/kotoquiz-api/config"
	"github.com/xanagit/kotoquiz-api/initialisation"
	"github.com/xanagit/kotoquiz-api/middlewares"
	"log"
	"os"
	"os/exec"
	"sync"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type MockAuthMiddleware struct{}

var (
	containerID string
	router      *gin.Engine
	setupOnce   sync.Once
	ready       sync.WaitGroup
	logger      *zap.Logger
)

func (m *MockAuthMiddleware) AuthRequired() gin.HandlerFunc {
	return func(c *gin.Context) {
		mockClaims := middlewares.Claims{
			Subject: "test-user",
			RealmAccess: struct {
				Roles []string `json:"roles"`
			}{
				Roles: []string{"user", "admin"},
			},
		}
		c.Set("claims", mockClaims)
		c.Next()
	}
}

func (m *MockAuthMiddleware) RequireRoles(_ ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next() // In test mode accept all roles
	}
}

type MockCorsMiddleware struct{}

func (cm *MockCorsMiddleware) HandleCORS() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()
	}
}

func InitializeMiddlewareComponents() (*initialisation.MiddlewareComponents, error) {
	return &initialisation.MiddlewareComponents{
		AuthMiddleware: &MockAuthMiddleware{},
		CORSMiddleware: &MockCorsMiddleware{},
	}, nil
}

// initLogger initialise le logger global Zap
func initLogger() {
	var err error
	logger, err = zap.NewProduction() // Utilisez zap.NewDevelopment() pour plus de verbosité en dev
	if err != nil {
		panic("Failed to initialize logger: " + err.Error())
	}
	defer func(logger *zap.Logger) {
		err := logger.Sync()
		if err != nil {

		}
	}(logger)
}

// startPostgresContainer démarre un conteneur PostgreSQL
func startPostgresContainer() (string, error) {
	cmd := exec.Command("docker", "run", "--rm", "-d", "-e", "POSTGRES_PASSWORD=password", "-e", "POSTGRES_DB=testdb", "-p", "5433:5432", "postgres:latest")
	output, err := cmd.CombinedOutput()
	if err != nil {
		logger.Error("Failed to start PostgreSQL container", zap.Error(err), zap.String("output", string(output)))
		return "", err
	}
	containerID := string(output[:len(output)-1]) // Remove the newline character
	logger.Info("PostgreSQL container started", zap.String("containerID", containerID))
	return containerID, nil
}

// stopPostgresContainer arrête le conteneur PostgreSQL
func stopPostgresContainer() error {
	cmd := exec.Command("docker", "stop", containerID)
	output, err := cmd.CombinedOutput()
	if err != nil {
		logger.Error("Failed to stop PostgreSQL container", zap.Error(err), zap.String("output", string(output)))
		return err
	}
	logger.Info("PostgreSQL container stopped", zap.String("containerID", containerID))
	return nil
}

// setupEnvironment initialise l'environnement global
func setupEnvironment() {
	setupOnce.Do(func() {
		logger.Info("Setting up environment...")
		ready.Add(1) // Ajouter une tâche pour la synchronisation

		// Lancer le conteneur PostgreSQL
		var err error
		containerID, err = startPostgresContainer()
		if err != nil {
			logger.Error("Failed to start PostgreSQL container", zap.Error(err))
			ready.Done()
			return
		}

		// Vérifier si la base est prête
		dsn := "host=localhost user=postgres password=password dbname=testdb port=5433 sslmode=disable"
		if err := waitForDatabase(dsn, 10); err != nil {
			logger.Error("Database did not become ready", zap.Error(err))
			stopErr := stopPostgresContainer()
			if stopErr != nil {
				logger.Error("Failed to stop PostgreSQL container after failure", zap.Error(stopErr))
			}
			ready.Done()
			return
		}

		ready.Done()
	})
}

// waitForDatabase vérifie si la base est prête
func waitForDatabase(dsn string, retries int) error {
	time.Sleep(2 * time.Second) // Attendre un peu avant de vérifier la base
	for i := 0; i < retries; i++ {
		db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
		if err == nil {
			sqlDB, _ := db.DB()
			if err := sqlDB.Ping(); err == nil {
				logger.Info("Database is ready!")
				return nil
			}
		}
		logger.Info("Waiting for database...", zap.Int("attempt", i+1))
		time.Sleep(2 * time.Second)
	}
	logger.Error("Database not ready after retries", zap.Int("retries", retries))
	return errors.New("database not ready after " + string(rune(retries)) + " retries")
}

// setupRouter configure le routeur Gin
func setupRouter() (*gin.Engine, error) {
	// S'assurer que l'environnement est prêt
	setupEnvironment()
	ready.Wait() // Attendre que l'environnement soit prêt

	if containerID == "" {
		return nil, errors.New("PostgreSQL container was not initialized")
	}

	dsn := "host=localhost user=postgres password=password dbname=testdb port=5433 sslmode=disable"
	db, err := initialisation.DatabaseConnection(dsn)
	if err != nil {
		logger.Error("Failed to migrate database", zap.Error(err))
		return nil, err
	}

	cfg := &config.Config{
		Database: config.DatabaseConfig{
			Host:     "localhost",
			User:     "postgres",
			Password: "password",
			Name:     "testdb",
			Port:     5433,
		},
	}

	components := initialisation.InitializeAppComponents(db, cfg)
	middlewareComponents, mcErr := InitializeMiddlewareComponents()
	if mcErr != nil {
		log.Fatalf("Failed to initialize app components: %v", err)
	}
	// Gin application configuration
	r := gin.Default()
	initialisation.ConfigureRoutes(r, components, middlewareComponents)

	logger.Info("router initialized", zap.Any("router", router))

	return r, nil
}

func TestMain(m *testing.M) {
	initLogger()
	defer func(logger *zap.Logger) {
		_ = logger.Sync()
	}(logger)

	defer func() {
		// recover from panic if one occurred. Set err to nil otherwise.
		if recover() != nil {
			logger.Error("An error occured during tests")
			logger.Info("Stopping PostgreSQL container...")
			err := stopPostgresContainer()
			if err != nil {
				logger.Error("Failed to stop PostgreSQL container", zap.Error(err))
			}
		}
	}()

	var err error
	router, err = setupRouter()
	if err != nil || router == nil {
		logger.Error("failed to setup router: %v", zap.Error(err))
		os.Exit(1)
	}
	logger.Info("router setup complete", zap.Any("router", router))

	// Exécution des tests
	exitCode := m.Run()
	print(exitCode)

	// Nettoyage après l'exécution des tests
	logger.Info("Stopping PostgreSQL container...")
	err = stopPostgresContainer()
	if err != nil {
		logger.Error("Failed to stop PostgreSQL container", zap.Error(err))
	}

	// Quitter avec le code de sortie des tests
	os.Exit(exitCode)
}
