package testrfc

/*
#include "sapnwrfc.h"
*/
import "C"

type rfcErrorInfo struct {
	Errorinfo C.RFC_ERROR_INFO
}

type RfcError struct {
	error
	errstr       string
	rfcErrorinfo *C.RFC_ERROR_INFO
}

func (e RfcError) Error() string {
	if e.rfcErrorinfo != nil {
		return rfcSapUcToUtf8(&e.rfcErrorinfo.message[0], 0)
	}
	return e.errstr
}

func NewRfcErrorErrorinfo(ei *C.RFC_ERROR_INFO) *RfcError {
	err := new(RfcError)
	if ei == nil {
		return err
	}
	err.rfcErrorinfo = ei
	return err
}

var RFC_RC = map[C.RFC_RC]string{
	C.RFC_OK: "RFC_OK",
}

type RfcRc struct {
	rc C.RFC_RC
}
