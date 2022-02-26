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

	other := controllers.NewWebSocketController(controllers.NewDiceController())
	r.GET("/dice/:name/ws", other.WebSocket)

	rock := controllers.NewWebSocketController(controllers.NewRockController())
	r.GET("/rock/:name/ws", rock.WebSocket)

	r.POST("/rock/:name/reset", rock.Reset)
	// r.GET("/rock/:name/close", rock.Close)

	return r
}
