package testrfc

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
	cp := makeConnectionParameters()
	_, err := NewRfcConnection(cp)
	if err != nil {
		t.Error(err.Error())
	}
}

func TestRfcPing(t *testing.T) {
	cp := makeConnectionParameters()
	c, err := NewRfcConnection(cp)
	if err != nil {
		t.Error(err.Error())
	}
	RfcPing(c)
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
