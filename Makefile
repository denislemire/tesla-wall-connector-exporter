VERSION ?= v0.1.0
IMAGE ?= ghcr.io/denislemire/tesla-wall-connector-exporter
PLATFORM ?= linux/amd64
CHART_DIR := helm/tesla-wall-connector-exporter

.PHONY: test build run docker-build docker-push helm-package helm-lint clean

test:
	go test ./...

build:
	CGO_ENABLED=0 go build -trimpath -ldflags="-s -w" -o bin/tesla-wall-connector-exporter ./cmd/tesla-wall-connector-exporter

run: build
	@test -n "$(TWC_ADDRESS)" || (echo "Set TWC_ADDRESS=hostname-or-ip" && exit 1)
	./bin/tesla-wall-connector-exporter --twc.address=$(TWC_ADDRESS)

docker-build:
	docker buildx build --platform $(PLATFORM) -t $(IMAGE):$(VERSION) .

docker-push:
	docker buildx build --platform $(PLATFORM) -t $(IMAGE):$(VERSION) --push .

helm-lint:
	helm lint $(CHART_DIR)

helm-package:
	helm package $(CHART_DIR) --version $(VERSION:v%=%) --app-version $(VERSION:v%=%)

clean:
	rm -rf bin/ *.tgz
