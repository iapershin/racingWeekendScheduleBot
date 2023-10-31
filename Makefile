lint:
	golangci-lint run;

test:
	go test -race ./...;

run:
	go run cmd/race-weekend-bot/main.go;

build:
	go build -o cmd/race-weekend-bot/main cmd/race-weekend-bot/main.go;

docker-compose-up:
	docker-compose -f deployment/docker-compose.yml  up -d --remove-orphans --build;

docker-compose-down:
	docker-compose -f deployment/docker-compose.yml  down;