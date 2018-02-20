package rplanlib

import (
	"fmt"
	"testing"
)

//
// Testing for input_params.go
//

func TestGetIPIntValue(t *testing.T) {
	tests := []struct {
		str    string
		expect float64
		strerr string
	}{
		{ // case 0
			str:    "",
			expect: 0,
			strerr: "",
		},
		{ // case 1
			str:    "453",
			expect: 453,
			strerr: "",
		},
		{ // case 3
			str:    "453.705",
			expect: 453.705,
			strerr: "strconv.Atoi: parsing \"453.705\": invalid syntax",
		},
	}
	for i, elem := range tests {
		func() {
			defer func() {
				r := recover()
				if r == nil && elem.strerr != "" {
					t.Errorf("getIPIntValue() case %d should have panicked", i)
				} else if elem.strerr == "" && r != nil {
					t.Errorf("getIPIntValue() case %d should not have panicked", i)
				} else if r != nil {
					errstr := fmt.Sprintf("%s", r)
					if errstr != elem.strerr {
						t.Errorf("getIPIntValue() case %d panicked! with err '%v' but should have err '%v'", i, errstr, elem.strerr)
					}
				}
			}()
			// This function may cause a panic
			val := getIPIntValue(elem.str)
			if float64(val) != elem.expect {
				t.Errorf("GetIPIntValue() case %d: Failed - Expected %d but found %d\n", i, int(elem.expect), val)
			}
		}()
	}
}

func TestGetIPFloatValue(t *testing.T) {
	tests := []struct {
		str    string
		expect float64
		strerr string
	}{
		{ // case 0
			str:    "",
			expect: 0,
			strerr: "",
		},
		{ // case 1
			str:    "453",
			expect: 453,
			strerr: "",
		},
		{ // case 3
			str:    "453.705",
			expect: 453.705,
			strerr: "",
		},
		{ // case 3
			str:    "453,705",
			expect: 453.705,
			strerr: "strconv.ParseFloat: parsing \"453,705\": invalid syntax",
		},
	}
	for i, elem := range tests {
		func() {
			defer func() {
				r := recover()
				if r == nil && elem.strerr != "" {
					t.Errorf("getIPFloatValue() case %d should have panicked", i)
				} else if elem.strerr == "" && r != nil {
					t.Errorf("getIPFloatValue() case %d should not have panicked", i)
				} else if r != nil {
					errstr := fmt.Sprintf("%s", r)
					if errstr != elem.strerr {
						t.Errorf("getIPFloatValue() case %d panicked! with err '%v' but should have err '%v'", i, errstr, elem.strerr)
					}
				}
			}()
			// This function may cause a panic
			val := getIPFloatValue(elem.str)
			if val != elem.expect {
				t.Errorf("GetIPFloatValue() case %d: Failed - Expected %f but found %f\n", i, elem.expect, val)
			}
		}()
	}
}
func TestGetIPBoolValue(t *testing.T) {
	tests := []struct {
		str    string
		expect bool
		strerr string
	}{
		{ // case 0
			str:    "",
			expect: false,
			strerr: "",
		},
		{ // case 1
			str:    "trUe",
			expect: true,
			strerr: "",
		},
		{ // case 3
			str:    "faLse",
			expect: false,
			strerr: "",
		},
	}
	for i, elem := range tests {
		func() {
			defer func() {
				r := recover()
				if r == nil && elem.strerr != "" {
					t.Errorf("getIPBoolValue() case %d should have panicked", i)
				} else if elem.strerr == "" && r != nil {
					t.Errorf("getIPBoolValue() case %d should not have panicked", i)
				} else if r != nil {
					errstr := fmt.Sprintf("%s", r)
					if errstr != elem.strerr {
						t.Errorf("getIPBoolValue() case %d panicked! with err '%v' but should have err '%v'", i, errstr, elem.strerr)
					}
				}
			}()
			// This function may cause a panic
			val := getIPBoolValue(elem.str)
			if val != elem.expect {
				t.Errorf("GetIPBoolValue() case %d: Failed - Expected %v but found %v\n", i, elem.expect, val)
			}
		}()
	}
}

