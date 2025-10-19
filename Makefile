.PHONY: help build run test clean docker-build docker-run k8s-deploy k8s-delete

help:
	@echo "Available commands:"
	@echo "  make build         - Build the Go application"
	@echo "  make run           - Run the application locally"
	@echo "  make test          - Run tests"
	@echo "  make clean         - Clean build artifacts"
	@echo "  make docker-build  - Build Docker image"
	@echo "  make docker-run    - Run with Docker Compose"
	@echo "  make k8s-deploy    - Deploy to Kubernetes"
	@echo "  make k8s-delete    - Delete Kubernetes resources"

build:
	@echo "Building Go application..."
	go build -o netconf-checker ./cmd/main.go

run:
	@echo "Running application..."
	go run cmd/main.go --router-address=localhost:830

test:
	@echo "Running tests..."
	go test -v ./...

clean:
	@echo "Cleaning build artifacts..."
	rm -f netconf-checker
	rm -rf build/

docker-build:
	@echo "Building Docker image..."
	docker build -t netconf-k8s-inspector:local .

docker-run:
	@echo "Starting services with Docker Compose..."
	docker-compose up --build

docker-stop:
	@echo "Stopping Docker Compose services..."
	docker-compose down

k8s-deploy:
	@echo "Deploying to Kubernetes..."
	kubectl apply -f k8s/

k8s-delete:
	@echo "Deleting Kubernetes resources..."
	kubectl delete -f k8s/

k8s-logs-router:
	@echo "Viewing router logs..."
	kubectl logs -l app=netconf-router --tail=50 -f

k8s-logs-checker:
	@echo "Viewing checker logs..."
	kubectl logs -l app=netconf-checker --tail=50

k8s-status:
	@echo "Checking Kubernetes resources status..."
	kubectl get all
	@echo ""
	@echo "CronJob details:"
	kubectl get cronjob netconf-checker-cronjob
	@echo ""
	@echo "Recent jobs:"
	kubectl get jobs --sort-by=.metadata.creationTimestamp | tail -5

fmt:
	@echo "Formatting Go code..."
	go fmt ./...

lint:
	@echo "Linting Go code..."
	golangci-lint run || echo "golangci-lint not installed, skipping..."

deps:
	@echo "Downloading dependencies..."
	go mod download
	go mod tidy

all: clean deps fmt build docker-build
