package ExtPortSrv

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"net"
	"net/http"
	"os"
	"time"
)

func ExternService() {
	fmt.Println("\n\n\n[INFO]start server....", time.Now().Format(time.UnixDate))
	TcpServer()
}

//Client Device Server,provide robot connect service
func TcpServer() {

	//pprof service
	go func() {
		logrus.Println(http.ListenAndServe("localhost:6060", nil))
	}()

	listener, err := net.Listen("tcp", "127.0.0.1:8888")
	if err != nil {
		fmt.Println("error listening:", err.Error())
		os.Exit(1)
	}
	defer listener.Close()
	fmt.Printf("Client Device Server running ...\n")

	//var cur_conn_num int = 0
	conn_chan := make(chan net.Conn)

	go func() {
		for {
			conn, err := listener.Accept()
			if err != nil {
				println("Error accept:", err.Error())
				return
			}
			conn_chan <- conn
		}
	}()

	for {
		conn := <-conn_chan
		go TcpFrameProcessor(conn)
	}
}
