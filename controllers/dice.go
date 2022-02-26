package controllers

import (
	"encoding/json"
	"log"
	"math/rand"

	"jysim/game/models"

	"github.com/gin-gonic/gin"
	"gopkg.in/olahol/melody.v1"
)

type DiceRouter struct {
	m map[string]*DiceController
}

type DiceController struct {
	m        *melody.Melody
	status   string
	sessions map[*melody.Session][5]int
}

func (h DiceController) handleDisconnect(s *melody.Session) {
	delete(h.sessions, s)
}

func (h DiceController) handleConnect(s *melody.Session) {
	h.sessions[s] = [5]int{1, 1, 1, 1, 1}
	/*
		content, _ := ioutil.ReadFile("room/" + name.(string) + ".json")
		var list []interface{}
		if err := json.Unmarshal(content, &list); err != nil {
		}
		for _, x := range list {
			str, _ := json.Marshal(x)
			s.Write([]byte(str))
		} */
}

func (h DiceController) reset() {
	h.sessions = make(map[*melody.Session][5]int)
}

func (h DiceController) open() {
	fullList := make([]int, 0)
	for s, _ := range h.sessions {
		list := h.sessions[s]
		fullList = append(fullList, list[:]...)
	}
	ret := models.Command{
		"start",
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
		h.sessions[s] = list
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
	}
}

func (h DiceController) HandleRequest(c *gin.Context) {
	h.m.HandleRequest(c.Writer, c.Request)
}

func (h DiceController) Close() {
	h.m.Close()
}

func NewDiceController() func() RoomControllerInterface {
	return func() RoomControllerInterface {
		h := DiceController{}
		m := melody.New()
		h.m = m
		h.sessions = make(map[*melody.Session][5]int)
		m.HandleConnect(h.handleConnect)
		m.HandleDisconnect(h.handleDisconnect)
		m.HandleMessage(h.handleMessage)
		return h
	}
}
