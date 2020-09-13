# Create Fantasy Sheets
=======================

This script is to generate a csv file containing data obtained from the [fantasy football website](https://fantasy.premierleague.com/) for all teams within a league during a particular round.

### Installation

As the script is written in Go, make sure you have version 1.13.0 or greater installed.

Using [Homebrew](https://brew.sh/) to install go
* Run `brew install go` or `brew upgrade go`

* Run `go run main.go -event-week=<event>` the event value should be a number representing the round you are interested in. This shall take several seconds to complete
