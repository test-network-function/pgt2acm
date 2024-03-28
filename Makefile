GO_PACKAGES=$(shell go list ./... | grep -v vendor)
.PHONY: all clean test
.PHONY: \
	lint \
	build \
	fetch-schema \
	vet
# Runs configured linters
lint:
	golangci-lint run --timeout 10m0s
build: build-plugins
	go build -o pgt2acm
build-linux-amd64: build-plugins-linux-amd64
	GOOS=linux GOARCH=amd64 go build -o pgt2acm
build-darwin-arm64: build-plugins-darwin-arm64
	GOOS=darwin GOARCH=arm64 go build -o pgt2acm
test: build
	scripts/test.sh
clean:
	rm -r test/acmgen-output || true
	rm pgt2acm || true
	rm -rf build || true
fetch-schema:
	kustomize openapi fetch > test/cluster-schema.json
vet:
	go vet ${GO_PACKAGES}
build-plugins-darwin-arm64: clean
	GOOS=darwin GOARCH=arm64 scripts/build-plugins.sh
build-plugins-linux-amd64: clean
	GOOS=linux GOARCH=amd64 scripts/build-plugins.sh
build-plugins: clean
	scripts/build-plugins.sh