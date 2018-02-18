package rplanlib

import (
	"fmt"
	"math"
	"os"
)

type retiree struct { // TODO limit the fields here maxContribution is bigest (only) user
	age                     int
	ageAtStart              int
	throughAge              int
	mykey                   string
	definedContributionPlan bool
	dcpBuckets              []float64
}
type account struct {
	bal       float64
	origbal   float64
	basis     float64
	origbasis float64
	//estateTax     float64
	contributions []float64
	contrib       float64
	rRate         float64
	acctype       string
	mykey         string
}

// ModelSpecs struct contains the needed info for building an RPlanner constraint model
type ModelSpecs struct {
	ip    InputParams
	vindx VectorVarIndex
	ti    Taxinfo
	ao    AppOutput

	allowTdraRothraDeposits bool

	// The following was through 'S'
	illiquidAssetPlanStart float64
	illiquidAssetPlanEnd   float64
	accounttable           []account
	retirees               []retiree

	SS         [][]float64 // SS[0] is combined, SS[1] for retiree1 ...
	SStag      []string    // ...
	income     [][]float64 // income[0] is combined, income[1] first income stream...
	incometag  []string    // ...
	assetSale  [][]float64 // assetSale[0] combined, assetSale[1] first asset
	assettag   []string    // ...
	expenses   [][]float64 // expenses[0] combined, expensee[1] first expense
	expensetag []string    // ...

	taxed        []float64
	cgAssetTaxed []float64

	verbose bool
	errfile *os.File
	logfile *os.File
	//csvfile   *os.File
	//tablefile *os.File

	OneK float64
}

func intMax(a, b int) int {
	if a > b {
		return a
	}
	return b
}
func intMin(a, b int) int {
	if a < b {
		return a
	}
	return b
}
func checkStrconvError(e error) { // TODO: should I remove this?
	if e != nil {
		//fmt.Fprintf(ms.logfile, "checkStrconvError(): %s\n", e)
		panic(e)
	}
}

func accessVector(v []float64, index int) float64 {
	if v != nil {
		return v[index]
	}
	return 0.0
}

// mergeVectors sums two vectors of equal length returning a third vector
func mergeVectors(v1, v2 []float64) ([]float64, error) {
	if len(v1) != len(v2) {
		err := fmt.Errorf("mergeVectors: Can not merge, lengths do not match, %d vs %d", len(v1), len(v2))
		return nil, err
	}
	v3 := make([]float64, len(v1))
	for i := 0; i < len(v1); i++ {
		v3[i] = v1[i] + v2[i]
	}
	return v3, nil
}

// genContrib generates the starting balance/basis as well as a vector
// of 'iRate' adjusted 'yearly' contributions
// iRate = 1.0 implies no inflation of contributions
// All age values must be consistant, ie, in turms of retiree 1 or 2 but
// not mixed
func genContrib(yearly int,
	startAge int,
	endAge int,
	vecStartAge int,
	vecEndAge int,
	iRate float64,
	rRate float64,
	baseAge int) ([]float64, float64, float64, error) {

	zeroVector := false
	//verify that startAge and endAge are within vecStart and end
	if vecStartAge > vecEndAge {
		err := fmt.Errorf("vec start age (%d) is greater than vec end age (%d)", vecStartAge, vecEndAge)
		return nil, 0.0, 0.0, err
	}
	if startAge > endAge {
		err := fmt.Errorf("start age (%d) is greater than end age (%d)", startAge, endAge)
		return nil, 0.0, 0.0, err
	}
	if startAge > vecEndAge {
		zeroVector = true
	}
	if endAge < vecStartAge {
		zeroVector = true
	}
	if endAge > vecEndAge && startAge < vecEndAge {
		endAge = vecEndAge
	}
	precontribs := 0.0
	precontribsPlusReturns := 0.0
	for age := startAge; age < vecStartAge; age++ {
		preyears := age - baseAge
		// capture all contributions before start of retirement
		b := float64(yearly) * math.Pow(iRate, float64(preyears))
		precontribs += b
		precontribsPlusReturns += b
		precontribsPlusReturns *= rRate
	}
	var vec []float64
	if !zeroVector {
		vecSize := vecEndAge - vecStartAge
		vec := make([]float64, vecSize)
		for i := 0; i < vecSize; i++ {
			if i >= startAge-vecStartAge && i <= endAge-vecStartAge {
				to := float64(startAge - baseAge + i)
				adj := math.Pow(iRate, to)
				vec[i] = float64(yearly) * adj // or something like this FIXME TODO
			}
		}
	}
	return vec, precontribsPlusReturns, precontribs, nil
}

