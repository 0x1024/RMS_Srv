package Protocol

type PackTag struct {
	Phead   uint16
	Plen    uint16
	Pserial uint16
	Pcmd    uint16
	Ppara   uint32
}

//frame cmd type list
const (
	Fc_fileTrans = 0x10
	Fc_fileTranD = 0x11
	Fc_dataTrans = 0x20

	Fc_HB = 0xF0
)

// fc file paras
const (
	Fcp_fileName = 0x01
	Fcp_fileEOF  = 0x02
	Fcp_fileSize = 0x03
	Fcp_filedata = 0x04
)

//tcp sender cmd
const (
	TSC_SendFile = 0x10
)
