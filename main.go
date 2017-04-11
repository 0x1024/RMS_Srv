package main

import (
	"RMS_Srv/DataBase_SAL"
	"RMS_Srv/Public"
	"RMS_Srv/WEB_IO"
	"golang.org/x/net/websocket"
)

var RMS_EXIT chan int

func main() {
	Public.LoginUser = make(map[*websocket.Conn]*Public.LoginType)

	go DataBase_SAL.DB_Init()

	go WEB_IO.Http_init()

	<-RMS_EXIT
}
