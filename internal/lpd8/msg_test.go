package lpd8

import (
	"encoding/hex"
	"encoding/json"
	"reflect"
	"testing"
)

func TestReadProgMessage(t *testing.T) {
	dataHex := "f0477f7563003a01060c0d0e000f1011001213140115161700000102000304050006070800090a0b0018191a1b1c1d1e1f202122232425262728292a2b2c2d2e2ff7"
	data, _ := hex.DecodeString(dataHex)

	var testProg Program
	if err := json.Unmarshal([]byte(testProgramJSON), &testProg); err != nil {
		t.Fatal(err)
	}

	resp, err := DecodeReadProgramResponse(data)
	if err != nil {
		t.Fatal(err)
	}
	if resp.Index != 1 {
		t.Fatalf("wrong progam index %d in decoded response", resp.Index)
	}
	if !reflect.DeepEqual(resp.Program, testProg) {
		t.Fatal("wrong program data in decoded response")
	}
}

var testProgramJSON = `
{
  "channel": 7,
  "pads": [
    {
      "note": 12,
      "pc": 13,
      "cc": 14,
      "toggle": false
    },
    {
      "note": 15,
      "pc": 16,
      "cc": 17,
      "toggle": false
    },
    {
      "note": 18,
      "pc": 19,
      "cc": 20,
      "toggle": true
    },
    {
      "note": 21,
      "pc": 22,
      "cc": 23,
      "toggle": false
    },
    {
      "note": 0,
      "pc": 1,
      "cc": 2,
      "toggle": false
    },
    {
      "note": 3,
      "pc": 4,
      "cc": 5,
      "toggle": false
    },
    {
      "note": 6,
      "pc": 7,
      "cc": 8,
      "toggle": false
    },
    {
      "note": 9,
      "pc": 10,
      "cc": 11,
      "toggle": false
    }
  ],
  "knobs": [
    {
      "cc": 24,
      "min": 25,
      "max": 26
    },
    {
      "cc": 27,
      "min": 28,
      "max": 29
    },
    {
      "cc": 30,
      "min": 31,
      "max": 32
    },
    {
      "cc": 33,
      "min": 34,
      "max": 35
    },
    {
      "cc": 36,
      "min": 37,
      "max": 38
    },
    {
      "cc": 39,
      "min": 40,
      "max": 41
    },
    {
      "cc": 42,
      "min": 43,
      "max": 44
    },
    {
      "cc": 45,
      "min": 46,
      "max": 47
    }
  ]
}
`
