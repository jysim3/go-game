package controllers

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"

	"jysim/game/models"

	"github.com/gin-gonic/gin"
	"gopkg.in/olahol/melody.v1"
)

type RoomController struct {
	m *melody.Melody
}

func (h RoomController) write(name string, transaction models.Transaction) {
	go func() {
		content, _ := ioutil.ReadFile("room/" + name + ".json")
		var list []interface{}
		if err := json.Unmarshal(content, &list); err != nil {

		}
		list = append(list, transaction)
		content, _ = json.Marshal(list)
		ioutil.WriteFile("room/"+name+".json", content, 0644)

	}()
	content, _ := json.Marshal(transaction)
	h.m.BroadcastFilter([]byte(content), func(q *melody.Session) bool {
		session_game, _ := q.Get("name")
		return session_game == name
	})
}

func (_ RoomController) New() RoomController {
	h := RoomController{}

	m := melody.New()
	h.m = m
	m.HandleConnect(func(s *melody.Session) {
		name, _ := s.Get("name")
		content, _ := ioutil.ReadFile("room/" + name.(string) + ".json")
		var list []interface{}
		if err := json.Unmarshal(content, &list); err != nil {
		}
		for _, x := range list {
			str, _ := json.Marshal(x)
			s.Write([]byte(str))
		}
	})

	m.HandleMessage(func(s *melody.Session, msg []byte) {
		name, _ := s.Get("name")
		var x models.Transaction
		if err := json.Unmarshal(msg, &x); err != nil {
			log.Print(err)
			return
		}
		log.Print(x)
		h.write(name.(string), x)
	})
	return h
}

func (h RoomController) Status(c *gin.Context) {
	c.String(http.StatusOK, "Working!")
}

func (h RoomController) WebSocket(c *gin.Context) {
	h.m.HandleRequestWithKeys(c.Writer, c.Request, map[string]interface{}{"name": c.Param("name")})
}
