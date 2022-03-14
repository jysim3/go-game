package controllers

import (
	"log"
	"net/http"

	//"jysim/game/models"

	"github.com/gin-gonic/gin"
	"gopkg.in/olahol/melody.v1"
)

type RoomControllerInterface interface {
	handleDisconnect(s *melody.Session)
	handleConnect(s *melody.Session)
	handleMessage(s *melody.Session, msg []byte)
	HandleRequest(c *gin.Context)
	Close()
}

type WebSocketController struct {
	m map[string]RoomControllerInterface
	n func() RoomControllerInterface
}

func (h WebSocketController) Close(c *gin.Context) {
	if val, ok := h.m[c.Param("name")]; ok {
		val.Close()
		h.m[c.Param("name")] = h.n()
	}
}

func (h WebSocketController) Reset(c *gin.Context) {
	if val, ok := h.m[c.Param("name")]; ok {
		val.Close()
		dice := h.n()
		h.m[c.Param("name")] = dice
		// dice.HandleRequest(c)
		c.String(http.StatusOK, "hello")
		return
	}
	c.String(http.StatusOK, "true")
}

func (h WebSocketController) WebSocket(c *gin.Context) {
	if val, ok := h.m[c.Param("name")]; ok {
		log.Printf("New request for existing room %s", c.Param("name"))
		val.HandleRequest(c)
	} else {
		dice := h.n()
		h.m[c.Param("name")] = dice
		dice.HandleRequest(c)
	}
}

func NewWebSocketController(n func() RoomControllerInterface) *WebSocketController {
	h := new(WebSocketController)
	h.m = make(map[string]RoomControllerInterface)
	h.n = n
	return h
}
