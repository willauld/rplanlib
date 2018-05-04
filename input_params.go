package rplanlib

import (
	"fmt"
	"strconv"
	"strings"
)

// MaxStreams is the maximam number of streams for each: income, expense, asset
const MaxStreams = 10

// MaxWhat is what is to be maximized while creating a plan (Spending, PlusEstate)
type MaxWhat int

const (
	UnknownMaxWhat MaxWhat = iota
	Spending       MaxWhat = iota
	PlusEstate     MaxWhat = iota
)

func (a MaxWhat) String() string {
	names := [...]string{
		"Unknown",
		"Spending",
		"PlusEstate",
	}
	if a < UnknownMaxWhat || a > PlusEstate {
		return "Unknown"
	}
	return names[a]
}

// TaxStatus is tax type for the plan to be created
type TaxStatus int

const (
	UnknownTaxStatus TaxStatus = iota
	Joint            TaxStatus = iota
	Mseparate        TaxStatus = iota
	Single           TaxStatus = iota
)

func (a TaxStatus) String() string {
	names := [...]string{
		"Unknown",
		"Joint",
		"Mseparate",
		"Single",
	}
	if a < UnknownTaxStatus || a > Single {
		return "Unknown"
	}
	return names[a]
}

type Acctype int

const (
	IRA            Acctype = iota
	Roth           Acctype = iota
	Aftertax       Acctype = iota
	UnknownAcctype Acctype = iota
)

func (a Acctype) String() string {
	names := [...]string{
		"IRA",
		"Roth",
		"Aftertax",
		"UnknownAcctype",
	}
	if a < IRA || a > UnknownAcctype {
		return "Unknown"
	}
	return names[a]
}

// InputParams are the model params constructed from driver program string input
type InputParams struct {
	FilingStatus                  TaxStatus
	MyKey1                        string
	MyKey2                        string
	Age1                          int
	Age2                          int
	RetireAge1                    int
	RetireAge2                    int
	PlanThroughAge1               int
	PlanThroughAge2               int
	DefinedContributionPlanStart1 int
	DefinedContributionPlanStart2 int
	DefinedContributionPlanEnd1   int
	DefinedContributionPlanEnd2   int
	PIA1                          int
	PIA2                          int
	SSStart1                      int
	SSStart2                      int
	TDRA1                         int
	TDRA2                         int
	TDRARate1                     float64
	TDRARate2                     float64
	TDRAContrib1                  int
	TDRAContrib2                  int
	TDRAContribStart1             int
	TDRAContribStart2             int
	TDRAContribEnd1               int
	TDRAContribEnd2               int
	TDRAContribInflate1           bool
	TDRAContribInflate2           bool
	Roth1                         int
	Roth2                         int
	RothRate1                     float64
	RothRate2                     float64
	RothContrib1                  int
	RothContrib2                  int
	RothContribStart1             int
	RothContribStart2             int
	RothContribEnd1               int
	RothContribEnd2               int
	RothContribInflate1           bool
	RothContribInflate2           bool
	Aftatax                       int
	AftataxBasis                  int // TODO Add to Mobile
	AftataxRate                   float64
	AftataxContrib                int
	AftataxContribStart           int
	AftataxContribEnd             int
	AftataxContribInflate         bool
	IRatePercent                  float64 // TODO add to Mobile
	IRate                         float64 // local, not mobile
	RRatePercent                  float64 // TODO add to Mobile
	RRate                         float64 // local, not mobile
	Maximize                      MaxWhat // "Spending" or "PlusEstate" // TODO add to Mobile
	Min                           int     // TODO add to Mobile
	Max                           int     // TODO add to Mobile

	PrePlanYears int
	StartPlan    int
	EndPlan      int
	AgeDelta     int
	Numyr        int
	Accmap       map[Acctype]int
	Numacc       int

	Income  []stream
	Expense []stream
	Assets  []asset
}

type stream struct {
	Tag      string
	Amount   int
	StartAge int
	EndAge   int
	Inflate  bool
	Tax      bool
}
type asset struct {
	Tag                 string
	Value               int     // current value of the asset
	CostAndImprovements int     // purchase price plus improvment cost
	AgeToSell           int     // age at which to sell the asset
	OwedAtAgeToSell     int     // amount owed at time of sell (ageToSell)
	PrimaryResidence    bool    // Primary residence gets tax break
	AssetRRatePercent   float64 // avg rate of return (defaults to global rate)
	AssetRRate          float64 // avg rate of return (defaults to global rate)
	BrokeragePercent    float64 // avg rate paid for brokerage services
}

