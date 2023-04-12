.DEFAULT_GOAL := help
SHELL = /bin/sh

# GIT SPECIFICS

GIT_HOOKS = .git/hooks/commit-msg .git/hooks/pre-commit .git/hooks/pre-push .git/hooks/prepare-commit-msg

$(GIT_HOOKS): .git/hooks/%: .githooks/%

.githooks/%:
	touch $@

.git/hooks/%:
	cp $< $@

.PHONY: remove-git-configs
remove-git-configs: ## Remove Git Configs
	echo 'remove-git-configs'

.PHONY: add-git-configs
add-git-configs: remove-git-configs ## Add Git Configs
	git config --global branch.autosetuprebase always
	git config --global color.branch true
	git config --global color.diff true
	git config --global color.interactive true
	git config --global color.status true
	git config --global color.ui true
	git config --global commit.gpgsign true
	git config --global core.autocrlf input
	git config --global core.editor 'code --wait'
	git config --global diff.tool code
	git config --global difftool.code.cmd 'code --diff $$LOCAL $$REMOTE --wait'
	git config --global gpg.program gpg
	git config --global before.defaultbranch main
	git config --global log.date relative
	git config --global merge.tool code
	git config --global mergetool.code.cmd 'code --wait $$MERGED'
	git config --global pull.default current
	git config --global pull.rebase true
	git config --global push.autoSetupRemote true
	git config --global push.default current
	git config --global rebase.autostash true
	git config --global rerere.enabled true
	git config --global stash.showpatch true
	git config --global tag.gpgsign true

.PHONY: remove-git-hooks
remove-git-hooks: ## Remove Git Hooks
	rm -fr $(GIT_HOOKS)

.PHONY: add-git-hooks
add-git-hooks: remove-git-hooks $(GIT_HOOKS) ## Add Git Hooks

.PHONY: remove-git
remove-git: remove-git-configs remove-git-hooks ## Remove Git Configs & Hooks

.PHONY: add-git
add-git: add-git-configs add-git-hooks ## Add Git Configs & Hooks

.PHONY: help
help: ## Help
	@grep --extended-regexp '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) \
  | sort \
  | awk 'BEGIN { FS = ":.*?## " }; { printf "\033[36m%-33s\033[0m %s\n", $$1, $$2 }'

# LANGUAGE SPECIFICS

GO := $(shell go env GOPATH)
# GO_FILES := $(shell find . -name '*.go' | grep --invert-match /vendor/ | grep --invert-match _test.go)
# PROJECT_NAME := server-services-next
# PKG := gitlab.playpod.ir/alpha/backend/$(PROJECT_NAME)
# PKG_LIST := $(shell go list $(PKG)/... | grep --invert-match /vendor/)
PKG_LIST := ./...
TAGS := unit

CI := false
MIGRATE_DSN := cockroachdb://root@127.0.0.1:26257/defaultdb?sslmode=disable
MIGRATE_NAME := migrate_name
MIGRATE_VER := development
ifeq ("$(CI)", "true")
	MIGRATE_SOURCE := gitlab://$(GITLAB_USER):$(GITLAB_TOKEN)@$(GITLAB_URL)/$(PROJECT_ID)/$(PROJECT_PATH)/$(MIGRATE_VER)\#$(PROJECT_REF)
	MIGRATE_TAGS := cockroachdb gitlab
else
	MIGRATE_SOURCE := file://db/migration/$(MIGRATE_VER)
	MIGRATE_TAGS := cockroachdb
endif
OUTPUT := build
PUML := config
REPORTER := local

.PHONY: build
build: ## Build
	mkdir -p $(OUTPUT)
	go mod vendor
	CGO_ENABLED=0 go build \
	-ldflags="-s -w" \
  -a \
  -buildmode=exe \
  -mod=vendor \
  -trimpath \
  -o $(OUTPUT)

.PHONY: commitlint
commitlint: ## Commit Lint
	go install github.com/conventionalcommit/commitlint@v0.10.1
	$(GO)/bin/commitlint lint

.PHONY: coverage
coverage: test ## Coverage
	rm -fr sock
	go tool cover -func=profile/cover.txt
	go tool cover -func=profile/cover.txt -o=profile/cover.txt
	go install gotest.tools/gotestsum@v1.9.0
	$(GO)/bin/gotestsum --junitfile=profile/cover.xml -- -count=1 -covermode=set -tags="$(TAGS)" $(PKG_LIST)

