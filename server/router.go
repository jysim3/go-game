package server

import (
	"net/http"

	"jysim/game/controllers"

	"github.com/gin-gonic/gin"
)

func NewRouter() *gin.Engine {
	r := gin.Default()

	r.GET("/channel/:name", func(c *gin.Context) {
		http.ServeFile(c.Writer, c.Request, "chan.html")
	})

	room := new(controllers.RoomController).New()
	r.GET("/channel/:name/ws", room.WebSocket)

	r.GET("/", func(c *gin.Context) {
		http.ServeFile(c.Writer, c.Request, "www/build/index.html")
	})
	r.Static("/static", "./www/build/static")

	dice := new(controllers.DiceRouter).New()
	r.GET("/dice/:name/ws", dice.WebSocket)

	return r
}
