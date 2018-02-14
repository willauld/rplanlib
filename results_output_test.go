package rplanlib

import (
	"fmt"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/willauld/lpsimplex"
)

//
// Testing for results_output.go
//

var sipSingle = map[string]string{
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

	"eT_iRate":    "2.5",
	"eT_rRate":    "6",
	"eT_maximize": "Spending", // or "PlusEstate"
}

var sipJoint = map[string]string{
	"setName":                    "activeParams",
	"filingStatus":               "joint",
	"key1":                       "retiree1",
	"key2":                       "retiree2",
	"eT_Age1":                    "54",
	"eT_Age2":                    "54",
	"eT_RetireAge1":              "65",
	"eT_RetireAge2":              "65",
	"eT_PlanThroughAge1":         "75",
	"eT_PlanThroughAge2":         "75",
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

	"eT_iRate":    "2.5",
	"eT_rRate":    "6",
	"eT_maximize": "Spending", // or "PlusEstate"
}

//def precheck_consistancy():
func TestPreCheckConsistancy(t *testing.T) {
	fmt.Printf("Not Yet Implemented\n")
}

//def consistancy_check(res, years, taxbins, cgbins, accounts, accmap, vindx):
func TestCheckConsistancy(t *testing.T) {
	fmt.Printf("Not Yet Implemented\n")
}

//func (ms ModelSpecs) activitySummaryHeader(fieldwidth int)
func TestActivitySummaryHeader(t *testing.T) {
	tests := []struct {
		sip    map[string]string
		expect string
	}{
		{
			sip:    sipJoint,
			expect: "retiree1/retiree2\n    age     fIRA    tIRA  RMDref   fRoth   tRoth fAftaTx tAftaTx   o_inc      SS Expense TFedTax Spndble",
		},
		{
			sip:    sipSingle,
			expect: "retiree1\n age     fIRA    tIRA  RMDref   fRoth   tRoth fAftaTx tAftaTx   o_inc      SS Expense TFedTax Spndble",
		},
	}
	for i, elem := range tests {
		ip := NewInputParams(elem.sip)
		//fmt.Printf("InputParams: %#v\n", ip)
		ti := NewTaxInfo(ip.filingStatus)
		taxbins := len(*ti.Taxtable)
		cgbins := len(*ti.Capgainstable)
		vindx, err := NewVectorVarIndex(ip.numyr, taxbins,
			cgbins, ip.accmap, os.Stdout)
		if err != nil {
			t.Errorf("TestActivitySummaryHeader case %d: %s", i, err)
			continue
		}
		csvfile := (*os.File)(nil)
		verbose := false
		allowDeposits := false
		logfile := os.Stdout
		ms := NewModelSpecs(vindx, ti, ip, verbose,
			allowDeposits, os.Stderr, logfile, csvfile, logfile)

		mychan := make(chan string)
		oldout, w, err := ms.RedirectModelSpecsTable(mychan)
		if err != nil {
			t.Errorf("RedirectModelSpecsTable: %s\n", err)
			return // should this be continue?
		}

		fieldwidth := 7
		ms.activitySummaryHeader(fieldwidth)

		str := ms.RestoreModelSpecsTable(mychan, oldout, w)
		strn := strings.TrimSpace(str)
		//strn := stripWhitespace(str)
		//warningmes := stripWhitespace(elem.warningmes)
		if elem.expect != strn {
			t.Errorf("TestActivitySummaryHeader case %d:  expected output:\n\t '%s'\n\tbut found:\n\t'%s'\n", i, elem.expect, strn)
		}

	}
	//fmt.Printf("Not Yet Implemented\n")
}

//func (ms ModelSpecs) printActivitySummary(xp *[]float64)
func TestActivitySummary(t *testing.T) {
	fmt.Printf("Not Yet Implemented\n")
}

//def print_income_expense_details():
func TestIncomeExpenseDetails(t *testing.T) {
	fmt.Printf("Not Yet Implemented\n")
}

//def print_account_trans(res):
func TestPrintAccountTrans(t *testing.T) {
	fmt.Printf("Not Yet Implemented\n")
}

//def print_tax(res):
func TestPrintTax(t *testing.T) {
	fmt.Printf("Not Yet Implemented\n")
}

//def print_tax_brackets(res):
func TestPrintTaxBrackets(t *testing.T) {
	fmt.Printf("Not Yet Implemented\n")
}

//def print_cap_gains_brackets(res):
func TestPrintCapGainsBrackets(t *testing.T) {
	fmt.Printf("Not Yet Implemented\n")
}

//func (ms ModelSpecs) depositAmount(xp *[]float64, year int, index int) float64
func TestDepositAmount(t *testing.T) {
	fmt.Printf("Not Yet Implemented\n")
}

//func (ms ModelSpecs) ordinaryTaxable(year int, xp *[]float64) float64
func TestOrdinaryTaxable(t *testing.T) {
	fmt.Printf("Not Yet Implemented\n")
}

//func (ms ModelSpecs) IncomeSummary(year int, xp *[]float64) (T, spendable, tax, rate, ncgtax, earlytax float64, rothearly bool)
func TestIncomeSummary(t *testing.T) {
	fmt.Printf("Not Yet Implemented\n")
}

//func (ms ModelSpecs) getResultTotals(xp *[]float64) (twithd, tcombined, tT, ttax, tcgtax, tearlytax, tspendable, tbeginbal, tendbal float64)

func TestGetResultTotals(t *testing.T) {
	fmt.Printf("Not Yet Implemented\n")
}

//func (ms ModelSpecs) printBaseConfig(xp *[]float64)  // input is res.x
func TestPrintBaseConfig(t *testing.T) {
	fmt.Printf("Not Yet Implemented\n")
}

//def verifyInputs( c , A , b ):
func TestVerifyInputs(t *testing.T) {
	fmt.Printf("Not Yet Implemented\n")
}

func TestResultsOutput(t *testing.T) {
	tests := []struct {
		ip            map[string]string
		verbose       bool
		allowDeposits bool
		iRate         float64
	}{
		{ // Case 0 // case to match mobile.toml
			ip: map[string]string{
				"setName":                    "activeParams",
				"filingStatus":               "joint",
				"key1":                       "retiree1",
				"key2":                       "retiree2",
				"eT_Age1":                    "54",
				"eT_Age2":                    "54",
				"eT_RetireAge1":              "65",
				"eT_RetireAge2":              "65",
				"eT_PlanThroughAge1":         "75",
				"eT_PlanThroughAge2":         "75",
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

				"eT_iRate":    "2.5",
				"eT_rRate":    "6",
				"eT_maximize": "Spending", // or "PlusEstate"
			},
			verbose:       true,
			allowDeposits: false,
			iRate:         1.025,
		},
	}
	for i, elem := range tests {
		ip := NewInputParams(elem.ip)
		//fmt.Printf("InputParams: %#v\n", ip)
		ti := NewTaxInfo(ip.filingStatus)
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

		ms.printActivitySummary(&res.X)
		/*
			//ms.print_model_results(res.x)
				        if args.verboseincome:
				            print_income_expense_details()
				        if args.verboseaccounttrans:
				            print_account_trans(res)
				        if args.verbosetax:
				            print_tax(res)
				        if args.verbosetaxbrackets:
				            print_tax_brackets(res)
							print_cap_gains_brackets(res)
		*/
		ms.printBaseConfig(&res.X)
	}
}
