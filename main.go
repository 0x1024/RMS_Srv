package main

import (
	"RMS_Srv/DataBase_SAL"
	"RMS_Srv/ExtPortSrv"
	"RMS_Srv/Public"
	"RMS_Srv/WEB_IO"
	"iMQ"
)

var RMS_EXIT chan int

func main() {
	Public.Init()
	iMQ.Init()

	go DataBase_SAL.DB_Init()

	go WEB_IO.Http_init()

	go ExtPortSrv.TcpServerStarter()

	<-RMS_EXIT
}
