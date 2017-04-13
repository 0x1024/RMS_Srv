package WEB_IO

import (
	"RMS_Srv/DataBase_SAL"
	"RMS_Srv/Public"

	"encoding/json"
	"fmt"
	log "github.com/sirupsen/logrus"
	"golang.org/x/net/websocket"
	"io"
	rand2 "math/rand"
	"net/http"
	"time"
)

type cmd struct {
	Cmd  string      `json:"cmd"`
	Data interface{} `json:"data"`
}

var WEBIO_EXIT chan int

func Init() {
	Public.LoginUser = make(map[*websocket.Conn]*Public.LoginType)
}

func Http_init() {

	http.Handle("/", websocket.Handler(echoHandler))

	//no tls
	go http.ListenAndServe(":9003", nil)

	//tls addon test
	go http.ListenAndServeTLS(":9004", "sign.pem", "ssl.key", nil)

	<-WEBIO_EXIT
}

//generate privilege passport license
func GenPPL(ws *websocket.Conn) {

	var nopass bool = true
	var tmp uint64
	fmt.Printf("ppl in  \r\n\n")
	for nopass {
		nopass = false
		tmp = rand2.Uint64()
		for _, v := range Public.LoginUser {
			//fmt.Printf("member %q,,%q  \r\nnn", n, v)
			if v.PplId == tmp {
				nopass = true
			}
		}
	}
	//fmt.Printf(" %s PPL is %d \r\n\n\n", ws, tmp)
	Public.LoginUser[ws].PplId = tmp
}

//heart beat ,living check
func HB(ws *websocket.Conn) {
	Senders := new(Public.Senders)
	var send cmd
	send.Cmd = "HB"
	send.Data = ""
	rec, _ := json.Marshal(send)
	data_tmp := string(rec)

	for true {

		if Public.LoginUser[ws] != nil {
			Public.DB2Ret <- Senders
			Senders.Ws = ws
			Senders.Dat = data_tmp
			Public.LoginUser[ws].HBLife = Public.LoginUser[ws].HBLife + 1
			fmt.Println("HBL", Public.LoginUser[ws].HBLife)
			if Public.LoginUser[ws].HBLife > 10 {
				ws.Close()
			}
		}
		time.Sleep(5e9)
	}
}

func echoHandler(ws *websocket.Conn) {
	var err error
	var n int

	defer func() {
		if err := recover(); err != nil {
			strLog := "longweb:main recover error => " + fmt.Sprintln(err)
			//os.Stdout.Write([]byte(strLog))
			log.Error(strLog)

			//buf := make([]byte, 8192)
			//n := runtime.Stack(buf, true)
			//log.Error(string(buf[:n]))
			//os.Stdout.Write(buf[:n])
		}
	}()

	defer ws.Close()
	go sender()
	go HB(ws)

	fmt.Println("\n\n\n client addr :", ws.Request().RemoteAddr)

	//register current dialog
	if _, ok := Public.LoginUser[ws]; !ok {
		Public.LoginUser[ws] = new(Public.LoginType)
		Public.LoginUser[ws].Name = "匿名"
		Public.LoginUser[ws].Handle = ws
		go GenPPL(ws)

	}

	msg := make([]byte, 1024)
	for true {

		err = ws.SetDeadline(time.Now().Add(30e9))
		n, err = ws.Read(msg)
		//		err = ws.SetReadDeadline(time.Unix(0,0))
		if err != nil {
			fmt.Printf("errss %s\n", err)
			switch {
			case err == io.EOF:
				delete(Public.LoginUser, ws)
				fmt.Println("\n\n\nusers %q：\n\n\n", Public.LoginUser)
				ws.Close()
				goto out
			default:
				log.Fatal(err)
			}
		}

		fmt.Printf("Receive:[%s] %s\n", time.Now().Format(time.UnixDate), msg[:n])
		DataBase_SAL.ReqProcess(ws, string(msg[:n]))

	}
out:
}

// ws send port, input with channel
func sender() {

	for {
		rec := <-Public.DB2Ret
		fmt.Printf("sender to send :%\r\n", rec)
		_, err := rec.Ws.Write([]byte(rec.Dat))
		if err != nil {
			fmt.Printf("sender err %s\n", err)
			fmt.Print(err)
			switch {
			case err == io.EOF:
				goto Exit
			default:
				goto Exit
				log.Fatal("Fatal Err: %s \r\n", err)
			}
		}
	}
Exit:
}

//bits := 1024
//if err := GenRsaKey(bits); err != nil {
//	log.Fatal("密钥文件生成失败！")
//}
//log.Println("密钥文件生成成功！")
//
//initData := "abcdefghijklmnopq"
//init := []byte(initData)
////load_keys()
//
//data, err := RsaEncrypt(init)
//if err != nil {
//	panic(err)
//}
//pre := time.Now()
//origData, err := RsaDecrypt(data)
//if err != nil {
//	panic(err)
//}
//now := time.Now()
//fmt.Println(now.Sub(pre))
//fmt.Println(string(origData))
