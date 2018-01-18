package rplanlib

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

	Accountspecs map[string](map[string]float32)
	Contribspecs map[string]float32

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
		Accountspecs: map[string]map[string]float32{"IRA": {"tax": 0.85, "maxcontrib": 18000 + 5500*2},
			"roth":     {"tax": 1.0, "maxcontrib": 5500 * 2},
			"aftertax": {"tax": 0.9, "basis": 0}},

		// 401(k), 403(b) and TSP currently have the same limits
		Contribspecs: map[string]float32{"401k": 18000, "401kCatchup": 6000,
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
