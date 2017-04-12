package DataBase_SAL

import (
	"RMS_Srv/AUTH_SAL"
	"fmt"
	_ "github.com/lib/pq"
	log "github.com/sirupsen/logrus"
	"github.com/xormplus/xorm"
	"time"
)

type Pd_index struct {
	Pid     uint64    "1,Pid"
	Dtype   string    "2,Dtype"
	Client  string    "3,Client"
	Tags    string    "4,Tags"
	Passwd  string    "5,Passwd"
	Created time.Time "6,Created"
	Updated time.Time "7,Updated"
}

//user manage
type Um_index struct {
	Pid     uint64 "1,Pid"
	Uid     uint64
	Name    string    "2,Name"
	Passwd  string    "3,Passwd"
	Role    string    "4,Level"
	Tags    string    "5,Tags"
	Created time.Time "6,Created"
	Updated time.Time "7,Updated"
	Stamp   uint64
	Jail    time.Duration
}

//client group
type Customer struct {
	Cid     uint64 "1,Cid"
	Uid     []uint64
	Pid     []uint64
	Created time.Time "6,Created"
	Updated time.Time "7,Updated"
}

//role group
type Role_group struct {
	Rid     uint64 "1,Rid"
	Name    string "json:name"
	Priv    uint
	Wlist   string //string={'"asdf","fda","fff" '}
	Blist   string
	Created time.Time "6,Created"
	Updated time.Time "7,Updated"
}

const (
	OP_Null   = 0
	OP_Ping   = 1 << 1
	OP_Read   = 1 << 2
	OP_Write  = 1 << 3
	OP_Create = 1 << 4
	OP_Delete = 1 << 5
	OP_Manage = 1 << 6
	OP_SysLv0 = 1 << 7
	OP_SysLv1 = 1 << 8
	OP_SysLv2 = 1 << 9
	OP_SysLv3 = 1 << 10
)

var engine *xorm.Engine

var DB_EXIT chan int

func DB_Init() {
	var err error

	//=====================================================================
	//open db
	engine, err = xorm.NewPostgreSQL("postgres://rms:123@10.1.11.151:5432/RMSDB?sslmode=disable")
	if err != nil {
		log.Panic(err)
	}

	AUTH_SAL.AuthEng, err = xorm.NewPostgreSQL("postgres://rms:123@10.1.11.151:5432/AuthDB?sslmode=disable")
	if err != nil {
		log.Panic(err)
	}

	defer fmt.Printf("\n\n\n Db_Service is closed\n\n\n")
	defer engine.Close()
	defer AUTH_SAL.AuthEng.Close()

	//en func
	engine.ShowSQL(true)
	//	engine.Logger().SetLevel(core.LOG_DEBUG)
	engine.ShowExecTime(true)
	//	AUTH_SAL.AuthEng.Logger().SetLevel(core.LOG_DEBUG)
	AUTH_SAL.AuthEng.ShowSQL(true)

	//check table
	err = engine.CreateTables(new(Pd_index))
	if err != nil {
		panic(err)
	}
	err = AUTH_SAL.AuthEng.CreateTables(new(Um_index))
	if err != nil {
		panic(err)
	}
	err = AUTH_SAL.AuthEng.CreateTables(new(Customer))
	if err != nil {
		panic(err)
	}
	err = AUTH_SAL.AuthEng.CreateTables(new(Role_group))
	if err != nil {
		panic(err)
	}

	//hold on
	<-DB_EXIT
}
