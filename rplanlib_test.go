package rplanlib

import (
	"bytes"
	"fmt"
	"io"
	"math"
	"os"
	"strings"
	"testing"
	"unicode"
)

//
// Testing for vector_var_index.go
//

func TestVectorVarIndex(t *testing.T) {
	tests := []struct {
		years   int
		taxbins int
		cgbins  int
		accmap  map[string]int
	}{
		{ // case 0
			years:   10,
			taxbins: 8,
			cgbins:  3,
			accmap:  map[string]int{"IRA": 1, "roth": 1, "aftertax": 1},
		},
		{ // case 1
			years:   100,
			taxbins: 8,
			cgbins:  3,
			accmap:  map[string]int{"IRA": 2, "roth": 2, "aftertax": 1},
		},
	}
	for i, elem := range tests {
		vvindex, err := NewVectorVarIndex(elem.years, elem.taxbins,
			elem.cgbins, elem.accmap)
		if err != nil {
			t.Errorf("VectorVarIndex case %d: %s", i, err)
			continue
		}
		OK := checkIndexSequence(elem.years, elem.taxbins,
			elem.cgbins, elem.accmap, vvindex)
		if OK != true {
			t.Errorf("VectorVarIndex case %d: Failed\n", i)
		}
	}
}

//
// Testing for taxinfo.go
//

func TestTaxinfo(t *testing.T) {
	tests := []struct {
		filingStatus string
		//spot check info
		brackets          int
		thirdBracketStart float64
	}{
		{ // case 0
			filingStatus:      "single",
			brackets:          7,
			thirdBracketStart: 37950,
		},
		{ // case 1
			filingStatus:      "joint",
			brackets:          7,
			thirdBracketStart: 75900,
		},
		{ // case 2
			filingStatus:      "mseparate",
			brackets:          7,
			thirdBracketStart: 37950,
		},
	}
	for i, elem := range tests {
		ti := NewTaxInfo(elem.filingStatus)
		brackets := len(*ti.Taxtable)
		if brackets != elem.brackets {
			t.Errorf("Taxinfo case %d: Failed - Expected %d brackes but found %d\n", i, elem.brackets, brackets)
		}
		if (*ti.Taxtable)[2][0] != elem.thirdBracketStart {
			t.Errorf("Taxinfo case %d: Failed - Expected %f for third bracket start but found %f\n", i, elem.thirdBracketStart, (*ti.Taxtable)[2][0])
		}
	}

}

func TestMaxContribution(t *testing.T) {

	tests := []struct {
		filingStatus string
		retireeindx  int
		year         int
		irate        float64
	}{
		{ // case 0
			filingStatus: "single",
			retireeindx:  0,
			year:         5,
			irate:        1.025,
		},
		{ // case 0
			filingStatus: "mseparate",
			retireeindx:  1,
			year:         5,
			irate:        1.025,
		},
		{ // case 0
			filingStatus: "joint",
			retireeindx:  1,
			year:         5,
			irate:        1.025,
		},
		{ // case 0
			filingStatus: "joint",
			retireeindx:  2, // pass in an empty key
			year:         5,
			irate:        1.025,
		},
	}

	retirees := []retiree{
		{ // retireeindx == 0
			age:        56,
			ageAtStart: 57,
			throughAge: 100,
			mykey:      "retiree1",
			definedContributionPlan: false,
			dcpBuckets:              nil,
		},
		{ // retireeindx == 1
			age:        54,
			ageAtStart: 55,
			throughAge: 100,
			mykey:      "retiree2",
			definedContributionPlan: false,
			dcpBuckets:              nil,
		},
		{ // retireeindx == 2 // fake retiree for getting empty mykey
			age:        0,
			ageAtStart: 0,
			throughAge: 0,
			mykey:      "", // empty mykey
			definedContributionPlan: false,
			dcpBuckets:              nil,
		},
	}
	for i, elem := range tests {
		ti := NewTaxInfo(elem.filingStatus)
		retiree := retirees[elem.retireeindx]
		prePlanYears := retiree.ageAtStart - retiree.age
		m := ti.maxContribution(elem.year, prePlanYears+elem.year, retirees, retiree.mykey, elem.irate)
		//fmt.Printf("m: %f, year: %d, prePlanYears: %d, key: %s, irate: %f\n", m, elem.year, prePlanYears, retiree.mykey, elem.irate)
		inflateYears := retiree.ageAtStart - retiree.age + elem.year
		memax := ti.Contribspecs["TDRA"] + ti.Contribspecs["TDRACatchup"]
		emax := memax * math.Pow(elem.irate, float64(inflateYears)) // adjust for inflation ??? current ????
		if retiree.mykey == "" {
			emax *= 2
		}
		//fmt.Printf("memax: %f, emax: %f\n", memax, emax)
		if emax != m {
			t.Errorf("maxContribution case %d: Failed - Expected %f but found %f\n", i, emax, m)
		}
	}
}

func TestApplyEarlyPenalty(t *testing.T) {
	retiree1 := retiree{
		age:        56,
		ageAtStart: 57,
		throughAge: 100,
		mykey:      "retiree1",
		definedContributionPlan: false,
		dcpBuckets:              nil,
	}

	tests := []struct {
		filingStatus string
		year         int
		response     bool
		retireer     *retiree
	}{
		{ // case 0
			filingStatus: "single",
			year:         2,
			response:     true,
			retireer:     &retiree1,
		},
		{ // case 1
			filingStatus: "single",
			year:         10,
			response:     false,
			retireer:     &retiree1,
		},
		{ // case 2
			filingStatus: "single",
			year:         3,
			response:     false,
			retireer:     &retiree1,
		},
		{ // case 3
			filingStatus: "single",
			year:         1,
			response:     false,
			retireer:     nil,
		},
	}
	for i, elem := range tests {
		ti := NewTaxInfo(elem.filingStatus)
		response := ti.applyEarlyPenalty(elem.year, elem.retireer)
		if response != elem.response {
			t.Errorf("applyEarlyPenalty case %d: Failed - Expected %v but found %v\n", i, elem.response, response)
		}
	}
}

/*
	tests := []struct {
	}{
		{},
	}
	for i, elem := range tests {
	}
*/

