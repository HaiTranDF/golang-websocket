package main

import (
	"log"
	"net/http"
	"websocket/helpers"

	"github.com/gin-gonic/gin"
	"github.com/pkg/browser"
)

func main() {
	router := gin.Default()

	go func() {
		browser.OpenURL("http://localhost:8080")
	}()

	// CREATE ENDPOIND FOR CONNECT WEBSOCKET
	router.GET("/ws", func(c *gin.Context) {
		helpers.WSHandler(c.Writer, c.Request)
	})

	// SEND YOU HTML FILE FOR OPEN TO BROWSER
	router.GET("/", func(c *gin.Context) {
		http.ServeFile(c.Writer, c.Request, "index.html")
	})

	go helpers.BroadcastCurrentTime()

	log.Println("Listening...")
	router.Run(":8080")
}
