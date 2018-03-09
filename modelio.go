package rplanlib

import (
	"encoding/binary"
	"fmt"
	"log"
	"math"
	"os"
)

func TersPrintArray(a []float64) {
	l := len(a)
	if l > 8 {
		//fmt.Printf("TPA:\n")
		fmt.Printf("[ %f\t%f\t%f,,,\t%f\t%f\t%f\t]\n", a[0], a[1], a[2], a[l-3], a[l-2], a[l-1])
	} else {
		fmt.Printf("TPA: %v", a)
	}
}
func TersPrintMatrix(a [][]float64) error {
	r := len(a)
	c := len(a[0])
	fmt.Printf("Rows: %d, Cols: %d\n", len(a), len(a[0]))
	fmt.Printf("[\n")
	if r > 8 {
		for i := 0; i < 3; i++ {
			if len(a[i]) != c {
				return fmt.Errorf("TersPrintMatrix: inconsistant row lenth")
			}
			TersPrintArray(a[i])
		}
		for i := 0; i < 2; i++ {
			fmt.Printf("		.				.				.\n")
		}
		for i := 3; i > 0; i-- {
			if len(a[r-i]) != c {
				return fmt.Errorf("TersPrintMatrix: inconsistant row lenth")
			}
			TersPrintArray(a[r-i])
		}
	} else {
		for i := 0; i < r; i++ {
			if len(a[i]) != c {
				return fmt.Errorf("TersPrintMatrix: inconsistant row lenth")
			}
			TersPrintArray(a[i])
		}
	}
	fmt.Printf("]\n")
	return nil
}

// BinDumpModel writes a comparable binary version of the retirement plan model
func BinDumpModel(c []float64, A [][]float64, b []float64, x []float64,
	vid *[]int32, fname string) error {
	if fname == "" {
		fname = "./RPlanModelgo.dat"
	}
	filem, err := os.Create(fname)
	if err != nil {
		e := fmt.Errorf("os.Open: %s", err)
		return e
	}
	Endian := binary.LittleEndian
	//stream = open(fname, 'wb')

	// write C array
	//header := make([]uint32, 3)
	header := []uint32{uint32(len(c)), 0, 0xDEADBEEF}
	err = binary.Write(filem, Endian, &header)
	if err != nil {
		e := fmt.Errorf("BinDumpModel(1): %s", err)
		return e
	}
	err = binary.Write(filem, Endian, &c)
	if err != nil {
		e := fmt.Errorf("BinDumpModel(2): %s", err)
		return e
	}
	//fmt.Printf("c Header Rows: %d, Cols: %d, code: %#X\n", header[0], header[1], header[2])

	// write A matrix
	header = []uint32{uint32(len(A)), uint32(len(A[0])), 0xDEADBEEF}
	//fmt.Printf("A length: %d, %d, dumping\n", (len(A), len(A[0])))
	err = binary.Write(filem, Endian, &header)
	if err != nil {
		e := fmt.Errorf("BinDumpModel(3): %s", err)
		return e
	}
	for _, row := range A {
		err = binary.Write(filem, Endian, &row)
		if err != nil {
			e := fmt.Errorf("BinDumpModel(4): %s", err)
			return e
		}
	}

	// write b vector
	header = []uint32{uint32(len(b)), 0, 0xDEADBEEF}
	//fmt.Printf("b length: %d, dumping\n", len(b))
	err = binary.Write(filem, Endian, &header)
	if err != nil {
		e := fmt.Errorf("BinDumpModel(5): %s", err)
		return e
	}
	err = binary.Write(filem, Endian, &b)
	if err != nil {
		e := fmt.Errorf("BinDumpModel(6): %s", err)
		return e
	}

	// Write x vector
	xlen := 0
	xOverhead := 0
	if x != nil {
		header = []uint32{uint32(len(x)), 0, 0xDEADBEEF}
		//fmt.Printf("x length: %d, dumping\n", len(x))
		err = binary.Write(filem, Endian, &header)
		if err != nil {
			e := fmt.Errorf("BinDumpModel(7): %s", err)
			return e
		}
		err = binary.Write(filem, Endian, &x)
		if err != nil {
			e := fmt.Errorf("BinDumpModel(8): %s", err)
			return e
		}
		xOverhead = 12
		xlen = len(x)
	}

	// Write Vector index data
	vidlen := 0
	vidOverhead := 0
	if vid != nil {
		header = []uint32{uint32(len(*vid)), 0, 0xDEADBEEF}
		fmt.Printf("vid length: %d, dumping\n", len(*vid))
		err = binary.Write(filem, Endian, &header)
		if err != nil {
			e := fmt.Errorf("BinDumpModel(9): %s", err)
			return e
		}
		err = binary.Write(filem, Endian, vid)
		if err != nil {
			e := fmt.Errorf("BinDumpModel(10): %s", err)
			return e
		}
		vidOverhead = 12
		vidlen = len(*vid)
	}
	//filem.Close()

	stats, err := filem.Stat()
	if err != nil {
		e := fmt.Errorf("os.Stat: %s", err)
		return e
	}
	fsize := stats.Size()

	overhead := 3*12 + xOverhead + vidOverhead
	csize := 8 * len(c)
	Asize := 8 * len(A) * len(A[0])
	bsize := 8 * len(b)
	xsize := 8 * xlen
	vidsize := 4 * vidlen
	fmt.Printf("xOverhead: %d, vidOverhead: %d, xlen: %d, vidlen: %d\n", xOverhead, vidOverhead, xlen, vidlen)
	calcsize := int64(overhead + csize + Asize + bsize + xsize + vidsize)
	diff := fsize - calcsize
	if fsize != calcsize {
		fmt.Printf("BinDumpModel Error - dump file size error, filesize: %d,\n\tlen(c): %d,\n\tlen(A): %d,\n\tLen(A[0]): %d,\n\tlen(b): %d,\n\tlen(x): %d,\n\tvidlen: %d,\n\t\tdiff: %d\n", fsize, len(c), len(A), len(A[0]), len(b), xlen, vidlen, diff)
	}
	binDumpCheck(c, A, b, x, vid, fname)
	return nil
}

