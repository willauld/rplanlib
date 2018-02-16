package rplanlib

import (
	"fmt"
	"os"
	//"regexp"
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
			sip: sipJoint,
			expect: `retiree1/retiree2
    age     fIRA    tIRA  RMDref   fRoth   tRoth fAftaTx tAftaTx   o_inc      SS Expense TFedTax Spndble`,
		},
		{
			sip: sipSingle,
			expect: `retiree1
 age     fIRA    tIRA  RMDref   fRoth   tRoth fAftaTx tAftaTx   o_inc      SS Expense TFedTax Spndble`,
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
		DoNothing := false //true
		oldout, w, err := ms.RedirectModelSpecsTable(mychan, DoNothing)
		if err != nil {
			t.Errorf("RedirectModelSpecsTable: %s\n", err)
			return // should this be continue?
		}

		fieldwidth := 7
		ms.activitySummaryHeader(fieldwidth)

		str := ms.RestoreModelSpecsTable(mychan, oldout, w, DoNothing)
		strn := strings.TrimSpace(str)
		//strn := stripWhitespace(str)
		//warningmes := stripWhitespace(elem.warningmes)
		if elem.expect != strn {
			t.Errorf("TestActivitySummaryHeader case %d:  expected output:\n\t '%s'\n\tbut found:\n\t'%s'\n", i, elem.expect, strn)
		}
	}
}

//func (ms ModelSpecs) printActivitySummary(xp *[]float64)
func TestActivitySummary(t *testing.T) {
	tests := []struct {
		expect string
		sip    map[string]string
		xp     *[]float64
	}{
		{ // Case 0
			expect: `Activity Summary:

 retiree1
 age     fIRA    tIRA  RMDref   fRoth   tRoth fAftaTx tAftaTx   o_inc      SS Expense TFedTax Spndble
  65:   40594       0       0       0       0       0       0       0       0       0    3431   37164 
  66:   41609       0       0       0       0       0       0       0       0       0    3516   38093 
  67:   42650       0       0       0       0       0       0       0       0       0    3604   39045 
  68:   43716       0       0       0       0       0       0       0       0       0    3694   40021 
  69:   44809       0       0       0       0       0       0       0       0       0    3787   41022 
  70:   45929       0    9263       0       0       0       0       0       0       0    3881   42048 
  71:   47077       0    8315       0       0       0       0       0       0       0    3978   43099 
  72:   48254       0    7174       0       0       0       0       0       0       0    4078   44176 
  73:   49460       0    5811       0       0       0       0       0       0       0    4180   45281 
  74:   50697       0    4190       0       0       0       0       0       0       0    4284   46413 
  75:   51964       0    2269       0       0       0       0       0       0       0    4391   47573 
 retiree1
 age     fIRA    tIRA  RMDref   fRoth   tRoth fAftaTx tAftaTx   o_inc      SS Expense TFedTax Spndble`,
			sip: sipSingle,
			xp:  xpSingle,
		},
		{ // Case 1
			expect: `Activity Summary:

retiree1/retiree2
    age     fIRA    tIRA  RMDref   fRoth   tRoth fAftaTx tAftaTx   o_inc      SS Expense TFedTax Spndble
 65/ 65:   40594       0       0       0       0       0       0       0       0       0    1330   39264 
 66/ 66:   41609       0       0       0       0       0       0       0       0       0    1364   40246 
 67/ 67:   42650       0       0       0       0       0       0       0       0       0    1398   41252 
 68/ 68:   43716       0       0       0       0       0       0       0       0       0    1433   42283 
 69/ 69:   44809       0       0       0       0       0       0       0       0       0    1468   43340 
 70/ 70:   45929       0    9263       0       0       0       0       0       0       0    1505   44424 
 71/ 71:   47077       0    8315       0       0       0       0       0       0       0    1543   45534 
 72/ 72:   48254       0    7174       0       0       0       0       0       0       0    1581   46673 
 73/ 73:   49460       0    5811       0       0       0       0       0       0       0    1621   47840 
 74/ 74:   50697       0    4190       0       0       0       0       0       0       0    1661   49036 
 75/ 75:   51964       0    2269       0       0       0       0       0       0       0    1703   50261 
retiree1/retiree2
    age     fIRA    tIRA  RMDref   fRoth   tRoth fAftaTx tAftaTx   o_inc      SS Expense TFedTax Spndble`,
			sip: sipJoint,
			xp:  xpJoint,
		},
	}
	for i, elem := range tests {
		//fmt.Printf("================ CASE %d ==================\n", i)
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
		donothing := false
		oldout, w, err := ms.RedirectModelSpecsTable(mychan, donothing)
		if err != nil {
			t.Errorf("RedirectModelSpecsTable: %s\n", err)
			return // should this be continue?
		}

		//xp := &[]float64{0.0, 0.0}
		ms.printActivitySummary(elem.xp)

		str := ms.RestoreModelSpecsTable(mychan, oldout, w, donothing)
		strnn := strings.TrimSpace(str)
		expect := elem.expect
		//re := regexp.MustCompile("\r")
		//strnn := re.ReplaceAllString(strn, "")
		//expect := re.ReplaceAllString(elem.expect, "")
		//r:=NewReplacer(U+0010, '')
		//strnn := r.Replace(strn)
		//expect := r.Replace(elem.expect)
		//strn := stripWhitespace(str)
		//warningmes := stripWhitespace(elem.warningmes)
		if expect != strnn {
			showStrMismatch(expect, strnn)
			t.Errorf("TestActivitySummary case %d:  expected output:\n\t '%s'\n\tbut found:\n\t'%s'\n", i, elem.expect, strnn)
		}
	}
}

