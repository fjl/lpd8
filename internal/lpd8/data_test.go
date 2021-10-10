package lpd8

import (
	"bytes"
	"encoding/hex"
	"testing"
)

func TestBinary(t *testing.T) {
	var prog = Program{Channel: 8}
	prog.Pads[3].CC = 10
	prog.Pads[4].Note = 44
	prog.Pads[6].PC = 22
	prog.Knobs[4].CC = 99
	prog.Knobs[4].Min = 10
	prog.Knobs[4].Max = 50
	want := "0700000000000000000000000000000a002c000000000000000016000000000000000000000000000000000000630a32000000000000000000"

	// Check encoding.
	enc, err := prog.MarshalBinary()
	if err != nil {
		t.Fatal("MarshalBinary error:", err)
	}
	t.Logf("enc: %x", enc)
	if bin, _ := hex.DecodeString(want); !bytes.Equal(enc, bin) {
		t.Errorf("wrong encoding: %x", bin)
	}

	// Check decoding.
	var dec Program
	if err := dec.UnmarshalBinary(enc); err != nil {
		t.Fatal("UnmarshalBinary error:", err)
	}
	if dec != prog {
		t.Errorf("decoded program does not match encoder input")
	}
}
