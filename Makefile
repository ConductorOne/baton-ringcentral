CONNECTOR_NAME := baton-ringcentral

ifeq ($(BATON_LAMBDA_SUPPORT),true)
	BUILD_TAGS := -tags baton_lambda_support
endif

.PHONY: build
build:
	go build ${BUILD_TAGS} -o bin/${CONNECTOR_NAME} ./cmd/${CONNECTOR_NAME}

.PHONY: update-deps
update-deps:
	go get -d -u ./...
	go mod tidy -v

.PHONY: add-dep
add-dep:
	go get -d -u "$v"

.PHONY: lint
lint:
	golangci-lint run

.PHONY: generate
generate:
	go generate ./...
