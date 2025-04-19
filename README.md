# KotoQuiz API

A Golang-based RESTful API for the japanese language learning application KotoQuiz, managing vocabulary, quizzes, and user progress tracking.

## Architecture

The KotoQuiz API follows a clean architecture pattern with clear separation of concerns:

- **Controllers**: Handle HTTP requests and responses
- **Services**: Implement business logic and orchestrate operations
- **Repositories**: Manage data persistence and retrieval
- **Models**: Define data structures and relationships
- **DTOs**: Data Transfer Objects for API communication
- **Middlewares**: Handle cross-cutting concerns like authentication and CORS

The application uses a layered approach where:
1. Controllers receive requests and delegate to appropriate services
2. Services implement business logic and interact with repositories
3. Repositories handle database operations
4. Models define the core domain entities

## Project Structure

```
kotoquiz-api/
├── cmd/                 # Application entry point
├── config/              # Configuration files and loading logic
├── controllers/         # HTTP request handlers
├── dto/                 # Data Transfer Objects
├── initialisation/      # App initialization and dependency wiring
├── middlewares/         # Request processing middleware
├── models/              # Domain models
├── repositories/        # Data access layer
├── services/            # Business logic implementation
├── e2e/                 # End-to-end tests
├── swagger.yaml         # API documentation
```

## Technologies

- **Language**: Go 1.23.3
- **Web Framework**: Gin Gonic
- **ORM**: GORM with PostgreSQL driver
- **Authentication**: Keycloak with OpenID Connect
- **Configuration**: Viper for config management
- **Logging**: Zap logger
- **Documentation**: Swagger/OpenAPI
- **Containerization**: Docker and Docker Compose
- **Testing**: Go's testing package with Testify

## Setup Environment

### Install Docker

#### Install Brew
```bash
/bin/bash -c "$(curl -fsSL https://raw.githubusercontent.com/Homebrew/install/HEAD/install.sh)"
```

#### Install Colima
```bash
brew install colima
```

#### Install Docker CLI
```bash
brew install docker
```

#### Start Docker
Start Colima with Docker buildkit (required for building images)
```bash
DOCKER_BUILDKIT=1 colima start
```

#### Test Docker
```bash
docker run hello-world
```

### Install Docker Compose
```zsh
DOCKER_CONFIG=${DOCKER_CONFIG:-$HOME/.docker}
mkdir -p $DOCKER_CONFIG/cli-plugins
curl -SL https://github.com/docker/compose/releases/download/v2.29.2/docker-compose-darwin-aarch64 -o $DOCKER_CONFIG/cli-plugins/docker-compose
chmod +x $DOCKER_CONFIG/cli-plugins/docker-compose
```

Add DOCKER_CONFIG to .zshrc
```zsh
echo 'export DOCKER_CONFIG=${DOCKER_CONFIG:-$HOME/.docker}' >> ~/.zshrc
source ~/.zshrc
```

Test:
```zsh
docker compose version
Docker Compose version v2.29.2
```

### Install Docker Buildkit
```zsh
curl -SL https://github.com/docker/buildx/releases/download/v0.17.0/buildx-v0.17.0.darwin-arm64 -o $DOCKER_CONFIG/cli-plugins/docker-buildx
chmod +x $DOCKER_CONFIG/cli-plugins/docker-buildx
```

Export the DOCKER_BUILDKIT variable
```zsh 
export DOCKER_BUILDKIT=1
```

### Start Postgres
```zsh
docker compose up -d
```

### Go Configuration
#### Configure env variables
```bash
go env -w GOSUMDB=sum.golang.org
go env -w GOPROXY=direct
```

## Development

### Recreate Database
```zsh
psql -h localhost -U admin postgres
> DROP DATABASE kotoquiz;
> CREATE DATABASE kotoquiz;
```

## Configuration for Flutter Application

[POST] http://localhost:8180/realms/kotoquiz/protocol/openid-connect/registrations
```json
{
    "username": "newuser",
    "email": "user@example.com",
    "firstName": "John",
    "lastName": "Doe",
    "password": "password",
    "enabled": true
}
```