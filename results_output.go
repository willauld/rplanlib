package rplanlib

import (
	"bytes"
	"fmt"
	"math"
	"os"
)

const (
	// used for formating with appoutput.output
	ampv = "&&&&&&&&&&&&&&&&&&&&&&&&&&&&&&&&&&&&&&&&&&&&&&&&&&&&&&&&&&&&&&&&&"
	atv  = "@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@"
)

/*
const __version__ = "0.3-rc2"
*/

/*
def precheck_consistancy():
    print("\nDoing Pre-check:")
    # check that there is income for all contibutions
        #tcontribs = 0
    for year in range(S.numyr):
        t = 0
        for j in range(len(S.accounttable)):
            if S.accounttable[j]['acctype'] != 'aftertax':
                v = S.accounttable[j]
                c = v.get('contributions', None)
                if c is not None:
                    t += c[year]
        if t > S.income[year]:
            print("year: %d, total contributions of (%.0f) to all Retirement accounts exceeds other earned income (%.0f)"%(year, t, S.income[year]))
            print("Please change the contributions in the toml file to be less than non-SS income.")
            exit(1)
    return True
*/

func (ms ModelSpecs) activitySummaryHeader(fieldwidth int) {
	var ageWidth int

	names := ""
	if ms.Ip.MyKey2 != "" && ms.Ip.FilingStatus == Joint {
		names = fmt.Sprintf("%s/%s\n", ms.Ip.MyKey1, ms.Ip.MyKey2)
		ageWidth = 8
	} else {
		if ms.Ip.MyKey1 != "nokey" {
			names = fmt.Sprintf("%s\n", ms.Ip.MyKey1)
		}
		ageWidth = 5
	}
	if names != "" {
		format := fmt.Sprintf("%%%ds", 2*ageWidth)
		str := fmt.Sprintf(format, names)
		ms.Ao.Output(str)
	}
	format := fmt.Sprintf("%%%d.%ds", ageWidth, ageWidth)
	str := fmt.Sprintf(format, "age ")
	ms.Ao.Output(str)
	headers := []string{"fIRA", "tIRA", "RMDref", "fRoth", "tRoth", "fAftaTx", "tAftaTx", "o_inc", "SS", "Expense", "TFedTax", "Spndble"}
	for _, s := range headers {
		format := fmt.Sprintf("&@%%%d.%ds", fieldwidth, fieldwidth)
		str := fmt.Sprintf(format, s)
		ms.Ao.Output(str)
	}
	ms.Ao.Output("\n")
}

func (ms ModelSpecs) PrintActivitySummary(xp *[]float64) {

	ms.Ao.Output("\nActivity Summary:\n")
	ms.Ao.Output("\n")
	fieldwidth := 7
	ms.activitySummaryHeader(fieldwidth)
	for year := 0; year < ms.Ip.Numyr; year++ {
		//T, spendable, tax, rate, cgtax, earlytax, rothearly := ms.IncomeSummary(year, xp)
		_, spendable, tax, _, cgtax, earlytax, _ := ms.IncomeSummary(year, xp)

		rmdref := 0.0
		for j := 0; j < intMin(2, len(ms.Accounttable)); j++ { // at most the first two accounts are type IRA w/ RMD requirement
			if ms.Accounttable[j].acctype == IRA {
				rmd := ms.Ti.rmdNeeded(year, ms.matchRetiree(ms.Accounttable[j].mykey, year, true))
				if rmd > 0 {
					rmdref += (*xp)[ms.Vindx.B(year, j)] / rmd
				}
			}
		}
		withdrawal := map[Acctype]float64{IRA: 0, Roth: 0, Aftertax: 0}
		deposit := map[Acctype]float64{IRA: 0, Roth: 0, Aftertax: 0}
		for j := 0; j < len(ms.Accounttable); j++ {
			withdrawal[ms.Accounttable[j].acctype] += (*xp)[ms.Vindx.W(year, j)]
			deposit[ms.Accounttable[j].acctype] += ms.depositAmount(xp, year, j)
		}

		if ms.Ip.MyKey2 != "" && ms.Ip.FilingStatus == Joint {
			//delta := ms.Ip.Age1 - ms.Ip.Age2
			ms.Ao.Output(fmt.Sprintf("%3d/%3d:", year+ms.Ip.StartPlan, year+ms.Ip.StartPlan-ms.Ip.AgeDelta))
		} else {
			ms.Ao.Output(fmt.Sprintf(" %3d:", year+ms.Ip.StartPlan))
		}
		items := []float64{withdrawal[IRA] / ms.OneK, deposit[IRA] / ms.OneK, rmdref / ms.OneK, // IRA
			withdrawal[Roth] / ms.OneK, deposit[Roth] / ms.OneK, // Roth
			withdrawal[Aftertax] / ms.OneK, deposit[Aftertax] / ms.OneK, //D, // AftaTax
			AccessVector(ms.Income[0], year) / ms.OneK, AccessVector(ms.SS[0], year) / ms.OneK, AccessVector(ms.Expenses[0], year) / ms.OneK,
			(tax + cgtax + earlytax) / ms.OneK}
		for _, f := range items {
			format := fmt.Sprintf("&@%%%d.0f", fieldwidth)
			str := fmt.Sprintf(format, f)
			ms.Ao.Output(str)
			//ao.output("&@{:>{width}.0f}".format(i, width=fieldwidth))
		}
		s := (*xp)[ms.Vindx.S(year)] / ms.OneK
		star := ' '
		if spendable+0.1 < (*xp)[ms.Vindx.S(year)] || spendable-0.1 > (*xp)[ms.Vindx.S(year)] {
			// replace the model ouput with actual value and add star
			// to indicate that we did so
			s = spendable / ms.OneK
			star = '*'
		}
		ms.Ao.Output(fmt.Sprintf("&@%7.0f%c", s, star))
		ms.Ao.Output("\n")
	}
	ms.activitySummaryHeader(fieldwidth)
}

func (ms ModelSpecs) printIncomeHeader(headerkeylist []string, countlist []int, incomeCat []string, fieldwidth int) {
	if len(countlist) != len(incomeCat) {
		e := fmt.Errorf("printIncomeHearder: lenth of countlist(%d) != length of incomeCat(%d)", len(countlist), len(incomeCat))
		panic(e)
	}
	var ageWidth int
	names := ""
	if ms.Ip.MyKey2 != "" && ms.Ip.FilingStatus == Joint {
		names = fmt.Sprintf("%s/%s", ms.Ip.MyKey1, ms.Ip.MyKey2)
		ageWidth = 8
	} else {
		if ms.Ip.MyKey1 != "nokey" {
			names = fmt.Sprintf("%s", ms.Ip.MyKey1)
		}
		ageWidth = 5
	}
	str := fmt.Sprintf("%[1]*.[1]*[2]s", ageWidth, names)
	ms.Ao.Output(str)
	for i := 0; i < len(countlist); i++ {
		if countlist[i] > 0 {
			ats := 1 // number of '@' to add
			if i > 0 {
				ats = countlist[i-1]
			}
			totalspace := fieldwidth*countlist[i] + countlist[i] - 1 // -1 is for the &
			str = fmt.Sprintf("&%s%-[3]*.[3]*[2]s", atv[:ats], incomeCat[i], totalspace)
			ms.Ao.Output(str)
		}
	}
	ms.Ao.Output("\n")
	str = fmt.Sprintf("%[1]*[2]s", ageWidth, "age ")
	ms.Ao.Output(str)
	for _, str := range headerkeylist {
		if str == "nokey" { // HAACCKKK
			str = "  "
		}
		ms.Ao.Output(fmt.Sprintf("&@%[2]*.[2]*[1]s", str, fieldwidth))
	}
	ms.Ao.Output("\n")
}

/**/
func (ms ModelSpecs) getSSIncomeAssetExpenseList() ([]string, []int, [][]float64) {
	typeList := []string{"SocialSecurity", "income", "asset", "expense"}
	headerlist := make([]string, 0)
	countlist := make([]int, 0)
	datamatrix := make([][]float64, 0)
	var vp *[][]float64
	var vt *[]string
	for _, t := range typeList {
		switch t {
		case "SocialSecurity":
			vp = &ms.SS
			vt = &ms.SStags
		case "income":
			vp = &ms.Income
			vt = &ms.Incometags
		case "asset":
			vp = &ms.AssetSale
			vt = &ms.Assettags
		case "expense":
			vp = &ms.Expenses
			vt = &ms.Expensetags
		}
		count := 0
		for elem := 1; elem < len(*vt); elem++ {
			/*
							if len(*vp) != len(*vt) {
								e := fmt.Errorf("getSSIncomeAssetExpenseList: %s vector lengths do not match (%d vs. %d)", t, len(*vp), len(*vt))
								fmt.Printf("*vp: %#v\n", *vp)
								fmt.Printf("*vt: %#v\n", *vt)
								panic(e)
				            }
			*/
			if *vp == nil {
				e := fmt.Errorf("getSSIncomeAssetExpenseList: %s data vector %s is nil", t, (*vt)[elem])
				panic(e)
			}
			datamatrix = append(datamatrix, (*vp)[elem])
			headerlist = append(headerlist, (*vt)[elem])
			count++
		}
		countlist = append(countlist, count)
	}
	return headerlist, countlist, datamatrix
}

