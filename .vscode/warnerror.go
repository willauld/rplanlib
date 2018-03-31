package rplanlib

type warnErrorList struct {
	warnList *string[]
	errorList *string[]
}

func NewWarnErrorList() *warnErrorList {
	return &warnErrorList{
		warnList: &string[]{},
		errorList: nil,
	}
}

func (s *warnErrorList)AppendWarning(warning string) {
	if s != nil && s.warnList != nil {
		*s.warnList = append(*s.warnList, warning)
	}
}

func (s *warnErrorList)AppendError(errorstr string) {
	if s != nil && s.errorList != nil {
		*s.errorList = append(*s.errorList, errorstr)
	}
}

func (s *warnErrorList)GetWarningCount() int {
	if s != nil && s.warnList != nil {
	return len(*s.warnList)
	}
	return 0
}

func (s *warnErrorList)GetErrorCount() int {
	if s != nil && s.errorList != nil {
	return len(*s.errorList)
	}
	return 0
}


func (s *warnErrorList)GetWarning(i int) string {
	if s!=nil && s.warnList != nil {
		if i >= 0 && i < len(*s.warnList) {
			return *s.warnList[i]
		}
	}
}

func (s *warnErrorList)GetError(i int) string {
	if s!=nil && s.errorList != nil {
		if i >= 0 && i < len(*s.errorList) {
			return *s.errorList[i]
		}
	}
}