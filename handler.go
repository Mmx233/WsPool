package ws

import (
	"encoding/json"
	"github.com/Mmx233/tool"
	"io"
)

type Msg struct {
	Type  uint8       `json:"type"`
	Value interface{} `json:"value,omitempty"`
}

func MsgHandler(conn *Conn, handler func(msgType int, r io.Reader)) {
	go func() {
		defer tool.Recover()
		defer func(conn *Conn) {
			_ = conn.Clear()
		}(conn)

		for {
			t, r, e := conn.NextReader()
			if e != nil {
				return
			}
			handler(t, r)
		}
	}()
}

func DecodeStandardTextMsg(r io.Reader) (*Msg, error) {
	var msg Msg
	return &msg, json.NewDecoder(r).Decode(&msg)
}