//func (ms ModelSpecs) printIncomeHeader(headerkeylist []string, countlist []int, incomeCat []string, fieldwidth int)
func TestPrintIncomeHeader(t *testing.T) {
	tests := []struct {
		sip        map[string]string
		expect     string
		headerlist []string
		countlist  []int
		csvfile    *os.File
		tablefile  *os.File
	}{
		{ //case 0
			sip:       sipJoint,
			countlist: []int{0, 0, 0, 0},
			expect: `retiree1
    age`,
			csvfile:   (*os.File)(nil),
			tablefile: os.Stdout,
		},
		{ //case 1
			sip:       sipSingle,
			countlist: []int{0, 0, 0, 0},
			expect: `retir
 age`,
			csvfile:   (*os.File)(nil),
			tablefile: os.Stdout,
		},
		{ //case 2
			sip:        sipJoint,
			countlist:  []int{1, 0, 0, 0},
			headerlist: []string{"nokey"},
			expect: `retiree1 SSincome
    age        SS`,
			csvfile:   (*os.File)(nil),
			tablefile: os.Stdout,
		},
		{ //case 3
			sip:       sipJoint,
			countlist: []int{3, 3, 3, 3},
			headerlist: []string{
				"SS1", "SS2", "SS3", "income1", "income2",
				"income3", "asset1", "asset2", "asset3",
				"expense1", "expense2", "expense3",
			},
			expect: `retiree1 SSincome:                  Income:                    AssetSale:                 Expense:                  
    age       SS1      SS2      SS3  income1  income2  income3   asset1   asset2   asset3 expense1 expense2 expense3`,
			csvfile:   (*os.File)(nil),
			tablefile: os.Stdout,
		},
	}
	for i, elem := range tests {
		//fmt.Printf("=============== Case %d =================\n", i)
		ip := NewInputParams(elem.sip)
		ms := ModelSpecs{
			ip:        ip,
			logfile:   os.Stdout,
			errfile:   os.Stderr,
			ao:        NewAppOutput(elem.csvfile, elem.tablefile),
			assetSale: make([][]float64, 0),
		}
		ms.assetSale = append(ms.assetSale, make([]float64, ip.numyr))

		mychan := make(chan string)
		DoNothing := false //true
		oldout, w, err := ms.RedirectModelSpecsTable(mychan, DoNothing)
		if err != nil {
			t.Errorf("RedirectModelSpecsTable: %s\n", err)
			return // should this be continue?
		}

		incomeCat := []string{"SSincome:", "Income:", "AssetSale:", "Expense:"}
		fieldwidth := 8
		ms.printIncomeHeader(elem.headerlist, elem.countlist, incomeCat, fieldwidth)

		str := ms.RestoreModelSpecsTable(mychan, oldout, w, DoNothing)
		strn := strings.TrimSpace(str)
		if elem.expect != strn {
			//showStrMismatch(elem.expect, strn)
			t.Errorf("TestPrintIncomeHeader case %d:  expected output:\n\t '%s'\n\tbut found:\n\t'%s'\n", i, elem.expect, strn)
		}
	}
}

