package lpd8

import (
	"bytes"
	"encoding/binary"
	"fmt"
)

// Program is the configuration data of an LPD-8 'program'.
type Program struct {
	Channel byte    `json:"channel"` // note: 1-based!
	Pads    [8]Pad  `json:"pads"`
	Knobs   [8]Knob `json:"knobs"`
}

// Pad is the configuration of a pad in a program.
type Pad struct {
	Note   byte `json:"note"`
	PC     byte `json:"pc"`
	CC     byte `json:"cc"`
	Toggle bool `json:"toggle"`
}

func (pad *Pad) toggle() byte {
	if pad.Toggle {
		return 1
	}
	return 0
}

// Knob is the configuration of a single knob in a program.
type Knob struct {
	CC  byte `json:"cc"`
	Min byte `json:"min"`
	Max byte `json:"max"`
}

// MarshalBinary encodes the program.
func (p *Program) MarshalBinary() ([]byte, error) {
	if err := p.Validate(); err != nil {
		return nil, err
	}
	pd := *p
	pd.Channel--
	var b bytes.Buffer
	binary.Write(&b, binary.BigEndian, &pd)
	return b.Bytes(), nil
}

// UnmarshalBinary decodes a program.
func (p *Program) UnmarshalBinary(b []byte) error {
	if len(b) != programSize {
		return fmt.Errorf("wrong encoded program size %d, want %d bytes", len(b), programSize)
	}
	var pd Program
	if err := binary.Read(bytes.NewReader(b), binary.BigEndian, &pd); err != nil {
		return err
	}
	pd.Channel++
	if err := pd.Validate(); err != nil {
		return err
	}
	*p = pd
	return nil
}

// Validate checks for consistency errors.
func (p *Program) Validate() error {
	if p.Channel == 0 || p.Channel > 16 {
		return fmt.Errorf("invalid MIDI channel %d", p.Channel)
	}
	for i, pad := range &p.Pads {
		if pad.CC > 127 {
			return fmt.Errorf("pad %d has invalid CC %d", i, pad.CC)
		}
		if pad.Note > 127 {
			return fmt.Errorf("pad %d has invalid note value %d", i, pad.Note)
		}
		if pad.PC > 127 {
			return fmt.Errorf("pad %d has invalid PC %d", i, pad.PC)
		}
	}
	for i, knob := range &p.Knobs {
		if knob.CC > 127 {
			return fmt.Errorf("knob %d has invalid CC %d", i, knob.CC)
		}
		if knob.Min > 127 {
			return fmt.Errorf("knob %d has invalid min value %d", i, knob.Min)
		}
		if knob.Max > 127 {
			return fmt.Errorf("knob %d has invalid max value %d", i, knob.Max)
		}
		if knob.Min > knob.Max {
			return fmt.Errorf("knob %d has min value > max value", i)
		}
	}
	return nil
}
