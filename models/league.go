package models

// Resource ...
type Resource struct {
	League    *League    `json:"league"`
	Standings *Standings `json:"standings"`
}

// League ...
type League struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

// Standings ...
type Standings struct {
	HasNext bool      `json:"has_next"`
	Results []*Result `json:"results"`
}

// Result ...
type Result struct {
	Entry      int    `json:"entry"`
	PlayerName string `json:"player_name"`
	TeamName   string `json:"entry_name"`
}