/*/TODO: REMOVE ME
func kgetIPIntValue(str string) int {
	return 1000 * getIPIntValue(str)
}
*/

func getIPIntValue(str string) int {
	if str == "" {
		return 0
	}
	n, e := strconv.Atoi(str)
	if e != nil {
		//fmt.Printf("GetIPIntValue(): %s\n", e)
		panic(e)
	}
	return n
}

func getIPFloatValue(str string, notSetVal float64) float64 {
	if str == "" {
		return notSetVal
	}
	n, e := strconv.ParseFloat(str, 64)
	if e != nil {
		//fmt.Printf("GetIPFloatValue(): %s\n", e)
		panic(e)
	}
	return n
}

//TODO: TESTME
func getIPBoolValue(str string) bool {
	if str == "" {
		return false
	}
	l := strings.ToLower(strings.TrimSpace(str))
	b, e := strconv.ParseBool(l)
	if e != nil {
		if l == "yes" {
			return true
		}
		if l == "no" {
			return false
		}
		//fmt.Printf("GetIPFloatValue(): %s\n", e)
		//TODO FIXME get ride of panic for production
		panic(e)
	}
	return b
}
func verifyMaximize(s string) (MaxWhat, error) {
	if s == "Spending" {
		return Spending, nil
	}
	if s == "PlusEstate" {
		return PlusEstate, nil
	}
	e := fmt.Errorf("verifyMaximize: %s is not a valid option", s)
	return UnknownMaxWhat, e
}

func verifyFilingStatus(s string) (TaxStatus, error) {
	if s == "joint" {
		return Joint, nil
	}
	if s == "mseparate" {
		return Mseparate, nil
	}
	if s == "single" {
		return Single, nil
	}
	e := fmt.Errorf("verifyFilingStatus: '%s' is not a valid option", s)
	return UnknownTaxStatus, e
}

// Default values is not defined
const ReturnRatePercent = 6.0
const InflactionRatePercent = 2.5
const MaximizeDefault = Spending
const FilingStatusDefault = Joint
const BrokeragePercentDefault = 4.0
const InflateContribDefault = false

