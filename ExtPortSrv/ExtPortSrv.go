package ExtPortSrv

import (
	"RMS_Srv/ProtProcessor"
	"RMS_Srv/Public"
	"fmt"
	"github.com/sirupsen/logrus"
	"net"
	"net/http"
	"os"
	"time"
)

//Client Device Server,provide robot connect service
func TcpServerStarter() {
	fmt.Println("\n\n\n[INFO]start server....", time.Now().Format(time.UnixDate))
	go ProtProcessor.DevOnlineManage()
	//pprof service
	go func() {
		logrus.Println(http.ListenAndServe("localhost:6060", nil))
	}()

	listener, err := net.Listen("tcp", ":8866")
	if err != nil {
		fmt.Println("error tcp listening:", err.Error())
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

//Client node ,provide robot connect service
func NodeStarter() {
	var conn net.Conn
	var err error
	//to online
	for {
		//local debug
		//conn1, err = net.Dial("tcp", "127.0.0.1:8866")
		//if err != nil {
		//
		//} else {
		//	defer conn1.Close()
		//	go TcpFrameProcessor(conn1)
		//}
		//server apply
		conn, err = net.Dial("tcp", "118.178.138.192:8866")
		//conn, err = net.Dial("tcp", "127.0.0.1:8866")
		if err != nil {
			fmt.Println("connect server fail！", err.Error())
			time.Sleep(10e9)
			continue
		}
		defer conn.Close()
		Public.LocalNode.NodeIPP = conn

		TcpFrameProcessor(conn)
		break
	}

	EXIT_ch := make(chan bool)
	<-EXIT_ch

}
