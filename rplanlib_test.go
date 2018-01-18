package rplanlib

import (
	"testing"
)

func TestVectorVarIndex(t *testing.T) {
	tests := []struct {
		years   int
		taxbins int
		cgbins  int
		accnum  int
		accmap  map[string]int
	}{
		{
			years:   10,
			taxbins: 8,
			cgbins:  3,
			accnum:  3,
			accmap:  map[string]int{"IRA": 1, "Roth": 1, "Aftertax": 1},
		},
		{
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

func TestTaxinfo(t *testing.T) {
	tests := []struct {
		filingStatus string
		//spot check info
		brackets          int
		thirdBracketStart float64
	}{
		{
			filingStatus:      "single",
			brackets:          7,
			thirdBracketStart: 37950,
		},
		{
			filingStatus:      "joint",
			brackets:          7,
			thirdBracketStart: 75900,
		},
		{
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