func TestRmdNeeded(t *testing.T) {
	retiree1 := retiree{
		age:        56,
		ageAtStart: 57,
		throughAge: 100,
		mykey:      "retiree1",
		definedContributionPlan: false,
		dcpBuckets:              nil,
	}

	tests := []struct {
		filingStatus string
		year         int
		response     float64
		retireer     *retiree
	}{
		{ // case 0
			filingStatus: "joint",
			year:         5,
			response:     0,
			retireer:     &retiree1,
		},
		{ // case 1
			filingStatus: "joint",
			year:         5,
			response:     0,
			retireer:     &retiree1,
		},
		{ // case 2
			filingStatus: "joint",
			year:         12,
			response:     0,
			retireer:     &retiree1,
		},
		{ // case 3
			filingStatus: "joint",
			year:         13, // should start here
			response:     27.4,
			retireer:     &retiree1,
		},
		{ // case 3
			filingStatus: "joint",
			year:         23, // should start here
			response:     18.7,
			retireer:     &retiree1,
		},
	}
	for i, elem := range tests {
		ti := NewTaxInfo(elem.filingStatus)
		response := ti.rmdNeeded(elem.year, elem.retireer)
		if response != elem.response {
			t.Errorf("rmdNeeded case %d: Failed - Expected %v but found %v\n", i, elem.response, response)
		}
	}

}

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
	}
	for i, elem := range tests {
		modelip := NewInputParams(elem.ip)
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

//
// Testing for lp_constraint_model.go
//

func TestIntMax(t *testing.T) {
	tests := []struct {
		a   int
		b   int
		max int
	}{
		{ // case 0
			a:   5,
			b:   6,
			max: 6,
		},
		{ // case 1
			a:   7,
			b:   6,
			max: 7,
		},
		{ // case 2
			a:   6,
			b:   6,
			max: 6,
		},
		{ // case 3
			a:   -10,
			b:   -6,
			max: -6,
		},
	}
	for i, elem := range tests {
		rmax := intMax(elem.a, elem.b)
		if rmax != elem.max {
			t.Errorf("intMax case %d: Failed - Expected %v but found %v\n", i, elem.max, rmax)
		}
	}
}

func TestIntMin(t *testing.T) {
	tests := []struct {
		a   int
		b   int
		min int
	}{
		{ // case 0
			a:   5,
			b:   6,
			min: 5,
		},
		{ // case 1
			a:   7,
			b:   6,
			min: 6,
		},
		{ // case 2
			a:   6,
			b:   6,
			min: 6,
		},
		{ // case 3
			a:   -10,
			b:   -6,
			min: -10,
		},
	}
	for i, elem := range tests {
		rmin := intMin(elem.a, elem.b)
		if rmin != elem.min {
			t.Errorf("intMin case %d: Failed - Expected %v but found %v\n", i, elem.min, rmin)
		}
	}
}

func TestCheckStrconvError(t *testing.T) {
	tests := []struct {
		//err    error
		errstr string
	}{
		{ // case 0
			errstr: "case 0",
		},
		{ // case 1
			errstr: "case 1",
		},
		{ // case 2
			errstr: "",
		},
	}
	for i, elem := range tests {
		var err error
		err = nil
		if elem.errstr != "" {
			err = fmt.Errorf(elem.errstr)
		}
		func() {
			defer func() {
				r := recover()
				if r == nil && elem.errstr != "" {
					t.Errorf("checkStrconvError case %d.a should have panicked", i)
				} else if elem.errstr == "" && r != nil {
					t.Errorf("checkStrconvError case %d.b should have panicked", i)
				} else if r != nil {
					errstr := fmt.Sprintf("%s", r)
					if errstr != elem.errstr {
						t.Errorf("checkStrconvError case %d panicked with err '%s' but should have err '%s'", i, errstr, elem.errstr)
					}
				}
			}()
			// This function should cause a panic
			checkStrconvError(err)
		}()
	}
}

func TestMergeVectors(t *testing.T) {
	tests := []struct {
		a      []float64
		b      []float64
		errstr string
	}{
		{ // Case 0
			a:      []float64{5, 2, -2, 388886, 0},
			b:      []float64{20, 30, 40, 50, 60},
			errstr: "",
		},
		{ // Case 1
			a:      []float64{5, 2, -2, 388886, 0, 20},
			b:      []float64{20, 30, 40, 50, 60},
			errstr: "mergeVectors: Can not merge, lengths do not match, 6 vs 5",
		},
	}
	for i, elem := range tests {
		newv, err := mergeVectors(elem.a, elem.b)
		if err != nil {
			if len(elem.a) == len(elem.b) {
				t.Errorf("mergeVectors case %d failed but should not have!", i)
			}
			s := fmt.Sprintf("%v", err)
			if s != elem.errstr {
				t.Errorf("mergeVectors case %d failed with incorrect err string\n\tExpected: '%s' but found: '%s'", i, elem.errstr, s)
			}
			continue
		}
		for i := 0; i < len(newv); i++ {
			if newv[i] != elem.a[i]+elem.b[i] {
				t.Errorf("mergeVectors case %d merged values do no sum", i)
			}
		}
	}
}

func TestBuildVector(t *testing.T) {
	tests := []struct {
		yearly      int
		startAge    int
		endAge      int
		vecStartAge int
		vecEndAge   int
		rate        float64
		baseAge     int
		errstr      string
	}{
		{ // case 0 // over begining of vec
			yearly:      1,
			startAge:    45,
			endAge:      66,
			vecStartAge: 62,
			vecEndAge:   100,
			rate:        1.025,
			baseAge:     40,
			errstr:      "",
		},
		{ // case 1 // over ending of vec
			yearly:      1,
			startAge:    70,
			endAge:      102,
			vecStartAge: 62,
			vecEndAge:   100,
			rate:        1.025,
			baseAge:     40,
			errstr:      "",
		},
		{ // case 2 // in the middle of vec
			yearly:      1,
			startAge:    66,
			endAge:      68,
			vecStartAge: 62,
			vecEndAge:   100,
			rate:        1.025,
			baseAge:     40,
			errstr:      "",
		},
		{ // case 3 // all above vec
			yearly:      1,
			startAge:    145,
			endAge:      166,
			vecStartAge: 62,
			vecEndAge:   100,
			rate:        1.025,
			baseAge:     40,
			errstr:      "",
		},
		{ // case 4 // all below vec
			yearly:      1,
			startAge:    45,
			endAge:      60,
			vecStartAge: 62,
			vecEndAge:   100,
			rate:        1.025,
			baseAge:     40,
			errstr:      "",
		},
		{ // case 5 // all match vec
			yearly:      1,
			startAge:    62,
			endAge:      100,
			vecStartAge: 62,
			vecEndAge:   100,
			rate:        1.025,
			baseAge:     40,
			errstr:      "",
		},
		{ // case 6 // vec start > vec end
			yearly:      1,
			startAge:    62,
			endAge:      100,
			vecStartAge: 100,
			vecEndAge:   62,
			rate:        1.025,
			baseAge:     40,
			errstr:      "vec start age (100) is greater than vec end age (62)",
		},
		{ // case 7 // start age > end age
			yearly:      1,
			startAge:    100,
			endAge:      62,
			vecStartAge: 62,
			vecEndAge:   100,
			rate:        1.025,
			baseAge:     40,
			errstr:      "start age (100) is greater than end age (62)",
		},
	}
	for i, elem := range tests {
		newv, err := buildVector(elem.yearly, elem.startAge, elem.endAge, elem.vecStartAge, elem.vecEndAge, elem.rate, elem.baseAge)
		if err != nil {
			es := fmt.Sprintf("%s", err)
			if elem.errstr != es {
				t.Errorf("buildVector case %d: expected errstr '%s', found '%s'", i, elem.errstr, es)
			}
			// tbd TODO fix this
			//fmt.Printf("&&&&&&&&&& buildVector() returned and err for case %d: %s\n", i, err)
			continue
		}
		fnz := -1
		if elem.startAge < elem.vecEndAge && elem.startAge >= elem.vecStartAge {
			fnz = elem.startAge - elem.vecStartAge
		} else if elem.startAge < elem.vecStartAge && elem.endAge > elem.vecStartAge /*elem.vecEndAge*/ {
			fnz = 0
		}
		lnz := len(newv) + 1
		if elem.endAge < elem.vecEndAge && elem.endAge >= elem.vecStartAge {
			lnz = elem.endAge - elem.vecStartAge
		} else if elem.endAge >= elem.vecEndAge && elem.startAge <= elem.vecEndAge {
			lnz = len(newv) - 1
		}
		//fmt.Printf("CASE %d: ===================================\n", i)
		//fmt.Printf("### endAge(%d) < vecEndAge(%d) && endAge(%d) >= vecStartAge(%d)\n", elem.endAge, elem.vecEndAge, elem.endAge, elem.vecStartAge)
		//fmt.Printf("*** endAge(%d) >= vecEndAge(%d) && startAge(%d) <= vecEndAge(%d)\n", elem.endAge, elem.vecEndAge, elem.startAge, elem.vecEndAge)
		firstNonZero := -1
		lastNonZero := len(newv) + 1
		for f := 0; f < len(newv); f++ {
			if newv[f] != 0 && firstNonZero < 0 {
				firstNonZero = f
			}
			if newv[f] != 0 && firstNonZero >= 0 {
				lastNonZero = f
			}
		}
		if fnz != firstNonZero {
			t.Errorf("buildVector case %d: firstNonZero is incorrect, expected %d, found %d", i, fnz, firstNonZero)
		}
		if lnz != lastNonZero {
			t.Errorf("buildVector case %d: lastNonZero is incorrect, expected %d, found %d", i, lnz, lastNonZero)
		}
		//fmt.Printf("Case %d: newv len:%d ============\n", i, len(newv))
		//fmt.Printf("firstNonZero: %d, lastNonZero: %d\n", firstNonZero, lastNonZero)
		//fmt.Printf("fnz: %d, lnz: %d\n", fnz, lnz)
		//fmt.Printf("Case %d: %v\n", i, newv)
	}
}

func TestNewModelSpecs(t *testing.T) {
	tests := []struct {
		years         int
		accmap        map[string]int
		ip            map[string]string
		verbose       bool
		allowDeposits bool
		iRate         float64
	}{
		{ // Case 0 // joint
			years:  10,
			accmap: map[string]int{"IRA": 2, "roth": 1, "aftertax": 1},
			ip: map[string]string{
				"setName":                    "activeParams",
				"filingStatus":               "joint",
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
				"eT_Roth2":                   "50", // 50k
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
			verbose:       false,
			allowDeposits: false,
			iRate:         1.025,
		},
		{ // Case 1 // mseparate
			years:  10,
			accmap: map[string]int{"IRA": 1, "roth": 1, "aftertax": 1},
			ip: map[string]string{
				"setName":                    "activeParams",
				"filingStatus":               "mseparate",
				"eT_Age1":                    "",
				"eT_Age2":                    "",
				"eT_RetireAge1":              "",
				"eT_RetireAge2":              "",
				"eT_PlanThroughAge1":         "",
				"eT_PlanThroughAge2":         "",
				"eT_PIA1":                    "",
				"eT_PIA2":                    "",
				"eT_SS_Start1":               "",
				"eT_SS_Start2":               "",
				"eT_TDRA1":                   "",
				"eT_TDRA2":                   "100", //100k
				"eT_TDRA_Rate1":              "",
				"eT_TDRA_Rate2":              "",
				"eT_TDRA_Contrib1":           "",
				"eT_TDRA_Contrib2":           "",
				"eT_TDRA_ContribStartAge1":   "",
				"eT_TDRA_ContribStartAge2":   "",
				"eT_TDRA_ContribEndAge1":     "",
				"eT_TDRA_ContribEndAge2":     "",
				"eT_Roth1":                   "100", //100k
				"eT_Roth2":                   "",
				"eT_Roth_Rate1":              "",
				"eT_Roth_Rate2":              "",
				"eT_Roth_Contrib1":           "",
				"eT_Roth_Contrib2":           "",
				"eT_Roth_ContribStartAge1":   "",
				"eT_Roth_ContribStartAge2":   "",
				"eT_Roth_ContribEndAge1":     "",
				"eT_Roth_ContribEndAge2":     "",
				"eT_Aftatax":                 "30", //30k
				"eT_Aftatax_Rate":            "",
				"eT_Aftatax_Contrib":         "",
				"eT_Aftatax_ContribStartAge": "",
				"eT_Aftatax_ContribEndAge":   "",
			},
			verbose:       false,
			allowDeposits: false,
			iRate:         1.025,
		},
		{ // Case 2 // single
			years:  10,
			accmap: map[string]int{"IRA": 1, "roth": 1, "aftertax": 1},
			//ip:            map[string]string{"filingStatus": "single"},
			ip: map[string]string{
				"setName":                    "activeParams",
				"filingStatus":               "single",
				"eT_Age1":                    "",
				"eT_Age2":                    "",
				"eT_RetireAge1":              "",
				"eT_RetireAge2":              "",
				"eT_PlanThroughAge1":         "",
				"eT_PlanThroughAge2":         "",
				"eT_PIA1":                    "",
				"eT_PIA2":                    "",
				"eT_SS_Start1":               "",
				"eT_SS_Start2":               "",
				"eT_TDRA1":                   "40", // 40k
				"eT_TDRA2":                   "",
				"eT_TDRA_Rate1":              "",
				"eT_TDRA_Rate2":              "",
				"eT_TDRA_Contrib1":           "",
				"eT_TDRA_Contrib2":           "",
				"eT_TDRA_ContribStartAge1":   "",
				"eT_TDRA_ContribStartAge2":   "",
				"eT_TDRA_ContribEndAge1":     "",
				"eT_TDRA_ContribEndAge2":     "",
				"eT_Roth1":                   "40", // 40k
				"eT_Roth2":                   "",
				"eT_Roth_Rate1":              "",
				"eT_Roth_Rate2":              "",
				"eT_Roth_Contrib1":           "",
				"eT_Roth_Contrib2":           "",
				"eT_Roth_ContribStartAge1":   "",
				"eT_Roth_ContribStartAge2":   "",
				"eT_Roth_ContribEndAge1":     "",
				"eT_Roth_ContribEndAge2":     "",
				"eT_Aftatax":                 "20", // 20k
				"eT_Aftatax_Rate":            "",
				"eT_Aftatax_Contrib":         "",
				"eT_Aftatax_ContribStartAge": "",
				"eT_Aftatax_ContribEndAge":   "",
			},
			verbose:       false,
			allowDeposits: false,
			iRate:         1.025,
		},
	}
	for i, elem := range tests {
		ip := NewInputParams(elem.ip)
		ti := NewTaxInfo(ip.filingStatus)
		taxbins := len(*ti.Taxtable)
		cgbins := len(*ti.Capgainstable)
		vindx, err := NewVectorVarIndex(ip.numyr, taxbins, cgbins, elem.accmap)
		if err != nil {
			t.Errorf("NewModelSpecs case %d: %s", i, err)
			continue
		}
		ms := NewModelSpecs(vindx, ti, ip, elem.verbose,
			elem.allowDeposits)
		if ms.iRate != elem.iRate {
			t.Errorf("NewModelSpecs case %d: iRate expected %f, found %f", i, elem.iRate, ms.iRate)
		}
	}
}

func TestBuildModel(t *testing.T) {
	tests := []struct {
		years         int
		accmap        map[string]int
		ip            map[string]string
		verbose       bool
		allowDeposits bool
		iRate         float64
	}{
		{ // Case 0 // joint
			years:  10,
			accmap: map[string]int{"IRA": 1, "roth": 1, "aftertax": 1},
			ip: map[string]string{
				"setName":                    "activeParams",
				"filingStatus":               "joint",
				"eT_Age1":                    "",
				"eT_Age2":                    "",
				"eT_RetireAge1":              "",
				"eT_RetireAge2":              "",
				"eT_PlanThroughAge1":         "",
				"eT_PlanThroughAge2":         "",
				"eT_PIA1":                    "",
				"eT_PIA2":                    "",
				"eT_SS_Start1":               "",
				"eT_SS_Start2":               "",
				"eT_TDRA1":                   "40", // 40k
				"eT_TDRA2":                   "",
				"eT_TDRA_Rate1":              "",
				"eT_TDRA_Rate2":              "",
				"eT_TDRA_Contrib1":           "",
				"eT_TDRA_Contrib2":           "",
				"eT_TDRA_ContribStartAge1":   "",
				"eT_TDRA_ContribStartAge2":   "",
				"eT_TDRA_ContribEndAge1":     "",
				"eT_TDRA_ContribEndAge2":     "",
				"eT_Roth1":                   "40", // 40k
				"eT_Roth2":                   "",
				"eT_Roth_Rate1":              "",
				"eT_Roth_Rate2":              "",
				"eT_Roth_Contrib1":           "",
				"eT_Roth_Contrib2":           "",
				"eT_Roth_ContribStartAge1":   "",
				"eT_Roth_ContribStartAge2":   "",
				"eT_Roth_ContribEndAge1":     "",
				"eT_Roth_ContribEndAge2":     "",
				"eT_Aftatax":                 "20", // 20k
				"eT_Aftatax_Rate":            "",
				"eT_Aftatax_Contrib":         "",
				"eT_Aftatax_ContribStartAge": "",
				"eT_Aftatax_ContribEndAge":   "",
			},
			verbose:       false,
			allowDeposits: false,
			iRate:         1.025,
		},
	}
	for i, elem := range tests {
		ti := NewTaxInfo(elem.ip["filingStatus"])
		ip := NewInputParams(elem.ip)
		taxbins := len(*ti.Taxtable)
		cgbins := len(*ti.Capgainstable)
		vindx, err := NewVectorVarIndex(elem.years, taxbins,
			cgbins, ip.accmap)
		if err != nil {
			t.Errorf("BuildModel case %d: %s", i, err)
			continue
		}
		ms := NewModelSpecs(vindx, ti, ip, elem.verbose,
			elem.allowDeposits)
		/**/
		c, A, b := ms.BuildModel()
		ms.printModelMatrix(c, A, b, nil, false)
		/**/
		if ms.iRate != elem.iRate {
			t.Errorf("BuildModel case %d: iRate expected %f, found %f", i, elem.iRate, ms.iRate)
		}
	}
}

func TestAccountOwnerAge(t *testing.T) {
	tests := []struct {
		ms    ModelSpecs
		index int
		year  int
	}{
		{ // case 0
			ms: ModelSpecs{
				retirees: []retiree{
					{ // retireeindx == 0
						age:        56,
						ageAtStart: 57,
						throughAge: 100,
						mykey:      "retiree1",
						definedContributionPlan: false,
						dcpBuckets:              nil,
					},
					{ // retireeindx == 1
						age:        54,
						ageAtStart: 55,
						throughAge: 100,
						mykey:      "retiree2",
						definedContributionPlan: false,
						dcpBuckets:              nil,
					},
				},
				accounttable: []account{
					{
						bal:           30,
						basis:         0,
						estateTax:     0.85,
						contributions: []float64{},
						rRate:         1.06,
						acctype:       "IRA",
						mykey:         "retiree2",
					},
				},
			},
			index: 1,
			year:  10,
		},
		{ // case 1
			ms: ModelSpecs{
				retirees: []retiree{
					{ // retireeindx == 0
						age:        56,
						ageAtStart: 57,
						throughAge: 100,
						mykey:      "retiree1",
						definedContributionPlan: false,
						dcpBuckets:              nil,
					},
					{ // retireeindx == 1
						age:        54,
						ageAtStart: 55,
						throughAge: 100,
						mykey:      "retiree2",
						definedContributionPlan: false,
						dcpBuckets:              nil,
					},
				},
				accounttable: []account{
					{
						bal:           30,
						basis:         0,
						estateTax:     0.85,
						contributions: []float64{},
						rRate:         1.06,
						acctype:       "IRA",
						mykey:         "retiree2",
					},
				},
			},
			index: 1,
			year:  7,
		},
	}
	for i, elem := range tests {
		ownerAge := elem.ms.accountOwnerAge(elem.year, elem.ms.accounttable[0])
		calcage := elem.ms.retirees[elem.index].ageAtStart + elem.year
		if ownerAge != calcage {
			t.Errorf("AccountOwnerAge case %d: age does not match, expected %d, found %d", i, calcage, ownerAge)
		}
	}
}

func TestMatchRetiree(t *testing.T) {
	tests := []struct {
		ms  ModelSpecs
		key string
		age int
	}{
		{ // case 0
			ms: ModelSpecs{
				retirees: []retiree{
					{ // retireeindx == 0
						age:        56,
						ageAtStart: 57,
						throughAge: 100,
						mykey:      "retiree1",
						definedContributionPlan: false,
						dcpBuckets:              nil,
					},
					{ // retireeindx == 1
						age:        54,
						ageAtStart: 55,
						throughAge: 100,
						mykey:      "retiree2",
						definedContributionPlan: false,
						dcpBuckets:              nil,
					},
				},
			},
			key: "retiree2",
			age: 54,
		},
		{ // case 1
			ms: ModelSpecs{
				retirees: []retiree{
					{ // retireeindx == 0
						age:        56,
						ageAtStart: 57,
						throughAge: 100,
						mykey:      "retiree1",
						definedContributionPlan: false,
						dcpBuckets:              nil,
					},
					{ // retireeindx == 1
						age:        54,
						ageAtStart: 55,
						throughAge: 100,
						mykey:      "retiree2",
						definedContributionPlan: false,
						dcpBuckets:              nil,
					},
				},
			},
			key: "retiree1",
			age: 56,
		},
	}
	for i, elem := range tests {
		r := elem.ms.matchRetiree(elem.key)
		if r.mykey != elem.key {
			t.Errorf("MatchRetiree case %d: key does not match, expected %s, found %s", i, elem.key, r.mykey)
		}
		if r.age != elem.age {
			t.Errorf("MatchRetiree case %d: age does not match, expected %d, found %d", i, elem.age, r.age)
		}
	}
}

func TestCgTaxableFraction(t *testing.T) { /* TODO:FIXME:IMPLEMENTME */
	tests := []struct {
		ms      ModelSpecs
		expectf float64
		year    int
	}{
		{ // case 0
			ms: ModelSpecs{
				retirees: []retiree{
					{ // retireeindx == 0
						age:        56,
						ageAtStart: 57,
						throughAge: 100,
						mykey:      "retiree1",
						definedContributionPlan: false,
						dcpBuckets:              nil,
					},
					{ // retireeindx == 1
						age:        54,
						ageAtStart: 55,
						throughAge: 100,
						mykey:      "retiree2",
						definedContributionPlan: false,
						dcpBuckets:              nil,
					},
				},
				accounttable: []account{
					{
						bal:           30,
						basis:         20,
						estateTax:     0.85,
						contributions: []float64{},
						rRate:         1.06,
						acctype:       "IRA",
						mykey:         "retiree2",
					},
				},
				accmap: map[string]int{
					"IRA":      1,
					"roth":     0,
					"aftertax": 0,
				},
			},
			expectf: 1, //no aftertax account
			year:    10,
		},
		{ // case 1
			ms: ModelSpecs{
				retirees: []retiree{
					{ // retireeindx == 0
						age:        56,
						ageAtStart: 57,
						throughAge: 100,
						mykey:      "retiree1",
						definedContributionPlan: false,
						dcpBuckets:              nil,
					},
					{ // retireeindx == 1
						age:        54,
						ageAtStart: 55,
						throughAge: 100,
						mykey:      "retiree2",
						definedContributionPlan: false,
						dcpBuckets:              nil,
					},
				},
				accounttable: []account{
					{
						bal:           30,
						basis:         10,
						estateTax:     0.85,
						contributions: []float64{},
						rRate:         1.06,
						acctype:       "IRA",
						mykey:         "retiree2",
					},
					{
						bal:           30,
						basis:         10,
						estateTax:     0.85,
						contributions: []float64{},
						rRate:         1.06,
						acctype:       "aftertax",
						mykey:         "retiree2",
					},
				},
				accmap: map[string]int{
					"IRA":      1,
					"roth":     0,
					"aftertax": 1,
				},
			},
			expectf: -1, //no aftertax account
			year:    7,
		},
	}
	for i, elem := range tests {
		f := elem.ms.cgTaxableFraction(elem.year)
		fprime := elem.expectf
		if elem.expectf < 0 {
			fprime = 1 - (elem.ms.accounttable[0].basis / (elem.ms.accounttable[0].bal * math.Pow(elem.ms.accounttable[0].rRate, float64(elem.year+elem.ms.prePlanYears))))
		}
		if f != fprime {
			t.Errorf("cgTaxableFraction case %d: expected %f, found %f", i, fprime, f)
		}
	}
	/*
		func (ms ModelSpecs) cgTaxableFraction(year int) float64 {
	*/
}

func TestPrintModelMatrix(t *testing.T) { /* TODO:FIXME:IMPLEMENTME */
	/*
		tests := []struct {
		}{
			{},
		}
		for i, elem := range tests {
		}
	*/
}

func TestPrintConstraint(t *testing.T) { /* TODO:FIXME:IMPLEMENTME */
	tests := []struct {
		ip        map[string]string
		row       []float64
		b         float64
		expectstr string
		testcase  string
	}{
		{ // Case 0
			ip: map[string]string{
				"setName":                    "activeParams",
				"filingStatus":               "single",
				"eT_Age1":                    "60",
				"eT_Age2":                    "",
				"eT_RetireAge1":              "65",
				"eT_RetireAge2":              "",
				"eT_PlanThroughAge1":         "70",
				"eT_PlanThroughAge2":         "",
				"eT_PIA1":                    "20", // 20k
				"eT_PIA2":                    "",
				"eT_SS_Start1":               "70",
				"eT_SS_Start2":               "",
				"eT_TDRA1":                   "10", // 10k
				"eT_TDRA2":                   "",
				"eT_TDRA_Rate1":              "",
				"eT_TDRA_Rate2":              "",
				"eT_TDRA_Contrib1":           "",
				"eT_TDRA_Contrib2":           "",
				"eT_TDRA_ContribStartAge1":   "",
				"eT_TDRA_ContribStartAge2":   "",
				"eT_TDRA_ContribEndAge1":     "",
				"eT_TDRA_ContribEndAge2":     "",
				"eT_Roth1":                   "5", // 5k
				"eT_Roth2":                   "",
				"eT_Roth_Rate1":              "",
				"eT_Roth_Rate2":              "",
				"eT_Roth_Contrib1":           "",
				"eT_Roth_Contrib2":           "",
				"eT_Roth_ContribStartAge1":   "",
				"eT_Roth_ContribStartAge2":   "",
				"eT_Roth_ContribEndAge1":     "",
				"eT_Roth_ContribEndAge2":     "",
				"eT_Aftatax":                 "15", // 15k
				"eT_Aftatax_Rate":            "",
				"eT_Aftatax_Contrib":         "",
				"eT_Aftatax_ContribStartAge": "",
				"eT_Aftatax_ContribEndAge":   "",
			},
			row:       []float64{},
			b:         743.027,
			expectstr: "Row: [0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0] b: 743.027 <= b[]: 743.03 ",
			testcase:  "allzeros",
		},
		{ // Case 1
			ip: map[string]string{
				"setName":                    "activeParams",
				"filingStatus":               "single",
				"eT_Age1":                    "60",
				"eT_Age2":                    "",
				"eT_RetireAge1":              "65",
				"eT_RetireAge2":              "",
				"eT_PlanThroughAge1":         "70",
				"eT_PlanThroughAge2":         "",
				"eT_PIA1":                    "20", // 20k
				"eT_PIA2":                    "",
				"eT_SS_Start1":               "70",
				"eT_SS_Start2":               "",
				"eT_TDRA1":                   "10", // 10k
				"eT_TDRA2":                   "",
				"eT_TDRA_Rate1":              "",
				"eT_TDRA_Rate2":              "",
				"eT_TDRA_Contrib1":           "",
				"eT_TDRA_Contrib2":           "",
				"eT_TDRA_ContribStartAge1":   "",
				"eT_TDRA_ContribStartAge2":   "",
				"eT_TDRA_ContribEndAge1":     "",
				"eT_TDRA_ContribEndAge2":     "",
				"eT_Roth1":                   "5", // 5k
				"eT_Roth2":                   "",
				"eT_Roth_Rate1":              "",
				"eT_Roth_Rate2":              "",
				"eT_Roth_Contrib1":           "",
				"eT_Roth_Contrib2":           "",
				"eT_Roth_ContribStartAge1":   "",
				"eT_Roth_ContribStartAge2":   "",
				"eT_Roth_ContribEndAge1":     "",
				"eT_Roth_ContribEndAge2":     "",
				"eT_Aftatax":                 "15", // 15k
				"eT_Aftatax_Rate":            "",
				"eT_Aftatax_Contrib":         "",
				"eT_Aftatax_ContribStartAge": "",
				"eT_Aftatax_ContribEndAge":   "",
			},
			row:       []float64{},
			b:         743.027,
			expectstr: "Row: [0 1 2 3 4 5 6 7 8 9 10 11 12 13 14 15 16 17 18 19 20 21 22 23 24 25 26 27 28 29 30 31 32 33 34 35 36 37 38 39 40 41 42 43 44 45 46 47 48 49 50 51 52 53 54 55 56 57 58 59 60 61 62 63 64 65 66 67 68 69 70 71 72 73 74 75 76 77 78 79 80 81 82 83 84 85 86 87 88 89 90 91 92 93 94 95 96 97 98 99 100 101 102 103 104 105 106 107 108 109 110 111 112 113 114 115 116 117 118 119 120 121 122 123 124 125 126 127 128 129 130 131 132 133 134 135 136 137 138 139 140 141 142 143 144 145 146 147 148 149 150 151 152 153 154 155 156 157 158 159 160 161 162 163 164 165 166 167 168 169 170 171 172 173 174 175 176 177 178 179 180 181 182 183 184 185 186 187 188 189 190 191 192 193 194 195 196 197 198 199 200 201 202 203 204 205 206 207 208 209 210 211 212 213 214 215 216 217 218 219 220 221 222] b: 743.027 x[0,1]= 1.000, x[0,2]= 2.000, x[0,3]= 3.000, x[0,4]= 4.000, x[0,5]= 5.000, x[0,6]= 6.000, x[1,0]= 7.000, x[1,1]= 8.000, x[1,2]= 9.000, x[1,3]=10.000, x[1,4]=11.000, x[1,5]=12.000, x[1,6]=13.000, x[2,0]=14.000, x[2,1]=15.000, x[2,2]=16.000, x[2,3]=17.000, x[2,4]=18.000, x[2,5]=19.000, x[2,6]=20.000, x[3,0]=21.000, x[3,1]=22.000, x[3,2]=23.000, x[3,3]=24.000, x[3,4]=25.000, x[3,5]=26.000, x[3,6]=27.000, x[4,0]=28.000, x[4,1]=29.000, x[4,2]=30.000, x[4,3]=31.000, x[4,4]=32.000, x[4,5]=33.000, x[4,6]=34.000, x[5,0]=35.000, x[5,1]=36.000, x[5,2]=37.000, x[5,3]=38.000, x[5,4]=39.000, x[5,5]=40.000, x[5,6]=41.000, x[6,0]=42.000, x[6,1]=43.000, x[6,2]=44.000, x[6,3]=45.000, x[6,4]=46.000, x[6,5]=47.000, x[6,6]=48.000, x[7,0]=49.000, x[7,1]=50.000, x[7,2]=51.000, x[7,3]=52.000, x[7,4]=53.000, x[7,5]=54.000, x[7,6]=55.000, x[8,0]=56.000, x[8,1]=57.000, x[8,2]=58.000, x[8,3]=59.000, x[8,4]=60.000, x[8,5]=61.000, x[8,6]=62.000, x[9,0]=63.000, x[9,1]=64.000, x[9,2]=65.000, x[9,3]=66.000, x[9,4]=67.000, x[9,5]=68.000, x[9,6]=69.000, x[10,0]=70.000, x[10,1]=71.000, x[10,2]=72.000, x[10,3]=73.000, x[10,4]=74.000, x[10,5]=75.000, x[10,6]=76.000, y[0,0]=77.000, y[0,1]=78.000, y[0,2]=79.000, y[1,0]=80.000, y[1,1]=81.000, y[1,2]=82.000, y[2,0]=83.000, y[2,1]=84.000, y[2,2]=85.000, y[3,0]=86.000, y[3,1]=87.000, y[3,2]=88.000, y[4,0]=89.000, y[4,1]=90.000, y[4,2]=91.000, y[5,0]=92.000, y[5,1]=93.000, y[5,2]=94.000, y[6,0]=95.000, y[6,1]=96.000, y[6,2]=97.000, y[7,0]=98.000, y[7,1]=99.000, y[7,2]=100.000, y[8,0]=101.000, y[8,1]=102.000, y[8,2]=103.000, y[9,0]=104.000, y[9,1]=105.000, y[9,2]=106.000, y[10,0]=107.000, y[10,1]=108.000, y[10,2]=109.000, w[0,0]=110.000, w[0,1]=111.000, w[0,2]=112.000, w[1,0]=113.000, w[1,1]=114.000, w[1,2]=115.000, w[2,0]=116.000, w[2,1]=117.000, w[2,2]=118.000, w[3,0]=119.000, w[3,1]=120.000, w[3,2]=121.000, w[4,0]=122.000, w[4,1]=123.000, w[4,2]=124.000, w[5,0]=125.000, w[5,1]=126.000, w[5,2]=127.000, w[6,0]=128.000, w[6,1]=129.000, w[6,2]=130.000, w[7,0]=131.000, w[7,1]=132.000, w[7,2]=133.000, w[8,0]=134.000, w[8,1]=135.000, w[8,2]=136.000, w[9,0]=137.000, w[9,1]=138.000, w[9,2]=139.000, w[10,0]=140.000, w[10,1]=141.000, w[10,2]=142.000, b[0,0]=143.000, b[0,1]=144.000, b[0,2]=145.000, b[1,0]=146.000, b[1,1]=147.000, b[1,2]=148.000, b[2,0]=149.000, b[2,1]=150.000, b[2,2]=151.000, b[3,0]=152.000, b[3,1]=153.000, b[3,2]=154.000, b[4,0]=155.000, b[4,1]=156.000, b[4,2]=157.000, b[5,0]=158.000, b[5,1]=159.000, b[5,2]=160.000, b[6,0]=161.000, b[6,1]=162.000, b[6,2]=163.000, b[7,0]=164.000, b[7,1]=165.000, b[7,2]=166.000, b[8,0]=167.000, b[8,1]=168.000, b[8,2]=169.000, b[9,0]=170.000, b[9,1]=171.000, b[9,2]=172.000, b[10,0]=173.000, b[10,1]=174.000, b[10,2]=175.000, b[11,0]=176.000, b[11,1]=177.000, b[11,2]=178.000, s[0]=179.000, s[1]=180.000, s[2]=181.000, s[3]=182.000, s[4]=183.000, s[5]=184.000, s[6]=185.000, s[7]=186.000, s[8]=187.000, s[9]=188.000, s[10]=189.000, D[0,0]=190.000, D[0,1]=191.000, D[0,2]=192.000, D[1,0]=193.000, D[1,1]=194.000, D[1,2]=195.000, D[2,0]=196.000, D[2,1]=197.000, D[2,2]=198.000, D[3,0]=199.000, D[3,1]=200.000, D[3,2]=201.000, D[4,0]=202.000, D[4,1]=203.000, D[4,2]=204.000, D[5,0]=205.000, D[5,1]=206.000, D[5,2]=207.000, D[6,0]=208.000, D[6,1]=209.000, D[6,2]=210.000, D[7,0]=211.000, D[7,1]=212.000, D[7,2]=213.000, D[8,0]=214.000, D[8,1]=215.000, D[8,2]=216.000, D[9,0]=217.000, D[9,1]=218.000, D[9,2]=219.000, D[10,0]=220.000, D[10,1]=221.000, D[10,2]=222.000, <= b[]: 743.03 ",
			testcase:  "counting",
		},
	}
	for i, elem := range tests {
		ip := NewInputParams(elem.ip)
		ti := NewTaxInfo(ip.filingStatus)
		taxbins := len(*ti.Taxtable)
		cgbins := len(*ti.Capgainstable)
		//fmt.Printf("ACCMAP: %v", ip.accmap)
		vindx, err := NewVectorVarIndex(ip.numyr, taxbins, cgbins, ip.accmap)
		if err != nil {
			t.Errorf("PrintConstraint case %d: %s", i, err)
			continue
		}
		numaccounts := 0
		for _, acc := range ip.accmap {
			numaccounts += acc
		}
		ms := ModelSpecs{
			ip:     ip,
			vindx:  vindx,
			ti:     ti,
			numyr:  ip.numyr,
			accmap: ip.accmap,
			numacc: numaccounts,
		}

		row := make([]float64, vindx.Vsize)
		switch elem.testcase {
		case "allones":
			for indx := 0; indx < vindx.Vsize; indx++ {
				row[indx] = float64(indx)
			}
		case "counting":
			for indx := 0; indx < vindx.Vsize; indx++ {
				row[indx] = float64(indx)
			}
		case "allzeros":
			// nothing to change
		default:
			t.Errorf("TestPrintModelRow: Unexpected test case '%s'\n", elem.testcase)
			continue
		}
		//fmt.Printf("Vsize: %d\n", vindx.Vsize)

		mychan := make(chan string)
		oldout, w, err := RedirectStdout(mychan)
		if err != nil {
			t.Errorf("RedirectStdout: %s\n", err)
			return // should this be continue?
		}
		fmt.Printf("Row: %v\n", row)
		fmt.Printf("b: %v\n", elem.b)
		//ms.printModelRow(row, elem.suppressNewline)
		ms.printConstraint(row, elem.b)

		str := RestoreStdout(mychan, oldout, w)
		strn := stripWhitespace(str)
		strexpect := stripWhitespace(elem.expectstr)
		if strn != strexpect {
			t.Errorf("PrintModelRow Case %d: expected '%s', found '%s'", i, elem.expectstr, str)
		}
	}
}

func TestPrintModelRow(t *testing.T) {
	tests := []struct {
		ip              map[string]string
		row             []float64
		suppressNewline bool
		expectstr       string
		testcase        string
	}{
		{ // Case 0
			ip: map[string]string{
				"setName":                    "activeParams",
				"filingStatus":               "single",
				"eT_Age1":                    "60",
				"eT_Age2":                    "",
				"eT_RetireAge1":              "65",
				"eT_RetireAge2":              "",
				"eT_PlanThroughAge1":         "70",
				"eT_PlanThroughAge2":         "",
				"eT_PIA1":                    "20", // 20k
				"eT_PIA2":                    "",
				"eT_SS_Start1":               "70",
				"eT_SS_Start2":               "",
				"eT_TDRA1":                   "10", // 10k
				"eT_TDRA2":                   "",
				"eT_TDRA_Rate1":              "",
				"eT_TDRA_Rate2":              "",
				"eT_TDRA_Contrib1":           "",
				"eT_TDRA_Contrib2":           "",
				"eT_TDRA_ContribStartAge1":   "",
				"eT_TDRA_ContribStartAge2":   "",
				"eT_TDRA_ContribEndAge1":     "",
				"eT_TDRA_ContribEndAge2":     "",
				"eT_Roth1":                   "5", // 5k
				"eT_Roth2":                   "",
				"eT_Roth_Rate1":              "",
				"eT_Roth_Rate2":              "",
				"eT_Roth_Contrib1":           "",
				"eT_Roth_Contrib2":           "",
				"eT_Roth_ContribStartAge1":   "",
				"eT_Roth_ContribStartAge2":   "",
				"eT_Roth_ContribEndAge1":     "",
				"eT_Roth_ContribEndAge2":     "",
				"eT_Aftatax":                 "15", // 15k
				"eT_Aftatax_Rate":            "",
				"eT_Aftatax_Contrib":         "",
				"eT_Aftatax_ContribStartAge": "",
				"eT_Aftatax_ContribEndAge":   "",
			},
			row:             []float64{},
			suppressNewline: false,
			expectstr:       "Row: [0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0] ",
			testcase:        "allzeros",
		},
		{ // Case 1
			ip: map[string]string{
				"setName":                    "activeParams",
				"filingStatus":               "single",
				"eT_Age1":                    "60",
				"eT_Age2":                    "",
				"eT_RetireAge1":              "65",
				"eT_RetireAge2":              "",
				"eT_PlanThroughAge1":         "70",
				"eT_PlanThroughAge2":         "",
				"eT_PIA1":                    "20", // 20k
				"eT_PIA2":                    "",
				"eT_SS_Start1":               "70",
				"eT_SS_Start2":               "",
				"eT_TDRA1":                   "10", // 10k
				"eT_TDRA2":                   "",
				"eT_TDRA_Rate1":              "",
				"eT_TDRA_Rate2":              "",
				"eT_TDRA_Contrib1":           "",
				"eT_TDRA_Contrib2":           "",
				"eT_TDRA_ContribStartAge1":   "",
				"eT_TDRA_ContribStartAge2":   "",
				"eT_TDRA_ContribEndAge1":     "",
				"eT_TDRA_ContribEndAge2":     "",
				"eT_Roth1":                   "5", // 5k
				"eT_Roth2":                   "",
				"eT_Roth_Rate1":              "",
				"eT_Roth_Rate2":              "",
				"eT_Roth_Contrib1":           "",
				"eT_Roth_Contrib2":           "",
				"eT_Roth_ContribStartAge1":   "",
				"eT_Roth_ContribStartAge2":   "",
				"eT_Roth_ContribEndAge1":     "",
				"eT_Roth_ContribEndAge2":     "",
				"eT_Aftatax":                 "15", // 15k
				"eT_Aftatax_Rate":            "",
				"eT_Aftatax_Contrib":         "",
				"eT_Aftatax_ContribStartAge": "",
				"eT_Aftatax_ContribEndAge":   "",
			},
			row:             []float64{},
			suppressNewline: true,
			expectstr: "Row: [0 1 2 3 4 5 6 7 8 9 10 11 12 13 14 15 16 17 18 19 20 21 22 23 24 25 26 27 28 29 30 31 32 33 34 35 36 37 38 39 40 41 42 43 44 45 46 47 48 49 50 51 52 53 54 55 56 57 58 59 60 61 62 63 64 65 66 67 68 69 70 71 72 73 74 75 76 77 78 79 80 81 82 83 84 85 86 87 88 89 90 91 92 93 94 95 96 97 98 99 100 101 102 103 104 105 106 107 108 109 110 111 112 113 114 115 116 117 118 119 120 121 122 123 124 125 126 127 128 129 130 131 132 133 134 135 136 137 138 139 140 141 142 143 144 145 146 147 148 149 150 151 152 153 154 155 156 157 158 159 160 161 162 163 164 165 166 167 168 169 170 171 172 173 174 175 176 177 178 179 180 181 182 183 184 185 186 187 188 189 190 191 192 193 194 195 196 197 198 199 200 201 202 203 204 205 206 207 208 209 210 211 212 213 214 215 216 217 218 219 220 221 222]" +
				"x[0,1]= 1.000, x[0,2]= 2.000, x[0,3]= 3.000, x[0,4]= 4.000, x[0,5]= 5.000, x[0,6]= 6.000, x[1,0]= 7.000, x[1,1]= 8.000, x[1,2]= 9.000, x[1,3]=10.000, x[1,4]=11.000, x[1,5]=12.000, x[1,6]=13.000, x[2,0]=14.000, x[2,1]=15.000, x[2,2]=16.000, x[2,3]=17.000, x[2,4]=18.000, x[2,5]=19.000, x[2,6]=20.000, x[3,0]=21.000, x[3,1]=22.000, x[3,2]=23.000, x[3,3]=24.000, x[3,4]=25.000, x[3,5]=26.000, x[3,6]=27.000, x[4,0]=28.000, x[4,1]=29.000, x[4,2]=30.000, x[4,3]=31.000, x[4,4]=32.000, x[4,5]=33.000, x[4,6]=34.000, x[5,0]=35.000, x[5,1]=36.000, x[5,2]=37.000, x[5,3]=38.000, x[5,4]=39.000, x[5,5]=40.000, x[5,6]=41.000, x[6,0]=42.000, x[6,1]=43.000, x[6,2]=44.000, x[6,3]=45.000, x[6,4]=46.000, x[6,5]=47.000, x[6,6]=48.000, x[7,0]=49.000, x[7,1]=50.000, x[7,2]=51.000, x[7,3]=52.000, x[7,4]=53.000, x[7,5]=54.000, x[7,6]=55.000, x[8,0]=56.000, x[8,1]=57.000, x[8,2]=58.000, x[8,3]=59.000, x[8,4]=60.000, x[8,5]=61.000, x[8,6]=62.000, x[9,0]=63.000, x[9,1]=64.000, x[9,2]=65.000, x[9,3]=66.000, x[9,4]=67.000, x[9,5]=68.000, x[9,6]=69.000, x[10,0]=70.000, x[10,1]=71.000, x[10,2]=72.000, x[10,3]=73.000, x[10,4]=74.000, x[10,5]=75.000, x[10,6]=76.000, y[0,0]=77.000, y[0,1]=78.000, y[0,2]=79.000, y[1,0]=80.000, y[1,1]=81.000, y[1,2]=82.000, y[2,0]=83.000, y[2,1]=84.000, y[2,2]=85.000, y[3,0]=86.000, y[3,1]=87.000, y[3,2]=88.000, y[4,0]=89.000, y[4,1]=90.000, y[4,2]=91.000, y[5,0]=92.000, y[5,1]=93.000, y[5,2]=94.000, y[6,0]=95.000, y[6,1]=96.000, y[6,2]=97.000, y[7,0]=98.000, y[7,1]=99.000, y[7,2]=100.000, y[8,0]=101.000, y[8,1]=102.000, y[8,2]=103.000, y[9,0]=104.000, y[9,1]=105.000, y[9,2]=106.000, y[10,0]=107.000, y[10,1]=108.000, y[10,2]=109.000, w[0,0]=110.000, w[0,1]=111.000, w[0,2]=112.000, w[1,0]=113.000, w[1,1]=114.000, w[1,2]=115.000, w[2,0]=116.000, w[2,1]=117.000, w[2,2]=118.000, w[3,0]=119.000, w[3,1]=120.000, w[3,2]=121.000, w[4,0]=122.000, w[4,1]=123.000, w[4,2]=124.000, w[5,0]=125.000, w[5,1]=126.000, w[5,2]=127.000, w[6,0]=128.000, w[6,1]=129.000, w[6,2]=130.000, w[7,0]=131.000, w[7,1]=132.000, w[7,2]=133.000, w[8,0]=134.000, w[8,1]=135.000, w[8,2]=136.000, w[9,0]=137.000, w[9,1]=138.000, w[9,2]=139.000, w[10,0]=140.000, w[10,1]=141.000, w[10,2]=142.000, b[0,0]=143.000, b[0,1]=144.000, b[0,2]=145.000, b[1,0]=146.000, b[1,1]=147.000, b[1,2]=148.000, b[2,0]=149.000, b[2,1]=150.000, b[2,2]=151.000, b[3,0]=152.000, b[3,1]=153.000, b[3,2]=154.000, b[4,0]=155.000, b[4,1]=156.000, b[4,2]=157.000, b[5,0]=158.000, b[5,1]=159.000, b[5,2]=160.000, b[6,0]=161.000, b[6,1]=162.000, b[6,2]=163.000, b[7,0]=164.000, b[7,1]=165.000, b[7,2]=166.000, b[8,0]=167.000, b[8,1]=168.000, b[8,2]=169.000, b[9,0]=170.000, b[9,1]=171.000, b[9,2]=172.000, b[10,0]=173.000, b[10,1]=174.000, b[10,2]=175.000, b[11,0]=176.000, b[11,1]=177.000, b[11,2]=178.000, s[0]=179.000, s[1]=180.000, s[2]=181.000, s[3]=182.000, s[4]=183.000, s[5]=184.000, s[6]=185.000, s[7]=186.000, s[8]=187.000, s[9]=188.000, s[10]=189.000, D[0,0]=190.000, D[0,1]=191.000, D[0,2]=192.000, D[1,0]=193.000, D[1,1]=194.000, D[1,2]=195.000, D[2,0]=196.000, D[2,1]=197.000, D[2,2]=198.000, D[3,0]=199.000, D[3,1]=200.000, D[3,2]=201.000, D[4,0]=202.000, D[4,1]=203.000, D[4,2]=204.000, D[5,0]=205.000, D[5,1]=206.000, D[5,2]=207.000, D[6,0]=208.000, D[6,1]=209.000, D[6,2]=210.000, D[7,0]=211.000, D[7,1]=212.000, D[7,2]=213.000, D[8,0]=214.000, D[8,1]=215.000, D[8,2]=216.000, D[9,0]=217.000, D[9,1]=218.000, D[9,2]=219.000, D[10,0]=220.000, D[10,1]=221.000, D[10,2]=222.000, ",
			testcase: "counting",
		},
	}
	for i, elem := range tests {
		ip := NewInputParams(elem.ip)
		ti := NewTaxInfo(ip.filingStatus)
		taxbins := len(*ti.Taxtable)
		cgbins := len(*ti.Capgainstable)
		//fmt.Printf("ACCMAP: %v", ip.accmap)
		vindx, err := NewVectorVarIndex(ip.numyr, taxbins, cgbins, ip.accmap)
		if err != nil {
			t.Errorf("PrintModelRow case %d: %s", i, err)
			continue
		}
		numaccounts := 0
		for _, acc := range ip.accmap {
			numaccounts += acc
		}
		ms := ModelSpecs{
			ip:     ip,
			vindx:  vindx,
			ti:     ti,
			numyr:  ip.numyr,
			accmap: ip.accmap,
			numacc: numaccounts,
		}

		row := make([]float64, vindx.Vsize)
		switch elem.testcase {
		case "allones":
			for indx := 0; indx < vindx.Vsize; indx++ {
				row[indx] = float64(indx)
			}
		case "counting":
			for indx := 0; indx < vindx.Vsize; indx++ {
				row[indx] = float64(indx)
			}
		case "allzeros":
			// nothing to change
		default:
			t.Errorf("TestPrintModelRow: Unexpected test case '%s'\n", elem.testcase)
			continue
		}
		//fmt.Printf("Vsize: %d\n", vindx.Vsize)

		mychan := make(chan string)
		oldout, w, err := RedirectStdout(mychan)
		if err != nil {
			t.Errorf("RedirectStdout: %s\n", err)
			return // should this be continue?
		}
		fmt.Printf("Row: %v\n", row)
		ms.printModelRow(row, elem.suppressNewline)

		str := RestoreStdout(mychan, oldout, w)
		strn := stripWhitespace(str)
		strexpect := stripWhitespace(elem.expectstr)
		if strn != strexpect {
			t.Errorf("PrintModelRow Case %d: expected '%s', found '%s'", i, elem.expectstr, str)
		}
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

func TestRedirectOutput(t *testing.T) {
	mychan := make(chan string)
	oldout, w, err := RedirectStdout(mychan)
	if err != nil {
		t.Errorf("RedirectStdout: %s\n", err)
		return
	}
	outstr := "This will be captured for comparisons later\nAnd this too\n"
	fmt.Printf("%s", outstr)
	str := RestoreStdout(mychan, oldout, w)
	if str != outstr {
		t.Errorf("Capured output fails: expected '%s', found '%s'", outstr, str)
	}
}

func RedirectStdout(mechan chan string) (*os.File, *os.File, error) {
	oldStdout := os.Stdout
	readPipe, writePipe, err := os.Pipe()
	if err != nil {
		return os.Stdout, nil, err
	}
	os.Stdout = writePipe
	//mechan := make(chan string)
	go func() {
		var buf bytes.Buffer
		_, err := io.Copy(&buf, readPipe)
		readPipe.Close()
		if err != nil {
			fmt.Fprintf(os.Stderr, "func() copyPipe: %v\n", err)
			return
		}
		mechan <- buf.String()
	}()
	return oldStdout, writePipe, nil
}

func RestoreStdout(mechan chan string, oldStdout *os.File, writePipe *os.File) string {
	// Reset the output again
	writePipe.Close()
	os.Stdout = oldStdout
	str := <-mechan
	return str
}

func stripWhitespace(str string) string {
	return strings.Map(func(r rune) rune {
		if unicode.IsSpace(r) {
			// if the character is any white space, drop it
			return -1
		}
		// else keep it in the string
		return r
	}, str)
}
