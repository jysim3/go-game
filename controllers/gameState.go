package controllers

type GameState interface {
	Move(player_id string, command string)
	GetUpdate(player_id string) interface{}
	GetAllUpdate() map[string]interface{}
}
