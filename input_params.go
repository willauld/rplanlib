package rplanlib

import (
	"fmt"
	"strconv"
	"strings"
)

// MaxStreams is the maximam number of streams for each: income, expense, asset
const MaxStreams = 10

// InputParams are the model params constructed from driver program string input
type InputParams struct {
	FilingStatus        string
	MyKey1              string
	MyKey2              string
	Age1                int
	Age2                int
	RetireAge1          int
	RetireAge2          int
	PlanThroughAge1     int
	PlanThroughAge2     int
	PIA1                int
	PIA2                int
	SSStart1            int
	SSStart2            int
	TDRA1               int
	TDRA2               int
	TDRARate1           float64
	TDRARate2           float64
	TDRAContrib1        int
	TDRAContrib2        int
	TDRAContribStart1   int
	TDRAContribStart2   int
	TDRAContribEnd1     int
	TDRAContribEnd2     int
	Roth1               int
	Roth2               int
	RothRate1           float64
	RothRate2           float64
	RothContrib1        int
	RothContrib2        int
	RothContribStart1   int
	RothContribStart2   int
	RothContribEnd1     int
	RothContribEnd2     int
	Aftatax             int
	AftataxBasis        int // TODO Add to Mobile
	AftataxRate         float64
	AftataxContrib      int
	AftataxContribStart int
	AftataxContribEnd   int
	IRatePercent        float64 // TODO add to Mobile
	IRate               float64 // local, not mobile
	RRatePercent        float64 // TODO add to Mobile
	RRate               float64 // local, not mobile
	Maximize            string  // "Spending" or "PlusEstate" // TODO add to Mobile
	Min                 int     // TODO add to Mobile
	Max                 int     // TODO add to Mobile

	PrePlanYears int
	StartPlan    int
	EndPlan      int
	AgeDelta     int
	Numyr        int
	Accmap       map[string]int
	Numacc       int

	// PROTOTYPE
	Income  []stream
	Expense []stream
	Assets  []asset
}

