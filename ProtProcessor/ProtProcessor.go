package ProtProcessor

import (
	"RMS_Srv/FileSrv"
	ptb "RMS_Srv/Protocol"
)

func Process(pt ptb.PackTag, rec []byte) {

	switch pt.Pcmd {
	case ptb.Fc_fileTrans:
		FileSrv.FileReciever(pt, rec)
	case ptb.Fc_fileTranD:
		FileSrv.FileReciever(pt, rec)
	case ptb.Fc_HB:

	default:

	}
}
