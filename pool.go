package pool

import (
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"net/http"
	"sync"
	"sync/atomic"
)

func New(upper *websocket.Upgrader) *Pool {
	return &Pool{
		Upper:  upper,
		pool:   &sync.Map{},
		length: &atomic.Int64{},
	}
}

type Pool struct {
	Upper  *websocket.Upgrader
	pool   *sync.Map
	length *atomic.Int64
}

func (a *Pool) Len() int {
	return int(a.length.Load())
}

func (a *Pool) Range(f func(any, *Conn) bool) {
	a.pool.Range(func(key, value any) bool {
		return f(key, value.(*Conn))
	})
}

func (a *Pool) Load(key any) (*Conn, bool) {
	value, ok := a.pool.Load(key)
	if !ok {
		return nil, false
	}
	return value.(*Conn), true
}

func (a *Pool) NewConn(c *gin.Context, key any, resHeader http.Header) (*Conn, error) {
	var newConn = &Conn{
		Pool: a,
		Key:  key,
	}
	conn, loaded := a.pool.LoadOrStore(key, newConn)
	if loaded {
		return conn.(*Conn), nil
	}

	ws, err := a.Upper.Upgrade(c.Writer, c.Request, resHeader)
	if err != nil {
		return nil, err
	}
	newConn.Conn = ws
	a.length.Add(1)
	return newConn, nil
}
