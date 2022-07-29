package server

import (
	"fmt"
	"net/http"
	"net/http/httputil"
	"os"

	"jysim/game/controllers"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func NewRouter() *gin.Engine {
	r := gin.Default()

	config := cors.DefaultConfig()
	if os.Getenv("SERVER_ENV") == "dev" {
		config.AllowOrigins = []string{"http://*.jysim3.com", "http://jysim3.com", "http://localhost:8080"}
	} else {
		config.AllowOrigins = []string{"http://*.jysim3.com", "http://jysim3.com"}
	}
	r.Use(cors.New(config))
	fmt.Println(os.Getenv("SERVER_ENV"))
	fmt.Println(config.AllowOrigins)
	r.GET("/channel/:name", func(c *gin.Context) {
		http.ServeFile(c.Writer, c.Request, "chan.html")
	})

	room := new(controllers.RoomController).New()
	r.GET("/channel/:name/ws", room.WebSocket)

	if os.Getenv("SERVER_ENV") == "dev" {

		r.GET("/", func(c *gin.Context) {

			if _, err := c.Cookie("sessionId"); err != nil {
				c.SetCookie("sessionId", uuid.New().String() /* maxAge= */, 3600 /* path= */, "/" /* domain= */, "localhost" /* httpOnly */, false, true)
			}
			director := func(req *http.Request) {
				req.URL.Scheme = "http"
				req.URL.Host = "vue:8080"

				cookie := &http.Cookie{
					Name:   "sessionId",
					Value:  uuid.New().String(),
					MaxAge: 300,
				}
				req.AddCookie(cookie)
			}
			proxy := &httputil.ReverseProxy{Director: director}
			proxy.ServeHTTP(c.Writer, c.Request)
		})

		r.GET("/static/*path", func(c *gin.Context) {

			director := func(req *http.Request) {
				req.URL.Scheme = "http"
				req.URL.Host = "vue:8080"
			}
			proxy := &httputil.ReverseProxy{Director: director}
			proxy.ServeHTTP(c.Writer, c.Request)
		})
	} else {
		r.GET("/", func(c *gin.Context) {

			if _, err := c.Cookie("sessionId"); err != nil {
				c.SetCookie("sessionId", uuid.New().String() /* maxAge= */, 3600 /* path= */, "/" /* domain= */, "game.jysim3.com" /* httpOnly */, false, true)
			}
			http.ServeFile(c.Writer, c.Request, "www/build/index.html")
		})
		r.Static("/static", "./www/build/static")
	}

	other := controllers.NewWebSocketController(controllers.NewDiceController())
	r.GET("/dice/:name/ws", other.WebSocket)

	rock := controllers.NewWebSocketController(controllers.NewRockController())
	r.GET("/rock/:name/ws", rock.WebSocket)
	r.POST("/rock/:name/reset", rock.Reset)
	// r.GET("/rock/:name/close", rock.Close)

	joker := controllers.NewWebSocketController(controllers.NewJokerController())
	r.GET("/joker/:name/ws", joker.WebSocket)
	r.POST("/joker/:name/reset", joker.Reset)

	pyramid := controllers.NewWebSocketController(controllers.NewPyramidController())
	r.GET("/pyramid/:name/ws", pyramid.WebSocket)
	r.POST("/pyramid/:name/reset", pyramid.Reset)

	r.GET("/summary", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"dice":    other.Summary(),
			"rock":    rock.Summary(),
			"joker":   joker.Summary(),
			"pyramid": pyramid.Summary(),
		})
	})
	return r
}