// NewInputParams takes string inputs and converts them to model inputs
// assigning default values where needed
func NewInputParams(ip map[string]string, warnList *WarnErrorList) (*InputParams, error) {

	var err error

	dollarsInThousands, ok := ip["dollarsInThousands"]
	if !ok {
		e := fmt.Errorf("NewInputParams: key: 'dollarsInThousands' not defined in API map")
		return nil, e
	}
	multiplier := 1
	if dollarsInThousands == "true" {
		multiplier = 1000
	}

	rip := InputParams{}

	// Add this first so they are available for error msg
	rip.MyKey1 = ip["key1"]
	if rip.MyKey1 == "" {
		rip.MyKey1 = "Retiree1"
	}
	rip.MyKey2 = ip["key2"]

	rip.RRatePercent = getIPFloatValue(ip["eT_rRatePercent"], -1.0)
	if rip.RRatePercent < 0 { // Changing so user can assign 0
		rip.RRatePercent = ReturnRatePercent
	}
	rip.IRatePercent = getIPFloatValue(ip["eT_iRatePercent"], -1.0)
	if rip.IRatePercent < 0 { // Changing so user can assign 0
		rip.IRatePercent = InflactionRatePercent
	}
	rip.RRate = 1 + rip.RRatePercent/100.0
	rip.IRate = 1 + rip.IRatePercent/100.0
	rip.Min = getIPIntValue(ip["eT_DesiredIncome"]) // TODO add to mobile
	rip.Max = getIPIntValue(ip["eT_MaxIncome"])     // TODO add to mobile
	//	maximize:                "Spending", // or "PlusEstate"
	rip.Maximize = MaximizeDefault
	if ip["eT_maximize"] != "" {
		rip.Maximize, err = verifyMaximize(ip["eT_maximize"])
		if err != nil {
			return nil, err
		}
	}

	rip.Accmap = map[Acctype]int{IRA: 0, Roth: 0, Aftertax: 0}

	rip.FilingStatus = FilingStatusDefault
	if ip["filingStatus"] != "" {
		rip.FilingStatus, err = verifyFilingStatus(ip["filingStatus"])
		if err != nil {
			return nil, err
		}
	}

	//rip.MyKey1 = ip["key1"]
	if ip["eT_Age1"] == "" ||
		ip["eT_RetireAge1"] == "" ||
		ip["eT_PlanThroughAge1"] == "" {
		//if rip.MyKey1 == "" {
		//	rip.MyKey1 = "retiree1"
		//}
		e := fmt.Errorf("NewInputParams: retiree '%s' age, retirement age and plan through age must all be specified", rip.MyKey1)
		return nil, e
	}
	rip.Age1 = getIPIntValue(ip["eT_Age1"])
	rip.RetireAge1 = getIPIntValue(ip["eT_RetireAge1"])
	if rip.RetireAge1 < rip.Age1 {
		rip.RetireAge1 = rip.Age1
	}
	rip.PlanThroughAge1 = getIPIntValue(ip["eT_PlanThroughAge1"])
	yearsToRetire1 := rip.RetireAge1 - rip.Age1
	through1 := rip.PlanThroughAge1 - rip.Age1

	if ip["eT_DefinedContributionPlanStart1"] != "" ||
		ip["eT_DefinedContributionPlanEnd1"] != "" {
		if ip["eT_DefinedContributionPlanStart1"] == "" ||
			ip["eT_DefinedContributionPlanEnd1"] == "" {
			e := fmt.Errorf("NewInputParams: retiree '%s' defined contribution plan start and end ages must be specified", rip.MyKey1)
			return nil, e
		}
	}
	rip.DefinedContributionPlanStart1 = getIPIntValue(ip["eT_DefinedContributionPlanStart1"])
	rip.DefinedContributionPlanEnd1 = getIPIntValue(ip["eT_DefinedContributionPlanEnd1"])

	if rip.DefinedContributionPlanStart1 > rip.DefinedContributionPlanEnd1 {
		// An error TODO FIXME
	}

	if ip["eT_PIA1"] != "" || ip["eT_SS_Start1"] != "" {
		if ip["eT_PIA1"] == "" || ip["eT_SS_Start1"] == "" {
			e := fmt.Errorf("NewInputParams: retiree '%s' social security PIA and start age both must be specified", rip.MyKey1)
			return nil, e
		}
	}
	rip.PIA1 = getIPIntValue(ip["eT_PIA1"]) * multiplier
	rip.SSStart1 = getIPIntValue(ip["eT_SS_Start1"])

	//fmt.Printf("contrib str: '%s', val: %d\n", ip["eT_TDRA_Contrib1"], getIPIntValue(ip["eT_TDRA_Contrib1"]))
	if (ip["eT_TDRA_Contrib1"] != "" && getIPIntValue(ip["eT_TDRA_Contrib1"]) != 0) ||
		ip["eT_TDRA_ContribStartAge1"] != "" ||
		ip["eT_TDRA_ContribEndAge1"] != "" {
		if ip["eT_TDRA_Contrib1"] == "" ||
			ip["eT_TDRA_ContribStartAge1"] == "" ||
			ip["eT_TDRA_ContribEndAge1"] == "" {
			e := fmt.Errorf("NewInputParams: retiree '%s' TDRA contribution requires contribution amount, start and end age for contributions be specified", rip.MyKey1)
			return nil, e
		}
	}
	rip.TDRA1 = getIPIntValue(ip["eT_TDRA1"]) * multiplier
	rip.TDRARate1 = getIPFloatValue(ip["eT_TDRA_Rate1"], -1.0)
	rip.TDRAContrib1 = getIPIntValue(ip["eT_TDRA_Contrib1"]) * multiplier
	rip.TDRAContribStart1 = getIPIntValue(ip["eT_TDRA_ContribStartAge1"])
	rip.TDRAContribEnd1 = getIPIntValue(ip["eT_TDRA_ContribEndAge1"])
	rip.TDRAContribInflate1 = getIPBoolValue(ip["eT_TDRA_ContribInflate1"])
	if rip.TDRA1 > 0 || rip.TDRAContrib1 > 0 {
		rip.Accmap[IRA]++
	}
	if rip.TDRARate1 < 0 {
		rip.TDRARate1 = rip.RRate // TODO FIXME this is a percent not adj factor
	} else {
		rip.TDRARate1 = 1 + (rip.TDRARate1 / 100.0)
	}

	if (ip["eT_Roth_Contrib1"] != "" && getIPIntValue(ip["eT_Roth_Contrib1"]) != 0) ||
		ip["eT_Roth_ContribStartAge1"] != "" ||
		ip["eT_Roth_ContribEndAge1"] != "" {
		if ip["eT_Roth_Contrib1"] == "" ||
			ip["eT_Roth_ContribStartAge1"] == "" ||
			ip["eT_Roth_ContribEndAge1"] == "" {
			e := fmt.Errorf("NewInputParams: retiree '%s' Roth contribution requires contribution amount, start and end age for contributions be specified", rip.MyKey1)
			return nil, e
		}
	}
	rip.Roth1 = getIPIntValue(ip["eT_Roth1"]) * multiplier
	rip.RothRate1 = getIPFloatValue(ip["eT_Roth_Rate1"], -1.0)
	rip.RothContrib1 = getIPIntValue(ip["eT_Roth_Contrib1"]) * multiplier
	rip.RothContribStart1 = getIPIntValue(ip["eT_Roth_ContribStartAge1"])
	rip.RothContribEnd1 = getIPIntValue(ip["eT_Roth_ContribEndAge1"])
	rip.RothContribInflate1 = getIPBoolValue(ip["eT_Roth_ContribInflate1"])
	if rip.Roth1 > 0 || rip.RothContrib1 > 0 {
		rip.Accmap[Roth]++
	}
	if rip.RothRate1 < 0 {
		rip.RothRate1 = rip.RRate
	} else {
		rip.RothRate1 = 1 + (rip.RothRate1 / 100.0)
	}

	var through2 int
	var yearsToRetire2 int
	needRetiree2 := false
	if rip.FilingStatus == Joint {

		if ip["eT_DefinedContributionPlanStart2"] != "" ||
			ip["eT_DefinedContributionPlanEnd2"] != "" {
			if ip["eT_DefinedContributionPlanStart2"] == "" ||
				ip["eT_DefinedContributionPlanEnd2"] == "" {
				e := fmt.Errorf("NewInputParams: retiree '%s' defined contribution plan start and end ages must be specified", rip.MyKey2)
				return nil, e
			}
			needRetiree2 = true
		}
		rip.DefinedContributionPlanStart2 = getIPIntValue(ip["eT_DefinedContributionPlanStart2"])
		rip.DefinedContributionPlanEnd2 = getIPIntValue(ip["eT_DefinedContributionPlanEnd2"])

		if rip.DefinedContributionPlanStart2 > rip.DefinedContributionPlanEnd2 {
			// An error TODO FIXME
		}

		if ip["eT_PIA2"] != "" || ip["eT_SS_Start2"] != "" {
			if ip["eT_PIA2"] == "" || ip["eT_SS_Start2"] == "" {
				e := fmt.Errorf("NewInputParams: retiree '%s' social security PIA and start age both must be specified", rip.MyKey2)
				return nil, e
			}
			needRetiree2 = true
		}
		if ip["eT_PIA2"] != "" || ip["eT_PIA1"] != "" {
			if ip["eT_PIA2"] == "" || ip["eT_PIA1"] == "" {
				// if any SS set both must be specified
				str := fmt.Sprintf("Warning - Both retiree social security PIA (-1 or 0 for spousal benefits) and start age should be specified or you may be leaving money on the table")
				warnList.AppendWarning(str)
			}
		}
		rip.PIA2 = getIPIntValue(ip["eT_PIA2"]) * multiplier
		rip.SSStart2 = getIPIntValue(ip["eT_SS_Start2"])

		if (ip["eT_TDRA_Contrib2"] != "" && getIPIntValue(ip["eT_TDRA_Contrib2"]) != 0) ||
			ip["eT_TDRA_ContribStartAge2"] != "" ||
			ip["eT_TDRA_ContribEndAge2"] != "" {
			if ip["eT_TDRA_Contrib2"] == "" ||
				ip["eT_TDRA_ContribStartAge2"] == "" ||
				ip["eT_TDRA_ContribEndAge2"] == "" {
				e := fmt.Errorf("NewInputParams: retiree '%s' TDRA contribution requires contribution amount, start and end age for contributions be specified", rip.MyKey2) //TODO should TDRA be changed to IRA for should Toml IRA be changed to TDRA
				return nil, e
			}
			needRetiree2 = true
		}
		rip.TDRA2 = getIPIntValue(ip["eT_TDRA2"]) * multiplier
		rip.TDRARate2 = getIPFloatValue(ip["eT_TDRA_Rate2"], -1.0)
		rip.TDRAContrib2 = getIPIntValue(ip["eT_TDRA_Contrib2"]) * multiplier
		rip.TDRAContribStart2 = getIPIntValue(ip["eT_TDRA_ContribStartAge2"])
		rip.TDRAContribEnd2 = getIPIntValue(ip["eT_TDRA_ContribEndAge2"])
		rip.TDRAContribInflate2 = getIPBoolValue(ip["eT_TDRA_ContribInflate2"])
		if rip.TDRA2 > 0 || rip.TDRAContrib2 > 0 {
			rip.Accmap[IRA]++
			needRetiree2 = true
		}
		if rip.TDRARate2 < 0 {
			rip.TDRARate2 = rip.RRate
		} else {
			rip.TDRARate2 = 1 + (rip.TDRARate2 / 100.0)
		}

		if (ip["eT_Roth_Contrib2"] != "" && getIPIntValue(ip["eT_Roth_Contrib2"]) != 0) ||
			ip["eT_Roth_ContribStartAge2"] != "" ||
			ip["eT_Roth_ContribEndAge2"] != "" {
			if ip["eT_Roth_Contrib2"] == "" ||
				ip["eT_Roth_ContribStartAge2"] == "" ||
				ip["eT_Roth_ContribEndAge2"] == "" {
				e := fmt.Errorf("NewInputParams: retiree '%s' Roth contribution requires contribution amount, start and end age for contributions be specified", rip.MyKey2)
				return nil, e
			}
			needRetiree2 = true
		}
		rip.Roth2 = getIPIntValue(ip["eT_Roth2"]) * multiplier
		rip.RothRate2 = getIPFloatValue(ip["eT_Roth_Rate2"], -1.0)
		rip.RothContrib2 = getIPIntValue(ip["eT_Roth_Contrib2"]) * multiplier
		rip.RothContribStart2 = getIPIntValue(ip["eT_Roth_ContribStartAge2"])
		rip.RothContribEnd2 = getIPIntValue(ip["eT_Roth_ContribEndAge2"])
		rip.RothContribInflate2 = getIPBoolValue(ip["eT_Roth_ContribInflate2"])
		if rip.Roth2 > 0 || rip.RothContrib2 > 0 {
			rip.Accmap[Roth]++
			needRetiree2 = true
		}
		if rip.RothRate2 < 0 {
			rip.RothRate2 = rip.RRate
		} else {
			rip.RothRate2 = 1 + (rip.RothRate2 / 100.0)
		}
		//rip.MyKey2 = ip["key2"]
		if needRetiree2 ||
			(ip["eT_Age2"] != "" ||
				ip["eT_RetireAge2"] != "" ||
				ip["eT_PlanThroughAge2"] != "") {
			if ip["eT_Age2"] == "" ||
				ip["eT_RetireAge2"] == "" ||
				ip["eT_PlanThroughAge2"] == "" {
				//if rip.MyKey2 == "" {
				//	rip.MyKey2 = "retiree2"
				//}
				e := fmt.Errorf("NewInputParams: retiree '%s' age, retirement age and plan through age must all be specified", rip.MyKey2)
				return nil, e
			}
			needRetiree2 = true
		}
		rip.Age2 = getIPIntValue(ip["eT_Age2"])
		rip.RetireAge2 = getIPIntValue(ip["eT_RetireAge2"])
		if rip.RetireAge2 < rip.Age2 {
			rip.RetireAge2 = rip.Age2
		}
		rip.PlanThroughAge2 = getIPIntValue(ip["eT_PlanThroughAge2"])
		yearsToRetire2 = rip.RetireAge2 - rip.Age2
		through2 = rip.PlanThroughAge2 - rip.Age2
	}
	// the following must be after "joint" section
	rip.PrePlanYears = yearsToRetire1
	rip.StartPlan = rip.PrePlanYears + rip.Age1
	rip.EndPlan = through1 + 1 + rip.Age1
	rip.AgeDelta = 0
	rip.Numyr = rip.EndPlan - rip.StartPlan
	//fmt.Printf("NEED RETIREE2: %#v\n", needRetiree2)
	if needRetiree2 {
		rip.PrePlanYears = intMin(yearsToRetire1, yearsToRetire2)
		rip.StartPlan = rip.PrePlanYears + rip.Age1
		rip.EndPlan = intMax(through1, through2) + 1 + rip.Age1
		rip.AgeDelta = rip.Age1 - rip.Age2
		rip.Numyr = rip.EndPlan - rip.StartPlan
	}

	if (ip["eT_Aftatax_Contrib"] != "" && getIPIntValue(ip["eT_Aftatax_Contrib"]) != 0) ||
		ip["eT_Aftatax_ContribStartAge"] != "" ||
		ip["eT_Aftatax_ContribEndAge"] != "" {
		if ip["eT_Aftatax_Contrib"] == "" ||
			ip["eT_Aftatax_ContribStartAge"] == "" ||
			ip["eT_Aftatax_ContribEndAge"] == "" {
			e := fmt.Errorf("NewInputParams: retiree After tax account contribution requires contribution amount, start and end age for contributions to be specified")
			return nil, e
		}
	}
	rip.Aftatax = getIPIntValue(ip["eT_Aftatax"]) * multiplier
	rip.AftataxBasis = getIPIntValue(ip["eT_Aftatax_Basis"]) * multiplier
	rip.AftataxRate = getIPFloatValue(ip["eT_Aftatax_Rate"], -1.0)
	rip.AftataxContrib = getIPIntValue(ip["eT_Aftatax_Contrib"]) * multiplier
	rip.AftataxContribStart = getIPIntValue(ip["eT_Aftatax_ContribStartAge"])
	rip.AftataxContribEnd = getIPIntValue(ip["eT_Aftatax_ContribEndAge"])
	rip.AftataxContribInflate = getIPBoolValue(ip["eT_Aftatax_ContribInflate"])
	if rip.Aftatax > 0 || rip.AftataxContrib > 0 {
		rip.Accmap[Aftertax]++
	}
	if rip.AftataxRate < 0 {
		rip.AftataxRate = rip.RRate
	} else {
		rip.AftataxRate = 1 + (rip.AftataxRate / 100.0)
	}

	rip.Numacc = 0
	for _, v := range rip.Accmap {
		rip.Numacc += v
	}

	rip.Min = getIPIntValue(ip["eT_DesiredIncome"]) * multiplier
	rip.Max = getIPIntValue(ip["eT_MaxIncome"]) * multiplier
	if rip.Min > 0 && rip.Maximize != PlusEstate {
		e := fmt.Errorf("Error - [min.income] ($%d) is only valid with 'maximize=\"PlusEstate\"' however maximize currently set to '%s'",
			rip.Min, rip.Maximize.String())
		return nil, e
	}
	if rip.Max > 0 && rip.Maximize != Spending {
		e := fmt.Errorf("Error - [max.income] ($%d) is only valid with 'maximize=\"Spinding\"' however maximize currently set to '%s'",
			rip.Max, rip.Maximize.String())
		return nil, e
	}
	rip.Income = make([]stream, 0)
	for i := 1; i < MaxStreams+1; i++ {
		if ip[fmt.Sprintf("eT_Income%d", i)] != "" ||
			ip[fmt.Sprintf("eT_IncomeAmount%d", i)] != "" ||
			ip[fmt.Sprintf("eT_IncomeStartAge%d", i)] != "" ||
			ip[fmt.Sprintf("eT_IncomeEndAge%d", i)] != "" {
			if ip[fmt.Sprintf("eT_Income%d", i)] == "" ||
				ip[fmt.Sprintf("eT_IncomeAmount%d", i)] == "" ||
				ip[fmt.Sprintf("eT_IncomeStartAge%d", i)] == "" ||
				ip[fmt.Sprintf("eT_IncomeEndAge%d", i)] == "" {
				e := fmt.Errorf("NewInputParams: retiree income stream '%s' requires name/tag, amount, start and end age all to be specified", ip[fmt.Sprintf("eT_Income%d", i)])
				return nil, e
			}
			sp := stream{
				Tag:      ip[fmt.Sprintf("eT_Income%d", i)],
				Amount:   getIPIntValue(ip[fmt.Sprintf("eT_IncomeAmount%d", i)]) * multiplier,
				StartAge: getIPIntValue(ip[fmt.Sprintf("eT_IncomeStartAge%d", i)]),
				EndAge:   getIPIntValue(ip[fmt.Sprintf("eT_IncomeEndAge%d", i)]),
				Inflate:  getIPBoolValue(ip[fmt.Sprintf("eT_IncomeInflate%d", i)]),
				Tax:      getIPBoolValue(ip[fmt.Sprintf("eT_IncomeTax%d", i)]),
			}
			rip.Income = append(rip.Income, sp)
		}
	}
	rip.Expense = make([]stream, 0)
	for i := 1; i < MaxStreams+1; i++ {
		if ip[fmt.Sprintf("eT_Expense%d", i)] != "" ||
			ip[fmt.Sprintf("eT_ExpenseAmount%d", i)] != "" ||
			ip[fmt.Sprintf("eT_ExpenseStartAge%d", i)] != "" ||
			ip[fmt.Sprintf("eT_ExpenseEndAge%d", i)] != "" {
			if ip[fmt.Sprintf("eT_Expense%d", i)] == "" ||
				ip[fmt.Sprintf("eT_ExpenseAmount%d", i)] == "" ||
				ip[fmt.Sprintf("eT_ExpenseStartAge%d", i)] == "" ||
				ip[fmt.Sprintf("eT_ExpenseEndAge%d", i)] == "" {
				e := fmt.Errorf("NewInputParams: retiree expense stream '%s' requires name/tag, amount, start and end age all to be specified", ip[fmt.Sprintf("eT_Expense%d", i)])
				return nil, e
			}
			sp := stream{
				Tag:      ip[fmt.Sprintf("eT_Expense%d", i)],
				Amount:   getIPIntValue(ip[fmt.Sprintf("eT_ExpenseAmount%d", i)]) * multiplier,
				StartAge: getIPIntValue(ip[fmt.Sprintf("eT_ExpenseStartAge%d", i)]),
				EndAge:   getIPIntValue(ip[fmt.Sprintf("eT_ExpenseEndAge%d", i)]),
				Inflate:  getIPBoolValue(ip[fmt.Sprintf("eT_ExpenseInflate%d", i)]),
				Tax:      false,
			}
			rip.Expense = append(rip.Expense, sp)
		}
	}
	rip.Assets = make([]asset, 0)
	for i := 1; i < MaxStreams; i++ {
		if ip[fmt.Sprintf("eT_Asset%d", i)] != "" ||
			ip[fmt.Sprintf("eT_AssetValue%d", i)] != "" ||
			ip[fmt.Sprintf("eT_AssetAgeToSell%d", i)] != "" ||
			ip[fmt.Sprintf("eT_AssetOwedAtAgeToSell%d", i)] != "" ||
			ip[fmt.Sprintf("eT_AssetPrimaryResidence%d", i)] != "" ||
			ip[fmt.Sprintf("eT_AssetCostAndImprovements%d", i)] != "" {
			if ip[fmt.Sprintf("eT_Asset%d", i)] == "" ||
				ip[fmt.Sprintf("eT_AssetValue%d", i)] == "" ||
				ip[fmt.Sprintf("eT_AssetAgeToSell%d", i)] == "" ||
				ip[fmt.Sprintf("eT_AssetOwedAtAgeToSell%d", i)] == "" ||
				ip[fmt.Sprintf("eT_AssetPrimaryResidence%d", i)] == "" ||
				ip[fmt.Sprintf("eT_AssetCostAndImprovements%d", i)] == "" {
				e := fmt.Errorf("NewInputParams: retiree assets '%s' requires name/tag, value, age to sell, amount owed at age to sell, cost plus improvements and whether the asset is the primary residence, all to be specified", ip[fmt.Sprintf("eT_Asset%d", i)])
				return nil, e
			}
			ap := asset{
				Tag:                 ip[fmt.Sprintf("eT_Asset%d", i)],
				Value:               getIPIntValue(ip[fmt.Sprintf("eT_AssetValue%d", i)]) * multiplier,
				AgeToSell:           getIPIntValue(ip[fmt.Sprintf("eT_AssetAgeToSell%d", i)]),
				CostAndImprovements: getIPIntValue(ip[fmt.Sprintf("eT_AssetCostAndImprovements%d", i)]) * multiplier,
				OwedAtAgeToSell:     getIPIntValue(ip[fmt.Sprintf("eT_AssetOwedAtAgeToSell%d", i)]) * multiplier,
				PrimaryResidence:    getIPBoolValue(ip[fmt.Sprintf("eT_AssetPrimaryResidence%d", i)]),
				AssetRRatePercent:   getIPFloatValue(ip[fmt.Sprintf("eT_AssetRRatePercent%d", i)], -1.0),
				BrokeragePercent:    getIPFloatValue(ip[fmt.Sprintf("eT_AssetBrokeragePercent%d", i)], -1.0),
			}
			if ap.AssetRRatePercent < 0 {
				ap.AssetRRatePercent = rip.RRatePercent
			}
			ap.AssetRRate = 1 + ap.AssetRRatePercent/100.0
			if ap.BrokeragePercent < 0 {
				ap.BrokeragePercent = BrokeragePercentDefault
			}
			rip.Assets = append(rip.Assets, ap)
		}
	}
	//
	// some more tests for consistancy
	//
	if rip.DefinedContributionPlanEnd2 >= rip.RetireAge2 {
		str := fmt.Sprintf("Warning - Normally a define contribution plan ends with employment prior to retirement begining; %s's defined contribution plan ends at age %d while retirement begins at age %d",
			rip.MyKey2, rip.DefinedContributionPlanEnd2, rip.RetireAge2)
		warnList.AppendWarning(str)
	}
	if rip.DefinedContributionPlanEnd1 >= rip.RetireAge1 {
		str := fmt.Sprintf("Warning - Normally a define contribution plan ends with employment prior to retirement begining; %s's defined contribution plan ends at age %d while retirement begins at age %d",
			rip.MyKey1, rip.DefinedContributionPlanEnd1, rip.RetireAge1)
		warnList.AppendWarning(str)
	}

	//fmt.Printf("\n&&&&\n%v\n&&&&\n", rip)

	return &rip, nil
}

