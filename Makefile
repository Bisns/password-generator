.PHONY: build list github

VERSION := "0.0.1"

VARS    := -w -s -H=windowsgui

build:
	@go build -ldflags "${VARS}" -o "bin/v${VERSION}/password-generator-${VERSION}.exe" && upx --best "bin/v${VERSION}/password-generator-${VERSION}.exe"

github:
	@go build -ldflags "${VARS}" -o "bin/password-generator-${VERSION}.exe"

list:
	@echo ${VERSION}
	@echo ${VARS}

vet:
	go vet
