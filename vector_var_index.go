package rplanlib

import (
	"fmt"
	"os"
)

func checkIndexSequence(years, taxbins, cgbins int, accmap map[Acctype]int, varindex VectorVarIndex, errfile *os.File) bool {
	accounts := 0
	for _, acc := range accmap {
		accounts += acc
	}
	// varindex.?() functions are laid out to index a vector of variables
	// laid out in the order:
	//   x(i,k),
	//   y(i,l),
	//   w(i,j),
	//   b(i,j),
	//   s(i),
	//   D(i,j),
	//   it(i),
	//   ti(i),
	//   cg(i),
	//   cgt(i),
	//   sgts(i,l)
	//   boolvar(i,l) // 6 for each year
	// The idea for this test is to step through the entire vector,
	// cell by cell, based on the generated indexes from the access
	// functions. By walking through each access function one at a time
	// we should get ever index of the vector in order. That is 0..vector
	// length -1. If we find any access function not returning an index
	// matching the current vector index position, signal an error
	passOk := true
	ky := 0
	//row = [0] * nvars
	if !GetPiecewiseChoice() {
		for i := 0; i < years; i++ {
			for k := 0; k < taxbins; k++ {
				if varindex.X(i, k) != ky {
					passOk = false
					fmt.Fprintf(errfile, "varindex.x(%d,%d) is %d not %d as it should be\n",
						i, k, varindex.X(i, k), ky)
				}
				ky++
			}
		}
		if accmap[Aftertax] > 0 {
			for i := 0; i < years; i++ {
				for l := 0; l < cgbins; l++ {
					if varindex.Sy(i, l) != ky {
						passOk = false
						fmt.Fprintf(errfile,
							"varindex.y(%d,%d) is %d not %d as it should be\n",
							i, l, varindex.Sy(i, l), ky)
					}
					ky++
				}
			}
			for i := 0; i < years; i++ {
				for l := 0; l < cgbins; l++ {
					if varindex.Y(i, l) != ky {
						passOk = false
						fmt.Fprintf(errfile, "varindex.y(%d,%d) is %d not %d as it should be\n",
							i, l, varindex.Y(i, l), ky)
					}
					ky++
				}
			}
		}
	}
	for i := 0; i < years; i++ {
		for j := 0; j < accounts; j++ {
			if varindex.W(i, j) != ky {
				passOk = false
				fmt.Fprintf(errfile, "varindex.w(%d,%d) is %d not %d as it should be\n",
					i, j, varindex.W(i, j), ky)
			}
			ky++
		}
	}
	for i := 0; i < years+1; i++ { // b[] has an extra year
		for j := 0; j < accounts; j++ {
			if varindex.B(i, j) != ky {
				passOk = false
				fmt.Fprintf(errfile, "varindex.b(%d,%d) is %d not %d as it should be\n",
					i, j, varindex.B(i, j), ky)
			}
			ky++
		}
	}
	for i := 0; i < years; i++ {
		if varindex.S(i) != ky {
			passOk = false
			fmt.Fprintf(errfile, "varindex.s(%d) is %d not %d as it should be\n",
				i, varindex.S(i), ky)
		}
		ky++
	}
	if accmap[Aftertax] > 0 {
		for i := 0; i < years; i++ {
			if varindex.D(i) != ky {
				passOk = false
				fmt.Fprintf(errfile, "varindex.D(%d) is %d not %d as it should be\n",
					i, varindex.D(i), ky)
			}
			ky++
		}
	}
	if GetPiecewiseChoice() {
		for i := 0; i < years; i++ {
			if varindex.IT(i) != ky {
				passOk = false
				fmt.Fprintf(errfile, "varindex.IT(%d) is %d not %d as it should be\n",
					i, varindex.IT(i), ky)
			}
			ky++
		}
		for i := 0; i < years; i++ {
			if varindex.TI(i) != ky {
				passOk = false
				fmt.Fprintf(errfile, "varindex.TI(%d) is %d not %d as it should be\n",
					i, varindex.TI(i), ky)
			}
			ky++
		}
		if accmap[Aftertax] > 0 {
			for i := 0; i < years; i++ {
				if varindex.CG(i) != ky {
					passOk = false
					fmt.Fprintf(errfile, "varindex.CG(%d) is %d not %d as it should be\n",
						i, varindex.CG(i), ky)
				}
				ky++
			}
			for i := 0; i < years; i++ {
				if varindex.CGT(i) != ky {
					passOk = false
					fmt.Fprintf(errfile, "varindex.CGT(%d) is %d not %d as it should be\n",
						i, varindex.CGT(i), ky)
				}
				ky++
			}
			for i := 0; i < years; i++ { // cgts[]
				for l := 0; l < cgbins; l++ {
					if varindex.CGTS(i, l) != ky {
						passOk = false
						fmt.Fprintf(errfile, "varindex.cgts(%d,%d) is %d not %d as it should be\n",
							i, l, varindex.CGTS(i, l), ky)
					}
					ky++
				}
			}
			for i := 0; i < years; i++ { // Boolvar[]
				for b := 0; b < varindex.Boolvars; b++ {
					if varindex.Boolvar(i, b) != ky {
						passOk = false
						fmt.Fprintf(errfile, "varindex.Boolvar(%d,%d) is %d not %d as it should be\n",
							i, b, varindex.Boolvar(i, b), ky)
					}
					ky++
				}
			}
		}
		if ky != varindex.Vsize {
			fmt.Fprintf(errfile, "varindex testing only %d cells but vector length is %d\n",
				ky, varindex.Vsize)

		}
	}
	return passOk
}