type stream struct {
	// PROTOTYPE
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

//TODO: TESTME
func kgetIPIntValue(str string) int {
	return 1000 * getIPIntValue(str)
}

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

func getIPFloatValue(str string) float64 {
	if str == "" {
		return 0
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
	b, e := strconv.ParseBool(strings.ToLower(str))
	if e != nil {
		//fmt.Printf("GetIPFloatValue(): %s\n", e)
		panic(e)
	}
	return b
}
func verifyMaximize(s string) error {
	e := error(nil)
	if s != "Spending" && s != "PlusEstate" {
		e = fmt.Errorf("verifyMaximize: %s is not a valid option", s)
	}
	return e
}
func verifyFilingStatus(s string) error {
	e := error(nil)
	if s != "joint" && s != "mseparate" && s != "single" {
		e = fmt.Errorf("verifyFilingStatus: %s is not a valid option", s)
	}
	return e
}

// NewInputParams takes string inputs and converts them to model inputs
func NewInputParams(ip map[string]string) (*InputParams, error) {

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

	rip.RRatePercent = getIPFloatValue(ip["eT_rRatePercent"]) // TODO add to mobile
	rip.IRatePercent = getIPFloatValue(ip["eT_iRatePercent"]) // TODO add to mobile
	rip.RRate = 1 + rip.RRatePercent/100.0
	rip.IRate = 1 + rip.IRatePercent/100.0
	rip.Min = 0 // = getIPIntValue(ip["eT_min"]) // TODO add to mobile
	rip.Max = 0 // = getIPIntValue(ip["eT_max"]) // TODO add to mobile
	//	maximize:                "Spending", // or "PlusEstate"
	rip.Maximize = "Spending" // = ip["eT_Maximize"] // TODO add to mobile
	err = verifyMaximize(rip.Maximize)
	if err != nil {
		return nil, err
	}

	rip.Accmap = map[string]int{"IRA": 0, "roth": 0, "aftertax": 0}
	err = verifyFilingStatus(ip["filingStatus"])
	if err != nil {
		return nil, err
	}
	rip.FilingStatus = ip["filingStatus"]

	rip.MyKey1 = ip["key1"]
	if ip["eT_Age1"] == "" ||
		ip["eT_RetireAge1"] == "" ||
		ip["eT_PlanThroughAge1"] == "" {
		e := fmt.Errorf("NewInputParams: retiree age, retirement age and plan through age must all be specified")
		return nil, e
	}
	rip.Age1 = getIPIntValue(ip["eT_Age1"])
	rip.RetireAge1 = getIPIntValue(ip["eT_RetireAge1"])
	if rip.RetireAge1 < rip.Age1 {
		rip.RetireAge1 = rip.Age1
	}
	rip.PlanThroughAge1 = getIPIntValue(ip["eT_PlanThroughAge1"])
	yearsToRetire1 := rip.RetireAge1 - rip.Age1
	rip.PrePlanYears = yearsToRetire1
	through1 := rip.PlanThroughAge1 - rip.Age1

	if ip["eT_PIA1"] != "" || ip["eT_SS_Start1"] != "" {
		if ip["eT_PIA1"] == "" || ip["eT_SS_Start1"] == "" {
			e := fmt.Errorf("NewInputParams: retiree social security PIA and start age both must be specified")
			return nil, e
		}
	}
	rip.PIA1 = getIPIntValue(ip["eT_PIA1"]) * multiplier
	rip.SSStart1 = getIPIntValue(ip["eT_SS_Start1"])

	if ip["eT_TDRA_Contrib1"] != "" ||
		ip["eT_TDRA_ContribStartAge1"] != "" ||
		ip["eT_TDRA_ContribEndAge1"] != "" {
		if ip["eT_TDRA_Contrib1"] == "" ||
			ip["eT_TDRA_ContribStartAge1"] == "" ||
			ip["eT_TDRA_ContribEndAge1"] == "" {
			e := fmt.Errorf("NewInputParams: retiree TDRA contribution requires contribution amount, start and end age for contributions be specified")
			return nil, e
		}
	}
	rip.TDRA1 = getIPIntValue(ip["eT_TDRA1"]) * multiplier
	rip.TDRARate1 = getIPFloatValue(ip["eT_TDRA_Rate1"])
	rip.TDRAContrib1 = getIPIntValue(ip["eT_TDRA_Contrib1"]) * multiplier
	rip.TDRAContribStart1 = getIPIntValue(ip["eT_TDRA_ContribStartAge1"])
	rip.TDRAContribEnd1 = getIPIntValue(ip["eT_TDRA_ContribEndAge1"])
	if rip.TDRA1 > 0 || rip.TDRAContrib1 > 0 {
		rip.Accmap["IRA"]++
	}

	if ip["eT_Roth_Contrib1"] != "" ||
		ip["eT_Roth_ContribStartAge1"] != "" ||
		ip["eT_Roth_ContribEndAge1"] != "" {
		if ip["eT_Roth_Contrib1"] == "" ||
			ip["eT_Roth_ContribStartAge1"] == "" ||
			ip["eT_Roth_ContribEndAge1"] == "" {
			e := fmt.Errorf("NewInputParams: retiree Roth contribution requires contribution amount, start and end age for contributions be specified")
			return nil, e
		}
	}
	rip.Roth1 = getIPIntValue(ip["eT_Roth1"]) * multiplier
	rip.RothRate1 = getIPFloatValue(ip["eT_Roth_Rate1"])
	rip.RothContrib1 = getIPIntValue(ip["eT_Roth_Contrib1"]) * multiplier
	rip.RothContribStart1 = getIPIntValue(ip["eT_Roth_ContribStartAge1"])
	rip.RothContribEnd1 = getIPIntValue(ip["eT_Roth_ContribEndAge1"])
	if rip.Roth1 > 0 || rip.RothContrib1 > 0 {
		rip.Accmap["roth"]++
	}

	var through2 int
	if rip.FilingStatus == "joint" {
		rip.MyKey2 = ip["key2"]
		if ip["eT_Age2"] == "" ||
			ip["eT_RetireAge2"] == "" ||
			ip["eT_PlanThroughAge2"] == "" {
			e := fmt.Errorf("NewInputParams: retiree age, retirement age and plan through age must all be specified")
			return nil, e
		}
		rip.Age2 = getIPIntValue(ip["eT_Age2"])
		rip.RetireAge2 = getIPIntValue(ip["eT_RetireAge2"])
		if rip.RetireAge2 < rip.Age2 {
			rip.RetireAge2 = rip.Age2
		}
		rip.PlanThroughAge2 = getIPIntValue(ip["eT_PlanThroughAge2"])
		yearsToRetire2 := rip.RetireAge2 - rip.Age2
		rip.PrePlanYears = intMin(yearsToRetire1, yearsToRetire2)
		through2 = rip.PlanThroughAge2 - rip.Age2

		if ip["eT_PIA2"] != "" || ip["eT_SS_Start2"] != "" ||
			ip["eT_PIA1"] != "" {
			if ip["eT_PIA1"] == "" {
				// if any SS set both must be specified
				e := fmt.Errorf("NewInputParams: both retiree social security PIA and start age must be specified if either retiree is")
				return nil, e
			}
			if ip["eT_PIA2"] == "" || ip["eT_SS_Start2"] == "" {
				e := fmt.Errorf("NewInputParams: retiree social security PIA and start age both must be specified")
				return nil, e
			}
		}
		rip.PIA2 = getIPIntValue(ip["eT_PIA2"]) * multiplier
		rip.SSStart2 = getIPIntValue(ip["eT_SS_Start2"])

		if ip["eT_TDRA_Contrib2"] != "" ||
			ip["eT_TDRA_ContribStartAge2"] != "" ||
			ip["eT_TDRA_ContribEndAge2"] != "" {
			if ip["eT_TDRA_Contrib2"] == "" ||
				ip["eT_TDRA_ContribStartAge2"] == "" ||
				ip["eT_TDRA_ContribEndAge2"] == "" {
				e := fmt.Errorf("NewInputParams: retiree TDRA contribution requires contribution amount, start and end age for contributions be specified")
				return nil, e
			}
		}
		rip.TDRA2 = getIPIntValue(ip["eT_TDRA2"]) * multiplier
		rip.TDRARate2 = getIPFloatValue(ip["eT_TDRA_Rate2"])
		rip.TDRAContrib2 = getIPIntValue(ip["eT_TDRA_Contrib2"]) * multiplier
		rip.TDRAContribStart2 = getIPIntValue(ip["eT_TDRA_ContribStartAge2"])
		rip.TDRAContribEnd2 = getIPIntValue(ip["eT_TDRA_ContribEndAge2"])
		if rip.TDRA2 > 0 || rip.TDRAContrib2 > 0 {
			rip.Accmap["IRA"]++
		}

		if ip["eT_Roth_Contrib2"] != "" ||
			ip["eT_Roth_ContribStartAge2"] != "" ||
			ip["eT_Roth_ContribEndAge2"] != "" {
			if ip["eT_Roth_Contrib2"] == "" ||
				ip["eT_Roth_ContribStartAge2"] == "" ||
				ip["eT_Roth_ContribEndAge2"] == "" {
				e := fmt.Errorf("NewInputParams: retiree Roth contribution requires contribution amount, start and end age for contributions be specified")
				return nil, e
			}
		}
		rip.Roth2 = getIPIntValue(ip["eT_Roth2"]) * multiplier
		rip.RothRate2 = getIPFloatValue(ip["eT_Roth_Rate2"])
		rip.RothContrib2 = getIPIntValue(ip["eT_Roth_Contrib2"]) * multiplier
		rip.RothContribStart2 = getIPIntValue(ip["eT_Roth_ContribStartAge2"])
		rip.RothContribEnd2 = getIPIntValue(ip["eT_Roth_ContribEndAge2"])
		if rip.Roth2 > 0 || rip.RothContrib2 > 0 {
			rip.Accmap["roth"]++
		}
	}
	// the following must be after "joint" section
	rip.StartPlan = rip.PrePlanYears + rip.Age1
	rip.EndPlan = intMax(through1, through2) + 1 + rip.Age1
	rip.AgeDelta = rip.Age1 - rip.Age2
	rip.Numyr = rip.EndPlan - rip.StartPlan

	if ip["eT_Aftatax_Contrib"] != "" ||
		ip["eT_Aftatax_ContribStartAge"] != "" ||
		ip["eT_Aftatax_ContribEndAge"] != "" {
		if ip["eT_Aftatax_Contrib"] == "" ||
			ip["eT_Aftatax_ContribStartAge"] == "" ||
			ip["eT_Aftatax_ContribEndAge"] == "" {
			e := fmt.Errorf("NewInputParams: retiree After tax account contribution requires contribution amount, start and end age for contributions be specified")
			return nil, e
		}
	}
	rip.Aftatax = getIPIntValue(ip["eT_Aftatax"]) * multiplier
	rip.AftataxRate = getIPFloatValue(ip["eT_Aftatax_Rate"])
	rip.AftataxContrib = getIPIntValue(ip["eT_Aftatax_Contrib"]) * multiplier
	rip.AftataxContribStart = getIPIntValue(ip["eT_Aftatax_ContribStartAge"])
	rip.AftataxContribEnd = getIPIntValue(ip["eT_Aftatax_ContribEndAge"])
	if rip.Aftatax > 0 || rip.AftataxContrib > 0 {
		rip.Accmap["aftertax"]++
	}

	rip.Numacc = 0
	for _, v := range rip.Accmap {
		rip.Numacc += v
	}

	//PROTOTYPE WORK ? FLAT ? LINK LIST ?
	rip.Income = make([]stream, 0)
	for i := 1; i < MaxStreams; i++ {
		if ip[fmt.Sprintf("eT_Income%d", i)] != "" ||
			ip[fmt.Sprintf("eT_IncomeAmount%d", i)] != "" ||
			ip[fmt.Sprintf("eT_IncomeStartAge%d", i)] != "" ||
			ip[fmt.Sprintf("eT_IncomeEndAge%d", i)] != "" {
			if ip[fmt.Sprintf("eT_Income%d", i)] == "" ||
				ip[fmt.Sprintf("eT_IncomeAmount%d", i)] == "" ||
				ip[fmt.Sprintf("eT_IncomeStartAge%d", i)] == "" ||
				ip[fmt.Sprintf("eT_IncomeEndAge%d", i)] == "" {
				e := fmt.Errorf("NewInputParams: retiree income stream requires name/tag, amount, start and end age all to be specified")
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
	for i := 1; i < MaxStreams; i++ {
		if ip[fmt.Sprintf("eT_Expense%d", i)] != "" ||
			ip[fmt.Sprintf("eT_ExpenseAmount%d", i)] != "" ||
			ip[fmt.Sprintf("eT_ExpenseStartAge%d", i)] != "" ||
			ip[fmt.Sprintf("eT_ExpenseEndAge%d", i)] != "" {
			if ip[fmt.Sprintf("eT_Expense%d", i)] == "" ||
				ip[fmt.Sprintf("eT_ExpenseAmount%d", i)] == "" ||
				ip[fmt.Sprintf("eT_ExpenseStartAge%d", i)] == "" ||
				ip[fmt.Sprintf("eT_ExpenseEndAge%d", i)] == "" {
				e := fmt.Errorf("NewInputParams: retiree expense stream requires name/tag, amount, start and end age all to be specified")
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
				e := fmt.Errorf("NewInputParams: retiree assets requires name/tag, value, age to sell, amount owed at age to sell, cost plus improvements and whether the asset is the primary residence, all to be specified")
				return nil, e
			}
			ap := asset{
				Tag:                 ip[fmt.Sprintf("eT_Asset%d", i)],
				Value:               getIPIntValue(ip[fmt.Sprintf("eT_AssetValue%d", i)]) * multiplier,
				AgeToSell:           getIPIntValue(ip[fmt.Sprintf("eT_AssetAgeToSell%d", i)]),
				CostAndImprovements: getIPIntValue(ip[fmt.Sprintf("eT_AssetCostAndImprovements%d", i)]) * multiplier,
				OwedAtAgeToSell:     getIPIntValue(ip[fmt.Sprintf("eT_AssetOwedAtAgeToSell%d", i)]) * multiplier,
				PrimaryResidence:    getIPBoolValue(ip[fmt.Sprintf("eT_AssetPrimaryResidence%d", i)]),
				AssetRRatePercent:   getIPFloatValue(ip[fmt.Sprintf("eT_AssetRRatePercent%d", i)]),
				BrokeragePercent:    getIPFloatValue(ip[fmt.Sprintf("eT_AssetBrokeragePercent%d", i)]),
			}
			ap.AssetRRate = 1 + ap.AssetRRatePercent/100.0
			rip.Assets = append(rip.Assets, ap)
		}
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
