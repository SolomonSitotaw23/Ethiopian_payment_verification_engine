.PHONY: run build test clean docker-build docker-run lint

APP_NAME = server

run:
	go run main.go

build:
	CGO_ENABLED=0 go build -ldflags="-w -s" -o $(APP_NAME) main.go

test:
	go test -v -cover ./...

clean:
	rm -f $(APP_NAME)

docker-build:
	docker build -t payment-verifier-go .

docker-run:
	docker run -p 5000:5000 --env-file .env payment-verifier-go
