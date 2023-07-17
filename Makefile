export GOBIN := ${PWD}/bin
export PATH  := ${GOBIN}:${PATH}

BINARY_NAME = 'gotagv'
BINARY_PATH = ${GOBIN}/${BINARY_NAME}

# all src packages without generated code
PACKAGES = $(shell go list ./...)

help:
	@echo 'Usage: make <TARGETS> ... <OPTIONS>'
	@echo ''
	@echo 'Available targets are:'
	@echo ''
	@echo '    help               Show this help'
	@echo '    clean              Remove binaries'
	@echo '    download-deps      Download and install dependencies'
	@echo '    tidy               Perform go tidy steps'
	@echo '    lint               Run all linters'
	@echo '    build              Compile packages and dependencies'
	@echo '    migrate-up         Applies migrations on database'
	@echo '    migrate-down       Rollbacks migrations on database'
	@echo ''

clean:
	@echo "[cleaning]"
	@go clean
	@if [ -f ${BINARY_PATH} ] ; then rm -v ${BINARY_PATH} ; fi
	@rm -rfv ${GOBIN}

download-deps:
	@echo "[download dependencies]"
	@go mod download
	@go mod download -modfile=tools/go.mod

tidy:
	@echo "[tidying]"
	@go mod tidy

lint:
	@echo "[lint]"
	@go install -modfile=tools/go.mod github.com/golangci/golangci-lint/cmd/golangci-lint
	@${GOBIN}/golangci-lint run

build:
	@echo "[build]"
	@go build -a -o ${BINARY_PATH}
	@echo "Compiled successfully!"
	@echo "Output directory: ${GOBIN}"

migrate-up:
	@echo "[migrate db up]"
	@make build
	@${BINARY_PATH} psql up --postgresql_password secretpassword --schema public

migrate-down:
	@echo "[migrate db down]"
	@make build
	@${BINARY_PATH} psql down --postgresql_password secretpassword --schema public
