package controllers

import (
	"encoding/json"
	"jysim/game/controllers/gameState"
	"jysim/game/models"
	"log"
	"sync"

	"github.com/gin-gonic/gin"
	"gopkg.in/olahol/melody.v1"
)

type Joker struct {
	lock     *sync.Mutex
	m        *melody.Melody
	state    gameState.GameState
	sessions map[*melody.Session]int
}

func (h *Joker) handleDisconnect(s *melody.Session) {
	h.lock.Lock()
	defer h.lock.Unlock()
	delete(h.sessions, s)
	i := 0
	for session, _ := range h.sessions {
		h.sessions[session] = i
		i += 1
	}
}

func (h *Joker) updatePlayers() {
	for session, id := range h.sessions {
		state := h.state.GetGameStateForPlayer(id)
		if r, err := json.Marshal(state); err == nil {
			session.Write(r)
		} else {
		}
	}
}

func (h *Joker) handleConnect(s *melody.Session) {
	h.lock.Lock()
	defer h.lock.Unlock()
	if len(h.sessions) >= 2 {
		s.CloseWithMsg(melody.FormatCloseMessage(1000, "SERVER_FULL"))
		return
	}
	h.sessions[s] = 0
	i := 0
	for session := range h.sessions {
		h.sessions[session] = i
		i += 1
	}
	if len(h.sessions) == 2 {
		h.updatePlayers()
	}
}

func (h *Joker) handleMessage(s *melody.Session, msg []byte) {
	// name, _ := s.Get("name")
	var command models.Command
	if err := json.Unmarshal(msg, &command); err != nil {
		return
	}
	if command.Code == "send" {
		err := h.state.Move(h.sessions[s], command.Data.(string))
		log.Print(err)
		h.updatePlayers()
	} else if command.Code == "start" {
		if h.state.GetGameStatus() != -1 {
			h.lock.Lock()
			defer h.lock.Unlock()
			for session, id := range h.sessions {
				h.sessions[session] = (id + 1) % 2
			}
		}
		h.state = new(gameState.JokerState)
		h.updatePlayers()
	}
}

func (h *Joker) Close() {
	h.m.Close()
}

func (h *Joker) HandleRequest(c *gin.Context) {
	h.m.HandleRequest(c.Writer, c.Request)
}

func NewJokerController() func() RoomControllerInterface {
	return func() RoomControllerInterface {
		h := Joker{}
		m := melody.New()
		h.lock = new(sync.Mutex)
		h.m = m
		h.state = new(gameState.JokerState)
		h.sessions = make(map[*melody.Session]int)
		m.HandleConnect(h.handleConnect)
		m.HandleDisconnect(h.handleDisconnect)
		m.HandleMessage(h.handleMessage)
		return &h
	}
}