func (ms ModelSpecs) PrintIncomeExpenseDetails() {
	if ms.OneK < 1 {
		e := fmt.Errorf("printIncomeExpenseDetails: ms.OneK is %f which is not allowed", ms.OneK)
		panic(e)
	}
	ms.Ao.Output("\nIncome and Expense Summary:\n\n")
	headerlist, countlist, datamatrix := ms.getSSIncomeAssetExpenseList()
	incomeCat := []string{"SSincome:", "Income:", "AssetSale:", "Expense:"}
	fieldwidth := 8
	ms.printIncomeHeader(headerlist, countlist, incomeCat, fieldwidth)
	for year := 0; year < ms.Ip.Numyr; year++ {
		if ms.Ip.MyKey2 != "" && ms.Ip.FilingStatus == Joint {
			ms.Ao.Output(fmt.Sprintf("%3d/%3d:", year+ms.Ip.StartPlan, year+ms.Ip.StartPlan-ms.Ip.AgeDelta))
		} else {
			ms.Ao.Output(fmt.Sprintf(" %3d:", year+ms.Ip.StartPlan))
		}
		for i := 0; i < len(datamatrix); i++ {
			str := fmt.Sprintf("&@%[2]*.0[1]f", datamatrix[i][year]/ms.OneK, fieldwidth)
			ms.Ao.Output(str)
		}
		ms.Ao.Output("\n")
	}
	ms.printIncomeHeader(headerlist, countlist, incomeCat, fieldwidth)
}

/*
for _, p := range proverbs {
	n, err := writer.Write([]byte(p))
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	if n != len(p) {
		fmt.Println("failed to write data")
		os.Exit(1)
	}
}
*/
// Print out the active input parameters (string string map)
func PrintInputParamsStrMapToBuffer(m map[string]string) string {
	var writer bytes.Buffer

	writer.Write([]byte("InputParamsStrMap:\n"))
	for i, v := range InputStrDefs {
		if m[v] != "" {
			writer.Write([]byte(fmt.Sprintf("%3d&@'%32s':&@'%s'\n", i, v, m[v])))
		}
	}
	for j := 1; j < MaxStreams+1; j++ {
		for i, v := range InputStreamStrDefs {
			lineno := i + len(InputStrDefs) +
				(j-1)*len(InputStreamStrDefs)
			k := fmt.Sprintf("%s%d", v, j)
			if m[k] != "" {
				writer.Write([]byte(fmt.Sprintf("%3d&@'%32s':&@'%s'\n", lineno, k, m[k])))
			}
		}
	}
	writer.Write([]byte("\n"))
	return writer.String()
}

// Write to file the active input parameters (string string map)
func WriteFileInputParamsStrMap(f *os.File, m map[string]string) {
	if f != nil {
		ao := NewAppOutput(nil, f)
		ao.Output(PrintInputParamsStrMapToBuffer(m))
	}
}

// Print out the active input parameters (string string map)
func (ms ModelSpecs) PrintInputParamsStrMap(m map[string]string) {
	ms.Ao.Output(PrintInputParamsStrMapToBuffer(m))
}

func (ms ModelSpecs) printAccHeader() {
	if ms.Ip.FilingStatus == Joint && ms.Ip.MyKey2 != "" {
		ms.Ao.Output(fmt.Sprintf("%s/%s\n", ms.Ip.MyKey1, ms.Ip.MyKey2))
		ms.Ao.Output("    age ")
	} else {
		if ms.Ip.MyKey1 != "nokey" {
			ms.Ao.Output(fmt.Sprintf("%s\n", ms.Ip.MyKey1))
		}
		ms.Ao.Output(" age ")
	}
	if ms.Ip.Accmap[IRA] > 1 {
		str := fmt.Sprintf("&@%7s&@%7s&@%7s&@%7s&@%7s&@%7s&@%7s&@%7s",
			"IRA1", "fIRA1", "tIRA1", "RMDref1", "IRA2", "fIRA2",
			"tIRA2", "RMDref2")
		ms.Ao.Output(str)
	} else if ms.Ip.Accmap[IRA] == 1 {
		str := fmt.Sprintf("&@%7s&@%7s&@%7s&@%7s",
			"IRA", "fIRA", "tIRA", "RMDref")
		ms.Ao.Output(str)
	}
	if ms.Ip.Accmap[Roth] > 1 {
		str := fmt.Sprintf("&@%7s&@%7s&@%7s&@%7s&@%7s&@%7s",
			"Roth1", "fRoth1", "tRoth1", "Roth2", "fRoth2", "tRoth2")
		ms.Ao.Output(str)
	} else if ms.Ip.Accmap[Roth] == 1 {
		str := fmt.Sprintf("&@%7s&@%7s&@%7s", "Roth", "fRoth", "tRoth")
		ms.Ao.Output(str)
	}
	if ms.Ip.Accmap[Aftertax] == 1 {
		str := fmt.Sprintf("&@%7s&@%7s&@%7s", "AftaTx", "fAftaTx", "tAftaTx")
		ms.Ao.Output(str)
	}
	ms.Ao.Output("\n")
}