func binDumpCheck(c []float64, A [][]float64, b []float64, x []float64,
	vid *[]int32, ftocheck string) {
	xlen := 0
	xOverhead := 0
	if x != nil {
		xlen = len(x)
		xOverhead = 12
	}
	vidlen := 0
	vidOverhead := 0
	if vid != nil {
		vidlen = len(*vid)
		vidOverhead = 12
	}
	filex, err := os.Open(ftocheck) // For read access.
	if err != nil {
		log.Fatal(err)
	}
	stats, err := filex.Stat()
	if err != nil {
		fmt.Printf("os.Stat: %s\n", err)
		os.Exit(1)
	}
	fsize := stats.Size()
	filex.Close()
	calcsize := xOverhead + vidOverhead + 3*12 + 8*len(c) + 8*len(A)*len(A[0]) + 8*len(b) + 8*xlen + 4*vidlen
	diff := int(fsize) - calcsize
	if int(fsize) != calcsize {
		fmt.Printf("binDumpCheck Error - dump file size error, filesize: %d,\n\tlen(c): %d,\n\tlen(A): %d,\n\tLen(A[0]): %d,\n\tlen(b): %d,\n\txlen: %d, vidlen: %d\n\toverhead: %d\n\t\tdiff: %d\n", fsize, len(c), len(A), len(A[0]), len(b), xlen, vidlen, vidOverhead+xOverhead+3*12, diff)
	}
	c1, A1, b1, x1, vid1 := BinLoadModel(ftocheck)
	// Check loaded C vector
	if len(c) != len(c1) {
		fmt.Printf("modelio error: len(c): %d does not match len(c1) %d\n", len(c), len(c1))
	}
	for i := 0; i < len(c); i++ {
		if c[i] != c1[i] {
			fmt.Printf("c[%d] is %g but found %g\n", i, c[i], c1[i])
		}
	}
	// Checking A matrix
	if len(A) != len(A1) {
		fmt.Printf("modelio error: len(A): %d does not match len(A1) %d\n", len(A), len(A1))
	}
	for i := 0; i < len(A); i++ {
		if len(A[0]) != len(A1[i]) {
			fmt.Printf("modelio error: len(A[0]): %d does not match len(A1[%d]) %d\n", i, len(A[0]), len(A1[i]))
		}
		for j := 0; j < len(A[0]); j++ {
			if A[i][j] != A1[i][j] {
				fmt.Printf("A[%d][%d] is %g but found %g\n", i, j, A[i][j], A1[i][j])
			}
		}
	}
	// Checking b vector
	if len(b) != len(b1) {
		fmt.Printf("modelio error: len(b): %d does not match len(b1) %d\n", len(b), len(b1))
	}
	for i := 0; i < len(b); i++ {
		if b[i] != b1[i] {
			fmt.Printf("b[%d] is %g but found %g\n", i, b[i], b1[i])
		}
	}
	if x != nil && x1 != nil {
		// Checking x vector
		if len(x) != len(x1) {
			fmt.Printf("modelio error: len(x): %d does not match len(x1) %d\n", len(x), len(x1))
		}
		for i := 0; i < len(x); i++ {
			if x[i] != x1[i] {
				fmt.Printf("x[%d] is %g but found %g\n", i, x[i], x1[i])
			}
		}
	}
	if vid != nil && vid1 != nil {
		if (*vid)[0] != (*vid1)[0] ||
			(*vid)[1] != (*vid1)[1] ||
			(*vid)[2] != (*vid1)[2] ||
			(*vid)[3] != (*vid1)[3] ||
			(*vid)[4] != (*vid1)[4] ||
			(*vid)[5] != (*vid1)[5] {
			fmt.Printf("vid: %#v not eaqual to vid1: %#v\n", vid, vid1)
		}
	}
}

