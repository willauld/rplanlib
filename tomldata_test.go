package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"testing"

	"github.com/willauld/rplanlib"
)

func TestGetTomlData(t *testing.T) {
	tests := []struct {
		toml []byte
		ipsm map[string]string
	}{
		{ // Case 0
			toml: []byte(`
title = "activeParams"
retirement_type = 'single' # joint, single or mseparate (married filing separately)
returns = 6.67		# return rate of investments, defaults to 6%
inflation = 3.5	# yearly inflation rate, defaults to 0%
maximize = "PlusEstate"
				`),
			ipsm: map[string]string{
				"setName":            "activeParams",
				"filingStatus":       "single",
				"eT_iRatePercent":    "3.5",
				"eT_rRatePercent":    "6.67",
				"eT_maximize":        "PlusEstate",
				"dollarsInThousands": "false",
			},
		},
		{ // Case 1
			toml: []byte(`
[iam.retiree1]  # iam (for each) is required in some joint cases (".xxx" use to match accounts IRA/roth)
primary = true  # retiree to have age listed first in the output (must choose one)
age = 54        # your current age
retire = 65     # age you plan to retire
through = 75    # age you want to plan through
definedContributionPlan = "54-65"
				`),
			ipsm: map[string]string{
				"key1":                             "retiree1",
				"eT_Age1":                          "54",
				"eT_RetireAge1":                    "65",
				"eT_PlanThroughAge1":               "75",
				"eT_DefinedContributionPlanStart1": "54",
				"eT_DefinedContributionPlanEnd1":   "65",
				"dollarsInThousands":               "false",
			},
		},
		{ // Case 2
			toml: []byte(`
[iam]  # iam (for each) is required in some joint cases (".xxx" use to match accounts IRA/roth)
primary = true  # retiree to have age listed first in the output (must choose one)
age = 54        # your current age
retire = 65     # age you plan to retire
through = 75    # age you want to plan through
definedContributionPlan = "54-65"
				`),
			ipsm: map[string]string{
				"key1":                             "nokey",
				"eT_Age1":                          "54",
				"eT_RetireAge1":                    "65",
				"eT_PlanThroughAge1":               "75",
				"eT_DefinedContributionPlanStart1": "54",
				"eT_DefinedContributionPlanEnd1":   "65",
				"dollarsInThousands":               "false",
			},
		},
		{ // Case 3
			toml: []byte(`
[iam]  # iam (for each) is required in some joint cases (".xxx" use to match accounts IRA/roth)
primary = true  # retiree to have age listed first in the output (must choose one)
age = 54        # your current age
retire = 65     # age you plan to retire
through = 75    # age you want to plan through
definedContributionPlan = "54-65"
[iam.retiree2]  # iam (for each) is required in some joint cases (".xxx" use to match accounts IRA/roth)
age = 55        # your current age
retire = 66     # age you plan to retire
through = 76    # age you want to plan through
definedContributionPlan = "55-66"
				`),
			ipsm: map[string]string{
				"key1":                             "nokey",
				"eT_Age1":                          "54",
				"eT_RetireAge1":                    "65",
				"eT_PlanThroughAge1":               "75",
				"eT_DefinedContributionPlanStart1": "54",
				"eT_DefinedContributionPlanEnd1":   "65",
				"key2":                             "retiree2",
				"eT_Age2":                          "55",
				"eT_RetireAge2":                    "66",
				"eT_PlanThroughAge2":               "76",
				"eT_DefinedContributionPlanStart2": "55",
				"eT_DefinedContributionPlanEnd2":   "66",
				"dollarsInThousands":               "false",
			},
		},
		{ // Case 4
			toml: []byte(`
[iam]  # iam (for each) is required in some joint cases (".xxx" use to match accounts IRA/roth)
age = 54        # your current age
retire = 65     # age you plan to retire
through = 75    # age you want to plan through
definedContributionPlan = "54-65"
[iam.retiree2]  # iam (for each) is required in some joint cases (".xxx" use to match accounts IRA/roth)
primary = true  # retiree to have age listed first in the output (must choose one)
age = 55        # your current age
retire = 66     # age you plan to retire
through = 76    # age you want to plan through
definedContributionPlan = "55-66"
				`),
			ipsm: map[string]string{
				"key1":                             "retiree2",
				"key2":                             "nokey",
				"eT_Age1":                          "55",
				"eT_Age2":                          "54",
				"eT_RetireAge1":                    "66",
				"eT_RetireAge2":                    "65",
				"eT_PlanThroughAge1":               "76",
				"eT_PlanThroughAge2":               "75",
				"eT_DefinedContributionPlanStart1": "55",
				"eT_DefinedContributionPlanStart2": "54",
				"eT_DefinedContributionPlanEnd1":   "66",
				"eT_DefinedContributionPlanEnd2":   "65",
				"dollarsInThousands":               "false",
			},
		},
		{ // Case 5
			toml: []byte(`
[iam]  # iam (for each) is required in some joint cases (".xxx" use to match accounts IRA/roth)
age = 54        # your current age
retire = 65     # age you plan to retire
through = 75    # age you want to plan through
definedContributionPlan = "54-65"
[iam.retiree2]  # iam (for each) is required in some joint cases (".xxx" use to match accounts IRA/roth)
primary = true  # retiree to have age listed first in the output (must choose one)
age = 55        # your current age
retire = 66     # age you plan to retire
through = 76    # age you want to plan through
definedContributionPlan = "55-66"
[SocialSecurity.retiree2]
FRA = 67            # your full retirement age (FRA) according to the IRS
amount =  20_000    # estimated yearly amount at Full Retirement Age (FRA); Assumes inflation, 85% taxed
age = "70-"         # period you expect to receive SS ("68-" indicates start at 68 and continue)

				`),
			ipsm: map[string]string{
				"key1":                             "retiree2",
				"key2":                             "nokey",
				"eT_Age1":                          "55",
				"eT_Age2":                          "54",
				"eT_RetireAge1":                    "66",
				"eT_RetireAge2":                    "65",
				"eT_PlanThroughAge1":               "76",
				"eT_PlanThroughAge2":               "75",
				"eT_DefinedContributionPlanStart1": "55",
				"eT_DefinedContributionPlanStart2": "54",
				"eT_DefinedContributionPlanEnd1":   "66",
				"eT_DefinedContributionPlanEnd2":   "65",
				"eT_PIA1":                          "20000", //20K
				"eT_PIA2":                          "",
				"eT_SS_Start1":                     "70",
				"eT_SS_Start2":                     "",
				"dollarsInThousands":               "false",
			},
		},
		{ // Case 6
			toml: []byte(`
[iam]  # iam (for each) is required in some joint cases (".xxx" use to match accounts IRA/roth)
age = 54        # your current age
retire = 65     # age you plan to retire
through = 75    # age you want to plan through
definedContributionPlan = "54-65"
[iam.retiree2]  # iam (for each) is required in some joint cases (".xxx" use to match accounts IRA/roth)
primary = true  # retiree to have age listed first in the output (must choose one)
age = 55        # your current age
retire = 66     # age you plan to retire
through = 76    # age you want to plan through
definedContributionPlan = "55-66"
[SocialSecurity]
FRA = 67            # your full retirement age (FRA) according to the IRS
amount =  20_000    # estimated yearly amount at Full Retirement Age (FRA); Assumes inflation, 85% taxed
age = "70-"         # period you expect to receive SS ("68-" indicates start at 68 and continue)

				`),
			ipsm: map[string]string{
				"key1":                             "retiree2",
				"key2":                             "nokey",
				"eT_Age1":                          "55",
				"eT_Age2":                          "54",
				"eT_RetireAge1":                    "66",
				"eT_RetireAge2":                    "65",
				"eT_PlanThroughAge1":               "76",
				"eT_PlanThroughAge2":               "75",
				"eT_DefinedContributionPlanStart1": "55",
				"eT_DefinedContributionPlanStart2": "54",
				"eT_DefinedContributionPlanEnd1":   "66",
				"eT_DefinedContributionPlanEnd2":   "65",
				"eT_PIA2":                          "20000", //20K
				"eT_PIA1":                          "",
				"eT_SS_Start2":                     "70",
				"eT_SS_Start1":                     "",
				"dollarsInThousands":               "false",
			},
		},
		{ // Case 7
			toml: []byte(`
[iam]  # iam (for each) is required in some joint cases (".xxx" use to match accounts IRA/roth)
primary = true  # retiree to have age listed first in the output (must choose one)
age = 54        # your current age
retire = 65     # age you plan to retire
through = 75    # age you want to plan through
definedContributionPlan = "54-65"
[IRA]
bal = 200_000       # current balance 
rate = 7.25        # defaults to global rate set above
#contrib = 0        # Annual contribution you will make for period (below)
inflation = true  # Will the contribution rise with inflation?
#period = '56-60'   # period you will be making the contributions
				`),
			ipsm: map[string]string{
				"key1":                             "nokey",
				"eT_Age1":                          "54",
				"eT_RetireAge1":                    "65",
				"eT_PlanThroughAge1":               "75",
				"eT_DefinedContributionPlanStart1": "54",
				"eT_DefinedContributionPlanEnd1":   "65",
				"eT_TDRA1":                         "200000", // 200k
				"eT_TDRA2":                         "",
				"eT_TDRA_Rate1":                    "7.25",
				"eT_TDRA_Rate2":                    "",
				"eT_TDRA_Contrib1":                 "",
				"eT_TDRA_Contrib2":                 "",
				"eT_TDRA_ContribStartAge1":         "",
				"eT_TDRA_ContribStartAge2":         "",
				"eT_TDRA_ContribEndAge1":           "",
				"eT_TDRA_ContribEndAge2":           "",
				"eT_TDRA_ContribInflate1":          "true",
				"dollarsInThousands":               "false",
			},
		},
		{ // Case 8
			toml: []byte(`
[iam]  # iam (for each) is required in some joint cases (".xxx" use to match accounts IRA/roth)
primary = true  # retiree to have age listed first in the output (must choose one)
age = 54        # your current age
retire = 65     # age you plan to retire
through = 75    # age you want to plan through
definedContributionPlan = "54-65"
[aftertax]
bal = 200_000       # current balance 
rate = 7.25        # defaults to global rate set above
#contrib = 0        # Annual contribution you will make for period (below)
inflation = true  # Will the contribution rise with inflation?
#period = '56-60'   # period you will be making the contributions
				`),
			ipsm: map[string]string{
				"key1":                             "nokey",
				"eT_Age1":                          "54",
				"eT_RetireAge1":                    "65",
				"eT_PlanThroughAge1":               "75",
				"eT_DefinedContributionPlanStart1": "54",
				"eT_DefinedContributionPlanEnd1":   "65",
				"eT_Aftatax":                       "200000", // 200k
				"eT_Aftatax_Rate":                  "7.25",
				"eT_Aftatax_Contrib":               "",
				"eT_Aftatax_ContribStartAge":       "",
				"eT_Aftatax_ContribEndAge":         "",
				"eT_Aftatax_ContribInflate":        "true",
				"dollarsInThousands":               "false",
			},
		},
		{ // Case 9
			toml: []byte(`
[iam]  # iam (for each) is required in some joint cases (".xxx" use to match accounts IRA/roth)
primary = true  # retiree to have age listed first in the output (must choose one)
age = 54        # your current age
retire = 65     # age you plan to retire
through = 75    # age you want to plan through
definedContributionPlan = "54-65"
[income]
amount = 1000
age = "63-67"
inflation = true
tax = true
				`),
			ipsm: map[string]string{
				"key1":                             "nokey",
				"eT_Age1":                          "54",
				"eT_RetireAge1":                    "65",
				"eT_PlanThroughAge1":               "75",
				"eT_DefinedContributionPlanStart1": "54",
				"eT_DefinedContributionPlanEnd1":   "65",
				"eT_Income1":                       "nokey",
				"eT_IncomeAmount1":                 "1000",
				"eT_IncomeStartAge1":               "63",
				"eT_IncomeEndAge1":                 "67",
				"eT_IncomeInflate1":                "true",
				"eT_IncomeTax1":                    "true",
				"dollarsInThousands":               "false",
			},
		},
		{ // Case 10
			toml: []byte(`
[iam]  # iam (for each) is required in some joint cases (".xxx" use to match accounts IRA/roth)
primary = true  # retiree to have age listed first in the output (must choose one)
age = 54        # your current age
retire = 65     # age you plan to retire
through = 75    # age you want to plan through
definedContributionPlan = "54-65"
[income]
amount = 1000
age = "63-67"
inflation = true
tax = true
[income.two]
amount = 1000
age = "63-67"
inflation = true
tax = true
				`),
			ipsm: map[string]string{
				"key1":                             "nokey",
				"eT_Age1":                          "54",
				"eT_RetireAge1":                    "65",
				"eT_PlanThroughAge1":               "75",
				"eT_DefinedContributionPlanStart1": "54",
				"eT_DefinedContributionPlanEnd1":   "65",
				"eT_Income1":                       "nokey",
				"eT_IncomeAmount1":                 "1000",
				"eT_IncomeStartAge1":               "63",
				"eT_IncomeEndAge1":                 "67",
				"eT_IncomeInflate1":                "true",
				"eT_IncomeTax1":                    "true",
				"eT_Income2":                       "two",
				"eT_IncomeAmount2":                 "1000",
				"eT_IncomeStartAge2":               "63",
				"eT_IncomeEndAge2":                 "67",
				"eT_IncomeInflate2":                "true",
				"eT_IncomeTax2":                    "true",
				"dollarsInThousands":               "false",
			},
		},
		{ // Case 11
			toml: []byte(`
title = "activeParams"

retirement_type = 'joint' # defaults to joint, could be single, joint (married) or mseparate (married filing separately)

returns = 6		# return rate of investments, defaults to 6%
inflation = 2.5	# yearly inflation rate, defaults to 0%

# what to optimize for? 'Spending' or spending 'PlusEstate', defaults to Spending
maximize = "Spending"
#maximize = "PlusEstate"

[iam.retiree1]  # iam (for each) is required in some joint cases (".xxx" use to match accounts IRA/roth)
primary = true  # retiree to have age listed first in the output (must choose one)
age = 54        # your current age
retire = 65     # age you plan to retire
through = 75    # age you want to plan through

[iam.retiree2]  # iam (for each) is required in some joint cases (".xxx" use to match accounts IRA/roth)
age = 54        # your current age
retire = 65     # age you plan to retire
through = 75    # age you want to plan through

[SocialSecurity.retiree1]
FRA = 67            # your full retirement age (FRA) according to the IRS
amount =  20_000    # estimated yearly amount at Full Retirement Age (FRA); Assumes inflation, 85% taxed
age = "70-"         # period you expect to receive SS ("68-" indicates start at 68 and continue)

[SocialSecurity.retiree2]
FRA = 67            # your full retirement age (FRA) according to the IRS
amount = -1         # -1 for default spousal benefit amount, amount at Full Retirement Age (FRA); Assumes inflation, 85% taxed
age = "70-"         # period you expect to receive SS ("68-" indicates start at 68 and continue)

[income.rental1]
amount = 1000
age = "63-67"
inflation = true
tax = true

[income.rental2]
amount = 2000
age = "62-70"
inflation = false
tax = true

[expense.exp1]
amount = 1000
age = "63-67"
inflation = true
tax = true

[expense.exp2]
amount = 2000
age = "62-70"
inflation = false
tax = true

[asset.ass1]
value = 100_000                 # current value of the asset
costAndImprovements = 20_000   # purchase price plus improvement cost
ageToSell = 73                  # age at which to sell the asset
owedAtAgeToSell = 10_000       # amount owed at time of sell (ageToSell)
primaryResidence = true         # Primary residence gets tax break
rate = 4                        # avg rate of return (defaults to global rate)
brokerageRate = 4               # brokerage fee percentage (defaults to 4%)

[asset.ass2]
value = 100_000                 # current value of the asset
costAndImprovements = 20_000   # purchase price plus improvement cost
ageToSell = 73                  # age at which to sell the asset
owedAtAgeToSell = 10_000       # amount owed at time of sell (ageToSell)
primaryResidence = false         # Primary residence gets tax break
rate = 6.0                        # avg rate of return (defaults to global rate)
#brokerageRate = 4               # brokerage fee percentage (defaults to 4%)

#[income.taxfreeNoneInflationAdjustedAnuity]
#amount = 3000      # yearly amount
#age = "65-70"      # period you expect to receive it
#inflation = false  # not inflation adjusted
#tax = false        # not federally taxable

#[income.InflationAdjustedAnuity]
#amount = 3000      # yearly amount
#age = "65-70"      # period you expect to receive it
#inflation = true   # inflation adjusted
#tax = true         # federally taxable

#[income.reversemortgage]
#amount = 12000      # yearly amount, 1000/mo
#age = '70-'         # period to receive payments
#inflation = false   # payment is not inflation adjusted
#tax = false         # payment/loan is not taxable

#[income.rental]
#amount = 5_000     # yearly amount
#age = "67-"        # period you expect to receive it
#inflation = true   # inflation adjusted
#tax = true         # federally taxable

#[asset.home]
#value = 550_000                 # current value of the asset
#costAndImprovements = 300_000   # purchase price plus improvement cost
#ageToSell = 72                  # age at which to sell the asset
#owedAtAgeToSell = 100_000       # amount owed at time of sell (ageToSell)
#primaryResidence = true         # Primary residence gets tax break
#rate = 4                        # avg rate of return (defaults to global rate)
#brokerageRate = 4               # brokerage fee percentage (defaults to 4%)

#[asset.rental]
#value = 250_000                 # current value of the asset
#costAndImprovements = 150_000   # purchase price plus improvement cost
#ageToSell = 72                  # age at which to sell the asset
#owedAtAgeToSell = 100_000       # amount owed at time of sell (ageToSell)
#primaryResidence = false        # Primary residence gets tax break
#rate = 4                        # avg rate of return (defaults to global rate)

#[min.income]   # used when maximize = "PlusEstate"
#amount = 45_000    # retirement first year income

#[max.income]       # used when maximize = "Spendable" (default)
#amount = 100_000   # retirement first year income

# pre-tax IRA accounts (TDRA)
[IRA.retiree1]
bal = 200_000       # current balance 
#rate = 7.25        # defaults to global rate set above
#contrib = 0        # Annual contribution you will make for period (below)
inflation = true  # Will the contribution rise with inflation?
#period = '56-60'   # period you will be making the contributions

[IRA.retiree2]
#bal = 0 #100_000      # current balance
#rate = 7.25        # defaults to global rate set above
contrib = 5_000        # Annual contribution you will make for period (below)
inflation = true  # Will the contribution rise with inflation?
period = '63-64'   # period you will be making the contributions

# roth IRA accounts (RothRA)
[roth.retiree1]
bal = 5_000       # current balance
#rate = 7.25        # defaults to global rate set above
#contrib = 0        # Annual contribution you will make for period (below)
#inflation = false  # Will the contribution rise with inflation?
#period = '56-60'   # period you will be making the contributions

[roth.retiree2]
bal = 20_000       # current balance
#rate = 7.25        # defaults to global rate set above
#contrib = 0        # Annual contribution you will make for period (below)
#inflation = false  # Will the contribution rise with inflation?
#period = '56-60'   # period you will be making the contributions

# after tax savings accounts (ATRSI)
[aftertax]
bal =   60_000    # current balance
#basis = 50_000	    # Contributions to total, for capital gains tax
#rate = 7.25        # defaults to global rate set above
contrib = 10_000        # Annual contribution you will make for period (below)
inflation = true  # Will the contribution rise with inflation?
period = '63-67'   # period you will be making the contributions
			`),
			ipsm: map[string]string{
				"setName":                    "activeParams",
				"filingStatus":               "joint",
				"key1":                       "retiree1",
				"key2":                       "retiree2",
				"eT_Age1":                    "54",
				"eT_Age2":                    "54",
				"eT_RetireAge1":              "65",
				"eT_RetireAge2":              "65",
				"eT_PlanThroughAge1":         "75",
				"eT_PlanThroughAge2":         "75",
				"eT_PIA1":                    "20000", //20K
				"eT_PIA2":                    "-1",
				"eT_SS_Start1":               "70",
				"eT_SS_Start2":               "70",
				"eT_TDRA1":                   "200000", // 200k
				"eT_TDRA2":                   "",
				"eT_TDRA_Rate1":              "",
				"eT_TDRA_Rate2":              "",
				"eT_TDRA_Contrib1":           "",
				"eT_TDRA_Contrib2":           "5000", // contribute 5k per year
				"eT_TDRA_ContribStartAge1":   "",
				"eT_TDRA_ContribStartAge2":   "63",
				"eT_TDRA_ContribEndAge1":     "",
				"eT_TDRA_ContribEndAge2":     "64",
				"eT_TDRA_ContribInflate1":    "true",
				"eT_TDRA_ContribInflate2":    "true",
				"eT_Roth1":                   "5000",
				"eT_Roth2":                   "20000", //20k
				"eT_Roth_Rate1":              "",
				"eT_Roth_Rate2":              "",
				"eT_Roth_Contrib1":           "",
				"eT_Roth_Contrib2":           "",
				"eT_Roth_ContribStartAge1":   "",
				"eT_Roth_ContribStartAge2":   "",
				"eT_Roth_ContribEndAge1":     "",
				"eT_Roth_ContribEndAge2":     "",
				"eT_Aftatax":                 "60000", //60k
				"eT_Aftatax_Rate":            "",
				"eT_Aftatax_Contrib":         "10000", //10K
				"eT_Aftatax_ContribStartAge": "63",
				"eT_Aftatax_ContribEndAge":   "67",
				"eT_Aftatax_ContribInflate":  "true",

				"eT_iRatePercent":    "2.5",
				"eT_rRatePercent":    "6",
				"eT_maximize":        "Spending", // or "PlusEstate"
				"dollarsInThousands": "false",

				//prototype entries below
				"eT_Income1":         "rental1",
				"eT_IncomeAmount1":   "1000",
				"eT_IncomeStartAge1": "63",
				"eT_IncomeEndAge1":   "67",
				"eT_IncomeInflate1":  "true",
				"eT_IncomeTax1":      "true",

				//prototype entries below
				"eT_Income2":         "rental2",
				"eT_IncomeAmount2":   "2000",
				"eT_IncomeStartAge2": "62",
				"eT_IncomeEndAge2":   "70",
				"eT_IncomeInflate2":  "false",
				"eT_IncomeTax2":      "true",

				//prototype entries below
				"eT_Expense1":         "exp1",
				"eT_ExpenseAmount1":   "1000",
				"eT_ExpenseStartAge1": "63",
				"eT_ExpenseEndAge1":   "67",
				"eT_ExpenseInflate1":  "true",
				"eT_ExpenseTax1":      "true", //ignored, or should be

				//prototype entries below
				"eT_Expense2":         "exp2",
				"eT_ExpenseAmount2":   "2000",
				"eT_ExpenseStartAge2": "62",
				"eT_ExpenseEndAge2":   "70",
				"eT_ExpenseInflate2":  "false",
				"eT_ExpenseTax2":      "true", //ignored, or should be

				//prototype entries below
				"eT_Asset1":                    "ass1",
				"eT_AssetValue1":               "100000",
				"eT_AssetAgeToSell1":           "73",
				"eT_AssetCostAndImprovements1": "20000",
				"eT_AssetOwedAtAgeToSell1":     "10000",
				"eT_AssetPrimaryResidence1":    "true",
				"eT_AssetRRatePercent1":        "4",
				"eT_AssetBrokeragePercent1":    "4",

				//prototype entries below
				"eT_Asset2":                    "ass2",
				"eT_AssetValue2":               "100000",
				"eT_AssetAgeToSell2":           "73",
				"eT_AssetCostAndImprovements2": "20000",
				"eT_AssetOwedAtAgeToSell2":     "10000",
				"eT_AssetPrimaryResidence2":    "false",
				"eT_AssetRRatePercent2":        "6", // python defaults to global rate
				"eT_AssetBrokeragePercent2":    "",
			},
		},
		{ // Case 12
			toml: []byte(`
[iam]  # iam (for each) is required in some joint cases (".xxx" use to match accounts IRA/roth)
primary = true  # retiree to have age listed first in the output (must choose one)
age = 54        # your current age
retire = 65     # age you plan to retire
through = 75    # age you want to plan through
definedContributionPlan = "54-65"
[IRA]
bal = 200_000       # current balance 
rate = 7.25        # defaults to global rate set above
#contrib = 0        # Annual contribution you will make for period (below)
inflation = true  # Will the contribution rise with inflation?
#period = '56-60'   # period you will be making the contributions
[max.income]
amount = 100_000
				`),
			ipsm: map[string]string{
				"key1":                             "nokey",
				"eT_Age1":                          "54",
				"eT_RetireAge1":                    "65",
				"eT_PlanThroughAge1":               "75",
				"eT_DefinedContributionPlanStart1": "54",
				"eT_DefinedContributionPlanEnd1":   "65",
				"eT_TDRA1":                         "200000", // 200k
				"eT_TDRA2":                         "",
				"eT_TDRA_Rate1":                    "7.25",
				"eT_TDRA_Rate2":                    "",
				"eT_TDRA_Contrib1":                 "",
				"eT_TDRA_Contrib2":                 "",
				"eT_TDRA_ContribStartAge1":         "",
				"eT_TDRA_ContribStartAge2":         "",
				"eT_TDRA_ContribEndAge1":           "",
				"eT_TDRA_ContribEndAge2":           "",
				"eT_TDRA_ContribInflate1":          "true",
				"eT_MaxIncome":                     "100000",
				"dollarsInThousands":               "false",
			},
		},
		{ // Case 13
			toml: []byte(`
[iam]  # iam (for each) is required in some joint cases (".xxx" use to match accounts IRA/roth)
primary = true  # retiree to have age listed first in the output (must choose one)
age = 54        # your current age
retire = 65     # age you plan to retire
through = 75    # age you want to plan through
definedContributionPlan = "54-65"
[IRA]
bal = 200_000       # current balance 
rate = 7.25        # defaults to global rate set above
#contrib = 0        # Annual contribution you will make for period (below)
inflation = true  # Will the contribution rise with inflation?
#period = '56-60'   # period you will be making the contributions
[min.income]
amount = 100_000
				`),
			ipsm: map[string]string{
				"key1":                             "nokey",
				"eT_Age1":                          "54",
				"eT_RetireAge1":                    "65",
				"eT_PlanThroughAge1":               "75",
				"eT_DefinedContributionPlanStart1": "54",
				"eT_DefinedContributionPlanEnd1":   "65",
				"eT_TDRA1":                         "200000", // 200k
				"eT_TDRA2":                         "",
				"eT_TDRA_Rate1":                    "7.25",
				"eT_TDRA_Rate2":                    "",
				"eT_TDRA_Contrib1":                 "",
				"eT_TDRA_Contrib2":                 "",
				"eT_TDRA_ContribStartAge1":         "",
				"eT_TDRA_ContribStartAge2":         "",
				"eT_TDRA_ContribEndAge1":           "",
				"eT_TDRA_ContribEndAge2":           "",
				"eT_TDRA_ContribInflate1":          "true",
				"eT_DesiredIncome":                 "100000",
				"dollarsInThousands":               "false",
			},
		},
		{ // Case 14
			toml: []byte(`
title = "activeParams"
retirement_type = 'joint'
maximize = "PlusEstate"
[iam.joe]
age = 54
retire = 65
through = 75

[IRA.joe]
bal = 200_000
				`),
			ipsm: map[string]string{
				"setName":            "activeParams",
				"filingStatus":       "joint",
				"key1":               "joe",
				"eT_Age1":            "54",
				"eT_RetireAge1":      "65",
				"eT_PlanThroughAge1": "75",
				"eT_TDRA1":           "200000", // 200k
				"eT_maximize":        "PlusEstate",
				"dollarsInThousands": "false",

				//"eT_maximize":        "PlusEstate",
			},
		},
	}
	onlyOnce := 0
	for i, elem := range tests {
		//fmt.Printf("------ Case %d -----------\n", i)
		f, err := ioutil.TempFile("", "tom")
		if err != nil {
			t.Errorf("TestGetTomlData case %d: %s", i, err)
		}
		tfile := f.Name()
		defer os.Remove(tfile) // clean up

		_, err = f.Write(elem.toml)
		if err != nil {
			t.Errorf("TestGetTomlData case %d: %s", i, err)
		}
		f.Close()
		if err != nil {
			t.Errorf("TestGetTomlData case %d: %s", i, err)
		}
		ms, err := getInputStringsMapFromToml(tfile)
		if err != nil {
			t.Errorf("TestGetTomlData case %d: ms is nil: %s", i, err)
			continue
		}
		if ms == nil {
			//This should not happen withour err being set!!!!
			t.Errorf("TestGetTomlData case %d: ms is nil: %s", i, "err is empty :: ms should never be nil unless err has something!!!")
			continue
		}
		//ipsm is just a subset of the values in ms to compare
		//if len(elem.ipsm) != len(*ms) {
		//	t.Errorf("TestGetTomlData case %d: len(ms): %d != len(ipms): %d",
		//		i, len(*ms), len(elem.ipsm))
		//}
		for k, v := range *ms {
			foundIssue := false
			if elem.ipsm[k] != (*ms)[k] {
				t.Errorf("TestGetTomlData case %d: For '%s', expected: '%s', but found: '%s'", i, k, elem.ipsm[k], v)
				foundIssue = true
			}
			if foundIssue && onlyOnce < 1 {
				onlyOnce++
				fmt.Printf("TestGetTomlData fails sometimes because of the somewhere random ordering of Map access (and TomlTree access) so these failures are sometimes false positives. If you run the tests a few times and the go away in some of the runs its OK. Need to improve this test checking to eliminate this false positive issue.")
			}
		}
		for _, v := range rplanlib.InputStrDefs {
			r, ok := (*ms)[v]
			if !ok {
				t.Errorf("TestGetTomlData case %d: missing ms[%s]", i, v)
			}
			if r != "" {
				// fmt.Printf("    %s: '%s'\n", v, r)
			}
		}
		for x := 1; x < rplanlib.MaxStreams+1; x++ {
			for _, v := range rplanlib.InputStreamStrDefs {
				r, ok := (*ms)[fmt.Sprintf("%s%d", v, x)]
				if !ok {
					t.Errorf("TestGetTomlData case %d: missing ms[%s]", i, v)
				}
				if r != "" {
					// fmt.Printf("    %s%d: '%s'\n", v, x, r)
				}
			}
		}
	}
}
