package rplanlib


type modelSpecs struct {
    InputParams map[string]string
    vindx vectorVarIndex

    // The following was through 'S'
    numyr int
    maximize string // "Spending" or "PlusEstate"
    accounttable []map[string]string
    accmap map[string]int
    retiree []map[string]string

    income []float64
    SS []float64
    expenses []float64
    taxed []float64
    asset_sale []float64
    cg_asset_taxed []float64

    iRate float64
    rRate float64

    min float64
    max float64

    preplanyears int


} 

func NewModelSpecs(vindx vectorVarIndex/*, 
                    taxtable, capgainstable, penalty, stded, SS_taxable, verbose, no_TDRA_ROTHRA_DEPOSITS*/
                                        ) modelSpecs {
    return modelSpecs{
    InputParams: map[string]string
    vindx: vindx,
    numyr: int
    maximize: string // "Spending" or "PlusEstate"
    accounttable: []map[string]string
    accmap: map[string]int
    retiree: []map[string]string
    income: []float64
    SS: []float64
    expenses: []float64
    taxed: []float64
    asset_sale: []float64
    cg_asset_taxed: []float64
    iRate: float64
    rRate: float64
    min: float64
    max: float64
    }
}
/*
class lp_constraint_model:
    def __init__(self, S, vindx, taxtable, capgainstable, penalty, stded, SS_taxable, verbose, no_TDRA_ROTHRA_DEPOSITS):
        self.S = S
        self.var_index = vindx
        self.taxtable = taxtable
        self.cgtaxtable = capgainstable
        self.penalty = penalty
        self.stded = stded
        self.ss_taxable = SS_taxable
        self.verbose = verbose
        self.noTdraRothraDeposits = no_TDRA_ROTHRA_DEPOSITS

*/

