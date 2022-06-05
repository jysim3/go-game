package controllers

import (
	"log"
	"net/http"
  "time"


	//"jysim/game/models"

	"github.com/gin-gonic/gin"
	"gopkg.in/olahol/melody.v1"
)

type RoomControllerInterface interface {
	handleDisconnect(s *melody.Session)
	handleConnect(s *melody.Session)
	handleMessage(s *melody.Session, msg []byte)
	HandleRequest(c *gin.Context)
GetCount() int
	Close()
}

type WebSocketController struct {
	m map[string]RoomControllerInterface
	newController func() RoomControllerInterface
}

func (h WebSocketController) Close(c *gin.Context) {
	if val, ok := h.m[c.Param("name")]; ok {
		val.Close()
		h.m[c.Param("name")] = h.newController()
	}
}

func (h WebSocketController) Reset(c *gin.Context) {
	if val, ok := h.m[c.Param("name")]; ok {
		val.Close()
		dice := h.newController()
		h.m[c.Param("name")] = dice
		// dice.HandleRequest(c)
		c.String(http.StatusOK, "hello")
		return
	}
	c.String(http.StatusOK, "true")
}

func (h WebSocketController) cleanRoom() {
      for range time.Tick(time.Second * 60) {

        for name, session := range h.m {
          if (session.GetCount() == 0) {
            session.Close()
            delete(h.m, name)
          }
        }
        if (len(h.m) == 0) {
          break
        }
      }
}

func (h WebSocketController) WebSocket(c *gin.Context) {
	if val, ok := h.m[c.Param("name")]; ok {
		log.Printf("New request for existing room %s", c.Param("name"))
		val.HandleRequest(c)
	} else {
		dice := h.newController()
      if (len(h.m) == 0) {
       go h.cleanRoom()
      }
		h.m[c.Param("name")] = dice
		dice.HandleRequest(c)
	}
}

func (h WebSocketController) Summary() map[string]int {
  roomSessions := make(map[string]int)
  for name, session := range h.m {
    roomSessions[name] = session.GetCount()
  }
return roomSessions
}

func NewWebSocketController(newController func() RoomControllerInterface) *WebSocketController {
	h := new(WebSocketController)
	h.m = make(map[string]RoomControllerInterface)
	h.newController = newController
	return h
}