// buildVector creates a vector with a 'rate' adjusted 'yearly' amount in buckets between start and end age
func buildVector(yearly, startAge, endAge, vecStartAge, vecEndAge int, rate float64, baseAge int) ([]float64, error) {
	zeroVector := false
	//verify that startAge and endAge are within vecStart and end
	if vecStartAge > vecEndAge {
		err := fmt.Errorf("vec start age (%d) is greater than vec end age (%d)", vecStartAge, vecEndAge)
		return nil, err
	}
	if startAge > endAge {
		err := fmt.Errorf("start age (%d) is greater than end age (%d)", startAge, endAge)
		return nil, err
	}

	if startAge < vecStartAge && endAge >= vecStartAge {
		startAge = vecStartAge
	}
	if startAge > vecEndAge {
		zeroVector = true
	}
	if endAge < vecStartAge {
		zeroVector = true
	}
	if endAge > vecEndAge && startAge < vecEndAge {
		endAge = vecEndAge
	}
	vecSize := vecEndAge - vecStartAge
	vec := make([]float64, vecSize)
	if !zeroVector {
		for i := 0; i < vecSize; i++ {
			if i >= startAge-vecStartAge && i <= endAge-vecStartAge {
				to := float64(startAge - baseAge + i)
				adj := math.Pow(rate, to)
				vec[i] = float64(yearly) * adj // or something like this FIXME TODO
			}
		}
	}
	return vec, nil
}

