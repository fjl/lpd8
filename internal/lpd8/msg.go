package lpd8

import (
	"bytes"
	"fmt"
)

const (
	sysexBegin = 0xf0
	sysexEnd   = 0xf7
)

// Message types.
const (
	msgtypeWriteProg     = 0x61
	msgtypeSetActiveProg = 0x62
	msgtypeReadProg      = 0x63
	msgtypeGetActiveProg = 0x64
)

const (
	programSize    = 57
	programMsgSize = programSize + 1
)

// EncodeWriteProgram creates a SysEx message that writes a program
// to the given index. There is no response for this message type.
// Note: progIndex is 1-based.
func EncodeWriteProgram(progIndex int, data Program) ([]byte, error) {
	if progIndex > 4 {
		return nil, fmt.Errorf("invalid program index")
	}
	prog, err := data.MarshalBinary()
	if err != nil {
		return nil, err
	}
	msg := messageHeader(msgtypeWriteProg, programMsgSize)
	msg = append(msg, byte(progIndex))
	msg = append(msg, prog...)
	msg = append(msg, sysexEnd)
	return msg, nil
}

// EncodeReadProgram creates a 'program read' MIDI message.
// Note: progIndex is 1-based.
func EncodeReadProgram(progIndex int) ([]byte, error) {
	if progIndex > 4 {
		return nil, fmt.Errorf("invalid program index")
	}
	msg := messageHeader(msgtypeReadProg, 1)
	msg = append(msg, byte(progIndex))
	msg = append(msg, sysexEnd)
	return msg, nil
}

// ProgramResponse is the response to a 'program read' message.
type ProgramResponse struct {
	Index   byte
	Program Program
}

// DecodeReadProgramResponse decodes the LPD8's MIDI response to a 'program read' message.
func DecodeReadProgramResponse(msg []byte) (ProgramResponse, error) {
	var r ProgramResponse
	msgdata, err := decodeFrame(msg, msgtypeReadProg, programMsgSize)
	if err != nil {
		return r, err
	}
	r.Index = msgdata[0]
	err = r.Program.UnmarshalBinary(msgdata[1:])
	return r, err
}

// EncodeActiveProgram creates a 'get active program' MIDI message.
func EncodeActiveProgram() []byte {
	msg := messageHeader(msgtypeGetActiveProg, 0)
	msg = append(msg, sysexEnd)
	return msg
}

// DecodeActiveProgramResponse decodes the response to a 'get active program' message.
func DecodeActiveProgramResponse(msg []byte) (int, error) {
	msgdata, err := decodeFrame(msg, msgtypeGetActiveProg, 1)
	if err != nil {
		return 0, err
	}
	return int(msgdata[0]), nil
}

// EncodeSetActiveProgram creates a MIDI message that sets the active program index.
// There is no response for this message type.
func EncodeSetActiveProgram(progIndex int) ([]byte, error) {
	if progIndex > 4 {
		return nil, fmt.Errorf("invalid program index")
	}
	msg := messageHeader(msgtypeSetActiveProg, 1)
	msg = append(msg, byte(progIndex))
	msg = append(msg, sysexEnd)
	return msg, nil
}

func messageHeader(msgType byte, length int) []byte {
	return append(make([]byte, 0, length+8),
		sysexBegin,
		0x47, 0x7f, 0x75, // manufacturer + modelsysexBegin,
		msgType, 0x00, byte(length),
	)
}

func decodeFrame(msg []byte, msgType byte, msgSize int) ([]byte, error) {
	prefix := messageHeader(msgType, msgSize)
	if !bytes.HasPrefix(msg, prefix) {
		return nil, fmt.Errorf("doesn't match LPD8 response prefix")
	}
	if msg[len(msg)-1] != sysexEnd {
		return nil, fmt.Errorf("missing sysex termination byte")
	}
	msg = msg[len(prefix) : len(msg)-1] // drop message frame

	if len(msg) != msgSize {
		return nil, fmt.Errorf("wrong message size %d (want %d)", len(msg), 1)
	}
	return msg, nil
}
