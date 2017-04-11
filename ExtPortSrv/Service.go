package ExtPortSrv

import (
	"PackFrame"
	"bytes"
	"encoding/binary"
	"fmt"
	"github.com/sirupsen/logrus"
	"io"
	"net"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"syscall"
	"time"
)

//frame cmd type list
const (
	fc_filehead = 0x10
	fc_filebody = 0x11
)

// fc file paras
const (
	fcp_fileName = 0x01
	fcp_fileEOF  = 0x02
	fcp_fileSize = 0x03
)

//service for dev, provide communicate
func EchoFunc(tcpcon net.Conn) {

	fmt.Println("Client connect", tcpcon.RemoteAddr())
	defer tcpcon.Close()

	var pts_last uint16
	var fo *os.File
	var file_len int64 = 1024
	var file_block_counter int = 0

	var data []byte
	data = make([]byte, 16384)

	var pt PackFrame.PackTag
	var rec []byte
	tcpcon.SetDeadline(time.Unix(5, 0))
	for {
		n, err := tcpcon.Read(data)

		if err != io.EOF && err != nil {

			rr := reflect.ValueOf(err).Elem().FieldByName("Err").Interface()
			ff := reflect.ValueOf(rr).Elem().FieldByName("Err")

			switch ff.Interface().(error) {

			case net.ErrWriteToConnected:
				fallthrough
			case io.ErrUnexpectedEOF:
				fallthrough

			case syscall.ENETRESET:
				fmt.Println("\nclient closed")
				return
			case syscall.ECONNABORTED:
			case syscall.ECONNRESET:
			case syscall.Errno(10054):
				fmt.Println("\n10054 closed", tcpcon.RemoteAddr())
				return
			default:
				fmt.Printf("%v \n%+v \n%q\n %t\n\n", err, err, err, err)
				return

			}
			return
		} //if err != io.EOF && err != nil

		if n == 0 {
			time.Sleep(time.Microsecond * 10)
		}

		for n > 0 {

			//fmt.Println(tcpcon.RemoteAddr(),n)
			file_block_counter++
		_L_next_head:
			if data[0] == 0x55 && data[1] == 0xAA {
				pt, rec, err = PackFrame.Depack(data[:12+int(data[2])+int(data[3])*256])
				if err != nil {

				}
				tmp := 13 + int(data[2]) + int(data[3])*256
				data = data[tmp-1:]
				n = n - tmp

			} else {
				for i, chk := range data[:n] {
					if chk == 0x55 && data[i+1] == 0xAA {
						data = data[i:]
						n = n - i
						fmt.Println("seek head next")
						goto _L_next_head
					}

				}
				n = 0
				data = nil
				continue
			}

			if pt.Pserial == pts_last {
				//what? the same pack?
				logrus.Panic("problem catch, Err:[pt.Pserial] %d two same", pt.Pserial)
			}
			pts_last = pt.Pserial

			switch pt.Pcmd {

			case fc_filehead:
				if pt.Ppara == fcp_fileName {
					fmt.Println(getCurrentDirectory())
					fo, err = os.Create(getCurrentDirectory() + "//rec//" + string(rec))
					//fo, err = os.Create( "e://"  + string(rec))
				} else if pt.Ppara == fcp_fileSize {
					bb := bytes.NewBuffer(rec)
					binary.Read(bb, binary.LittleEndian, &file_len)

				} else if pt.Ppara == fcp_fileEOF {
					return
				} else {
					panic("how did reach the empty filehead?? ")
				}
			case fc_filebody:
				//write to the file
				jn, err := fo.Write(rec)
				fmt.Printf("%d \t/ %d \r", file_block_counter, file_len/1024)
				jn = jn
				err = err
			default:

			}
		}

	}
}

func getCurrentDirectory() string {
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		panic(err)
	}
	return strings.Replace(dir, "\\", "/", -1)
}