// NewModelSpecs creates a ModelSpecs object
func NewModelSpecs(vindx VectorVarIndex,
	ti Taxinfo,
	ip InputParams,
	verbose bool,
	allowDeposits bool,
	errfile *os.File,
	logfile *os.File,
	csvfile *os.File,
	tablefile *os.File) ModelSpecs {

	ms := ModelSpecs{
		ip:    ip,
		vindx: vindx,
		ti:    ti,
		ao:    NewAppOutput(csvfile, tablefile),
		allowTdraRothraDeposits: allowDeposits,
		illiquidAssetPlanStart:  0.0,
		illiquidAssetPlanEnd:    0.0,
		verbose:                 verbose,
		errfile:                 errfile,
		logfile:                 logfile,
		//csvfile:                 csvfile,
		//tablefile:               tablefile,
		OneK: 1.0, //1000.0,
	}

	retirees := []retiree{
		{
			age:        ip.age1,
			ageAtStart: ip.age1 + ip.prePlanYears,
			throughAge: ip.planThroughAge1,
			mykey:      ip.myKey1,
			definedContributionPlan: false,
			dcpBuckets:              nil,
		},
	}
	if ip.filingStatus == "joint" {
		r2 := retiree{
			age:        ip.age2,
			ageAtStart: ip.age2 + ip.prePlanYears,
			throughAge: ip.planThroughAge2,
			mykey:      ip.myKey2,
			definedContributionPlan: false,
			dcpBuckets:              nil,
		}
		retirees = append(retirees, r2)
	}
	ms.retirees = retirees
	//fmt.Fprintf(ms.logfile, "NewModelSpec: numacc: %d, accmap: %v\n", ms.ip.numacc, ms.ip.accmap)

	var err error
	var dbal float64
	const maxPossibleAccounts = 5
	if ip.TDRA1 > 0 {
		a := account{}
		a.rRate = ip.rRate
		if ip.TDRARate1 != 0.0 {
			a.rRate = ip.TDRARate1
		}
		a.acctype = "IRA"
		a.mykey = "retiree1" // need to make this definable for pc versions
		a.origbal = float64(ip.TDRA1)
		a.contrib = float64(ip.TDRAContrib1)
		a.contributions, dbal, _, err = genContrib(ip.TDRAContrib1,
			ms.convertAge(ip.TDRAContribStart1, a.mykey),
			ms.convertAge(ip.TDRAContribEnd1, a.mykey),
			ip.startPlan, ip.endPlan, ip.iRate, a.rRate, ip.age1)
		if err != nil {
			panic(err)
		}
		a.bal = a.origbal*math.Pow(a.rRate, float64(ip.prePlanYears)) + dbal
		ms.accounttable = append(ms.accounttable, a)
	}
	if ip.TDRA2 > 0 {
		a := account{}
		a.rRate = ip.rRate
		if ip.TDRARate2 != 0 {
			a.rRate = ip.TDRARate2
		}
		a.acctype = "IRA"
		a.mykey = "retiree2" // need to make this definable for pc versions
		a.origbal = float64(ip.TDRA2)
		a.contrib = float64(ip.TDRAContrib2)
		a.contributions, dbal, _, err = genContrib(ip.TDRAContrib2,
			ms.convertAge(ip.TDRAContribStart2, a.mykey),
			ms.convertAge(ip.TDRAContribEnd2, a.mykey),
			ip.startPlan, ip.endPlan, ip.iRate, a.rRate, ip.age1)
		if err != nil {
			panic(err)
		}
		a.bal = a.origbal*math.Pow(a.rRate, float64(ip.prePlanYears)) + dbal
		ms.accounttable = append(ms.accounttable, a)
	}
	if ip.Roth1 > 0 {
		a := account{}
		a.rRate = ip.rRate
		if ip.RothRate1 != 0 {
			a.rRate = ip.RothRate1
		}
		a.acctype = "roth"
		a.mykey = "retiree1" // need to make this definable for pc versions
		a.origbal = float64(ip.Roth1)
		a.contrib = float64(ip.RothContrib1)
		a.contributions, dbal, _, err = genContrib(ip.RothContrib1,
			ms.convertAge(ip.RothContribStart1, a.mykey),
			ms.convertAge(ip.RothContribEnd1, a.mykey),
			ip.startPlan, ip.endPlan, ip.iRate, a.rRate, ip.age1)
		if err != nil {
			panic(err)
		}
		a.bal = a.origbal*math.Pow(a.rRate, float64(ip.prePlanYears)) + dbal
		ms.accounttable = append(ms.accounttable, a)
	}
	if ip.Roth2 > 0 {
		a := account{}
		a.rRate = ip.rRate
		if ip.RothRate2 != 0 {
			a.rRate = ip.RothRate2
		}
		a.acctype = "roth"
		a.mykey = "retiree2" // need to make this definable for pc versions
		a.origbal = float64(ip.Roth2)
		a.contrib = float64(ip.RothContrib2)
		a.contributions, dbal, _, err = genContrib(ip.RothContrib1,
			ms.convertAge(ip.RothContribStart2, a.mykey),
			ms.convertAge(ip.RothContribEnd2, a.mykey),
			ip.startPlan, ip.endPlan, ip.iRate, a.rRate, ip.age1)
		if err != nil {
			panic(err)
		}
		a.bal = a.origbal*math.Pow(a.rRate, float64(ip.prePlanYears)) + dbal
		ms.accounttable = append(ms.accounttable, a)
	}
	if ip.Aftatax > 0 {
		var dbasis float64
		a := account{}
		a.rRate = ms.ip.rRate
		if ip.AftataxRate != 0 {
			a.rRate = ip.AftataxRate
		}
		a.acctype = "aftertax"
		a.mykey = "" // need to make this definable for pc versions
		a.origbal = float64(ip.Aftatax)
		a.origbasis = float64(ip.AftataxBasis)
		a.contrib = float64(ip.AftataxContrib)
		a.contributions, dbal, dbasis, err = genContrib(ip.AftataxContrib,
			ms.convertAge(ip.AftataxContribStart, a.mykey),
			ms.convertAge(ip.AftataxContribEnd, a.mykey),
			ip.startPlan, ip.endPlan, ip.iRate, a.rRate, ip.age1)
		if err != nil {
			panic(err)
		}
		a.bal = a.origbal*math.Pow(a.rRate, float64(ip.prePlanYears)) + dbal
		a.basis = a.origbasis + dbasis
		fmt.Printf("aftertax: %#v\n", a)
		ms.accounttable = append(ms.accounttable, a)
	}
	if len(ms.accounttable) != ms.ip.numacc {
		e := fmt.Errorf("NewModelSpecs: len(accounttable): %d not equal to numacc: %d", len(ms.accounttable), ms.ip.numacc)
		panic(e)
	}

	//income: []float64 // TODO add real income vector, dummy for now
	income1 := 0
	incomeStart1 := ip.startPlan
	incomeEnd1 := ip.endPlan
	income, err := buildVector(income1, incomeStart1, incomeEnd1, ip.startPlan, ip.endPlan, ms.ip.iRate, ip.age1)
	if err != nil {
		fmt.Fprintf(errfile, "BuildVector Failed: %s\n", err)
	}
	ms.income = make([][]float64, 0)
	ms.income = append(ms.income, income)

	/*
		ms.SS1, err := buildVector(ip.PIA1, ip.SSStart1, ip.endPlan, ip.startPlan, ip.endPlan, ms.ip.iRate, ip.age1)
		if err != nil {
			fmt.Fprintf(errfile, "BuildVector Failed: %s\n", err)
		}
		ms.SS2, err := buildVector(ip.PIA2, ip.SSStart2, ip.endPlan, ip.startPlan, ip.endPlan, ms.ip.iRate, ip.age1)
		if err != nil {
			fmt.Fprintf(errfile, "BuildVector Failed: %s\n", err)
		}
		ms.SS, err = mergeVectors(SS1, SS2)
		if err != nil {
			fmt.Fprintf(errfile, "mergeVector Failed: %s\n", err)
		}
	*/
	ms.SS = make([][]float64, 0)
	SS, SS1, SS2 := processSS(&ip)
	ms.SS = append(ms.SS, SS)
	ms.SS = append(ms.SS, SS1)
	ms.SS = append(ms.SS, SS2)

	//fmt.Printf("SS1: %v\n", ms.SS1)
	//fmt.Printf("SS2: %v\n", ms.SS2)
	//fmt.Printf("SS: %v\n", ms.SS)

	//expenses: []float64 // TODO add real income vector, dummy for now
	exp1 := 0
	expStart1 := ip.startPlan
	expEnd1 := ip.endPlan
	expenses, err := buildVector(exp1, expStart1, expEnd1, ip.startPlan, ip.endPlan, ms.ip.iRate, ip.age1)
	if err != nil {
		fmt.Fprintf(errfile, "BuildVector Failed: %s\n", err)
	}
	ms.expenses = make([][]float64, 0)
	ms.expenses = append(ms.expenses, expenses)

	//taxed: []float64 // TODO add real income vector, dummy for now
	tax1 := 0
	taxStart1 := ip.startPlan
	taxEnd1 := ip.endPlan
	taxed, err := buildVector(tax1, taxStart1, taxEnd1, ip.startPlan, ip.endPlan, ms.ip.iRate, ip.age1)
	if err != nil {
		fmt.Fprintf(errfile, "BuildVector Failed: %s\n", err)
	}
	ms.taxed = taxed

	//asset_sale: []float64
	assetSale, err := buildVector(0, ip.startPlan, ip.endPlan, ip.startPlan, ip.endPlan, ms.ip.iRate, ip.age1)
	if err != nil {
		fmt.Fprintf(errfile, "BuildVector Failed: %s\n", err)
	}
	ms.assetSale = make([][]float64, 0)
	ms.assetSale = append(ms.assetSale, assetSale)

	//cg_asset_taxed: []float64 // TODO add real income vector, dummy for now
	cgtax1 := 0
	cgtaxStart1 := ip.startPlan
	cgtaxEnd1 := ip.endPlan
	cgtaxed, err := buildVector(cgtax1, cgtaxStart1, cgtaxEnd1, ip.startPlan, ip.endPlan, ms.ip.iRate, ip.age1)
	if err != nil {
		fmt.Fprintf(errfile, "BuildVector Failed: %s\n", err)
	}
	ms.cgAssetTaxed = cgtaxed

	if ip.filingStatus == "joint" {
		// do nothing
	} else { // single or mseparate zero retiree2 info
		// TODO FIXME
	}
	return ms
}

