Fantasy Football
================

A list of scripts to pull data from the fantasy premier league (fpl) API.

## Installation

As the script is written in Go, make sure you have version 1.13.0 or greater installed.

Using [Homebrew](https://brew.sh/) to install go
* Run `brew install go` or `brew upgrade go`

Download repository
* Copy repo using clone with ssh or https in `Code` dropdown
* On commandline run `git clone <paste copy of url>`

## Create Fantasy Sheets

This script is to generate a csv file containing data obtained from the [fantasy football website](https://fantasy.premierleague.com/) for all teams within a league during a particular round.

Once you have setup go, you can run the script using either of the following:

* Run `go run main.go -event-week={event} -league-id={league id}` the event value should be a number representing the round you are interested in. 

Or

* Setup environment variables, using `export EVENT_WEEK={integer value}` and `export LEAGUE_ID={the league ID you are interested in}`
* Then run `make script`

This shall take several seconds to complete and will generate a file (default extension to csv) containing a list of all teams of a division and the total points they received in that week. The file will contain data on the following headers for each team:

* Entry ID
* Player
* Team Name
* Number of Transfers
* Transfer Cost (the number of points fpl deduct due to the number of transfers made)
* Points

There are two other flags for the script, these are to describe the file extension `FILE_EXTENSION` and `FANTASY_LEAGUE_URL`, see Makefile for default values and how they map to defined flags in `cmd/createFanatasySheets/main.go`.
