package helpers

import (
	"sync"

	"github.com/gorilla/websocket"
)

var connectionMutex = &sync.Mutex{}
var userConnections = make(map[string]*websocket.Conn)

func AddConnection(userID string, conn *websocket.Conn) {
	connectionMutex.Lock()
	defer connectionMutex.Unlock()
	userConnections[userID] = conn
}

func RemoveConnection(userID string) {
	connectionMutex.Lock()
	defer connectionMutex.Unlock()
	delete(userConnections, userID)
}

func GetConnection(userID string) (*websocket.Conn, bool) {
	connectionMutex.Lock()
	defer connectionMutex.Unlock()
	conn, ok := userConnections[userID]
	return conn, ok
}
