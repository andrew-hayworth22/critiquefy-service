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