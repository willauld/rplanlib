package rplanlib

import (
	"fmt"
	"math"
	"testing"
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
			accmap:  map[string]int{"IRA": 1, "Roth": 1, "Aftertax": 1},
		},
		{ // case 1
			years:   100,
			taxbins: 8,
			cgbins:  3,
			accmap:  map[string]int{"IRA": 2, "Roth": 2, "Aftertax": 1},
		},
	}
	for i, elem := range tests {
		vvindex := NewVectorVarIndex(elem.years, elem.taxbins,
			elem.cgbins, elem.accmap)
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
				t.Errorf("GetIPIntValue() case %d: Failed - Expected %d but found %d\n", i, elem.expect, val)
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
			startPlan:    64,
			endPlan:      101,
			numyr:        37,
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
			prePlanYears: 0,
			startPlan:    65,
			endPlan:      98,
			numyr:        33,
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
			accmap: map[string]int{"IRA": 1, "Roth": 1, "Aftatax": 1},
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
			verbose:       false,
			allowDeposits: false,
			iRate:         1.025,
		},
		{ // Case 1 // mseparate
			years:  10,
			accmap: map[string]int{"IRA": 1, "Roth": 1, "Aftatax": 1},
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
			verbose:       false,
			allowDeposits: false,
			iRate:         1.025,
		},
		{ // Case 2 // single
			years:  10,
			accmap: map[string]int{"IRA": 1, "Roth": 1, "Aftatax": 1},
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
		accnum := 0
		for _, acc := range elem.accmap { //TODO FIXME accmap should come from ImputParams
			accnum += acc
		}
		vindx := NewVectorVarIndex(ip.numyr, taxbins, cgbins, elem.accmap)
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
			years:         10,
			accmap:        map[string]int{"IRA": 1, "Roth": 1, "Aftatax": 1},
			ip:            map[string]string{"filingStatus": "joint"},
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
		accnum := 0
		for _, acc := range elem.accmap {
			accnum += acc
		}
		vindx := NewVectorVarIndex(elem.years, taxbins,
			cgbins, elem.accmap)
		ms := NewModelSpecs(vindx, ti, ip, elem.verbose,
			elem.allowDeposits)
		/**/
		c, A, b := ms.BuildModel()
		ms.printModelMatrix(c, A, b, nil, false)
		/**/
		if ms.iRate != elem.iRate {
			t.Errorf("NewModelSpecs case %d: iRate expected %f, found %f", i, elem.iRate, ms.iRate)
		}
	}
}

func TestAccountOwnerAge(t *testing.T) {
	/*
			tests := []struct {
			}{
				{},
			}
			for i, elem := range tests {
		func (ms ModelSpecs) accountOwnerAge(year int, acc account) int {
			}
	*/
}

func TestMatchRetiree(t *testing.T) { /* TODO:FIXME:IMPLEMENTME */ }

func TestCgTaxableFraction(t *testing.T) { /* TODO:FIXME:IMPLEMENTME */ }

func TestPrintModelMatrix(t *testing.T) { /* TODO:FIXME:IMPLEMENTME */ }

func TestPrintConstraint(t *testing.T) { /* TODO:FIXME:IMPLEMENTME */ }

func TestPrintModelRow(t *testing.T) {
	/*
		tests := []struct {
			row []float64
			suppressNewline bool
		}{
			{},
		}
		for i, elem := range tests {
			ms.printModelRow(elem.row, elem.suppressNewline)
		}
	*/
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
				"eT_Age1 INTEGER":            "",
				"eT_Age2 INTEGER":            "",
				"eT_RetireAge1 INTEGER":      "",
				"eT_RetireAge2 INTEGER":      "",
				"eT_PlanThroughAge1 INTEGER": "",
				"eT_PlanThroughAge2 INTEGER": "",
				"eT_PIA1 INTEGER":            "",
				"eT_PIA2 INTEGER":            "",
				"eT_SS_Start1 INTEGER":       "",
				"eT_SS_Start2 INTEGER":       "",
				"eT_TDRA1 INTEGER":           "",
				"eT_TDRA2 INTEGER":           "",
				"eT_TDRA_Rate1 REAL":         "",
				"eT_TDRA_Rate2 REAL":         "",
				"eT_TDRA_Contrib1 INTEGER":   "",
				"eT_TDRA_Contrib2 INTEGER":   "",
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