func TestNewInputParams(t *testing.T) {
	tests := []struct {
		ip           map[string]string
		prePlanYears int
		startPlan    int
		endPlan      int
		numyr        int
		accmap       map[string]int
	}{
		{ // case 0
			ip: map[string]string{
				"setName":                    "activeParams",
				"filingStatus":               "joint",
				"key1":                       "retiree1",
				"key2":                       "retiree2",
				"eT_Age1":                    "65",
				"eT_Age2":                    "63",
				"eT_RetireAge1":              "66",
				"eT_RetireAge2":              "66",
				"eT_PlanThroughAge1":         "100",
				"eT_PlanThroughAge2":         "100",
				"eT_PIA1":                    "30", // 30k
				"eT_PIA2":                    "-1",
				"eT_SS_Start1":               "70",
				"eT_SS_Start2":               "66",
				"eT_TDRA1":                   "200", // 200k
				"eT_TDRA2":                   "100", // 100k
				"eT_TDRA_Rate1":              "",
				"eT_TDRA_Rate2":              "",
				"eT_TDRA_Contrib1":           "",
				"eT_TDRA_Contrib2":           "",
				"eT_TDRA_ContribStartAge1":   "",
				"eT_TDRA_ContribStartAge2":   "",
				"eT_TDRA_ContribEndAge1":     "",
				"eT_TDRA_ContribEndAge2":     "",
				"eT_Roth1":                   "",
				"eT_Roth2":                   "",
				"eT_Roth_Rate1":              "",
				"eT_Roth_Rate2":              "",
				"eT_Roth_Contrib1":           "",
				"eT_Roth_Contrib2":           "",
				"eT_Roth_ContribStartAge1":   "",
				"eT_Roth_ContribStartAge2":   "",
				"eT_Roth_ContribEndAge1":     "",
				"eT_Roth_ContribEndAge2":     "",
				"eT_Aftatax":                 "50", // 50k
				"eT_Aftatax_Rate":            "7.25",
				"eT_Aftatax_Contrib":         "",
				"eT_Aftatax_ContribStartAge": "",
				"eT_Aftatax_ContribEndAge":   "",
			},
			prePlanYears: 1,
			startPlan:    66,
			endPlan:      103,
			numyr:        37,
			accmap:       map[string]int{"IRA": 2, "roth": 0, "aftertax": 1},
		},
		{ // case 1 // switch retirees
			ip: map[string]string{
				"setName":                    "activeParams",
				"filingStatus":               "joint",
				"key1":                       "retiree1",
				"key2":                       "retiree2",
				"eT_Age1":                    "63",
				"eT_Age2":                    "65",
				"eT_RetireAge1":              "66",
				"eT_RetireAge2":              "66",
				"eT_PlanThroughAge1":         "100",
				"eT_PlanThroughAge2":         "100",
				"eT_PIA1":                    "30", // 30k
				"eT_PIA2":                    "-1",
				"eT_SS_Start1":               "70",
				"eT_SS_Start2":               "66",
				"eT_TDRA1":                   "200", // 200k
				"eT_TDRA2":                   "",
				"eT_TDRA_Rate1":              "",
				"eT_TDRA_Rate2":              "",
				"eT_TDRA_Contrib1":           "",
				"eT_TDRA_Contrib2":           "",
				"eT_TDRA_ContribStartAge1":   "",
				"eT_TDRA_ContribStartAge2":   "",
				"eT_TDRA_ContribEndAge1":     "",
				"eT_TDRA_ContribEndAge2":     "",
				"eT_Roth1":                   "",
				"eT_Roth2":                   "",
				"eT_Roth_Rate1":              "",
				"eT_Roth_Rate2":              "",
				"eT_Roth_Contrib1":           "",
				"eT_Roth_Contrib2":           "",
				"eT_Roth_ContribStartAge1":   "",
				"eT_Roth_ContribStartAge2":   "",
				"eT_Roth_ContribEndAge1":     "",
				"eT_Roth_ContribEndAge2":     "",
				"eT_Aftatax":                 "50", // 50k
				"eT_Aftatax_Rate":            "7.25",
				"eT_Aftatax_Contrib":         "",
				"eT_Aftatax_ContribStartAge": "",
				"eT_Aftatax_ContribEndAge":   "",
			},
			prePlanYears: 1,
			startPlan:    64,
			endPlan:      101,
			numyr:        37,
			accmap:       map[string]int{"IRA": 1, "roth": 0, "aftertax": 1},
		},
		{ // case 2 // switch retirees
			ip: map[string]string{
				"setName":                    "activeParams",
				"filingStatus":               "joint",
				"key1":                       "retiree1",
				"key2":                       "retiree2",
				"eT_Age1":                    "65",
				"eT_Age2":                    "55",
				"eT_RetireAge1":              "65",
				"eT_RetireAge2":              "67",
				"eT_PlanThroughAge1":         "85",
				"eT_PlanThroughAge2":         "87",
				"eT_PIA1":                    "30", // 30k
				"eT_PIA2":                    "-1",
				"eT_SS_Start1":               "70",
				"eT_SS_Start2":               "66",
				"eT_TDRA1":                   "200", // 200k
				"eT_TDRA2":                   "100", // 100k
				"eT_TDRA_Rate1":              "",
				"eT_TDRA_Rate2":              "",
				"eT_TDRA_Contrib1":           "",
				"eT_TDRA_Contrib2":           "",
				"eT_TDRA_ContribStartAge1":   "",
				"eT_TDRA_ContribStartAge2":   "",
				"eT_TDRA_ContribEndAge1":     "",
				"eT_TDRA_ContribEndAge2":     "",
				"eT_Roth1":                   "10", // 10K
				"eT_Roth2":                   "",
				"eT_Roth_Rate1":              "",
				"eT_Roth_Rate2":              "",
				"eT_Roth_Contrib1":           "",
				"eT_Roth_Contrib2":           "",
				"eT_Roth_ContribStartAge1":   "",
				"eT_Roth_ContribStartAge2":   "",
				"eT_Roth_ContribEndAge1":     "",
				"eT_Roth_ContribEndAge2":     "",
				"eT_Aftatax":                 "",
				"eT_Aftatax_Rate":            "7.25",
				"eT_Aftatax_Contrib":         "",
				"eT_Aftatax_ContribStartAge": "",
				"eT_Aftatax_ContribEndAge":   "",
			},
			prePlanYears: 0,
			startPlan:    65,
			endPlan:      98,
			numyr:        33,
			accmap:       map[string]int{"IRA": 2, "roth": 1, "aftertax": 0},
		},
		{ // case 2 // switch retirees
			ip: map[string]string{
				"setName":                    "activeParams",
				"filingStatus":               "single",
				"key1":                       "retiree1",
				"key2":                       "retiree2",
				"eT_Age1":                    "45",
				"eT_Age2":                    "",
				"eT_RetireAge1":              "65",
				"eT_RetireAge2":              "",
				"eT_PlanThroughAge1":         "85",
				"eT_PlanThroughAge2":         "",
				"eT_PIA1":                    "30", // 30k
				"eT_PIA2":                    "-1",
				"eT_SS_Start1":               "70",
				"eT_SS_Start2":               "66",
				"eT_TDRA1":                   "200", // 200k
				"eT_TDRA2":                   "",
				"eT_TDRA_Rate1":              "",
				"eT_TDRA_Rate2":              "",
				"eT_TDRA_Contrib1":           "",
				"eT_TDRA_Contrib2":           "",
				"eT_TDRA_ContribStartAge1":   "",
				"eT_TDRA_ContribStartAge2":   "",
				"eT_TDRA_ContribEndAge1":     "",
				"eT_TDRA_ContribEndAge2":     "",
				"eT_Roth1":                   "10", // 10K
				"eT_Roth2":                   "",
				"eT_Roth_Rate1":              "",
				"eT_Roth_Rate2":              "",
				"eT_Roth_Contrib1":           "",
				"eT_Roth_Contrib2":           "",
				"eT_Roth_ContribStartAge1":   "",
				"eT_Roth_ContribStartAge2":   "",
				"eT_Roth_ContribEndAge1":     "",
				"eT_Roth_ContribEndAge2":     "",
				"eT_Aftatax":                 "",
				"eT_Aftatax_Rate":            "7.25",
				"eT_Aftatax_Contrib":         "",
				"eT_Aftatax_ContribStartAge": "",
				"eT_Aftatax_ContribEndAge":   "",
			},
			prePlanYears: 20,
			startPlan:    65,
			endPlan:      86,
			numyr:        21,
			accmap:       map[string]int{"IRA": 1, "roth": 1, "aftertax": 0},
		},
	}
	for i, elem := range tests {
		modelip, err := NewInputParams(elem.ip)
		if err != nil {
			fmt.Printf("TestNewInputParents: %s\n", err)
			continue
		}
		if modelip.prePlanYears != elem.prePlanYears {
			t.Errorf("NewInputParams case %d: Failed - prePlanYears Expected %v but found %v\n", i, elem.prePlanYears, modelip.prePlanYears)
		}
		if modelip.startPlan != elem.startPlan {
			t.Errorf("NewInputParams case %d: Failed - startPlan Expected %v but found %v\n", i, elem.startPlan, modelip.startPlan)
		}
		if modelip.endPlan != elem.endPlan {
			t.Errorf("NewInputParams case %d: Failed - endPlan Expected %v but found %v\n", i, elem.endPlan, modelip.endPlan)
		}
		if modelip.numyr != elem.numyr {
			t.Errorf("NewInputParams case %d: Failed - numyr Expected %v but found %v\n", i, elem.numyr, modelip.numyr)
		}
		if modelip.accmap["IRA"] != elem.accmap["IRA"] {
			t.Errorf("NewInputParams case %d: Failed - IRA accounts Expected %v but found %v\n", i, elem.accmap["IRA"], modelip.accmap["IRA"])
		}
		if modelip.accmap["roth"] != elem.accmap["roth"] {
			t.Errorf("NewInputParams case %d: Failed - roth accounts Expected %v but found %v\n", i, elem.accmap["roth"], modelip.accmap["roth"])
		}
		if modelip.accmap["aftertax"] != elem.accmap["aftertax"] {
			t.Errorf("NewInputParams case %d: Failed - aftertax accounts Expected %v but found %v\n", i, elem.accmap["aftertax"], modelip.accmap["aftertax"])
		}
	}
}
