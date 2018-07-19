run: fmt vet
	@go run main.go
fmt:
	@gofmt -s -w -l .
vet:
	@go vet ./...

.PHONY: fmt vet
