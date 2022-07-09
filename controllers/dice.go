package controllers

import (
	"encoding/json"
	"log"
	"math/rand"

	"jysim/game/models"

	"github.com/gin-gonic/gin"
	"gopkg.in/olahol/melody.v1"
)

type DiceController struct {
	m        *melody.Melody
	status   string
	sessions map[*melody.Session]map[string]interface{}
}

func (h DiceController) handleDisconnect(s *melody.Session) {
	delete(h.sessions, s)
	ret := models.Command{
		"players",
		len(h.sessions),
	}
	for s, _ := range h.sessions {
		r, _ := json.Marshal(ret)
		s.Write(r)
	}
}

func (h DiceController) handleConnect(s *melody.Session) {
	h.sessions[s] = make(map[string]interface{})
	var list [5]int
	for i, _ := range list {
		list[i] = rand.Intn(6) + 1
	}
	h.sessions[s]["dice"] = list
	ret := models.Command{
		"start",
		list,
	}
	r, _ := json.Marshal(ret)
	s.Write(r)

	ret = models.Command{
		"players",
		len(h.sessions),
	}
	for s, _ := range h.sessions {
		r, _ := json.Marshal(ret)
		s.Write(r)
	}
}

func (h DiceController) reset() {
	for s, _ := range h.sessions {
		s.Close()
	}
	h.sessions = make(map[*melody.Session]map[string]interface{})
}

func (h DiceController) open() {
	fullList := make(map[string][5]int)
	for s, _ := range h.sessions {
		list, name := h.sessions[s]["dice"].([5]int), h.sessions[s]["name"].(string)
		fullList[name] = list
	}
	ret := models.Command{
		"open",
		fullList,
	}
	r, _ := json.Marshal(ret)
	for s, _ := range h.sessions {
		s.Write(r)
	}
	h.status = ""
}

func (h DiceController) restartGame() {
	for s, _ := range h.sessions {
		var list [5]int
		for i, _ := range list {
			list[i] = rand.Intn(6) + 1
		}
		h.sessions[s]["dice"] = list
		ret := models.Command{
			"start",
			list,
		}
		r, _ := json.Marshal(ret)
		s.Write(r)
	}
	h.status = "started"
}

func (h DiceController) handleMessage(s *melody.Session, msg []byte) {
	// name, _ := s.Get("name")
	var command models.Command
	if err := json.Unmarshal(msg, &command); err != nil {
		log.Print(err)
		return
	}
	if command.Code == "reset" {
		h.reset()
	} else if command.Code == "start" {
		if h.status == "" {
			h.restartGame()
		} else {
			s.Write([]byte("Fuck off"))
		}
	} else if command.Code == "open" {
		if h.status == "" {
			h.open()
		} else {
			h.status = ""
		}
	} else if command.Code == "setName" {
		h.sessions[s]["name"] = command.Data.(string)
	} else if command.Code == "backdoor" {
		a := command.Data.([]interface{})
		var temp [5]int
		for i := range a {
			temp[i] = int(a[i].(float64))
		}
		h.sessions[s]["dice"] = temp
		ret := models.Command{
			"start",
			h.sessions[s]["dice"],
		}
		r, _ := json.Marshal(ret)
		s.Write(r)
	}
}

func (h DiceController) HandleRequest(c *gin.Context) {
	h.m.HandleRequest(c.Writer, c.Request)
}

func (h DiceController) Close() {
	h.m.Close()
}

func (h DiceController) GetCount() int {
	return len(h.sessions)
}

func NewDiceController() func() RoomControllerInterface {
	return func() RoomControllerInterface {
		h := DiceController{}
		m := melody.New()
		h.m = m
		h.sessions = make(map[*melody.Session]map[string]interface{})
		m.HandleConnect(h.handleConnect)
		m.HandleDisconnect(h.handleDisconnect)
		m.HandleMessage(h.handleMessage)
		return h
	}
}
