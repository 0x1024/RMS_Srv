package DataBase_SAL

import (
	"RMS_Srv/AUTH_SAL"
	"RMS_Srv/FileSrv"
	"RMS_Srv/Public"
	"encoding/json"
	"errors"
	"fmt"
	"golang.org/x/net/websocket"
	"reflect"
	"strconv"
	"strings"
	"time"
)

type cmd struct {
	Cmd  string      `json:"cmd"`
	Data interface{} `json:"data"`
}

func ReqProcess(ws *websocket.Conn, dat string) {
	var err error
	var send cmd

	var dats map[string]string

	var rrr string
	rrr = strings.TrimRight(dat, "\x00")

	if err := json.Unmarshal([]byte(rrr), &dats); err == nil {
		//fmt.Printf("\r\nReqProcess   %q \r\n", dats) //debug
		//fmt.Println("cmd: ", dats["cmd"])            //debug
	}

	Senders := new(Public.Senders)
	Senders.Ws = ws

	//auth login=====================================
	if Public.LoginUser[ws].Logined == false {

		if dats["cmd"] == "auth_req" {

			auth_tmp := new(Um_index)
			_, err = AUTH_SAL.AuthEng.Where("name=?", dats["user"]).Get(auth_tmp)

			if err != nil {
				fmt.Println(err)
				send.Cmd = "auth_failed"
				send.Data = ""

			} else if (auth_tmp.Name == "") || (auth_tmp.Name == "0") {
				send.Cmd = "auth_name_fault"
				send.Data = ""

			} else if (dats["pswd"] == auth_tmp.Passwd) && (dats["pswd"] != "") {
				send.Cmd = "auth_ok"
				send.Data = ""
				Public.LoginUser[ws].Logined = true

				rolegroup := new(Role_group)

				Public.LoginUser[ws].Role = auth_tmp.Role
				_, err = AUTH_SAL.AuthEng.Table("role_group").Where("name=?", auth_tmp.Role).Get(rolegroup)
				fmt.Println(rolegroup)
				Public.LoginUser[ws].Priv = rolegroup.Priv
				Public.LoginUser[ws].Wlist = rolegroup.Wlist
				Public.LoginUser[ws].Blist = rolegroup.Blist

			} else if dats["cmd"] == "HB" {
				Public.LoginUser[ws].HBLife = 0
			} else {
				send.Cmd = "auth_pwd_fault"
				send.Data = ""

			}

			rec, _ := json.Marshal(send)
			fmt.Printf("json %q \r\n==========%q\r\n  \r\n", rec, send)
			data_tmp := string(rec)
			Senders.Dat = data_tmp
			Public.DB2Ret <- Senders
		} else if dats["cmd"] == "HB" {
			Public.LoginUser[ws].HBLife = 0
		}
		//end of    auth login=====================================
	} else {

		//fmt.Print("priv   ")
		//fmt.Printf("%b", Public.LoginUser[ws].Priv)
		switch dats["cmd"] {
		case "req":
			fmt.Printf("\npriv  %X  ", Public.LoginUser[ws].Priv)
			delete(dats, "cmd")
			if Public.LoginUser[ws].Priv&OP_Read != 0 {
				sa1 := new(Pd_index)
				re, _ := strconv.Atoi(dats["pid"])
				_, err = engine.Where("pid=?", re).Get(sa1)
				if err != nil {
					fmt.Println(err)
				}

				send.Cmd = "data_single"
				send.Data = sa1
				rec, _ := json.Marshal(send)
				data_tmp := string(rec)

				Senders.Dat = data_tmp
				Public.DB2Ret <- Senders
			} else {
				send.Cmd = dats["cmd"]
				authAct_NoPermition(Senders, send)
			}
		case "all":
			delete(dats, "cmd")
			if Public.LoginUser[ws].Priv&OP_Read != 0 {
				sa2 := new([]Pd_index)
				err = engine.Find(sa2)

				send.Cmd = "data_all"
				send.Data = sa2

				rec, _ := json.Marshal(send)

				data_tmp := string(rec)

				Senders.Dat = data_tmp
				Public.DB2Ret <- Senders

			} else {
				send.Cmd = dats["cmd"]
				authAct_NoPermition(Senders, send)
			}
		case "comitone":
			delete(dats, "cmd")
			if Public.LoginUser[ws].Priv&OP_Write != 0 {
				var recs int64
				//prepare data
				result := &Pd_index{}
				fmt.Printf("%q \r\n\n", dats) //debug
				err = FillStruct(dats, result)

				//check exsit
				sa1 := new(Pd_index)
				_, err = engine.Where("pid=?", result.Pid).Get(sa1)
				if sa1.Pid == 0 {
					//no item,insert new one
					recs, err = engine.InsertOne(result)
				} else {
					//no err means there is item
					recs, err = engine.Update(result, &Pd_index{Pid: result.Pid})

				}

				fmt.Println(recs, err) //debug

				send.Cmd = "respond"
				send.Data = nil

				rec, _ := json.Marshal(send)

				data_tmp := string(rec)

				Senders.Dat = data_tmp
				Public.DB2Ret <- Senders

			} else {
				send.Cmd = dats["cmd"]
				authAct_NoPermition(Senders, send)
			}
		case "update":
			delete(dats, "cmd")
			if Public.LoginUser[ws].Priv&OP_Write != 0 {
				result := &Pd_index{}
				fmt.Printf("%q \r\n\n", dats) //debug

				err = FillStruct(dats, result)
				recs, err := engine.Update(result, &Pd_index{Pid: result.Pid})
				fmt.Println(recs, err) //debug

				send.Cmd = "respond"
				send.Data = nil

				rec, _ := json.Marshal(send)

				data_tmp := string(rec)

				Senders.Dat = data_tmp
				Public.DB2Ret <- Senders

			} else {
				send.Cmd = dats["cmd"]
				authAct_NoPermition(Senders, send)
			}
		case "delete_id":
			delete(dats, "cmd")
			if Public.LoginUser[ws].Priv&OP_Delete != 0 {

				result := &Pd_index{}
				fmt.Printf("%q \r\n\n", dats) //debug

				err = FillStruct(dats, result)

				n, err := engine.Delete(result)
				if err != nil {
					fmt.Println(n, err)
				}

				send.Cmd = "respond"
				send.Data = nil

				rec, _ := json.Marshal(send)

				data_tmp := string(rec)

				Senders.Dat = data_tmp
				Public.DB2Ret <- Senders

			} else {
				send.Cmd = dats["cmd"]
				authAct_NoPermition(Senders, send)
			}
		case "HB":
			delete(dats, "cmd")

			//if Public.LoginUser[ws].Priv&OP_Ping != 0 {
			Public.LoginUser[ws].HBLife = 0

		//}
		case "updateFW":
			delete(dats, "cmd")
			if Public.LoginUser[ws].Priv&OP_Manage != 0 {

				re, _ := strconv.Atoi(dats["pid"])
				for _, v := range Public.OnlineNodes {
					if v.PID == uint64(re) {
						FileSrv.FileOpener(v.NodeIPP)

					}
				}

			} else {
				send.Cmd = dats["cmd"]
				authAct_NoPermition(Senders, send)
			}

		default:
			send.Cmd = "respond"
			send.Data = nil

			rec, _ := json.Marshal(send)
			delete(dats, "cmd")

			data_tmp := string(rec)
			Senders.Dat = data_tmp
			Public.DB2Ret <- Senders

		} //end of switch
	} //end of  if logined

}

