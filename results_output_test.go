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
var sipSingle3Acc = map[string]string{
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
	"eT_Roth1":                   "10", //10K
	"eT_Roth2":                   "",
	"eT_Roth_Rate1":              "",
	"eT_Roth_Rate2":              "",
	"eT_Roth_Contrib1":           "",
	"eT_Roth_Contrib2":           "",
	"eT_Roth_ContribStartAge1":   "",
	"eT_Roth_ContribStartAge2":   "",
	"eT_Roth_ContribEndAge1":     "",
	"eT_Roth_ContribEndAge2":     "",
	"eT_Aftatax":                 "50", //50k
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
		onek           float64
	}{
		{ //case 0
			incomeStreams:  3,
			expenseStreams: 3,
			SSStreams:      3,
			AssetStreams:   3,
			sip:            sipJoint,
			expect: `Income and Expense Summary:

retiree1 SSincome:                  Income:                    AssetSale:                 Expense:                  
    age       SS1      SS2      SS3  income1  income2  income3   asset1   asset2   asset3 expense1 expense2 expense3
 65/ 65:        1        2        3        1        2        3        1        2        3        1        2        3
 66/ 66:     1000     1000     1000     1000     1000     1000     1000     1000     1000     1000     1000     1000
 67/ 67:     2000     2000     2000     2000     2000     2000     2000     2000     2000     2000     2000     2000
 68/ 68:     3000     3000     3000     3000     3000     3000     3000     3000     3000     3000     3000     3000
 69/ 69:     4000     4000     4000     4000     4000     4000     4000     4000     4000     4000     4000     4000
 70/ 70:     5000     5000     5000     5000     5000     5000     5000     5000     5000     5000     5000     5000
 71/ 71:     6000     6000     6000     6000     6000     6000     6000     6000     6000     6000     6000     6000
 72/ 72:     7000     7000     7000     7000     7000     7000    50000     7000     7000     7000     7000     7000
 73/ 73:     8000     8000     8000     8000     8000     8000     8000     8000     8000     8000     8000     8000
 74/ 74:     9000     9000     9000     9000     9000     9000     9000     9000     9000     9000     9000     9000
 75/ 75:    10000    10000    10000    10000    10000    10000    10000    10000    10000    10000    10000    10000
retiree1 SSincome:                  Income:                    AssetSale:                 Expense:                  
    age       SS1      SS2      SS3  income1  income2  income3   asset1   asset2   asset3 expense1 expense2 expense3`,
			onek: 1,
		},
		{ //case 1
			incomeStreams:  1,
			expenseStreams: 0,
			SSStreams:      2,
			AssetStreams:   4,
			sip:            sipSingle,
			expect: `Income and Expense Summary:

retir SSincome:         Income:  AssetSale:                         
 age       SS1      SS2  income1   asset1   asset2   asset3   asset4
  65:        0        0        0        0        0        0        0
  66:        1        1        1        1        1        1        1
  67:        2        2        2        2        2        2        2
  68:        3        3        3        3        3        3        3
  69:        4        4        4        4        4        4        4
  70:        5        5        5        5        5        5        5
  71:        6        6        6        6        6        6        6
  72:        7        7        7       50        7        7        7
  73:        8        8        8        8        8        8        8
  74:        9        9        9        9        9        9        9
  75:       10       10       10       10       10       10       10
retir SSincome:         Income:  AssetSale:                         
 age       SS1      SS2  income1   asset1   asset2   asset3   asset4`,
			onek: 1000,
		},
	}
	for i, elem := range tests {
		//fmt.Printf("=============== Case %d =================\n", i)
		ip := NewInputParams(elem.sip)
		csvfile := (*os.File)(nil)
		tablefile := os.Stdout
		ms := ModelSpecs{
			ip:      ip,
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

			OneK: elem.onek,
		}
		for i := 0; i <= elem.SSStreams; i++ {
			v := make([]float64, ip.numyr)
			for j := 1; j < ip.numyr; j++ {
				v[j] = float64(j * 1000)
			}
			v[0] = float64(i)
			ms.SS = append(ms.SS, v)
			str := fmt.Sprintf("SS%d", i)
			ms.SStag = append(ms.SStag, str)
		}
		for i := 0; i <= elem.incomeStreams; i++ {
			v := make([]float64, ip.numyr)
			for j := 1; j < ip.numyr; j++ {
				v[j] = float64(j * 1000)
			}
			v[0] = float64(i)
			ms.income = append(ms.income, v)
			str := fmt.Sprintf("income%d", i)
			ms.incometag = append(ms.incometag, str)
		}
		for i := 0; i <= elem.AssetStreams; i++ {
			v := make([]float64, ip.numyr)
			for j := 1; j < ip.numyr; j++ {
				v[j] = float64(j * 1000)
			}
			v[0] = float64(i)
			ms.assetSale = append(ms.assetSale, v)
			str := fmt.Sprintf("asset%d", i)
			ms.assettag = append(ms.assettag, str)
		}
		for i := 0; i <= elem.expenseStreams; i++ {
			v := make([]float64, ip.numyr)
			for j := 1; j < ip.numyr; j++ {
				v[j] = float64(j * 1000)
			}
			v[0] = float64(i)
			ms.expenses = append(ms.expenses, v)
			str := fmt.Sprintf("expense%d", i)
			ms.expensetag = append(ms.expensetag, str)
		}
		ms.assetSale[1][7] = 50000.0

		//headerlist, countlist, matrix := ms.getSSIncomeAssetExpenseList()
		//fmt.Printf("headerlist: %#v\n", headerlist)
		//fmt.Printf("countlist: %#v\n", countlist)
		//fmt.Printf("matrix: %#v\n", matrix)

		mychan := make(chan string)
		DoNothing := false //true
		oldout, w, err := ms.RedirectModelSpecsTable(mychan, DoNothing)
		if err != nil {
			t.Errorf("RedirectModelSpecsTable: %s\n", err)
			return // should this be continue?
		}

		ms.printIncomeExpenseDetails()

		str := ms.RestoreModelSpecsTable(mychan, oldout, w, DoNothing)
		strn := strings.TrimSpace(str)
		if elem.expect != strn {
			//showStrMismatch(elem.expect, strn)
			t.Errorf("TestTestPrintIncomeExpenseDetails case %d:  expected output:\n\t '%s'\n\tbut found:\n\t'%s'\n", i, elem.expect, strn)
		}
	}
}

