package ProtProcessor

import (
	"PackFrame"
	"RMS_Srv/FileSrv"
	ptb "RMS_Srv/Protocol"
	"RMS_Srv/Public"
	"encoding/binary"
	"fmt"
	"net"
)

//frame dispatch
func RecProcess(pt ptb.PackTag, rec []byte) {

	switch pt.Pcmd {
	case ptb.Fc_fileTrans:
		FileSrv.FileReciever(pt, rec)
	case ptb.Fc_fileTranD:
		FileSrv.FileReciever(pt, rec)
	case ptb.Fc_HB:
		fmt.Printf("fc HB : %s", rec)
	default:

	}
}

//send with chan : Public.TcpSender_Ch ( TcpTrucker)
//type TcpTrucker struct {
//	Cmd int
//	Dat interface{}
//	Ext []interface{}
//}
func SenderProcess(tcpcon net.Conn) {
	for {
		switch c := <-Public.TcpSender_Ch; c.Cmd {
		case ptb.TSC_SendFile:
			FileSrv.Sendfile(tcpcon, c)
		case ptb.Fc_HB:
			//send file name
			var ss []byte = make([]byte, 8)
			binary.BigEndian.PutUint64(ss, uint64(c.Dat.(int)))
			ready, err := PackFrame.Dopack(ss,
				ptb.Fc_HB, 0)

			fmt.Printf("%s", ready)
			_, err = tcpcon.Write(ready)
			if err != nil {
			}
		default:

		}

	}

}
