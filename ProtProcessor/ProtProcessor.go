package ProtProcessor

import (
	"RMS_Srv/DataBase_SAL"
	"RMS_Srv/FileSrv"
	ptb "RMS_Srv/Protocol"
	"RMS_Srv/Public"
	"fmt"
	"iMQ"
	"net"
	"time"
)

//frame dispatch
func RecProcess(pt ptb.PackTag, rec []byte, tcpcon net.Conn) {

	switch pt.Pcmd {
	case ptb.Fc_fileTrans:
		FileSrv.FileReciever(pt, rec)
	case ptb.Fc_fileTranD:
		FileSrv.FileReciever(pt, rec)
	case ptb.Fc_HB:
		fmt.Printf("client %v say: %s \n", pt, rec)
		//iMQ.Imqsrv.PublishMessage("NodeStat", rec)
	//Public.DB2Ret <- s
	case ptb.Fc_HC:
		fmt.Printf("client %v say: % X \n", pt, rec)

		Public.OnlineNodes[tcpcon].McuId = rec
		Public.DevNodes_Ch <- *Public.OnlineNodes[tcpcon]
		//Public.DB2Ret <- s
		//
	default:
		fmt.Printf("RecProcess say WTF: \n", pt, rec)

	}
}

func DevOnlineRegin() {
	//reg id
	for {
		tmp := <-Public.DevNodes_Ch

		ack := DataBase_SAL.QueryMCUID(Public.OnlineNodes[tmp.NodeIPP].McuId)

		Public.OnlineNodes[tmp.NodeIPP].PID = ack

		fmt.Printf("dev %d online from %s", ack, Public.OnlineNodes[tmp.NodeIPP].NodeIPP.RemoteAddr())

		//n := fmt.Sprintf("%d",Public.OnlineNodes[tmp.NodeIPP].PID )
		//iMQ.Imqsrv.PublishMessage("NodeStat",[]byte(n) )

	}
}

func DevOnlineManage() {
	//iMQ.Imqsrv.PublishMessage("NodeStat", rec)

	go DevOnlineRegin()

	//query id
	for {
		for k, v := range Public.OnlineNodes {
			n := fmt.Sprintf("%d %s", v.PID, k.RemoteAddr())
			iMQ.Imqsrv.PublishMessage("NodeStat", []byte(n))
		}
		time.Sleep(10e9)
	}

}

//send with chan : Public.TcpSender_Ch ( TcpTrucker)
//type TcpTrucker struct {
//	Cmd int
//	Dat interface{}
//	Ext []interface{}
//}

func SenderProcess() {
	for {
		c, _ := <-Public.TcpSender_Ch
		tcpcon := c.Ip
		switch c.Cmd {
		case ptb.TSC_SendFile:
			FileSrv.Sendfile(tcpcon, c)
		case ptb.Fc_HB:
			//send file name
			//var ss []byte = make([]byte, 8)
			//binary.BigEndian.PutUint64(ss, uint64(c.Dat.(int)))
			ready, err := ptb.Dopack(c.Dat.([]byte),
				ptb.Fc_HB, 0)

			fmt.Printf("Fc_HB %s\n", ready)
			_, err = tcpcon.Write(ready)
			if err != nil {
			}
		case ptb.Fc_HC:
			//send file name
			//var ss []byte = make([]byte, 8)
			//binary.BigEndian.PutUint64(ss, uint64(c.Dat.(int)))
			ready, err := ptb.Dopack(c.Dat.([]byte),
				ptb.Fc_HC, 0)

			fmt.Printf("Fc_HC %s\n", ready)
			_, err = tcpcon.Write(ready)
			if err != nil {
			}
		default:
			fmt.Println("send cmd not found,it's ", c.Cmd)
		}

	}

}
