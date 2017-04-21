package FileSrv

import (
	ptb "RMS_Srv/Protocol"
	"bytes"
	"encoding/binary"
	"fmt"
	"github.com/sirupsen/logrus"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"strings"
	"time"
)

var fo *os.File
var file_len int64 = 1024
var file_block_counter uint = 0
var TimeEslaped time.Time

//
func FileReciever(pt ptb.PackTag, rec []byte) {
	var err error
	switch pt.Ppara & 0xF {
	case ptb.Fcp_fileName:
		TimeEslaped = time.Now()
		fo, err = os.Create(getCurrentDirectory() + "/rec/" + string(rec))
		if err != nil {
			err = os.Mkdir(getCurrentDirectory()+"/rec/", os.ModeDir)
			if err != nil {
				logrus.Panic(err)
			}
			fo, err = os.Create(getCurrentDirectory() + "/rec/" + string(rec))
			if err != nil {
				logrus.Panic(err)
			}
		}
		file_block_counter = 0

	case ptb.Fcp_fileSize:
		bb := bytes.NewBuffer(rec)
		binary.Read(bb, binary.LittleEndian, &file_len)

	case ptb.Fcp_fileEOF:
		//fmt.Printf("%d \t/ %d \r", file_block_counter, file_len/1024)
		logrus.Info("file rec end")
		fmt.Print(path.Dir(fo.Name()))
		cmd := exec.Command("explorer", strings.Replace(path.Dir(fo.Name()), "/", "\\", -1))
		//cmd:=exec.Command("explorer",")
		err := cmd.Run()
		if err == nil {
			fmt.Println("\nerr:", err)
		}
		fo.Close()
		fmt.Println("\ntime cost(sec) :", float32(time.Now().Sub(TimeEslaped).Nanoseconds())/1e9)
		return

	default:

		panic("how did reach the empty filehead?? ")

	case ptb.Fcp_filedata:
		//write to the file
		jn, err := fo.Write(rec)
		if err != nil {
			logrus.Panic(err)
		}
		//logrus.Info("%d \t/ %d \r", file_block_counter, file_len/1024 )
		//		fmt.Printf("%d \t/ %d \r", file_block_counter, file_len/1024)
		fmt.Println("file rec : ", file_block_counter)
		file_block_counter++
		jn = jn
		err = err
	}
}

func getCurrentDirectory() string {
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		panic(err)
	}
	return strings.Replace(dir, "\\", "/", -1)
}
