package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"github.com/docopt/docopt-go"
	"github.com/fjl/lpd8/internal/lpd8"
)

var usage = `LPD8 tool.

Usage:
  lpd8 backup      [options] <backup.json>
  lpd8 restore     [options] <backup.json>
  lpd8 read-prog   [options] <program>
  lpd8 write-prog  [options] <program> <prog.json>

Options:
  -h --help      Show this screen.
  -v --verbose   Enable logging.
  --device=<dev> MIDI device name [default: LPD8].
`

type cliOptions struct {
	// Commands.
	CmdBackup  bool `docopt:"backup"`
	CmdRestore bool `docopt:"restore"`
	CmdRead    bool `docopt:"read-prog"`
	CmdWrite   bool `docopt:"write-prog"`

	// Arguments.
	ProgIndex  int    `docopt:"<program>"`
	ProgFile   string `docopt:"<prog.json>"`
	BackupFile string `docopt:"<backup.json>"`

	// Global options.
	Device  string `docopt:"--device"`
	Verbose bool   `docopt:"--verbose"`
}

func main() {
	rawOpt, err := docopt.ParseDoc(usage)
	if err != nil {
		fatal(err)
	}
	var opt cliOptions
	if err := rawOpt.Bind(&opt); err != nil {
		fatal(err)
	}
	if opt.Device == "" {
		opt.Device = "LPD8"
	}
	if opt.Verbose {
		log.SetOutput(os.Stderr)
	} else {
		log.SetOutput(ioutil.Discard)
	}

	conn, err := open(opt.Device)
	if err != nil {
		fatal(err)
	}
	defer conn.close()

	switch {
	case opt.CmdBackup:
		doBackup(conn, opt.BackupFile)
	case opt.CmdRestore:
		doRestore(conn, opt.BackupFile)
	case opt.CmdRead:
		doReadProgram(conn, opt.ProgIndex)
	case opt.CmdWrite:
		doWriteProgram(conn, opt.ProgIndex, opt.ProgFile)
	}
}

type backupJSON struct {
	Programs map[int]lpd8.Program `json:"programs"`
}

func doBackup(c *conn, file string) {
	var data backupJSON
	data.Programs = make(map[int]lpd8.Program)

	for i := 1; i <= 4; i++ {
		prog, err := c.readProgram(i)
		if err != nil {
			fatal(err)
		}
		data.Programs[i] = *prog
	}

	fmt.Println("Writing backup:", file)
	if err := writeJSON(file, &data); err != nil {
		fatal(err)
	}
}

func doRestore(c *conn, file string) {
	var data backupJSON
	if err := readJSON(file, &data); err != nil {
		fatal(err)
	}
	if data.Programs == nil {
		fatal(fmt.Errorf("missing 'programs' key in file"))
	}

	for i := 1; i <= 4; i++ {
		prog, ok := data.Programs[i]
		if !ok {
			fmt.Println("Note: program", i, "is missing in backup file.")
			continue
		}
		if err := c.writeProgram(i, prog); err != nil {
			fatal(err)
		}
	}
	fmt.Println("OK")
}

func doReadProgram(c *conn, progIndex int) {
	prog, err := c.readProgram(progIndex)
	if err != nil {
		fatal(err)
	}
	text, _ := json.MarshalIndent(&prog, "", "  ")
	fmt.Println(string(text))
}

func doWriteProgram(c *conn, progIndex int, file string) {
	var prog lpd8.Program
	if err := readJSON(file, &prog); err != nil {
		fatal(err)
	}
	if err := prog.Validate(); err != nil {
		fatal(err)
	}
	if err := c.writeProgram(progIndex, prog); err != nil {
		fatal(err)
	}
	fmt.Println("OK")
}

func readJSON(file string, v interface{}) error {
	text, err := ioutil.ReadFile(file)
	if err != nil {
		return err
	}
	return json.Unmarshal(text, v)
}

func writeJSON(file string, v interface{}) error {
	text, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		return err
	}
	return ioutil.WriteFile(file, text, 0644)
}

func fatal(err error) {
	fmt.Fprintln(os.Stderr, "Error:", err)
	os.Exit(1)
}
