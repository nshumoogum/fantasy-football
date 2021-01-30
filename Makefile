SHELL=bash

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

.PHONEY: script-build
script-build:
	@mkdir -p $(BUILD)/$(BIN_DIR)
	go build -o $(BUILD)/$(BIN_DIR)/$(CREATE_FANTASY_SHEETS) cmd/$(CREATE_FANTASY_SHEETS)/main.go

.PHONEY: api-build
api-build:
	@mkdir -p $(BUILD)/$(BIN_DIR)
	go build -o $(BUILD)/$(BIN_DIR)/$(FANTASY_FOOTBAL_API) cmd/$(FANTASY_FOOTBAL_API)/main.go

.PHONEY: script
script: script-build
	HUMAN_LOG=1 go run -race cmd/$(CREATE_FANTASY_SHEETS)/main.go -event-week=$(EVENT) -league-id=$(LID) -file-extension=$(EXTENSION) -url=$(URL)

.PHONEY: api
api: script-build
	HUMAN_LOG=1 go run -race cmd/$(FANTASY_FOOTBAL_API)/main.go

.PHONEY: test
test:
	go test -race -cover ./...