//func printAccHeader()
func TestPrintAccHeader(t *testing.T) {
	tests := []struct {
		sip    map[string]string
		expect string
		onek   float64
	}{
		{ //case 0
			sip: sipJoint,
			expect: `retiree1/retiree2
    age      IRA    fIRA    tIRA  RMDref`,
			onek: 1,
		},
		{ //case 1
			sip: sipSingle,
			expect: `retiree1
 age      IRA    fIRA    tIRA  RMDref`,
			onek: 1,
		},
		{ //case 2
			sip: sipSingle3Acc,
			expect: `retiree1
 age      IRA    fIRA    tIRA  RMDref    Roth   fRoth   tRoth  AftaTx fAftaTx tAftaTx`,
			onek: 1,
		},
	}
	for i, elem := range tests {
		//fmt.Printf("=============== Case %d =================\n", i)
		ip := NewInputParams(elem.sip)
		csvfile := (*os.File)(nil)
		tablefile := os.Stdout
		ms := ModelSpecs{
			ip:      ip,
			logfile: os.Stdout,
			errfile: os.Stderr,
			ao:      NewAppOutput(csvfile, tablefile),
			OneK:    elem.onek,
		}

		mychan := make(chan string)
		DoNothing := false //true
		oldout, w, err := ms.RedirectModelSpecsTable(mychan, DoNothing)
		if err != nil {
			t.Errorf("RedirectModelSpecsTable: %s\n", err)
			return // should this be continue?
		}

		ms.printAccHeader()

		str := ms.RestoreModelSpecsTable(mychan, oldout, w, DoNothing)
		strn := strings.TrimSpace(str)
		if elem.expect != strn {
			showStrMismatch(elem.expect, strn)
			t.Errorf("TestTestPrintIncomeExpenseDetails case %d:  expected output:\n\t '%s'\n\tbut found:\n\t'%s'\n", i, elem.expect, strn)
		}
	}
}

