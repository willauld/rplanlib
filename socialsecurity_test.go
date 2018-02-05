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
		ip       map[string]string
		retirees []retiree
	}{
		{ // case 0
			ip: map[string]string{
				"setName":                    "",
				"filingStatus":               "joint",
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
			},
			retirees: []retiree{
				{ // retireeindx == 0
					age:        66,
					ageAtStart: 66,
					throughAge: 76,
					mykey:      "retiree1",
					definedContributionPlan: false,
					dcpBuckets:              nil,
				},
				{ // retireeindx == 1
					age:        64,
					ageAtStart: 64,
					throughAge: 76,
					mykey:      "retiree2",
					definedContributionPlan: false,
					dcpBuckets:              nil,
				},
			},
		},
		{ // case 1
			ip: map[string]string{
				"setName":                    "",
				"filingStatus":               "joint",
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
			},
			retirees: []retiree{
				{ // retireeindx == 0
					age:        66,
					ageAtStart: 66,
					throughAge: 76,
					mykey:      "retiree1",
					definedContributionPlan: false,
					dcpBuckets:              nil,
				},
				{ // retireeindx == 1
					age:        64,
					ageAtStart: 64,
					throughAge: 76,
					mykey:      "retiree2",
					definedContributionPlan: false,
					dcpBuckets:              nil,
				},
			},
		},
	}
	for i, elem := range tests {
		ip := NewInputParams(elem.ip)
		iRate := 1.03
		ss, ss1, ss2, err := processSS(ip, elem.retirees, iRate)
		if err != nil {
			t.Errorf("TestProcessSS case %d: %s\n", i, err)
		}
		if ip.filingStatus != "joint" {
			if len(ss) != len(ss1) {
				t.Errorf("TestProcessSS case %d: Social Security vectors are not the same lens as required\n", i)
			}
			for j := 0; j < len(ss); j++ {
				if ss[j] != ss1[j] {
					t.Errorf("TestProcessSS case %d:  SS[j] must equal SS1[j]\n", i)

				}
			}
		} else {
			if len(ss) != len(ss1) || len(ss) != len(ss2) {
				t.Errorf("TestProcessSS case %d: Social Security vectors are not the same lens as required\n", i)
			}
			for j := 0; j < len(ss); j++ {
				if ss[j] != ss1[j]+ss2[j] {
					t.Errorf("TestProcessSS case %d:  SS[j] must equal SS1[j] + SS2[j]\n", i)

				}
			}
		}
		zeros := ip.SSStart1 - ip.startPlan
		// Verify years before starting SS have zero SS income
		for j := 0; j < zeros; j++ {
			if ss1[j] != 0 {
				t.Errorf("TestProcessSS case %d:  ss1[%d]: %f should equal zero, it's before starting SS\n", i, j, ss1[j])
			}
		}
		// varify that years after retiree's planthrough have zero SS income
		r1end := ip.planThroughAge1 - elem.retirees[0].ageAtStart + 1
		for j := r1end; j < ip.numyr; j++ {
			if ss1[j] != 0 {
				t.Errorf("TestProcessSS case %d:  ss1[%d]: %f should equal zero, it's after planThrough age\n", i, j, ss1[j])
			}
		}

		delta := ip.age2 - ip.age1
		zeros = ip.SSStart2 - delta - ip.startPlan // convert to prime age
		// Verify years before starting SS have zero SS income
		for j := 0; j < zeros; j++ {
			if int(ss2[j]) != 0 {
				t.Errorf("TestProcessSS case %d:  ss2[%d]: %f should equal zero, it's before starting SS\n", i, j, ss2[j])
			}
		}
		// varify that years after retiree's planthrough have zero SS income
		r2end := ip.planThroughAge2 - elem.retirees[1].ageAtStart + 1
		for j := r2end; j < ip.numyr; j++ {
			if ss2[j] != 0 {
				t.Errorf("TestProcessSS case %d:  ss2[%d]: %f should equal zero, it's after planThrough age\n", i, j, ss2[j])
			}
		}
		if r1end < r2end {
			// Verify retiree2 gets greater SS after retiree1 is gone
			woulda := iRate * ss1[r1end-1]
			if ss2[r1end] < woulda {
				t.Errorf("TestProcessSS case %d:  ss2[%d]: %f should have gotten %f after spouses death\n", i, r1end, ss2[r1end], woulda)
			}
		} else { // equal case does not muck up this test
			// Verify retiree1 gets greater SS after retiree2 is gone
			woulda := iRate * ss2[r2end-1]
			if ss1[r1end] < woulda {
				t.Errorf("TestProcessSS case %d:  ss1[%d]: %f should have gotten %f after spouses death\n", i, r2end, ss2[r2end], woulda)
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