.PHONY: doc
doc: ## Documentation
	go mod vendor
	go doc -all .

.PHONY: fieldalignment
fieldalignment: ## Field Alignment
	go install golang.org/x/tools/go/analysis/passes/fieldalignment/cmd/fieldalignment@v0.7.0
	$(GO)/bin/fieldalignment -fix -tags="$(TAGS)" $(PKG_LIST)

.PHONY: fix
fix: ## Fix
	go fix $(PKG_LIST)

.PHONY: fmt
fmt: ## Format
	go fmt $(PKG_LIST)

.PHONY: generate
generate: ## Generate
	go generate $(PKG_LIST)

.PHONY: git-chglog
git-chglog: ## Git Changelog
	go install github.com/git-chglog/git-chglog/cmd/git-chglog@v0.15.4
	$(GO)/bin/git-chglog --config=.chglog/config.yaml --output=CHANGELOG.md $(COMMIT_TAG)

.PHONY: golangci-lint
golangci-lint: ## GolangCI Lint
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@v1.52.2
	$(GO)/bin/golangci-lint run --build-tags=$(TAGS) --config=./.golangci.yaml $(PKG_LIST)

.PHONY: golines
golines: ## GoLines
	go install github.com/segmentio/golines@v0.11.0
	$(GO)/bin/golines --reformat-tags --write-output .

.PHONY: migrate-create
migrate-create: ## Migrate Create
	go install -tags "$(MIGRATE_TAGS)" github.com/golang-migrate/migrate/v4/cmd/migrate@v4.15.2
	$(GO)/bin/migrate -database="$(MIGRATE_DSN)" -source="$(MIGRATE_SOURCE)" -verbose create -dir="db/migration/$(MIGRATE_VER)" -ext=sql -seq "$(MIGRATE_NAME)"

.PHONY: migrate-down
migrate-down: ## Migrate Down
	go install -tags "$(MIGRATE_TAGS)" github.com/golang-migrate/migrate/v4/cmd/migrate@v4.15.2
	$(GO)/bin/migrate -database="$(MIGRATE_DSN)" -source="$(MIGRATE_SOURCE)" -verbose down 1

.PHONY: migrate-drop
migrate-drop: ## Migrate Drop
	go install -tags "$(MIGRATE_TAGS)" github.com/golang-migrate/migrate/v4/cmd/migrate@v4.15.2
	$(GO)/bin/migrate -database="$(MIGRATE_DSN)" -source="$(MIGRATE_SOURCE)" -verbose drop -f

.PHONY: migrate-up
migrate-up: ## Migrate Up
	go install -tags "$(MIGRATE_TAGS)" github.com/golang-migrate/migrate/v4/cmd/migrate@v4.15.2
	$(GO)/bin/migrate -database="$(MIGRATE_DSN)" -source="$(MIGRATE_SOURCE)" -verbose up

.PHONY: migrate-version
migrate-version: ## Migrate Version
	go install -tags "$(MIGRATE_TAGS)" github.com/golang-migrate/migrate/v4/cmd/migrate@v4.15.2
	$(GO)/bin/migrate -database="$(MIGRATE_DSN)" -source="$(MIGRATE_SOURCE)" -verbose version

.PHONY: reviewdog
reviewdog: ## Review Dog
	go install github.com/reviewdog/reviewdog/cmd/reviewdog@v0.14.1
	$(GO)/bin/reviewdog -conf=.reviewdog.yaml -fail-on-error=true -filter-mode=nofilter -reporter="$(REPORTER)"

.PHONY: shfmt
shfmt: ## Shell Formatter
	go install mvdan.cc/sh/v3/cmd/shfmt@v3.6.0
	$(GO)/bin/shfmt --case-indent --indent=2 --write script/*.sh

.PHONY: swagger
swagger: ## Swagger
	mkdir -p doc
	go install github.com/go-swagger/go-swagger/cmd/swagger@latest
	$(GO)/bin/swagger generate spec --output=doc/swagger.json

.PHONY: test
test: ## Test
	rm -fr sock
	mkdir -p profile
	go mod vendor
	go test $(PKG_LIST) -count=1 -covermode=set -coverprofile=profile/cover.txt -tags="$(TAGS)"

.PHONY: stringer
stringer: ## Stringer
	go install golang.org/x/tools/cmd/stringer@v0.7.0
	$(GO)/bin/stringer -output=./object/const_enum_string.go -type=ESRBType,GalleryType,TagType ./object
