package rfc2go

/*
#include "sapnwrfc.h"
*/
import "C"

type rfcErrorInfo struct {
	Errorinfo C.RFC_ERROR_INFO
}

func (ei *rfcErrorInfo) String() string {
	return rfcSapUcToUtf8(&ei.Errorinfo.message[0], 0)
}

type RfcError struct {
	error
	errstr string
	//	rfcErrorinfo *C.RFC_ERROR_INFO
}

func (e RfcError) Error() string {
	return e.errstr
}

func NewRfcErrorErrorinfo(ei *rfcErrorInfo) *RfcError {
	err := new(RfcError)
	if ei == nil {
		return err
	}
	//	err.rfcErrorinfo = ei.Errorinfo
	err.SetErrorInfo(ei)
	return err
}

func (e *RfcError) SetErrorInfo(ei *rfcErrorInfo) {
	if ei != nil {
		e.errstr = ei.String()
	}
}

const (
	RFC_OK                    = C.RFC_OK
	RFC_COMMUNICATION_FAILURE = C.RFC_COMMUNICATION_FAILURE
	RFC_LOGON_FAILURE         = C.RFC_LOGON_FAILURE
	RFC_ABAP_RUNTIME_FAILURE  = C.RFC_ABAP_RUNTIME_FAILURE
)
