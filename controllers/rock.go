package controllers

import (
	"encoding/json"
	"jysim/game/models"

	"github.com/gin-gonic/gin"
	"gopkg.in/olahol/melody.v1"
)

type Rock struct {
	m        *melody.Melody
	status   string
	sessions map[*melody.Session]string
}

func (h Rock) handleDisconnect(s *melody.Session) {
	delete(h.sessions, s)
}

func (h Rock) handleConnect(s *melody.Session) {
	if len(h.sessions) >= 2 {
		s.CloseWithMsg(melody.FormatCloseMessage(1000, "SERVER_FULL"))
	} else {
		h.sessions[s] = ""
	}
}

func (h Rock) handleMessage(s *melody.Session, msg []byte) {
	// name, _ := s.Get("name")
	var command models.Command
	if err := json.Unmarshal(msg, &command); err != nil {
		return
	}
	if command.Code == "send" {
		h.sessions[s] = command.Data.(string)
		end := false
		var opponent *melody.Session
		if len(h.sessions) == 2 {
			end = true
			for session, data := range h.sessions {
				if data == "" {
					end = false
				} else if session != s {
					opponent = session
				}
			}
		}
		if end {
			ret := models.Command{
				"open",
				h.sessions[opponent],
			}
			r, _ := json.Marshal(ret)
			s.Write(r)

			ret = models.Command{
				"open",
				h.sessions[s],
			}
			r, _ = json.Marshal(ret)
			opponent.Write(r)
		}
	} else if command.Code == "start" {

		ret := models.Command{
			"start",
			"",
		}
		r, _ := json.Marshal(ret)
		for s := range h.sessions {
			h.sessions[s] = ""
			s.Write(r)
		}
	}
}

func (h Rock) Close() {
	h.m.Close()
}

func (h Rock) HandleRequest(c *gin.Context) {
	h.m.HandleRequest(c.Writer, c.Request)
}

func NewRockController() func() RoomControllerInterface {
	return func() RoomControllerInterface {
		h := Rock{}
		m := melody.New()
		h.m = m
		h.sessions = make(map[*melody.Session]string)
		m.HandleConnect(h.handleConnect)
		m.HandleDisconnect(h.handleDisconnect)
		m.HandleMessage(h.handleMessage)
		return h
	}
}
