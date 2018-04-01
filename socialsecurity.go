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

type ssI struct {
	fraamount  int
	fraage     int
	startSSAge int
	endAge     int
	key        string
	ageAtStart int
	currAge    int
	bucket     []float64
}

func processSS(ip *InputParams, warnList *warnErrorList) (SS, SS1, SS2 []float64, tags []string) {

	//fmt.Printf("PIA1: %d, PIA2: %d\n", ip.PIA1, ip.PIA2)
	ssi := make([]ssI, 2)
	if ip.PIA1 <= 0 && ip.PIA2 <= 0 {
		//e := fmt.Errorf("processSS: both PIA1: %d and PIA2: %d non-positive", ip.PIA1, ip.PIA2)
		return nil, nil, nil, nil
	}
	SS = make([]float64, ip.Numyr) // = [0] * self.numyr
	tags = make([]string, 2)
	tags[0] = "combined"
	tags[1] = ip.MyKey1

	index := 0
	sections := 1
	dt := ssI{
		fraamount:  ip.PIA1,            // fraamount := v["amount"]
		fraage:     fra(ip.Age1),       // fraage := v["FRA"]
		startSSAge: ip.SSStart1,        // agestr := v["age"]
		endAge:     ip.PlanThroughAge1, // agestr := v["age"]
		key:        ip.MyKey1,
		ageAtStart: ip.Age1 + ip.PrePlanYears,
		currAge:    ip.Age1,
	}
	if dt.fraamount <= 0 && ip.FilingStatus != Joint {
		return nil, nil, nil, nil
	}
	if dt.fraamount <= 0 { // place default spousal support in second slot
		ssi[1] = dt
	} else {
		ssi[index] = dt
		index++
	}
	if ip.FilingStatus == Joint && ip.Age2 != 0 && ip.SSStart2 > 0 {
		sections = 2
		dt = ssI{
			fraamount:  ip.PIA2,            // fraamount := v["amount"]
			fraage:     fra(ip.Age2),       // fraage := v["FRA"]
			startSSAge: ip.SSStart2,        // agestr := v["age"]
			endAge:     ip.PlanThroughAge2, // agestr := v["age"]
			key:        ip.MyKey2,
			ageAtStart: ip.Age2 + ip.PrePlanYears,
			currAge:    ip.Age2,
		}
		ssi[index] = dt
		tags = append(tags, ip.MyKey2)
	}
	//fmt.Printf("ssi[0]: %#v\n\n", ssi[0])
	//fmt.Printf("ssi[1]: %#v\n\n", ssi[1])
	//
	// spousal benefit can not start before SS primary starts taking SS
	//
	firstdisperseyear := ssi[0].startSSAge - ssi[0].ageAtStart
	//fmt.Printf("firstdispersyear: %d\n", firstdisperseyear)

	var amount float64
	for i := 0; i < sections; i++ {
		disperseage := ssi[i].startSSAge
		//fraage := ssi[i].fraage
		//fmt.Printf("FRA age[%d]: %d\n", i, ssi[i].fraage)
		//fraamount := ssi[i].fraamount
		//ageAtStart := ssi[i].ageAtStart
		//currAge := ssi[i].currAge
		if ssi[i].fraamount > 0 { // TODO check if this needs to be able to equal zero ; FIXME maybe explicitly check for zero and return nils
			// alter amount for start age vs fra (minus if before fra and + is after)
			amount = adjPIA(float64(ssi[i].fraamount), ssi[i].fraage, disperseage)
		} else {
			if i != 1 {
				e := fmt.Errorf("Error: Assert i == 1 failed (1212121)")
				panic(e)
			}
			name := ssi[i].key
			if firstdisperseyear > ssi[i].startSSAge-ssi[i].ageAtStart {
				disperseage = firstdisperseyear + ssi[i].ageAtStart
				str := fmt.Sprintf("Warning - Social Security spousal benefit can only be claimed\n\tafter the spouse claims benefits.\n\tPlease correct %s's SS age in the configuration file to '%d'.", name, disperseage)
				warnList.AppendWarning(str)
			} else if ssi[i].startSSAge > ssi[i].fraage && firstdisperseyear != ssi[i].startSSAge-ssi[i].ageAtStart {
				if firstdisperseyear <= ssi[i].fraage-ssi[i].ageAtStart {
					disperseage = ssi[i].fraage
					str := fmt.Sprintf("Warning - Social Security spousal benefits do not increase after FRA,\n\tresetting benefits start to FRA.\n\tPlease correct %s's SS age in the configuration file to '%d'.", name, ssi[i].fraage)
					warnList.AppendWarning(str)
				} else {
					disperseage = firstdisperseyear + ssi[i].ageAtStart
					str := fmt.Sprintf("Warning - Social Security spousal benefits do not increase after FRA,\n\tresetting benefits start to spouse claim year.\n\tPlease correct %s's age in the configuration file to '%d'.", name, disperseage)
					warnList.AppendWarning(str)
				}
			}
			//fraamount := ssi[0].fraamount / 2 // spousal benefit is 1/2 spouses at FRA
			// alter amount for start age vs fra (minus if before fra)
			amount = adjPIA(float64(ssi[0].fraamount)/2, ssi[i].fraage, intMin(disperseage, ssi[i].fraage))
		}
		ssi[i].bucket = make([]float64, ip.Numyr) // = [0] * self.numyr
		endage := ip.Numyr + ssi[i].ageAtStart
		//fmt.Printf("section: %d, disperseage: %d\n", i, disperseage)
		for age := disperseage; age < endage; age++ {
			year := age - ssi[i].ageAtStart //self.startage
			//fmt.Printf("year: %d, age: %d, ageAtStart: %d, name: %s\n", year, age, ssi[i].ageAtStart, ssi[i].key)
			if year < 0 {
				// ERROR if ever happens
				fmt.Printf("ERROR - this should never happen. local code 11111\n")
				fmt.Printf("age: %d, year: %d, endage: %d, ageAtStart: %d, name: %s\n", age, year, endage, ssi[i].ageAtStart, ssi[i].key)
				fmt.Printf("ssi[%d]: %#v\n", i, ssi[i])
				break
			} else if year >= ip.Numyr {
				// ERROR if ever happens
				fmt.Printf("ERROR - this should never happen. local code 22222\n\tage: %d, year: %d, ip.numyr: %d, endPlan: %d, startPlan: %d\n", age, year, ip.Numyr, ip.EndPlan, ip.StartPlan)
				break
			} else {
				adjAmount := amount * math.Pow(ip.IRate, float64(age-ssi[i].currAge)) //year
				//print("age %d, year %d, SS: %6.0f += amount %6.0f" %(age, year, SS[year], adj_amount))
				SS[year] += adjAmount
				ssi[i].bucket[year] = adjAmount
			}
		}
		if ssi[i].key == ip.MyKey1 {
			ip.SSStart1 = disperseage
		} else {
			ip.SSStart2 = disperseage
		}
	}
	if sections > 1 {
		//
		// Must fix up SS for period after one spouse dies
		//
		d := make([]int, 2)
		d[0] = ssi[0].endAge - ssi[0].ageAtStart
		d[1] = ssi[1].endAge - ssi[1].ageAtStart
		firstToDie, secondToDie := 0, 1
		if d[0] > d[1] {
			firstToDie, secondToDie = 1, 0
		}
		for year := d[firstToDie] + 1; year < ip.Numyr; year++ {
			greater := ssi[1].bucket[year]
			if ssi[0].bucket[year] > ssi[1].bucket[year] {
				greater = ssi[0].bucket[year]
			}
			ssi[firstToDie].bucket[year] = 0
			ssi[secondToDie].bucket[year] = greater
			SS[year] = greater
		}
	}
	//fmt.Printf("ssi[0]: %#v\n", ssi[0])
	//fmt.Printf("ssi[1]: %#v\n", ssi[1])
	SS1 = ssi[0].bucket
	SS2 = ssi[1].bucket
	if ssi[0].key != ip.MyKey1 {
		SS1 = ssi[1].bucket
		SS2 = ssi[0].bucket
	}
	return SS, SS1, SS2, tags
}
