SHELL=bash

GREEN  := $(shell tput -Txterm setaf 2)
YELLOW := $(shell tput -Txterm setaf 3)
WHITE  := $(shell tput -Txterm setaf 7)
CYAN   := $(shell tput -Txterm setaf 6)
RESET  := $(shell tput -Txterm sgr0)

BUILD=build
BIN_DIR?=.

DEFAULT_EVENT=0
DEFAULT_URL=https://fantasy.premierleague.com/api
DEFAULT_EXTENSION=.csv
DEFAULT_LID=414763

EVENT = $(if $(EVENT_WEEK),$(EVENT_WEEK),$(DEFAULT_EVENT))
URL = $(if $(FANTASY_LEAGUE_URL),$(FANTASY_LEAGUE_URL),$(DEFAULT_URL))
EXTENSION = $(if $(FILE_EXTENSION),$(FILE_EXTENSION),$(DEFAULT_EXTENSION))
LID = $(if $(LEAGUE_ID),$(LEAGUE_ID),$(DEFAULT_LID))

CREATE_FANTASY_SHEETS=createFantasySheets
FANTASY_FOOTBAL_API=fantasyFootballAPI

.PHONY: all
all: delimiter-AUDIT audit delimiter-LINTERS lint delimiter-UNIT-TESTS test delimiter-FINISH ## Runs multiple targets, audit, lint and test

.PHONEY: api-build
api-build: ## Builds binary of the fpl api and stores in build directory
	@mkdir -p $(BUILD)/$(BIN_DIR)
	go build -o $(BUILD)/$(BIN_DIR)/$(FANTASY_FOOTBAL_API) cmd/$(FANTASY_FOOTBAL_API)/main.go

.PHONEY: api
api: api-build ## Builds and runs fpl api locally
	HUMAN_LOG=1 go run -race cmd/$(FANTASY_FOOTBAL_API)/main.go

.PHONY: audit
audit: ## Runs checks for security vulnerabilities on dependencies (including transient ones)
	go list -m all | nancy sleuth

.PHONY: convey
convey: ## Runs unit test suite and outputs results on http://127.0.0.1:8080/
	goconvey ./...

.PHONY: delimiter-%
delimiter-%:
	@echo '===================${GREEN} $* ${RESET}==================='

.PHONEY: lint
lint: ## Use locally to run linters against Go code
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@v1.50.1
	golangci-lint run ./...

.PHONEY: fmt
fmt: ## Run Go formatting on code
	go fmt ./...

.PHONEY: run
run:
	$(BUILD)/$(BIN_DIR)/$(FANTASY_FOOTBAL_API)

.PHONEY: script-build
script-build: ## Builds binary of the fpl script and stores in build directory
	@mkdir -p $(BUILD)/$(BIN_DIR)
	go build -o $(BUILD)/$(BIN_DIR)/$(CREATE_FANTASY_SHEETS) cmd/$(CREATE_FANTASY_SHEETS)/main.go

.PHONEY: script
script: script-build ## Runs fpl script with local configurations
	HUMAN_LOG=1 go run -race cmd/$(CREATE_FANTASY_SHEETS)/main.go -event-week=$(EVENT) -league-id=$(LID) -file-extension=$(EXTENSION) -url=$(URL)

.PHONEY: test
test: ## Runs unit tests including checks for race conditions and returns coverage
	go test -count=1 -race -cover ./...

.PHONY: help
help: ## Show help page for list of make targets
	@echo ''
	@echo 'Usage:'
	@echo '  ${YELLOW}make${RESET} ${GREEN}<target>${RESET}'
	@echo ''
	@echo 'Targets:'
	@awk 'BEGIN {FS = ":.*?## "} { \
		if (/^[a-zA-Z_-]+:.*?##.*$$/) {printf "    ${YELLOW}%-20s${GREEN}%s${RESET}\n", $$1, $$2} \
		else if (/^## .*$$/) {printf "  ${CYAN}%s${RESET}\n", substr($$1,4)} \
		}' $(MAKEFILE_LIST)