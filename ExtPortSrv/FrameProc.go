package ExtPortSrv

import (
	"RMS_Srv/ProtProcessor"
	ptb "RMS_Srv/Protocol"
	"RMS_Srv/Public"
	"fmt"
	"io"
	"net"
	"reflect"
	"syscall"
	"time"
)

//deal with the tcp connect,recieve data,
//unpack it and putting into frame
//then for the correct frame ,call the protocol dealing
func TcpFrameProcessor(tcpcon net.Conn) {

	var pt ptb.PackTag
	var rec []byte
	var pts_last uint16 = 0xffff
	var dat = make([]byte, 16384)
	var data []byte

	fmt.Println("TCP Client connect", tcpcon.RemoteAddr())
	defer tcpcon.Close()

	//reg new come node ip
	Public.OnlineNodes[tcpcon] = new(Public.NodeStats)
	Public.OnlineNodes[tcpcon].NodeIPP = tcpcon
	defer delete(Public.OnlineNodes, tcpcon)

	go ProtProcessor.SenderProcess()

	//tcpcon.SetReadDeadline(time.Now().Add(time.Second * 5))

	//test file download on connect
	//FileSrv.FileOpener()

	// reciever
newdata:
	for {
		n, err := tcpcon.Read(dat)
		if err == io.EOF || err != nil {
			//if errHandling(tcpcon, err) > 0 {
			//return
			//}
		}

		if n == 0 {
			time.Sleep(time.Microsecond)
		}

		data = append(data, dat[:n]...)
		n = len(data)
		for n > 0 {
			for i, chk := range data {
				if chk == 0x55 && data[i+1] == 0xAA {
					data = data[i:]
					n = n - i
					break
				}
			}

			//check head stx and lens
			if data[0] == 0x55 && data[1] == 0xAA && (n >= 12+(int(data[2])+int(data[3])*256)) {
				pt, rec, _ = ptb.Depack(data[:12+int(data[2])+int(data[3])*256])

				//regroup remained data,consider it's next package
				tmp := 12 + int(data[2]) + int(data[3])*256
				data = data[tmp:]
				n = n - tmp

			} else {
				continue newdata
			}

			if pt.Pserial == pts_last {
				//what? the same pack?
				//dump dumply data
				//if len(rec) >10 {
				//	fmt.Printf("%+v \t %x \n ",pt,rec[:10])
				//}else {
				//	fmt.Printf("%+v \t %x \n ",pt,rec)
				//}
				//data = data[pt.Plen + 12:]
				//n = n - int(pt.Plen + 12)
				continue newdata
				//logrus.Panic("problem catch, Err:[pt.Pserial] two same", pt.Pserial)
			}
			pts_last = pt.Pserial

			ProtProcessor.RecProcess(pt, rec, tcpcon)
		}
	}
}

func errHandling(tcpcon net.Conn, err error) int {
	rr := reflect.ValueOf(err).Elem().FieldByName("Err").Interface()
	ff := reflect.ValueOf(rr).Elem().FieldByName("Err")

	switch ff.Interface().(error) {

	case net.ErrWriteToConnected:
		fallthrough
	case io.ErrUnexpectedEOF:
		fallthrough

	case syscall.ENETRESET:
		fmt.Println("\nclient closed")
		return 1
	case syscall.ECONNABORTED:
	case syscall.ECONNRESET:
	case syscall.Errno(10054):
		fmt.Println("\n10054 closed", tcpcon.RemoteAddr())
		return 1
	default:
		fmt.Printf("%v \n%+v \n%q\n %t\n\n", err, err, err, err)
		return 1

	}
	return 0
}
