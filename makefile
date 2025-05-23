include /$(PWD)/.env

##############################################3
# Variables

MONOLITH_IMAGE := critiquefy-service/monolith/$(CRITIQUEFY_VERSION_NUMBER)

##############################################3
# Building and Running

run:
	go run ./api/monolith

build-docker:
	docker build \
		-f zarf/docker/dockerfile.monolith \
		-t $(MONOLITH_IMAGE) \
		--build-arg BUILD_REF=$(CRITIQUEFY_VERSION_NUMBER) \
		--build-arg BUILD_DATE=$(date -u +"%Y-%m-%dT%H:%M:%SZ") \
		.

run-docker:
	docker run -p 3000:3000 -p 3010:3010 $(MONOLITH_IMAGE) 

##############################################3
# Testing

test:
	go test ./...

test-race:
	CGO_ENABLED=1 go test -race -count=1 ./...

test-coverage:
	go test -cover ./...

test-coverage-profile:
	go test ./... -coverprofile=.profs/coverage.out
	go tool cover -html=.profs/coverage.out

lint:
	go vet ./...

##############################################3
# Administrative Tasks

admin-genkey:
	go run tooling/admin/main.go genkey

pgcli:
	pgcli $(CRITIQUEFY_DB_URL_CLI)