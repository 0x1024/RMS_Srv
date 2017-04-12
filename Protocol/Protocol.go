package Protocol

import (
	"RMS_Srv/FileSrv"
	ptb "RMS_Srv/Prot_table"
)

func Process(pt ptb.PackTag, rec []byte) {

	switch pt.Pcmd {
	case ptb.Fc_fileTrans:
		FileSrv.FileTransfer(pt, rec)
	case ptb.Fc_fileTranD:
		FileSrv.FileTransfer(pt, rec)
	default:

	}
}
