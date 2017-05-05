package Public

import (
	"golang.org/x/net/websocket"
	"net"
	"time"
)

var Signal = make(chan *Senders, 2)
var DB2Ret = make(chan *Senders, 2)

var TcpSender_Ch chan TcpTrucker = make(chan TcpTrucker, 16)

type TcpTrucker struct {
	Cmd int
	Ip  net.Conn
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

type NodeStats struct {
	NodeIPP net.Conn
	McuId   []byte
	PID     uint64
}

var LocalNode NodeStats = NodeStats{PID: 1}
var DevNodes_Ch chan NodeStats = make(chan NodeStats, 16)

var OnlineNodes map[net.Conn]*NodeStats

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

func Init() {
	OnlineNodes = make(map[net.Conn]*NodeStats)
	LoginUser = make(map[*websocket.Conn]*LoginType)

}