func BinCheckModelFiles(f1, f2 string) {
	var v *VectorVarIndex
	filex, err := os.Open(f1) // For read access.
	if err != nil {
		log.Fatal(err)
	}
	stats, err := filex.Stat()
	if err != nil {
		fmt.Printf("os.Stat: %s\n", err)
		os.Exit(1)
	}
	f1size := stats.Size()
	filex.Close()
	filex, err = os.Open(f2) // For read access.
	if err != nil {
		log.Fatal(err)
	}
	stats, err = filex.Stat()
	if err != nil {
		fmt.Printf("os.Stat: %s\n", err)
		os.Exit(1)
	}
	f2size := stats.Size()
	filex.Close()
	if f1size != f2size {
		fmt.Printf("Error - file sizes do not match %d vs %d\n", f1size, f2size)
	}
	c, A, b, x, vid := BinLoadModel(f2)
	c1, A1, b1, x1, vid1 := BinLoadModel(f1)

	// Check vid
	if vid != nil && vid1 != nil {
		if (*vid)[0] != (*vid1)[0] ||
			(*vid)[1] != (*vid1)[1] ||
			(*vid)[2] != (*vid1)[2] ||
			(*vid)[3] != (*vid1)[3] ||
			(*vid)[4] != (*vid1)[4] ||
			(*vid)[5] != (*vid1)[5] {
			fmt.Printf("vid: %#v not eaqual to vid1: %#v\n", vid, vid1)
			//TODO come up with a v to work in this case
		} else {

			years := int((*vid)[0])
			taxbins := int((*vid)[1])
			cgbins := int((*vid)[2])
			m := map[string]int{
				"IRA":      int((*vid)[3]),
				"roth":     int((*vid)[4]),
				"aftertax": int((*vid)[5]),
			}
			tmp, err := NewVectorVarIndex(years, taxbins,
				cgbins, m, os.Stdout)
			v = &tmp
			if err != nil {
				fmt.Printf("BinCheckModelFiles: %s\n", err)
				os.Exit(1)
			}
		}
	}

	// Check loaded C vector
	if len(c) != len(c1) {
		fmt.Printf("modelio error: len(c): %d does not match len(c1) %d\n", len(c), len(c1))
	}
	for i := 0; i < len(c); i++ {

		//if AlmostEqualRelativeAndAbs(c[i], c1[i], 0, 0)
		if c[i] > 0 && c1[i] < 0 || c[i] < 0 && c1[i] > 0 ||
			math.Abs(c[i]-c1[i]) > 0.00000001 {
			fmt.Printf("c[%d] is:\nf2: %.[4]*[2]g\nf1: %.[4]*[3]g\n", i, c[i], c1[i], 20)
			if v != nil {
				fmt.Printf("	%s\n", v.Varstr(i))
			}
		}
	}
	// Checking A matrix
	if len(A) != len(A1) {
		fmt.Printf("modelio error: len(A): %d does not match len(A1) %d\n", len(A), len(A1))
	}
	for i := 0; i < len(A); i++ {
		if len(A[0]) != len(A1[i]) {
			fmt.Printf("modelio error: len(A[0]): %d does not match len(A1[%d]) %d\n", i, len(A[0]), len(A1[i]))
		}
		for j := 0; j < len(A[0]); j++ {
			if A[i][j] > 0 && A1[i][j] < 0 || A[i][j] < 0 && A1[i][j] > 0 ||
				math.Abs(A[i][j]-A1[i][j]) > 0.00000001 {
				//if A[i][j] != A1[i][j]
				fmt.Printf("A[%d][%d] is:\nf2: %g\nf1: %g\n", i, j, A[i][j], A1[i][j])
				if v != nil {
					fmt.Printf("	%s\n", v.Varstr(j))
				}
			}
		}
	}
	// Checking b vector
	if len(b) != len(b1) {
		fmt.Printf("modelio error: len(b): %d does not match len(b1) %d\n", len(b), len(b1))
	}
	for i := 0; i < len(b); i++ {
		if b[i] > 0 && b1[i] < 0 || b[i] < 0 && b1[i] > 0 ||
			math.Abs(b[i]-b1[i]) > 0.00000001 {
			//if b[i] != b1[i]
			fmt.Printf("b[%d] is:\nf2: %g\nf1: %g\n", i, b[i], b1[i])
			fmt.Printf("	%s\n", "<= b")
		}
	}
	if x != nil && x1 != nil {
		// Checking x vector
		if len(x) != len(x1) {
			fmt.Printf("modelio error: len(x): %d does not match len(x1) %d\n", len(x), len(x1))
		}
		for i := 0; i < len(x); i++ {
			if x[i] > 0 && x1[i] < 0 || x[i] < 0 && x1[i] > 0 ||
				math.Abs(x[i]-x1[i]) > 0.00000001 {
				fmt.Printf("x[%d] is:\nf2: %g\nf1: %g\n", i, x[i], x1[i])
				if v != nil {
					fmt.Printf("	%s\n", v.Varstr(i))
				}
			}
		}
	} else {
		fmt.Printf("modelio warning: x or x1 is not present\n")
		fmt.Printf("x: %v\n", x)
		fmt.Printf("x1: %v\n", x1)
	}
}

