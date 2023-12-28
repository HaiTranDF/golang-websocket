package helpers

import (
	"log"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

// Message struct lưu trữ thông tin giờ và status
type Message struct {
	Time   string `json:"time"`
	Status bool   `json:"status"`
}

var (
	upgrader     = websocket.Upgrader{ReadBufferSize: 1024, WriteBufferSize: 1024}
	clients      = make(map[*websocket.Conn]bool)
	clientsMutex = sync.Mutex{}
	messageMu    = sync.Mutex{}
	status       = false
)

func WSHandler(c *gin.Context) {
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Println(err)
		return
	}
	defer conn.Close()

	// Add client to the clients map
	clientsMutex.Lock()
	clients[conn] = true
	clientsMutex.Unlock()

	// Send initial message to the new client
	sendMessage(conn, getCurrentMessageWithoutStatus())

	for {
		// Read message from the client (if needed)
		_, _, err := conn.ReadMessage()
		if err != nil {
			// Remove client when connection is closed
			removeClient(conn)
			break
		}
	}
}

func removeClient(conn *websocket.Conn) {
	clientsMutex.Lock()
	delete(clients, conn)
	clientsMutex.Unlock()
}

func getCurrentMessageWithoutStatus() Message {
	messageMu.Lock()
	defer messageMu.Unlock()

	currentTime := time.Now().Format("2006-01-02 15:04:05")

	return Message{
		Time:   currentTime,
		Status: status,
	}
}

func BroadcastCurrentMessage() {
	for {
		message := getCurrentMessageWithoutStatus()

		clientsMutex.Lock()
		for client := range clients {
			sendMessage(client, message)
		}
		clientsMutex.Unlock()

		// Ngủ 1 giây để giữ nguyên giờ hiện tại mà không cập nhật status liên tục
		time.Sleep(time.Second)
	}
}

func sendMessage(conn *websocket.Conn, message Message) {
	messageMu.Lock()
	defer messageMu.Unlock()

	if err := conn.WriteJSON(message); err != nil {
		log.Println(err)
		removeClient(conn)
	}
}

func UpdateStatusRoutine() {
	for {
		// Cập nhật status mỗi 5 giây
		time.Sleep(5 * time.Second)
		updateStatus()
	}
}

func updateStatus() {
	messageMu.Lock()
	defer messageMu.Unlock()

	// Đảo ngược giá trị của status
	status = !status
}