func (ms ModelSpecs) PrintAccountTrans(xp *[]float64) {

	ms.Ao.Output("\nAccount Transactions Summary:\n\n")
	ms.printAccHeader()
	//
	// Print pre-plan info
	//
	var index int
	if ms.Ip.MyKey2 != "" && ms.Ip.FilingStatus == Joint {
		ms.Ao.Output(fmt.Sprintf("%3d/%3d:", ms.Ip.Age1, ms.Ip.Age1-ms.Ip.AgeDelta))
	} else {
		ms.Ao.Output(fmt.Sprintf(" %3d:", ms.Ip.Age1))
	}
	for i := 0; i < ms.Ip.Accmap[IRA]; i++ {
		str := fmt.Sprintf("&@%7.0f&@%7.0f&@%7.0f&@%7.0f",
			ms.Accounttable[i].Origbal/ms.OneK, 0.0,
			ms.Accounttable[i].Contrib/ms.OneK, 0.0) // IRAn
		ms.Ao.Output(str)
	}
	for i := 0; i < ms.Ip.Accmap[Roth]; i++ {
		index = ms.Ip.Accmap[IRA] + i
		str := fmt.Sprintf("&@%7.0f&@%7.0f&@%7.0f",
			ms.Accounttable[index].Origbal/ms.OneK, 0.0,
			ms.Accounttable[index].Contrib/ms.OneK) // rothn
		ms.Ao.Output(str)
	}
	index = ms.Ip.Accmap[IRA] + ms.Ip.Accmap[Roth]
	if index == len(ms.Accounttable)-1 {
		str := fmt.Sprintf("&@%7.0f&@%7.0f&@%7.0f",
			ms.Accounttable[index].Origbal/ms.OneK, 0.0,
			ms.Accounttable[index].Contrib/ms.OneK) // aftertax
		ms.Ao.Output(str)
	}
	ms.Ao.Output("\n")
	ms.Ao.Output("Plan Start: ---------\n")
	//
	// Print plan info for each year
	// TODO clean up the if/else below to follow the above forloop pattern
	//
	for year := 0; year < ms.Ip.Numyr; year++ {
		rmdref := make([]float64, 2)
		for j := 0; j < intMin(2, len(ms.Accounttable)); j++ { // only first two accounts are type IRA w/ RMD
			if ms.Accounttable[j].acctype == IRA {
				rmd := ms.Ti.rmdNeeded(year, ms.matchRetiree(ms.Accounttable[j].mykey, year, true))
				if rmd > 0 {
					rmdref[j] = (*xp)[ms.Vindx.B(year, j)] / rmd
				}
			}
		}

		if ms.Ip.MyKey2 != "" && ms.Ip.FilingStatus == Joint {
			ms.Ao.Output(fmt.Sprintf("%3d/%3d:", year+ms.Ip.StartPlan, year+ms.Ip.StartPlan-ms.Ip.AgeDelta))
		} else {
			ms.Ao.Output(fmt.Sprintf(" %3d:", year+ms.Ip.StartPlan))
		}
		if ms.Ip.Accmap[IRA] > 1 {
			str := fmt.Sprintf("&@%7.0f&@%7.0f&@%7.0f&@%7.0f&@%7.0f&@%7.0f&@%7.0f&@%7.0f",
				(*xp)[ms.Vindx.B(year, 0)]/ms.OneK,
				(*xp)[ms.Vindx.W(year, 0)]/ms.OneK,
				ms.depositAmount(xp, year, 0)/ms.OneK,
				rmdref[0]/ms.OneK, // IRA1
				(*xp)[ms.Vindx.B(year, 1)]/ms.OneK,
				(*xp)[ms.Vindx.W(year, 1)]/ms.OneK,
				ms.depositAmount(xp, year, 1)/ms.OneK,
				rmdref[1]/ms.OneK) // IRA2
			ms.Ao.Output(str)
		} else if ms.Ip.Accmap[IRA] == 1 {
			str := fmt.Sprintf("&@%7.0f&@%7.0f&@%7.0f&@%7.0f",
				(*xp)[ms.Vindx.B(year, 0)]/ms.OneK,
				(*xp)[ms.Vindx.W(year, 0)]/ms.OneK,
				ms.depositAmount(xp, year, 0)/ms.OneK,
				rmdref[0]/ms.OneK) // IRA1
			ms.Ao.Output(str)
		}
		index := ms.Ip.Accmap[IRA]
		if ms.Ip.Accmap[Roth] > 1 {
			str := fmt.Sprintf("&@%7.0f&@%7.0f&@%7.0f&@%7.0f&@%7.0f&@%7.0f",
				(*xp)[ms.Vindx.B(year, index)]/ms.OneK,
				(*xp)[ms.Vindx.W(year, index)]/ms.OneK,
				ms.depositAmount(xp, year, index)/ms.OneK, // roth1
				(*xp)[ms.Vindx.B(year, index+1)]/ms.OneK,
				(*xp)[ms.Vindx.W(year, index+1)]/ms.OneK,
				ms.depositAmount(xp, year, index+1)/ms.OneK) // roth2
			ms.Ao.Output(str)
		} else if ms.Ip.Accmap[Roth] == 1 {
			str := fmt.Sprintf("&@%7.0f&@%7.0f&@%7.0f",
				(*xp)[ms.Vindx.B(year, index)]/ms.OneK,
				(*xp)[ms.Vindx.W(year, index)]/ms.OneK,
				ms.depositAmount(xp, year, index)/ms.OneK) // roth1
			ms.Ao.Output(str)
		}
		index = ms.Ip.Accmap[IRA] + ms.Ip.Accmap[Roth]
		//assert index == len(S.accounttable)-1
		if index == len(ms.Accounttable)-1 {
			str := fmt.Sprintf("&@%7.0f&@%7.0f&@%7.0f",
				(*xp)[ms.Vindx.B(year, index)]/ms.OneK,
				(*xp)[ms.Vindx.W(year, index)]/ms.OneK,
				ms.depositAmount(xp, year, index)/ms.OneK) // aftertax account
			ms.Ao.Output(str)
		}
		ms.Ao.Output("\n")
	}
	ms.Ao.Output("Plan End: -----------\n")
	//
	// Post plan info
	//
	year := ms.Ip.Numyr
	if ms.Ip.MyKey2 != "" && ms.Ip.FilingStatus == Joint {
		ms.Ao.Output(fmt.Sprintf("%3d/%3d:", year+ms.Ip.StartPlan, ms.Ip.Numyr+ms.Ip.StartPlan-ms.Ip.AgeDelta))
	} else {
		ms.Ao.Output(fmt.Sprintf(" %3d:", year+ms.Ip.StartPlan))
	}
	for i := 0; i < ms.Ip.Accmap[IRA]; i++ {
		str := fmt.Sprintf("&@%7.0f&@%7.0f&@%7.0f&@%7.0f",
			(*xp)[ms.Vindx.B(year, i)]/ms.OneK, 0.0, 0.0, 0.0) // IRAn
		ms.Ao.Output(str)
	}
	for i := 0; i < ms.Ip.Accmap[Roth]; i++ {
		index = ms.Ip.Accmap[IRA] + i
		str := fmt.Sprintf("&@%7.0f&@%7.0f&@%7.0f",
			(*xp)[ms.Vindx.B(year, index)]/ms.OneK, 0.0, 0.0) // rothn
		ms.Ao.Output(str)
	}
	index = ms.Ip.Accmap[IRA] + ms.Ip.Accmap[Roth]
	if index == len(ms.Accounttable)-1 {
		str := fmt.Sprintf("&@%7.0f&@%7.0f&@%7.0f",
			(*xp)[ms.Vindx.B(year, index)]/ms.OneK, 0.0, 0.0) // aftertax
		ms.Ao.Output(str)
	}
	ms.Ao.Output("\n")
	ms.printAccHeader()
}

func (ms ModelSpecs) printAccWithdHeader() {
	if ms.Ip.MyKey2 != "" && ms.Ip.FilingStatus == Joint {
		ms.Ao.Output(fmt.Sprintf("%s/%s\n", ms.Ip.MyKey1, ms.Ip.MyKey2))
		ms.Ao.Output("    age ")
	} else {
		if ms.Ip.MyKey1 != "nokey" {
			ms.Ao.Output(fmt.Sprintf("%s\n", ms.Ip.MyKey1))
		}
		ms.Ao.Output(" age ")
	}
	str := fmt.Sprintf("&@%7s&@%7s&@%8s&@%8s", "fACC", "Real", "%%Liq", "%%All")
	ms.Ao.Output(str)
	ms.Ao.Output("\n")
}

func (ms ModelSpecs) PrintAccountWithdrawals(xp *[]float64) {

	ms.Ao.Output("\nAccount Withdrawals Summary:\n\n")
	ms.printAccWithdHeader()
	//
	// Print plan withdrawals for each year
	// TODO clean up the if/else below to follow the above forloop pattern
	//
	for year := 0; year < ms.Ip.Numyr; year++ {
		adjInf := math.Pow(ms.Ip.IRate, float64(year))
		if ms.Ip.MyKey2 != "" && ms.Ip.FilingStatus == Joint {
			ms.Ao.Output(fmt.Sprintf("%3d/%3d:", year+ms.Ip.StartPlan, year+ms.Ip.StartPlan-ms.Ip.AgeDelta))
		} else {
			ms.Ao.Output(fmt.Sprintf(" %3d:", year+ms.Ip.StartPlan))
		}
		totWithdrawals := 0.0
		for j := 0; j < ms.Ip.Numacc; j++ {
			totWithdrawals += (*xp)[ms.Vindx.W(year, j)]
		}
		realWithdrawals := totWithdrawals / adjInf
		realPercentOfOrigLiquidBal := 100 * realWithdrawals / ms.LiquidAssetPlanStart
		realPercentOfOrigAllAssets := 100 * realWithdrawals / (ms.LiquidAssetPlanStart + ms.IlliquidAssetPlanStart)

		str := fmt.Sprintf("&@%7.0f&@%7.0f&@%7.2f&@%7.2f",
			totWithdrawals, realWithdrawals,
			realPercentOfOrigLiquidBal, realPercentOfOrigAllAssets)
		ms.Ao.Output(str)
		ms.Ao.Output("\n")
	}
	ms.Ao.Output("\n")
	ms.printAccWithdHeader()
}

func (ms ModelSpecs) printHeaderTax() {
	if ms.Ip.MyKey2 != "" && ms.Ip.FilingStatus == Joint {
		ms.Ao.Output(fmt.Sprintf("%s/%s\n", ms.Ip.MyKey1, ms.Ip.MyKey2))
		ms.Ao.Output("    age ")
	} else {
		if ms.Ip.MyKey1 != "nokey" {
			ms.Ao.Output(fmt.Sprintf("%s\n", ms.Ip.MyKey1))
		}
		ms.Ao.Output(" age ")
	}
	str := fmt.Sprintf("&@%7s&@%7s&@%7s&@%7s&@%7s&@%7s&@%7s&@%7s&@%8s&@%7s&@%7s&@%8s&@%7s&@%7s&@%7s",
		"fIRA", "tIRA", "TxbleO", "TxbleSS", "deduct",
		"T_inc", "earlyP", "fedtax", "mTaxB%%", "fAftaTx",
		"tAftaTx", "cgTax%%", "cgTax", "TFedTax", "spndble")
	ms.Ao.Output(str)
	ms.Ao.Output("\n")
}

