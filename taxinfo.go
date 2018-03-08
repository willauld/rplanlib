package rplanlib

import (
	"math"
)

type tableref2d *[][]float64
type tableref1d *[]float64

// 2017 table (predict it moves with inflation?)
// married joint, married separate, single
// Table Columns:
// [braket $ start,
//  bracket size,
//  marginal rate,
//  total tax from all lower brackets ] ### TODO check if this field is used delete if not!
var marriedjointtax2017 = &[][]float64{
	{0, 18650, 0.10, 0},
	{18650, 57250, 0.15, 1865},
	{75900, 77200, 0.25, 10452.5},
	{153100, 80250, 0.28, 29752.5},
	{233350, 183350, 0.33, 52222.5},
	{416700, 54000, 0.35, 112728},
	{470700, -2, 0.396, 131628},
}

var marriedseparatetax2017 = &[][]float64{
	{0, 9325, 0.10, 0},
	{9325, 28625, 0.15, 932.5},
	{37950, 38600, 0.25, 5226.25},
	{76550, 40125, 0.28, 14876.25},
	{116675, 91675, 0.33, 26111.25},
	{208350, 27000, 0.35, 56364.00},
	{235350, -2, 0.396, 65814.00},
}

var singletax2017 = &[][]float64{
	{0, 9325, 0.10, 0},
	{9325, 28625, 0.15, 932.5},
	{37950, 53950, 0.25, 5226.25},
	{91900, 99750, 0.28, 18713.75},
	{191650, 225050, 0.33, 46643.75},
	{416700, 1700, 0.35, 120910.25},
	{418400, -2, 0.396, 121505.25},
}

// Table Columns:
// [braket $ start,
//  bracket size,
//  marginal rate ]
var marriedjointcapitalgains2017 = &[][]float64{
	{0, 75900, 0.0},
	{75900, 394800, 0.15},
	{470700, -3, 0.20},
}

var marriedseparatecapitalgains2017 = &[][]float64{
	{0, 76550, 0.0},
	{76550, 158800, 0.15},
	{235350, -3, 0.20},
}

var singlecapitalgains2017 = &[][]float64{
	{0, 37950, 0.0},
	{37950, 380450, 0.15},
	{418400, -3, 0.20},
}

// Required Minimal Distributions from IRA starting with age 70
// https://www.irs.gov/publications/p590b#en_US_2016_publink1000231258
// Using appendix B table III in all cases.
var marriedjointRMD = &[]float64{
	27.4, 26.5, 25.6, 24.7, 23.8, 22.9, 22.0, 21.2, 20.3, 19.5, //age 70-79
	18.7, 17.9, 17.1, 16.3, 15.5, 14.8, 14.1, 13.4, 12.7, 12.0, //age 80-89
	11.4, 10.8, 10.2, 9.6, 9.1, 8.6, 8.1, 7.6, 7.1, 6.7, //age 90-99
	6.3, 5.9, 5.5, 5.2, 4.9, 4.5, 4.2, 3.9, 3.7, 3.4, //age 100+
	3.1, 2.9, 2.6, 2.4, 2.1, 1.9, 1.9, 1.9, 1.9, 1.9,
}

var marriedseparateRMD = marriedjointRMD
var singleRMD = marriedjointRMD

var marriedjointstded2017 = 12700 + 2*4050 //std dedction + 2 prsonal exemptions
var marriedseparatestded2017 = 9350 + 4050 //std dedction + 1 prsonal exemptions
var singlestded2017 = 6350 + 4050          //std dedction + 1 personal exmptions

var jointprimeresidence2017 = 500000
var singleprimresidence2017 = 250000

// Taxinfo contains the centeral tax information
type Taxinfo struct {
	Taxtable       tableref2d // income tax table
	Capgainstable  tableref2d // capital gains tax table
	RMD            tableref1d // Required Minimum Distribution table
	Stded          float64    // standard deduction
	Primeresidence float64    // exclusion for prime residence

	AccountEstateTax map[string]float64
	Contribspecs     map[string]float64

	Penalty      float64 // for early withdrawal
	SStaxable    float64 // taxable portion of SS
	SSnotTaxable float64 // non-taxable portion of SS
}

// TODO: Should I merge NewTaxInfo() and set_retirement_staus() I'm thinking I should!!!

