default: test

.PHONY: test index delete

test:
	@CONFIG=config/test.config.json GOCACHE=off go test ./...

index:
	@go run main.go -logtostderr=true

delete:
	@go run main.go -logtostderr=true -process=delete