// ModelNote contains section information for the constraint model
type ModelNote struct {
	index int
	note  string
}

// BuildModel for:
// Minimize: c^T * x
// Subject to: A_ub * x <= b_ub
// all vars positive
func (ms ModelSpecs) BuildModel() ([]float64, [][]float64, []float64, []ModelNote) {

	nvars := ms.vindx.Vsize
	A := make([][]float64, 0)
	b := make([]float64, 0)
	c := make([]float64, nvars)
	notes := make([]ModelNote, 0)

	//fmt.Printf("\nms.accounttable len: %d\n", len(ms.accounttable))

	//
	// Add objective function (S1') becomes (R1') if PlusEstate is added
	//
	for year := 0; year < ms.ip.numyr; year++ {
		c[ms.vindx.S(year)] = -1
	}
	//
	// Add objective function tax bracket forcing function
	//
	for year := 0; year < ms.ip.numyr; year++ {
		for k := 0; k < len(*ms.ti.Taxtable); k++ {
			// Multiplies the impact of higher brackets opposite to
			// optimization. The intent here is to pressure higher
			// brackets more and pack the lower brackets
			c[ms.vindx.X(year, k)] = float64(k) / 10.0
		}
	}
	//
	// Adder objective function (R1') when PlusEstate is added
	//
	if ms.ip.maximize == "PlusEstate" {
		for j := 0; j < len(ms.accounttable); j++ {
			estateTax := ms.ti.AccountEstateTax[ms.accounttable[j].mykey]
			c[ms.vindx.B(ms.ip.numyr, j)] = -1 * estateTax // account discount rate
		}
		fmt.Fprintf(ms.logfile, "\nConstructing Spending + Estate Model:\n")
		notes = append(notes, ModelNote{-1, "Objective function R1':"})
	} else {
		fmt.Fprintf(ms.logfile, "\nConstructing Spending Model:\n")

		startamount := 0.0
		for j := 0; j < len(ms.accounttable); j++ {
			startamount += ms.accounttable[j].bal
		}
		balancer := 1.0 / (startamount)
		for j := 0; j < len(ms.accounttable); j++ {
			estateTax := ms.ti.AccountEstateTax[ms.accounttable[j].acctype]
			c[ms.vindx.B(ms.ip.numyr, j)] = -1 * balancer * estateTax // balance and discount rate
		}
		notes = append(notes, ModelNote{-1, "Objective function S1':"})
	}
	//
	// Add constraint (2')
	//
	notes = append(notes, ModelNote{len(A), "Constraints 2':"})
	for year := 0; year < ms.ip.numyr; year++ {
		row := make([]float64, nvars)
		for j := 0; j < len(ms.accounttable); j++ {
			p := 1.0
			if ms.accounttable[j].acctype != "aftertax" {
				if ms.ti.applyEarlyPenalty(year, ms.matchRetiree(ms.accounttable[j].mykey)) { // TODO: should applyEarlyPenalty() return the penalty amount, spimplifying things?
					p = 1 - ms.ti.Penalty
				}
			}
			row[ms.vindx.W(year, j)] = -1 * p
		}
		for k := 0; k < len(*ms.ti.Taxtable); k++ {
			row[ms.vindx.X(year, k)] = (*ms.ti.Taxtable)[k][2] // income tax
		}
		if ms.ip.accmap["aftertax"] > 0 {
			for l := 0; l < len(*ms.ti.Capgainstable); l++ {
				row[ms.vindx.Y(year, l)] = (*ms.ti.Capgainstable)[l][2] // cap gains tax
			}
		}
		for j := 0; j < len(ms.accounttable); j++ {
			row[ms.vindx.D(year, j)] = 1
		}
		row[ms.vindx.S(year)] = 1
		A = append(A, row)
		b = append(b, accessVector(ms.income[0], year)+accessVector(ms.SS[0], year)-accessVector(ms.expenses[0], year))
	}
	//
	// Add constraint (3a')
	//
	notes = append(notes, ModelNote{len(A), "Constraints 3a':"})
	for year := 0; year < ms.ip.numyr-1; year++ {
		row := make([]float64, nvars)
		row[ms.vindx.S(year+1)] = 1
		row[ms.vindx.S(year)] = -1 * ms.ip.iRate
		A = append(A, row)
		b = append(b, 0)
	}
	//
	// Add constraint (3b')
	//
	notes = append(notes, ModelNote{len(A), "Constraints 3b':"})
	for year := 0; year < ms.ip.numyr-1; year++ {
		row := make([]float64, nvars)
		row[ms.vindx.S(year)] = ms.ip.iRate
		row[ms.vindx.S(year+1)] = -1
		A = append(A, row)
		b = append(b, 0)
	}
	//
	// Add constrant (4') rows - not needed if [desired.income] is not defined in input
	//
	notes = append(notes, ModelNote{len(A), "Constraints 4':"})
	if ms.ip.min != 0 {
		for year := 0; year < 1; year++ { // Only needs setting at the beginning
			row := make([]float64, nvars)
			row[ms.vindx.S(year)] = -1
			A = append(A, row)
			b = append(b, float64(-ms.ip.min)) // [- d_i]
		}
	}

	//
	// Add constraints for (5') rows - not added if [max.income] is
	// not defined in input
	//
	notes = append(notes, ModelNote{len(A), "Constraints 5':"})
	if ms.ip.max != 0 {
		for year := 0; year < 1; year++ { // Only needs to be set at the beginning
			row := make([]float64, nvars)
			row[ms.vindx.S(year)] = 1
			A = append(A, row)
			b = append(b, float64(ms.ip.max)) // [ dm_i]
		}
	}

	//
	// Add constaints for (6') rows
	//
	notes = append(notes, ModelNote{len(A), "Constraints 6':"})
	for year := 0; year < ms.ip.numyr; year++ {
		row := make([]float64, nvars)
		for j := 0; j < len(ms.accounttable); j++ {
			if ms.accounttable[j].acctype != "aftertax" {
				row[ms.vindx.D(year, j)] = 1 // TODO if this is not executed, DONT register this constrain, DONT add to A and b
			}
		}
		A = append(A, row)
		//b+=[min(ms.income[year],ms.ti.maxContribution(year,None))]
		// using ms.taxed rather than ms.income because income could
		// include non-taxed anueities that don't count.
		None := ""
		infyears := ms.ip.prePlanYears + year
		b = append(b, math.Min(ms.taxed[year], ms.ti.maxContribution(year, infyears, ms.retirees, None, ms.ip.iRate)))
	}
	//
	// Add constaints for (7') rows
	//
	notes = append(notes, ModelNote{len(A), "Constraints 7':"})
	for year := 0; year < ms.ip.numyr; year++ {
		// TODO this is not needed when there is only one retiree
		infyears := ms.ip.prePlanYears + year
		for _, v := range ms.retirees {
			row := make([]float64, nvars)
			for j := 0; j < len(ms.accounttable); j++ {
				if v.mykey == ms.accounttable[j].mykey {
					// ["acctype"] != "aftertax": no "mykey" in aftertax
					// (this will either break or just not match - we
					// will see)
					row[ms.vindx.D(year, j)] = 1 // TODO if this is not executed, DONT register this constraint, DONT add to A and b
				}
			}
			A = append(A, row)
			b = append(b, ms.ti.maxContribution(year, infyears, ms.retirees, v.mykey, ms.ip.iRate))
		}
	}
	//
	// Add constaints for (8') rows
	//
	notes = append(notes, ModelNote{len(A), "Constraints 8':"})
	for year := 0; year < ms.ip.numyr; year++ {
		for j := 0; j < len(ms.accounttable); j++ {
			v := ms.accounttable[j].contributions
			if v != nil {
				if v[year] > 0 {
					row := make([]float64, nvars)
					row[ms.vindx.D(year, j)] = -1
					A = append(A, row)
					b = append(b, -1*v[year])
				}
			}
		}
	}
	//
	// Add constaints for (9') rows
	//
	notes = append(notes, ModelNote{len(A), "Constraints 9':"})
	for year := 0; year < ms.ip.numyr; year++ {
		for j := 0; j < intMin(2, len(ms.accounttable)); j++ {
			// at most the first two accounts are type IRA w/
			// RMD requirement
			if ms.accounttable[j].acctype == "IRA" {
				ownerage := ms.accountOwnerAge(year, ms.accounttable[j])
				if ownerage >= 70 {
					row := make([]float64, nvars)
					row[ms.vindx.D(year, j)] = 1
					A = append(A, row)
					b = append(b, 0)
				}
			}
		}
	}
	//
	// Add constaints for (N') rows
	//
	notes = append(notes, ModelNote{len(A), "Constraints N':"})
	if !ms.allowTdraRothraDeposits {
		for year := 0; year < ms.ip.numyr; year++ {
			for j := 0; j < len(ms.accounttable); j++ {
				v := ms.accounttable[j].contributions
				max := 0.0
				if v != nil {
					max = v[year]
				}
				if ms.accounttable[j].acctype != "aftertax" { //Todo: move this if statement up just under the for to remove all unnessasary work
					row := make([]float64, nvars)
					row[ms.vindx.D(year, j)] = 1
					A = append(A, row)
					b = append(b, max)
				}
			}
		}
	}
	//
	// Add constaints for (10') rows
	//
	notes = append(notes, ModelNote{len(A), "Constraints 10':"})
	for year := 0; year < ms.ip.numyr; year++ {
		for j := 0; j < intMin(2, len(ms.accounttable)); j++ {
			// at most the first two accounts are type IRA
			// w/ RMD requirement
			if ms.accounttable[j].acctype == "IRA" {
				rmd := ms.ti.rmdNeeded(year, ms.matchRetiree(ms.accounttable[j].mykey))
				if rmd > 0 {
					row := make([]float64, nvars)
					row[ms.vindx.B(year, j)] = 1 / rmd
					row[ms.vindx.W(year, j)] = -1
					A = append(A, row)
					b = append(b, 0)
				}
			}
		}
	}

	//
	// Add constraints for (11')
	//
	notes = append(notes, ModelNote{len(A), "Constraints 11':"})
	for year := 0; year < ms.ip.numyr; year++ {
		adjInf := math.Pow(ms.ip.iRate, float64(ms.ip.prePlanYears+year))
		row := make([]float64, nvars)
		for j := 0; j < intMin(2, len(ms.accounttable)); j++ {
			// IRA can only be in the first two accounts
			if ms.accounttable[j].acctype == "IRA" {
				row[ms.vindx.W(year, j)] = 1  // Account 0 is TDRA
				row[ms.vindx.D(year, j)] = -1 // Account 0 is TDRA
			}
		}
		for k := 0; k < len(*ms.ti.Taxtable); k++ {
			row[ms.vindx.X(year, k)] = -1
		}
		A = append(A, row)
		b = append(b, ms.ti.Stded*adjInf-accessVector(ms.taxed, year)-ms.ti.SStaxable*accessVector(ms.SS[0], year))
	}
	//
	// Add constraints for (12')
	//
	notes = append(notes, ModelNote{len(A), "Constraints 12':"})
	for year := 0; year < ms.ip.numyr; year++ {
		for k := 0; k < len(*ms.ti.Taxtable)-1; k++ {
			row := make([]float64, nvars)
			row[ms.vindx.X(year, k)] = 1
			A = append(A, row)
			b = append(b, ((*ms.ti.Taxtable)[k][1])*math.Pow(ms.ip.iRate, float64(ms.ip.prePlanYears+year))) // inflation adjusted
		}
	}
	//
	// Add constraints for (13a')
	//
	notes = append(notes, ModelNote{len(A), "Constraints 13a':"})
	if ms.ip.accmap["aftertax"] > 0 {
		for year := 0; year < ms.ip.numyr; year++ {
			f := ms.cgTaxableFraction(year)
			row := make([]float64, nvars)
			for l := 0; l < len(*ms.ti.Capgainstable); l++ {
				row[ms.vindx.Y(year, l)] = 1
			}
			// Awful Hack! If year of asset sale, assume w(i,j)-D(i,j) is
			// negative so taxable from this is zero
			if ms.cgAssetTaxed[year] <= 0 { // i.e., no sale
				j := len(ms.accounttable) - 1 // last Acc is investment / stocks
				row[ms.vindx.W(year, j)] = -1 * f
				row[ms.vindx.D(year, j)] = f
			}
			A = append(A, row)
			b = append(b, ms.cgAssetTaxed[year])
		}
	}
	//
	// Add constraints for (13b')
	//
	notes = append(notes, ModelNote{len(A), "Constraints 13b':"})
	if ms.ip.accmap["aftertax"] > 0 {
		for year := 0; year < ms.ip.numyr; year++ {
			f := ms.cgTaxableFraction(year)
			row := make([]float64, nvars)
			////// Awful Hack! If year of asset sale, assume w(i,j)-D(i,j) is
			////// negative so taxable from this is zero
			if ms.cgAssetTaxed[year] <= 0 { // i.e., no sale
				j := len(ms.accounttable) - 1 // last Acc is investment / stocks
				row[ms.vindx.W(year, j)] = f
				row[ms.vindx.D(year, j)] = -f
			}
			for l := 0; l < len(*ms.ti.Capgainstable); l++ {
				row[ms.vindx.Y(year, l)] = -1
			}
			A = append(A, row)
			b = append(b, -ms.cgAssetTaxed[year])
		}
	}
	//
	// Add constraints for (14')
	//
	notes = append(notes, ModelNote{len(A), "Constraints 14':"})
	if ms.ip.accmap["aftertax"] > 0 {
		for year := 0; year < ms.ip.numyr; year++ {
			adjInf := math.Pow(ms.ip.iRate, float64(ms.ip.prePlanYears+year))
			for l := 0; l < len(*ms.ti.Capgainstable)-1; l++ {
				row := make([]float64, nvars)
				row[ms.vindx.Y(year, l)] = 1
				for k := 0; k < len(*ms.ti.Taxtable)-1; k++ {
					if (*ms.ti.Taxtable)[k][0] >= (*ms.ti.Capgainstable)[l][0] && (*ms.ti.Taxtable)[k][0] < (*ms.ti.Capgainstable)[l+1][0] {
						row[ms.vindx.X(year, k)] = 1
					}
				}
				A = append(A, row)
				b = append(b, (*ms.ti.Capgainstable)[l][1]*adjInf) // mcg[i,l] inflation adjusted
			}
		}
	}
	//
	// Add constraints for (15a')
	//
	notes = append(notes, ModelNote{len(A), "Constraints 15a':"})
	for year := 0; year < ms.ip.numyr; year++ {
		for j := 0; j < len(ms.accounttable); j++ {
			//j = len(ms.accounttable)-1 // nl the last account, the investment account
			row := make([]float64, nvars)
			row[ms.vindx.B(year+1, j)] = 1 // b[i,j] supports an extra year
			row[ms.vindx.B(year, j)] = -1 * ms.accounttable[j].rRate
			row[ms.vindx.W(year, j)] = ms.accounttable[j].rRate
			row[ms.vindx.D(year, j)] = -1 * ms.accounttable[j].rRate
			A = append(A, row)
			// In the event of a sell of an asset for the year
			temp := 0.0
			if ms.accounttable[j].acctype == "aftertax" {
				temp = accessVector(ms.assetSale[0], year) * ms.accounttable[j].rRate //TODO test
			}
			b = append(b, temp)
		}
	}
	//
	// Add constraints for (15b')
	//
	notes = append(notes, ModelNote{len(A), "Constraints 15b':"})
	for year := 0; year < ms.ip.numyr; year++ {
		for j := 0; j < len(ms.accounttable); j++ {
			//j = len(ms.accounttable)-1 // nl the last account, the investment account
			row := make([]float64, nvars)
			row[ms.vindx.B(year, j)] = ms.accounttable[j].rRate
			row[ms.vindx.W(year, j)] = -1 * ms.accounttable[j].rRate
			row[ms.vindx.D(year, j)] = ms.accounttable[j].rRate
			row[ms.vindx.B(year+1, j)] = -1 ////// b[i,j] supports an extra year
			A = append(A, row)
			temp := 0.0
			if ms.accounttable[j].acctype == "aftertax" {
				temp = -1 * accessVector(ms.assetSale[0], year) * ms.accounttable[j].rRate //TODO test
			}
			b = append(b, temp)
		}
	}
	//
	// Constraint for (16a')
	//   Set the begining b[1,j] balances
	//
	notes = append(notes, ModelNote{len(A), "Constraints 16a':"})
	for j := 0; j < len(ms.accounttable); j++ {
		row := make([]float64, nvars)
		row[ms.vindx.B(0, j)] = 1
		A = append(A, row)
		b = append(b, ms.accounttable[j].bal)
	}
	//
	// Constraint for (16b')
	//   Set the begining b[1,j] balances
	//
	notes = append(notes, ModelNote{len(A), "Constraints 16b':"})
	for j := 0; j < len(ms.accounttable); j++ {
		row := make([]float64, nvars)
		row[ms.vindx.B(0, j)] = -1
		A = append(A, row)
		b = append(b, -1*ms.accounttable[j].bal)
	}
	//
	// Constrant for (17') is default for sycpy so no code is needed
	//
	notes = append(notes, ModelNote{len(A), "Constraints 17':"})
	if ms.verbose {
		fmt.Fprintf(ms.logfile, "Num vars: %d\n", len(c))
		fmt.Fprintf(ms.logfile, "Num contraints: %d\n", len(b))
		fmt.Fprintf(ms.logfile, "\n")
	}

	return c, A, b, notes
}

