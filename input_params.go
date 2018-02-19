package rplanlib

import (
	"fmt"
	"strconv"
)

// InputParams are the model params constructed from driver program string input
type InputParams struct {
	filingStatus        string
	myKey1              string
	myKey2              string
	age1                int
	age2                int
	retireAge1          int
	retireAge2          int
	planThroughAge1     int
	planThroughAge2     int
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
	iRate               float64 // TODO add to Mobile
	rRate               float64 // TODO add to Mobile
	maximize            string  // "Spending" or "PlusEstate" // TODO add to Mobile
	min                 int     // TODO add to Mobile
	max                 int     // TODO add to Mobile

	prePlanYears int
	startPlan    int
	endPlan      int
	ageDelta     int
	numyr        int
	accmap       map[string]int
	numacc       int
}

//TODO: TESTME
func kgetIPIntValue(str string) int {
	return 1000 * getIPIntValue(str)
}

//TODO: TESTME
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

//TODO: TESTME
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
	rip := InputParams{}

	rip.rRate = 1.06  // = getIPFloatValue(ip["eT_Gen_rRate"]) // TODO add to mobile
	rip.iRate = 1.025 // = getIPFloatValue(ip["eT_Gen_iRate"]) // TODO add to mobile
	rip.min = 0       // = getIPIntValue(ip["eT_min"]) // TODO add to mobile
	rip.max = 0       // = getIPIntValue(ip["eT_max"]) // TODO add to mobile
	//	maximize:                "Spending", // or "PlusEstate"
	rip.maximize = "Spending" // = ip["eT_Maximize"] // TODO add to mobile

	rip.accmap = map[string]int{"IRA": 0, "roth": 0, "aftertax": 0}
	err = verifyFilingStatus(ip["filingStatus"])
	if err != nil {
		return nil, err
	}
	rip.filingStatus = ip["filingStatus"]
	rip.myKey1 = ip["key1"] // TODO: these two should come from ip FIXME
	rip.age1 = getIPIntValue(ip["eT_Age1"])
	rip.retireAge1 = getIPIntValue(ip["eT_RetireAge1"])
	if rip.retireAge1 < rip.age1 {
		rip.retireAge1 = rip.age1
	}
	rip.planThroughAge1 = getIPIntValue(ip["eT_PlanThroughAge1"])
	yearsToRetire1 := rip.retireAge1 - rip.age1
	rip.prePlanYears = yearsToRetire1
	through1 := rip.planThroughAge1 - rip.age1

	rip.PIA1 = kgetIPIntValue(ip["eT_PIA1"])
	rip.SSStart1 = getIPIntValue(ip["eT_SS_Start1"])

	rip.TDRA1 = kgetIPIntValue(ip["eT_TDRA1"])
	rip.TDRARate1 = getIPFloatValue(ip["eT_TDRA_Rate1"])
	rip.TDRAContrib1 = kgetIPIntValue(ip["eT_TDRA_Contrib1"])
	rip.TDRAContribStart1 = getIPIntValue(ip["eT_TDRA_ContribStartAge1"])
	rip.TDRAContribEnd1 = getIPIntValue(ip["eT_TDRA_ContribEndAge1"])
	if rip.TDRA1 > 0 {
		rip.accmap["IRA"]++
	}

	rip.Roth1 = kgetIPIntValue(ip["eT_Roth1"])
	rip.RothRate1 = getIPFloatValue(ip["eT_Roth_Rate1"])
	rip.RothContrib1 = kgetIPIntValue(ip["eT_Roth_Contrib1"])
	rip.RothContribStart1 = getIPIntValue(ip["eT_Roth_ContribStartAge1"])
	rip.RothContribEnd1 = getIPIntValue(ip["eT_Roth_ContribEndAge1"])
	if rip.Roth1 > 0 {
		rip.accmap["roth"]++
	}

	var through2 int
	if rip.filingStatus == "joint" {
		rip.myKey2 = ip["key2"]
		rip.age2 = getIPIntValue(ip["eT_Age2"])
		rip.retireAge2 = getIPIntValue(ip["eT_RetireAge2"])
		if rip.retireAge2 < rip.age2 {
			rip.retireAge2 = rip.age2
		}
		rip.planThroughAge2 = getIPIntValue(ip["eT_PlanThroughAge2"])
		yearsToRetire2 := rip.retireAge2 - rip.age2
		rip.prePlanYears = intMin(yearsToRetire1, yearsToRetire2)
		through2 = rip.planThroughAge2 - rip.age2

		rip.PIA2 = kgetIPIntValue(ip["eT_PIA2"])
		rip.SSStart2 = getIPIntValue(ip["eT_SS_Start2"])

		rip.TDRA2 = kgetIPIntValue(ip["eT_TDRA2"])
		rip.TDRARate2 = getIPFloatValue(ip["eT_TDRA_Rate2"])
		rip.TDRAContrib2 = kgetIPIntValue(ip["eT_TDRA_Contrib2"])
		rip.TDRAContribStart2 = getIPIntValue(ip["eT_TDRA_ContribStartAge2"])
		rip.TDRAContribEnd2 = getIPIntValue(ip["eT_TDRA_ContribEndAge2"])
		if rip.TDRA2 > 0 {
			rip.accmap["IRA"]++
		}

		rip.Roth2 = kgetIPIntValue(ip["eT_Roth2"])
		rip.RothRate2 = getIPFloatValue(ip["eT_Roth_Rate2"])
		rip.RothContrib2 = kgetIPIntValue(ip["eT_Roth_Contrib2"])
		rip.RothContribStart2 = getIPIntValue(ip["eT_Roth_ContribStartAge2"])
		rip.RothContribEnd2 = getIPIntValue(ip["eT_Roth_ContribEndAge2"])
		if rip.Roth2 > 0 {
			rip.accmap["roth"]++
		}
	}
	// the following must be after "joint" section
	rip.startPlan = rip.prePlanYears + rip.age1
	rip.endPlan = intMax(through1, through2) + 1 + rip.age1
	rip.ageDelta = rip.age1 - rip.age2
	rip.numyr = rip.endPlan - rip.startPlan

	rip.Aftatax = kgetIPIntValue(ip["eT_Aftatax"])
	rip.AftataxRate = getIPFloatValue(ip["eT_Aftatax_Rate"])
	rip.AftataxContrib = kgetIPIntValue(ip["eT_Aftatax_Contrib"])
	rip.AftataxContribStart = getIPIntValue(ip["eT_Aftatax_ContribStartAge"])
	rip.AftataxContribEnd = getIPIntValue(ip["eT_Aftatax_ContribEndAge"])
	if rip.Aftatax > 0 {
		rip.accmap["aftertax"]++
	}

	rip.numacc = 0
	for _, v := range rip.accmap {
		rip.numacc += v
	}

	//fmt.Printf("\n&&&&\n%v\n&&&&\n", rip)

	return &rip, nil
}
