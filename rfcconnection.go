package rfc2go

/*
#include "sapnwrfc.h"
*/
import "C"

import (
	//	"os"
	"unsafe"
)

type rfcConnectionHandle struct {
	h C.RFC_CONNECTION_HANDLE
}

type rfcAttributes struct {
	Attrs C.RFC_ATTRIBUTES
}

type RfcConnection struct {
	pms        *RfcConnectionParameters
	Handle     rfcConnectionHandle
	Attributes rfcAttributes
	Opened     bool
}

//todo fill ConnectionAttributes
//todo add error checks
func NewRfcConnection(cp *RfcConnectionParameters) (*RfcConnection, *RfcError) {
	con := new(RfcConnection)
	con.pms = cp
	con.Opened = false
	return con, nil
}

func (c *RfcConnection) Connect() *RfcError {
	ei := new(rfcErrorInfo)
	h := C.RfcOpenConnection((*C.RFC_CONNECTION_PARAMETER)(unsafe.Pointer(&c.pms.cparameters[0])), C.unsigned(c.pms.count), &ei.Errorinfo)
	if h == nil {
		err := NewRfcErrorErrorinfo(ei)
		return err

	}
	c.Handle.h = h
	c.Opened = false
	return nil
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

func RfcPing(c *RfcConnection) *RfcError {
	if c.Opened == true {
		h := c.Handle
		ei := new(rfcErrorInfo)
		rc := C.RfcPing(h.h, &ei.Errorinfo)
		if RFC_OK != rc {
			return NewRfcErrorErrorinfo(ei)
		}
	}
	return nil
}
