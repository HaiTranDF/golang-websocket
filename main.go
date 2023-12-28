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
	router.GET("/ws", helpers.WSHandler)

	// SEND YOU HTML FILE FOR OPEN TO BROWSER
	router.GET("/", func(c *gin.Context) {
		http.ServeFile(c.Writer, c.Request, "index.html")
	})

	// Goroutine để gửi tin nhắn realtime
	go helpers.BroadcastCurrentMessage()

	// Goroutine để cập nhật status mỗi 5 giây
	go helpers.UpdateStatusRoutine()

	log.Println("Listening...")
	router.Run(":8080")
}
