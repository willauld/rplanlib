package rplanlib

import (
	"testing"
)

func TestVectorVarIndex(t *testing.T) {

	iyears := 10
	itaxbins := 8
	icgbins := 3
	iaccounts := 3
	iaccmap := map[string]int{"IRA": 1, "Roth": 1, "Aftertax": 1}
	vvindex := NewVectorVarIndex(iyears, itaxbins, icgbins, iaccounts, iaccmap)
	OK := my_check_index_sequence(iyears, itaxbins, icgbins, iaccounts, iaccmap, vvindex)
	if OK != true {
		t.Errorf("VectorVarIndex: Failed\n")
	}
}
