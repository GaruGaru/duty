fmt:
	go fmt ./...

test:
	go test ./...

deps:
	go mod vendor
	go mod verify
