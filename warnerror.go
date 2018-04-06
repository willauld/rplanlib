package rplanlib

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
