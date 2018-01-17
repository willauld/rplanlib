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
		OK := my_check_index_sequence(elem.years, elem.taxbins,
			elem.cgbins, elem.accnum, elem.accmap, vvindex)
		if OK != true {
			t.Errorf("VectorVarIndex case %d: Failed\n", i)
		}
	}
}
