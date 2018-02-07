package rplanlib

import (
	"fmt"
	"os"
	"strings"
)

type appOutput struct {
	csvFile   *os.File
	tableFile *os.File
}

func (ao *appOutput) NewAppOutput(cvsfile, tablefile string) error {
	var err error
	ao.csvFile = nil
	ao.tableFile = os.Stdout
	if cvsfile != "" {
		ao.csvFile, err = os.Create(cvsfile)
		if err != nil {
			return err
		}
	}
	if tablefile != "" {
		ao.tableFile, err = os.Create(tablefile)
		if err != nil {
			return err
		}
	}
	return nil
}

// How to auto close a file on exit in go?
/*
   def __del__(self):
       if self.csv_file is not None and not self.csv_file.closed:
           self.csv_file.close()
*/

func (ao appOutput) output(str string) { // TODO move to a better place
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
