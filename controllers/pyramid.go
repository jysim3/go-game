package controllers

import (
	"encoding/json"
	"log"
	"math/rand"

	"jysim/game/models"

	"github.com/gin-gonic/gin"
	"gopkg.in/olahol/melody.v1"
)
type card struct {
	number int
	suit int
}
type PyramidController struct {
	m           *melody.Melody
	sessionData map[string]*melody.Session
	cardsDeck []card
}

func dealCards(cards []card) card {
	randInt = rand.Intn(len(cards))
	card = cards[randInt]
	cards[randInt] = cards[len(cards)-1]
	cards = cards[:len(cards) - 1]
	return card
}

func (h DiceController) reset() {
	h.m.Close()
	h.m = melody.New()
}

func (h DiceController) restartGame() {
	h.card = make([]card, 56)
	suit, number := 0, 0
	for i := range h.card {
		h.card[i] = card{suit, number}
		number += 1
		if number == 13 {
			suit += 1
			number = 0
		}
	}

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
		h.restartGame()
	} else if command.Code == "open" {
		h.open()
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
	if val, err := c.Cookie("sessionId"); err == nil {
		h.m.HandleRequestWithKeys(c.Writer, c.Request, map[string]interface{}{"sessionId": val})
		return
	}
	h.m.HandleRequest(c.Writer, c.Request)
}

func NewPyramidController() func() RoomControllerInterface {
	return func() RoomControllerInterface {
		h := PyramidController{}
		m := melody.New()
		h.m = m
		m.HandleConnect(h.handleConnect)
		m.HandleDisconnect(h.handleDisconnect)
		m.HandleMessage(h.handleMessage)
		return h
	}
}
