package rplanlib

import (
	"os"
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
			elem.cgbins, elem.accmap, os.Stdout)
		if err != nil {
			t.Errorf("VectorVarIndex case %d: %s", i, err)
			continue
		}
		OK := checkIndexSequence(elem.years, elem.taxbins,
			elem.cgbins, elem.accmap, vvindex, os.Stdout)
		if OK != true {
			t.Errorf("VectorVarIndex case %d: Failed\n", i)
		}
	}
}
