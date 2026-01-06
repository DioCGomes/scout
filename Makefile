# Scout - Dependency Vulnerability Scanner

BINARY := scout
BUILD_DIR := ./cmd/scout
OUTPUT_DIR := ./bin

.PHONY: build clean test run docker-build

build:
	@mkdir -p $(OUTPUT_DIR)
	CGO_ENABLED=1 go build -o $(OUTPUT_DIR)/$(BINARY) $(BUILD_DIR)

clean:
	rm -rf $(OUTPUT_DIR)
	rm -f scout_report.*

test:
	go test -v ./...

run: build
	$(OUTPUT_DIR)/$(BINARY) $(ARGS)

docker-build:
	docker build -t $(BINARY) .