var InputStrDefs = []string{
	"setName",
	"filingStatus",
	"key1",
	"key2",
	"eT_Age1",
	"eT_Age2",
	"eT_RetireAge1",
	"eT_RetireAge2",
	"eT_PlanThroughAge1",
	"eT_PlanThroughAge2",
	"eT_DefinedContributionPlanStart1",
	"eT_DefinedContributionPlanStart2",
	"eT_DefinedContributionPlanEnd1",
	"eT_DefinedContributionPlanEnd2",
	"eT_PIA1",
	"eT_PIA2",
	"eT_SS_Start1",
	"eT_SS_Start2",
	"eT_TDRA1",
	"eT_TDRA2",
	"eT_TDRA_Rate1",
	"eT_TDRA_Rate2",
	"eT_TDRA_Contrib1",
	"eT_TDRA_Contrib2",
	"eT_TDRA_ContribStartAge1",
	"eT_TDRA_ContribStartAge2",
	"eT_TDRA_ContribEndAge1",
	"eT_TDRA_ContribEndAge2",
	"eT_TDRA_ContribInflate1",
	"eT_TDRA_ContribInflate2",
	"eT_Roth1",
	"eT_Roth2",
	"eT_Roth_Rate1",
	"eT_Roth_Rate2",
	"eT_Roth_Contrib1",
	"eT_Roth_Contrib2",
	"eT_Roth_ContribStartAge1",
	"eT_Roth_ContribStartAge2",
	"eT_Roth_ContribEndAge1",
	"eT_Roth_ContribEndAge2",
	"eT_Roth_ContribInflate1",
	"eT_Roth_ContribInflate2",
	"eT_Aftatax",
	"eT_Aftatax_Basis",
	"eT_Aftatax_Rate",
	"eT_Aftatax_Contrib",
	"eT_Aftatax_ContribStartAge",
	"eT_Aftatax_ContribEndAge",
	"eT_Aftatax_ContribInflate",

	"eT_DesiredIncome",
	"eT_MaxIncome",

	"eT_iRatePercent",
	"eT_rRatePercent",
	"eT_maximize",
	"dollarsInThousands",
}
var InputStreamStrDefs = []string{
	"eT_Income",
	"eT_IncomeAmount",
	"eT_IncomeStartAge",
	"eT_IncomeEndAge",
	"eT_IncomeInflate",
	"eT_IncomeTax",
	"eT_Expense",
	"eT_ExpenseAmount",
	"eT_ExpenseStartAge",
	"eT_ExpenseEndAge",
	"eT_ExpenseInflate",
	"eT_ExpenseTax",
	"eT_Asset",
	"eT_AssetValue",
	"eT_AssetAgeToSell",
	"eT_AssetCostAndImprovements",
	"eT_AssetOwedAtAgeToSell",
	"eT_AssetPrimaryResidence",
	"eT_AssetRRatePercent",
	"eT_AssetBrokeragePercent",
}

func UpdateInputStringsMap(ipsm *map[string]string, key, value string) error {
	if *ipsm == nil {
		//fmt.Printf("ipsm is nil, Calling NewInputStringsMap()\n")
		*ipsm = NewInputStringsMap()
		//fmt.Printf("ipsm new size is: %d\n", len(*ipsm))
	}
	//fmt.Printf("ipsm type: %T, %#v\n", ipsm, ipsm)
	_, ok := (*ipsm)[key]
	if !ok {
		e := fmt.Errorf("UpdateInputStringsMap: Attempting to update non-present field: %q", key)
		return e
	}
	(*ipsm)[key] = value
	return nil
}

// NewInputStringsMap returns a map with all available settings set to the empty string
func NewInputStringsMap() map[string]string {
	//This functions is the One Source of Truth for the avalable settings
	ipsm := map[string]string{}

	for _, v := range InputStrDefs {
		ipsm[v] = ""
	}

	for i := 1; i < MaxStreams+1; i++ {
		for _, v := range InputStreamStrDefs {
			ipsm[fmt.Sprintf("%s%d", v, i)] = ""
		}
	}
	return ipsm
}
