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
```bash
colima start
```

### Test docker
```bash
docker run hello-world
```

## Go configuration
### Configure env variables
```bash
go env -w GOSUMDB=sum.golang.org
go env -w GOPROXY=direct
```

