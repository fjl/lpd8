package lpd8

import (
	"encoding/hex"
	"testing"
)

func TestReadProgMessage(t *testing.T) {
	dataHex := "f0477f7563003a01060c0d0e000f1011001213140115161700000102000304050006070800090a0b0018191a1b1c1d1e1f202122232425262728292a2b2c2d2e2ff7"
	data, _ := hex.DecodeString(dataHex)

	resp, err := DecodeReadProgramResponse(data)
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("%+v", resp)
}