// VectorVarIndex contains the index information to convert from variable index to vector index
type VectorVarIndex struct {
	// inplements the vector var index functions
	Years        int
	Taxbins      int
	Cgbins       int
	Accounts     int
	Accmap       map[Acctype]int
	Accname      []string
	Boolvars     int
	Xcount       int
	Sycount      int
	Ycount       int
	Wcount       int
	Bcount       int
	Scount       int
	Dcount       int
	Itcount      int
	Ticount      int
	Cgcount      int
	Cgtcount     int
	Cgtscount    int
	Boolvarcount int
	Vsize        int
	Systart      int
	Ystart       int
	Wstart       int
	Bstart       int
	Sstart       int
	Dstart       int
	Itstart      int
	Tistart      int
	Cgstart      int
	Cgtstart     int
	Cgtsstart    int
	Boolvarstart int
	errfile      *os.File
}

// NewVectorVarIndex creates an object for index translation
func NewVectorVarIndex(iyears, itaxbins, icgbins int,
	iaccmap map[Acctype]int, errfile *os.File) (VectorVarIndex, error) {

	boolvars := 6

	if iyears < 1 || iyears > 100 {
		e := fmt.Errorf("NewVectorVarIndex: invalid value for year, %d", iyears)
		return VectorVarIndex{}, e
	}
	if itaxbins < 0 {
		e := fmt.Errorf("NewVectorVarIntex: invalid value, taxbins < 0")
		return VectorVarIndex{}, e
	}
	if icgbins < 0 {
		e := fmt.Errorf("NewVectorVarIndex: invalid value, cgbins < 0")
		return VectorVarIndex{}, e
	}
	if len(iaccmap) != 3 {
		e := fmt.Errorf("NewVectorVarIndex: invalid value, accmap length != 3 but rather %d, accmap: %v", len(iaccmap), iaccmap)
		return VectorVarIndex{}, e
	}
	_, okIRA := iaccmap[IRA]
	_, okRoth := iaccmap[Roth]
	_, okAftertax := iaccmap[Aftertax]
	if !okIRA || !okRoth || !okAftertax {
		e := fmt.Errorf("NewVectorVarIndex: accmap missing one of (IRA, roth, aftertax) key values, has: %v", iaccmap)
		return VectorVarIndex{}, e
	}

	iaccounts := 0
	accname := make([]string, 5)
	for j := Acctype(0); j < Acctype(len(iaccmap)); j++ {
		n := iaccmap[j]
		for i := 0; i < n; i++ {
			accname[iaccounts+i] = fmt.Sprintf("%s%d", j.String(), i+1)
		}
		iaccounts += n
	}
	//fmt.Printf("iaccounts: %d, iaccmap: %v\n", iaccounts, iaccmap)
	xcount := 0
	ycount := 0
	sycount := 0
	itcount := 0
	ticount := 0
	cgcount := 0
	cgtcount := 0
	cgtscount := 0
	boolvarcount := 0
	if iaccmap[Aftertax] > 0 { // no cgbins if no aftertax account
		if GetPiecewiseChoice() {
			cgcount = iyears
			cgtcount = iyears
			cgtscount = iyears * icgbins
			boolvarcount = iyears * boolvars

		} else {
			ycount = iyears * icgbins
			sycount = iyears * icgbins
		}
	}
	if !GetPiecewiseChoice() {
		xcount = iyears * itaxbins
	}
	wcount := iyears * iaccounts
	bcount := (iyears + 1) * iaccounts
	scount := iyears
	dcount := iyears
	if GetPiecewiseChoice() {
		itcount = iyears
		ticount = iyears
	}
	vsize := xcount + sycount + ycount + wcount + bcount + scount + dcount + itcount + ticount + cgcount + cgtcount + cgtscount + boolvarcount
	systart := xcount
	ystart := systart + sycount
	wstart := ystart + ycount
	bstart := wstart + wcount
	sstart := bstart + bcount
	dstart := sstart + scount
	itstart := dstart + dcount
	tistart := itstart + itcount
	cgstart := tistart + ticount
	cgtstart := cgstart + cgcount
	cgtsstart := cgtstart + cgtcount
	boolvarstart := cgtsstart + cgtscount

	return VectorVarIndex{
		Years:    iyears,
		Taxbins:  itaxbins,
		Cgbins:   icgbins,
		Accounts: iaccounts,
		Accmap:   iaccmap,
		Accname:  accname,
		Boolvars: boolvars,
		Xcount:   xcount,
		Sycount:  sycount,
		Ycount:   ycount,
		Wcount:   wcount,
		// final balances in years+1,
		Bcount:       bcount,
		Scount:       scount,
		Dcount:       dcount,
		Itcount:      itcount,
		Ticount:      ticount,
		Cgcount:      cgcount,
		Cgtcount:     cgtcount,
		Cgtscount:    cgtscount,
		Boolvarcount: boolvarcount,
		Vsize:        vsize,

		//xstart = 0
		Systart:      systart,
		Ystart:       ystart,
		Wstart:       wstart,
		Bstart:       bstart,
		Sstart:       sstart,
		Dstart:       dstart,
		Itstart:      itstart,
		Tistart:      tistart,
		Cgstart:      cgstart,
		Cgtstart:     cgtstart,
		Cgtsstart:    cgtsstart,
		Boolvarstart: boolvarstart,

		errfile: errfile,
	}, nil
}