// accountOwnerAge finds the age of the retiree who owns the account
// Only valid in plan years
func (ms ModelSpecs) accountOwnerAge(year int, acc account) int {
	age := 0
	retireekey := acc.mykey
	v := ms.matchRetiree(retireekey)
	if v != nil {
		age = v.ageAtStart + year
	}
	return age
}

// matchRetiree searches retirees by key returning nil if not found
func (ms ModelSpecs) matchRetiree(retireekey string) *retiree {
	for _, v := range ms.retirees {
		if v.mykey == retireekey {
			return &v
		}
	}
	return nil
}

// TODO unit test me :-)
// convertAge converts an age for key1 to an age in the primary timeline
func (ms ModelSpecs) convertAge(age int, key string) int {
	index := -1
	for i, v := range ms.retirees {
		if v.mykey == key {
			index = i
		}
	}
	if index == -1 {
		return -1000 // TODO find a better response
	} else if index == 0 {
		return age
	}
	//delta := ms.retirees[0].age - ms.retirees[1].age
	return age + ms.ip.ageDelta
}

// cgTaxableFraction estimates the portion of capital gains not from basis
func (ms ModelSpecs) cgTaxableFraction(year int) float64 {
	// applies only in Plan years
	f := 1.0
	if ms.ip.accmap["aftertax"] > 0 {
		for _, v := range ms.accounttable {
			if v.acctype == "aftertax" {
				if v.bal > 0 {
					f = 1 - (v.basis / (v.bal * math.Pow(v.rRate, float64(year+ms.ip.prePlanYears))))
				}
				break // should be the last entry anyway but...
			}
		}
	}
	return f
}

