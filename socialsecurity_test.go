package rplanlib

import (
	"testing"
)

//
// Testing for socialsecurity.go
//

func TestFra(t *testing.T) {
	tests := []struct {
		age int
		fra int
	}{
		{ // case 0
			age: 10,
			fra: 67,
		},
		{ // case 1
			age: 50,
			fra: 67,
		},
		{ // case 1
			age: 70,
			fra: 66,
		},
		{ // case 1
			age: 80,
			fra: 65,
		},
	}
	for i, elem := range tests {
		fra := fra(elem.age)
		if elem.fra != fra {
			t.Errorf("TestFra case %d: expected %d, found %d\n", i, elem.fra, fra)
		}
	}
}

func TestAdjPIA(t *testing.T) {
	tests := []struct {
		PIA      float64
		fra      int
		startAge int
		adjPIA   float64
	}{
		{ // case 0
			PIA:      10000,
			fra:      67,
			startAge: 67,
			adjPIA:   10000,
		},
		{ // case 1
			PIA:      10000,
			fra:      66,
			startAge: 70,
			adjPIA:   13604,
		},
		{ // case 2
			PIA:      10000,
			fra:      65,
			startAge: 62,
			adjPIA:   8232,
		},
	}
	for i, elem := range tests {
		amt := adjPIA(elem.PIA, elem.fra, elem.startAge)
		if int(elem.adjPIA) != int(amt) {
			t.Errorf("TestAdjPIA case %d: expected %f, found %f\n", i, elem.adjPIA, amt)
		}
	}
}