// X (i,k) returns the variable index in a variable vector
func (v VectorVarIndex) X(i, k int) int {
	//assert i >= 0 and i < v.Years
	//assert k >= 0 and k < v.Taxbins
	if v.Xcount == 0 {
		panic("Xcount is zero so you can't index to it")
	}
	return i*v.Taxbins + k
}

// Sy (i,l) returns the variable index in a variable vector
func (v VectorVarIndex) Sy(i, l int) int {
	//assert v.Accmap["aftertax"] > 0
	//assert i >= 0 and i < v.Years
	//assert l >= 0 and l < v.Cgbins
	if v.Sycount == 0 {
		panic("Sycount is zero so you can't index to it")
	}
	return v.Systart + i*v.Cgbins + l
}

// Y (i,l) returns the variable index in a variable vector
func (v VectorVarIndex) Y(i, l int) int {
	//assert v.Accmap["aftertax"] > 0
	//assert i >= 0 and i < v.Years
	//assert l >= 0 and l < v.Cgbins
	if v.Ycount == 0 {
		panic("Ycount is zero so you can't index to it")
	}
	return v.Ystart + i*v.Cgbins + l
}

// W (i,j) returns the variable index in a variable vector
func (v VectorVarIndex) W(i, j int) int {
	//assert i >= 0 and i < v.Years
	//assert j >= 0 and j < v.Accounts
	return v.Wstart + i*v.Accounts + j
}

// B (i,j) returns the variable index in a variable vector
func (v VectorVarIndex) B(i, j int) int {
	//assert i >= 0 and i < v.Years + 1  // b has an extra year on the end
	//assert j >= 0 and j < v.Accounts
	return v.Bstart + i*v.Accounts + j
}