// TODO: FIXME: Create UNIT tests: last two parameters need s vector (s is output from simplex run)
// printModelMatrix prints to object function (cx) and constraint matrix (Ax<=b)
func (ms ModelSpecs) printModelMatrix(c []float64, A [][]float64, b []float64, notes []ModelNote, s []float64, nonBindingOnly bool) {
	note := ""
	notesIndex := 0
	nextModelIndex := len(A) + 1 // beyond the end of A
	if notes != nil {
		nextModelIndex = notes[notesIndex].index
		note = notes[notesIndex].note
		notesIndex++
	}
	if nextModelIndex < 0 { // Object function index -1
		from := nextModelIndex
		nextModelIndex = notes[notesIndex].index
		to := nextModelIndex - 1
		fmt.Fprintf(ms.logfile, "\n##== [%d-%d]: %s ==##\n", from, to, note)
		note = notes[notesIndex].note
		notesIndex++
	}
	fmt.Fprintf(ms.logfile, "c: ")
	ms.printModelRow(c, false)
	fmt.Fprintf(ms.logfile, "\n")
	if !nonBindingOnly {
		fmt.Fprintf(ms.logfile, "B?  i: A_ub[i]: b[i]\n")
		for constraint := 0; constraint < len(A); constraint++ {
			if nextModelIndex == constraint {
				from := nextModelIndex
				nextModelIndex = notes[notesIndex].index
				to := nextModelIndex - 1
				for to < from {
					fmt.Fprintf(ms.logfile, "\n##== [%d-%d]: %s ==##\n", from, to, note)
					note = notes[notesIndex].note
					notesIndex++
					from = nextModelIndex
					nextModelIndex = notes[notesIndex].index
					to = nextModelIndex - 1
				}
				fmt.Fprintf(ms.logfile, "\n##== [%d-%d]: %s ==##\n", from, to, note)
				note = notes[notesIndex].note
				notesIndex++
			}
			if s == nil || s[constraint] > 0 {
				fmt.Fprintf(ms.logfile, "  ")
			} else {
				fmt.Fprintf(ms.logfile, "B ")
			}
			fmt.Fprintf(ms.logfile, "%3d: ", constraint)
			ms.printConstraint(A[constraint], b[constraint])
		}
	} else {
		fmt.Fprintf(ms.logfile, "  i: A_ub[i]: b[i]\n")
		j := 0
		for constraint := 0; constraint < len(A); constraint++ {
			if nextModelIndex == constraint {
				from := nextModelIndex
				nextModelIndex = notes[notesIndex].index
				to := nextModelIndex - 1
				for to < from {
					fmt.Fprintf(ms.logfile, "\n##== [%d-%d]: %s ==##\n", from, to, note)
					note = notes[notesIndex].note
					notesIndex++
					from = nextModelIndex
					nextModelIndex = notes[notesIndex].index
					to = nextModelIndex - 1
				}
				fmt.Fprintf(ms.logfile, "\n##== [%d-%d]: %s ==##\n", from, to, note)
				note = notes[notesIndex].note
				notesIndex++
			}
			if s[constraint] > 0 {
				j++
				fmt.Fprintf(ms.logfile, "%3d: ", constraint)
				ms.printConstraint(A[constraint], b[constraint])
			}
		}
		fmt.Fprintf(ms.logfile, "\n\n%d non-binding constrains printed\n", j)
	}
	fmt.Fprintf(ms.logfile, "\n")
}

