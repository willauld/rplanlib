package rplanlib

import (
	"fmt"
	"os"
)

type WarnErrorList struct {
	warnList  *[]string
	errorList *[]string
}

func NewWarnErrorList() *WarnErrorList {
	return &WarnErrorList{
		warnList:  &[]string{},
		errorList: &[]string{},
	}
}

func (s *WarnErrorList) AppendWarning(warning string) {
	if s != nil && s.warnList != nil {
		*s.warnList = append(*s.warnList, warning)
	}
}

func (s *WarnErrorList) AppendError(errorstr string) {
	if s != nil && s.errorList != nil {
		*s.errorList = append(*s.errorList, errorstr)
	}
}

func (s *WarnErrorList) GetWarningCount() int {
	if s != nil && s.warnList != nil {
		return len(*s.warnList)
	}
	return 0
}

func (s *WarnErrorList) GetErrorCount() int {
	if s != nil && s.errorList != nil {
		return len(*s.errorList)
	}
	return 0
}

func (s *WarnErrorList) GetWarning(i int) string {
	if s != nil && s.warnList != nil {
		if i >= 0 && i < len(*s.warnList) {
			return (*s.warnList)[i]
		}
	}
	return ""
}

func (s *WarnErrorList) GetError(i int) string {
	if s != nil && s.errorList != nil {
		if i >= 0 && i < len(*s.errorList) {
			return (*s.errorList)[i]
		}
	}
	return ""
}

func (s *WarnErrorList) ClearWarnings() {
	s.warnList = &[]string{}
}

func (s *WarnErrorList) ClearErrors() {
	s.errorList = &[]string{}
}

func PrintAndClearMsg(f *os.File, msgList *WarnErrorList) {
	// FIXME TODO
	ec := msgList.GetErrorCount()
	if ec > 0 {
		fmt.Fprintf(f, "%d Error(s) found:\n", ec)
		for i := 0; i < ec; i++ {
			fmt.Fprintf(f, "%s\n", msgList.GetError(i))
		}
	}
	msgList.ClearErrors()
	//

	wc := msgList.GetWarningCount()
	if wc > 0 {
		fmt.Fprintf(f, "%d Warning(s) found:\n", wc)
		for i := 0; i < wc; i++ {
			fmt.Fprintf(f, "%s\n", msgList.GetWarning(i))
		}
	}
	msgList.ClearWarnings()
}
