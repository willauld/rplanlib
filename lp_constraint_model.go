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

	SS          [][]float64 // SS[0] is combined, SS[1] for retiree1 ...
	SStags      []string    // ...
	income      [][]float64 // income[0] is combined, income[1] first income stream...
	incometags  []string    // ...
	assetSale   [][]float64 // assetSale[0] combined, assetSale[1] first asset
	assettags   []string    // ...
	expenses    [][]float64 // expenses[0] combined, expensee[1] first expense
	expensetags []string    // ...

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
	if v1 == nil && v2 != nil {
		return v2, nil
	}
	if v1 != nil && v2 == nil {
		return v1, nil
	}
	if v1 == nil && v2 == nil {
		err := fmt.Errorf("mergeVectors: Can not merge two nil vectors")
		return nil, err
	}
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

	//fmt.Printf("yearly: %d, startAge %d, endAge %d, vsAge %d, veAge %d, irate %f, rrate %f, bage %d\n", yearly, startAge, endAge, vecStartAge, vecEndAge, iRate, rRate, baseAge)
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
	b := 0.0
	for age := baseAge; age < vecStartAge; age++ {
		preyears := age - baseAge
		// capture all contributions before start of retirement
		if age >= startAge && age <= endAge {
			b = float64(yearly) * math.Pow(iRate, float64(preyears))
			precontribs += b
			precontribsPlusReturns += b
		}
		precontribsPlusReturns *= rRate
		//fmt.Printf("preplan: age %d, preyears %d, yearly %d, thisyear %f, precontribs %f, precontribsPlusReturns %f\n", age, preyears, yearly, b, precontribs, precontribsPlusReturns)
	}
	var vec []float64
	if !zeroVector {
		vecSize := vecEndAge - vecStartAge
		vec = make([]float64, vecSize)
		for i := 0; i < vecSize; i++ {
			if i >= startAge-vecStartAge && i <= endAge-vecStartAge {
				to := float64(vecStartAge - baseAge + i)
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

func (ms ModelSpecs) verifyTaxableIncomeCoversContrib() error {
	//   and contrib is less than max
	// TODO: before return check the following:
	// - No TDRA contributions after reaching age 70
	// - Contributions are below legal maximums
	//   - Sum of contrib for all retirees is less than taxable income
	//   - Sum of contrib for all retirees is less thansum of legal max's
	//   - IRA+ROTH+401(k) Contributions are less than other taxable income for each
	//   - IRA+ROTH+401(k) Contributions are less than legal max for each
	// - this is all talking about uc(ij)
	for year := 0; year < ms.ip.numyr; year++ {
		contrib := 0.0
		jointMaxContrib := ms.ti.maxContribution(year, year+ms.ip.prePlanYears, ms.retirees, "", ms.ip.iRate)
		//print("jointMaxContrib: ", jointMaxContrib)
		//jointMaxContrib = maxContribution(year, None)
		for _, acc := range ms.accounttable {
			if acc.acctype != "aftertax" {
				carray := acc.contributions
				if carray != nil && carray[year] > 0 {
					contrib += carray[year]
					ownerage := ms.accountOwnerAge(year, acc)
					if ownerage >= 70 {
						if acc.acctype == "IRA" {
							e := fmt.Errorf("Error - IRS does not allow contributions to TDRA accounts after age 70.\n\tPlease correct the contributions at age %d,\n\tfrom the PRIMARY age line, for account type %s and owner %s", ms.ip.startPlan+year, acc.acctype, acc.mykey)
							return e
						}
					}
				}
			}
		}
		//print("Contrib amount: ", contrib)
		if accessVector(ms.taxed, year) < contrib {
			e := fmt.Errorf("Error - IRS requires contributions to retirement accounts\n\tbe less than your ordinary taxable income.\n\tHowever, contributions of $%.0f at age %d,\n\tfrom the PRIMARY age line, exceeds the taxable\n\tincome of $%.0f", contrib, ms.ip.startPlan+year, accessVector(ms.taxed, year))
			return e
		}
		if jointMaxContrib < contrib {
			e := fmt.Errorf("Error - IRS requires contributions to retirement accounts\n\tbe less than a define maximum.\n\tHowever, contributions of $%.0f at age %d,\n\tfrom the PRIMARY age line, exceeds your maximum amount\n\tof $%.0f", contrib, ms.ip.startPlan+year, jointMaxContrib)
			return e
		}
		//print("TYPE: ", self.retirement_type)
		if ms.ip.filingStatus == "joint" {
			// Need to check each retiree individually
			for _, v := range ms.retirees {
				//print(v)
				contrib := 0.0
				personalstartage := v.ageAtStart
				MaxContrib := ms.ti.maxContribution(year, year+ms.ip.prePlanYears, ms.retirees, v.mykey, ms.ip.iRate)
				//print("MaxContrib: ", MaxContrib, v.mykey)
				for _, acc := range ms.accounttable {
					//print(acc)
					if acc.acctype != "aftertax" {
						if acc.mykey == v.mykey {
							carray := acc.contributions
							if carray != nil {
								contrib += accessVector(carray, year)
							}
						}
					}
				}
				//print("Contrib amount: ", contrib)
				if MaxContrib < contrib {
					e := fmt.Errorf("Error - IRS requires contributions to retirement accounts be less than\n\ta define maximum.\n\tHowever, contributions of $%.0f at age %d, of the account owner's\n\tage line, exceeds the maximum personal amount of $%.0f for %s", contrib, personalstartage+year, MaxContrib, v.mykey)
					return e
				}
			}
		}
	}
	return nil
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
	tablefile *os.File) (*ModelSpecs, error) {

	//fmt.Printf("InputParams: %#v\n", ip)
	ms := ModelSpecs{
		ip:    ip,
		vindx: vindx,
		ti:    ti,
		ao:    NewAppOutput(csvfile, tablefile),
		allowTdraRothraDeposits: allowDeposits,
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
	if ip.TDRA1 > 0 || ip.TDRAContrib1 > 0 {
		a := account{}
		a.rRate = ip.rRate
		if ip.TDRARate1 != 0.0 {
			a.rRate = ip.TDRARate1
		}
		a.acctype = "IRA"
		a.mykey = ip.myKey1
		a.origbal = float64(ip.TDRA1)
		a.contrib = float64(ip.TDRAContrib1)
		a.contributions, dbal, _, err = genContrib(ip.TDRAContrib1,
			ms.convertAge(ip.TDRAContribStart1, a.mykey),
			ms.convertAge(ip.TDRAContribEnd1, a.mykey),
			ip.startPlan, ip.endPlan, ip.iRate, a.rRate, ip.age1)
		if err != nil {
			return nil, err
		}
		a.bal = a.origbal*math.Pow(a.rRate, float64(ip.prePlanYears)) + dbal
		ms.accounttable = append(ms.accounttable, a)
	}
	if ip.TDRA2 > 0 || ip.TDRAContrib2 > 0 {
		a := account{}
		a.rRate = ip.rRate
		if ip.TDRARate2 != 0 {
			a.rRate = ip.TDRARate2
		}
		a.acctype = "IRA"
		a.mykey = ip.myKey2
		a.origbal = float64(ip.TDRA2)
		a.contrib = float64(ip.TDRAContrib2)
		a.contributions, dbal, _, err = genContrib(ip.TDRAContrib2,
			ms.convertAge(ip.TDRAContribStart2, a.mykey),
			ms.convertAge(ip.TDRAContribEnd2, a.mykey),
			ip.startPlan, ip.endPlan, ip.iRate, a.rRate, ip.age1)
		if err != nil {
			return nil, err
		}
		a.bal = a.origbal*math.Pow(a.rRate, float64(ip.prePlanYears)) + dbal
		ms.accounttable = append(ms.accounttable, a)
	}
	if ip.Roth1 > 0 || ip.RothContrib1 > 0 {
		a := account{}
		a.rRate = ip.rRate
		if ip.RothRate1 != 0 {
			a.rRate = ip.RothRate1
		}
		a.acctype = "roth"
		a.mykey = ip.myKey1
		a.origbal = float64(ip.Roth1)
		a.contrib = float64(ip.RothContrib1)
		a.contributions, dbal, _, err = genContrib(ip.RothContrib1,
			ms.convertAge(ip.RothContribStart1, a.mykey),
			ms.convertAge(ip.RothContribEnd1, a.mykey),
			ip.startPlan, ip.endPlan, ip.iRate, a.rRate, ip.age1)
		if err != nil {
			return nil, err
		}
		a.bal = a.origbal*math.Pow(a.rRate, float64(ip.prePlanYears)) + dbal
		//fmt.Printf("Roth acc: %#v\n", a)
		ms.accounttable = append(ms.accounttable, a)
	}
	if ip.Roth2 > 0 || ip.RothContrib2 > 0 {
		a := account{}
		a.rRate = ip.rRate
		if ip.RothRate2 != 0 {
			a.rRate = ip.RothRate2
		}
		a.acctype = "roth"
		a.mykey = ip.myKey2
		a.origbal = float64(ip.Roth2)
		a.contrib = float64(ip.RothContrib2)
		a.contributions, dbal, _, err = genContrib(ip.RothContrib1,
			ms.convertAge(ip.RothContribStart2, a.mykey),
			ms.convertAge(ip.RothContribEnd2, a.mykey),
			ip.startPlan, ip.endPlan, ip.iRate, a.rRate, ip.age1)
		if err != nil {
			return nil, err
		}
		a.bal = a.origbal*math.Pow(a.rRate, float64(ip.prePlanYears)) + dbal
		ms.accounttable = append(ms.accounttable, a)
	}
	if ip.Aftatax > 0 || ip.AftataxContrib > 0 {
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
			return nil, err
		}
		a.bal = a.origbal*math.Pow(a.rRate, float64(ip.prePlanYears)) + dbal
		a.basis = a.origbasis + dbasis
		//fmt.Printf("aftertax accout: %#v\n", a)
		ms.accounttable = append(ms.accounttable, a)
	}
	if len(ms.accounttable) != ms.ip.numacc {
		e := fmt.Errorf("NewModelSpecs: len(accounttable): %d not equal to numacc: %d", len(ms.accounttable), ms.ip.numacc)
		return nil, e
	}

	ms.SS = make([][]float64, 0)
	SS, SS1, SS2, tags := processSS(&ip)
	ms.SS = append(ms.SS, SS)
	ms.SS = append(ms.SS, SS1)
	ms.SS = append(ms.SS, SS2)
	ms.SStags = tags

	//fmt.Printf("SS1: %v\n", ms.SS1)
	//fmt.Printf("SS2: %v\n", ms.SS2)
	//fmt.Printf("SS: %v\n", ms.SS)

	ms.income = make([][]float64, 1)
	ms.incometags = append(ms.incometags, "combined income")
	for i := 0; i < len(ip.income); i++ {
		tag := ip.income[i].Tag
		amount := ip.income[i].Amount
		startage := ip.income[i].StartAge
		endage := ip.income[i].EndAge
		infr := ms.ip.iRate
		if !ip.income[i].Inflate {
			infr = 1.0
		}
		income, err := buildVector(amount, startage, endage, ip.startPlan, ip.endPlan, infr, ip.age1)
		if err != nil {
			return nil, err
		}
		ms.income[0], err = mergeVectors(ms.income[0], income)
		if err != nil {
			return nil, err
		}
		if ip.income[i].Tax {
			ms.taxed, err = mergeVectors(ms.taxed, income)
			if err != nil {
				return nil, err
			}
		}
		ms.income = append(ms.income, income)
		ms.incometags = append(ms.incometags, tag)
	}

	ms.expenses = make([][]float64, 1)
	ms.expensetags = append(ms.expensetags, "combined expense")
	for i := 0; i < len(ip.expense); i++ {
		tag := ip.expense[i].Tag
		amount := ip.expense[i].Amount
		startage := ip.expense[i].StartAge
		endage := ip.expense[i].EndAge
		infr := ms.ip.iRate
		if !ip.expense[i].Inflate {
			infr = 1.0
		}
		expense, err := buildVector(amount, startage, endage, ip.startPlan, ip.endPlan, infr, ip.age1)
		if err != nil {
			return nil, err
		}
		ms.expenses[0], err = mergeVectors(ms.expenses[0], expense)
		if err != nil {
			return nil, err
		}
		ms.expenses = append(ms.expenses, expense)
		ms.expensetags = append(ms.expensetags, tag)
	}

	ms.assetSale = make([][]float64, 1)
	ms.assettags = append(ms.assettags, "combined assets")
	ms.illiquidAssetPlanStart = 0.0
	ms.illiquidAssetPlanEnd = 0.0
	noSell := false
	for i := 0; i < len(ip.assets); i++ {
		tag := ip.assets[i].tag
		value := float64(ip.assets[i].value)
		ageToSell := ip.assets[i].ageToSell
		if ageToSell < ip.startPlan || ageToSell > ip.endPlan {
			noSell = true
		}
		brokerageRate := ip.assets[i].brokeragePercent / 100.0
		if brokerageRate == 0 {
			brokerageRate = 0.04 // default to 4%
		}
		assetRRate := ip.assets[i].assetRRate
		costAndImprovements := float64(ip.assets[i].costAndImprovements)
		owedAtAgeToSell := float64(ip.assets[i].owedAtAgeToSell)
		primaryResidence := ip.assets[i].primaryResidence

		infr := ms.ip.iRate // asset rate of return defaults to inflation rate
		if assetRRate != 0 {
			infr = assetRRate
		}
		//fmt.Printf("tag: %s, value: %.0f, ageToSell: %d, brokerageRate: %f, infr: %f\n", tag, value, ageToSell, brokerageRate, infr)
		//fmt.Printf("owedAtAgeToSell: %.0f\n", owedAtAgeToSell)
		if ageToSell < ip.startPlan && ageToSell != 0 {
			err := fmt.Errorf("NewModelSpecs: Assets to be sold before plan start are not allow unless the age to sell is zero")
			return nil, err
		}
		ms.illiquidAssetPlanStart += value * math.Pow(infr, float64(ip.startPlan-ip.age1))
		temp := 0.0
		if ageToSell > ip.endPlan || ageToSell == 0 {
			// age after plan ends or zero cause value to remain
			temp = value * math.Pow(infr, float64(ip.startPlan+ip.numyr-ip.age1))
		}
		ms.illiquidAssetPlanEnd += temp

		assvec := make([]float64, ip.endPlan-ip.startPlan)
		tempvec := make([]float64, ip.endPlan-ip.startPlan)

		if !noSell {
			sellprice := value * math.Pow(infr, float64(ageToSell-ip.age1))
			income := sellprice*(1-brokerageRate) - owedAtAgeToSell
			if income < 0 {
				income = 0
			}
			cgtaxable := sellprice*(1-brokerageRate) - costAndImprovements
			//fmt.Printf("Asset sell price $%.0f, brokerageRate: %%%d, income $%.0f, cgtaxable $%.0f\n", sellprice, int(100*brokerageRate), income, cgtaxable)
			if primaryResidence {
				cgtaxable -= ti.Primeresidence * math.Pow(ip.iRate, float64(ageToSell-ip.age1))
				//fmt.Printf("cgtaxable: ", cgtaxable)
			}
			if cgtaxable < 0 {
				cgtaxable = 0
			}
			if income > 0 && ip.accmap["aftertax"] <= 0 {
				e := fmt.Errorf("Error - Assets to be sold must have an 'aftertax' investment\naccount into which to deposit the net proceeds. Please\nadd an 'aftertax' account to yourn configuration; the bal may be zero")
				return nil, e
			}
			year := ageToSell - ip.startPlan
			assvec[year] = income
			tempvec[year] = cgtaxable
		}
		ms.assetSale[0], err = mergeVectors(ms.assetSale[0], assvec)
		if err != nil {
			return nil, err
		}
		if ip.income[i].Tax {
			ms.cgAssetTaxed, err = mergeVectors(ms.cgAssetTaxed, tempvec)
			if err != nil {
				return nil, err
			}
		}
		ms.assetSale = append(ms.assetSale, assvec)
		ms.assettags = append(ms.assettags, tag)
	}
	/*
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
	*/

	err = ms.verifyTaxableIncomeCoversContrib()
	if err != nil {
		return nil, err
	}

	if ip.filingStatus == "joint" {
		// do nothing
	} else { // single or mseparate zero retiree2 info
		// TODO FIXME
	}
	return &ms, nil
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
		inc := accessVector(ms.income[0], year)
		ss := accessVector(ms.SS[0], year)
		exp := accessVector(ms.expenses[0], year)
		b = append(b, inc+ss-exp)
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
	if index <= 0 {
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
		//TODO: FIXME REMOVE THIS LOOP
		for _, v := range ms.accounttable {
			if v.acctype == "aftertax" {
				if v.bal > 0 { // don't want to divide by zero
					//
					// v.bal includes the rRate and v.basis includes
					// the additional contributions up until
					// startPlan so no need to inflate for ms.ip.prePlanYears
					//
					f = 1 - (v.basis / (v.bal * math.Pow(v.rRate, float64(year))))
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
