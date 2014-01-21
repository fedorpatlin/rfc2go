package rfc2go

import (
	"fmt"
	"testing"
	"time"
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
	err := c.Open()
	if err != nil {
		t.Error(err.Error())
	}
	err = c.Close()
	if err != nil {
		t.Error(err.Error())
	}
}

func TestRfcPing(t *testing.T) {
	c := getRfcConnection(makeConnectionParameters())
	err := c.Open()
	if err != nil {
		t.Error(err.Error())
	}
	err = c.ListenAndDispatch(2)
	if err != nil {
		t.Error(err.Error())
	}
	err = RfcPing(c)
	if err != nil {
		t.Error(err.Error())
	}
	err = c.Close()
	if err != nil {
		t.Error(err.Error())
	}
}

func TestCallback(t *testing.T) {

	ei := new(rfcErrorInfo)

	err := RfcInstallPasswordChangeHandler(ei)
	if err != nil {
		t.Error(ei.String())
	}
	c := getRfcConnection(makeConnectionParameters())

	err = c.Open()

	if err != nil {
		time.Sleep(2 * time.Second)
		c.Close()
		t.Error(err.Error())
	}
	err = c.Close()
	if err != nil {
		t.Error(err.Error())
	}
}

func getRfcConnection(cp *RfcConnectionParameters) *RfcConnection {
	c, _ := NewRfcConnection(cp)
	return c
}

func makeConnectionParameters() *RfcConnectionParameters {
	cp := NewRfcConnectionParameters()
	cp.Add("ASHOST", "10.142.70.101")
	cp.Add("sysnr", "00")
	cp.Add("client", "001")
	cp.Add("user", "dev")
	cp.Add("lang", "ru")
	cp.Add("passwd", "123456789")
	return cp
}
