package rplanlib

import (
	"fmt"
	"math"
	"testing"
)

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
		status, err := verifyFilingStatus(elem.filingStatus)
		if err != nil {
			fmt.Printf("TestNewModelSpecs: %s\n", err)
			continue
		}
		ti := NewTaxInfo(status, 2017)
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
		{ // case 1
			filingStatus: "mseparate",
			retireeindx:  1,
			year:         5,
			irate:        1.025,
		},
		{ // case 2
			filingStatus: "joint",
			retireeindx:  1,
			year:         5,
			irate:        1.025,
		},
		{ // case 3
			filingStatus: "joint",
			retireeindx:  -1, // pass in an empty key
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
		},
		{ // retireeindx == 1
			age:        54,
			ageAtStart: 55,
			throughAge: 100,
			mykey:      "retiree2",
		},
	}
	for i, elem := range tests {
		status, err := verifyFilingStatus(elem.filingStatus)
		if err != nil {
			fmt.Printf("TestNewModelSpecs: %s\n", err)
			continue
		}
		ti := NewTaxInfo(status, 2017)
		retireekey := ""
		if elem.retireeindx > 0 {
			retireekey = retirees[elem.retireeindx].mykey
		}
		prePlanYears := retirees[0].ageAtStart - retirees[0].age
		m := ti.maxContribution(elem.year, prePlanYears+elem.year,
			retirees, retireekey, elem.irate)
		//fmt.Printf("m: %f, year: %d, prePlanYears: %d, key: %s, irate: %f\n", m, elem.year, prePlanYears, retiree.mykey, elem.irate)
		inflateYears := prePlanYears + elem.year
		memax := ti.Contribspecs["TDRA"] + ti.Contribspecs["TDRACatchup"]
		emax := memax * math.Pow(elem.irate, float64(inflateYears)) // adjust for inflation ??? current ????
		if retireekey == "" {
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
		status, err := verifyFilingStatus(elem.filingStatus)
		if err != nil {
			fmt.Printf("TestNewModelSpecs: %s\n", err)
			continue
		}
		ti := NewTaxInfo(status, 2017)
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
		status, err := verifyFilingStatus(elem.filingStatus)
		if err != nil {
			fmt.Printf("TestNewModelSpecs: %s\n", err)
			continue
		}
		ti := NewTaxInfo(status, 2017)
		response := ti.rmdNeeded(elem.year, elem.retireer)
		if response != elem.response {
			t.Errorf("rmdNeeded case %d: Failed - Expected %v but found %v\n", i, elem.response, response)
		}
	}

}
