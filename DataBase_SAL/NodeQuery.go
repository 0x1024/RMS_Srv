package DataBase_SAL

import (
	"fmt"
)

func QueryMCUID(mcuid []byte) uint64 {
	fmt.Println(mcuid)
	var re uint64 = 0
	sa1 := new(Pd_index)
	for k, v := range mcuid {
		re = re + uint64(v)*(1<<(12-uint64(k)))
	}
	n, err := engine.Where("mcuid=?", re).Get(sa1)
	if err != nil {
		fmt.Println(err)
	}
	if n == true {
		return sa1.Pid

	} else {

	}

	return uint64(1)

}
