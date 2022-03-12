package gameState

import (
	"testing"
)

func TestJoker_testGameEnd_King(t *testing.T) {
	j := JokerState{}
	j.Move(0, "N")
	j.Move(1, "n")
	j.Move(0, "X")
	j.Move(1, "n")

	state := j.GetGameStateForPlayer(0)
	if state.Status != "king won" {
		t.Errorf("Game ended with King winning but Status is %s", j.GetGameStateForPlayer(0).Status)
	}
	history := state.Data.(map[string]interface{})["history"].([]string)
	if history[0] != "Nn" {
		t.Errorf("Turn 1 expected to be Nn, but it is %s", history[0])
	}
	if history[1] != "Xn" {
		t.Errorf("Turn 1 expected to be Nx, but it is %s", history[1])
	}
}

func TestJoker_testWaiting(t *testing.T) {
	j := JokerState{}
	j.Move(0, "N")
	j.Move(1, "n")
	j.Move(0, "N")
	j.Move(1, "n")
	state := j.GetGameStateForPlayer(0)
	if state.Status != "ready" {
		t.Errorf("Game ready winning but Status is %s", state.Status)
	}
	state = j.GetGameStateForPlayer(1)
	if state.Status != "waiting" {
		t.Errorf("Game waiting for King winning but Status is %s", state.Status)
	}
	history := state.Data.(map[string]interface{})["history"].([]string)
	if history[0] != "Nn" {
		t.Errorf("Turn 1 expected to be Nn, but it is %s", history[0])
	}
}

func TestJoker_testGameEnd_JokerWin(t *testing.T) {
	j := JokerState{}
	j.Move(0, "N")
	j.Move(1, "n")
	j.Move(0, "N")
	j.Move(1, "n")
	j.Move(0, "X")
	j.Move(1, "x")
	state := j.GetGameStateForPlayer(0)
	if state.Status != "joker won" {
		t.Errorf("Game ended with joker winning but Status is %s", j.GetGameStateForPlayer(0).Status)
	}
	history := state.Data.(map[string]interface{})["history"].([]string)
	if history[0] != "Nn" {
		t.Errorf("Turn 1 expected to be Nn, but it is %s", history[0])
	}
	if history[2] != "Xx" {
		t.Errorf("Turn 3 expected to be Nx, but it is %s", history[1])
	}
}

func TestJoker_testGameEnd_Joker(t *testing.T) {
	j := JokerState{}
	j.Move(0, "N")
	j.Move(1, "n")
	j.Move(0, "N")
	j.Move(1, "x")

	state := j.GetGameStateForPlayer(0)
	if state.Status != "king won" {
		t.Errorf("Game ended with joker winning but Status is %s", j.GetGameStateForPlayer(0).Status)
	}
	history := state.Data.(map[string]interface{})["history"].([]string)
	if history[0] != "Nn" {
		t.Errorf("Turn 1 expected to be Nn, but it is %s", history[0])
	}
	if history[1] != "Nx" {
		t.Errorf("Turn 1 expected to be Nx, but it is %s", history[1])
	}
}

func TestJoker_1(t *testing.T) {
	j := JokerState{}
	j.Move(0, "N")
	j.Move(1, "n")
	if err := j.Move(1, "n"); err == nil {
		t.Error("Error expected")
	}
}
