NAME		:= i2
PACKAGE    	:= github.com/ipedrazas/$(NAME)
VERSION		?= v0.1.0
GIT_REV    ?= $(shell git rev-parse --short HEAD)
LDFLAGS		?= "-X '${PACKAGE}/pkg/api.Version=${VERSION}' -X '${PACKAGE}/pkg/api.BuildDate=`date +%FT%T%z`' -X '${PACKAGE}/pkg/api.GitCommit=${GIT_REV}'"
REPO_NAME	:= ipedrazas
BINARY_NAME := fl

.PHONY: test
test:
	@echo "Running tests..."
	@go test -v ./...

.PHONY: build
# build: gen-docs test
build:
	@echo "Building..."
	@go build -o dist/ -ldflags=${LDFLAGS} ./...


.PHONY: api
api: build
	@echo "Running api..."
	@./dist/flex api

.PHONY: install
install: build
	@echo "Installing..."
	@go install -ldflags=${LDFLAGS} ./...
	@mv ${GOPATH}/bin/${NAME} ${GOPATH}/bin/${BINARY_NAME}

.PHONY: dist
dist:
	GOARCH=amd64 GOOS=darwin go build -ldflags=${LDFLAGS} -o dist/darwin/${BINARY_NAME} main.go
	GOARCH=amd64 GOOS=linux go build -ldflags=${LDFLAGS} -o dist/linux/${BINARY_NAME} main.go
	GOARCH=amd64 GOOS=windows go build -ldflags=${LDFLAGS} -o dist/windows/${BINARY_NAME} main.go

.PHONY: docker
docker:
	@echo "docker build -t ${REPO_NAME}/${NAME}:${VERSION}"
	@docker build --platform linux/amd64,linux/arm64 --build-arg TAG=$(VERSION) -t ${REPO_NAME}/${NAME}:${VERSION} . 
	@docker tag ${REPO_NAME}/${NAME}:${VERSION} ${REPO_NAME}/${NAME}:latest
	@docker push ${REPO_NAME}/${NAME}:${VERSION}
	@docker push ${REPO_NAME}/${NAME}:latest	

.PHONY: dbc
dbc:
	@echo "docker buildx build --platform linux/amd64,linux/arm64 --builder cloud-docker-labs -t ${REPO_NAME}/${NAME}:${VERSION} --attest type=sbom,generator=docker/scout-sbom-indexer:latest . --push"
	@docker buildx build --platform linux/amd64,linux/arm64 --builder cloud-docker-labs -t ${REPO_NAME}/${NAME}:${VERSION} \
	--attest type=sbom,generator=docker/scout-sbom-indexer:latest . --push
	@docker pull ${REPO_NAME}/${NAME}:${VERSION}
	


.PHONY: cover
cover:
	@go clean -testcache
	@go test ./... -coverprofile=coverage.out
	@go tool cover -html=coverage.out

.PHONY: tidy
tidy:
	@go fmt ./...
	@go mod tidy -v

.PHONY: gen-docs
gen-docs:
	@echo "Generating docs..."
	@swag init -g pkg/api/api.go -o pkg/docs