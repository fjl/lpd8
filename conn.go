package main

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/fjl/lpd8/internal/lpd8"
	"gitlab.com/gomidi/midi"
	driver "gitlab.com/gomidi/rtmididrv"
)

type conn struct {
	in  midi.In
	out midi.Out

	packetCh chan []byte
	closeCh  chan struct{}
}

func open(devname string) (*conn, error) {
	in, out, err := findDevice(devname)
	if err != nil {
		return nil, err
	}
	log.Println("midi input:", in)
	log.Println("midi output:", out)
	if err := in.Open(); err != nil {
		return nil, fmt.Errorf("can't open MIDI input: %v", err)
	}
	if err := out.Open(); err != nil {
		in.Close()
		return nil, fmt.Errorf("can't open MIDI output: %v", err)
	}

	c := &conn{
		in:       in,
		out:      out,
		packetCh: make(chan []byte, 256),
		closeCh:  make(chan struct{}),
	}
	in.SetListener(func(msg []byte, deltaT int64) {
		if !isSysex(msg) {
			return
		}
		select {
		case c.packetCh <- msg:
		case <-c.closeCh:
		}
	})
	return c, nil
}

func (c *conn) close() {
	close(c.closeCh)
	c.in.Close()
	c.out.Close()
}

func findDevice(devname string) (midi.In, midi.Out, error) {
	drv, err := driver.New()
	if err != nil {
		return nil, nil, err
	}
	inputs, err := drv.Ins()
	if err != nil {
		return nil, nil, fmt.Errorf("can't list MIDI inputs: %v", err)
	}
	outputs, err := drv.Outs()
	if err != nil {
		return nil, nil, fmt.Errorf("can't list MIDI outputs: %v", err)
	}
	if len(inputs) == 0 {
		return nil, nil, fmt.Errorf("no MIDI inputs")
	}

	// Find a matching input device.
	var selectedIn midi.In
	if devname == "" {
		selectedIn = inputs[0]
	} else {
		var inputNames []string
		for _, in := range inputs {
			name := in.String()
			inputNames = append(inputNames, name)
			if strings.Contains(strings.ToLower(name), strings.ToLower(devname)) {
				selectedIn = in
				break
			}
		}
		if selectedIn == nil {
			return nil, nil, fmt.Errorf("can't find MIDI input device %q, have %v", devname, inputNames)
		}
	}

	// Find the output device matching input.
	var selectedOut midi.Out
	var outputNames []string
	for _, out := range outputs {
		outputNames = append(outputNames, out.String())
		if out.String() == selectedIn.String() {
			selectedOut = out
			break
		}
	}
	if selectedOut == nil {
		return nil, nil, fmt.Errorf("can't find MIDI output device %q, have %v", selectedIn.String(), outputNames)
	}

	// Found it.
	return selectedIn, selectedOut, nil
}

// readProgram retrieves the program at progIndex.
func (c *conn) readProgram(progIndex int) (prog *lpd8.Program, err error) {
	requestMsg, err := lpd8.EncodeReadProgram(progIndex)
	if err != nil {
		return nil, err
	}

	log.Println("requesting program", progIndex)
	err = c.sysexReqResp(requestMsg, func(msg []byte) bool {
		resp, err := lpd8.DecodeReadProgramResponse(msg)
		if err != nil {
			log.Println("ignoring message:", err)
			return false
		}
		if int(resp.Index) != progIndex {
			log.Println("ignoring response for wrong program index", resp.Index)
			return false
		}
		log.Println("got program", progIndex)
		prog = &resp.Program
		return true
	})
	return prog, err
}

// writeProgram writes a program to the LPD8.
func (c *conn) writeProgram(progIndex int, prog lpd8.Program) error {
	msg, err := lpd8.EncodeWriteProgram(progIndex, prog)
	if err != nil {
		return err
	}
	log.Println("writing program", progIndex)
	_, err = c.out.Write(msg)
	return err
}

func (c *conn) sysexReqResp(requestMsg []byte, decode func(msg []byte) bool) error {
	// Write the request.
	if _, err := c.out.Write(requestMsg); err != nil {
		return err
	}

	// Wait for response.
	timeout := time.NewTimer(3 * time.Second)
	defer timeout.Stop()
	for {
		select {
		case packet := <-c.packetCh:
			if decode(packet) {
				return nil
			}
		case <-timeout.C:
			return fmt.Errorf("timed out waiting for valid LPD-8 response")
		}
	}
}

func isSysex(msg []byte) bool {
	return len(msg) > 0 && msg[0] == 0xf0 && msg[len(msg)-1] == 0xf7
}
