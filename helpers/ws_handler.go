package helpers

import (
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

var (
	upgrader      = websocket.Upgrader{ReadBufferSize: 1024, WriteBufferSize: 1024}
	clients       = make(map[*websocket.Conn]bool)
	clientsMutex  = sync.Mutex{}
	currentTimeMu = sync.Mutex{}
)

func WSHandler(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}
	defer conn.Close()

	// Add client to the clients map
	clientsMutex.Lock()
	clients[conn] = true
	clientsMutex.Unlock()

	// Send initial time to the new client
	sendCurrentTimeToClient(conn)

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

func sendCurrentTimeToClient(conn *websocket.Conn) {
	currentTimeMu.Lock()
	currentTime := time.Now().Format("2006-01-02 15:04:05")
	currentTimeMu.Unlock()

	sendMessage(conn, currentTime)
}

func BroadcastCurrentTime() {
	for {
		currentTimeMu.Lock()
		currentTime := time.Now().Format("2006-01-02 15:04:05")
		currentTimeMu.Unlock()

		clientsMutex.Lock()
		for client := range clients {
			sendMessage(client, currentTime)
		}
		clientsMutex.Unlock()

		time.Sleep(time.Second)
	}
}

func sendMessage(conn *websocket.Conn, message string) {
	if err := conn.WriteMessage(websocket.TextMessage, []byte(message)); err != nil {
		log.Println(err)
		removeClient(conn)
	}
}
