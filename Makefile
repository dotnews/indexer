default: test

.PHONY: test compose index delete

test:
	@CONFIG=config/test.config.json GOCACHE=off go test ./...

compose:
	@docker-compose up

index:
	@go run main.go -logtostderr=true

delete:
	@go run main.go -logtostderr=true -process=delete