//func (ms ModelSpecs) getSSIncomeAssetExpenseList() ([]string, []int, [][]float64)
func TestGetIncomeAssetExpenseList(t *testing.T) {
	tests := []struct {
		sip            map[string]string
		incomeStreams  int
		expenseStreams int
		SSStreams      int
		AssetStreams   int
	}{
		{ //case 0
			incomeStreams:  3,
			expenseStreams: 3,
			SSStreams:      3,
			AssetStreams:   3,
			sip:            sipJoint,
		},
		{ //case 1
			incomeStreams:  1,
			expenseStreams: 0,
			SSStreams:      2,
			AssetStreams:   4,
			sip:            sipJoint,
		},
	}
	for i, elem := range tests {
		//fmt.Printf("=============== Case %d =================\n", i)
		ip := NewInputParams(elem.sip)
		ms := ModelSpecs{
			//ip:      ip,
			//logfile: os.Stdout,
			//errfile: os.Stderr,
			//ao:      NewAppOutput(elem.csvfile, elem.tablefile),

			SS:    make([][]float64, 0),
			SStag: make([]string, 0),

			income:    make([][]float64, 0),
			incometag: make([]string, 0),

			assetSale: make([][]float64, 0),
			assettag:  make([]string, 0),

			expenses:   make([][]float64, 0),
			expensetag: make([]string, 0),
		}
		for i := 0; i <= elem.SSStreams; i++ {
			ms.SS = append(ms.SS, make([]float64, ip.numyr))
			str := fmt.Sprintf("SS%d", i)
			ms.SStag = append(ms.SStag, str)
		}
		for i := 0; i <= elem.incomeStreams; i++ {
			ms.income = append(ms.income, make([]float64, ip.numyr))
			str := fmt.Sprintf("income%d", i)
			ms.incometag = append(ms.incometag, str)
		}
		for i := 0; i <= elem.AssetStreams; i++ {
			ms.assetSale = append(ms.assetSale, make([]float64, ip.numyr))
			str := fmt.Sprintf("asset%d", i)
			ms.assettag = append(ms.assettag, str)
		}
		for i := 0; i <= elem.expenseStreams; i++ {
			ms.expenses = append(ms.expenses, make([]float64, ip.numyr))
			str := fmt.Sprintf("expense%d", i)
			ms.expensetag = append(ms.expensetag, str)
		}
		ms.assetSale[1][7] = 50000

		headerlist, countlist, matrix := ms.getSSIncomeAssetExpenseList()
		//fmt.Printf("headerlist: %#v\n", headerlist)
		//fmt.Printf("countlist: %#v\n", countlist)
		//fmt.Printf("matrix: %#v\n", matrix)
		htot := elem.SSStreams + elem.incomeStreams + elem.AssetStreams + elem.expenseStreams
		if htot != len(headerlist) {
			t.Errorf("TestGetIncomeAssetExpenseList case %d: expected %d headers but found %d\n", i, htot, len(headerlist))
		}
		if htot != len(matrix) {
			t.Errorf("TestGetIncomeAssetExpenseList case %d: expected %d vectors but found %d\n", i, htot, len(matrix))
		}
		if elem.SSStreams != countlist[0] {
			t.Errorf("TestGetIncomeAssetExpenseList case %d:  expected %d SS streams but found %d streams\n", i, elem.SSStreams, countlist[0])
		}
		if elem.incomeStreams != countlist[1] {
			t.Errorf("TestGetIncomeAssetExpenseList case %d:  expected %d income streams but found %d streams\n", i, elem.SSStreams, countlist[0])
		}
		if elem.AssetStreams != countlist[2] {
			t.Errorf("TestGetIncomeAssetExpenseList case %d:  expected %d asset streams but found %d streams\n", i, elem.SSStreams, countlist[0])
		}
		if elem.expenseStreams != countlist[3] {
			t.Errorf("TestGetIncomeAssetExpenseList case %d:  expected %d expense streams but found %d streams\n", i, elem.SSStreams, countlist[0])
		}
	}
}