//TODO: Make the err messages uniform
//TODO: unit test

//TODO: Make the err messages uniform
//TODO: unit test

// BinLoadModel reads a binary file, extracting c, A, b of a Linear Program
func BinLoadModel(filename string) ([]float64, [][]float64, []float64,
	[]float64, *[]int32) {
	if filename == "" {
		filename = "./RPlanModelpython.dat"
	}
	filem, err := os.Open(filename)
	if err != nil {
		fmt.Printf("os.Open: %s\n", err)
		os.Exit(1)
	}
	stats, err := filem.Stat()
	if err != nil {
		fmt.Printf("os.Stat: %s\n", err)
		os.Exit(1)
	}
	size := stats.Size()
	//fmt.Printf("File size is %d\n", size)

	Endian := binary.LittleEndian

	// Load C array
	header := make([]uint32, 3)
	err = binary.Read(filem, Endian, &header)
	if header[2] != 0xDEADBEEF {
		fmt.Printf("header code is not 0xDEADBEEF: %#X\n", header[2])
		os.Exit(1)
	}
	//fmt.Printf("c Header Rows: %d, Cols: %d, code: %#X\n", header[0], header[1], header[2])
	c := make([]float64, header[0])
	err = binary.Read(filem, Endian, &c)
	if err != nil {
		fmt.Printf("c binary.Read failed: %s\n", err)
		os.Exit(1)
	}

	// Load A matrix
	//header = make([]uint32, 3) // can I reuse the other header???
	err = binary.Read(filem, Endian, &header)
	if header[2] != 0xDEADBEEF {
		fmt.Printf("header code is not 0xDEADBEEF: %#X\n", header[2])
		os.Exit(1)
	}
	//fmt.Printf("A Header Rows: %d, Cols: %d, code: %#X\n", header[0], header[1], header[2])
	A := make([][]float64, 0)
	for i := 0; i < int(header[0]); i++ {
		row := make([]float64, header[1])
		err = binary.Read(filem, Endian, &row)
		if err != nil {
			fmt.Printf("A[%d] binary.Read failed: %s\n", i, err)
			os.Exit(1)
		}
		A = append(A, row)
	}

	// Load b array
	//header = make([]uint32, 3) // can I reuse the other header???
	err = binary.Read(filem, Endian, &header)
	if header[2] != 0xDEADBEEF {
		fmt.Printf("header code is not 0xDEADBEEF: %#X\n", header[2])
		os.Exit(1)
	}
	//fmt.Printf("b Header Rows: %d, Cols: %d, code: %#X\n", header[0], header[1], header[2])
	b := make([]float64, header[0])
	err = binary.Read(filem, Endian, &b)
	if err != nil {
		fmt.Printf("b binary.Read failed: %s\n", err)
		os.Exit(1)
	}

	// Load x array
	//header = make([]uint32, 3) // can I reuse the other header???
	x := []float64(nil)
	xlen := 0
	xOverhead := 0
	err = binary.Read(filem, Endian, &header)
	if err != nil {
		fmt.Printf("x binary.Read failed: %s\n", err)
	} else {
		if header[2] != 0xDEADBEEF {
			fmt.Printf("header code is not 0xDEADBEEF: %#X\n", header[2])
			os.Exit(1)
		}
		//fmt.Printf("x Header Rows: %d, Cols: %d, code: %#X\n", header[0], header[1], header[2])
		x = make([]float64, header[0])
		err = binary.Read(filem, Endian, &x)
		if err != nil {
			fmt.Printf("x binary.Read failed: %s\n", err)
			os.Exit(1)
		}
		xlen = len(x)
		xOverhead = 12
		//fmt.Printf("X: %v\n", x)
	}

	//
	// Currently Assumes X or no VIndexData // FIXME TODO
	//

	// Load Vector Index Data array
	vida := make([]int32, 6)
	vid := &vida
	vidlen := 0
	vidOverhead := 0
	err = binary.Read(filem, Endian, &header)
	if err != nil {
		fmt.Printf("x binary.Read failed: %s\n", err)
	} else {
		if header[2] != 0xDEADBEEF {
			fmt.Printf("header code is not 0xDEADBEEF: %#X\n", header[2])
			os.Exit(1)
		}
		//fmt.Printf("x Header Rows: %d, Cols: %d, code: %#X\n", header[0], header[1], header[2])
		err = binary.Read(filem, Endian, vid)
		if err != nil {
			fmt.Printf("x binary.Read failed: %s\n", err)
			os.Exit(1)
		}
		vidlen = len(*vid)
		vidOverhead = 12
		//fmt.Printf("X: %v\n", x)
	}
	filem.Close()
	contentsize := xOverhead + vidOverhead + 3*12 + len(c)*8 + len(A)*len(A[0])*8 + len(b)*8 + 8*xlen + 4*vidlen
	if size != int64(contentsize) {
		fmt.Printf("BinLoadModel: file size (%d) and content size (%d) do not match\n", size, contentsize)
	}
	return c, A, b, x, vid
}

