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

	csvtag "github.com/artonge/go-csv-tag"
	"github.com/willauld/lpsimplex"
	"github.com/willauld/rplanlib"
)

type errorCode int

type acase struct {
	Testfile         string    `csv:"Testfile"`
	Logfile          string    `csv:"Logfile"`
	ErrorType        errorCode `csv:"Error Type"`
	SpendableAtLeast int       `csv:"Spendable At Least"`
	ModelM           int       `csv:"Model M"`
	ModelN           int       `csv:"Model N"`
	Iterations       int       `csv:"Iterations"`
}

type errorType struct {
	ErrorType errorCode `csv:"Error Code"`
	ErrorStr  string    `csv:"Error String"`
	CustomStr string    `csv:"Custom String"`
}

func createErrorTypeCSV() {
	errorTypeTable := []errorType{
		{
			ErrorType: 0,
			ErrorStr:  "no error",
		},
		{
			ErrorType: 1,
			ErrorStr:  "check spenable and no aftertax accounts",
		},
		{
			ErrorType: 2,
			ErrorStr:  "check spenable with aftertax accounts",
		},
		{
			ErrorType: 3,
			ErrorStr:  "configuration input error",
			CustomStr: "checkNames: name",
		},
	}
	err := csvtag.DumpToFile(errorTypeTable, "testdata/errortypes.csv")
	if err != nil {
		// cry
	}
}