//func (ms ModelSpecs) printIncomeExpenseDetails()
func TestPrintIncomeExpenseDetails(t *testing.T) {
	tests := []struct {
		sip            map[string]string
		incomeStreams  int
		expenseStreams int
		SSStreams      int
		AssetStreams   int
		expect         string
	}{
		{ //case 0
			incomeStreams:  3,
			expenseStreams: 3,
			SSStreams:      3,
			AssetStreams:   3,
			sip:            sipJoint,
			expect:         ``,
		},
		{ //case 1
			incomeStreams:  1,
			expenseStreams: 0,
			SSStreams:      2,
			AssetStreams:   4,
			sip:            sipJoint,
			expect:         ``,
		},
	}
	for i, elem := range tests {
		fmt.Printf("=============== Case %d =================\n", i)
		ip := NewInputParams(elem.sip)
		csvfile := (*os.File)(nil)
		tablefile := os.Stdout
		ms := ModelSpecs{
			//ip:      ip,
			logfile: os.Stdout,
			errfile: os.Stderr,
			ao:      NewAppOutput(csvfile, tablefile),

			SS:    make([][]float64, 0),
			SStag: make([]string, 0),

			income:    make([][]float64, 0),
			incometag: make([]string, 0),

			assetSale: make([][]float64, 0),
			assettag:  make([]string, 0),

			expenses:   make([][]float64, 0),
			expensetag: make([]string, 0),
		}
		for i := 0; i <= elem.SSStreams; i++ {
			ms.SS = append(ms.SS, make([]float64, ip.numyr))
			str := fmt.Sprintf("SS%d", i)
			ms.SStag = append(ms.SStag, str)
		}
		for i := 0; i <= elem.incomeStreams; i++ {
			ms.income = append(ms.income, make([]float64, ip.numyr))
			str := fmt.Sprintf("income%d", i)
			ms.incometag = append(ms.incometag, str)
		}
		for i := 0; i <= elem.AssetStreams; i++ {
			ms.assetSale = append(ms.assetSale, make([]float64, ip.numyr))
			str := fmt.Sprintf("asset%d", i)
			ms.assettag = append(ms.assettag, str)
		}
		for i := 0; i <= elem.expenseStreams; i++ {
			ms.expenses = append(ms.expenses, make([]float64, ip.numyr))
			str := fmt.Sprintf("expense%d", i)
			ms.expensetag = append(ms.expensetag, str)
		}
		ms.assetSale[1][7] = 50000

		//headerlist, countlist, matrix := ms.getSSIncomeAssetExpenseList()
		//fmt.Printf("headerlist: %#v\n", headerlist)
		//fmt.Printf("countlist: %#v\n", countlist)
		//fmt.Printf("matrix: %#v\n", matrix)

		mychan := make(chan string)
		DoNothing := true //false //true
		oldout, w, err := ms.RedirectModelSpecsTable(mychan, DoNothing)
		if err != nil {
			t.Errorf("RedirectModelSpecsTable: %s\n", err)
			return // should this be continue?
		}

		ms.printIncomeExpenseDetails()

		str := ms.RestoreModelSpecsTable(mychan, oldout, w, DoNothing)
		strn := strings.TrimSpace(str)
		if elem.expect != strn && i == 10 {
			//showStrMismatch(elem.expect, strn)
			t.Errorf("TestPrintIncomeHeader case %d:  expected output:\n\t '%s'\n\tbut found:\n\t'%s'\n", i, elem.expect, strn)
		}
	}
}

func showStrMismatch(s1, s2 string) { // TODO move to Utility functions
	for i := 0; i < len(s1); i++ {
		if s1[i] != s2[i] {
			fmt.Printf("Char#: %d, CharVals1: %c, CharInts1: %d, CharVals2: %c, CharInts2: %d\n", i, s1[i], s1[i], s2[i], s2[i])
			fmt.Printf("expect: '%s'\n", s1[:i])
			fmt.Printf(" strnn: '%s'\n", s2[:i])
			break
		}
	}
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
	tests := []struct {
		year     int
		index    int
		expected float64
	}{
		{ // CASE 0
			year:     5,
			index:    0,
			expected: 0.0,
		},
		{ // CASE 1
			year:     5,
			index:    1,
			expected: 20000.0,
		},
	}
	for i, elem := range tests {
		ip := NewInputParams(sipSingle)
		ti := NewTaxInfo(ip.filingStatus)
		taxbins := len(*ti.Taxtable)
		cgbins := len(*ti.Capgainstable)
		vindx, err := NewVectorVarIndex(ip.numyr, taxbins, cgbins, ip.accmap, os.Stdout)
		if err != nil {
			t.Errorf("PrintModelRow case %d: %s", i, err)
			continue
		}
		ms := ModelSpecs{
			ip:      ip,
			vindx:   vindx,
			ti:      ti,
			logfile: os.Stdout,
			errfile: os.Stderr,
			accounttable: []account{
				{
					acctype: "IRA",
				},
				{
					acctype: "aftertax",
				},
			},
			assetSale: make([][]float64, 0),
		}
		ms.assetSale = append(ms.assetSale, make([]float64, ip.numyr))
		ms.assetSale[0][5] = 20000

		damount := ms.depositAmount(xpSingle, elem.year, elem.index)
		if damount != elem.expected {
			t.Errorf("PrintdepositAmount case %d: expected: %f found %f\n", i, elem.expected, damount)
		}
	}
}