func (ms ModelSpecs) PrintTax(xp *[]float64) {
	ms.Ao.Output("\nTax Summary:\n\n")
	ms.printHeaderTax()
	for year := 0; year < ms.Ip.Numyr; year++ {
		age := year + ms.Ip.StartPlan
		iMul := math.Pow(ms.Ip.IRate, float64(ms.Ip.PrePlanYears+year))
		//T, spendable, tax, rate, cgtax, earlytax, rothearly := ms.IncomeSummary(year, xp)
		T, _, tax, rate, cgtax, earlytax, rothearly := ms.IncomeSummary(year, xp)
		f := ms.cgTaxableFraction(year)
		ttax := tax + cgtax + earlytax
		withdrawal := map[Acctype]float64{IRA: 0, Roth: 0, Aftertax: 0}
		deposit := map[Acctype]float64{IRA: 0, Roth: 0, Aftertax: 0}
		for j := 0; j < len(ms.Accounttable); j++ {
			withdrawal[ms.Accounttable[j].acctype] += (*xp)[ms.Vindx.W(year, j)]
			deposit[ms.Accounttable[j].acctype] += ms.depositAmount(xp, year, j)
		}
		if ms.Ip.MyKey2 != "" && ms.Ip.FilingStatus == Joint {
			ms.Ao.Output(fmt.Sprintf("%3d/%3d:", age, age-ms.Ip.AgeDelta))
		} else {
			ms.Ao.Output(fmt.Sprintf(" %3d:", year+ms.Ip.StartPlan))
		}
		star := ' '
		if rothearly {
			star = '*'
		}
		str := fmt.Sprintf("&@%7.0f&@%7.0f&@%7.0f&@%7.0f&@%7.0f&@%7.0f",
			withdrawal[IRA]/ms.OneK,
			deposit[IRA]/ms.OneK,
			AccessVector(ms.Taxed, year)/ms.OneK,
			ms.Ti.SStaxable*AccessVector(ms.SS[0], year)/ms.OneK,
			ms.Ti.Stded*iMul/ms.OneK,
			T/ms.OneK)
		str += fmt.Sprintf("&@%6.0f%c", earlytax/ms.OneK, star)
		str += fmt.Sprintf("&@%7.0f&@%7.0f&@%7.0f&@%7.0f&@%7.0f&@%7.0f&@%7.0f&@%7.0f",
			tax/ms.OneK,
			rate*100,
			withdrawal[Aftertax]/ms.OneK,
			deposit[Aftertax]/ms.OneK,
			f*100,
			cgtax/ms.OneK,
			ttax/ms.OneK,
			(*xp)[ms.Vindx.S(year)]/ms.OneK)
		ms.Ao.Output(str)
		ms.Ao.Output("\n")
	}
	ms.printHeaderTax()
}

func (ms ModelSpecs) printHeaderTaxBrackets() {
	spaces := 44
	if ms.Ip.MyKey2 != "" && ms.Ip.FilingStatus == Joint {
		//ao.output("@@@@@@@%64s" % "Marginal Rate(%):")
		spaces = 47
	}
	ampWidth := spaces
	atWidth := 7
	str := fmt.Sprintf("%s%sMarginal Rate(%s):", ampv[:ampWidth], atv[:atWidth], "%%")
	ms.Ao.Output(str)
	for k := 0; k < len(*ms.Ti.Taxtable); k++ {
		//(cut, size, rate, base) = ms.Ti.taxtable[k]
		rate := (*ms.Ti.Taxtable)[k][2]
		ms.Ao.Output(fmt.Sprintf("&@%6.0f", rate*100))
	}
	ms.Ao.Output("\n")
	if ms.Ip.MyKey2 != "" && ms.Ip.FilingStatus == Joint {
		ms.Ao.Output(fmt.Sprintf("%s/%s\n", ms.Ip.MyKey1, ms.Ip.MyKey2))
		ms.Ao.Output("    age ")
	} else {
		if ms.Ip.MyKey1 != "nokey" {
			ms.Ao.Output(fmt.Sprintf("%s\n", ms.Ip.MyKey1))
		}
		ms.Ao.Output(" age ")
	}
	str = fmt.Sprintf("&@%7s&@%7s&@%7s&@%7s&@%7s&@%7s&@%7s",
		"fIRA", "tIRA", "TxbleO", "TxbleSS", "deduct",
		"T_inc", "fedtax")
	ms.Ao.Output(str)
	for k := 0; k < len(*ms.Ti.Taxtable); k++ {
		ms.Ao.Output(fmt.Sprintf("&@brckt%d", k))
	}
	ms.Ao.Output("&@brkTot\n")
}

func (ms ModelSpecs) PrintTaxBrackets(xp *[]float64) {
	// For the bracket output don't do any rounding (ms.OneK)
	ms.Ao.Output("\nOverall Tax Bracket Summary:\n")
	ms.printHeaderTaxBrackets()
	colstrlen := 0
	for year := 0; year < ms.Ip.Numyr; year++ {
		age := year + ms.Ip.StartPlan
		iMul := math.Pow(ms.Ip.IRate, float64(ms.Ip.PrePlanYears+year))
		//T, spendable, tax, rate, cgtax, earlytax, rothearly := ms.IncomeSummary(year, xp)
		T, _, tax, _, _, _, _ := ms.IncomeSummary(year, xp)
		//ttax := tax + cgtax
		if ms.Ip.MyKey2 != "" && ms.Ip.FilingStatus == Joint {
			ms.Ao.Output(fmt.Sprintf("%3d/%3d:", age, age-ms.Ip.AgeDelta))
		} else {
			ms.Ao.Output(fmt.Sprintf(" %3d:", age))
		}
		withdrawal := map[Acctype]float64{IRA: 0, Roth: 0, Aftertax: 0}
		deposit := map[Acctype]float64{IRA: 0, Roth: 0, Aftertax: 0}
		for j := 0; j < len(ms.Accounttable); j++ {
			withdrawal[ms.Accounttable[j].acctype] += (*xp)[ms.Vindx.W(year, j)]
			deposit[ms.Accounttable[j].acctype] += ms.depositAmount(xp, year, j)
		}
		str := fmt.Sprintf("&@%7.0f&@%7.0f&@%7.0f&@%7.0f&@%7.0f&@%7.0f&@%7.0f",
			withdrawal[IRA],                              // /ms.OneK,
			deposit[IRA],                                 // /ms.OneK,
			AccessVector(ms.Taxed, year),                 // /ms.OneK,
			ms.Ti.SStaxable*AccessVector(ms.SS[0], year), // /ms.OneK,
			ms.Ti.Stded*iMul /*/ms.OneK*/, T /*/ms.OneK*/, tax /*/ms.OneK*/)
		ms.Ao.Output(str)
		colstrlen = len(str)
		bt := 0.0
		for k := 0; k < len(*ms.Ti.Taxtable); k++ {
			ms.Ao.Output(fmt.Sprintf("&@%6.0f", (*xp)[ms.Vindx.X(year, k)]))
			bt += (*xp)[ms.Vindx.X(year, k)]
		}
		ms.Ao.Output(fmt.Sprintf("&@%6.0f\n", bt))
	}
	if ms.DeveloperInfo {
		var agestr string
		ms.Ao.Output("Yearly ordinary income bracket boundaries:\n")
		for year := 0; year < ms.Ip.Numyr; year++ {
			age := year + ms.Ip.StartPlan
			if ms.Ip.MyKey2 != "" && ms.Ip.FilingStatus == Joint {
				agestr = fmt.Sprintf("%3d/%3d:", age, age-ms.Ip.AgeDelta)
			} else {
				agestr = fmt.Sprintf(" %3d:", age)
			}
			ms.Ao.Output(agestr)
			adjInf := math.Pow(ms.Ip.IRate, float64(ms.Ip.PrePlanYears+year))
			format := fmt.Sprintf("%s%ds", "%", colstrlen-10)
			ms.Ao.Output(fmt.Sprintf(format, "@@@@@@"))
			ms.Ao.Output(fmt.Sprintf("Bracket size:@&&"))
			for k := 0; k < len(*ms.Ti.Taxtable)-1; k++ {
				bsize := (*ms.Ti.Taxtable)[k][1] * adjInf // mcg[i,l] inflation adjusted
				ms.Ao.Output(fmt.Sprintf("&&@%6.0f", bsize))
			}
			ms.Ao.Output("@&inf\n")
		}
	}
	ms.printHeaderTaxBrackets()
}

