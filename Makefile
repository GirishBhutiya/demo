SHELL=bash
FRONT_END_BINARY=frontApp
BACK_END_BINARY=backendApp


## up: starts all containers in the background without forcing build
up:
	@echo Starting Docker images...
	docker-compose up -d
	@echo Docker images started!

## up_build: stops docker-compose (if running), builds all projects and starts docker compose
up_build: build_backend 
	@echo Stopping docker images if running...
	docker-compose down
	@echo Building when required and starting docker images...
	docker-compose up --build -d
	@echo Docker images built and started!

## down: stop docker compose
down:
	@echo Stopping docker compose...
	docker-compose down
	@echo Done!

## build_backend: builds the backend binary as a linux executable
build_backend:
	@echo Building backend binary...
	cd backend-service && set GOOS=linux && set GOARCH=amd64 && set CGO_ENABLED=0 && go build -o ${BACK_END_BINARY} ./cmd/api
	@echo Done!


## build_front: builds the broker binary as a linux executable
build_front:
	@echo Building front end linux binary...
	cd ./front-end && set GOOS=linux&& set GOARCH=amd64&& set CGO_ENABLED=0 && go build -o ${FRONT_END_BINARY} ./cmd/web
	@echo Done!

## build_mailer: builds the mailer binary as a linux executable
build_mailer:
	@echo Building mailer binary...
	cd ./mail-service && set GOOS=linux&& set GOARCH=amd64&& set CGO_ENABLED=0 && go build -o ${MAILER_BINARY} ./cmd/api
	@echo Done!

## start: starts the front end
start: build_front
	@echo Starting front end
	cd ./front-end && ./${FRONT_END_BINARY} &

## stop: stop the front end
stop:
	@echo Stopping front end...
	@-pkill -SIGTERM -f "./${FRONT_END_BINARY}"
	@echo "Stopped front end!"

migrateup:
	migrate -path ./backend-service/db/migration/ -database "postgresql://postgres:password@localhost:5432/postcreation?sslmode=disable" -verbose up

migratedown:
	migrate -path ./backend-service/db/migration/ -database "postgresql://postgres:password@localhost:5432/postcreation?sslmode=disable" -verbose down