//func (ms ModelSpecs) printAccountTrans(xp *[]float64)
func TestPrintAccountTrans(t *testing.T) {
	tests := []struct {
		sip    map[string]string
		sxp    *[]float64
		expect string
	}{
		{ //case 0
			sip: sipJoint,
			sxp: xpJoint,
			expect: `Account Transactions Summary:

retiree1/retiree2
    age      IRA    fIRA    tIRA  RMDref
 54/ 54:  200000       0       0       0
Plan Start: ---------
 65/ 65:  379660   40594       0       0
 66/ 66:  359409   41609       0       0
 67/ 67:  336868   42650       0       0
 68/ 68:  311871   43716       0       0
 69/ 69:  284245   44809       0       0
 70/ 70:  253803   45929       0    9263
 71/ 71:  220346   47077       0    8315
 72/ 72:  183665   48254       0    7174
 73/ 73:  143536   49460       0    5811
 74/ 74:   99720   50697       0    4190
 75/ 75:   51964   51964       0    2269
Plan End: -----------
 76/ 76:       0       0       0       0
retiree1/retiree2
    age      IRA    fIRA    tIRA  RMDref`,
		},
		{ //case 1
			sip: sipSingle,
			sxp: xpSingle,
			expect: `Account Transactions Summary:

retiree1
 age      IRA    fIRA    tIRA  RMDref
  54:  200000       0       0       0
Plan Start: ---------
  65:  379660   40594       0       0
  66:  359409   41609       0       0
  67:  336868   42650       0       0
  68:  311871   43716       0       0
  69:  284245   44809       0       0
  70:  253803   45929       0    9263
  71:  220346   47077       0    8315
  72:  183665   48254       0    7174
  73:  143536   49460       0    5811
  74:   99720   50697       0    4190
  75:   51964   51964       0    2269
Plan End: -----------
  76:       0       0       0       0
retiree1
 age      IRA    fIRA    tIRA  RMDref`,
		},
		{ //case 2
			sip: sipSingle3Acc,
			sxp: xpSingle3Acc,
			expect: `Account Transactions Summary:

retiree1
 age      IRA    fIRA    tIRA  RMDref    Roth   fRoth   tRoth  AftaTx fAftaTx tAftaTx
  54:  200000       0       0       0   10000       0       0   50000       0       0
Plan Start: ---------
  65:  379660   54922       0       0   18983       0       0   94915       0       0
  66:  344222   56295       0       0   20122       0       0  100610       0       0
  67:  305203   57702       0       0   21329       0       0  106646      -0       0
  68:  262350   59145       0       0   22609       0       0  113045      -0       0
  69:  215398   60623       0       0   23966       0       0  119828       0      -0
  70:  164061   31531       0    5988   25404    1791       0  127018   24225       0
  71:  140481   30014       0    5301   25029       0       0  108960   28627       0
  72:  117095   30764       0    4574   26531       0       0   85153   29343       0
  73:   91511   31533       0    3705   28123       0       0   59159   30076       0
  74:   63576   32322       0    2671   29810       0       0   30828   30828       0
  75:   33130   33130       0    1447   31599   31599       0       0       0       0
Plan End: -----------
  76:       0       0       0       0       0       0       0       0       0       0
retiree1
 age      IRA    fIRA    tIRA  RMDref    Roth   fRoth   tRoth  AftaTx fAftaTx tAftaTx`,
		},
	}
	for i, elem := range tests {
		//fmt.Printf("=============== Case %d =================\n", i)
		ip := NewInputParams(elem.sip)
		ti := NewTaxInfo(ip.filingStatus)
		taxbins := len(*ti.Taxtable)
		cgbins := len(*ti.Capgainstable)
		vindx, err := NewVectorVarIndex(ip.numyr, taxbins,
			cgbins, ip.accmap, os.Stdout)
		if err != nil {
			t.Errorf("TestPrintAccountTrans case %d: %s", i, err)
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

		ms.printAccountTrans(elem.sxp)

		str := ms.RestoreModelSpecsTable(mychan, oldout, w, DoNothing)
		strn := strings.TrimSpace(str)
		if elem.expect != strn {
			showStrMismatch(elem.expect, strn)
			t.Errorf("TestPrintAccountTrans case %d:  expected output:\n\t '%s'\n\tbut found:\n\t'%s'\n", i, elem.expect, strn)
		}
	}
}

//func (ms ModelSpecs) printheaderTax()
func TestPrintHeaderTax(t *testing.T) {
	tests := []struct {
		sip    map[string]string
		expect string
	}{
		{ // Case 0
			sip: sipSingle,
			expect: `retiree1
 age     fIRA    tIRA  TxbleO TxbleSS  deduct   T_inc  earlyP  fedtax  mTaxB% fAftaTx tAftaTx  cgTax%   cgTax TFedTax spndble`,
		},
		{ // Case 1
			sip: sipJoint,
			expect: `retiree1/retiree2
    age     fIRA    tIRA  TxbleO TxbleSS  deduct   T_inc  earlyP  fedtax  mTaxB% fAftaTx tAftaTx  cgTax%   cgTax TFedTax spndble`,
		},
	}
	for i, elem := range tests {
		//fmt.Printf("=============== Case %d =================\n", i)
		ip := NewInputParams(elem.sip)
		csvfile := (*os.File)(nil)
		tablefile := os.Stdout
		ms := ModelSpecs{
			ip:      ip,
			logfile: os.Stdout,
			errfile: os.Stderr,
			ao:      NewAppOutput(csvfile, tablefile),
		}

		mychan := make(chan string)
		DoNothing := false //true
		oldout, w, err := ms.RedirectModelSpecsTable(mychan, DoNothing)
		if err != nil {
			t.Errorf("RedirectModelSpecsTable: %s\n", err)
			return // should this be continue?
		}

		ms.printHeaderTax()

		str := ms.RestoreModelSpecsTable(mychan, oldout, w, DoNothing)
		strn := strings.TrimSpace(str)
		if elem.expect != strn {
			showStrMismatch(elem.expect, strn)
			t.Errorf("TestPrintHeaderTax case %d:  expected output:\n\t '%s'\n\tbut found:\n\t'%s'\n", i, elem.expect, strn)
		}
	}
}

//def print_tax(res):
func TestPrintTax(t *testing.T) {
	tests := []struct {
		sip    map[string]string
		sxp    *[]float64
		expect string
	}{
		{ // Case 0
			sip: sipSingle,
			sxp: xpSingle,
			expect: `Tax Summary:

retiree1
 age     fIRA    tIRA  TxbleO TxbleSS  deduct   T_inc  earlyP  fedtax  mTaxB% fAftaTx tAftaTx  cgTax%   cgTax TFedTax spndble
  65:   40594       0       0       0   13646   26949      0     3431      15       0       0     100       0    3431   37164
  66:   41609       0       0       0   13987   27622      0     3516      15       0       0     100       0    3516   38093
  67:   42650       0       0       0   14337   28313      0     3604      15       0       0     100       0    3604   39045
  68:   43716       0       0       0   14695   29021      0     3694      15       0       0     100       0    3694   40021
  69:   44809       0       0       0   15062   29746      0     3787      15       0       0     100       0    3787   41022
  70:   45929       0       0       0   15439   30490      0     3881      15       0       0     100       0    3881   42048
  71:   47077       0       0       0   15825   31252      0     3978      15       0       0     100       0    3978   43099
  72:   48254       0       0       0   16220   32034      0     4078      15       0       0     100       0    4078   44176
  73:   49460       0       0       0   16626   32834      0     4180      15       0       0     100       0    4180   45281
  74:   50697       0       0       0   17042   33655      0     4284      15       0       0     100       0    4284   46413
  75:   51964       0       0       0   17468   34497      0     4391      15       0       0     100       0    4391   47573
retiree1
 age     fIRA    tIRA  TxbleO TxbleSS  deduct   T_inc  earlyP  fedtax  mTaxB% fAftaTx tAftaTx  cgTax%   cgTax TFedTax spndble`,
		},
		{ // Case 1
			sip: sipJoint,
			sxp: xpJoint,
			expect: `Tax Summary:

retiree1/retiree2
    age     fIRA    tIRA  TxbleO TxbleSS  deduct   T_inc  earlyP  fedtax  mTaxB% fAftaTx tAftaTx  cgTax%   cgTax TFedTax spndble
 65/ 65:   40594       0       0       0   27291   13303      0     1330      10       0       0     100       0    1330   39264
 66/ 66:   41609       0       0       0   27974   13636      0     1364      10       0       0     100       0    1364   40246
 67/ 67:   42650       0       0       0   28673   13977      0     1398      10       0       0     100       0    1398   41252
 68/ 68:   43716       0       0       0   29390   14326      0     1433      10       0       0     100       0    1433   42283
 69/ 69:   44809       0       0       0   30125   14684      0     1468      10       0       0     100       0    1468   43340
 70/ 70:   45929       0       0       0   30878   15051      0     1505      10       0       0     100       0    1505   44424
 71/ 71:   47077       0       0       0   31650   15427      0     1543      10       0       0     100       0    1543   45534
 72/ 72:   48254       0       0       0   32441   15813      0     1581      10       0       0     100       0    1581   46673
 73/ 73:   49460       0       0       0   33252   16208      0     1621      10       0       0     100       0    1621   47840
 74/ 74:   50697       0       0       0   34083   16614      0     1661      10       0       0     100       0    1661   49036
 75/ 75:   51964       0       0       0   34935   17029      0     1703      10       0       0     100       0    1703   50261
retiree1/retiree2
    age     fIRA    tIRA  TxbleO TxbleSS  deduct   T_inc  earlyP  fedtax  mTaxB% fAftaTx tAftaTx  cgTax%   cgTax TFedTax spndble`,
		},
		{ // Case 2
			sip: sipSingle3Acc,
			sxp: xpSingle3Acc,
			expect: `Tax Summary:

retiree1
 age     fIRA    tIRA  TxbleO TxbleSS  deduct   T_inc  earlyP  fedtax  mTaxB% fAftaTx tAftaTx  cgTax%   cgTax TFedTax spndble
  65:   54922       0       0       0   13646   41276      0     5580      15       0       0     100       0    5580   49342
  66:   56295       0       0       0   13987   42308      0     5719      15       0       0     100       0    5719   50576
  67:   57702       0       0       0   14337   43366      0     5862      15      -0       0     100       0    5862   51840
  68:   59145       0       0       0   14695   44450      0     6009      15      -0       0     100       0    6009   53136
  69:   60623       0       0       0   15062   45561      0     6159      15       0      -0     100       0    6159   54465
  70:   31531       0       0       0   15439   16093      0     1722      15   24225       0     100       0    1722   55826
  71:   30014       0       0       0   15825   14189      0     1419      10   28627       0     100       0    1419   57222
  72:   30764       0       0       0   16220   14544      0     1454      10   29343       0     100       0    1454   58652
  73:   31533       0       0       0   16626   14907      0     1491      10   30076       0     100       0    1491   60119
  74:   32322       0       0       0   17042   15280      0     1528      10   30828       0     100       0    1528   61622
  75:   33130       0       0       0   17468   15662      0     1566      10       0       0     100       0    1566   63162
retiree1
 age     fIRA    tIRA  TxbleO TxbleSS  deduct   T_inc  earlyP  fedtax  mTaxB% fAftaTx tAftaTx  cgTax%   cgTax TFedTax spndble`,
		},
	}
	for i, elem := range tests {
		//fmt.Printf("=============== Case %d =================\n", i)
		ip := NewInputParams(elem.sip)
		ti := NewTaxInfo(ip.filingStatus)
		taxbins := len(*ti.Taxtable)
		cgbins := len(*ti.Capgainstable)
		vindx, err := NewVectorVarIndex(ip.numyr, taxbins,
			cgbins, ip.accmap, os.Stdout)
		if err != nil {
			t.Errorf("TestPrintAccountTrans case %d: %s", i, err)
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

		ms.printTax(elem.sxp)

		str := ms.RestoreModelSpecsTable(mychan, oldout, w, DoNothing)
		strn := strings.TrimSpace(str)
		if elem.expect != strn {
			showStrMismatch(elem.expect, strn)
			t.Errorf("TestPrintTax case %d:  expected output:\n\t '%s'\n\tbut found:\n\t'%s'\n", i, elem.expect, strn)
		}
	}
}

//func (ms ModelSpecs) printHeaderTaxBrackets()
func TestPrintHeaderTaxBrackets(t *testing.T) {
	tests := []struct {
		sip    map[string]string
		expect string
	}{
		{ // Case 0
			sip: sipSingle,
			expect: `Marginal Rate(%):     10     15     25     28     33     35     40
retiree1
 age     fIRA    tIRA  TxbleO TxbleSS  deduct   T_inc  fedtax brckt0 brckt1 brckt2 brckt3 brckt4 brckt5 brckt6 brkTot`,
		},
		{ // Case 1
			sip: sipJoint,
			expect: `Marginal Rate(%):     10     15     25     28     33     35     40
retiree1/retiree2
    age     fIRA    tIRA  TxbleO TxbleSS  deduct   T_inc  fedtax brckt0 brckt1 brckt2 brckt3 brckt4 brckt5 brckt6 brkTot`,
		},
	}
	for i, elem := range tests {
		//fmt.Printf("=============== Case %d =================\n", i)
		ip := NewInputParams(elem.sip)
		ti := NewTaxInfo(ip.filingStatus)
		csvfile := (*os.File)(nil)
		tablefile := os.Stdout
		ms := ModelSpecs{
			ip:      ip,
			ti:      ti,
			logfile: os.Stdout,
			errfile: os.Stderr,
			ao:      NewAppOutput(csvfile, tablefile),
		}

		mychan := make(chan string)
		DoNothing := false //true
		oldout, w, err := ms.RedirectModelSpecsTable(mychan, DoNothing)
		if err != nil {
			t.Errorf("RedirectModelSpecsTable: %s\n", err)
			return // should this be continue?
		}

		ms.printHeaderTaxBrackets()

		str := ms.RestoreModelSpecsTable(mychan, oldout, w, DoNothing)
		strn := strings.TrimSpace(str)
		if elem.expect != strn {
			showStrMismatch(elem.expect, strn)
			t.Errorf("TestPrintHeaderTaxBrackests case %d:  expected output:\n\t '%s'\n\tbut found:\n\t'%s'\n", i, elem.expect, strn)
		}
	}
}

//def print_tax_brackets(res):
func TestPrintTaxBrackets(t *testing.T) {
	tests := []struct {
		sip    map[string]string
		sxp    *[]float64
		expect string
	}{
		{ // Case 0
			sip: sipSingle,
			sxp: xpSingle,
			expect: `Overall Tax Bracket Summary:
                                            Marginal Rate(%):     10     15     25     28     33     35     40
retiree1
 age     fIRA    tIRA  TxbleO TxbleSS  deduct   T_inc  fedtax brckt0 brckt1 brckt2 brckt3 brckt4 brckt5 brckt6 brkTot
  65:   40594       0       0       0   13646   26949    3431  12235  14714      0      0      0      0      0  26949
  66:   41609       0       0       0   13987   27622    3516  12541  15081      0      0      0      0      0  27622
  67:   42650       0       0       0   14337   28313    3604  12855  15458      0      0      0      0      0  28313
  68:   43716       0       0       0   14695   29021    3694  13176  15845      0      0      0      0      0  29021
  69:   44809       0       0       0   15062   29746    3787  13505  16241      0      0      0      0      0  29746
  70:   45929       0       0       0   15439   30490    3881  13843  16647      0      0      0      0      0  30490
  71:   47077       0       0       0   15825   31252    3978  14189  17063      0      0      0      0      0  31252
  72:   48254       0       0       0   16220   32034    4078  14544  17490      0      0      0      0      0  32034
  73:   49460       0       0       0   16626   32834    4180  14907  17927      0      0      0      0      0  32834
  74:   50697       0       0       0   17042   33655    4284  15280  18375      0      0      0      0      0  33655
  75:   51964       0       0       0   17468   34497    4391  15662  18835      0      0      0      0      0  34497
                                            Marginal Rate(%):     10     15     25     28     33     35     40
retiree1
 age     fIRA    tIRA  TxbleO TxbleSS  deduct   T_inc  fedtax brckt0 brckt1 brckt2 brckt3 brckt4 brckt5 brckt6 brkTot`,
		},
		{ // Case 1
			sip: sipJoint,
			sxp: xpJoint,
			expect: `Overall Tax Bracket Summary:
                                               Marginal Rate(%):     10     15     25     28     33     35     40
retiree1/retiree2
    age     fIRA    tIRA  TxbleO TxbleSS  deduct   T_inc  fedtax brckt0 brckt1 brckt2 brckt3 brckt4 brckt5 brckt6 brkTot
 65/ 65:   40594       0       0       0   27291   13303    1330  13303      0      0      0      0      0      0  13303
 66/ 66:   41609       0       0       0   27974   13636    1364  13636      0      0      0      0      0      0  13636
 67/ 67:   42650       0       0       0   28673   13977    1398  13977      0      0      0      0      0      0  13977
 68/ 68:   43716       0       0       0   29390   14326    1433  14326      0      0      0      0      0      0  14326
 69/ 69:   44809       0       0       0   30125   14684    1468  14684      0      0      0      0      0      0  14684
 70/ 70:   45929       0       0       0   30878   15051    1505  15051      0      0      0      0      0      0  15051
 71/ 71:   47077       0       0       0   31650   15427    1543  15427      0      0      0      0      0      0  15427
 72/ 72:   48254       0       0       0   32441   15813    1581  15813      0      0      0      0      0      0  15813
 73/ 73:   49460       0       0       0   33252   16208    1621  16208      0      0      0      0      0      0  16208
 74/ 74:   50697       0       0       0   34083   16614    1661  16614      0      0      0      0      0      0  16614
 75/ 75:   51964       0       0       0   34935   17029    1703  17029      0      0      0      0      0      0  17029
                                               Marginal Rate(%):     10     15     25     28     33     35     40
retiree1/retiree2
    age     fIRA    tIRA  TxbleO TxbleSS  deduct   T_inc  fedtax brckt0 brckt1 brckt2 brckt3 brckt4 brckt5 brckt6 brkTot`,
		},
		{ // Case 2
			sip: sipSingle3Acc,
			sxp: xpSingle3Acc,
			expect: `Overall Tax Bracket Summary:
                                            Marginal Rate(%):     10     15     25     28     33     35     40
retiree1
 age     fIRA    tIRA  TxbleO TxbleSS  deduct   T_inc  fedtax brckt0 brckt1 brckt2 brckt3 brckt4 brckt5 brckt6 brkTot
  65:   54922       0       0       0   13646   41276    5580  12235  29041      0      0      0      0      0  41276
  66:   56295       0       0       0   13987   42308    5719  12541  29767      0      0      0      0      0  42308
  67:   57702       0       0       0   14337   43366    5862  12855  30511      0      0      0      0      0  43366
  68:   59145       0       0       0   14695   44450    6009  13176  31274      0      0      0      0      0  44450
  69:   60623       0       0       0   15062   45561    6159  13505  32056      0      0      0      0      0  45561
  70:   31531       0       0       0   15439   16093    1722  13843   2250      0      0      0      0      0  16093
  71:   30014       0       0       0   15825   14189    1419  14189      0      0      0      0      0      0  14189
  72:   30764       0       0       0   16220   14544    1454  14544      0      0      0      0      0      0  14544
  73:   31533       0       0       0   16626   14907    1491  14907      0      0      0      0      0      0  14907
  74:   32322       0       0       0   17042   15280    1528  15280      0      0      0      0      0      0  15280
  75:   33130       0       0       0   17468   15662    1566  15662      0      0      0      0      0      0  15662
                                            Marginal Rate(%):     10     15     25     28     33     35     40
retiree1
 age     fIRA    tIRA  TxbleO TxbleSS  deduct   T_inc  fedtax brckt0 brckt1 brckt2 brckt3 brckt4 brckt5 brckt6 brkTot`,
		},
	}
	for i, elem := range tests {
		//fmt.Printf("=============== Case %d =================\n", i)
		ip := NewInputParams(elem.sip)
		ti := NewTaxInfo(ip.filingStatus)
		taxbins := len(*ti.Taxtable)
		cgbins := len(*ti.Capgainstable)
		vindx, err := NewVectorVarIndex(ip.numyr, taxbins,
			cgbins, ip.accmap, os.Stdout)
		if err != nil {
			t.Errorf("TestPrintTaxBrackets case %d: %s", i, err)
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

		ms.printTaxBrackets(elem.sxp)

		str := ms.RestoreModelSpecsTable(mychan, oldout, w, DoNothing)
		strn := strings.TrimSpace(str)
		if elem.expect != strn {
			showStrMismatch(elem.expect, strn)
			t.Errorf("TestPrintTax case %d:  expected output:\n\t '%s'\n\tbut found:\n\t'%s'\n", i, elem.expect, strn)
		}
	}
}

func TestPrintHeaderCapGainsBrackets(t *testing.T) {
	tests := []struct {
		sip    map[string]string
		expect string
	}{
		{ // Case 0
			sip: sipSingle,
			expect: `Marginal Rate(%):      0     15     20
retiree1
 age  fAftaTx tAftaTx  cgTax% cgTaxbl   T_inc   cgTax brckt0 brckt1 brckt2 brkTot`,
		},
		{ // Case 1
			sip: sipJoint,
			expect: `Marginal Rate(%):      0     15     20
retiree1/retiree2
    age  fAftaTx tAftaTx  cgTax% cgTaxbl   T_inc   cgTax brckt0 brckt1 brckt2 brkTot`,
		},
	}
	for i, elem := range tests {
		//fmt.Printf("=============== Case %d =================\n", i)
		ip := NewInputParams(elem.sip)
		ti := NewTaxInfo(ip.filingStatus)
		csvfile := (*os.File)(nil)
		tablefile := os.Stdout
		ms := ModelSpecs{
			ip:      ip,
			ti:      ti,
			logfile: os.Stdout,
			errfile: os.Stderr,
			ao:      NewAppOutput(csvfile, tablefile),
		}

		mychan := make(chan string)
		DoNothing := false //true
		oldout, w, err := ms.RedirectModelSpecsTable(mychan, DoNothing)
		if err != nil {
			t.Errorf("RedirectModelSpecsTable: %s\n", err)
			return // should this be continue?
		}

		ms.printHeaderCapgainsBrackets()

		str := ms.RestoreModelSpecsTable(mychan, oldout, w, DoNothing)
		strn := strings.TrimSpace(str)
		if elem.expect != strn {
			showStrMismatch(elem.expect, strn)
			t.Errorf("TestPrintHeaderCapgainsBrackests case %d:  expected output:\n\t '%s'\n\tbut found:\n\t'%s'\n", i, elem.expect, strn)
		}
	}
}

//def print_cap_gains_brackets(res):
func TestPrintCapGainsBrackets(t *testing.T) {
	tests := []struct {
		sip    map[string]string
		sxp    *[]float64
		expect string
	}{
		{ // Case 0
			sip:    sipSingle,
			sxp:    xpSingle,
			expect: `Overall Capital Gains Bracket Summary:
                                    Marginal Rate(%):      0     15     20
retiree1
 age  fAftaTx tAftaTx  cgTax% cgTaxbl   T_inc   cgTax brckt0 brckt1 brckt2 brkTot
  65:       0       0     100       0   26949       0      0      0      0      0
  66:       0       0     100       0   27622       0      0      0      0      0
  67:       0       0     100       0   28313       0      0      0      0      0
  68:       0       0     100       0   29021       0      0      0      0      0
  69:       0       0     100       0   29746       0      0      0      0      0
  70:       0       0     100       0   30490       0      0      0      0      0
  71:       0       0     100       0   31252       0      0      0      0      0
  72:       0       0     100       0   32034       0      0      0      0      0
  73:       0       0     100       0   32834       0      0      0      0      0
  74:       0       0     100       0   33655       0      0      0      0      0
  75:       0       0     100       0   34497       0      0      0      0      0
                                    Marginal Rate(%):      0     15     20
retiree1
 age  fAftaTx tAftaTx  cgTax% cgTaxbl   T_inc   cgTax brckt0 brckt1 brckt2 brkTot`,
		},
		{ // Case 1
			sip:    sipJoint,
			sxp:    xpJoint,
			expect: `Overall Capital Gains Bracket Summary:
                                       Marginal Rate(%):      0     15     20
retiree1/retiree2
    age  fAftaTx tAftaTx  cgTax% cgTaxbl   T_inc   cgTax brckt0 brckt1 brckt2 brkTot
 65/ 65:       0       0     100       0   13303       0      0      0      0      0
 66/ 66:       0       0     100       0   13636       0      0      0      0      0
 67/ 67:       0       0     100       0   13977       0      0      0      0      0
 68/ 68:       0       0     100       0   14326       0      0      0      0      0
 69/ 69:       0       0     100       0   14684       0      0      0      0      0
 70/ 70:       0       0     100       0   15051       0      0      0      0      0
 71/ 71:       0       0     100       0   15427       0      0      0      0      0
 72/ 72:       0       0     100       0   15813       0      0      0      0      0
 73/ 73:       0       0     100       0   16208       0      0      0      0      0
 74/ 74:       0       0     100       0   16614       0      0      0      0      0
 75/ 75:       0       0     100       0   17029       0      0      0      0      0
                                       Marginal Rate(%):      0     15     20
retiree1/retiree2
    age  fAftaTx tAftaTx  cgTax% cgTaxbl   T_inc   cgTax brckt0 brckt1 brckt2 brkTot`,
		},
		{ // Case 2
			sip:    sipSingle3Acc,
			sxp:    xpSingle3Acc,
			expect: `Overall Capital Gains Bracket Summary:
                                    Marginal Rate(%):      0     15     20
retiree1
 age  fAftaTx tAftaTx  cgTax% cgTaxbl   T_inc   cgTax brckt0 brckt1 brckt2 brkTot
  65:       0       0     100       0   41276       0      0      0      0      0
  66:       0       0     100       0   42308       0      0      0      0      0
  67:      -0       0     100       0   43366       0      0      0      0      0
  68:      -0       0     100       0   44450       0      0      0      0      0
  69:       0      -0     100       0   45561       0      0      0      0      0
  70:   24225       0     100   24225   16093       0  24225      0      0  24225
  71:   28627       0     100   28627   14189       0  28627      0      0  28627
  72:   29343       0     100   29343   14544       0  29343      0      0  29343
  73:   30076       0     100   30076   14907       0  30076      0      0  30076
  74:   30828       0     100   30828   15280       0  30828      0      0  30828
  75:       0       0     100       0   15662       0      0      0      0      0
                                    Marginal Rate(%):      0     15     20
retiree1
 age  fAftaTx tAftaTx  cgTax% cgTaxbl   T_inc   cgTax brckt0 brckt1 brckt2 brkTot`,
		},
	}
	for i, elem := range tests {
		//fmt.Printf("=============== Case %d =================\n", i)
		ip := NewInputParams(elem.sip)
		ti := NewTaxInfo(ip.filingStatus)
		taxbins := len(*ti.Taxtable)
		cgbins := len(*ti.Capgainstable)
		vindx, err := NewVectorVarIndex(ip.numyr, taxbins,
			cgbins, ip.accmap, os.Stdout)
		if err != nil {
			t.Errorf("TestPrintTaxBrackets case %d: %s", i, err)
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

		ms.printCapGainsBrackets(elem.sxp)

		str := ms.RestoreModelSpecsTable(mychan, oldout, w, DoNothing)
		strn := strings.TrimSpace(str)
		if elem.expect != strn {
			showStrMismatch(elem.expect, strn)
			t.Errorf("TestPrintCapGainsBrackets case %d:  expected output:\n\t '%s'\n\tbut found:\n\t'%s'\n", i, elem.expect, strn)
		}
	}
}

func showStrMismatch(s1, s2 string) { // TODO move to Utility functions
	for i := 0; i < intMin(len(s1), len(s2)); i++ {
		if s1[i] != s2[i] {
			fmt.Printf("Char#: %d, CharVals1: %c, CharInts1: %d, CharVals2: %c, CharInts2: %d\n", i, s1[i], s1[i], s2[i], s2[i])
			fmt.Printf("expect: '%s'\n", s1[:i])
			fmt.Printf(" strnn: '%s'\n", s2[:i])
			break
		}
	}
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
		expect int
	}{
		{
			sip:    sipSingle,
			sxp:    xpSingle,
			year:   7,
			expect: 32033,
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
		ot := ms.ordinaryTaxable(elem.year, elem.sxp)
		if int(ot) != elem.expect {
			t.Errorf("TestOrdinaryTaxable case %d: expected %d, found %d\n", i, elem.expect, int(ot))
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
		fmt.Printf("======== CASE %d ========\n", i)
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
		fmt.Printf("Called LPSimplex() for m:%d x n:%d model\n", len(a), len(a[0]))

		ms.printActivitySummary(&res.X)
		//ms.printIncomeExpenseDetails()
		ms.printAccountTrans(&res.X)
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

var xpSingle3Acc = &[]float64{
	12235.208083996806, 29040.983272280755, 0, 0, 0,
	0, 0, 12541.088286096725, 29767.007854088544, 0,
	0, 0, 0, 0, 12854.615493249143,
	30511.18305044068, 0, 0, 0, 0,
	0, 13175.98088058037, 31273.96262670182, 0, 0,
	0, 0, 0, 13505.380402594881, 32055.811692369327,
	0, 0, 0, 0, 0,
	13843.014912659752, 2249.6221275416715, 0, 0, 0,
	0, 0, 14189.090285476243, 0, 0,
	0, 0, 0, 0, 14543.81754261315,
	0, 0, 0, 0, 0,
	0, 14907.412981178479, 0, 0, 0,
	0, 0, 0, 15280.09830570794, 0,
	0, 0, 0, 0, 0,
	15662.100763350636, 0, 0, 0, 0,
	0, 0, 0, 0, 0,
	0, 0, 0, 0, 0,
	0, 0, 0, 0, 0,
	0, 0, 24225.29974225708, 0, 0,
	28626.84158540119, 0, 0, 29342.512625036125, 0,
	0, 30076.075440661876, 0, 0, 30827.977326678607,
	0, 0, 0, 0, 0,
	54921.89259741077, 0, 0, 56294.93991234677, 0,
	0, 57702.313410155344, 0, -6.322885561649867e-13, 59144.87124540936,
	0, -2.692016170914929e-13, 60623.49302654456, 0, 0,
	31531.495495071264, 1791.147386309299, 24225.29974225708, 30013.920201717847, 0,
	28626.84158540119, 30764.268206760797, 0, 29342.512625036125, 31533.374911929805,
	0, 30076.075440661876, 32321.709284728055, 0, 30827.977326678607,
	33129.75201684625, 31598.676759845548, 0, 379659.71166708524, 18982.98558335426,
	94914.92791677131, 344222.088213855, 20121.96471835553, 100609.82359177759, 305202.7771995983,
	21329.28260145675, 106646.41300728427, 262350.49161680974, 22609.03955754421, 113045.19778772136,
	215397.95759368438, 23965.581930996854, 119827.90965498464, 164060.93244116823, 25403.51684685662,
	127017.58423428373, 140481.2031628629, 25029.1116281802, 108959.82156154828, 117095.31993881382,
	26530.858325871006, 85152.95877471592, 91510.91483597607, 28122.709825423255, 59159.07291866057,
	63576.192319488706, 29810.072414948645, 30827.97732667861, 33129.752016846294, 31598.676759845548,
	0, 0, 0, 0, 49342.22429816952,
	50575.779905623785, 51840.17440326433, 53136.17876334605, 54464.583232429664, 55826.19781324043,
	57221.8527585714, 58652.39907753561, 60118.709054473824, 61621.67678083585, 63162.21870035677,
	0, 0, 3.637978807091714e-12, 0, 0,
	0, 0, 0, 0, 0,
	0, 0, 0, 0, -2.2240121136268788e-13,
	0, 0, 0, 0, 0,
	0, 0, 0, 0, 0,
	0, 0, 0, 0, 0,
	0, 0, 0,
}
