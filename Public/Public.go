package Public

import (
	"golang.org/x/net/websocket"
	"time"
)

var Signal = make(chan *Senders, 2)
var DB2Ret = make(chan *Senders, 2)

var TcpSender_Ch chan TcpTrucker = make(chan TcpTrucker, 16)

type TcpTrucker struct {
	Cmd int
	Dat interface{}
	Ext []interface{}
}

//server ws
//	Ws  *websocket.Conn
//      Dat string
type Senders struct {
	Ws  *websocket.Conn
	Dat string
}

var LoginUser map[*websocket.Conn]*LoginType

type LoginType struct {
	Handle  *websocket.Conn
	Name    string
	InDT    time.Time
	RsaPri  []byte
	RsaPub  []byte
	PplId   uint64
	HBLife  int
	Logined bool
	Role    string
	Priv    uint
	Wlist   string //string={'"asdf","fda","fff" '}
	Blist   string
}
