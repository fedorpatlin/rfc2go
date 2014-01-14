package rfc2go

import (
	"fmt"
	"testing"
)

func TestRfcInit(t *testing.T) {
	err := RfcInit()
	if err != nil {
		t.Error(err.Error())
	}
}

func TestRfcGetVersion(t *testing.T) {
	ver, err := RfcGetVersion()
	if err != nil {
		t.Error(err.Error())
	}
	fmt.Println(ver.VersionString)
}

func TestOpenConnection(t *testing.T) {
	c := getRfcConnection(makeConnectionParameters())
	err := c.Connect()
	if err != nil {
		t.Error(err.Error())
	}
}

func TestRfcPing(t *testing.T) {
	c := getRfcConnection(makeConnectionParameters())
	c.Connect()
	rc := RfcPing(c)
	if rc != nil {
		t.Error(rc.Error())
	}
}

func getRfcConnection(cp *RfcConnectionParameters) *RfcConnection {
	c, _ := NewRfcConnection(cp)
	return c
}

func makeConnectionParameters() *RfcConnectionParameters {
	cp := NewRfcConnectionParameters()
	cp.add("ASHOST", "10.142.70.101")
	cp.add("sysnr", "00")
	cp.add("client", "100")
	cp.add("user", "dev")
	cp.add("lang", "ru")
	cp.add("passwd", "1234")
	return cp
}
