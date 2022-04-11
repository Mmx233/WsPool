package ws

import (
	"github.com/gorilla/websocket"
	"sync"
	"time"
)

type Conn struct {
	*websocket.Conn
	sync.Mutex
	pool *Pool
	key  any
}

func (a *Conn) WriteJSON(v interface{}) error {
	a.Lock()
	defer a.Unlock()
	return a.Conn.WriteJSON(v)
}

func (a *Conn) WriteMessage(messageType int, data []byte) error {
	a.Lock()
	defer a.Unlock()
	return a.Conn.WriteMessage(messageType, data)
}

func (a *Conn) WriteControl(messageType int, data []byte, deadline time.Time) error {
	a.Lock()
	defer a.Unlock()
	return a.Conn.WriteControl(messageType, data, deadline)
}

func (a *Conn) WritePreparedMessage(pm *websocket.PreparedMessage) error {
	a.Lock()
	defer a.Unlock()
	return a.Conn.WritePreparedMessage(pm)
}

func (a *Conn) Clear() error {
	a.pool.Lock()
	ws, ok := a.pool.conn[a.key]
	if ok && ws == a {
		delete(a.pool.conn, a.key)
		a.pool.Unlock()
		a.Lock()
		defer a.Unlock()
		return a.Close()
	} else {
		a.pool.Unlock()
	}
	return nil
}
