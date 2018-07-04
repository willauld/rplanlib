package rplanlib_test

import (
	"fmt"
	"os"
	"strings"
	//"regexp"
	//"strings"
	"path/filepath"
	"testing"
	"time"

	"github.com/willauld/lpsimplex"
	"github.com/willauld/rplanlib"
)

func TestE2E(t *testing.T) {
	/*
		tests := []struct {
			ip            map[string]string
			verbose       bool
			allowDeposits bool
			iRate         float64
		}{
			{ // Case 3 // case to match mobile.toml
				ip: map[string]string{
					"setName":                    "activeParams",
					"filingStatus":               "single",
					"key1":                       "retiree1",
					"key2":                       "",
					"eT_Age1":                    "54",
					"eT_Age2":                    "",
					"eT_RetireAge1":              "65",
					"eT_RetireAge2":              "",
					"eT_PlanThroughAge1":         "75",
					"eT_PlanThroughAge2":         "",
					"eT_PIA1":                    "",
					"eT_PIA2":                    "",
					"eT_SS_Start1":               "",
					"eT_SS_Start2":               "",
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
					"eT_Aftatax":                 "",
					"eT_Aftatax_Rate":            "",
					"eT_Aftatax_Contrib":         "",
					"eT_Aftatax_ContribStartAge": "",
					"eT_Aftatax_ContribEndAge":   "",

					"eT_iRatePercent": "2.5",
					"eT_rRatePercent": "6",
					"eT_maximize":     "Spending", // or "PlusEstate"
				},
				verbose:       true,
				allowDeposits: false,
				iRate:         1.025,
			},
		}
	*/
	if !(testing.Short() && testing.Verbose()) { //Skip unless set "-v -short"
		t.Skip("TestResultsOutput() (full runs): skipping unless set '-v -short'")
	}
	paramfiles, err := filepath.Glob("./testdata/strmap/*.strmap")
	if err != nil {
		t.Errorf("TestE2E Error: %s", err)
	}
	for i, ifile := range paramfiles {
		ifilecore := strings.TrimSuffix(filepath.Base(ifile), filepath.Ext(ifile))
		ifileext := filepath.Ext(ifile)
		//ifiledirpath := strings.TrimSuffix(ifile, filepath.Base(ifile))
		if i == -1 {
			//fmt.Printf("Testing input files, i: %d, ifile: %s\n", i, ifile)
			//fmt.Printf("core: %s\n", ifilecore)
			//fmt.Printf("dirpath: %s\n", ifiledirpath)
			break
		}
		fmt.Printf("======== CASE %d - %s ========\n", i, ifile)
		var ipsmp *map[string]string

		// ifile can be .toml or .strmap, Toml file is assumed
		if filepath.Ext(ifile) == ".strmap" {
			ipsmp, err = rplanlib.GetInputStrStrMapFromFile(ifile)
		} else {
			ipsmp, err = rplanlib.GetInputStringsMapFromToml(ifile)
		}
		if err != nil {
			t.Errorf("reading file (%s): %s", ifile, err)
			//printMsgAndExit(msgList, e)
		}
		ip, err := rplanlib.NewInputParams(*ipsmp, nil)
		if err != nil {
			t.Errorf("TestE2E case %d: %s", i, err)
			continue
		}
		//fmt.Printf("InputParams: %#v\n", ip)
		ti := rplanlib.NewTaxInfo(ip.FilingStatus, 2017)
		taxbins := len(*ti.Taxtable)
		cgbins := len(*ti.Capgainstable)
		vindx, err := rplanlib.NewVectorVarIndex(ip.Numyr, taxbins,
			cgbins, ip.Accmap, os.Stdout)
		if err != nil {
			t.Errorf("TestE2E case %d: %s", i, err)
			continue
		}
		logname := "./testdata/" + ifileext[1:] + "_test_output/" + ifilecore + ".log"
		logfile, err := os.Create(logname)
		if err != nil {
			t.Errorf("TestE2E case %d: %s", i, err)
			continue
		}
		//csvfile := (*os.File)(nil)
		csvname := "./testdata/" + ifileext[1:] + "_test_output/" + ifilecore + ".csv"
		csvfile, err := os.Create(csvname)
		if err != nil {
			t.Errorf("TestE2E case %d: %s", i, err)
			continue
		}
		RoundToOneK := false
		allowDeposits := false
		ms, err := rplanlib.NewModelSpecs(vindx, ti, *ip,
			allowDeposits, RoundToOneK, false, false,
			os.Stderr, logfile, csvfile, logfile, nil)
		if err != nil {
			t.Errorf("TestE2E case %d: %s", i, err)
			continue
		}
		//fmt.Printf("ModelSpecs: %#v\n", ms)

		c, a, b, notes := ms.BuildModel()

		Optstart := time.Now()
		aprime, bprime, oinfo := ms.OptimizeLPModel(&a, &b)
		Optelapsed := time.Since(Optstart)

		ms.PrintModelMatrix(c, a, b, notes, nil, false, oinfo) // TODO FIXME need to make this print somewhere else for examining the optimized model

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

		Ostart := time.Now()
		resPrime := lpsimplex.LPSimplex(c, *aprime, *bprime, nil, nil, nil, callback, disp, maxiter, tol, bland)
		Oelapsed := time.Since(Ostart)
		/*
			err = BinDumpModel(c, a, b, res.X, "./RPlanModelgo.datX")
			if err != nil {
				t.Errorf("TestE2E case %d: %s", i, err)
				continue
			}
			BinCheckModelFiles("./RPlanModelgo.datX", "./RPlanModelpython.datX", &vindx)
		*/

		//fmt.Printf("Res: %#v\n", res)
		str := fmt.Sprintf("Message: %v\n", res.Message)
		fmt.Printf(str)
		str = fmt.Sprintf("Message ResPrime: %v\n", resPrime.Message)
		fmt.Printf(str)
		str = fmt.Sprintf("Time: LPSimplex() took %s\n", elapsed)
		fmt.Printf(str)
		str = fmt.Sprintf("Time: Opt took %s, LPSimplex() took %s\n", Optelapsed, Oelapsed)
		fmt.Printf(str)
		fmt.Printf("Called LPSimplex() for m:%d x n:%d model\n", len(a), len(a[0]))
		if res.Success {
			//OK := ms.ConsistencyCheck(os.Stdout, &res.X)
			OK := ms.ConsistencyCheck(logfile, &res.X)
			if !OK {
				t.Errorf("TestE2E case %d: ConsistencyCheck() found issues with %s", i, ifilecore+ifileext)
			}

			ms.PrintActivitySummary(&res.X)
			ms.PrintIncomeExpenseDetails()
			ms.PrintAccountTrans(&res.X)
			ms.PrintTax(&res.X)
			ms.PrintTaxBrackets(&res.X)
			ms.PrintShadowTaxBrackets(&res.X)
			ms.PrintCapGainsBrackets(&res.X)
			ms.PrintAssetSummary()
			ms.PrintBaseConfig(&res.X)

			ms.PrintAccountWithdrawals(&res.X) // TESTING TESTING TESTING FIXME TODO
		}
		//createDefX(&res.X)
	}
}