func (ms ModelSpecs) printHeaderShadowTaxBrackets() {
	if ms.Ip.MyKey2 != "" && ms.Ip.FilingStatus == Joint {
		ms.Ao.Output(fmt.Sprintf("%s/%s\n", ms.Ip.MyKey1, ms.Ip.MyKey2))
		ms.Ao.Output("    age ")
	} else {
		if ms.Ip.MyKey1 != "nokey" {
			ms.Ao.Output(fmt.Sprintf("%s\n", ms.Ip.MyKey1))
		}
		ms.Ao.Output(" age ")
	}
	str := fmt.Sprintf("&@%7s&@%7s&@%7s&@%7s&@%7s&@%7s&@%7s",
		"fIRA", "tIRA", "TxbleO", "TxbleSS", "deduct",
		"T_inc", "fedtax")
	ms.Ao.Output(str)
	for l := 0; l < len(*ms.Ti.Capgainstable); l++ {
		ms.Ao.Output(fmt.Sprintf("&@brckt%d", l))
	}
	ms.Ao.Output("&@brkTot\n")
}

func (ms ModelSpecs) PrintShadowTaxBrackets(xp *[]float64) {
	// For the bracket output don't do any rounding (ms.OneK)
	ms.Ao.Output("\nOverall Shadow Tax Bracket Summary:\n")
	ms.printHeaderShadowTaxBrackets()
	colstrlen := 0
	for year := 0; year < ms.Ip.Numyr; year++ {
		age := year + ms.Ip.StartPlan
		iMul := math.Pow(ms.Ip.IRate, float64(ms.Ip.PrePlanYears+year))
		//T, spendable, tax, rate, cgtax, earlytax, rothearly := ms.IncomeSummary(year, xp)
		T, _, tax, _, _, _, _ := ms.IncomeSummary(year, xp)
		//ttax := tax + cgtax
		if ms.Ip.MyKey2 != "" && ms.Ip.FilingStatus == Joint {
			ms.Ao.Output(fmt.Sprintf("%3d/%3d:", age, age-ms.Ip.AgeDelta))
		} else {
			ms.Ao.Output(fmt.Sprintf(" %3d:", age))
		}
		withdrawal := map[Acctype]float64{IRA: 0, Roth: 0, Aftertax: 0}
		deposit := map[Acctype]float64{IRA: 0, Roth: 0, Aftertax: 0}
		for j := 0; j < len(ms.Accounttable); j++ {
			withdrawal[ms.Accounttable[j].acctype] += (*xp)[ms.Vindx.W(year, j)]
			deposit[ms.Accounttable[j].acctype] += ms.depositAmount(xp, year, j)
		}
		str := fmt.Sprintf("&@%7.0f&@%7.0f&@%7.0f&@%7.0f&@%7.0f&@%7.0f&@%7.0f",
			withdrawal[IRA],                              // /ms.OneK,
			deposit[IRA],                                 // /ms.OneK,
			AccessVector(ms.Taxed, year),                 // /ms.OneK,
			ms.Ti.SStaxable*AccessVector(ms.SS[0], year), // /ms.OneK,
			ms.Ti.Stded*iMul /*/ms.OneK*/, T /*/ms.OneK*/, tax /*/ms.OneK*/)
		ms.Ao.Output(str)
		colstrlen = len(str)
		bt := 0.0
		for l := 0; l < len(*ms.Ti.Capgainstable); l++ {
			sy := 0.0
			if ms.Ip.Accmap[Aftertax] > 0 {
				sy = (*xp)[ms.Vindx.Sy(year, l)]
			}
			ms.Ao.Output(fmt.Sprintf("&@%6.0f", sy))
			bt += sy
		}
		ms.Ao.Output(fmt.Sprintf("&@%6.0f\n", bt))
	}
	if ms.DeveloperInfo {
		var agestr string
		ms.Ao.Output("Yearly Shadow bracket boundaries:\n")
		for year := 0; year < ms.Ip.Numyr; year++ {
			age := year + ms.Ip.StartPlan
			if ms.Ip.MyKey2 != "" && ms.Ip.FilingStatus == Joint {
				agestr = fmt.Sprintf("%3d/%3d:", age, age-ms.Ip.AgeDelta)
			} else {
				agestr = fmt.Sprintf(" %3d:", age)
			}
			ms.Ao.Output(agestr)
			adjInf := math.Pow(ms.Ip.IRate, float64(ms.Ip.PrePlanYears+year))
			format := fmt.Sprintf("%s%ds", "%", colstrlen-10)
			ms.Ao.Output(fmt.Sprintf(format, "@@@@@@"))
			ms.Ao.Output(fmt.Sprintf("Bracket Size:@&"))
			for l := 0; l < len(*ms.Ti.Capgainstable)-1; l++ {
				bsize := (*ms.Ti.Capgainstable)[l][1] * adjInf // mcg[i,l] inflation adjusted
				ms.Ao.Output(fmt.Sprintf("&&@%6.0f", bsize))
			}
			ms.Ao.Output("&@inf\n")
		}
	}
	ms.printHeaderShadowTaxBrackets()
}

func (ms ModelSpecs) printHeaderCapgainsBrackets() {
	spaces := 36
	if ms.Ip.MyKey2 != "" && ms.Ip.FilingStatus == Joint {
		spaces = 39
	}
	ampWidth := spaces
	atWidth := 6
	ms.Ao.Output(fmt.Sprintf("%s%sMarginal Rate(%s):", ampv[:ampWidth], atv[:atWidth], "%%"))
	for l := 0; l < len(*ms.Ti.Capgainstable); l++ {
		rate := (*ms.Ti.Capgainstable)[l][2]
		ms.Ao.Output(fmt.Sprintf("&@%6.0f", rate*100))
	}
	ms.Ao.Output("\n")
	if ms.Ip.MyKey2 != "" && ms.Ip.FilingStatus == Joint {
		ms.Ao.Output(fmt.Sprintf("%s/%s\n", ms.Ip.MyKey1, ms.Ip.MyKey2))
		ms.Ao.Output("    age ")
	} else {
		if ms.Ip.MyKey1 != "nokey" {
			ms.Ao.Output(fmt.Sprintf("%s\n", ms.Ip.MyKey1))
		}
		ms.Ao.Output(" age ")
	}
	str := fmt.Sprintf("&@%7s&@%7s&@%8s&@%7s&@%7s&@%7s",
		"fAftaTx", "TblASle", "cgTax%%", "cgTaxbl",
		"T_inc", "cgTax")
	ms.Ao.Output(str)
	for l := 0; l < len(*ms.Ti.Capgainstable); l++ {
		ms.Ao.Output(fmt.Sprintf("&@brckt%d", l))
	}
	ms.Ao.Output("&@brkTot\n")
}