func (ms ModelSpecs) printConstraint(row []float64, b float64) {
	ms.printModelRow(row, true)
	fmt.Fprintf(ms.logfile, "<= b[]: %6.2f\n", b)
}

func (ms ModelSpecs) printModelRow(row []float64, suppressNewline bool) {
	if ms.ip.numacc < 0 || ms.ip.numacc > 5 {
		e := fmt.Errorf("PrintModelRow: number of accounts is out of bounds, should be between [0, 5] but is %d", ms.ip.numacc)
		fmt.Fprintf(ms.logfile, "%s\n", e)
		return
	}
	for i := 0; i < ms.ip.numyr; i++ { // x[]
		for k := 0; k < len(*ms.ti.Taxtable); k++ {
			if row[ms.vindx.X(i, k)] != 0 {
				fmt.Fprintf(ms.logfile, "x[%d,%d]=%6.3f, ", i, k, row[ms.vindx.X(i, k)])
			}
		}
	}
	if ms.ip.accmap["aftertax"] > 0 {
		for i := 0; i < ms.ip.numyr; i++ { // y[]
			for l := 0; l < len(*ms.ti.Capgainstable); l++ {
				if row[ms.vindx.Y(i, l)] != 0 {
					fmt.Fprintf(ms.logfile, "y[%d,%d]=%6.3f, ", i, l, row[ms.vindx.Y(i, l)])
				}
			}
		}
	}
	for i := 0; i < ms.ip.numyr; i++ { // w[]
		for j := 0; j < ms.ip.numacc; j++ {
			if row[ms.vindx.W(i, j)] != 0 {
				fmt.Fprintf(ms.logfile, "w[%d,%d]=%6.3f, ", i, j, row[ms.vindx.W(i, j)])
			}
		}
	}
	for i := 0; i < ms.ip.numyr+1; i++ { // b[] has an extra year
		for j := 0; j < ms.ip.numacc; j++ {
			if row[ms.vindx.B(i, j)] != 0 {
				fmt.Fprintf(ms.logfile, "b[%d,%d]=%6.3f, ", i, j, row[ms.vindx.B(i, j)])
			}
		}
	}
	for i := 0; i < ms.ip.numyr; i++ { // s[]
		if row[ms.vindx.S(i)] != 0 {
			fmt.Fprintf(ms.logfile, "s[%d]=%6.3f, ", i, row[ms.vindx.S(i)])
		}
	}
	for i := 0; i < ms.ip.numyr; i++ { // D[]
		for j := 0; j < ms.ip.numacc; j++ {
			if row[ms.vindx.D(i, j)] != 0 {
				fmt.Fprintf(ms.logfile, "D[%d,%d]=%6.3f, ", i, j, row[ms.vindx.D(i, j)])
			}
		}
	}
	if !suppressNewline {
		fmt.Fprintf(ms.logfile, "\n")
	}
}
