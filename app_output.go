package rplanlib

import (
	"fmt"
	"os"
	"strings"
)

// AppOutput holds the info needed for writing output tables and csv
type AppOutput struct {
	csvFile   *os.File
	tableFile *os.File
}

// NewAppOutput creats an initialized AppOutput object
func NewAppOutput(csvfile, tablefile *os.File) AppOutput {
	ao := AppOutput{
		csvFile:   nil,
		tableFile: os.Stdout,
	}
	if csvfile != nil {
		ao.csvFile = csvfile
	}
	if tablefile != nil {
		ao.tableFile = tablefile
	}
	return ao
}

func (ao AppOutput) output(str string) { // TODO move to a better place
	//
	// output writes the information after doing two separate
	// transformations. One for standard out and the other for
	// writing the csv file.
	// For stdout, all '@' are removed and all '&' replaced with
	// a ' '.
	// For cvs, all '@' are replaced with ',' and all '&' are
	// removed.
	// The cvs wrok is done whenever the csv_file handle is not None
	//
	fmt.Fprintf(ao.tableFile, strings.Replace(strings.Replace(str, "@", "", -1), "&", " ", -1))
	if ao.csvFile != nil {
		fmt.Fprintf(ao.csvFile, strings.Replace(strings.Replace(str, "@", ",", -1), "&", "", -1))
	}
}
