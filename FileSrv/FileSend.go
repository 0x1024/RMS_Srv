package FileSrv

import (
	ptb "RMS_Srv/Protocol"
	"fmt"
	"io"
	"net"
	"os"
)

func Client() {
	//
	//open file
	//fmt.Println("send ur file to the destination", "input ur filename:")
	//reader := bufio.NewReader(os.Stdin)
	//input, _, _ := reader.ReadLine()
	//fmt.Println(string(input))

	input := "e:\\2.rom"
	fi, err := os.Open(string(input))
	if err != nil {
		panic(err)
	}
	defer fi.Close()
	fiinfo, err := fi.Stat()
	fmt.Println("the size of file is ", fiinfo.Size(), "bytes") //fiinfo.Size() return int64 type

	//to online
	conn, err := net.Dial("tcp", "127.0.0.1:8888")
	if err != nil {
		fmt.Println("connect server fail！", err.Error())
		return
	}
	defer conn.Close()

	//send file name
	ready, err := ptb.Dopack([]byte(fiinfo.Name()),
		ptb.Fc_fileTrans, ptb.Fcp_fileName)
	fmt.Printf("%s", ready)
	_, err = conn.Write(ready)
	if err != nil {
		fmt.Println("conn.Write", err.Error())
	}

	//send file size
	ready, err = ptb.Dopack(ptb.TypeToByte(fiinfo.Size()),
		ptb.Fc_fileTrans, ptb.Fcp_fileSize)
	_, err = conn.Write(ready)
	if err != nil {
		fmt.Println("conn.Write", err.Error())
	}

	//send file body
	var ctr uint32 = 0
	for {
		buff := make([]byte, 1024)
		n, err := fi.Read(buff)
		if err != nil && err != io.EOF {
			panic(err)
		}

		if n == 0 {
			ready, err = ptb.Dopack(buff[:n], ptb.Fc_dataTrans, ptb.Fcp_fileEOF)
			_, err = conn.Write(ready)
			//conn.Write([]byte("filerecvend"))
			fmt.Println("filerecvend")
			break
		}

		ready, err = ptb.Dopack(buff[:n], ptb.Fc_fileTrans, ptb.Fcp_filedata|(ctr<<4))
		_, err = conn.Write(ready)
		if err != nil {
			fmt.Println(err.Error())
		}
		ctr++
	}
}
