GOFLAGS += -cover -covermode=atomic -race -shuffle=on -mod=readonly -trimpath

export GOFLAGS

LDFLAGS += -s -w -buildid=

define GO_SOURCES
$(strip
$(wildcard *.go)
$(wildcard cmd/cocopilot/*.go)
)
endef

.PHONY: all
all: bin/cocopilot gomod2nix.toml

.PHONY: check
check:
	nix flake check -L

.PHONY: clean
clean:
	-rm -rf bin

.PHONY: configure
configure:
	cog install-hook --all --overwrite

.PHONY: delete-passwords
delete-passwords:
	-security delete-generic-password -s "https://github.com" -a 01ab8ac9400c4e429b23
	-security delete-generic-password -s "https://api.githubcopilot.com" -a 01ab8ac9400c4e429b23

.PHONY: dist
dist: GOFLAGS=
dist:
	goreleaser release --clean --snapshot --skip before --release-notes CHANGELOG.md

.PHONY: release
release:
	cog bump --auto

.PHONY: test
test:
	gotestdox ${GOFLAGS} ./...

## make magic, not war ;)

bin/%: ${GO_SOURCES} go.sum
	go build -ldflags='${LDFLAGS}' -o ./$@ ./cmd/$(patsubst bin/%,%,$@)/...

go.sum: GOFLAGS-=-mod-readonly
go.sum: ${GO_SOURCES} go.mod
	@go mod tidy -v -x
	@touch $@

gomod2nix.toml: go.sum
	gomod2nix generate
