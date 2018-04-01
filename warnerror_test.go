package rplanlib

import (
	"testing"
)

//
// Testing for warnerror.go
//

func TestNewWarnErrorList(t *testing.T) {
	tests := []struct {
		size  int
		key   string
		value string
	}{
		{ //case 0
			size:  255,
			key:   "eT_SS_Start1",
			value: "65",
		},
	}
	for i, elem := range tests {
		//fmt.Printf("----- Case %d -----\n", i)
		l := NewWarnErrorList()
		if l == nil {
			t.Errorf("TestNewWarnErrorList() case %d: Failed - warnErrorList should not be nil\n", i)
		}
		if l.warnList == nil {
			t.Errorf("TestNewWarnErrorList() case %d: Failed - warnList should not be nil\n", i)
		}
		if l.errorList == nil {
			t.Errorf("TestNewWarnErrorList() case %d: Failed - errorList should not be nil\n", i)
		}
		if elem.key == "help" {
			t.Errorf("TestNewWarnErrorList() case %d: Failed - elem.key is %q\n", i, elem.key)
		}
	}
}

func TestAppendWarnErrorList(t *testing.T) {
	tests := []struct {
		str1 string
		str2 string
	}{
		{ //case 0
			str1: "eT_SS_Start1",
			str2: "65",
		},
	}
	l := NewWarnErrorList()
	for i, elem := range tests {
		//fmt.Printf("----- Case %d -----\n", i)
		l.AppendWarning(elem.str1)
		l.AppendWarning(elem.str2)
		c := l.GetWarningCount()
		if c != 2 {
			t.Errorf("TestAppendWarnErrorList() case %d: Failed - GetWarningCount() should have return 2 but returned %d\n", i, c)
		}
		s1 := l.GetWarning(0)
		if s1 != elem.str1 {
			t.Errorf("TestAppendWarnErrorList() case %d: Failed - GetWarning(0) should have returned %q but returned %q\n", i, elem.str1, s1)
		}
		s2 := l.GetWarning(1)
		if s2 != elem.str2 {
			t.Errorf("TestAppendWarnErrorList() case %d: Failed - GetWarning(1) should have returned %q but returned %q\n", i, elem.str2, s2)
		}
		l.AppendError(elem.str1)
		l.AppendError(elem.str2)
		c = l.GetErrorCount()
		if c != 2 {
			t.Errorf("TestAppendWarnErrorList() case %d: Failed - GetErrorCount() should have return 2 but returned %d\n", i, c)
		}
		s1 = l.GetError(0)
		if s1 != elem.str1 {
			t.Errorf("TestAppendWarnErrorList() case %d: Failed - GetError(0) should have returned %q but returned %q\n", i, elem.str1, s1)
		}
		s2 = l.GetError(1)
		if s2 != elem.str2 {
			t.Errorf("TestAppendWarnErrorList() case %d: Failed - GetError(1) should have returned %q but returned %q\n", i, elem.str2, s2)
		}
	}
}