func TestProcessSS(t *testing.T) {
	tests := []struct {
		ip         map[string]string
		warningmes string
		expectnil  bool
	}{
		{ // case 0
			ip: map[string]string{
				"setName":                    "",
				"filingStatus":               "joint",
				"key1":                       "retiree1",
				"key2":                       "retiree2",
				"eT_Age1":                    "66",
				"eT_Age2":                    "64",
				"eT_RetireAge1":              "66",
				"eT_RetireAge2":              "66",
				"eT_PlanThroughAge1":         "76",
				"eT_PlanThroughAge2":         "76",
				"eT_PIA1":                    "20", //20k
				"eT_PIA2":                    "-1", // spousal benefit
				"eT_SS_Start1":               "66",
				"eT_SS_Start2":               "66",
				"eT_TDRA1":                   "",
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
				"eT_Aftatax":                 "",
				"eT_Aftatax_Rate":            "",
				"eT_Aftatax_Contrib":         "",
				"eT_Aftatax_ContribStartAge": "",
				"eT_Aftatax_ContribEndAge":   "",
				"dollarsInThousands":         "true",
			},
			warningmes: "",
			expectnil:  false,
		},
		{ // case 1
			ip: map[string]string{
				"setName":                    "",
				"filingStatus":               "joint",
				"key1":                       "retiree1",
				"key2":                       "retiree2",
				"eT_Age1":                    "66",
				"eT_Age2":                    "64",
				"eT_RetireAge1":              "66",
				"eT_RetireAge2":              "66",
				"eT_PlanThroughAge1":         "76",
				"eT_PlanThroughAge2":         "76",
				"eT_PIA1":                    "20", //20k
				"eT_PIA2":                    "-1", // spousal benefit
				"eT_SS_Start1":               "66",
				"eT_SS_Start2":               "68",
				"eT_TDRA1":                   "",
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
				"eT_Aftatax":                 "",
				"eT_Aftatax_Rate":            "",
				"eT_Aftatax_Contrib":         "",
				"eT_Aftatax_ContribStartAge": "",
				"eT_Aftatax_ContribEndAge":   "",
				"dollarsInThousands":         "true",
			},
			warningmes: "Warning-SocialSecurityspousalbenefitsdonotincreaseafterFRA,resettingbenefitsstarttoFRA.Pleasecorrectretiree2'sSSageintheconfigurationfileto'66'.",
			expectnil:  false,
		},
		{ // case 2
			ip: map[string]string{
				"setName":                    "",
				"filingStatus":               "joint",
				"key1":                       "retiree1",
				"key2":                       "retiree2",
				"eT_Age1":                    "66",
				"eT_Age2":                    "64",
				"eT_RetireAge1":              "66",
				"eT_RetireAge2":              "66",
				"eT_PlanThroughAge1":         "76",
				"eT_PlanThroughAge2":         "76",
				"eT_PIA1":                    "-1", //Spousal ben
				"eT_PIA2":                    "20", //20k
				"eT_SS_Start1":               "66",
				"eT_SS_Start2":               "70",
				"eT_TDRA1":                   "",
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
				"eT_Aftatax":                 "",
				"eT_Aftatax_Rate":            "",
				"eT_Aftatax_Contrib":         "",
				"eT_Aftatax_ContribStartAge": "",
				"eT_Aftatax_ContribEndAge":   "",
				"dollarsInThousands":         "true",
			},
			warningmes: "Warning - Social Security spousal benefit can only be claimed after the spouse claims benefits. Please correct retiree1's SS age in the configuration file to '72'.",
			expectnil:  false,
		},
		{ // Case 3 // case to match mobile.toml
			ip: map[string]string{
				"setName":                    "activeParams",
				"filingStatus":               "joint",
				"key1":                       "retiree1",
				"key2":                       "retiree2",
				"eT_Age1":                    "54",
				"eT_Age2":                    "54",
				"eT_RetireAge1":              "65",
				"eT_RetireAge2":              "65",
				"eT_PlanThroughAge1":         "75",
				"eT_PlanThroughAge2":         "75",
				"eT_PIA1":                    "",
				"eT_PIA2":                    "",
				"eT_SS_Start1":               "",
				"eT_SS_Start2":               "",
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
				"eT_Aftatax":                 "",
				"eT_Aftatax_Rate":            "",
				"eT_Aftatax_Contrib":         "",
				"eT_Aftatax_ContribStartAge": "",
				"eT_Aftatax_ContribEndAge":   "",

				"eT_iRate":           "2.5",
				"eT_rRate":           "6",
				"eT_maximize":        "Spending", // or "PlusEstate"
				"dollarsInThousands": "true",
			},
			warningmes: "",
			expectnil:  true,
		},
		{ // Case 4 // case to match mobile.toml
			ip: map[string]string{
				"setName":                    "activeParams",
				"filingStatus":               "single",
				"key1":                       "retiree1",
				"key2":                       "retiree2",
				"eT_Age1":                    "54",
				"eT_Age2":                    "",
				"eT_RetireAge1":              "65",
				"eT_RetireAge2":              "",
				"eT_PlanThroughAge1":         "75",
				"eT_PlanThroughAge2":         "",
				"eT_PIA1":                    "20", //20k
				"eT_PIA2":                    "",
				"eT_SS_Start1":               "67",
				"eT_SS_Start2":               "",
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
				"eT_Aftatax":                 "",
				"eT_Aftatax_Rate":            "",
				"eT_Aftatax_Contrib":         "",
				"eT_Aftatax_ContribStartAge": "",
				"eT_Aftatax_ContribEndAge":   "",

				"eT_iRate":           "2.5",
				"eT_rRate":           "6",
				"eT_maximize":        "Spending", // or "PlusEstate"
				"dollarsInThousands": "true",
			},
			warningmes: "",
			expectnil:  false,
		},
		{ // Case 5 // joint with independent ss incomes
			ip: map[string]string{
				"setName":                    "activeParams",
				"filingStatus":               "joint",
				"key1":                       "retiree1",
				"key2":                       "retiree2",
				"eT_Age1":                    "54",
				"eT_Age2":                    "57",
				"eT_RetireAge1":              "65",
				"eT_RetireAge2":              "66",
				"eT_PlanThroughAge1":         "75",
				"eT_PlanThroughAge2":         "80",
				"eT_PIA1":                    "20", //20k
				"eT_PIA2":                    "30",
				"eT_SS_Start1":               "67",
				"eT_SS_Start2":               "70",
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
				"eT_Aftatax":                 "",
				"eT_Aftatax_Rate":            "",
				"eT_Aftatax_Contrib":         "",
				"eT_Aftatax_ContribStartAge": "",
				"eT_Aftatax_ContribEndAge":   "",

				"eT_iRate":           "2.5",
				"eT_rRate":           "6",
				"eT_maximize":        "Spending", // or "PlusEstate"
				"dollarsInThousands": "true",
			},
			warningmes: "",
			expectnil:  false,
		},
	}
	for i, elem := range tests {
		/*
			if i != 5 {
				continue
			}
			fmt.Printf("\nCase %d::\n", i)
		*/
		ip, err := NewInputParams(elem.ip)
		if err != nil {
			t.Errorf("TestProcessSS: %s\n", err)
			continue
		}

		doNothing := false // Turn on/off Stdio redirection
		mychan := make(chan string)
		oldout, w, err := RedirectStdout(mychan, doNothing)
		if err != nil {
			t.Errorf("RedirectStdout: %s\n", err)
			return // should this be continue?
		}
		ss, ss1, ss2, tags := processSS(ip)

		str := RestoreStdout(mychan, oldout, w, doNothing)
		strn := stripWhitespace(str)
		warningmes := stripWhitespace(elem.warningmes)
		if warningmes != strn {
			t.Errorf("TestProcessSS case %d:  expected Warning '%s'\n\tbut found: '%s'\n", i, warningmes, strn)
		}

		if ss == nil {
			if !elem.expectnil {
				t.Errorf("TestProcessSS case %d: Social Security vector (SS) is unexpectedly nil\n", i)
			}
			continue
		}
		/*
			fmt.Printf("ss: %#v\n", ss)
			fmt.Printf("ss1: %#v\n", ss1)
			fmt.Printf("ss2: %#v\n", ss2)
			fmt.Printf("tags: %#v\n", tags)
		*/
		if len(ss) != len(ss1) {
			t.Errorf("TestProcessSS case %d: Social Security vectors are not the same lengths as required\n", i)
		}
		zeros := ip.SSStart1 - ip.StartPlan
		//fmt.Printf("zeros: %d, SSstart1: %d, startPlan: %d, endPlan: %d, planthrough1: %d\n", zeros, ip.SSStart1, ip.StartPlan, ip.EndPlan, ip.planThroughAge1)
		expt := "combined"
		if tags[0] != expt {
			t.Errorf("TestProcessSS case %d:  tags[0] should be (%s) but is (%s)\n", i, expt, tags[0])
		}
		if tags[1] != ip.MyKey1 {
			t.Errorf("TestProcessSS case %d:  tags[1] should be (%s) but is (%s)\n", i, ip.MyKey1, tags[1])
		}
		// Verify years before starting SS have zero SS income
		for j := 0; j < zeros; j++ {
			if ss1[j] != 0 {
				t.Errorf("TestProcessSS case %d:  ss1[%d]: %f should equal zero, it's before starting SS\n", i, j, ss1[j])
			}
		}
		if ip.FilingStatus != "joint" {
			for j := 0; j < len(ss); j++ { // TODO: FIXME if not joint to need ss1 separate from ss - remove it
				if ss[j] != ss1[j] {
					t.Errorf("TestProcessSS case %d:  SS[%d](%f) must equal SS1[%d](%f)\n", i, j, ss[j], j, ss1[j])
				}
			}
		} else { // is "joint"
			if tags[2] != ip.MyKey2 {
				t.Errorf("TestProcessSS case %d:  tags[2] should be (%s) but is (%s)\n", i, ip.MyKey2, tags[2])
			}
			if len(ss) != len(ss2) {
				t.Errorf("TestProcessSS case %d: Social Security vectors are not the same lengths as required\n", i)
			}
			for j := 0; j < len(ss); j++ {
				if ss[j] != ss1[j]+ss2[j] {
					t.Errorf("TestProcessSS case %d:  SS[j] must equal SS1[j] + SS2[j]\n", i)
				}
			}
			//delta := ip.Age2 - ip.Age1
			zeros = ip.SSStart2 + ip.AgeDelta - ip.StartPlan // convert to prime age
			// Verify years before starting SS have zero SS income
			for j := 0; j < zeros; j++ {
				if int(ss2[j]) != 0 {
					t.Errorf("TestProcessSS case %d:  ss2[%d]: %f should equal zero, it's before starting SS\n", i, j, ss2[j])
				}
			}
			// varify that years after retiree's planthrough have zero SS income
			r1end := ip.PlanThroughAge1 - ip.StartPlan + 1
			//fmt.Printf("r1end: %d, planThroughAge1: %d, retireAge1: %d\n", r1end, ip.planThroughAge1, ip.retireAge1)
			for j := r1end; j < ip.Numyr; j++ {
				if ss1[j] != 0 {
					t.Errorf("TestProcessSS case %d:  ss1[%d]: %f should equal zero, it's after planThrough age\n", i, j, ss1[j])
				}
			}
			// varify that years after retiree's planthrough have zero SS income
			r2end := ip.PlanThroughAge2 + ip.AgeDelta - ip.StartPlan + 1
			for j := r2end; j < ip.Numyr; j++ {
				if ss2[j] != 0 {
					t.Errorf("TestProcessSS case %d:  ss2[%d]: %f should equal zero, it's after planThrough age\n", i, j, ss2[j])
				}
			}
			if r1end < r2end {
				// Verify retiree2 gets greater SS after retiree1 is gone
				woulda := ip.IRate * ss1[r1end-1]
				if ss2[r1end] < woulda {
					t.Errorf("TestProcessSS case %d:  ss2[%d]: %f should have gotten %f after spouses death\n", i, r1end, ss2[r1end], woulda)
				}
			} else { // equal case does not muck up this test
				// Verify retiree1 gets greater SS after retiree2 is gone
				woulda := ip.IRate * ss2[r2end-1]
				if ss1[r1end] < woulda {
					t.Errorf("TestProcessSS case %d:  ss1[%d]: %f should have gotten %f after spouses death\n", i, r2end, ss2[r2end], woulda)
				}
			}
		}

		//fmt.Printf("len ss: %d\n", len(ss))
		//fmt.Printf(" ss: %v\n", ss)
		//fmt.Printf("ss1: %v\n", ss1)
		//fmt.Printf("ss2: %v\n", ss2)
	}
}

/*
	tests := []struct {
	}{
		{},
	}
	for i, elem := range tests {
	}

			ip: map[string]string{
				"setName":                    "",
				"filingStatus":               "",
				"key1":                       "retiree1",
				"key2":                       "retiree2",
				"eT_Age1":            "",
				"eT_Age2":            "",
				"eT_RetireAge1":      "",
				"eT_RetireAge2":      "",
				"eT_PlanThroughAge1": "",
				"eT_PlanThroughAge2": "",
				"eT_PIA1":            "",
				"eT_PIA2":            "",
				"eT_SS_Start1":       "",
				"eT_SS_Start2":       "",
				"eT_TDRA1":           "",
				"eT_TDRA2":           "",
				"eT_TDRA_Rate1":         "",
				"eT_TDRA_Rate2":         "",
				"eT_TDRA_Contrib1":   "",
				"eT_TDRA_Contrib2":   "",
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
				"eT_Aftatax":                 "",
				"eT_Aftatax_Rate":            "",
				"eT_Aftatax_Contrib":         "",
				"eT_Aftatax_ContribStartAge": "",
				"eT_Aftatax_ContribEndAge":   "",
			},
*/
