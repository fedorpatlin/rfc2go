// testrfc project testrfc.go
package rfc2go

/*
#cgo windows LDFLAGS: -L"d:/saprfc/rfcsdk/lib/"  -lsapnwrfc
#cgo windows CFLAGS:  -I"D:/saprfc/rfcsdk/include" -DSAPwithUNICODE -DSAPOnNT
#cgo linux LDFLAGS: -L/home/fedor/Downloads/rfcsdk/nwrfcsdk/lib -lsapnwrfc -lsapucum
#cgo linux CFLAGS: -I/home/fedor/Downloads/rfcsdk/nwrfcsdk/include -DSAPwithUNICODE
#include <stdlib.h>
#include <stdio.h>

#include "sapnwrfc.h"


//SAP_UC * mycu (SAP_UC * str){
//	return cU(str);
//}

RFC_BYTE * pchar_to_prfc_byte(char * from){
        return (RFC_BYTE*) from;
}

unsigned rfcuclength(SAP_UC *str){
	SAP_UC *tmp = str;
	if (tmp == 0){
		return 0;
	}
	int i=0;
	while (*tmp) {
		tmp++;
		i++;
	}
	return i;
}
*/
import "C"

import (
	//	"fmt"
	"os"
	"unicode/utf16"
	"unsafe"
)

type SapUc struct {
	str    []C.SAP_UC
	length uint
}

func zerosapchar() C.SAP_UC {
	var c C.SAP_UC
	l := unsafe.Sizeof(c)
	b := make([]byte, l)
	for i, _ := range b {
		b[i] = 0x0
	}
	c = *(*C.SAP_UC)(unsafe.Pointer(&b[0]))
	return c
}

func sapstrlen(sapstr *C.SAP_UC) ([]C.SAP_UC, uint) {
	zero := zerosapchar()
	sapstrbuf := (*[1000]C.SAP_UC)(unsafe.Pointer(sapstr))
	//	fmt.Printf("received %v\n", sapstr)
	for i, v := range *sapstrbuf {
		if v == zero {
			return sapstrbuf[:i+1], uint(i + 1)
		}
	}
	return make([]C.SAP_UC, 0), 0
}

func sapstrng(str *C.SAP_UC, strlen int) []C.SAP_UC {
	sapstrbuf := (*[1000]C.SAP_UC)(unsafe.Pointer(str))
	return sapstrbuf[:strlen+1]
}

func NewSapUc(str *C.SAP_UC, length uint) *SapUc {
	if str == nil {
		return nil
	}
	s := new(SapUc)
	var sapstr []C.SAP_UC
	if length == 0 {
		sapstr, length = sapstrlen(str)
	} else {
		sapstr = sapstrng(str, int(length))
	}
	s.str = make([]C.SAP_UC, int(length))
	copy(s.str, sapstr)
	//fmt.Printf("array %v of %d items\n", sapstr, length)
	s.length = length
	return s
}

func NewSapUcFromString(str string) *SapUc {
	s := new(SapUc)
	var err *RfcError
	s = rfcUtf8ToSapUc(str)
	if err != nil {
		return nil
	}
	return s
}

func (s *SapUc) String() string {
	return rfcSapUcToUtf8(s)
}

func rfcUtf16ToSapUc(ustr []uint16) *C.SAP_UC {
	//l := len(ustr)
	return (*C.SAP_UC)(unsafe.Pointer(&(append(ustr, 0)[0])))
}

//разобраться с этим мусором
func rfcUtf8ToSapUc(s string) *SapUc {
	var ucsize, uclength C.unsigned
	var ei = new(rfcErrorInfo)
	cs := (*C.RFC_BYTE)(unsafe.Pointer(C.CString(s)))
	defer C.free(unsafe.Pointer(cs))
	//буфер для строки *SAP_UC
	//размер буфера для результата *SAP_UC
	ucsize = C.uint(uint(len(s) * 2))
	sapstring := make([]C.SAP_UC, int(ucsize))
	result := (*C.SAP_UC)(unsafe.Pointer(&sapstring[0]))
	rc := C.RfcUTF8ToSAPUC(cs, C.uint(len(s)), result, &ucsize, &uclength, &ei.Errorinfo)
	if RFC_OK != rc {
		return nil
	}
	return NewSapUc(result, 0)
}

func trydecode(str *C.SAP_UC, l uint) string {
	buf := make([]uint16, 0, 256)
	var i uint = 0
	if str != nil {
		for s := uintptr(unsafe.Pointer(str)); ; s += 2 {
			u := *(*uint16)(unsafe.Pointer(s))
			if i >= l {
				return string(utf16.Decode(buf))
			}
			i++
			buf = append(buf, u)
		}
	}
	return string(utf16.Decode(buf))
}

func rfcSapUcToUtf8(ucstr *SapUc) string {
	if ucstr == nil || ucstr.length == 0 {
		return ""
	}
	var uclength, utf8bufsize, utf8length uint
	err := new(rfcErrorInfo)
	uclength = ucstr.length
	utf8bufsize = uclength * 2
	utf8buf := make([]byte, utf8bufsize)

	utf8bufp := (*C.RFC_BYTE)(&utf8buf[0])

	rc := C.RfcSAPUCToUTF8(&ucstr.str[0], C.unsigned(uclength), utf8bufp, (*C.unsigned)(unsafe.Pointer(&utf8bufsize)), (*C.unsigned)(unsafe.Pointer(&utf8length)), &err.Errorinfo)
	if RFC_OK != rc {
		return ""
	}
	//fmt.Printf("utf8buffer %v of len %d\n", utf8buf[:utf8length], utf8length)
	out := string(utf8buf[:utf8length-1])
	return out
}

func RfcInit() error {
	os.Setenv("SAP_CODEPAGE", "4103") //SAP Note 1021459
	os.Setenv("RFC_TRACE", "3")
	var rc C.RFC_RC
	rc = C.RfcInit()
	if rc != C.RFC_OK {
		err := RfcError{errstr: "Error while initializing library"}
		return err
	}
	return nil
}

type RfcVersion struct {
	maj, min, patch uint
	VersionString   string
}

func RfcGetVersion() (*RfcVersion, error) {
	var maj, min, patch C.unsigned
	var ver = new(RfcVersion)
	rc := NewSapUc(C.RfcGetVersion(&maj, &min, &patch), 0)
	if rc == nil {
		err := RfcError{errstr: "Error getting library version"}
		return nil, err
	}
	ver.VersionString = rc.String()
	ver.maj = uint(maj)
	ver.min = uint(min)
	ver.patch = uint(patch)
	return ver, nil
}
