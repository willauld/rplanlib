package rplanlib

import "fmt"

func my_check_index_sequence(years, taxbins, cgbins, accounts int, accmap map[string]int, varindex vectorVarIndex) bool {
	// varindex.?() functions are laid out to index a vector of variables
	// laid out in the order x(i,k), y(i,l), w(i,j), b(i,j), s(i), D(i,j), ns()
	pass_ok := true
	ky := 0
	//row = [0] * nvars
	for i := 0; i < years; i++ {
		for k := 0; k < taxbins; k++ {
			if varindex.X(i, k) != ky {
				pass_ok = false
				fmt.Printf("varindex.x(%d,%d) is %d not %d as it should be",
					i, k, varindex.X(i, k), ky)
			}
			ky += 1
		}
	}
	if accmap["aftertax"] > 0 {
		for i := 0; i < years; i++ {
			for l := 0; l < cgbins; l++ {
				if varindex.Y(i, l) != ky {
					pass_ok = false
					fmt.Printf("varindex.y(%d,%d) is %d not %d as it should be",
						i, l, varindex.Y(i, l), ky)
				}
				ky += 1
			}
		}
	}
	for i := 0; i < years; i++ {
		for j := 0; j < accounts; j++ {
			if varindex.W(i, j) != ky {
				pass_ok = false
				fmt.Printf("varindex.w(%d,%d) is %d not %d as it should be",
					i, j, varindex.W(i, j), ky)
			}
			ky += 1
		}
	}
	for i := 0; i < years+1; i++ { // b[] has an extra year
		for j := 0; j < accounts; j++ {
			if varindex.B(i, j) != ky {
				pass_ok = false
				fmt.Printf("varindex.b(%d,%d) is %d not %d as it should be",
					i, j, varindex.B(i, j), ky)
			}
			ky += 1
		}
	}
	for i := 0; i < years; i++ {
		if varindex.S(i) != ky {
			pass_ok = false
			fmt.Printf("varindex.s(%d) is %d not %d as it should be",
				i, varindex.S(i), ky)
		}
		ky += 1
	}
	if accmap["aftertax"] > 0 {
		for i := 0; i < years; i++ {
			for j := 0; j < accounts; j++ {
				if varindex.D(i, j) != ky {
					pass_ok = false
					fmt.Printf("varindex.D(%d,%d) is %d not %d as it should be",
						i, j, varindex.D(i, j), ky)
				}
				ky += 1
			}
		}
	}
	return pass_ok
}

type vectorVarIndex struct {
	// inplements the vector var index functions
	Years    int
	Taxbins  int
	Cgbins   int
	Accounts int
	Accmap   map[string]int
	Xcount   int
	Ycount   int
	Wcount   int
	Bcount   int
	Scount   int
	Dcount   int
	Vsize    int
	Ystart   int
	Wstart   int
	Bstart   int
	Sstart   int
	Dstart   int
}

func NewVectorVarIndex(iyears, itaxbins, icgbins, iaccounts int,
	iaccmap map[string]int) vectorVarIndex {

	ycount := 0
	if iaccmap["aftertax"] > 0 { // no cgbins if no aftertax account
		ycount = iyears * icgbins
	}
	xcount := iyears * itaxbins
	wcount := iyears * iaccounts
	bcount := (iyears + 1) * iaccounts
	dcount := iyears * iaccounts
	scount := iyears
	vsize := xcount + ycount + wcount + bcount + scount + dcount
	ystart := xcount
	wstart := ystart + ycount
	bstart := wstart + wcount
	sstart := bstart + bcount
	dstart := sstart + scount

	return vectorVarIndex{
		Years:    iyears,
		Taxbins:  itaxbins,
		Cgbins:   icgbins,
		Accounts: iaccounts,
		Accmap:   iaccmap,
		Xcount:   xcount,
		Ycount:   ycount,
		Wcount:   wcount,
		// final balances in years+1,
		Bcount: bcount,
		Scount: scount,
		Dcount: dcount,
		Vsize:  vsize,

		//xstart = 0
		Ystart: ystart,
		Wstart: wstart,
		Bstart: bstart,
		Sstart: sstart,
		Dstart: dstart,
	}
}

func (v vectorVarIndex) X(i, k int) int {
	//assert i >= 0 and i < v.Years
	//assert k >= 0 and k < v.Taxbins
	return i*v.Taxbins + k
}

func (v vectorVarIndex) Y(i, l int) int {
	//assert v.Accmap["aftertax"] > 0
	//assert i >= 0 and i < v.Years
	//assert l >= 0 and l < v.Cgbins
	return v.Ystart + i*v.Cgbins + l
}

func (v vectorVarIndex) W(i, j int) int {
	//assert i >= 0 and i < v.Years
	//assert j >= 0 and j < v.Accounts
	return v.Wstart + i*v.Accounts + j
}

func (v vectorVarIndex) B(i, j int) int {
	//assert i >= 0 and i < v.Years + 1  // b has an extra year on the end
	//assert j >= 0 and j < v.Accounts
	return v.Bstart + i*v.Accounts + j
}

func (v vectorVarIndex) S(i int) int {
	//assert i >= 0 and i < v.Years
	return v.Sstart + i
}

func (v vectorVarIndex) D(i, j int) int {
	//assert S.accmap["aftertax"] > 0
	//assert j >= 0 and j < v.Accounts
	//assert i >= 0 and i < v.Years
	return v.Dstart + i*v.Accounts + j
}

func (v vectorVarIndex) Varstr(indx int) string {
	var a, b, c int

	//assert indx < v.Vsize
	if indx < v.Xcount {
		a = indx // v.taxbins
		b = indx % v.Taxbins
		return fmt.Sprintf("x[%d,%d]", a, b) // add actual values for i,j
	} else if indx < v.Xcount+v.Ycount {
		c = indx - v.Xcount
		a = c // v.Cgbins
		b = c % v.Cgbins
		return fmt.Sprintf("y[%d,%d]", a, b) // add actual values for i,j
	} else if indx < v.Xcount+v.Ycount+v.Wcount {
		c = indx - (v.Xcount + v.Ycount)
		a = c // v.Accounts
		b = c % v.Accounts
		return fmt.Sprintf("w[%d,%d]", a, b) // add actual values for i,j
	} else if indx < v.Xcount+v.Ycount+v.Wcount+v.Bcount {
		c = indx - (v.Xcount + v.Ycount + v.Wcount)
		a = c // v.Accounts
		b = c % v.Accounts
		return fmt.Sprintf("b[%d,%d]", a, b) // add actual values for i,j
	} else if indx < v.Xcount+v.Ycount+v.Wcount+v.Bcount+v.Scount {
		c = indx - (v.Xcount + v.Ycount + v.Wcount + v.Bcount)
		//a = c // v.Years
		//b = c % v.Years
		return fmt.Sprintf("s[%d]", c) // add actual values for i,j
	} else if indx < v.Xcount+v.Ycount+v.Wcount+v.Bcount+v.Scount+v.Dcount {
		c = indx - (v.Xcount + v.Ycount + v.Wcount + v.Bcount + v.Scount)
		a = c // v.Accounts
		b = c % v.Accounts
		return fmt.Sprintf("D[%d,%d]", a, b) // add actual values for i,j
	} else {
		fmt.Printf("\nError -- varstr() corupted\n")
	}
	return "don't know"
}

//v.ycount = 0
//if v.accmap["aftertax"] > 0:  /// no cgbins if no aftertax account
//    v.ycount = v.years * v.cgbins
//v.wcount = v.years * v.accounts
// final balances in years+1
//v.bcount = (v.years + 1) * v.accounts
//v.scount = v.years
