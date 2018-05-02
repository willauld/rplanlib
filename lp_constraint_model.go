package rplanlib

import (
	"fmt"
	"math"
	"os"
)

type retiree struct {
	age                             int
	ageAtStart                      int
	throughAge                      int
	mykey                           string
	definedContributionPlanStartAge int
	definedContributionPlanEndAge   int
}
type account struct {
	bal       float64
	Origbal   float64
	basis     float64
	origbasis float64
	//estateTax     float64
	Contributions []float64
	Contrib       float64
	RRate         float64
	acctype       Acctype
	mykey         string
}

// ModelSpecs struct contains the needed info for building an RPlanner constraint model
type ModelSpecs struct {
	Ip    InputParams
	Vindx VectorVarIndex
	Ti    Taxinfo
	Ao    AppOutput

	AllowTdraRothraDeposits bool

	// The following was through 'S'
	LiquidAssetPlanStart   float64
	IlliquidAssetPlanStart float64
	IlliquidAssetPlanEnd   float64
	Accounttable           []account
	Retirees               []retiree

	SS          [][]float64 // SS[0] is combined, SS[1] for retiree1 ...
	SStags      []string    // ...
	Income      [][]float64 // income[0] is combined, income[1] first income stream...
	Incometags  []string    // ...
	AssetSale   [][]float64 // assetSale[0] combined, assetSale[1] first asset
	Assettags   []string    // ...
	Expenses    [][]float64 // expenses[0] combined, expensee[1] first expense
	Expensetags []string    // ...

	Taxed        []float64
	CgAssetTaxed []float64

	Errfile *os.File
	Logfile *os.File

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

func AccessVector(v []float64, index int) float64 {
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
				to := float64(vecStartAge - baseAge + i)
				adj := math.Pow(rate, to)
				vec[i] = float64(yearly) * adj
			}
		}
	}
	return vec, nil
}