//TODO complete TestOrdinaryTaxable after more print function are working
//func (ms ModelSpecs) ordinaryTaxable(year int, xp *[]float64) float64
func TestOrdinaryTaxable(t *testing.T) {
	tests := []struct {
		sip    map[string]string
		sxp    *[]float64
		year   int
		expect float64
	}{
		{
			sip:    sipSingle,
			sxp:    xpSingle,
			year:   7,
			expect: 0.0,
		},
	}
	for i, elem := range tests {
		fmt.Printf("======== CASE %d ========\n", i)
		ip := NewInputParams(elem.sip)
		//fmt.Printf("InputParams: %#v\n", ip)
		ti := NewTaxInfo(ip.filingStatus)
		taxbins := len(*ti.Taxtable)
		cgbins := len(*ti.Capgainstable)
		vindx, err := NewVectorVarIndex(ip.numyr, taxbins,
			cgbins, ip.accmap, os.Stdout)
		if err != nil {
			t.Errorf("TestOrdinaryTaxable case %d: %s", i, err)
			continue
		}
		logfile := os.Stdout
		csvfile := (*os.File)(nil)
		ms := NewModelSpecs(vindx, ti, ip, false,
			false, os.Stderr, logfile, csvfile, logfile)
		ot := ms.ordinaryTaxable(elem.year, elem.sxp)
		if ot != elem.expect {
			t.Errorf("TestOrdinaryTaxable case %d: expected %f, found %f\n", i, elem.expect, ot)
		}
	}
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
	tests := []struct {
		sip    map[string]string
		sxp    *[]float64
		expect string
	}{
		{ // case 0
			sip: sipSingle,
			sxp: xpSingle,
			expect: `======
Optimized for Spending with single status
	starting at age 65 with an estate of $379_660 liquid and $0 illiquid

No desired minium or maximum amount specified

After tax yearly income: $37_164 adjusting for inflation
	and final estate at age 76 with $0 liquid and $0 illiquid

total withdrawals: $506_759
total ordinary taxable income $336_414
total ordinary tax on all taxable income: $42_825 (12.7%) of taxable income
total income (withdrawals + other) $506_759
total cap gains tax: $0
total all tax on all income: $42_825 (8.5%)
Total spendable (after tax money): $463_934`,
		},
		{ // case 1
			sip: sipJoint,
			sxp: xpJoint,
			expect: `======
Optimized for Spending with joint status
	starting at age 65 with an estate of $379_660 liquid and $0 illiquid

No desired minium or maximum amount specified

After tax yearly income: $39_264 adjusting for inflation
	and final estate at age 76 with $0 liquid and $0 illiquid

total withdrawals: $506_759
total ordinary taxable income $166_068
total ordinary tax on all taxable income: $16_607 (10.0%) of taxable income
total income (withdrawals + other) $506_759
total cap gains tax: $0
total all tax on all income: $16_607 (3.3%)
Total spendable (after tax money): $490_153`,
		},
	}
	for i, elem := range tests {
		//fmt.Printf("======== CASE %d ========\n", i)
		ip := NewInputParams(elem.sip)
		//fmt.Printf("InputParams: %#v\n", ip)
		ti := NewTaxInfo(ip.filingStatus)
		taxbins := len(*ti.Taxtable)
		cgbins := len(*ti.Capgainstable)
		vindx, err := NewVectorVarIndex(ip.numyr, taxbins,
			cgbins, ip.accmap, os.Stdout)
		if err != nil {
			t.Errorf("TestOrdinaryTaxable case %d: %s", i, err)
			continue
		}
		logfile := os.Stdout
		csvfile := (*os.File)(nil)
		ms := NewModelSpecs(vindx, ti, ip, false,
			false, os.Stderr, logfile, csvfile, logfile)

		mychan := make(chan string)
		DoNothing := false //true
		oldout, w, err := ms.RedirectModelSpecsTable(mychan, DoNothing)
		if err != nil {
			t.Errorf("RedirectModelSpecsTable: %s\n", err)
			return // should this be continue?
		}

		ms.printBaseConfig(elem.sxp)

		str := ms.RestoreModelSpecsTable(mychan, oldout, w, DoNothing)
		strn := strings.TrimSpace(str)
		//strn := stripWhitespace(str)
		//ot := ms.ordinaryTaxable(elem.year, elem.sxp)
		if strn != elem.expect {
			//showStrMismatch(elem.expect, strn)
			t.Errorf("TestPrintBaseConfig case %d: expected\n'%s'\nfound '%s'\n", i, elem.expect, strn)
		}
	}
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
		{ // Case 1 // case to match mobile.toml
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
		//fmt.Printf("======== CASE %d ========\n", i)
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
		//createDefX(&res.X)
	}
}

