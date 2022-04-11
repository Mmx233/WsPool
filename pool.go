package pool

import (
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"net/http"
	"sync"
)

func NewPool(upper *websocket.Upgrader) *Pool {
	return &Pool{
		conn:  make(map[any]*Conn),
		upper: upper,
	}
}

type Pool struct {
	sync.RWMutex
	upper *websocket.Upgrader
	conn  map[any]*Conn
}

func (a *Pool) Range(f func(any, *Conn) bool) {
	a.RLock()
	defer a.RUnlock()
	for k, v := range a.conn {
		if !f(k, v) {
			break
		}
	}
}

func (a *Pool) Load(key any) (*Conn, bool) {
	a.RLock()
	defer a.RUnlock()
	conn, ok := a.conn[key]
	return conn, ok
}

func (a *Pool) NewConn(c *gin.Context, key any, resHeader http.Header) (*Conn, error) {
	ws, e := a.upper.Upgrade(c.Writer, c.Request, resHeader)
	if e != nil {
		return nil, e
	}
	conn := &Conn{
		Conn: ws,
		Pool: a,
		Key:  key,
	}
	a.Lock()
	_, ok := a.conn[key]
	if ok {
		a.conn[key].Lock()
		_ = a.conn[key].Conn.Close()
		a.conn[key].Unlock()
	}
	a.conn[key] = conn
	a.Unlock()
	return conn, nil
}
