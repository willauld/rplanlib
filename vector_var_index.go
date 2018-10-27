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
	// laid out in the order x(i,k), y(i,l), w(i,j), b(i,j), s(i), D(i,j)
	passOk := true
	ky := 0
	//row = [0] * nvars
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
	return passOk
}

// VectorVarIndex contains the index information to convert from variable index to vector index
type VectorVarIndex struct {
	// inplements the vector var index functions
	Years    int
	Taxbins  int
	Cgbins   int
	Accounts int
	Accmap   map[Acctype]int
	Accname  []string
	Xcount   int
	Sycount  int
	Ycount   int
	Wcount   int
	Bcount   int
	Scount   int
	Dcount   int
	Vsize    int
	Systart  int
	Ystart   int
	Wstart   int
	Bstart   int
	Sstart   int
	Dstart   int
	errfile  *os.File
}

// NewVectorVarIndex creates an object for index translation
func NewVectorVarIndex(iyears, itaxbins, icgbins int,
	iaccmap map[Acctype]int, errfile *os.File) (VectorVarIndex, error) {

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
	ycount := 0
	sycount := 0
	if iaccmap[Aftertax] > 0 { // no cgbins if no aftertax account
		ycount = iyears * icgbins
		sycount = iyears * icgbins
	}
	xcount := iyears * itaxbins
	wcount := iyears * iaccounts
	bcount := (iyears + 1) * iaccounts
	scount := iyears
	dcount := iyears
	vsize := xcount + 2*ycount + wcount + bcount + scount + dcount
	systart := xcount
	ystart := systart + sycount
	wstart := ystart + ycount
	bstart := wstart + wcount
	sstart := bstart + bcount
	dstart := sstart + scount

	return VectorVarIndex{
		Years:    iyears,
		Taxbins:  itaxbins,
		Cgbins:   icgbins,
		Accounts: iaccounts,
		Accmap:   iaccmap,
		Accname:  accname,
		Xcount:   xcount,
		Sycount:  sycount,
		Ycount:   ycount,
		Wcount:   wcount,
		// final balances in years+1,
		Bcount: bcount,
		Scount: scount,
		Dcount: dcount,
		Vsize:  vsize,

		//xstart = 0
		Systart: systart,
		Ystart:  ystart,
		Wstart:  wstart,
		Bstart:  bstart,
		Sstart:  sstart,
		Dstart:  dstart,

		errfile: errfile,
	}, nil
}

// X (i,k) returns the variable index in a variable vector
func (v VectorVarIndex) X(i, k int) int {
	//assert i >= 0 and i < v.Years
	//assert k >= 0 and k < v.Taxbins
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
	}
	if v.errfile != nil {
		fmt.Fprintf(v.errfile, "\nError -- varstr() corupted\n")
	}
	return "don't know"
}