func (ms ModelSpecs) PrintCapGainsBrackets(xp *[]float64) {
	// For the bracket output don't do any rounding (ms.OneK)
	ms.Ao.Output("\nOverall Capital Gains Bracket Summary:\n")
	ms.printHeaderCapgainsBrackets()
	colstrlen := 0
	for year := 0; year < ms.Ip.Numyr; year++ {
		age := year + ms.Ip.StartPlan
		//iMul := math.Pow(ms.Ip.IRate, float64(ms.Ip.PrePlanYears+year))
		f := 1.0
		atw := 0.0
		att := 0.0
		tas := 0.0
		if ms.Ip.Accmap[Aftertax] > 0 {
			f = ms.cgTaxableFraction(year)
			j := len(ms.Accounttable) - 1    // Aftertax / investment account always the last entry when present
			atw = (*xp)[ms.Vindx.W(year, j)] // Aftertax / investment account
			//
			// OK, this next bit can be confusing.
			// The sale of illiquid assets do not use the aftertax account
			// basis. They have been handled separately in ms.CgAssetTaxed.
			// Given this we only add to cg_taxable the withdrawals, as is
			// normal, plus the taxable amounts from asset sales.
			//
			att = (f * (*xp)[ms.Vindx.W(year, j)]) +
				AccessVector(ms.CgAssetTaxed, year) // non-basis fraction + cg taxable $
		}
		//T, spendable, tax, rate, cgtax, earlytax, rothearly := ms.IncomeSummary(year, xp)
		T, _, _, _, cgtax, _, _ := ms.IncomeSummary(year, xp)
		var agestr string
		if ms.Ip.MyKey2 != "" && ms.Ip.FilingStatus == Joint {
			agestr = fmt.Sprintf("%3d/%3d:", age, age-ms.Ip.AgeDelta)
		} else {
			agestr = fmt.Sprintf(" %3d:", age)
		}
		ms.Ao.Output(agestr)
		str := fmt.Sprintf("&@%7.0f&@%7.0f&@%7.0f&@%7.0f&@%7.0f&@%7.0f",
			atw, tas, f*100, att, T, cgtax)
		ms.Ao.Output(str)
		colstrlen = len(str)
		bt := 0.0
		bttax := 0.0
		for l := 0; l < len(*ms.Ti.Capgainstable); l++ {
			ty := 0.0
			if ms.Ip.Accmap[Aftertax] > 0 {
				ty = (*xp)[ms.Vindx.Y(year, l)]
			}
			ms.Ao.Output(fmt.Sprintf("&@%6.0f", ty))
			bt += ty
			bttax += ty * (*ms.Ti.Capgainstable)[l][2]
		}
		ms.Ao.Output(fmt.Sprintf("&@%6.0f\n", bt))
		/*
		   if args.verbosewga:
		       print(" cg bracket ttax %6.0f " % bttax, end='')
		       print("x->y[1]: %6.0f "% (res.x[vindx.x(year,0)]+res.x[vindx.x(year,1)]),end='')
		       print("x->y[2]: %6.0f "% (res.x[vindx.x(year,2)]+ res.x[vindx.x(year,3)]+ res.x[vindx.x(year,4)]+res.x[vindx.x(year,5)]),end='')
		       print("x->y[3]: %6.0f"% res.x[vindx.x(year,6)])
		*/
	}
	if ms.DeveloperInfo {
		var agestr string
		ms.Ao.Output("Yearly Capital Gains bracket boundaries:\n")
		for year := 0; year < ms.Ip.Numyr; year++ {
			age := year + ms.Ip.StartPlan
			if ms.Ip.MyKey2 != "" && ms.Ip.FilingStatus == Joint {
				agestr = fmt.Sprintf("%3d/%3d:", age, age-ms.Ip.AgeDelta)
			} else {
				agestr = fmt.Sprintf(" %3d:", age)
			}
			ms.Ao.Output(agestr)
			adjInf := math.Pow(ms.Ip.IRate, float64(ms.Ip.PrePlanYears+year))
			format := fmt.Sprintf("%s%ds", "%", colstrlen-10)
			ms.Ao.Output(fmt.Sprintf(format, "@@@@@"))
			ms.Ao.Output(fmt.Sprintf("Bracket Size:@&&"))
			for l := 0; l < len(*ms.Ti.Capgainstable)-1; l++ {
				bsize := (*ms.Ti.Capgainstable)[l][1] * adjInf // mcg[i,l] inflation adjusted
				ms.Ao.Output(fmt.Sprintf("&&@%6.0f", bsize))
			}
			ms.Ao.Output("@&inf\n")
		}
	}
	ms.printHeaderCapgainsBrackets()
}

func (ms ModelSpecs) printHeaderAssetSummary() {
	if ms.Ip.MyKey2 != "" && ms.Ip.FilingStatus == Joint {
		ms.Ao.Output(fmt.Sprintf("%s/%s\n", ms.Ip.MyKey1, ms.Ip.MyKey2))
		ms.Ao.Output("    age ")
	} else {
		if ms.Ip.MyKey1 != "nokey" {
			ms.Ao.Output(fmt.Sprintf("%s\n", ms.Ip.MyKey1))
		}
		ms.Ao.Output(" age ")
	}
	str := fmt.Sprintf("&@%20s&@%9s&@%9s&@%9s&@%9s&@%9s&@%9s\n",
		"Name", "Price", "BrkrFee", "Owed",
		"Net", "Basis", "Taxable")
	ms.Ao.Output(str)
}

func (ms ModelSpecs) PrintAssetSummary() {
	// For the bracket output don't do any rounding (ms.OneK)
	ms.Ao.Output("\nAsset Sales Summary:\n\n")
	ms.printHeaderAssetSummary()
	if ms.AssetSale != nil && len(ms.AssetSale[0]) >= ms.Ip.Numyr {
		for year := 0; year < ms.Ip.Numyr; year++ {
			age := year + ms.Ip.StartPlan

			//iMul := math.Pow(ms.Ip.IRate, float64(ms.Ip.PrePlanYears+year))

			if ms.AssetSale[0][year] != 0.0 {
				for indx := 1; indx < len(ms.AssetSale); indx++ {
					if ms.AssetSale[indx][year] != 0.0 {
						tag := ms.Assettags[indx]
						value, brate, assetRR, basis, owed, prime, _ := ms.AssetByTag(tag)
						price := value * math.Pow(assetRR, float64(age-ms.Ip.Age1))
						bfee := price * brate
						net := price*(1-brate) - owed
						if net < 0.0 {
							net = 0.0
						}
						taxable := price*(1-brate) - basis
						if prime == 1.0 {
							taxable -= ms.Ti.Primeresidence *
								math.Pow(ms.Ip.IRate, float64(age-ms.Ip.Age1))
							tag = "*" + tag
						}
						if taxable < 0.0 {
							taxable = 0.0
						}
						if ms.Ip.MyKey2 != "" && ms.Ip.FilingStatus == Joint {
							ms.Ao.Output(fmt.Sprintf("%3d/%3d:", age, age-ms.Ip.AgeDelta))
						} else {
							ms.Ao.Output(fmt.Sprintf(" %3d:", age))
						}
						str := fmt.Sprintf(
							"&@%20.20s&@%9.0f&@%9.0f&@%9.0f&@%9.0f&@%9.0f&@%9.0f\n",
							tag, price, bfee, owed,
							net, basis, taxable)
						ms.Ao.Output(str)
					}
				}
			}
		}
	}
}

// TODO FixMe this function should be place somewhere more appropriete
func (ms ModelSpecs) AssetByTag(name string) (value, brate, assetRR, basis, owed, prime, ageToSell float64) {
	for i := 0; i < len(ms.Ip.Assets); i++ {
		if ms.Ip.Assets[i].Tag == name {
			value = float64(ms.Ip.Assets[i].Value)
			brate = ms.Ip.Assets[i].BrokeragePercent / 100.0
			assetRR = ms.Ip.Assets[i].AssetRRate
			basis = float64(ms.Ip.Assets[i].CostAndImprovements)
			owed = float64(ms.Ip.Assets[i].OwedAtAgeToSell)
			prime = 0.0
			if ms.Ip.Assets[i].PrimaryResidence {
				prime = 1.0
			}
			ageToSell = float64(ms.Ip.Assets[i].AgeToSell)
			return
		}
	}
	return 0.0, 0.0, 0.0, 0.0, 0.0, 0.0, 0.0
}

// TODO FixMe this function should be place somewhere more appropriete
func (ms ModelSpecs) AssetByTagAndField(tag, field string) float64 {
	for i := 0; i < len(ms.Ip.Assets); i++ {
		var val float64
		if ms.Ip.Assets[i].Tag == tag {
			switch field {
			case "Value":
				val = float64(ms.Ip.Assets[i].Value)
			case "BrokeragePercent":
				val = ms.Ip.Assets[i].BrokeragePercent
			case "BroderFee":
				age := ms.Ip.Assets[i].AgeToSell
				assetRR := ms.Ip.Assets[i].AssetRRate
				value := float64(ms.Ip.Assets[i].Value)
				price := value * math.Pow(assetRR, float64(age-ms.Ip.Age1))
				brate := ms.Ip.Assets[i].BrokeragePercent / 100.0
				val = price * brate
			case "Taxable":
				age := ms.Ip.Assets[i].AgeToSell
				assetRR := ms.Ip.Assets[i].AssetRRate
				value := float64(ms.Ip.Assets[i].Value)
				price := value * math.Pow(assetRR, float64(age-ms.Ip.Age1))
				brate := ms.Ip.Assets[i].BrokeragePercent / 100.0
				//bfee := price * brate
				basis := float64(ms.Ip.Assets[i].CostAndImprovements)
				taxable := price*(1-brate) - basis
				if ms.Ip.Assets[i].PrimaryResidence {
					taxable -= ms.Ti.Primeresidence *
						math.Pow(ms.Ip.IRate, float64(age-ms.Ip.Age1))
				}
				if taxable < 0.0 {
					taxable = 0.0
				}
				val = taxable
			case "MaxTaxableExclution":
				if !ms.Ip.Assets[i].PrimaryResidence {
					val = 0.0
					break
				}
				age := ms.Ip.Assets[i].AgeToSell
				assetRR := ms.Ip.Assets[i].AssetRRate
				value := float64(ms.Ip.Assets[i].Value)
				price := value * math.Pow(assetRR, float64(age-ms.Ip.Age1))
				brate := ms.Ip.Assets[i].BrokeragePercent / 100.0
				//bfee := price * brate
				basis := float64(ms.Ip.Assets[i].CostAndImprovements)
				taxable := price*(1-brate) - basis
				taxableExclusion := ms.Ti.Primeresidence * math.Pow(ms.Ip.IRate, float64(age-ms.Ip.Age1))
				if taxableExclusion >= taxable {
					val = taxable
				} else {
					val = taxableExclusion
				}
			case "AssetRRate":
				val = ms.Ip.Assets[i].AssetRRate
			case "CostAndImprovements":
				val = float64(ms.Ip.Assets[i].CostAndImprovements)
			case "OwedAtAgeToSell":
				val = float64(ms.Ip.Assets[i].OwedAtAgeToSell)
			case "AgeToSell":
				val = float64(ms.Ip.Assets[i].AgeToSell)
			case "SellPrice":
				age := ms.Ip.Assets[i].AgeToSell
				assetRR := ms.Ip.Assets[i].AssetRRate
				value := float64(ms.Ip.Assets[i].Value)
				val = value * math.Pow(assetRR, float64(age-ms.Ip.Age1))
			case "SellNet":
				indx := 0
				for j := 0; j < len(ms.Assettags); j++ {
					if tag == ms.Assettags[j] {
						indx = j
						break
					}
				}
				val = AccessVector(ms.AssetSale[indx], ms.Ip.Assets[i].AgeToSell-ms.Ip.StartPlan)
			case "PrimaryResidence":
				val = 0.0
				if ms.Ip.Assets[i].PrimaryResidence {
					val = 1.0
				}
			}
			return val
		}
	}
	return 0.0
}