// BuildModel for:
// Minimize: c^T * x
// Subject to: A_ub * x <= b_ub
// all vars positive
func (ms modelSpecs) BuildModel() ([]float64, [][]float64, []float64) {
    
        // TODO integrate the following assignments into the code and remove them
        //S = ms.S
        taxtable = ms.taxtable
        capgainstable = ms.cgtaxtable
        penalty = ms.penalty
        stded = ms.stded
        SS_taxable = ms.ss_taxable
    
        nvars = ms.vindx.Vsize
        A = make([][]float64, 0)
        b = make([]float64, 0)
        c = make([]float64, nvars)
    
        //
        // Add objective function (S1') becomes (R1') if PlusEstate is added
        //
        for year = 0; year<ms.numyr; year++ {
            c[ms.vindx.S(year)] = -1
        }
        //
        // Add objective function tax bracket forcing function (EXPERIMENTAL)
        //
        for year = 0; year<ms.numyr; year++ {
            for k := 0; k<len(taxtable); k++ {
                // multiplies the impact of higher brackets opposite to
                // optimization the intent here is to pressure higher 
                // brackets more and pack the lower brackets
                c[ms.vindx.X(year,k)] = k/10 
            }
        }
        //
        // Adder objective function (R1') when PlusEstate is added
        //
        if ms.maximize == "PlusEstate" {
            for j := 0; j<len(ms.accounttable); j++ {
                c[ms.vindx.B(ms.numyr,j)] = -1*ms.accounttable[j]["estateTax"] // account discount rate
            }
            print("\nConstructing Spending + Estate Model:\n")
        } else {
            print("\nConstructing Spending Model:\n")
            startamount = 0
            for j := 0; j<len(ms.accounttable); j++ {
                startamount += ms.accounttable[j]["bal"]
            }
            balancer = 1/(startamount) 
            for j := 0; j<len(ms.accounttable); j++ {
                c[ms.vindx.B(ms.numyr,j)] = -1*balancer *ms.accounttable[j]["estateTax"] // balance and discount rate
            }
        }
        //
        // Add constraint (2')
        //
        for year = 0; year<ms.numyr; year++ {
            row = [0] * nvars
            for j := 0; j<len(ms.accounttable); j++ {
                p = 1
                if ms.accounttable[j]["acctype"] != "aftertax"{
                    if S.apply_early_penalty(year,ms.accounttable[j]["mykey"]){
                        p = 1-penalty
                    }
                }
                row[ms.vindx.W(year,j)] = -1*p 
            }
            for k:=0;  k<len(taxtable);k++ {
                row[ms.vindx.X(year,k)] = taxtable[k][2] // income tax
            } 
            if ms.accmap["aftertax"] > 0 {
                for l:=0; l<len(capgainstable); l++ {
                    row[ms.vindx.Y(year,l)] = capgainstable[l][2] // cap gains tax
                } 
                for j := 0; j<len(ms.accounttable); j++ {
                    row[ms.vindx.D(year,j)] = 1
                }
            }
            row[ms.vindx.S(year)] = 1
            A = append(A, row)
            b= append(b, ms.income[year] + ms.SS[year] - ms.expenses[year])
        }
        //
        // Add constraint (3a')
        //
        for year = 0; year<ms.numyr-1; year++ {
            row = [0] * nvars
            row[ms.vindx.S(year+1)] = 1
            row[ms.vindx.S(year)] = -1*ms.i_rate
            A = append(A, row)
            b = append(b, 0)
        }
        //
        // Add constraint (3b')
        //
        for year = 0; year<ms.numyr-1; year++ {
            row = [0] * nvars
            row[ms.vindx.S(year)] = ms.i_rate
            row[ms.vindx.S(year+1)] = -1
            A = append(A, row)
            b = append(b, 0)
        }
        //
        // Add constrant (4') rows - not needed if [desired.income] is not defined in input
        //
        if ms.min != 0{
            for year:=0; year<1; year++ { // Only needs setting at the beginning
                row = [0] * nvars
                row[ms.vindx.S(year)] = -1
                A = append(A, row)
                b = append(b, - ms.min )     // [- d_i]
            }
        }
    
        //
        // Add constraints for (5') rows - not added if [max.income] is 
        // not defined in input
        //
        if ms.max != 0{
            for year:=0; year<1; year++ { // Only needs to be set at the beginning
                row = [0] * nvars
                row[ms.vindx.S(year)] = 1
                A = append(A, row)
                b = append(b, ms.max )     // [ dm_i]
            }
        }
    
        //
        // Add constaints for (6') rows
        //
        for year = 0; year<ms.numyr; year++ {
            row = [0] * nvars
            for j := 0; j<len(ms.accounttable); j++ {
                if ms.accounttable[j]["acctype"] != "aftertax"{
                    row[ms.vindx.D(year,j)] = 1
                } 
            }
            A = append(A, row)
            //b+=[min(ms.income[year],ms.maxContribution(year,None))] 
            // using ms.taxed rather than ms.income because income could
            // include non-taxed anueities that don't count.
            b = append(b, min(ms.taxed[year],S.maxContribution(year,None)))
        }
        //
        // Add constaints for (7') rows
        //
        for year = 0; year<ms.numyr; year++ {
            // TODO this is not needed when there is only one retiree
            for _, v := range S.retiree {
                row = [0] * nvars
                for j := 0; j<len(ms.accounttable); j++ {
                    if v["mykey"] == ms.accounttable[j]["mykey"] {
                        // ["acctype"] != "aftertax": no "mykey" in aftertax
                        // (this will either break or just not match - we 
                        // will see)
                        row[ms.vindx.D(year,j)] = 1
                    } 
                }
                A = append(A, row)
                b = append(b, S.maxContribution(year,v["mykey"]))
            }
        } 
        //
        // Add constaints for (8') rows
        //
        for year = 0; year<ms.numyr; year++ {
            for j := 0; j<len(ms.accounttable); j++ {
                v = ms.accounttable[j].get("contributions", nil)
                if v != nil {
                    if v[year] > 0{
                        row = [0] * nvars
                        row[ms.vindx.D(year,j)] = -1
                        A = append(A, row)
                        b = append(b, -1*v[year])
                    }
                } 
            }
        }
        //
        // Add constaints for (9') rows
        //
        for year = 0; year<ms.numyr; year++ {
            for j := 0; j<min(2, len(ms.accounttable)); j++ {
                // at most the first two accounts are type IRA w/ 
                // RMD requirement
                if ms.accounttable[j]["acctype"] == "IRA"{
                    ownerage = S.account_owner_age(year, ms.accounttable[j])
                    if ownerage >= 70{
                        row = [0] * nvars
                        row[ms.vindx.D(year,j)] = 1
                        A = append(A, row)
                        b = append(b, 0)
                    }
                }
            } 
        }
        //
        // Add constaints for (N') rows
        //
        if ms.noTdraRothraDeposits{
            for year:=0; year<ms.numyr; year++ {
                for j := 0; j<len(ms.accounttable); j++ {
                    v = ms.accounttable[j].get("contributions", nil)
                    max = 0
                    if v != nil {
                        max = v[year]
                    } 
                    if ms.accounttable[j]["acctype"] != "aftertax"{
                        row = [0] * nvars
                        row[ms.vindx.D(year,j)] = 1
                        A = append(A, row)
                        b = append(b, max)
                    }
                }
            }
        }
        //
        // Add constaints for (10') rows
        //
        for year = 0; year<ms.numyr; year++ {
            for j := 0; j<min(2, len(ms.accounttable)); j++ {
                // at most the first two accounts are type IRA 
                // w/ RMD requirement
                if ms.accounttable[j]["acctype"] == "IRA"{
                    rmd = S.rmd_needed(year,ms.accounttable[j]["mykey"])
                    if rmd > 0{
                        row = [0] * nvars
                        row[ms.vindx.B(year,j)] = 1/rmd 
                        row[ms.vindx.W(year,j)] = -1
                        A = append(A, row)
                        b = append(b, 0)
                    }
                }
            } 
        }
    
        //
        // Add constraints for (11')
        //
        for year = 0; year<ms.numyr; year++ {
            adj_inf = ms.i_rate**(ms.preplanyears+year)
            row = [0] * nvars
            for j := 0; j<min(2, len(ms.accounttable)); j++ {
                // IRA can only be in the first two accounts
                if ms.accounttable[j]["acctype"] == "IRA"{
                    row[ms.vindx.W(year,j)] = 1 // Account 0 is TDRA
                    row[ms.vindx.D(year,j)] = -1 // Account 0 is TDRA
                }
            } 
            for k:=0;  k<len(taxtable);k++ {
                row[ms.vindx.X(year,k)] = -1
            }
            A = append(A, row)
            b = append(b, stded*adj_inf-ms.taxed[year]-SS_taxable*ms.SS[year])
        }
        //
        // Add constraints for (12')
        //
        for year = 0; year<ms.numyr; year++ {
            for k:=0;  k<len(taxtable)-1;k++ {
                row = [0] * nvars
                row[ms.vindx.X(year,k)] = 1
                A = append(A, row)
                b = append(b, (taxtable[k][1])*(ms.i_rate**(ms.preplanyears+year))) // inflation adjusted
            }
        }
        //
        // Add constraints for (13a')
        //
        if ms.accmap["aftertax"] > 0{
            for year:=0; year<ms.numyr; year++ {
                f = ms.cg_taxable_fraction(year)
                row = [0] * nvars
                for l:=0; l<len(capgainstable);l++{
                    row[ms.vindx.Y(year,l)] = 1
                }
                // Awful Hack! If year of asset sale, assume w(i,j)-D(i,j) is 
                // negative so taxable from this is zero
                if ms.cg_asset_taxed[year] <= 0{ // i.e., no sale
                    j = len(ms.accounttable)-1 // last Acc is investment / stocks
                    row[ms.vindx.W(year,j)] = -1*f 
                    row[ms.vindx.D(year,j)] = f 
                }
                A = append(A, row)
                b = append(b, ms.cg_asset_taxed[year])
            }
        }
        //
        // Add constraints for (13b')
        //
        if ms.accmap["aftertax"] > 0 {
            for year:=0; year<ms.numyr; year++ {
                f = ms.cg_taxable_fraction(year)
                row = [0] * nvars
                ////// Awful Hack! If year of asset sale, assume w(i,j)-D(i,j) is 
                ////// negative so taxable from this is zero
                if ms.cg_asset_taxed[year] <= 0{ // i.e., no sale
                    j = len(ms.accounttable)-1 // last Acc is investment / stocks
                    row[ms.vindx.W(year,j)] = f 
                    row[ms.vindx.D(year,j)] = -f 
                }
                for l:=0; l<len(capgainstable);l++{
                    row[ms.vindx.Y(year,l)] = -1
                }
                A = append(A, row)
                b = append(b, -ms.cg_asset_taxed[year])
            }
        }
        //
        // Add constraints for (14')
        //
        if ms.accmap["aftertax"] > 0{
            for year:=0; year<ms.numyr; year++ {
                adj_inf = ms.i_rate**(ms.preplanyears+year)
                for l:=0; l<len(capgainstable)-1;l++{
                    row = [0] * nvars
                    row[ms.vindx.Y(year,l)] = 1
                    for k:=0;  k<len(taxtable)-1;k++ {
                        if taxtable[k][0] >= capgainstable[l][0] && taxtable[k][0] < capgainstable[l+1][0] {
                            row[ms.vindx.X(year,k)] = 1
                        }
                    }
                    A = append(A, row)
                    b = append(b, capgainstable[l][1]*adj_inf) // mcg[i,l] inflation adjusted
                    //print_constraint( row, capgainstable[l][1]*adj_inf)
                }
            }
        }
        //
        // Add constraints for (15a')
        //
        for year = 0; year<ms.numyr; year++ {
            for j := 0; j<len(ms.accounttable); j++ {
                //j = len(ms.accounttable)-1 // nl the last account, the investment account
                row = [0] * nvars
                row[ms.vindx.B(year+1,j)] = 1 // b[i,j] supports an extra year
                row[ms.vindx.B(year,j)] = -1*ms.accounttable[j]["rate"]
                row[ms.vindx.W(year,j)] = ms.accounttable[j]["rate"]
                row[ms.vindx.D(year,j)] = -1*ms.accounttable[j]["rate"]
                A = append(A, row)
                // In the event of a sell of an asset for the year 
                temp = 0
                if ms.accounttable[j]["acctype"] == "aftertax"{
                    temp= ms.asset_sale[year] * ms.accounttable[j]["rate"] //TODO test
                }
                b = append(b, temp)
                //print("temp_a: ", temp, "rate", ms.accounttable[j]["rate"] , "asset sell price: ", ms.asset_sale[year]  )
            }
        } 
        //
        // Add constraints for (15b')
        //
        for year = 0; year<ms.numyr; year++ {
            for j := 0; j<len(ms.accounttable); j++ {
                //j = len(ms.accounttable)-1 // nl the last account, the investment account
                row = [0] * nvars
                row[ms.vindx.B(year,j)] = ms.accounttable[j]["rate"]
                row[ms.vindx.W(year,j)] = -1*ms.accounttable[j]["rate"]
                row[ms.vindx.D(year,j)] = ms.accounttable[j]["rate"]
                row[ms.vindx.B(year+1,j)] = -1  ////// b[i,j] supports an extra year
                A = append(A, row)
                temp = 0
                if ms.accounttable[j]["acctype"] == "aftertax"{
                    temp= -1 * ms.asset_sale[year] * ms.accounttable[j]["rate"]//TODO test
                }
                b = append(b, temp)
                //print("temp_b: ", temp, "rate", ms.accounttable[j]["rate"] , "asset sell price: ", ms.asset_sale[year]  )
            }
        }
        //
        // Constraint for (16a')
        //   Set the begining b[1,j] balances
        //
        for j := 0; j<len(ms.accounttable); j++ {
            row = [0] * nvars
            row[ms.vindx.B(0,j)] = 1
            A = append(A, row)
            b = append(b, ms.accounttable[j]["bal"])
        }
        //
        // Constraint for (16b')
        //   Set the begining b[1,j] balances
        //
        for j := 0; j<len(ms.accounttable); j++ {
            row = [0] * nvars
            row[ms.vindx.B(0,j)] = -1
            A = append(A, row)
            b = append(b, -1*ms.accounttable[j]["bal"])
        }
        //
        // Constrant for (17') is default for sycpy so no code is needed
        //
        if ms.verbose{
            print("Num vars: ", len(c))
            print("Num contraints: ", len(b))
            print()
        }
    
        return c, A, b
}
    
    
func (ms modelSpecs) cg_taxable_fraction(year int) float64 {
    f = 1
    if ms.accmap["aftertax"] > 0{
        for _, v := range ms.accounttable {
            if v["acctype"] == "aftertax" {
                if v["bal"] > 0 {
                    f = 1 - (v["basis"]/(v["bal"]*v["rate"]**year))
                }
                break // should be the last entry anyway but...
            }
        }
    }
    return f
}
    
    
func (ms modelSpecs) print_model_matrix(c []float64, A [][]float64, b []float64, s, non_binding_only bool){
        if ! non_binding_only {
            fmt.Printf("c: ")
            ms.print_model_row(c)
            fmt.Printf("\n")
            fmt.Printf("B? i: A_ub[i]: b[i]")
            for constraint:=0; constraint<len(A); constraint++ {
                if s == nil || s[constraint] > 0 {
                    fmt.Printf("  ")
                } else{
                    fmt.Printf("B ")
                }
                fmt.Printf("%d: ", constraint)
                ms.print_constraint( A[constraint], b[constraint])
            }
        } else {
            fmt.Printf(" i: A_ub[i]: b[i]")
            j = 0
            for constraint:=0;  contraint<len(A); constraint++ {
                if s[constraint] >0 {
                    j+=1
                    fmt.Printf("%d: ", constraint)
                    ms.print_constraint( A[constraint], b[constraint])
                }
            }
            fmt.Printf("\n\n%d non-binding constrains printed\n", j)
        }
        fmt.Printf("\n")
}
    
    
func (ms modelSpecs) print_constraint(row []float64, b float64) {
        ms.print_model_row(row, true)
        fmt.Printf("<= b[]: %6.2f", b)
}
    
