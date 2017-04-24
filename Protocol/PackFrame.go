package Protocol

import (
	"encoding/binary"
	"fmt"
	"math"
	"unsafe"
)

var send_serial uint16

func serialGen() uint16 {
	send_serial++
	return send_serial
}

//@return the packed data  []byte, error
//the do pack is to make the message loadto a packet frame......
//the frame is for application control and serial correct.....
// msg is the info body ,len is just sizeof msg......
// serial notes counts every times dopack called,to be reconige on client,.....
// cmd para used for function control
//
func Dopack(loads []byte, cmd uint16, para uint32) ([]byte, error) {
	var err error = nil
	var pk []byte
	var p PackTag

	p.Phead = 0xAA55
	p.Plen = uint16(len(loads))
	p.Pcmd = cmd
	p.Ppara = para
	p.Pserial = serialGen()

	l := unsafe.Sizeof(p)
	pb := (*[1024]byte)(unsafe.Pointer(&p))
	lenss := len(loads) + int(l)
	pk = make([]byte, lenss)
	copy(pk, (*pb)[:l])
	copy(pk[l:], loads)

	if err != nil {
		return nil, err
	}

	//ret1 := []byte(p)
	//ret := bytes.Join(ret1, loads)
	return pk, err
}

func Depack(d []byte) (PackTag, []byte, error) {
	var pt PackTag
	l := unsafe.Sizeof(pt)
	pb := (*[1024]byte)(unsafe.Pointer(&pt))
	copy((*pb)[:l], d[:l])
	//	fmt.Println(pt)
	return pt, d[l:], nil
}

//
//func main(){
//	dd :=[]byte{0,1,2,3,4,5}
//	tr,err :=Dopack(dd,0,1)
//	fmt.Printf("%X,   %q",tr,err)
//	var pt PackTag
//	pt,rec,err:=Depack(tr)
//	fmt.Printf("%+v,    %x,   %s",pt,rec,err)
//}

func Float32ToByte(float float32) []byte {
	bits := math.Float32bits(float)
	bytes := make([]byte, 4)
	binary.LittleEndian.PutUint32(bytes, bits)

	return bytes
}

func ByteToFloat32(bytes []byte) float32 {
	bits := binary.LittleEndian.Uint32(bytes)
	return math.Float32frombits(bits)
}

func Float64ToByte(float float64) []byte {
	bits := math.Float64bits(float)
	bytes := make([]byte, 8)
	binary.LittleEndian.PutUint64(bytes, bits)

	return bytes
}

func ByteToFloat64(bytes []byte) float64 {
	bits := binary.LittleEndian.Uint64(bytes)

	return math.Float64frombits(bits)
}

func TypeToByte(v interface{}) []byte {

	switch value := v.(type) {

	case []byte:
		return value

	case string:
		ret := []byte(value)
		return ret

	case int64:
		ret := make([]byte, 8)
		binary.LittleEndian.PutUint64(ret, uint64(value))
		return ret
	case uint64:
		ret := make([]byte, 8)
		binary.LittleEndian.PutUint64(ret, value)
		return ret

	case int32:
		ret := make([]byte, 8)
		binary.LittleEndian.PutUint32(ret, uint32(value))
		return ret
	case uint32:
		ret := make([]byte, 8)
		binary.LittleEndian.PutUint32(ret, value)
		return ret

	case float32:
		bits := math.Float32bits(value)
		ret := make([]byte, 4)
		binary.LittleEndian.PutUint32(ret, bits)
		return ret

	default:
		break
	}

	return nil
}

func ByteToType(src []byte, v interface{}) {
	switch va := v.(type) {
	case []byte:
		v = src
	case string:
		v = string(src)
	case int:
		v = int(binary.LittleEndian.Uint64(src))
	case int64:
		v = int64(binary.LittleEndian.Uint64(src))
	case uint64:
		v = uint64(binary.LittleEndian.Uint64(src))
	case int32:
		v = int32(binary.LittleEndian.Uint32(src))
	case uint32:
		v = uint32(binary.LittleEndian.Uint32(src))
	case float32:
		v = ByteToFloat32(src)
	case float64:
		v = ByteToFloat64(src)
	default:
		fmt.Printf("%t", va)
		break
	}
	return
}
