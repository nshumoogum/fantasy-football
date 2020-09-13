package csv

import (
	"context"
	"os"
	"strconv"
	"strings"

	"github.com/nshumoogum/fantasy-football/models"
)

// CreateTeam ...
func CreateTeam(ctx context.Context, connection *os.File, filename string, team *models.Result, teamData *models.EntryHistory) (err error) {
	teamLine := strconv.Itoa(team.Entry) + "," +
		removeCommas(team.PlayerName) + "," +
		removeCommas(team.TeamName) + "," +
		strconv.Itoa(teamData.EventTransfers) + "," +
		strconv.Itoa(teamData.EventTransfersCost) + "," +
		strconv.Itoa(teamData.Points)

	if err = writeToFile(ctx, connection, filename, teamLine); err != nil {
		return
	}

	return
}

func removeCommas(name string) string {
	newName := strings.Replace(name, ",", "", -1)
	return newName
}
