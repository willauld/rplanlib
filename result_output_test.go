package rplanlib

import (
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/willauld/lpsimplex"
)

//
// Testing for results_output.go
//

func TestResultsOutput(t *testing.T) {
	tests := []struct {
		ip            map[string]string
		verbose       bool
		allowDeposits bool
		iRate         float64
	}{
		{ // Case 0 // joint
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
			verbose:       true,
			allowDeposits: false,
			iRate:         1.025,
		},
	}
	for i, elem := range tests {
		ti := NewTaxInfo(elem.ip["filingStatus"])
		ip := NewInputParams(elem.ip)
		taxbins := len(*ti.Taxtable)
		cgbins := len(*ti.Capgainstable)
		vindx, err := NewVectorVarIndex(ip.numyr, taxbins,
			cgbins, ip.accmap, os.Stdout)
		if err != nil {
			t.Errorf("TestResultsOutput case %d: %s", i, err)
			continue
		}
		logfile, err := os.Create("ModelMatixPP.log")
		if err != nil {
			t.Errorf("TestResultsOutput case %d: %s", i, err)
			continue
		}
		csvfile := (*os.File)(nil)
		ms := NewModelSpecs(vindx, ti, ip, elem.verbose,
			elem.allowDeposits, os.Stderr, logfile, csvfile, logfile)
		/**/
		c, a, b, notes := ms.BuildModel()
		ms.printModelMatrix(c, a, b, notes, nil, false)
		/**/
		tol := 1.0e-7

		bland := false
		maxiter := 4000

		callback := lpsimplex.Callbackfunc(nil)
		//callback := lpsimplex.LPSimplexVerboseCallback
		//callback := lpsimplex.LPSimplexTerseCallback
		disp := true // false //true
		start := time.Now()
		res := lpsimplex.LPSimplex(c, a, b, nil, nil, nil, callback, disp, maxiter, tol, bland)
		elapsed := time.Since(start)
		var str string
		//str = fmt.Sprintf("Res: %+v\n", res)
		//fmt.Printf(str)
		//str = fmt.Sprintf("expeced opt: %v,          have opt: %v\n", expectOpt, res.Fun)
		//fmt.Printf(str)
		str = fmt.Sprintf("Message: %v\n", res.Message)
		fmt.Printf(str)
		str = fmt.Sprintf("Time: LPSimplex() took %s\n", elapsed)
		fmt.Printf(str)
		fmt.Printf("Calling LPSimplex() for m:%d x n:%d model\n", len(a), len(a[0]))

	}
}