/*
//
// The following two functions come from the following URL:
// https://randomascii.wordpress.com/2012/02/25/comparing-floating-point-numbers-2012-edition/
//
type Float_t struct {
    Float_t(float num = 0.0f) : f(num) {}
    // Portable extraction of components.

    int32_t i;
    float f;
    struct
    {   // Bitfields for exploration. Do not use in production code.
        uint32_t mantissa : 23;
        uint32_t exponent : 8;
        uint32_t sign : 1;
    } parts;
}
func bool Negative() const { return i < 0; }
func int32_t RawMantissa() const { return i & ((1 << 23) - 1); }
func int32_t RawExponent() const { return (i >> 23) & 0xFF; }

func AlmostEqualUlpsAndAbs(A, B, float64,
	maxDiff float64, maxUlpsDiff int) bool {
	// Check if the numbers are really close -- needed
	// when comparing numbers near zero.
	absDiff := fabs(A - B)
	if absDiff <= maxDiff {
		return true
	}

	Float_t uA(A);
	Float_t uB(B);

	// Different signs means they do not match.
	if uA.Negative() != uB.Negative() {
		return false;
	}

	// Find the difference in ULPs.
	int ulpsDiff = abs(uA.i - uB.i);
	if ulpsDiff <= maxUlpsDiff {
		return true;
	}
	return false;
}
*/

func AlmostEqualRelativeAndAbs(A float64, B float64, maxDiff float64, maxReldiff float64) bool { // FLT_EPSILON
	// Check if the numbers are really close -- needed
	// when comparing numbers near zero.
	if maxReldiff == 0.0 {
		maxReldiff = 1.19e-7
	}
	if maxDiff == 0.0 {
		maxDiff = 1.19e-7
	}
	diff := math.Abs(A - B)
	if diff <= maxDiff {
		return true
	}
	A = math.Abs(A)
	B = math.Abs(B)
	largest := A
	if B > A {
		largest = B
	}
	if diff <= largest*maxReldiff {
		return true
	}
	return false
}
