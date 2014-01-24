package rfc2go

/*
#include "sapnwrfc.h"

static uint mycpy (SAP_UC* to, SAP_UC* from, unsigned len){
	uint c = 0;
	do {
		*to = *from;
		to++;
		from++;
		len--;
		c++;
	} while(len);
	return c;
}

extern RFC_RC cb(void * data);

extern RFC_RC Cbcall(SAP_UC* p0,
				   SAP_UC* p1,
				   SAP_UC* p2,
				   SAP_UC* p3,
				   unsigned p4,
				   SAP_UC* p5,
				   unsigned p6,
				   RFC_ERROR_INFO* p7);

static RFC_RC registerCallback(RFC_ERROR_INFO* ei){
	return RfcInstallPasswordChangeHandler((RFC_ON_PASSWORD_CHANGE)Cbcall, ei);
}

*/
import "C"

import (
	"fmt"
	"time"
	"unsafe"
)

type rfcConnectionHandle struct {
	h C.RFC_CONNECTION_HANDLE
}

func NewRfcConnectionHandle(ch C.RFC_CONNECTION_HANDLE) *rfcConnectionHandle {
	h := new(rfcConnectionHandle)
	h.h = ch
	return h
}

type rfcAttributes struct {
	Attrs C.RFC_ATTRIBUTES
}

type RfcConnection struct {
	pms        *RfcConnectionParameters
	Handle     *rfcConnectionHandle
	Attributes *rfcAttributes
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

func (c *RfcConnection) Open() *RfcError {
	ei := new(rfcErrorInfo)
	fmt.Printf("param name %s , value %s\n", NewSapUc(c.pms.cparameters[2].name, 0), NewSapUc(c.pms.cparameters[2].value, 0))
	h := C.RfcOpenConnection((*C.RFC_CONNECTION_PARAMETER)(unsafe.Pointer(&c.pms.cparameters[0])), C.unsigned(c.pms.count), &ei.Errorinfo)
	if h == nil {
		err := NewRfcErrorErrorinfo(ei)
		return err
	}
	c.Handle = NewRfcConnectionHandle(h)
	c.Opened = true

	return nil
}

func (c *RfcConnection) ListenAndDispatch(timeout int) *RfcError {
	ei := new(rfcErrorInfo)
	if !c.Opened {
		err := new(RfcError)
		err.errstr = "Connection is closed"
		return err
	}
	rc := C.RfcListenAndDispatch(c.Handle.h, C.int(timeout), &ei.Errorinfo)
	if rc != RFC_OK {
		return NewRfcErrorErrorinfo(ei)
	}
	return nil
}

func (c *RfcConnection) Close() *RfcError {
	ei := new(rfcErrorInfo)
	if !c.Opened {
		return nil
	}
	rc := C.RfcCloseConnection(c.Handle.h, &ei.Errorinfo)
	if rc != RFC_OK {
		return NewRfcErrorErrorinfo(ei)
	}
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

func (cp *RfcConnectionParameters) Add(name, value string) error {
	cp.Parameters[name] = value
	var sapname = NewSapUcFromString(name)
	var sapvalue = NewSapUcFromString(value)
	cp.cparameters[cp.count].name = &(sapname.str)[0]
	cp.cparameters[cp.count].value = &(sapvalue.str)[0]
	cp.count++
	return nil
}

func RfcPing(c *RfcConnection) *RfcError {
	if c.Opened {
		h := c.Handle
		ei := new(rfcErrorInfo)
		rc := C.RfcPing(h.h, &ei.Errorinfo)
		if RFC_OK != rc {
			return NewRfcErrorErrorinfo(ei)
		}
	} else {
		return &RfcError{errstr: "Error: Closed connection"}
	}
	return nil
}

//export cb
func cb(data unsafe.Pointer) C.RFC_RC {
	return RFC_OK
}

//export Cbcall
func Cbcall(sysid, user, client, password *C.SAP_UC, pwlen C.unsigned, newpassword *C.SAP_UC, newpwlen C.unsigned, ei *C.RFC_ERROR_INFO) C.RFC_RC {
	fmt.Println("callback called")
	npwd := NewSapUcFromString("987654321")
	pwd := NewSapUcFromString("123456789")
	copy(*(*[]C.SAP_UC)(unsafe.Pointer(password)), pwd.str)
	copy(*(*[]C.SAP_UC)(unsafe.Pointer(newpassword)), npwd.str)
	pwlen = C.unsigned(pwd.length)
	newpwlen = C.unsigned(npwd.length)
	go func(v1, v2 *C.SAP_UC, v3, v4 C.unsigned) {
		time.Sleep(2 * time.Second)
		fmt.Printf("%v %d %v %d\n", v1, v3, v2, v4)

	}(password, newpassword, pwlen, newpwlen)
	time.Sleep(1 * time.Second)
	fmt.Println("callback returned")
	return C.RFC_OK
}

var f = Cbcall

func RfcInstallPasswordChangeHandler(ei *rfcErrorInfo) *RfcError {
	rc := C.registerCallback(&ei.Errorinfo)
	if rc != RFC_OK {
		return NewRfcErrorErrorinfo(ei)
	}
	return nil
}