func createDefX(xp *[]float64) {
	fmt.Printf("var xp = *[]float64{\n")
	count := 0
	for _, v := range *xp {
		count++
		fmt.Printf("%v, ", v)
		if count > 4 {
			count = 0
			fmt.Printf("\n")
		}

	}
	fmt.Printf("}\n")
}

var xpJoint = &[]float64{ // using sipJoint InputParameters
	13303.039872341637, 0, 0, 0, 0,
	0, 0, 13635.615869150632, 0, 0,
	0, 0, 0, 0, 13976.506265879383,
	0, 0, 0, 0, 0,
	0, 14325.918922526458, 0, 0, 0,
	0, 0, 0, 14684.066895589624, 0,
	0, 0, 0, 0, 0,
	15051.168567979354, 0, 0, 0, 0,
	0, 0, 15427.447782178877, 0, 0,
	0, 0, 0, 0, 15813.133976733143,
	0, 0, 0, 0, 0,
	0, 16208.462326151364, 0, 0, 0,
	0, 0, 0, 16613.673884305335, 0,
	0, 0, 0, 0, 0,
	17029.015731413012, 0, 0, 0, 0,
	0, 0, 40594.442354607934, 41609.30341347362, 42649.53599881044,
	43715.7743987808, 44808.66875875031, 45928.885477719065, 47077.107614662076, 48254.035305028425,
	49460.38618765402, 50696.895842345555, 51964.318238404245, 379659.7116670852, 359409.1854712258,
	336867.874981217, 311871.4393213511, 284245.0048179244, 253802.51622272478, 220346.0485897059,
	183665.07743354648, 143535.70465622915, 99719.83757668926, 51964.318238404245, 0,
	39264.13836737424, 40245.74182655855, 41251.8853722225, 42283.182506528145, 43340.262069191354,
	44423.768620921124, 45534.36283644418, 46672.72190735511, 47839.53995503888, 49035.52845391502,
	50261.416665262936, 0, 0, 0, 0,
	0, 0, 0, 0, 0,
	0, 0,
}

var xpSingle = &[]float64{ // using sipSingle input parameters
	12235.208083996806, 14713.53302947793, 0, 0, 0,
	0, 0, 12541.088286096725, 15081.371355215406, 0,
	0, 0, 0, 0, 12854.615493249144,
	15458.405639095778, 0, 0, 0, 0,
	0, 13175.980880580368, 15844.865780073274, 0, 0,
	0, 0, 0, 13505.380402594878, 16240.987424575116,
	0, 0, 0, 0, 0,
	13843.01491265975, 16647.012110189484, 0, 0, 0,
	0, 0, 14189.090285476244, 17063.18741294425, 0,
	0, 0, 0, 0, 14543.817542613147,
	17489.767098267657, 0, 0, 0, 0,
	0, 14907.412981178479, 17927.011275724217, 0, 0,
	0, 0, 0, 15280.09830570794, 18375.18655761753,
	0, 0, 0, 0, 0,
	15662.100763350638, 18834.566221558, 0, 0, 0,
	0, 0, 40594.44235460792, 41609.30341347362, 42649.535998810454,
	43715.77439878082, 44808.66875875035, 45928.885477719086, 47077.1076146621, 48254.035305028454,
	49460.386187654025, 50696.895842345584, 51964.31823840425, 379659.71166708524, 359409.18547122594,
	336867.8749812171, 311871.4393213513, 284245.00481792446, 253802.51622272475, 220346.04858970596,
	183665.0774335465, 143535.70465622915, 99719.8375766893, 51964.31823840426, 0,
	37163.891591787, 38092.98888158163, 39045.31360362118, 40021.446443711786, 41021.98260480459,
	42047.5321699247, 43098.720474172835, 44176.188486026986, 45280.59319817755, 46412.608028132156,
	47572.92322883549, 0, 0, 0, 0,
	0, 0, 0, 0, 0,
	0, 0,
}
