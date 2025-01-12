package helpers

import (
	"fmt"
	"sync"

	"github.com/gorilla/websocket"
)

// Connection represents a user connection in a room
type Room struct {
	Connections map[int64]*websocket.Conn // Map from userID to websocket connection
}

type Location struct {
	sync.Mutex
	Rooms map[int64]*Room // Map from roomID to Room
}

func NewLocation() *Location {
	return &Location{Rooms: make(map[int64]*Room)}
}

// AddUserToRoom adds a user and their connection to a room
func (l *Location) AddUserToRoom(roomID int64, userID int64, conn *websocket.Conn) {
	l.Lock()
	defer l.Unlock()

	// Check if the room exists, if not create it
	if _, exists := l.Rooms[roomID]; !exists {
		l.Rooms[roomID] = &Room{Connections: make(map[int64]*websocket.Conn)}
	}

	// Add or update the user in the room's connections map
	l.Rooms[roomID].Connections[userID] = conn
}

// RemoveUserFromRoom removes a user and their connection from the room
func (l *Location) RemoveUserFromRoom(roomID int64, userID int64) error {
	l.Lock()
	defer l.Unlock()

	room, exists := l.Rooms[roomID]
	if !exists {
		return fmt.Errorf("room with ID %d does not exist", roomID)
	}

	// Check if the user exists in the room
	if conn, exists := room.Connections[userID]; exists {
		// Close the WebSocket connection
		conn.Close()

		// Remove the user from the room's connections map
		delete(room.Connections, userID)

		// If no users remain in the room, delete the room
		if len(room.Connections) == 0 {
			delete(l.Rooms, roomID)
		}
		return nil
	}

	return fmt.Errorf("user with ID %d not found in room %d", userID, roomID)
}

// FetchUsersInRoom fetches all users and their connections for a specific room
func (l *Location) FetchUsersInRoom(roomID int64) (map[int64]*websocket.Conn, bool) {
	l.Lock()
	defer l.Unlock()

	room, exists := l.Rooms[roomID]
	if exists {
		return room.Connections, true
	}
	return nil, false
}

// FetchUserConn fetches the websocket connection for a specific user in a room
func (l *Location) FetchUserConn(roomID int64, userID int64) (*websocket.Conn, bool) {
	l.Lock()
	defer l.Unlock()

	room, exists := l.Rooms[roomID]
	if exists {
		conn, exists := room.Connections[userID]
		return conn, exists
	}
	return nil, false
}
