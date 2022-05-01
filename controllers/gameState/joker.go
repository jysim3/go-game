package gameState

import (
	"errors"
	"strings"
)

// use 'x' to move special
// use 'n' to move normal
type JokerState struct {
	history []string
}

func (j *JokerState) Move(player_id int, move string) error {
	/*
		if !j.isMoveValid(player_id, move) {
			return errors.New("Invalid Move")
		}
	*/
	if player_id > 1 {
		return errors.New("Invalid Player")
	}
	if len(j.history) == 0 && player_id != 0 {
		return errors.New("King starts first")
	}
	if len(j.history) > 0 && j.getPlayerFromMove(j.history[len(j.history)-1]) == player_id {
		return errors.New("Please wait for other player")
	}
	move = j.convertToPlayer(player_id, move)
	j.history = append(j.history, move)
	return nil
}
func (j JokerState) GetGameStateForPlayer(player_id int) State {
	opponent, player := 0, 0
	for _, move := range j.history {
		if player_id == j.getPlayerFromMove(move) {
			player += 1
		} else {
			opponent += 1
		}
	}
	ret := map[string]interface{}{}
	if player_id == 0 {
		ret["player"] = "king"
	} else {
		ret["player"] = "joker"
	}
	ret["player_past_moves"] = player
	ret["opponent_past_moves"] = opponent
	history := j.getFormattedHistory()
	if len(history) != 0 {
		if player_id == 1 {
			last_element := history[len(history)-1]
			if len(last_element) == 1 {
				history = history[:len(history)-1]
			}
		}
	}
	ret["history"] = history
	var status string

	switch s := j.GetGameStatus(); s {
	case -1:
		if (opponent == player) != (player_id != 0) {
			status = "ready"
		} else {
			status = "waiting"
		}
	case player_id:
		status = "won"
	default:
		status = "lost"
	}
	return State{
		Data:   ret,
		Status: status,
	}
}

func (j JokerState) convertToPlayer(player_id int, move string) string {
	if player_id == 0 {
		return strings.ToUpper(move)
	} else {
		return strings.ToLower(move)
	}
}
func (j JokerState) getPlayerFromMove(move string) int {
	if move == strings.ToUpper(move) {
		return 0
	} else {
		return 1
	}
}
func (j JokerState) isMoveValid(player_id int, move string) bool {
	return move == j.convertToPlayer(player_id, move)
}

func (j JokerState) getFormattedHistory() []string {
	if len(j.history) == 0 {
		return make([]string, 0)
	}
	cur_turn, history := "", make([]string, 0)
	for _, move := range j.history {
		if len(cur_turn) < 2 {
			cur_turn += move
		} else {
			history = append(history, cur_turn)
			cur_turn = move
		}
	}
	return append(history, cur_turn)
}

func (j JokerState) GetGameStatus() interface{} {
	cur_turn := ""
	end := false
	for _, move := range j.history {
		cur_turn += move
		if j.convertToPlayer(1, move) == "x" {
			end = true
		}
		if len(cur_turn) >= 2 {
			if end {
				if cur_turn == "Xx" {
					return 1
				} else {
					return 0
				}
			}
			cur_turn = ""
		}
	}
	return -1
}
