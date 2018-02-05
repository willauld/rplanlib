package rplanlib

import "math"
import "fmt"
import "time"

/*
* Social Security Full Retirement Age (FRA) table
* Year of Birth *	Full Retirement Age
*	1937 or earlier		65
*	1938				65 and 2 months
*	1939				65 and 4 months
*	1940				65 and 6 months < 1940 65, 1940 >= 66
*	1941				65 and 8 months
*	1942				65 and 10 months
*	1943--1954			66
*	1955				66 and 2 months
*	1956				66 and 4 months
*	1957				66 and 6 months < 1957 66, 1957 >= 67
*	1958				66 and 8 months
*	1959				66 and 10 months
*	1960 and later		67
 */

// fra returns the retiree social security fra given current age
func fra(currAge int) int {
	thisYear := time.Now().Year()
	born := thisYear - currAge
	if born < 1940 {
		return 65
	} else if born < 1957 {
		return 66
	}
	return 67
}

// adjPIA returns an adjusted PIA amount a retiree would recieve based on their PIA, proposed starting age and fra
func adjPIA(PIA float64, fra int, startAge int) float64 {
	if startAge > 70 {
		startAge = 70
	}
	if startAge < 62 {
		startAge = 62
	}
	if startAge < fra {
		return PIA / (math.Pow(1.067, float64(fra-startAge)))
	}
	//  start >= fra must be
	return PIA * (math.Pow(1.08, float64(startAge-fra)))
}

type SSInput struct {
	fraamount  int
	fraage     int
	startAge   int
	endAge     int
	key        string
	ageAtStart int
	currAge    int
	bucket     []float64
}

func processSS(ip InputParams, r []retiree, iRate float64) (SS, SS1, SS2 []float64, e error) {

	SSinput := make([]SSInput, 2)
	if ip.PIA1 < 0 && ip.PIA2 < 0 {
		e := fmt.Errorf("processSS: both PIA1: %d and PIA2: %d are negitive", ip.PIA1, ip.PIA2)
		return nil, nil, nil, e
	}
	SS = make([]float64, ip.numyr) // = [0] * self.numyr

	index := 0
	sections := 1
	dt := SSInput{
		fraamount:  ip.PIA1,            // fraamount := v["amount"]
		fraage:     fra(ip.age1),       // fraage := v["FRA"]
		startAge:   ip.SSStart1,        // agestr := v["age"]
		endAge:     ip.planThroughAge1, // agestr := v["age"]
		key:        ip.myKey1,
		ageAtStart: r[0].ageAtStart,
		currAge:    ip.age1,
	}
	if dt.fraamount < 0 { // place default spousal support in second slot
		SSinput[1] = dt
	} else {
		SSinput[index] = dt
		index++
	}
	if ip.filingStatus == "joint" {

		sections = 2
		dt = SSInput{
			fraamount:  ip.PIA2,            // fraamount := v["amount"]
			fraage:     fra(ip.age2),       // fraage := v["FRA"]
			startAge:   ip.SSStart2,        // agestr := v["age"]
			endAge:     ip.planThroughAge2, // agestr := v["age"]
			key:        ip.myKey2,
			ageAtStart: r[1].ageAtStart,
			currAge:    ip.age2,
		}
		SSinput[index] = dt
	}
	//fmt.Printf("SSinput[0]: %v\n", SSinput[0])
	//fmt.Printf("SSinput[1]: %v\n", SSinput[1])
	var disperseage int
	var firstdisperseyear int
	for i := 0; i < sections; i++ {
		disperseage = SSinput[i].startAge
		if i == 0 {
			firstdisperseyear = disperseage - SSinput[0].ageAtStart
		}
		fraage := SSinput[i].fraage
		fraamount := SSinput[i].fraamount
		ageAtStart := SSinput[i].ageAtStart
		currAge := SSinput[i].currAge
		var amount float64
		if fraamount >= 0 {
			// alter amount for start age vs fra (minus if before fra and + is after)
			amount = adjPIA(float64(fraamount), fraage, disperseage)
		} else {
			//assert i == 1
			name := SSinput[i].key
			if firstdisperseyear > disperseage-ageAtStart {
				disperseage = firstdisperseyear + ageAtStart
				fmt.Printf("Warning - Social Security spousal benefit can only be claimed\n\tafter the spouse claims benefits.\n\tPlease correct %s's SS age in the configuration file to '%d'.\n", name, disperseage)
			} else if disperseage > fraage && firstdisperseyear != disperseage-ageAtStart {
				if firstdisperseyear <= fraage-ageAtStart {
					disperseage = fraage
					fmt.Printf("Warning - Social Security spousal benefits do not increase after FRA,\n\tresetting benefits start to FRA.\n\tPlease correct %s's SS age in the configuration file to '%d'.\n", name, fraage)
				} else {
					disperseage = firstdisperseyear + ageAtStart
					fmt.Printf("Warning - Social Security spousal benefits do not increase after FRA,\n\tresetting benefits start to spouse claim year.\n\tPlease correct %s's age in the configuration file to '%d'.\n", name, disperseage)
				}
			}
			fraamount = SSinput[0].fraamount / 2 // spousal benefit is 1/2 spouses at FRA
			// alter amount for start age vs fra (minus if before fra)
			amount = adjPIA(float64(fraamount), fraage, intMin(disperseage, fraage))
		}
		SSinput[i].bucket = make([]float64, ip.numyr) // = [0] * self.numyr
		endage := ip.numyr + SSinput[i].ageAtStart
		//for age := SSinput[i].startAge; age < endage; age++
		for age := disperseage; age < endage; age++ {
			year := age - SSinput[i].ageAtStart //self.startage
			if year < 0 {
				// ERROR if ever happens
				fmt.Printf("ERROR - this should never happen. local code 11111\n")
				continue
			} else if year >= ip.numyr {
				// ERROR if ever happens
				fmt.Printf("ERROR - this should never happen. local code 22222\n\tage: %d, year: %d, ip.numyr: %d, endPlan: %d, startPlan: %d\n", age, year, ip.numyr, ip.endPlan, ip.startPlan)
				break
			} else {
				adjAmount := amount * math.Pow(iRate, float64(age-currAge)) //year
				//print("age %d, year %d, SS: %6.0f += amount %6.0f" %(age, year, SS[year], adj_amount))
				SS[year] += adjAmount
				SSinput[i].bucket[year] = adjAmount
			}
		}
	}
	if sections > 1 {
		//
		// Must fix up SS for period after one spouse dies
		//
		d := make([]int, 2)
		d[0] = SSinput[0].endAge - SSinput[0].ageAtStart
		d[1] = SSinput[1].endAge - SSinput[1].ageAtStart
		firstToDie, secondToDie := 0, 1
		if d[0] > d[1] {
			firstToDie, secondToDie = 1, 0
		}
		for year := d[firstToDie] + 1; year < ip.numyr; year++ {
			greater := SSinput[1].bucket[year]
			if SSinput[0].bucket[year] > SSinput[1].bucket[year] {
				greater = SSinput[0].bucket[year]
			}
			SSinput[firstToDie].bucket[year] = 0
			SSinput[secondToDie].bucket[year] = greater
			SS[year] = greater
		}
	}
	//fmt.Printf("SSinput[0]: %v\n", SSinput[0])
	//fmt.Printf("SSinput[1]: %v\n", SSinput[1])
	SS1 = SSinput[0].bucket
	SS2 = SSinput[1].bucket
	if SSinput[0].key != ip.myKey1 {
		SS1 = SSinput[1].bucket
		SS2 = SSinput[0].bucket
	}
	return SS, SS1, SS2, nil
}