func TestE2E(t *testing.T) {

	/*
		if !(testing.Short() && testing.Verbose()) { //Skip unless set "-v -short"
			t.Skip("TestResultsOutput() (full runs): skipping unless set '-v -short'")
		}
	*/
	//
	// Define the local Testing options
	//
	DisplayOutputAndTiming := true   //false  //true
	DoModelOptimizationTest := false //true //false
	DoScaleModel := true             // false            // true
	updateExpectFile := false        // true
	updateExpectFileInterationCounts := false
	updateExpectFileSpendableAtLeast := false
	updateExpectFileModelMxN := false
	updateExpectFileLogName := false
	//updateExpectFileExpectedError := false // Not good to automate this one
	ExecuteOnlyCase := -1 // -1 for all cases OR specific case number
	//
	// Bring this back in to make sure all configuration files are
	// accounted for. Need to code this up
	//
	strmapfiles, err := filepath.Glob("./testdata/strmap/*.strmap")
	if err != nil {
		t.Errorf("TestE2E Error: %s", err)
	}
	tomlfiles, err := filepath.Glob("./testdata/toml/*.toml")
	if err != nil {
		t.Errorf("TestE2E Error: %s", err)
	}
	paramfiles := append(tomlfiles, strmapfiles...)

	cases := []acase{}
	err = csvtag.Load(csvtag.Config{ // Load your csv with configuration
		Path: "testdata/expect.csv", // Path of the csv file
		Dest: &cases,                // A pointer to the create slice
	})
	for _, ifile := range paramfiles {
		match := false
		for _, thiscase := range cases {
			if ifile == thiscase.Testfile {
				match = true
				continue
			}
		}
		if match == false {
			var lfile string
			base := filepath.Base(ifile)
			extention := filepath.Ext(ifile)
			if extention == ".strmap" {
				lfile = "./testdata/strmap_test_output/" + base + ".log"

			} else if extention == ".toml" {
				lfile = "./testdata/toml_test_output/" + base + ".log"
			} else {
				// Error
			}
			c := acase{
				Testfile:         ifile,
				Logfile:          lfile,
				ErrorType:        0,
				SpendableAtLeast: 0,
			}
			cases = append(cases, c)
			err = csvtag.DumpToFile(cases, "testdata/expect.csv")
		}
	}
	//createErrorTypeCSV() // Uncomment to update errortypes.csv
	errorTypeTable := []errorType{}
	err = csvtag.Load(csvtag.Config{ // Load your csv with configuration
		Path: "testdata/errortypes.csv", // Path of the csv file
		Dest: &errorTypeTable,           // A pointer to the create slice
	})

	for i, curCase := range cases {
		if ExecuteOnlyCase >= 0 && i != ExecuteOnlyCase {
			continue
		}
		ifilecore := strings.TrimSuffix(filepath.Base(curCase.Testfile), filepath.Ext(curCase.Testfile))
		ifileext := filepath.Ext(curCase.Testfile)

		fmt.Printf("======== CASE %d - %s ========\n", i, curCase.Testfile)

		var ipsmp *map[string]string
		msgList := rplanlib.NewWarnErrorList()
		// curCase.testfile can be .toml or .strmap, Toml file is assumed
		if filepath.Ext(curCase.Testfile) == ".strmap" {
			ipsmp, err = rplanlib.GetInputStrStrMapFromFile(curCase.Testfile)
		} else {
			ipsmp, err = rplanlib.GetInputStringsMapFromToml(curCase.Testfile)
		}
		if err != nil {
			if curCase.ErrorType == 3 && strings.Contains(err.Error(), errorTypeTable[3].CustomStr) {
				// expected error
				continue
			}
			t.Errorf("TestE2E case %d: configuration file error (%s): %s", i, curCase.Logfile, err)
			rplanlib.PrintAndClearMsg(os.Stdout, msgList)
			continue
		}
		ip, err := rplanlib.NewInputParams(*ipsmp, msgList)
		if err != nil {
			t.Errorf("TestE2E case %d: %s", i, err)
			rplanlib.PrintAndClearMsg(os.Stdout, msgList)
			continue
		}
		//fmt.Printf("InputParams: %#v\n", ip)
		taxYear := 2018
		ti := rplanlib.NewTaxInfo(ip.FilingStatus, taxYear)
		taxbins := len(*ti.Taxtable)
		cgbins := len(*ti.Capgainstable)
		vindx, err := rplanlib.NewVectorVarIndex(ip.Numyr, taxbins,
			cgbins, ip.Accmap, os.Stdout)
		if err != nil {
			t.Errorf("TestE2E case %d: %s", i, err)
			continue
		}
		logname := curCase.Logfile
		if logname == "" {
			logname = "./testdata/" + ifileext[1:] + "_test_output/" + ifilecore + ".log"
		}
		if updateExpectFile || updateExpectFileLogName {
			cases[i].Logfile = logname
			curCase.Logfile = logname
		}
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
		developerInfo := true
		fourPercentRule := false
		ms, err := rplanlib.NewModelSpecs(vindx, ti, *ip,
			RoundToOneK, developerInfo, fourPercentRule,
			os.Stderr, logfile, csvfile, logfile, msgList)
		if err != nil {
			t.Errorf("TestE2E case %d: %s", i, err)
			rplanlib.PrintAndClearMsg(logfile, msgList)
			continue
		}
		//fmt.Printf("ModelSpecs: %#v\n", ms)

		c, a, b, notes := ms.BuildModel()
		//c, a, b, _ := ms.BuildModel()

		var aprime *[][]float64
		var bprime *[]float64
		var Optelapsed time.Duration
		if DoModelOptimizationTest {
			Optstart := time.Now()
			var oinfo *[]rplanlib.OptInfo
			aprime, bprime, oinfo = ms.OptimizeLPModel(&a, &b)
			//aprime, bprime, _ = ms.OptimizeLPModel(&a, &b)
			Optelapsed = time.Since(Optstart)

			ms.PrintModelMatrix(c, a, b, notes, nil, false, oinfo) // TODO FIXME need to make this print somewhere else for examining the optimized model
		}

		tol := 1.0e-7

		bland := false
		maxiter := 4000

		if DoScaleModel {
			//lpsimplex.LPSimplexSetNewBehavior(lpsimplex.NB_CMD_RESET | lpsimplex.NB_CMD_SCALEME | lpsimplex.NB_CMD_SCALEME_PIV_DIFF)
			lpsimplex.LPSimplexSetNewBehavior(lpsimplex.NB_CMD_RESET | lpsimplex.NB_CMD_SCALEME)
		} else {
			lpsimplex.LPSimplexSetNewBehavior(lpsimplex.NB_CMD_RESET)
		}

		callback := lpsimplex.Callbackfunc(nil)
		//callback := lpsimplex.LPSimplexVerboseCallback
		//callback := lpsimplex.LPSimplexTerseCallback
		disp := false //true
		start := time.Now()
		res := lpsimplex.LPSimplex(c, a, b, nil, nil, nil, callback, disp, maxiter, tol, bland)
		elapsed := time.Since(start)

		var Oelapsed time.Duration
		var resPrime lpsimplex.OptResult
		if DoModelOptimizationTest {
			Ostart := time.Now()
			resPrime = lpsimplex.LPSimplex(c, *aprime, *bprime, nil, nil, nil, callback, disp, maxiter, tol, bland)
			Oelapsed = time.Since(Ostart)
		}
		/*
			err = BinDumpModel(c, a, b, res.X, "./RPlanModelgo.datX")
			if err != nil {
				t.Errorf("TestE2E case %d: %s", i, err)
				continue
			}
			BinCheckModelFiles("./RPlanModelgo.datX", "./RPlanModelpython.datX", &vindx)
		*/

		if DisplayOutputAndTiming || DoModelOptimizationTest {
			//fmt.Printf("Res: %#v\n", res)
			str := fmt.Sprintf("Message: %v\n", res.Message)
			fmt.Printf(str)
			if DoModelOptimizationTest {
				str = fmt.Sprintf("Message ResPrime: %v\n", resPrime.Message)
				fmt.Printf(str)
				if res.Fun == resPrime.Fun {
					fmt.Printf("Object functions match: %f\n", res.Fun)
				} else {
					fmt.Printf("Object functions DO NOT match, Standard: %f, Optimized: %f\n", res.Fun, resPrime.Fun)
				}
			}
			str = fmt.Sprintf("Time: LPSimplex() took %s\n", elapsed)
			fmt.Printf(str)
			fmt.Printf("\tIterations: %d\n", res.Nitr)
			if DoModelOptimizationTest {
				str = fmt.Sprintf("Time: Opt took %s, LPSimplex() took %s\n", Optelapsed, Oelapsed)
				fmt.Printf(str)
			}
			fmt.Printf("Called LPSimplex() for m:%d x n:%d model\n", len(a), len(a[0]))
		}
		if res.Success {
			//OK := ms.ConsistencyCheck(os.Stdout, &res.X)
			OK := ms.ConsistencyCheckBrackets(&res.X)
			if !OK {
				t.Errorf("TestE2E case %d: Check Brackets found issues with %s", i, curCase.Logfile)
			}
			OK = ms.ConsistencyCheckSpendable(&res.X)
			if !(OK || curCase.ErrorType != 0) {
				if !OK {
					if ms.Ip.Accmap[rplanlib.Aftertax] > 0 {
						foundError := 2
						if curCase.ErrorType != 2 {
							t.Errorf("TestE2E case %d: %s for file %s", i, errorTypeTable[foundError].ErrorStr, curCase.Logfile)
						}
					} else {
						foundError := 1
						if curCase.ErrorType != 1 {
							// actual error does not match expected error
							t.Errorf("TestE2E case %d: %s for file %s", i, errorTypeTable[foundError].ErrorStr, curCase.Logfile)
						}
					}
				} else {
					t.Errorf("TestE2E case %d: did not generate expected error for file %s", i, curCase.Logfile)
				}
			} else {
				if curCase.ErrorType != 0 {
					// actual error does not match expected error
					t.Errorf("TestE2E case %d: Expected %s but DID NOT find it for file %s", i, errorTypeTable[curCase.ErrorType].ErrorStr, curCase.Logfile)
				}
			}
			//
			// Check expected Spendable
			//
			s := res.X[ms.Vindx.S(0)]
			if updateExpectFile ||
				updateExpectFileSpendableAtLeast {
				newVal := int(s - 5.0)
				cases[i].SpendableAtLeast = newVal
				curCase.SpendableAtLeast = newVal
			}
			if s < float64(curCase.SpendableAtLeast) {
				t.Errorf("TestE2E case %d: first year spendable (%6.0f) is less than expected (%d diff of %6.0f) for file %s", i, s, curCase.SpendableAtLeast, float64(curCase.SpendableAtLeast)-s, curCase.Logfile)
			}
			//
			// Check expected model size
			//
			m := len(a)
			n := len(a[0])
			if updateExpectFile ||
				updateExpectFileModelMxN {
				cases[i].ModelM = m
				curCase.ModelM = m
				cases[i].ModelN = n
				curCase.ModelN = n
			}
			if m != curCase.ModelM || n != curCase.ModelN {
				t.Errorf("TestE2E case %d: Expected m x n model of %d x %d but found %d x %d", i, curCase.ModelM, curCase.ModelN, m, n)
			}
			//
			// Check expected iteration count
			//
			nitr := res.Nitr
			if updateExpectFile || updateExpectFileInterationCounts {
				cases[i].Iterations = nitr
				curCase.Iterations = nitr
			}
			if nitr != curCase.Iterations {
				t.Errorf("TestE2E case %d: Expected %d Iterations but found %d", i, curCase.Iterations, nitr)
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
			ms.PrintObjectFunctionSolution(c, res.X)

			ModelAllBinding := true
			if ModelAllBinding {
				var bindingOnly bool
				slack := []float64(nil)
				if res.Success {
					slack = res.Slack
				}
				//bindingOnly = true
				//if !res.Success {
				bindingOnly = false
				//}
				ms.PrintModelMatrix(c, a, b, notes, slack, bindingOnly, nil)
			}
		} else {
			str := fmt.Sprintf("Message: %v\n", res.Message)
			ms.Ao.Output(str)
			str = fmt.Sprintf("Time: LPSimplex() took %s\n", elapsed)
			ms.Ao.Output(str)
			str = fmt.Sprintf("\tIterations: %d\n", res.Nitr)
			ms.Ao.Output(str)
			str = fmt.Sprintf("Called LPSimplex() for m:%d x n:%d model\n", len(a), len(a[0]))
			ms.Ao.Output(str)
			ms.Ao.Output("LPSimplex failed\n")
			varIndex := lpsimplex.LPSimplexNewBehaviorGetUnboundedVarNum()
			str = fmt.Sprintf(", unbounded at %s", ms.Vindx.Varstr(varIndex))
			if len(str) < 16 {
				str = ""
			}
			t.Errorf("TestE2E case %d: Unexpected simplex failure after %d iterations with msg: %s%s", i, res.Nitr, res.Message, str)

			ms.PrintModelMatrix(c, a, b, notes, nil, false, nil)
		}
		//createDefX(&res.X)
	}
	if updateExpectFile ||
		updateExpectFileInterationCounts ||
		updateExpectFileSpendableAtLeast ||
		updateExpectFileLogName ||
		updateExpectFileModelMxN {
		err = csvtag.DumpToFile(cases, "testdata/expect.csv")
		if err != nil {
			t.Errorf("***** Update of expect.csv failed *******")
			// TODO FIXME
		}
	}
}