func (ms ModelSpecs) depositAmount(xp *[]float64, year int, index int) float64 {
	amount := (*xp)[ms.Vindx.D(year, index)]
	if ms.Accounttable[index].acctype == Aftertax {
		amount += AccessVector(ms.AssetSale[0], year)
	}
	return amount
}

func (ms ModelSpecs) ordinaryTaxable(year int, xp *[]float64) float64 {
	withdrawals := 0.0
	deposits := 0.0
	for j := 0; j < intMin(2, len(ms.Accounttable)); j++ {
		if ms.Accounttable[j].acctype == IRA {
			withdrawals += (*xp)[ms.Vindx.W(year, j)]
			deposits += ms.depositAmount(xp, year, j)
		}
	}
	T := withdrawals - deposits + AccessVector(ms.Taxed, year) + ms.Ti.SStaxable*AccessVector(ms.SS[0], year) - (ms.Ti.Stded * math.Pow(ms.Ip.IRate, float64(ms.Ip.PrePlanYears+year)))
	if T < 0 {
		T = 0
	}
	return T
}

// EarlyPenaltyCharged returns the charged amount of penalty and whether any portion of the penalty was from a Roth account
func (ms ModelSpecs) EarlyPenaltyCharged(year int, xp *[]float64) (earlytax float64, rothearly bool) {
	earlytax = 0.0
	rothearly = false
	for j, acc := range ms.Accounttable {
		if acc.acctype != Aftertax {
			if ms.Ti.applyEarlyPenalty(year, ms.matchRetiree(acc.mykey, year, true)) {
				earlytax += (*xp)[ms.Vindx.W(year, j)] * ms.Ti.Penalty
				if (*xp)[ms.Vindx.W(year, j)] > 0 && acc.acctype == Roth {
					rothearly = true
				}
			}
		}
	}
	return earlytax, rothearly
}

// IncomeSummary returns key indicators to summarize income
func (ms ModelSpecs) IncomeSummary(year int, xp *[]float64) (T, spendable, tax, rate, ncgtax, earlytax float64, rothearly bool) {
	// TODO clean up and simplify this fuction
	//
	// return ordinaryTaxable, Spendable, Tax, Rate, CG_Tax
	// Need to account for withdrawals from IRA deposited in Investment account NOT SPENDABLE

	earlytax, rothearly = ms.EarlyPenaltyCharged(year, xp)

	T = ms.ordinaryTaxable(year, xp)
	ntax := 0.0
	rate = 0.0
	for k := 0; k < len(*ms.Ti.Taxtable); k++ {
		ntax += (*xp)[ms.Vindx.X(year, k)] * (*ms.Ti.Taxtable)[k][2]
		if (*xp)[ms.Vindx.X(year, k)] > 0 {
			rate = (*ms.Ti.Taxtable)[k][2]
		}
	}
	tax = ntax
	D := 0.0
	ncgtax = 0.0
	//if S.accmap["aftertax"] > 0:
	for j := 0; j < len(ms.Accounttable); j++ {
		D += ms.depositAmount(xp, year, j)
	}
	if ms.Ip.Accmap[Aftertax] > 0 {
		for l := 0; l < len(*ms.Ti.Capgainstable); l++ {
			ncgtax += (*xp)[ms.Vindx.Y(year, l)] * (*ms.Ti.Capgainstable)[l][2]
		}
	}
	totWithdrawals := 0.0
	for j := 0; j < len(ms.Accounttable); j++ {
		totWithdrawals += (*xp)[ms.Vindx.W(year, j)]
	}
	spendable = totWithdrawals - D + AccessVector(ms.Income[0], year) + AccessVector(ms.SS[0], year) - AccessVector(ms.Expenses[0], year) - tax - ncgtax - earlytax + AccessVector(ms.AssetSale[0], year)
	return T, spendable, tax, rate, ncgtax, earlytax, rothearly
}

//func (ro ResultsOutput) getResultTotals(x []float64)
func (ms ModelSpecs) getResultTotals(xp *[]float64) (twithd, tcombined, tT, ttax, tcgtax, tearlytax, tspendable, tbeginbal, tendbal float64) {
	tincome := 0.0
	for year := 0; year < ms.Ip.Numyr; year++ {
		//T, spendable, tax, rate, cg_tax, earlytax, rothearly := ms.IncomeSummary(year, xp)
		T, spendable, tax, _, cgtax, earlytax, _ := ms.IncomeSummary(year, xp)
		totWithdrawals := 0.0
		for j := 0; j < ms.Ip.Numacc; j++ {
			totWithdrawals += (*xp)[ms.Vindx.W(year, j)]
		}
		twithd += totWithdrawals
		tincome += AccessVector(ms.Income[0], year) + AccessVector(ms.SS[0], year) // + withdrawals
		ttax += tax
		tcgtax += cgtax
		tearlytax += earlytax
		tT += T
		tspendable += spendable
	}
	tbeginbal = 0
	tendbal = 0
	for j := 0; j < ms.Ip.Numacc; j++ {
		tbeginbal += (*xp)[ms.Vindx.B(0, j)]
		//balance for the year following the last year
		tendbal += (*xp)[ms.Vindx.B(ms.Ip.Numyr, j)]
	}

	tcombined = tincome + twithd
	return twithd, tcombined, tT, ttax, tcgtax, tearlytax, tspendable, tbeginbal, tendbal
}

func (ms ModelSpecs) PrintBaseConfig(xp *[]float64) { // input is res.x
	totwithd, tincome, tTaxable, tincometax, tcgtax, tearlytax, tspendable, tbeginbal, tendbal := ms.getResultTotals(xp)
	ms.Ao.Output("\n")
	ms.Ao.Output("======\n")
	// Probably should switch from using tbeginbal to ms.LiquidAssetPlanStart
	str := fmt.Sprintf("Optimized for %s with %s status\n\tstarting at age %d with an estate of $%s liquid and $%s illiquid\n", ms.Ip.Maximize, ms.Ip.FilingStatus /*retirement_type?*/, ms.Ip.StartPlan, RenderFloat("#_###.", tbeginbal), RenderFloat("#_###.", ms.IlliquidAssetPlanStart))
	ms.Ao.Output(str)
	ms.Ao.Output("\n")
	if ms.Ip.Min == 0 && ms.Ip.Max == 0 {
		ms.Ao.Output("No desired minimum or maximum amount specified\n")
	} else if ms.Ip.Min == 0 {
		// max specified
		ms.Ao.Output(fmt.Sprintf("Maximum desired: $%s\n", RenderFloat("#_###.", float64(ms.Ip.Max))))

	} else {
		// min specified
		ms.Ao.Output(fmt.Sprintf("Minium desired: $%s\n", RenderFloat("#_###.", float64(ms.Ip.Min))))
	}
	ms.Ao.Output("\n")
	str = fmt.Sprintf("After tax yearly income: $%s adjusting for inflation\n\tand final estate at age %d with $%s liquid and $%s illiquid\n", RenderFloat("#_###.", (*xp)[ms.Vindx.S(0)]), ms.Ip.RetireAge1+ms.Ip.Numyr, RenderFloat("#_###.", tendbal), RenderFloat("#_###.", ms.IlliquidAssetPlanEnd))
	ms.Ao.Output(str)
	ms.Ao.Output("\n")
	ms.Ao.Output(fmt.Sprintf("total withdrawals: $%s\n", RenderFloat("#_###.", totwithd)))
	ms.Ao.Output(fmt.Sprintf("total ordinary taxable income $%s\n", RenderFloat("#_###.", tTaxable)))
	if tTaxable > 0.0 {
		s1 := RenderFloat("#_###.", tincometax+tearlytax)
		s2 := RenderFloat("##.#", 100*(tincometax+tearlytax)/tTaxable)
		ms.Ao.Output(fmt.Sprintf("total ordinary tax on all taxable income: $%s (%s%s) of taxable income\n", s1, s2, "%%"))
	} else {
		s1 := RenderFloat("#_###.", tincometax+tearlytax)
		ms.Ao.Output(fmt.Sprintf("total ordinary tax on all taxable income: $%s\n", s1))
	}
	ms.Ao.Output(fmt.Sprintf("total income (withdrawals + other) $%s\n", RenderFloat("#_###.", tincome)))
	ms.Ao.Output(fmt.Sprintf("total cap gains tax: $%s\n", RenderFloat("#_###.", tcgtax)))
	if int(tincome) > 0 {
		s1 := RenderFloat("#_###.", tincometax+tcgtax+tearlytax)
		s2 := RenderFloat("##.#", 100*(tincometax+tcgtax+tearlytax)/tincome)
		ms.Ao.Output(fmt.Sprintf("total all tax on all income: $%s (%s%s)\n", s1, s2, "%%"))
	}
	ms.Ao.Output(fmt.Sprintf("Total spendable (after tax money): $%s\n", RenderFloat("#_###.", tspendable)))
	ms.Ao.Output("\n")
}