// S (i) returns the variable index in a variable vector
func (v VectorVarIndex) S(i int) int {
	//assert i >= 0 and i < v.Years
	return v.Sstart + i
}

// D (i) returns the variable index in a variable vector
func (v VectorVarIndex) D(i int) int {
	//assert S.accmap["aftertax"] > 0
	//assert i >= 0 and i < v.Years
	return v.Dstart + i
}

// IT (i) returns the variable index in a variable vector
func (v VectorVarIndex) IT(i int) int {
	//assert i >= 0 and i < v.Years
	if v.Itcount == 0 {
		panic("Itcount is zero so you can't index to it")
	}
	return v.Itstart + i
}

// TI (i) returns the variable index in a variable vector
func (v VectorVarIndex) TI(i int) int {
	//assert i >= 0 and i < v.Years
	if v.Ticount == 0 {
		panic("Ticount is zero so you can't index to it")
	}
	return v.Tistart + i
}

// CG (i) returns the variable index in a variable vector
func (v VectorVarIndex) CG(i int) int {
	//assert S.accmap["aftertax"] > 0
	//assert i >= 0 and i < v.Years
	if v.Cgcount == 0 {
		panic("Cgcount is zero so you can't index to it")
	}
	return v.Cgstart + i
}

// CGT (i) returns the variable index in a variable vector
func (v VectorVarIndex) CGT(i int) int {
	//assert S.accmap["aftertax"] > 0
	//assert i >= 0 and i < v.Years
	if v.Cgtcount == 0 {
		panic("Cgtcount is zero so you can't index to it")
	}
	return v.Cgtstart + i
}

// CGT (i,l) returns the variable index in a variable vector
func (v VectorVarIndex) CGTS(i, l int) int {
	//assert v.Accmap["aftertax"] > 0
	//assert i >= 0 and i < v.Years
	//assert l >= 0 and l < v.Cgbins
	if v.Cgtscount == 0 {
		panic("Cgtscount is zero so you can't index to it")
	}
	return v.Cgtsstart + i*v.Cgbins + l
}

// Boolvar (i,b) returns the variable index in a variable vector
func (v VectorVarIndex) Boolvar(i, b int) int {
	Assert(v.Accmap[Aftertax] > 0, "Vector Var Boolvar is not present unless there is an aftertax account")
	Assert(i >= 0 && i < v.Years, fmt.Sprintf("For Boolvar(%d,%d), year must be at least zero", i, b))
	Assert(b >= 0 && b < v.Boolvars, fmt.Sprintf("For Boolvar(%d,%d), second index must be at least zero", i, b))
	a := v.Boolvarstart + i*v.Boolvars + b
	Assert(a < v.Vsize,
		fmt.Sprintf("For Boolvar(%d,%d), return index (%d) must be less than Vsize (%d)", i, b, a, v.Vsize))
	if v.Boolvarcount == 0 {
		panic("Boolvarcount is zero so you can't index to it")
	}
	return a
}

