####################################################################################################
# Config
####################################################################################################

# Load a .env file if present
set dotenv-load := true

# The host to use for the local API Gateway emulator and the API docs preview
dev-host := env("VERDICT_DEV_HOST", "127.0.0.1")

# The tag (version) for the local DynamoDB image. See https://hub.docker.com/r/amazon/dynamodb-local
db-img-tag := "3.3.0"

# The name for the local DynamoDB container. Must match the one used in `initdev.go`.
dev-db := "verdict-dev-db"

# The name for the local Docker network used for DynamoDB
dev-network := "verdict-dev-network"

####################################################################################################
# Development
####################################################################################################

# Build and run the API with DynamoDB in Docker (default recipe)
dev: dev-db
    GOFLAGS="-tags=dev" sam build
    sam local start-api --host {{ dev-host }} --docker-network {{ dev-network }} \
        --parameter-overrides "StageName=dev FrontendUrl=http://localhost:8000"

# Create the local DynamoDB Docker network and container if they don't already exist
[private]
dev-db:
    if ! docker network inspect {{ dev-network }} >/dev/null 2>&1; then \
        docker network create {{ dev-network }}; \
    fi
    if ! docker inspect {{ dev-db }} >/dev/null 2>&1; then \
        docker run -d --rm -e 8000 --network {{ dev-network }} --name {{ dev-db }} \
            amazon/dynamodb-local:{{ db-img-tag }}; \
    fi

# Clean up the dev DB Docker container and network if they exist
dev-clean:
    if docker inspect {{ dev-db }} >/dev/null 2>&1; then \
        docker stop {{ dev-db }}; \
    fi
    if docker network inspect {{ dev-network }} >/dev/null 2>&1; then \
        docker network rm {{ dev-network }}; \
    fi

# Run a server to preview the docs and/or try out the API running locally
docs:
    python3 -m http.server -d docs -b {{ dev-host }}

# Run tests
test:
    go test -tags test ./... -v

####################################################################################################
# Deployment
####################################################################################################

# Run guided deployment
deploy:
    sam build
    sam deploy --guided