/*
def verifyInputs( c , A , b ):
    m = len(A)
    n = len(A[0])
    if len(c) != n :
        print("lp: c vector incorrect length")
    if len(b) != m :
        print("lp: b vector incorrect length")

	# Do some sanity checks so that ab does not become singular during the
	# simplex solution. If the ZeroRow checks are removed then the code for
	# finding a set of linearly indepent columns must be improved.

	# Check that if a row of A only has zero elements that corresponding
	# element in b is zero, otherwise the problem is infeasible.
	# Otherwise return ErrZeroRow.
    zeroRows = 0
    for i in range(m):
        isZero = True
        for j in range(n) :
            if A[i][j] != 0 :
                isZero = False
                break
        if isZero and b[i] != 0 :
            # Infeasible
            print("ErrInfeasible -- row[%d]\n"% i)
        elif isZero :
            zeroRows+=1
            print("ErrZeroRow -- row[%d]\n"% i)
    # Check that if a column only has zero elements that the respective C vector
    # is positive (otherwise unbounded). Otherwise return ErrZeroColumn.
    zeroColumns = 0
    for j in range( n) :
        isZero = True
        for i in range( m) :
            if A[i][j] != 0 :
                isZero = False
                break
        if isZero and c[j] < 0 :
            print("ErrUnbounded -- column[%d] %s\n"% (j, vindx.varstr(j)))
        elif isZero :
            zeroColumns+=1
            print("ErrZeroColumn -- column[%d] %s\n"% (j, vindx.varstr(j)))
    print("\nZero Rows: %d, Zero Columns: %d\n"%(zeroRows, zeroColumns))
*/

/*
# Program entry point
# Instantiate the parser
if __name__ == '__main__':
    parser = argparse.ArgumentParser(description='Create an optimized finacial plan for retirement.')
    parser.add_argument('-v', '--verbose', action='store_true',
                        help="Extra output from solver")
    parser.add_argument('-va', '--verboseaccounttrans', action='store_true',
                        help="Output detailed account transactions from solver")
    parser.add_argument('-vi', '--verboseincome', action='store_true',
                        help="Output detailed list of income as specified in social security, income, asset and expense sections")
    parser.add_argument('-vt', '--verbosetax', action='store_true',
                        help="Output detailed tax info from solver")
    parser.add_argument('-vtb', '--verbosetaxbrackets', action='store_true',
                        help="Output detailed tax bracket info from solver")
    parser.add_argument('-vw', '--verbosewga', action='store_true',
                        help="Extra wga output from solver")
    parser.add_argument('-vm', '--verbosemodel', action='store_true',
                        help="Output the binding constraints of the LP model")
    parser.add_argument('-mall', '--verbosemodelall', action='store_true',
                        help="Output the entire LP model - not just the binding constraints")
    parser.add_argument('-mdp', '--modeldumptable', nargs='?',
                        const='./RPlanModel.dat', default='',
                        help="Output the entire LP model as c, A, b to file MODELDUMPTABLE (default: ./RPlanModel.dat)")
    parser.add_argument('-mld', '--modelloadtable', nargs='?',
                         const='./RPlanModel.dat', default='',
                        help="Load the LP model as c, A, b from file MODELLOADTABLE (default: ./RPlanModel.dat)")
    parser.add_argument('-ts', '--timesimplex', action='store_true',
                        help="Measure and print the amount of time used by the simplex solver")
    parser.add_argument('-csv', '--csv', nargs='?', const='./a.csv', default='',
                        help="Additionally write the output to a csv file CVS (default: ./ .cvs)")
    parser.add_argument('-1k', '--noroundingoutput', action='store_true',
                        help="Do not round the output to thousands")
    parser.add_argument('-nd', '--notdrarothradeposits', action='store_true',
                        help="Do not allow deposits to TDRA or ROTHRA accounts beyond explicit contributions")
    parser.add_argument('-V', '--version', action='version', version='%(prog)s Version '+__version__,
                        help="Display the program version number and exit")
    parser.add_argument('conffile', help='Require configuration input toml file')
    args = parser.parse_args()

    csv_file_name = None
    if args.csv != '':
        csv_file_name = args.csv
    ao = app_out.app_output(csv_file_name)

    taxinfo = tif.taxinfo()
    S = tomldata.Data(taxinfo)
    S.load_toml_file(args.conffile)
    S.process_toml_info()

    #print("\naccounttable: ", S.accounttable)

    if S.accmap['IRA']+S.accmap['roth']+S.accmap['aftertax'] == 0:
        print('Error: This app optimizes the withdrawals from your retirement account(s); you must have at least one specified in the input toml file.')
        exit(0)

    if args.verbosewga:
        print("accounttable: ", S.accounttable)

    non_binding_only = True
    if args.verbosemodelall:
        non_binding_only = False

    ms.OneK = 1000.0
    if args.noroundingoutput:
        ms.OneK = 1

    years = S.numyr
    taxbins = len(taxinfo.taxtable)
    cgbins = len(taxinfo.capgainstable)
    accounts = len(S.accounttable)

    vindx = vvar.vector_var_index(years, taxbins, cgbins, accounts, S.accmap)

    if precheck_consistancy():

        if args.modelloadtable == '':
            model = lp.lp_constraint_model(S, vindx, taxinfo.taxtable, taxinfo.capgainstable, taxinfo.penalty, taxinfo.stded, taxinfo.SS_taxable, args.verbose, args.notdrarothradeposits)
            c, A, b = model.build_model()
            if args.modeldumptable != '':
                #modelio.dumpModel(c, A, b)
                modelio.binDumpModel(c, A, b, args.modeldumptable)
        else:
            print("Loadfile: ", args.modelloadtable)
            c, A, b = modelio.binLoadModel(args.modelloadtable)
        #verifyInputs( c , A , b )
        if args.timesimplex:
            t = time.process_time()
        res = scipy.optimize.linprog(c, A_ub=A, b_ub=b,
                                 options={"disp": args.verbose,
                                          #"bland": True,
                                          "tol": 1.0e-7,
                                          "maxiter": 4000})
        if args.timesimplex:
            elapsed_time = time.process_time() - t
            print("\nElapsed Simplex time: %s seconds" % elapsed_time)
        if args.verbosemodel or args.verbosemodelall:
            if res.success == False:
                model.print_model_matrix(c, A, b, None, False)
                print(res)
                exit(1)
            else:
                model.print_model_matrix(c, A, b, res.slack, non_binding_only)
        if args.verbosewga or res.success == False:
            print(res)
            if res.success == False:
                exit(1)
        consistancy_check(res, years, taxbins, cgbins, accounts, S.accmap, vindx)

        ms.printActivitySummary(res.x)
        //print_model_results(res)
        if args.verboseincome:
            print_income_expense_details()
        if args.verboseaccounttrans:
            print_account_trans(res)
        if args.verbosetax:
            print_tax(res)
        if args.verbosetaxbrackets:
            print_tax_brackets(res)
            print_cap_gains_brackets(res)
        print_base_config(res)
*/
