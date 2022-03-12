package gameState

type GameState interface {
	Move(player_id int, move string) error
	GetGameStateForPlayer(player_id int) State
	GetGameStatus() interface{}
	// GetGameState() map[int]State
}

type State struct {
	Data   interface{} `json:"data"`
	Status string      `json:"status"`
}