// Varstr returns the variable name and index(s) for the variable at indx in the variable vector
func (v VectorVarIndex) Varstr(indx int) string {
	var a, b, c int
	var name string

	//assert indx < v.Vsize
	if indx < v.Xcount {
		a = indx / v.Taxbins
		b = indx % v.Taxbins
		return fmt.Sprintf("x[%d,%d]", a, b)
	} else if indx < v.Xcount+v.Sycount {
		c = indx - v.Xcount
		a = c / v.Cgbins
		b = c % v.Cgbins
		return fmt.Sprintf("Sy[%d,%d]", a, b)
	} else if indx < v.Xcount+v.Sycount+v.Ycount {
		c = indx - v.Xcount + v.Sycount
		a = c / v.Cgbins
		b = c % v.Cgbins
		return fmt.Sprintf("y[%d,%d]", a, b)
	} else if indx < v.Xcount+v.Sycount+v.Ycount+v.Wcount {
		c = indx - (v.Xcount + v.Sycount + v.Ycount)
		a = c / v.Accounts
		b = c % v.Accounts
		name = v.Accname[b]
		return fmt.Sprintf("w[%d,%d=%s]", a, b, name)
	} else if indx < v.Xcount+v.Sycount+v.Ycount+v.Wcount+v.Bcount {
		c = indx - (v.Xcount + v.Sycount + v.Ycount + v.Wcount)
		a = c / v.Accounts
		b = c % v.Accounts
		name = v.Accname[b]
		return fmt.Sprintf("b[%d,%d=%s]", a, b, name)
	} else if indx < v.Xcount+v.Sycount+v.Ycount+v.Wcount+v.Bcount+v.Scount {
		c = indx - (v.Xcount + v.Sycount + v.Ycount + v.Wcount + v.Bcount)
		//a = c / v.Years
		//b = c % v.Years
		return fmt.Sprintf("s[%d]", c) // add actual values for i,j
	} else if indx < v.Xcount+v.Sycount+v.Ycount+v.Wcount+v.Bcount+v.Scount+v.Dcount {
		c = indx - (v.Xcount + v.Sycount + v.Ycount + v.Wcount + v.Bcount + v.Scount)
		//a = c / v.Accounts
		//b = c % v.Accounts
		//name = v.Accname[b]
		return fmt.Sprintf("d[%d]", c)
	} else if indx < v.Xcount+v.Sycount+v.Ycount+v.Wcount+v.Bcount+v.Scount+v.Dcount+v.Itcount {
		c = indx - (v.Xcount + v.Sycount + v.Ycount + v.Wcount + v.Bcount + v.Scount + v.Dcount)
		return fmt.Sprintf("it[%d]", c)
	} else if indx < v.Xcount+v.Sycount+v.Ycount+v.Wcount+v.Bcount+v.Scount+v.Dcount+v.Itcount+v.Ticount {
		c = indx - (v.Xcount + v.Sycount + v.Ycount + v.Wcount + v.Bcount + v.Scount + v.Dcount + v.Itcount)
		return fmt.Sprintf("ti[%d]", c)
	} else if indx < v.Xcount+v.Sycount+v.Ycount+v.Wcount+v.Bcount+v.Scount+v.Dcount+v.Itcount+v.Ticount+v.Cgcount {
		c = indx - (v.Xcount + v.Sycount + v.Ycount + v.Wcount + v.Bcount + v.Scount + v.Dcount + v.Itcount + v.Ticount)
		return fmt.Sprintf("cg[%d]", c)
	} else if indx < v.Xcount+v.Sycount+v.Ycount+v.Wcount+v.Bcount+v.Scount+v.Dcount+v.Itcount+v.Ticount+v.Cgcount+v.Cgtcount {
		c = indx - (v.Xcount + v.Sycount + v.Ycount + v.Wcount + v.Bcount + v.Scount + v.Dcount + v.Itcount + v.Ticount + v.Cgcount)
		return fmt.Sprintf("cgt[%d]", c)
	} else if indx < v.Xcount+v.Sycount+v.Ycount+v.Wcount+v.Bcount+v.Scount+v.Dcount+v.Itcount+v.Ticount+v.Cgcount+v.Cgtcount+v.Cgtscount {
		c = indx - (v.Xcount + v.Sycount + v.Ycount + v.Wcount + v.Bcount + v.Scount + v.Dcount + v.Itcount + v.Ticount + v.Cgcount + v.Cgtcount)
		a = c / v.Cgcount
		b = c % v.Cgcount
		return fmt.Sprintf("cgts[%d,%d]", a, b)
	} else if indx < v.Xcount+v.Sycount+v.Ycount+v.Wcount+v.Bcount+v.Scount+v.Dcount+v.Itcount+v.Ticount+v.Cgcount+v.Cgtcount+v.Cgtscount+v.Boolvarcount {
		c = indx - (v.Xcount + v.Sycount + v.Ycount + v.Wcount + v.Bcount + v.Scount + v.Dcount + v.Itcount + v.Ticount + v.Cgcount + v.Cgtcount + v.Boolvarcount)
		a = c / v.Boolvarcount
		b = c % v.Boolvarcount
		return fmt.Sprintf("boolvar[%d,%d]", a, b)
	}
	if v.errfile != nil {
		fmt.Fprintf(v.errfile, "\nError -- varstr() corupted\n")
	}
	return "don't know"
}
