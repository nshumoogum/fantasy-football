package csv

import (
	"context"
	"os"
	"strconv"

	"github.com/nshumoogum/fantasy-football/createFantasySheets/models"
)

// CreateTeam ...
func CreateTeam(ctx context.Context, connection *os.File, filename string, team *models.Result, teamData *models.EntryHistory) (err error) {
	teamLine := strconv.Itoa(team.Entry) + "," +
		team.PlayerName + "," +
		team.TeamName + "," +
		strconv.Itoa(teamData.EventTransfers) + "," +
		strconv.Itoa(teamData.EventTransfersCost) + "," +
		strconv.Itoa(teamData.Points)

	if err = writeToFile(ctx, connection, filename, teamLine); err != nil {
		return
	}

	return
}