func FillStruct(data map[string]string, obj interface{}) error {
	for k, v := range data {
		err := SetField(obj, k, v)
		if err != nil {
			return err
		}
	}
	return nil
}

//用map的值替换结构的值
func SetField(obj interface{}, name string, value interface{}) error {

	name = strings.Title(name)
	structValue := reflect.ValueOf(obj).Elem()        //结构体属性值
	structFieldValue := structValue.FieldByName(name) //结构体单个属性值

	if !structFieldValue.IsValid() {
		return fmt.Errorf("No such field: %s in obj", name)
	}

	if !structFieldValue.CanSet() {
		return fmt.Errorf("Cannot set %s field value", name)
	}

	structFieldType := structFieldValue.Type() //结构体的类型
	val := reflect.ValueOf(value)              //map值的反射值

	var err error
	if structFieldType != val.Type() {
		val, err = TypeConversion(fmt.Sprintf("%v", value), structFieldValue.Type().Name()) //类型转换
		if err != nil {
			return err
		}
	}

	structFieldValue.Set(val)
	return nil
}

//类型转换
func TypeConversion(value string, ntype string) (reflect.Value, error) {
	if ntype == "string" {
		return reflect.ValueOf(value), nil
	} else if ntype == "time.Time" {
		t, err := time.ParseInLocation("2006-01-02 15:04:05", value, time.Local)
		return reflect.ValueOf(t), err
	} else if ntype == "Time" {
		t, err := time.ParseInLocation("2006-01-02 15:04:05", value, time.Local)
		return reflect.ValueOf(t), err
	} else if ntype == "int" {
		i, err := strconv.Atoi(value)
		return reflect.ValueOf(i), err
	} else if ntype == "int8" {
		i, err := strconv.ParseInt(value, 10, 64)
		return reflect.ValueOf(int8(i)), err
	} else if ntype == "int32" {
		i, err := strconv.ParseInt(value, 10, 64)
		return reflect.ValueOf(int64(i)), err
	} else if ntype == "int64" {
		i, err := strconv.ParseInt(value, 10, 64)
		return reflect.ValueOf(i), err
	} else if ntype == "float32" {
		i, err := strconv.ParseFloat(value, 64)
		return reflect.ValueOf(float32(i)), err
	} else if ntype == "float64" {
		i, err := strconv.ParseFloat(value, 64)
		return reflect.ValueOf(i), err
	}

	//else if .......增加其他一些类型的转换

	return reflect.ValueOf(value), errors.New("未知的类型：" + ntype)
}

func authAct_NoPermition(Senders *Public.Senders, send cmd) {
	send.Data = fmt.Sprintf("%s,%s", send.Cmd, "No Permittion")
	send.Cmd = ""
	rec, _ := json.Marshal(send)
	data_tmp := string(rec)

	Senders.Dat = data_tmp
	Public.DB2Ret <- Senders
}
