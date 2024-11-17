# Setup environment
## Install docker
### Install brew
```bash
/bin/bash -c "$(curl -fsSL https://raw.githubusercontent.com/Homebrew/install/HEAD/install.sh)"
```

### Install Colima
```bash
brew install colima
```

### Install Docker cli
```bash
brew install docker
```

### Start Docker
Démarrer colima avec docker buildkit (nécessaire pour builder les images)
```bash
DOCKER_BUILDKIT=1 colima start
```

### Test docker
```bash
docker run hello-world
```

## Install docker compose
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

Test :
```zsh
docker compose version
Docker Compose version v2.29.2
```

## Install docker buildkit
```zsh
curl -SL https://github.com/docker/buildx/releases/download/v0.17.0/buildx-v0.17.0.darwin-arm64 -o $DOCKER_CONFIG/cli-plugins/docker-buildx
chmod +x $DOCKER_CONFIG/cli-plugins/docker-buildx
```

Exporter la variable DOCKER_BUILDKIT
```zsh 
export DOCKER_BUILDKIT=1
```

## Start Postgres
```zsh
docker compose up -d
```

## Go configuration
### Configure env variables
```bash
go env -w GOSUMDB=sum.golang.org
go env -w GOPROXY=direct
```
