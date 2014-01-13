package testrfc

/*
#include "sapnwrfc.h"
*/
import "C"

import (
	//	"os"
	"unsafe"
)

type rfcConnectionHandle struct {
	h    C.RFC_CONNECTION_HANDLE
	open bool
}

type rfcAttributes struct {
	Attrs C.RFC_ATTRIBUTES
}

type RfcConnection struct {
	pms        RfcConnectionParameters
	Handle     rfcConnectionHandle
	Attributes rfcAttributes
}

//todo fill ConnectionAttributes
//todo add error checks
func NewRfcConnection(cp *RfcConnectionParameters) (*RfcConnection, *RfcError) {
	ei := new(rfcErrorInfo)
	con := new(RfcConnection)
	h := C.RfcOpenConnection((*C.RFC_CONNECTION_PARAMETER)(unsafe.Pointer(&cp.cparameters[0])), C.unsigned(cp.count), &ei.Errorinfo)
	if h == nil {
		err := NewRfcErrorErrorinfo(&ei.Errorinfo)
		return nil, err

	}
	con.Handle.h = h
	con.Handle.open = true
	return con, nil
}

type RfcConnectionParameters struct {
	Parameters  map[string]string
	cparameters [10]C.RFC_CONNECTION_PARAMETER
	count       int
}

func NewRfcConnectionParameters() *RfcConnectionParameters {
	cp := new(RfcConnectionParameters)
	cp.Parameters = make(map[string]string)
	cp.count = 0
	return cp
}

func (cp *RfcConnectionParameters) add(name, value string) error {
	cp.Parameters[name] = value
	var sapname = NewSapUcFromString(name)
	var sapvalue = NewSapUcFromString(value)
	cp.cparameters[cp.count].name = sapname.str
	cp.cparameters[cp.count].value = sapvalue.str
	cp.count++
	return nil
}

func RfcPing(c *RfcConnection) error {
	if c != nil {
		h := c.Handle
		err := new(rfcErrorInfo)
		C.RfcPing(h.h, &err.Errorinfo)
		return nil
	}
	//fmt.Println("no connections established")
	return nil
}
