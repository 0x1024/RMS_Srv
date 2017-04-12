package main

import (
	"RMS_Srv/DataBase_SAL"
	"RMS_Srv/ExtPortSrv"
	"RMS_Srv/WEB_IO"
)

var RMS_EXIT chan int

func main() {
	WEB_IO.Init()

	go DataBase_SAL.DB_Init()

	go WEB_IO.Http_init()

	go ExtPortSrv.ExternService()

	<-RMS_EXIT
}
