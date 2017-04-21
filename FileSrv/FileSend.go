package FileSrv

import (
	"PackFrame"
	ptb "RMS_Srv/Protocol"
	"RMS_Srv/Public"
	"fmt"
	"io"
	"net"
	"os"
	"time"
)

func FileOpener() {
	var c Public.TcpTrucker
	input := "e:\\2.rom"
	fi, err := os.Open(string(input))
	if err != nil {
		panic(err)
	}
	c.Cmd = ptb.TSC_SendFile
	c.Dat = fi
	Public.TcpSender_Ch <- c
}

func Sendfile(conn net.Conn, c Public.TcpTrucker) {
	var err error
	var fi *os.File
	fi = c.Dat.(*os.File)
	defer fi.Close()
	fiinfo, err := fi.Stat()
	fmt.Println("the size of file is ", fiinfo.Size(), "bytes") //fiinfo.Size() return int64 type

	nn := time.Now()
	fmt.Println("now", nn.Format(time.RFC3339Nano))

	//send file name
	ready, err := PackFrame.Dopack([]byte(fiinfo.Name()),
		ptb.Fc_fileTrans, ptb.Fcp_fileName)

	fmt.Printf("%s", ready)
	_, err = conn.Write(ready)
	if err != nil {
		fmt.Println("conn.Write", err.Error())
	}
	time.Sleep(time.Microsecond * 5)
	//send file size
	ready, err = PackFrame.Dopack(PackFrame.TypeToByte(fiinfo.Size()),
		ptb.Fc_fileTrans, ptb.Fcp_fileSize)

	_, err = conn.Write(ready)
	time.Sleep(time.Microsecond * 5)
	//_, err = conn.Write([]byte(string(fiinfo.Size())))
	if err != nil {
		fmt.Println("conn.Write", err.Error())
	}
	time.Sleep(time.Microsecond * 5)

	var ctr uint32 = 0
	for {
		//fmt.Println("No ", ctr, "of total ", fiinfo.Size()/1024)
		buff := make([]byte, 1024*8)
		n, err := fi.Read(buff)
		if err != nil && err != io.EOF {
			panic(err)
		}

		if n == 0 {
			ready, err = PackFrame.Dopack(buff[:n], ptb.Fc_fileTrans, ptb.Fcp_fileEOF)
			_, err = conn.Write(ready)

			fmt.Println("time cost ", time.Now().Sub(nn))

			fmt.Println("\nfile send finished")
			break
		}

		ready, err = PackFrame.Dopack(buff[:n], ptb.Fc_fileTrans, ptb.Fcp_filedata|(ctr<<4))
		_, err = conn.Write(ready)
		//		_, err = conn.Write(buff)
		if err != nil {
			fmt.Println(err.Error())
		}
		fmt.Println("file send process", ctr)
		ctr++
		time.Sleep(1)
	}

}