func (ms ModelSpecs) verifyTaxableIncomeCoversContrib(mList *WarnErrorList) error {
	//   and contrib is less than max
	// TODO: before return check the following:
	// - No TDRA contributions after reaching age 70
	// - Contributions are below legal maximums
	//   - Sum of contrib for all retirees is less than taxable income
	//   - Sum of contrib for all retirees is less thansum of legal max's
	//   - IRA+ROTH+401(k) Contributions are less than other taxable income for each
	//   - IRA+ROTH+401(k) Contributions are less than legal max for each
	// - this is all talking about uc(ij)
	for year := 0; year < ms.Ip.Numyr; year++ {
		contrib := 0.0
		jointMaxContrib := ms.Ti.maxContribution(year, year+ms.Ip.PrePlanYears, ms.Retirees, "", ms.Ip.IRate)
		//print("jointMaxContrib: ", jointMaxContrib)
		//jointMaxContrib = maxContribution(year, None)
		for _, acc := range ms.Accounttable {
			if acc.acctype != Aftertax {
				carray := acc.Contributions
				if carray != nil && carray[year] > 0 {
					contrib += carray[year]
					ownerage := ms.accountOwnerAge(year, acc)
					if ownerage >= 70 {
						if acc.acctype == IRA {
							str := fmt.Sprintf("Error - IRS does not allow contributions to TDRA accounts after age 70.\n\tPlease correct the contributions at age %d,\n\tfrom the PRIMARY age line, for account type %s and owner %s", ms.Ip.StartPlan+year, acc.acctype, acc.mykey)
							mList.AppendError(str)
							e := fmt.Errorf(str)
							return e
						}
					}
				}
			}
		}
		//print("Contrib amount: ", contrib)
		if AccessVector(ms.Taxed, year) < contrib {
			str := fmt.Sprintf("Error - IRS requires contributions to retirement accounts\n\tbe less than your ordinary taxable income.\n\tHowever, contributions of $%.0f at age %d,\n\tfrom the PRIMARY age line, exceeds the taxable\n\tincome of $%.0f", contrib, ms.Ip.StartPlan+year, AccessVector(ms.Taxed, year))
			mList.AppendError(str)
			e := fmt.Errorf(str)
			return e
		}
		if jointMaxContrib < contrib {
			str := fmt.Sprintf("Error - IRS requires contributions to retirement accounts\n\tbe less than a define maximum.\n\tHowever, contributions of $%.0f at age %d,\n\tfrom the PRIMARY age line, exceeds your maximum amount\n\tof $%.0f", contrib, ms.Ip.StartPlan+year, jointMaxContrib)
			mList.AppendError(str)
			e := fmt.Errorf(str)
			return e
		}
		//print("TYPE: ", self.retirement_type)
		if ms.Ip.FilingStatus == Joint {
			// Need to check each retiree individually
			for _, v := range ms.Retirees {
				//fmt.Printf("Retriee: %#v\n", v)
				contrib := 0.0
				personalstartage := v.ageAtStart
				MaxContrib := ms.Ti.maxContribution(year, year+ms.Ip.PrePlanYears, ms.Retirees, v.mykey, ms.Ip.IRate)
				//print("MaxContrib: ", MaxContrib, v.mykey)
				for _, acc := range ms.Accounttable {
					//print(acc)
					if acc.acctype != Aftertax {
						if acc.mykey == v.mykey {
							carray := acc.Contributions
							if carray != nil {
								contrib += AccessVector(carray, year)
							}
						}
					}
				}
				//print("Contrib amount: ", contrib)
				if MaxContrib < contrib {
					str := fmt.Sprintf("Error - IRS requires contributions to retirement accounts be less than\n\ta define maximum.\n\tHowever, contributions of $%.0f at age %d, of the account owner's\n\tage line, exceeds the maximum personal amount of $%.0f for %s", contrib, personalstartage+year, MaxContrib, v.mykey)
					mList.AppendError(str)
					e := fmt.Errorf(str)
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
	allowDeposits bool,
	RoundToOneK bool,
	fourPercentRule bool,
	errfile *os.File,
	logfile *os.File,
	csvfile *os.File,
	tablefile *os.File,
	wel *WarnErrorList) (*ModelSpecs, error) {

	//fmt.Printf("InputParams: %#v\n", ip)
	ms := ModelSpecs{
		Ip:    ip,
		Vindx: vindx,
		Ti:    ti,
		Ao:    NewAppOutput(csvfile, tablefile),
		AllowTdraRothraDeposits: allowDeposits,
		Errfile:                 errfile,
		Logfile:                 logfile,
		//csvfile:                 csvfile,
		//tablefile:               tablefile,
		OneK: 1000.0,
	}
	if !RoundToOneK {
		ms.OneK = 1.0
	}

	retirees := []retiree{
		{
			age:        ip.Age1,
			ageAtStart: ip.Age1 + ip.PrePlanYears,
			throughAge: ip.PlanThroughAge1,
			mykey:      ip.MyKey1,
			definedContributionPlanStartAge: ip.DefinedContributionPlanStart1,
			definedContributionPlanEndAge:   ip.DefinedContributionPlanEnd1,
		},
	}
	if ip.FilingStatus == Joint {
		r2 := retiree{
			age:        ip.Age2,
			ageAtStart: ip.Age2 + ip.PrePlanYears,
			throughAge: ip.PlanThroughAge2,
			mykey:      ip.MyKey2,
			definedContributionPlanStartAge: ip.DefinedContributionPlanStart2,
			definedContributionPlanEndAge:   ip.DefinedContributionPlanEnd2,
		}
		retirees = append(retirees, r2)
	}
	ms.Retirees = retirees
	//fmt.Fprintf(ms.Logfile, "Retirees: %#v\n", retirees)
	//fmt.Fprintf(ms.Logfile, "NewModelSpec: numacc: %d, accmap: %v\n", ms.ip.numacc, ms.ip.accmap)

	var err error
	var dbal float64
	const maxPossibleAccounts = 5
	if ip.TDRA1 > 0 || ip.TDRAContrib1 > 0 {
		a := account{}
		//fmt.Printf("TDRARate1: %v, RRate: %v\n", ip.TDRARate1, ip.RRate)
		a.RRate = ip.TDRARate1
		infr := 1.0
		if ip.TDRAContribInflate1 == true {
			infr = ip.IRate
		}
		a.acctype = IRA
		a.mykey = ip.MyKey1
		a.Origbal = float64(ip.TDRA1)
		a.Contrib = float64(ip.TDRAContrib1)
		a.Contributions, dbal, _, err = genContrib(ip.TDRAContrib1,
			ms.convertAge(ip.TDRAContribStart1, a.mykey),
			ms.convertAge(ip.TDRAContribEnd1, a.mykey),
			ip.StartPlan, ip.EndPlan, infr, a.RRate, ip.Age1)
		if err != nil {
			return nil, err
		}
		a.bal = a.Origbal*math.Pow(a.RRate, float64(ip.PrePlanYears)) + dbal
		ms.Accounttable = append(ms.Accounttable, a)
	}
	if ip.TDRA2 > 0 || ip.TDRAContrib2 > 0 {
		a := account{}
		a.RRate = ip.TDRARate2
		infr := 1.0
		if ip.TDRAContribInflate2 == true {
			infr = ip.IRate
		}
		a.acctype = IRA
		a.mykey = ip.MyKey2
		a.Origbal = float64(ip.TDRA2)
		a.Contrib = float64(ip.TDRAContrib2)
		a.Contributions, dbal, _, err = genContrib(ip.TDRAContrib2,
			ms.convertAge(ip.TDRAContribStart2, a.mykey),
			ms.convertAge(ip.TDRAContribEnd2, a.mykey),
			ip.StartPlan, ip.EndPlan, infr, a.RRate, ip.Age1)
		if err != nil {
			return nil, err
		}
		a.bal = a.Origbal*math.Pow(a.RRate, float64(ip.PrePlanYears)) + dbal
		ms.Accounttable = append(ms.Accounttable, a)
	}
	if ip.Roth1 > 0 || ip.RothContrib1 > 0 {
		a := account{}
		a.RRate = ip.RothRate1
		infr := 1.0
		if ip.RothContribInflate1 == true {
			infr = ip.IRate
		}
		a.acctype = Roth
		a.mykey = ip.MyKey1
		a.Origbal = float64(ip.Roth1)
		a.Contrib = float64(ip.RothContrib1)
		a.Contributions, dbal, _, err = genContrib(ip.RothContrib1,
			ms.convertAge(ip.RothContribStart1, a.mykey),
			ms.convertAge(ip.RothContribEnd1, a.mykey),
			ip.StartPlan, ip.EndPlan, infr, a.RRate, ip.Age1)
		if err != nil {
			return nil, err
		}
		a.bal = a.Origbal*math.Pow(a.RRate, float64(ip.PrePlanYears)) + dbal
		//fmt.Printf("Roth acc: %#v\n", a)
		ms.Accounttable = append(ms.Accounttable, a)
	}
	if ip.Roth2 > 0 || ip.RothContrib2 > 0 {
		a := account{}
		a.RRate = ip.RothRate2
		infr := 1.0
		if ip.RothContribInflate2 == true {
			infr = ip.IRate
		}
		a.acctype = Roth
		a.mykey = ip.MyKey2
		a.Origbal = float64(ip.Roth2)
		a.Contrib = float64(ip.RothContrib2)
		a.Contributions, dbal, _, err = genContrib(ip.RothContrib2,
			ms.convertAge(ip.RothContribStart2, a.mykey),
			ms.convertAge(ip.RothContribEnd2, a.mykey),
			ip.StartPlan, ip.EndPlan, infr, a.RRate, ip.Age1)
		if err != nil {
			return nil, err
		}
		a.bal = a.Origbal*math.Pow(a.RRate, float64(ip.PrePlanYears)) + dbal
		ms.Accounttable = append(ms.Accounttable, a)
	}
	if ip.Aftatax > 0 || ip.AftataxContrib > 0 {
		var dbasis float64
		a := account{}
		a.RRate = ip.AftataxRate
		infr := 1.0
		if ip.AftataxContribInflate == true {
			infr = ip.IRate
		}
		a.acctype = Aftertax
		a.mykey = "" // need to make this definable for pc versions
		a.Origbal = float64(ip.Aftatax)
		a.origbasis = float64(ip.AftataxBasis)
		a.Contrib = float64(ip.AftataxContrib)
		a.Contributions, dbal, dbasis, err = genContrib(ip.AftataxContrib,
			ms.convertAge(ip.AftataxContribStart, a.mykey),
			ms.convertAge(ip.AftataxContribEnd, a.mykey),
			ip.StartPlan, ip.EndPlan, infr, a.RRate, ip.Age1)
		if err != nil {
			return nil, err
		}
		a.bal = a.Origbal*math.Pow(a.RRate, float64(ip.PrePlanYears)) + dbal
		a.basis = a.origbasis + dbasis
		//fmt.Printf("aftertax accout: %#v\n", a)
		ms.Accounttable = append(ms.Accounttable, a)
	}
	if len(ms.Accounttable) != ms.Ip.Numacc {
		e := fmt.Errorf("NewModelSpecs: len(accounttable): %d not equal to numacc: %d", len(ms.Accounttable), ms.Ip.Numacc)
		return nil, e
	}
	ms.LiquidAssetPlanStart = 0.0
	for _, a := range ms.Accounttable {
		ms.LiquidAssetPlanStart += a.bal
	}

	ms.SS = make([][]float64, 0)
	SS, SS1, SS2, tags := processSS(&ip, wel)
	ms.SS = append(ms.SS, SS)
	ms.SS = append(ms.SS, SS1)
	ms.SS = append(ms.SS, SS2)
	ms.SStags = tags

	//fmt.Printf("SS1: %v\n", ms.SS1)
	//fmt.Printf("SS2: %v\n", ms.SS2)
	//fmt.Printf("SS: %v\n", ms.SS)

	ms.Income = make([][]float64, 1)
	ms.Incometags = append(ms.Incometags, "combined income")
	for i := 0; i < len(ip.Income); i++ {
		tag := ip.Income[i].Tag
		amount := ip.Income[i].Amount
		startage := ip.Income[i].StartAge
		endage := ip.Income[i].EndAge
		infr := ms.Ip.IRate
		if !ip.Income[i].Inflate {
			infr = 1.0
		}
		//fmt.Printf("tag: %s, amount: %d, start: %d, end %d, infr: %.3f, splan: %d, eplan: %d, age1: %d\n", tag, amount, startage, endage, infr, ip.StartPlan, ip.EndPlan, ip.Age1)
		income, err := buildVector(amount, startage, endage, ip.StartPlan, ip.EndPlan, infr, ip.Age1)
		if err != nil {
			return nil, err
		}
		//fmt.Printf("income: %#v\n", income)
		ms.Income[0], err = mergeVectors(ms.Income[0], income)
		if err != nil {
			return nil, err
		}
		if ip.Income[i].Tax {
			ms.Taxed, err = mergeVectors(ms.Taxed, income)
			if err != nil {
				return nil, err
			}
		}
		ms.Income = append(ms.Income, income)
		ms.Incometags = append(ms.Incometags, tag)
	}

	ms.Expenses = make([][]float64, 1)
	ms.Expensetags = append(ms.Expensetags, "combined expense")
	for i := 0; i < len(ip.Expense); i++ {
		tag := ip.Expense[i].Tag
		amount := ip.Expense[i].Amount
		startage := ip.Expense[i].StartAge
		endage := ip.Expense[i].EndAge
		infr := ms.Ip.IRate
		if !ip.Expense[i].Inflate {
			infr = 1.0
		}
		expense, err := buildVector(amount, startage, endage, ip.StartPlan, ip.EndPlan, infr, ip.Age1)
		if err != nil {
			return nil, err
		}
		ms.Expenses[0], err = mergeVectors(ms.Expenses[0], expense)
		if err != nil {
			return nil, err
		}
		ms.Expenses = append(ms.Expenses, expense)
		ms.Expensetags = append(ms.Expensetags, tag)
	}

	ms.AssetSale = make([][]float64, 1)
	ms.Assettags = append(ms.Assettags, "combined assets")
	ms.IlliquidAssetPlanStart = 0.0
	ms.IlliquidAssetPlanEnd = 0.0
	var noSell bool
	for i := 0; i < len(ip.Assets); i++ {
		noSell = false
		tag := ip.Assets[i].Tag
		value := float64(ip.Assets[i].Value)
		ageToSell := ip.Assets[i].AgeToSell
		if ageToSell < ip.StartPlan || ageToSell > ip.EndPlan {
			noSell = true
		}
		brokerageRate := ip.Assets[i].BrokeragePercent / 100.0
		if brokerageRate == 0 {
			brokerageRate = 0.04 // default to 4% // TODO FIXME defaults should be set in NewInputParams
			e := fmt.Errorf("Default (non-zero) value for BrokeragePercent should be set prior to this in NewInputSpecs()")
			panic(e)
		}
		assetRRate := ip.Assets[i].AssetRRate
		costAndImprovements := float64(ip.Assets[i].CostAndImprovements)
		owedAtAgeToSell := float64(ip.Assets[i].OwedAtAgeToSell)
		primaryResidence := ip.Assets[i].PrimaryResidence

		infr := ms.Ip.IRate // asset rate of return defaults to inflation rate
		if assetRRate != 0 {
			infr = assetRRate
		}
		//fmt.Printf("tag: %s, value: %.0f, ageToSell: %d, brokerageRate: %f, infr: %f\n", tag, value, ageToSell, brokerageRate, infr)
		//fmt.Printf("owedAtAgeToSell: %.0f\n", owedAtAgeToSell)
		if ageToSell < ip.StartPlan && ageToSell != 0 {
			err := fmt.Errorf("NewModelSpecs: Assets to be sold before plan start are not allow unless the age to sell is zero")
			return nil, err
		}
		ms.IlliquidAssetPlanStart += value * math.Pow(infr, float64(ip.StartPlan-ip.Age1))
		temp := 0.0
		if ageToSell > ip.EndPlan || ageToSell == 0 {
			// age after plan ends or zero cause value to remain
			temp = value * math.Pow(infr, float64(ip.StartPlan+ip.Numyr-ip.Age1))
		}
		ms.IlliquidAssetPlanEnd += temp

		assvec := make([]float64, ip.EndPlan-ip.StartPlan)
		tempvec := make([]float64, ip.EndPlan-ip.StartPlan)

		if !noSell {
			sellprice := value * math.Pow(infr, float64(ageToSell-ip.Age1))
			income := sellprice*(1-brokerageRate) - owedAtAgeToSell
			if income < 0 {
				income = 0
			}
			cgtaxable := sellprice*(1-brokerageRate) - costAndImprovements
			//fmt.Printf("Asset sell price $%.0f, brokerageRate: %%%d, income $%.0f, cgtaxable $%.0f\n", sellprice, int(100*brokerageRate), income, cgtaxable)
			if primaryResidence {
				cgtaxable -= ti.Primeresidence * math.Pow(ip.IRate, float64(ageToSell-ip.Age1))
				//fmt.Printf("cgtaxable: ", cgtaxable)
			}
			if cgtaxable < 0 {
				cgtaxable = 0
			}
			if income > 0 && ip.Accmap[Aftertax] <= 0 {
				e := fmt.Errorf("Error - Assets to be sold must have an 'aftertax' investment\naccount into which to deposit the net proceeds. Please\nadd an 'aftertax' account to yourn configuration; the bal may be zero")
				return nil, e
			}
			year := ageToSell - ip.StartPlan
			assvec[year] = income
			tempvec[year] = cgtaxable
		}
		ms.AssetSale[0], err = mergeVectors(ms.AssetSale[0], assvec)
		if err != nil {
			return nil, err
		}
		//if ip.Income[i].Tax {
		ms.CgAssetTaxed, err = mergeVectors(ms.CgAssetTaxed, tempvec)
		if err != nil {
			return nil, err
		}
		//}
		ms.AssetSale = append(ms.AssetSale, assvec)
		ms.Assettags = append(ms.Assettags, tag)
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
	if fourPercentRule {
		// override any setting of [max.income]
		// William Bengen the origniator of the 4% rule now says
		// its more like 4.5
		ms.Ip.Max =
			int(0.045 * (ms.LiquidAssetPlanStart + ms.IlliquidAssetPlanStart))
	}

	err = ms.verifyTaxableIncomeCoversContrib(wel)
	if err != nil {
		return nil, err
	}

	if ip.FilingStatus == Joint {
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

	nvars := ms.Vindx.Vsize
	A := make([][]float64, 0)
	b := make([]float64, 0)
	c := make([]float64, nvars)
	notes := make([]ModelNote, 0)

	//fmt.Printf("\nms.accounttable len: %d\n", len(ms.accounttable))

	//
	// Add objective function (S1') becomes (R1') if PlusEstate is added
	//
	for year := 0; year < ms.Ip.Numyr; year++ {
		c[ms.Vindx.S(year)] = -1
	}
	//
	// Add objective function tax bracket forcing function
	//
	for year := 0; year < ms.Ip.Numyr; year++ {
		for k := 0; k < len(*ms.Ti.Taxtable); k++ {
			// Multiplies the impact of higher brackets opposite to
			// optimization. The intent here is to pressure higher
			// brackets more and pack the lower brackets
			c[ms.Vindx.X(year, k)] = float64(k) / 10.0
		}
	}
	//
	// Adder objective function (R1') when PlusEstate is added
	//
	//fmt.Printf("ms.Ip.Maximize: %#v\n", ms.Ip.Maximize)
	if ms.Ip.Maximize == PlusEstate {
		for j := 0; j < len(ms.Accounttable); j++ {
			estateTax := ms.Ti.AccountEstateTax[ms.Accounttable[j].acctype]
			c[ms.Vindx.B(ms.Ip.Numyr, j)] = -1 * estateTax // account discount rate
		}
		//fmt.Fprintf(ms.Logfile, "\nConstructing Spending + Estate Model:\n")
		notes = append(notes, ModelNote{-1, "Objective function R1':"})
	} else {
		//fmt.Fprintf(ms.Logfile, "\nConstructing Spending Model:\n")

		startamount := 0.0
		for j := 0; j < len(ms.Accounttable); j++ {
			startamount += ms.Accounttable[j].bal
		}
		balancer := 1.0 / (startamount)
		for j := 0; j < len(ms.Accounttable); j++ {
			estateTax := ms.Ti.AccountEstateTax[ms.Accounttable[j].acctype]
			c[ms.Vindx.B(ms.Ip.Numyr, j)] = -1 * balancer * estateTax // balance and discount rate
		}
		notes = append(notes, ModelNote{-1, "Objective function S1':"})
	}
	//
	// Add constraint (2')
	//
	notes = append(notes, ModelNote{len(A), "Constraints 2':"})
	for year := 0; year < ms.Ip.Numyr; year++ {
		row := make([]float64, nvars)
		for j := 0; j < len(ms.Accounttable); j++ {
			p := 1.0
			if ms.Accounttable[j].acctype != Aftertax {
				if ms.Ti.applyEarlyPenalty(year, ms.matchRetiree(ms.Accounttable[j].mykey, year, true)) { // TODO: should applyEarlyPenalty() return the penalty amount, spimplifying things?
					p = 1 - ms.Ti.Penalty
				}
			}
			row[ms.Vindx.W(year, j)] = -1 * p
		}
		for k := 0; k < len(*ms.Ti.Taxtable); k++ {
			row[ms.Vindx.X(year, k)] = (*ms.Ti.Taxtable)[k][2] // income tax
		}
		if ms.Ip.Accmap[Aftertax] > 0 {
			for l := 0; l < len(*ms.Ti.Capgainstable); l++ {
				row[ms.Vindx.Y(year, l)] = (*ms.Ti.Capgainstable)[l][2] // cap gains tax
			}
		}
		for j := 0; j < len(ms.Accounttable); j++ {
			row[ms.Vindx.D(year, j)] = 1
		}
		row[ms.Vindx.S(year)] = 1
		A = append(A, row)
		inc := AccessVector(ms.Income[0], year)
		ss := AccessVector(ms.SS[0], year)
		exp := AccessVector(ms.Expenses[0], year)
		b = append(b, inc+ss-exp)
	}
	//
	// Add constraint (3a')
	//
	notes = append(notes, ModelNote{len(A), "Constraints 3a':"})
	for year := 0; year < ms.Ip.Numyr-1; year++ {
		row := make([]float64, nvars)
		row[ms.Vindx.S(year+1)] = 1
		row[ms.Vindx.S(year)] = -1 * ms.Ip.IRate
		A = append(A, row)
		b = append(b, 0)
	}
	//
	// Add constraint (3b')
	//
	notes = append(notes, ModelNote{len(A), "Constraints 3b':"})
	for year := 0; year < ms.Ip.Numyr-1; year++ {
		row := make([]float64, nvars)
		row[ms.Vindx.S(year)] = ms.Ip.IRate
		row[ms.Vindx.S(year+1)] = -1
		A = append(A, row)
		b = append(b, 0)
	}
	//
	// Add constrant (4') rows - not needed if [desired.income] is not defined in input
	//
	notes = append(notes, ModelNote{len(A), "Constraints 4':"})
	if ms.Ip.Min != 0 {
		for year := 0; year < 1; year++ { // Only needs setting at the beginning
			row := make([]float64, nvars)
			row[ms.Vindx.S(year)] = -1
			A = append(A, row)
			b = append(b, float64(-ms.Ip.Min)) // [- d_i]
		}
	}

	//
	// Add constraints for (5') rows - not added if [max.income] is
	// not defined in input
	//
	notes = append(notes, ModelNote{len(A), "Constraints 5':"})
	if ms.Ip.Max != 0 {
		for year := 0; year < 1; year++ { // Only needs to be set at the beginning
			row := make([]float64, nvars)
			row[ms.Vindx.S(year)] = 1
			A = append(A, row)
			b = append(b, float64(ms.Ip.Max)) // [ dm_i]
		}
	}

	//
	// Add constaints for (6') rows
	//
	notes = append(notes, ModelNote{len(A), "Constraints 6':"})
	for year := 0; year < ms.Ip.Numyr; year++ {
		row := make([]float64, nvars)
		for j := 0; j < len(ms.Accounttable); j++ {
			if ms.Accounttable[j].acctype != Aftertax {
				row[ms.Vindx.D(year, j)] = 1 // TODO if this is not executed, DONT register this constrain, DONT add to A and b
			}
		}
		A = append(A, row)
		//b+=[min(ms.Income[year],ms.Ti.maxContribution(year,None))]
		// using ms.Taxed rather than ms.Income because income could
		// include non-taxed anueities that don't count.
		None := ""
		infyears := ms.Ip.PrePlanYears + year
		b = append(b, math.Min(AccessVector(ms.Taxed, year), ms.Ti.maxContribution(year, infyears, ms.Retirees, None, ms.Ip.IRate)))
	}
	//
	// Add constaints for (7') rows
	//
	notes = append(notes, ModelNote{len(A), "Constraints 7':"})
	for year := 0; year < ms.Ip.Numyr; year++ {
		// TODO this is not needed when there is only one retiree
		infyears := ms.Ip.PrePlanYears + year
		for _, v := range ms.Retirees {
			row := make([]float64, nvars)
			for j := 0; j < len(ms.Accounttable); j++ {
				if v.mykey == ms.Accounttable[j].mykey {
					// ["acctype"] != "aftertax": no "mykey" in aftertax
					// (this will either break or just not match - we
					// will see)
					row[ms.Vindx.D(year, j)] = 1 // TODO if this is not executed, DONT register this constraint, DONT add to A and b
				}
			}
			A = append(A, row)
			b = append(b, ms.Ti.maxContribution(year, infyears, ms.Retirees, v.mykey, ms.Ip.IRate))
		}
	}
	//
	// Add constaints for (8') rows
	//
	notes = append(notes, ModelNote{len(A), "Constraints 8':"})
	for year := 0; year < ms.Ip.Numyr; year++ {
		for j := 0; j < len(ms.Accounttable); j++ {
			v := ms.Accounttable[j].Contributions
			if v != nil {
				if v[year] > 0 {
					row := make([]float64, nvars)
					row[ms.Vindx.D(year, j)] = -1
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
	for year := 0; year < ms.Ip.Numyr; year++ {
		for j := 0; j < intMin(2, len(ms.Accounttable)); j++ {
			// at most the first two accounts are type IRA w/
			// RMD requirement
			if ms.Accounttable[j].acctype == IRA {
				ownerage := ms.accountOwnerAge(year, ms.Accounttable[j])
				if ownerage >= 70 {
					row := make([]float64, nvars)
					row[ms.Vindx.D(year, j)] = 1
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
	if !ms.AllowTdraRothraDeposits {
		for year := 0; year < ms.Ip.Numyr; year++ {
			for j := 0; j < len(ms.Accounttable); j++ {
				v := ms.Accounttable[j].Contributions
				max := 0.0
				if v != nil {
					max = v[year]
				}
				if ms.Accounttable[j].acctype != Aftertax { //Todo: move this if statement up just under the for to remove all unnessasary work
					row := make([]float64, nvars)
					row[ms.Vindx.D(year, j)] = 1
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
	for year := 0; year < ms.Ip.Numyr; year++ {
		for j := 0; j < intMin(2, len(ms.Accounttable)); j++ {
			// at most the first two accounts are type IRA
			// w/ RMD requirement
			if ms.Accounttable[j].acctype == IRA {
				rmd := ms.Ti.rmdNeeded(year, ms.matchRetiree(ms.Accounttable[j].mykey, year, true))
				if rmd > 0 {
					row := make([]float64, nvars)
					row[ms.Vindx.B(year, j)] = 1 / rmd
					row[ms.Vindx.W(year, j)] = -1
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
	for year := 0; year < ms.Ip.Numyr; year++ {
		adjInf := math.Pow(ms.Ip.IRate, float64(ms.Ip.PrePlanYears+year))
		row := make([]float64, nvars)
		for j := 0; j < intMin(2, len(ms.Accounttable)); j++ {
			// IRA can only be in the first two accounts
			if ms.Accounttable[j].acctype == IRA {
				row[ms.Vindx.W(year, j)] = 1  // Account 0 is TDRA
				row[ms.Vindx.D(year, j)] = -1 // Account 0 is TDRA
			}
		}
		for k := 0; k < len(*ms.Ti.Taxtable); k++ {
			row[ms.Vindx.X(year, k)] = -1
		}
		A = append(A, row)
		b = append(b, ms.Ti.Stded*adjInf-AccessVector(ms.Taxed, year)-ms.Ti.SStaxable*AccessVector(ms.SS[0], year))
	}
	//
	// Add constraints for (12')
	//
	notes = append(notes, ModelNote{len(A), "Constraints 12':"})
	for year := 0; year < ms.Ip.Numyr; year++ {
		for k := 0; k < len(*ms.Ti.Taxtable)-1; k++ {
			row := make([]float64, nvars)
			row[ms.Vindx.X(year, k)] = 1
			A = append(A, row)
			b = append(b, ((*ms.Ti.Taxtable)[k][1])*math.Pow(ms.Ip.IRate, float64(ms.Ip.PrePlanYears+year))) // inflation adjusted
		}
	}
	//
	// Add constraints for (13a')
	//
	notes = append(notes, ModelNote{len(A), "Constraints 13a':"})
	if ms.Ip.Accmap[Aftertax] > 0 {
		for year := 0; year < ms.Ip.Numyr; year++ {
			f := ms.cgTaxableFraction(year)
			row := make([]float64, nvars)
			for l := 0; l < len(*ms.Ti.Capgainstable); l++ {
				row[ms.Vindx.Y(year, l)] = 1
			}
			// Awful Hack! If year of asset sale, assume w(i,j)-D(i,j) is
			// negative so taxable from this is zero
			cgt := AccessVector(ms.CgAssetTaxed, year)
			if cgt <= 0 { // i.e., no sale
				j := len(ms.Accounttable) - 1 // last Acc is investment / stocks
				row[ms.Vindx.W(year, j)] = -1 * f
				row[ms.Vindx.D(year, j)] = f
			}
			A = append(A, row)
			b = append(b, cgt)
		}
	}
	//
	// Add constraints for (13b')
	//
	notes = append(notes, ModelNote{len(A), "Constraints 13b':"})
	if ms.Ip.Accmap[Aftertax] > 0 {
		for year := 0; year < ms.Ip.Numyr; year++ {
			f := ms.cgTaxableFraction(year)
			row := make([]float64, nvars)
			////// Awful Hack! If year of asset sale, assume w(i,j)-D(i,j) is
			////// negative so taxable from this is zero
			cgt := AccessVector(ms.CgAssetTaxed, year)
			if cgt <= 0 { // i.e., no sale
				j := len(ms.Accounttable) - 1 // last Acc is investment / stocks
				row[ms.Vindx.W(year, j)] = f
				row[ms.Vindx.D(year, j)] = -f
			}
			for l := 0; l < len(*ms.Ti.Capgainstable); l++ {
				row[ms.Vindx.Y(year, l)] = -1
			}
			A = append(A, row)
			b = append(b, -1*cgt)
		}
	}
	//
	// Add constraints for (14')
	//
	notes = append(notes, ModelNote{len(A), "Constraints 14':"})
	if ms.Ip.Accmap[Aftertax] > 0 {
		for year := 0; year < ms.Ip.Numyr; year++ {
			adjInf := math.Pow(ms.Ip.IRate, float64(ms.Ip.PrePlanYears+year))
			for l := 0; l < len(*ms.Ti.Capgainstable)-1; l++ {
				row := make([]float64, nvars)
				row[ms.Vindx.Y(year, l)] = 1
				for k := 0; k < len(*ms.Ti.Taxtable)-1; k++ {
					if (*ms.Ti.Taxtable)[k][0] >= (*ms.Ti.Capgainstable)[l][0] 
						&& 
					  (*ms.Ti.Taxtable)[k][0] < (*ms.Ti.Capgainstable)[l+1][0] {
						row[ms.Vindx.X(year, k)] = 1
					}
				}
				A = append(A, row)
				b = append(b, (*ms.Ti.Capgainstable)[l][1]*adjInf) // mcg[i,l] inflation adjusted
			}
		}
	}
	//
	// Add constraints for (15a')
	//
	notes = append(notes, ModelNote{len(A), "Constraints 15a':"})
	for year := 0; year < ms.Ip.Numyr; year++ {
		for j := 0; j < len(ms.Accounttable); j++ {
			//j = len(ms.Accounttable)-1 // nl the last account, the investment account
			row := make([]float64, nvars)
			row[ms.Vindx.B(year+1, j)] = 1 // b[i,j] supports an extra year
			row[ms.Vindx.B(year, j)] = -1 * ms.Accounttable[j].RRate
			row[ms.Vindx.W(year, j)] = ms.Accounttable[j].RRate
			row[ms.Vindx.D(year, j)] = -1 * ms.Accounttable[j].RRate
			A = append(A, row)
			// In the event of a sell of an asset for the year
			temp := 0.0
			if ms.Accounttable[j].acctype == Aftertax {
				temp = AccessVector(ms.AssetSale[0], year) * ms.Accounttable[j].RRate //TODO test
			}
			b = append(b, temp)
		}
	}
	//
	// Add constraints for (15b')
	//
	notes = append(notes, ModelNote{len(A), "Constraints 15b':"})
	for year := 0; year < ms.Ip.Numyr; year++ {
		for j := 0; j < len(ms.Accounttable); j++ {
			//j = len(ms.Accounttable)-1 // nl the last account, the investment account
			row := make([]float64, nvars)
			row[ms.Vindx.B(year, j)] = ms.Accounttable[j].RRate
			row[ms.Vindx.W(year, j)] = -1 * ms.Accounttable[j].RRate
			row[ms.Vindx.D(year, j)] = ms.Accounttable[j].RRate
			row[ms.Vindx.B(year+1, j)] = -1 ////// b[i,j] supports an extra year
			A = append(A, row)
			temp := 0.0
			if ms.Accounttable[j].acctype == Aftertax {
				temp = -1 * AccessVector(ms.AssetSale[0], year) * ms.Accounttable[j].RRate //TODO test
			}
			b = append(b, temp)
		}
	}
	//
	// Constraint for (16a')
	//   Set the begining b[1,j] balances
	//
	notes = append(notes, ModelNote{len(A), "Constraints 16a':"})
	for j := 0; j < len(ms.Accounttable); j++ {
		row := make([]float64, nvars)
		row[ms.Vindx.B(0, j)] = 1
		A = append(A, row)
		b = append(b, ms.Accounttable[j].bal)
	}
	//
	// Constraint for (16b')
	//   Set the begining b[1,j] balances
	//
	notes = append(notes, ModelNote{len(A), "Constraints 16b':"})
	for j := 0; j < len(ms.Accounttable); j++ {
		row := make([]float64, nvars)
		row[ms.Vindx.B(0, j)] = -1
		A = append(A, row)
		b = append(b, -1*ms.Accounttable[j].bal)
	}
	//
	// Constrant for (17') is default for sycpy so no code is needed
	//
	notes = append(notes, ModelNote{len(A), "Constraints 17':"})

	return c, A, b, notes
}

// accountOwnerAge finds the age of the retiree who owns the account
// Only valid in plan years
func (ms ModelSpecs) accountOwnerAge(year int, acc account) int {
	age := 0
	retireekey := acc.mykey
	v := ms.matchRetiree(retireekey, year, true)
	if v != nil {
		age = v.ageAtStart + year
	}
	return age
}

// matchRetiree searches retirees by key returning nil if not found
func (ms ModelSpecs) matchRetiree(retireekey string, year int, livingOnly bool) *retiree {
	// Assumes only one retiree could be dead else passed end plan
	var ov *retiree
	ov = nil
	for _, v := range ms.Retirees {
		if v.mykey == retireekey {
			if v.throughAge-v.ageAtStart+1 > year {
				//fmt.Printf("matchRetire: looking for %s in year %d and returning %s\n", retireekey, year, v.mykey)
				return &v
			}
		} else {
			ov = &v
		}
	}
	/*
		if ov != nil {
			fmt.Printf("matchRetire: looking for %s in year %d and returning %s\n", retireekey, year, ov.mykey)
		} else {
			fmt.Printf("matchRetire: looking for %s in year %d and returning nil\n", retireekey, year)
		}
	*/
	return ov
}

// TODO unit test me :-)
// convertAge converts an age for key1 to an age in the primary timeline
func (ms ModelSpecs) convertAge(age int, key string) int {
	index := -1
	for i, v := range ms.Retirees {
		if v.mykey == key {
			index = i
		}
	}
	if index <= 0 {
		return age
	}
	//delta := ms.Retirees[0].age - ms.Retirees[1].age
	return age + ms.Ip.AgeDelta
}

// cgTaxableFraction estimates the portion of capital gains not from basis
func (ms ModelSpecs) cgTaxableFraction(year int) float64 {
	// applies only in Plan years
	f := 1.0
	if ms.Ip.Accmap[Aftertax] > 0 {
		//TODO: FIXME REMOVE THIS LOOP
		for _, v := range ms.Accounttable {
			if v.acctype == Aftertax {
				if v.bal > 0 { // don't want to divide by zero
					//
					// v.bal includes the rRate and v.basis includes
					// the additional contributions up until
					// startPlan so no need to inflate for ms.Ip.PrePlanYears
					//
					f = 1 - (v.basis / (v.bal * math.Pow(v.RRate, float64(year))))
				}
				break // should be the last entry anyway but...
			}
		}
	}
	return f
}

// PrecheckConsistancy is, I think, checked elsewhere; delete?
func (ms ModelSpecs) PrecheckConsistancy() bool {
	fmt.Printf("\nDoing Pre-check:")
	// check that there is income for all contibutions
	//    #tcontribs = 0
	for year := 0; year < ms.Ip.Numyr; year++ {
		t := 0.0
		for j := 0; j < len(ms.Accounttable); j++ {
			v := ms.Accounttable[j]
			if v.acctype != Aftertax {
				if v.Contributions != nil && len(v.Contributions) > 0 {
					t += v.Contributions[year]
				}
			}
		}
		if t > AccessVector(ms.Taxed, year) { // was S.income[year]
			fmt.Printf("year: %d, total contributions of (%.0f) to all Retirement accounts exceeds other earned (i.e., taxable) income (%.0f)",
				year, t, AccessVector(ms.Taxed, year))
			// was S.income[year]
			fmt.Printf("Please change the contributions in the toml file to be less than non-SS income.")
			os.Exit(1) //TODO FIXME no exit allowed pass the error back!!!
		}
	}
	return true
}

/*
def consistancy_check(res, years, taxbins, cgbins, accounts, accmap, vindx):
    # check to see if the ordinary tax brackets are filled in properly
    print()
    print()
    print("Consistancy Checking:")
    print()

    result = vvar.my_check_index_sequence(years, taxbins, cgbins, accounts, accmap, vindx)

    for year in range(S.numyr):
        s = 0
        fz = False
        fnf = False
        i_mul = S.i_rate ** (S.preplanyears+year)
        for k in range(len(taxinfo.taxtable)):
            cut, size, rate, base = taxinfo.taxtable[k]
            size *= i_mul
            s += res.x[vindx.x(year,k)]
            if fnf and res.x[vindx.x(year,k)] > 0:
                print("Inproper packed brackets in year %d, bracket %d not empty while previous bracket not full." % (year, k))
            if res.x[vindx.x(year,k)]+1 < size:
                fnf = True
            if fz and res.x[vindx.x(year,k)] > 0:
                print("Inproperly packed tax brackets in year %d bracket %d" % (year, k))
            if res.x[vindx.x(year,k)] == 0.0:
                fz = True
        if S.accmap['aftertax'] > 0:
            scg = 0
            fz = False
            fnf = False
            for l in range(len(taxinfo.capgainstable)):
                cut, size, rate = taxinfo.capgainstable[l]
                size *= i_mul
                bamount = res.x[vindx.y(year,l)]
                scg += bamount
                for k in range(len(taxinfo.taxtable)-1):
                    if taxinfo.taxtable[k][0] >= taxinfo.capgainstable[l][0] and taxinfo.taxtable[k][0] < taxinfo.capgainstable[l+1][0]:
                        bamount += res.x[vindx.x(year,k)]
                if fnf and bamount > 0:
                    print("Inproper packed CG brackets in year %d, bracket %d not empty while previous bracket not full." % (year, l))
                if bamount+1 < size:
                    fnf = True
                if fz and bamount > 0:
                    print("Inproperly packed GC tax brackets in year %d bracket %d" % (year, l))
                if bamount == 0.0:
                    fz = True
        TaxableOrdinary = OrdinaryTaxable(year)
        if (TaxableOrdinary + 0.1 < s) or (TaxableOrdinary - 0.1 > s):
            print("Error: Expected (age:%d) Taxable Ordinary income %6.2f doesn't match bracket sum %6.2f" %
                (year + S.startage, TaxableOrdinary,s))

        for j in range(len(S.accounttable)):
            a = res.x[vindx.b(year+1,j)] - (res.x[vindx.b(year,j)] - res.x[vindx.w(year,j)] + deposit_amount(S, res, year, j))*S.accounttable[j]['rate']
            if a > 1:
                v = S.accounttable[j]
                print("account[%d], type %s, index %d, mykey %s" % (j, v['acctype'], v['index'], v['mykey']))
                print("account[%d] year to year balance NOT OK years %d to %d" % (j, year, year+1))
                print("difference is", a)

        T,spendable,tax,rate,cg_tax,earlytax,rothearly = IncomeSummary(year)
        if spendable + 0.1 < res.x[vindx.s(year)]  or spendable -0.1 > res.x[vindx.s(year)]:
            print("Calc Spendable %6.2f should equal s(year:%d) %6.2f"% (spendable, year, res.x[vindx.s(year)]))
            for j in range(len(S.accounttable)):
                print("+w[%d,%d]: %6.0f" % (year, j, res.x[vindx.w(year,j)]))
                print("-D[%d,%d]: %6.0f" % (year, j, deposit_amount(S, res, year, j)))
            print("+o[%d]: %6.0f +SS[%d]: %6.0f -tax: %6.0f -cg_tax: %6.0f" % (year, S.income[year] ,year, S.SS[year] , tax ,cg_tax))

        bt = 0
        for k in range(len(taxinfo.taxtable)):
            bt += res.x[vindx.x(year,k)] * taxinfo.taxtable[k][2]
        if tax + 0.1 < bt  or tax -0.1 > bt:
            print("Calc tax %6.2f should equal brackettax(bt)[]: %6.2f" % (tax, bt))
    print()
*/

// TODO: FIXME: Create UNIT tests: last two parameters need s vector (s is output from simplex run)

// PrintModelMatrix prints to object function (cx) and constraint matrix (Ax<=b)
func (ms ModelSpecs) PrintModelMatrix(c []float64, A [][]float64, b []float64, notes []ModelNote, s []float64, nonBindingOnly bool) {
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
		fmt.Fprintf(ms.Logfile, "\n##== [%d-%d]: %s ==##\n", from, to, note)
		note = notes[notesIndex].note
		notesIndex++
	}
	fmt.Fprintf(ms.Logfile, "c: ")
	ms.printModelRow(c, false)
	fmt.Fprintf(ms.Logfile, "\n")
	if !nonBindingOnly {
		fmt.Fprintf(ms.Logfile, "B?  i: A_ub[i]: b[i]\n")
		for constraint := 0; constraint < len(A); constraint++ {
			if nextModelIndex == constraint {
				from := nextModelIndex
				nextModelIndex = notes[notesIndex].index
				to := nextModelIndex - 1
				for to < from {
					fmt.Fprintf(ms.Logfile, "\n##== [%d-%d]: %s ==##\n", from, to, note)
					note = notes[notesIndex].note
					notesIndex++
					from = nextModelIndex
					nextModelIndex = notes[notesIndex].index
					to = nextModelIndex - 1
				}
				fmt.Fprintf(ms.Logfile, "\n##== [%d-%d]: %s ==##\n", from, to, note)
				note = notes[notesIndex].note
				notesIndex++
			}
			if s == nil || s[constraint] > 0 {
				fmt.Fprintf(ms.Logfile, "  ")
			} else {
				fmt.Fprintf(ms.Logfile, "B ")
			}
			fmt.Fprintf(ms.Logfile, "%3d: ", constraint)
			ms.printConstraint(A[constraint], b[constraint])
		}
	} else {
		fmt.Fprintf(ms.Logfile, "  i: A_ub[i]: b[i]\n")
		j := 0
		for constraint := 0; constraint < len(A); constraint++ {
			if nextModelIndex == constraint {
				from := nextModelIndex
				nextModelIndex = notes[notesIndex].index
				to := nextModelIndex - 1
				for to < from {
					fmt.Fprintf(ms.Logfile, "\n##== [%d-%d]: %s ==##\n", from, to, note)
					note = notes[notesIndex].note
					notesIndex++
					from = nextModelIndex
					nextModelIndex = notes[notesIndex].index
					to = nextModelIndex - 1
				}
				fmt.Fprintf(ms.Logfile, "\n##== [%d-%d]: %s ==##\n", from, to, note)
				note = notes[notesIndex].note
				notesIndex++
			}
			if s[constraint] > 0 {
				j++
				fmt.Fprintf(ms.Logfile, "%3d: ", constraint)
				ms.printConstraint(A[constraint], b[constraint])
			}
		}
		fmt.Fprintf(ms.Logfile, "\n\n%d non-binding constrains printed\n", j)
	}
	fmt.Fprintf(ms.Logfile, "\n")
}

func (ms ModelSpecs) printConstraint(row []float64, b float64) {
	ms.printModelRow(row, true)
	fmt.Fprintf(ms.Logfile, "<= b[]: %6.2f\n", b)
}

func (ms ModelSpecs) printModelRow(row []float64, suppressNewline bool) {
	if ms.Ip.Numacc < 0 || ms.Ip.Numacc > 5 {
		e := fmt.Errorf("PrintModelRow: number of accounts is out of bounds, should be between [0, 5] but is %d", ms.Ip.Numacc)
		fmt.Fprintf(ms.Logfile, "%s\n", e)
		return
	}
	for i := 0; i < ms.Ip.Numyr; i++ { // x[]
		for k := 0; k < len(*ms.Ti.Taxtable); k++ {
			if row[ms.Vindx.X(i, k)] != 0 {
				fmt.Fprintf(ms.Logfile, "x[%d,%d]=%6.3f, ", i, k, row[ms.Vindx.X(i, k)])
			}
		}
	}
	if ms.Ip.Accmap[Aftertax] > 0 {
		for i := 0; i < ms.Ip.Numyr; i++ { // y[]
			for l := 0; l < len(*ms.Ti.Capgainstable); l++ {
				if row[ms.Vindx.Y(i, l)] != 0 {
					fmt.Fprintf(ms.Logfile, "y[%d,%d]=%6.3f, ", i, l, row[ms.Vindx.Y(i, l)])
				}
			}
		}
	}
	for i := 0; i < ms.Ip.Numyr; i++ { // w[]
		for j := 0; j < ms.Ip.Numacc; j++ {
			if row[ms.Vindx.W(i, j)] != 0 {
				fmt.Fprintf(ms.Logfile, "w[%d,%d]=%6.3f, ", i, j, row[ms.Vindx.W(i, j)])
			}
		}
	}
	for i := 0; i < ms.Ip.Numyr+1; i++ { // b[] has an extra year
		for j := 0; j < ms.Ip.Numacc; j++ {
			if row[ms.Vindx.B(i, j)] != 0 {
				fmt.Fprintf(ms.Logfile, "b[%d,%d]=%6.3f, ", i, j, row[ms.Vindx.B(i, j)])
			}
		}
	}
	for i := 0; i < ms.Ip.Numyr; i++ { // s[]
		if row[ms.Vindx.S(i)] != 0 {
			fmt.Fprintf(ms.Logfile, "s[%d]=%6.3f, ", i, row[ms.Vindx.S(i)])
		}
	}
	for i := 0; i < ms.Ip.Numyr; i++ { // D[]
		for j := 0; j < ms.Ip.Numacc; j++ {
			if row[ms.Vindx.D(i, j)] != 0 {
				fmt.Fprintf(ms.Logfile, "D[%d,%d]=%6.3f, ", i, j, row[ms.Vindx.D(i, j)])
			}
		}
	}
	if !suppressNewline {
		fmt.Fprintf(ms.Logfile, "\n")
	}
}