//NewTaxInfo creates the applicable tax structure
func NewTaxInfo(status string) Taxinfo {
	sstaxable := 0.85
	ssnontaxable := 1 - sstaxable
	ti := Taxinfo{
		// Account specs contains some initial information # TODO if maxcontrib not used delete
		AccountEstateTax: map[string]float64{
			"IRA":      0.85,
			"roth":     1.0,
			"aftertax": 0.9,
		},

		// 401(k), 403(b) and TSP currently have the same limits
		Contribspecs: map[string]float64{"401k": 18000, "401kCatchup": 6000,
			"TDRA": 5500, "TDRACatchup": 1000, "CatchupAge": 50},

		Penalty:      0.1,       // 10% early withdrawal penalty
		SStaxable:    sstaxable, // maximum portion of SS that is taxable
		SSnotTaxable: ssnontaxable,
	}
	if status == "single" {
		ti.Taxtable = singletax2017
		ti.Capgainstable = singlecapitalgains2017
		ti.Stded = float64(singlestded2017)
		ti.RMD = singleRMD
		ti.Primeresidence = float64(singleprimresidence2017)
	} else if status == "mseparate" {
		ti.Taxtable = marriedseparatetax2017
		ti.Capgainstable = marriedseparatecapitalgains2017
		ti.Stded = float64(marriedseparatestded2017)
		ti.RMD = marriedseparateRMD
		ti.Primeresidence = float64(singleprimresidence2017)
	} else { // status == 'joint':
		ti.Taxtable = marriedjointtax2017
		ti.Capgainstable = marriedjointcapitalgains2017
		ti.Stded = float64(marriedjointstded2017)
		ti.RMD = marriedjointRMD
		ti.Primeresidence = float64(jointprimeresidence2017)
	}
	//print('taxtable:\n', self.taxtable, '\n')
	//print('capgainstable:\n', self.capgainstable, '\n')
	//print('stded:\n', self.stded, '\n')
	//print('RMD:\n', self.RMD, '\n')
	return ti
}

/* Keep???? TODO FIXME
func expandYears(numyr, ageAtStart, agestr) ([]float64) {
	bucket := make([]float64, numyr)
        for age := agelist(agestr) {
            year := age - ageAtStart
            if year < 0 {
                continue
			} else if year >= numyr {
                break
			} else {
                bucket[year] = 1
			}
		}
		return bucket
}
*/

// TODO: FIXME NEED UNIT TEST FOR THIS FUNCTION
// maxContributions returns the max allowable contributions for one or all retirees
func (ti Taxinfo) maxContribution(year int, yearsToInflateBy int, retirees []retiree, retireekey string, iRate float64) float64 {
	//FIXME: not currently handling 401K max contributions TODO
	max := 0.0
	for _, v := range retirees {
		if retireekey == "" || v.mykey == retireekey { // if "", Sum all retiree
			max += ti.Contribspecs["TDRA"]
			//fmt.Printf("max += tiContribspecs[TDRA]: %f\n", ti.Contribspecs["TDRA"])
			age := v.ageAtStart + year
			if age >= int(ti.Contribspecs["CatchupAge"]) {
				max += ti.Contribspecs["TDRACatchup"]
				//fmt.Printf("max += tiContribspecs[TDRACatchup]: %f\n", ti.Contribspecs["TDRACatchup"])
			}
			if v.definedContributionPlanStartAge > 0 &&
				v.definedContributionPlanEndAge >=
					v.definedContributionPlanStartAge {
				lower := v.definedContributionPlanStartAge - v.ageAtStart
				upper := v.definedContributionPlanEndAge - v.ageAtStart
				/* no lazy expantion in golang implementation, created in NewModelSpecs()
				a := v.dcpBuckets
				                    if a == nil {
				                        // lazy expantion of the definedContributionPlan info
				                        v.dcpBuckets = expandYears(startage, have_plan)
				                        a = v.dcpBuckets
									}
				*/
				if year >= lower && year <= upper {
					max += ti.Contribspecs["401k"]
					if age >= int(ti.Contribspecs["CatchupAge"]) {
						max += ti.Contribspecs["401kCatchup"]
					}
				}
			}
		}
	}
	max *= math.Pow(iRate, float64(yearsToInflateBy))
	//fmt.Printf("maxContribution: %6.0f, key: %s\n", max, retireekey)
	return max
}

// TODO: FIXME NEED UNIT TEST FOR THIS FUNCTION
// applyEarlyPenalty returns a bool indicating if an early penalty needs to be applied
func (ti Taxinfo) applyEarlyPenalty(year int, r *retiree) bool {
	response := false
	//v := r.match_retiree(retireekey)
	if r == nil {
		return response
	}
	age := r.ageAtStart + year
	if age < 60 { // IRA retirement account require penalty if withdrawn before age 59.5
		response = true
	}
	return response
}

// TODO: FIXME NEED UNIT TEST FOR THIS FUNCTION
// rmdNeeded returns the life expectance to use if needed or zero otherwise
func (ti Taxinfo) rmdNeeded(year int, r *retiree) float64 {
	rmd := 0.0
	//v = self.match_retiree(retireekey)
	if r == nil {
		//print("RMD_NEEDED() year: %d, rmd: %6.3f, Not Valid Retiree, retiree: %s" % (year, rmd, retireekey))
		return rmd
	}
	age := r.ageAtStart + year
	if age >= 70 { // IRA retirement: minimum distribution starting age 70.5
		rmd = (*ti.RMD)[age-70]
	}
	//print("RMD_NEEDED() year: %d, rmd: %6.3f, age: %d, retiree: %s" % (year, rmd, age, retireekey))
	return rmd
}
