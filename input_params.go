package rplanlib

import (
	"strconv"
)

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
	AftataxRate         float64
	AftataxContrib      int
	AftataxContribStart int
	AftataxContribEnd   int

	prePlanYears int
	startPlan    int
	endPlan      int
	numyr        int
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

// NewInputParams takes string inputs and converts them to model inputs
func NewInputParams(ip map[string]string) InputParams {

	const rRate = 1.06
	const iRate = 1.025

	rip := InputParams{}

	rip.filingStatus = ip["filingStatus"]
	rip.myKey1 = "retiree1"
	rip.myKey2 = "retiree2"
	rip.age1 = getIPIntValue(ip["eT_Age1"])
	rip.age2 = getIPIntValue(ip["eT_Age2"])
	rip.retireAge1 = getIPIntValue(ip["eT_RetireAge1"])
	if rip.retireAge1 < rip.age1 {
		rip.retireAge1 = rip.age1
	}
	rip.retireAge2 = getIPIntValue(ip["eT_RetireAge2"])
	if rip.retireAge2 < rip.age2 {
		rip.retireAge2 = rip.age2
	}
	rip.planThroughAge1 = getIPIntValue(ip["eT_PlanThroughAge1"])
	rip.planThroughAge2 = getIPIntValue(ip["eT_PlanThroughAge2"])
	yearsToRetire1 := rip.retireAge1 - rip.age1
	yearsToRetire2 := rip.retireAge2 - rip.age2
	rip.prePlanYears = intMin(yearsToRetire1, yearsToRetire2)
	rip.startPlan = rip.prePlanYears + rip.age1
	through1 := rip.planThroughAge1 - rip.age1
	through2 := rip.planThroughAge2 - rip.age2
	rip.endPlan = intMax(through1, through2) + 1 + rip.age1
	//delta := age1 - age2

	rip.numyr = rip.endPlan - rip.startPlan

	//accounttable: []map[string]string
	//accmap: map[string]int

	rip.PIA1 = getIPIntValue(ip["eT_PIA1"])
	rip.PIA2 = getIPIntValue(ip["eT_PIA2"])
	rip.SSStart1 = getIPIntValue(ip["eT_SS_Start1"])
	rip.SSStart2 = getIPIntValue(ip["eT_SS_Start2"])

	rip.TDRA1 = getIPIntValue(ip["eT_TDRA1"])
	rip.TDRA2 = getIPIntValue(ip["eT_TDRA2"])
	rip.TDRARate1 = getIPFloatValue(ip["eT_TDRA_Rate1"])
	rip.TDRARate2 = getIPFloatValue(ip["eT_TDRA_Rate2"])
	rip.TDRAContrib1 = getIPIntValue(ip["eT_TDRA_Contrib1"])
	rip.TDRAContrib2 = getIPIntValue(ip["eT_TDRA_Contrib2"])
	rip.TDRAContribStart1 = getIPIntValue(ip["eT_TDRA_ContribStartAge1"])
	rip.TDRAContribStart2 = getIPIntValue(ip["eT_TDRA_ContribStartAge2"])
	rip.TDRAContribEnd1 = getIPIntValue(ip["eT_TDRA_ContribEndAge1"])
	rip.TDRAContribEnd2 = getIPIntValue(ip["eT_TDRA_ContribEndAge2"])

	rip.Roth1 = getIPIntValue(ip["eT_Roth1"])
	rip.Roth2 = getIPIntValue(ip["eT_Roth2"])
	rip.RothRate1 = getIPFloatValue(ip["eT_Roth_Rate1"])
	rip.RothRate2 = getIPFloatValue(ip["eT_Roth_Rate2"])
	rip.RothContrib1 = getIPIntValue(ip["eT_Roth_Contrib1"])
	rip.RothContrib2 = getIPIntValue(ip["eT_Roth_Contrib2"])
	rip.RothContribStart1 = getIPIntValue(ip["eT_Roth_ContribStartAge1"])
	rip.RothContribStart2 = getIPIntValue(ip["eT_Roth_ContribStartAge2"])
	rip.RothContribEnd1 = getIPIntValue(ip["eT_Roth_ContribEndAge1"])
	rip.RothContribEnd2 = getIPIntValue(ip["eT_Roth_ContribEndAge2"])

	rip.Aftatax = getIPIntValue(ip["eT_Aftatax"])
	rip.AftataxRate = getIPFloatValue(ip["eT_Aftatax_Rate"])
	rip.AftataxContrib = getIPIntValue(ip["eT_Aftatax_Contrib"])
	rip.AftataxContribStart = getIPIntValue(ip["eT_Aftatax_ContribStartAge"])
	rip.AftataxContribEnd = getIPIntValue(ip["eT_Aftatax_ContribEndAge"])

	//fmt.Printf("\n&&&&\n%v\n&&&&\n", rip)

	return rip
}
