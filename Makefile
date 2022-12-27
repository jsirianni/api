.PHONY: install-tools
install-tools:
	cd internal/tools && go install github.com/mgechev/revive
	cd internal/tools && go install github.com/securego/gosec/v2/cmd/gosec
	cd internal/tools && go install golang.org/x/tools/cmd/goimports
	cd internal/tools && go install honnef.co/go/tools/cmd/staticcheck
	cd internal/tools && go install github.com/uw-labs/lichen

.PHONY: gomoddownload
gomoddownload:
	go mod download

.PHONY: test
test:
	go test ./...

.PHONY: vet
vet:
	go vet ./...

.PHONY: tidy
tidy:
	@go mod tidy
	@cd internal/tools && go mod tidy

.PHONY: lint
lint:
	revive -config .revive.toml -formatter friendly ./...

.PHONY: gosec
gosec:
	gosec -exclude-dir internal/tools ./...

.PHONY: staticcheck
staticcheck:
	staticcheck ./...
