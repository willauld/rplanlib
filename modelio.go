package rplanlib

import (
	"encoding/binary"
	"fmt"
	"log"
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

func BinDumpModel(c []float64, A [][]float64, b []float64, fname string) error {
	if fname == "" {
		fname = "./RPlanModelgo.dat"
	}
	filem, err := os.Create(fname)
	if err != nil {
		e := fmt.Errorf("os.Open: %s\n", err)
		return e
	}
	Endian := binary.LittleEndian
	//stream = open(fname, 'wb')

	// write C array
	//header := make([]uint32, 3)
	header := []uint32{uint32(len(c)), 0, 0xDEADBEEF}
	err = binary.Write(filem, Endian, &header)
	if err != nil {
		e := fmt.Errorf("BinDumpModel(1): %s\n", err)
		return e
	}
	err = binary.Write(filem, Endian, &c)
	if err != nil {
		e := fmt.Errorf("BinDumpModel(2): %s\n", err)
		return e
	}
	//fmt.Printf("c Header Rows: %d, Cols: %d, code: %#X\n", header[0], header[1], header[2])

	header = []uint32{uint32(len(A)), uint32(len(A[0])), 0xDEADBEEF}
	//fmt.Printf("A length: %d, %d, dumping", (len(A), len(A[0])))
	err = binary.Write(filem, Endian, &header)
	if err != nil {
		e := fmt.Errorf("BinDumpModel(3): %s\n", err)
		return e
	}
	for _, row := range A {
		err = binary.Write(filem, Endian, &row)
		if err != nil {
			e := fmt.Errorf("BinDumpModel(4): %s\n", err)
			return e
		}
	}
	header = []uint32{uint32(len(b)), 0, 0xDEADBEEF}
	//fmt.Printf("b length: %d, dumping", len(b))
	err = binary.Write(filem, Endian, &header)
	if err != nil {
		e := fmt.Errorf("BinDumpModel(5): %s\n", err)
		return e
	}
	err = binary.Write(filem, Endian, &b)
	if err != nil {
		e := fmt.Errorf("BinDumpModel(6): %s\n", err)
		return e
	}

	//filem.Close()

	stats, err := filem.Stat()
	if err != nil {
		e := fmt.Errorf("os.Stat: %s\n", err)
		return e
	}
	fsize := stats.Size()

	overhead := 3 * 12
	csize := 8 * len(c)
	Asize := 8 * len(A) * len(A[0])
	bsize := 8 * len(b)
	if fsize != int64(overhead+csize+Asize+bsize) {
		fmt.Printf("Error - dump file size error, filesize: %d, len(c): %d, len(A): %d, Len(A[0]): %d, len(b): %d\n", fsize, len(c), len(A), len(A[0]), len(b))
	}
	fmt.Printf("See Me - dump file size error, filesize: %d, len(c): %d, len(A): %d, Len(A[0]): %d, len(b): %d\n", fsize, len(c), len(A), len(A[0]), len(b))
	binDumpCheck(c, A, b, fname)
	return nil
}

func binDumpCheck(c []float64, A [][]float64, b []float64, ftocheck string) {
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
	if int(fsize) != 3*12+8*len(c)+8*len(A)*len(A[0])+8*len(b) {
		fmt.Printf("Error - dump file size error, filesize: %d, len(c): %d, len(A): %d, Len(A[0]): %d, len(b): %d\n", fsize, len(c), len(A), len(A[0]), len(b))
	}
	c1, A1, b1 := BinLoadModel(ftocheck)
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
}

func BinCheckModelFiles(f1, f2 string) {
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
	c, A, b := BinLoadModel(f2)
	c1, A1, b1 := BinLoadModel(f1)
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
}

//TODO: Make the err messages uniform
//TODO: unit test

//TODO: Make the err messages uniform
//TODO: unit test

// BinLoadModel reads a binary file, extracting c, A, b of a Linear Program
func BinLoadModel(filename string) ([]float64, [][]float64, []float64) {
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

	filem.Close()
	contentsize := 3*12 + len(c)*8 + len(A)*len(A[0])*8 + len(b)*8
	if size != int64(contentsize) {
		fmt.Printf("BinLoadModel: file size (%d) and content size (%d) do not match\n", size, contentsize)
	}
	return c, A, b
}
