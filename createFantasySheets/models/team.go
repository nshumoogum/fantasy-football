package models

// Team ...
type Team struct {
	EntryHistory *EntryHistory `json:"entry_history"`
}

// EntryHistory ...
type EntryHistory struct {
	Points             int `json:"points"`
	EventTransfers     int `json:"event_transfers"`
	EventTransfersCost int `json:"event_transfers_cost"`
}
