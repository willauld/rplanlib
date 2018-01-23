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
		accnum  int
		accmap  map[string]int
	}{
		{ // case 0
			years:   10,
			taxbins: 8,
			cgbins:  3,
			accnum:  3,
			accmap:  map[string]int{"IRA": 1, "Roth": 1, "Aftertax": 1},
		},
		{ // case 1
			years:   100,
			taxbins: 8,
			cgbins:  3,
			accnum:  5,
			accmap:  map[string]int{"IRA": 2, "Roth": 2, "Aftertax": 1},
		},
	}
	for i, elem := range tests {
		vvindex := NewVectorVarIndex(elem.years, elem.taxbins,
			elem.cgbins, elem.accnum, elem.accmap)
		OK := checkIndexSequence(elem.years, elem.taxbins,
			elem.cgbins, elem.accnum, elem.accmap, vvindex)
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

/*
func checkStrconvError(e error) {
	if e != nil {
		panic(e)
	}
}
*/
/*
	tests := []struct {
	}{
		{},
	}
	for i, elem := range tests {
	}
*/
func TestCheckStrconvError(t *testing.T) {
	tests := []struct {
		//err    error
		errstr string
	}{
		{errstr: "case 0"},
		{errstr: "case 1"},
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
					t.Errorf("checkStrcovError case %d.a should have panicked!", i)
				} else if elem.errstr == "" && r != nil {
					t.Errorf("checkStrcovError case %d.b should have panicked!", i)
				}
			}()
			// This function should cause a panic
			checkStrconvError(err)
		}()
	}
}

func TestMergeVectors(t *testing.T)      {}
func TestBuildVector(t *testing.T)       {}
func TestNewModelSpecs(t *testing.T)     {}
func TestBuildModel(t *testing.T)        {}
func TestAccountOwnerAge(t *testing.T)   {}
func TestMatchRetiree(t *testing.T)      {}
func TestCgTaxableFraction(t *testing.T) {}
func TestPrintModelMatrix(t *testing.T)  {}
func TestPrintConstraint(t *testing.T)   {}
func TestPrintModelRow(t *testing.T)     {}