func (ms modelSpecs) print_model_row(row []float64, suppress_newline bool) {
        S = ms.S
        taxtable = ms.taxtable
        capgainstable = ms.cgtaxtable
        penalty = ms.penalty
        stded = ms.stded
        SS_taxable = ms.ss_taxable
    
        for i:=0; i<ms.numyr; i++ {
            for k:=0;  k<len(taxtable);k++ {
                if row[ms.vindx.X(i, k)] != 0{
                    fmt.Printf("x[%d,%d]: %6.3f", i, k, row[ms.vindx.X(i, k)])
                }
            }
        }
        if ms.accmap["aftertax"] > 0 {
            for i:=0; i<ms.numyr; i++ {
                for l:=0; l<len(capgainstable);l++{
                    if row[ms.vindx.Y(i, l)] != 0{
                        fmt.Printf("y[%d,%d]: %6.3f ", i, l, row[ms.vindx.Y(i, l)])
                    }
                }
            }
        }
        for i:=0; i<ms.numyr; i++ {
            for j := 0; j<len(ms.accounttable); j++ {
                if row[ms.vindx.W(i, j)] != 0{
                    fmt.Printf("w[%d,%d]: %6.3f ", i, j, row[ms.vindx.W(i, j)])
                }
            }
        }
        for i:=0; i<ms.numyr+1; i++ { // b[] has an extra year
            for j := 0; j<len(ms.accounttable); j++ {
                if row[ms.vindx.B(i, j)] != 0{
                    fmt.Printf("b[%d,%d]: %6.3f ", i, j, row[ms.vindx.B(i, j)])
                }
            }
        }
        for i:=0; i<ms.numyr; i++ {
            if row[ms.vindx.S(i)] !=0{
                fmt.Printf("s[%d]: %6.3f ", i, row[ms.vindx.S(i)])
            }
        }
        if ms.accmap["aftertax"] > 0 {
            for i:=0; i<ms.numyr; i++ {
                for j := 0; j<len(ms.accounttable); j++ {
                    if row[ms.vindx.D(i,j)] !=0{
                        fmt.Printf("D[%d,%d]: %6.3f ", i, j, row[ms.vindx.D(i,j)])
                    }
                }
            }
        }
        if ! suppress_newline {
            fmt.Printf("\n")
        }
}
    