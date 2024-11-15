# Variables
BIN_DIR = bin
CONFIG_SRC = config/config-prod.yml
CONFIG_DEST = $(BIN_DIR)/config/config.yml
BINARY = $(BIN_DIR)/main

# Cibles
all: build

# Construction du binaire en mode développement
build:
	@echo "Building the binary..."
	@mkdir -p $(BIN_DIR)
	@go build -o $(BINARY) ./cmd

# Construction du binaire en mode production
build-arm64:
	@echo "Building the binary..."
	@mkdir -p $(BIN_DIR)
	@GOOS=linux GOARCH=arm64 go build -o $(BINARY) ./cmd

# Copie du fichier de configuration
prepare-config:
	@echo "Preparing configuration file..."
	@mkdir -p $(BIN_DIR)/config
	@cp $(CONFIG_SRC) $(CONFIG_DEST)

# Cible finale pour automatiser tout
release: build prepare-config
	@echo "Build and configuration setup completed."

release-arm64: build-arm64 prepare-config
	@echo "Build and configuration setup completed."

# Nettoyage des fichiers générés
clean:
	@echo "Cleaning up..."
	@rm -rf $(BIN_DIR)

.PHONY: all build prepare-config release clean
