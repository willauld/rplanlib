package rplanlib

import "math"
import "fmt"
import "os"

func startamount(amount float64, fra int, start int) float64 {
    if start > 70 {
        start = 70
	}
    if start < 62 {
        start = 62
	}
    if start < fra {
        return amount/(math.Pow(1.067,float64(fra-start)))
	}
	//  start >= fra must be
    return amount*(math.Pow(1.08,float64(start-fra)))
}

func do_SS_details(S map[string]string , bucket []float64) {
    sections := 0
    index := 0
    itype := "SocialSecurity"
    for k,v := range S.itype { //S.get( itype , {}).items()) 
                sections += 1
                r := match_retiree(k)
                if r == nil {
                    fmt.Printf("Error: [%s.%s] must match a retiree\n\t[%s.%s] should match [iam.%s] but there is no [iam.%s]\n", itype,k,itype,k,k,k)
                    os.Exit(1)
				}
                fraamount := v["amount"]
                fraage := v["FRA"]
                agestr := v["age"]
                dt := map[string]float64{"key": k, "amount": fraamount, "fra": fraage, "agestr": agestr, "ageAtStart": r["ageAtStart"], "currAge": r["age"], "throughAge": r["through"]}
                if fraamount < 0 && sections == 1 { // default spousal support in second slot 
                    self.SSinput[1] = dt
				} else {
                    self.SSinput[index] = dt
                    index += 1
				}
	}
            for i := range sections {
                //print("SSinput", self.SSinput)
                agestr = self.SSinput[i]["agestr"]
                firstage = agelist(agestr)
                disperseage = next(firstage)
                if i == 0 {
                    firstdisperseage = disperseage
                    firstdisperseyear = disperseage - self.SSinput[0]["ageAtStart"]
				}
                fraage = self.SSinput[i]["fra"]
                fraamount = self.SSinput[i]["amount"]
                ageAtStart = self.SSinput[i]["ageAtStart"]
                currAge = self.SSinput[i]["currAge"]
                if fraamount >= 0 {
                    // alter amount for start age vs fra (minus if before fra and + is after)
                    amount = self.startamount(fraamount, fraage, disperseage)
				} else {
                    //assert i == 1
                    name = self.SSinput[i]["key"]
                    if firstdisperseyear > disperseage - ageAtStart {
                        disperseage = firstdisperseyear + ageAtStart
                        agestr = "{}-".format(disperseage)
                        self.SSinput[i]["agestr"] = agestr 
                        fmt.Printf("Warning - Social Security spousal benefit can only be claimed\nafter the spouse claims benefits.\nPlease correct %s's SS age in the configuration file to '%s'.", name, agestr)
					} else if disperseage > fraage && firstdisperseyear != disperseage - ageAtStart {
                        if firstdisperseyear <= fraage - ageAtStart {
                            disperseage = fraage
                            agestr = "{}-".format(fraage)
                            self.SSinput[i]["agestr"] = agestr 
                            fmt.Printf("Warning - Social Security spousal benefits do not increase after FRA,\nresetting benefits start to FRA.\nPlease correct %s's SS age in the configuration file to '%s'.",name, agestr)
						} else {
                            disperseage = firstdisperseyear + ageAtStart
                            agestr = "{}-".format(disperseage)
                            self.SSinput[i]["agestr"] = agestr 
                            fmt.Printf("Warning - Social Security spousal benefits do not increase after FRA,\nresetting benefits start to spouse claim year.\nPlease correct %s's age in the configuration file to '%s'.",name, agestr)
						}
					}
                    fraamount = self.SSinput[0]["amount"]/2 // spousal benefit is 1/2 spouses at FRA 
                    // alter amount for start age vs fra (minus if before fra)
                    amount = self.startamount(fraamount, fraage, min(disperseage,fraage))
				}
                //print("FRA: %d, FRAamount: %6.0f, Age: %s, amount: %6.0f" % (fraage, fraamount, agestr, amount))
                self.SSinput[i]["bucket"] = [0] * self.numyr
                for age := range elist(agestr) {
                    year = age - ageAtStart //self.startage
                    if year < 0 {
                        continue
					} else if year >= self.numyr {
                        break
					} else {
                        adj_amount = amount * self.i_rate ** (age - currAge) //year
                        //print("age %d, year %d, bucket: %6.0f += amount %6.0f" %(age, year, bucket[year], adj_amount))
                        bucket[year] += adj_amount
                        self.SSinput[i]["bucket"][year] = adj_amount
					}
				}
			}
            if sections > 1 {
                //
                // Must fix up SS for period after one spouse dies
                //
                d := make([]float64, 2)
                d[0] = self.SSinput[0]["throughAge"]-self.SSinput[0]["ageAtStart"]
                d[1] = self.SSinput[1]["throughAge"]-self.SSinput[1]["ageAtStart"]
				if d[0] > d[1]  {
                firstToDie, secondToDie = 1, 0 
				} else  {
                firstToDie, secondToDie = 0, 1
				}
                for year := range a(d[firstToDie]+1, self.numyr) {
                    if self.SSinput[0]["bucket"][year] > self.SSinput[1]["bucket"][year] {
                        greater = self.SSinput[0]["bucket"][year]
					} else {
                        greater = self.SSinput[1]["bucket"][year]
					}
                    self.SSinput[firstToDie]["bucket"][year] = 0
                    self.SSinput[secondToDie]["bucket"][year] = greater
                    bucket[year] = greater
				}
			}
}