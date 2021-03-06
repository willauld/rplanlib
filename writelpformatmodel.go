package rplanlib

import (
	"fmt"
	"os"
)

//var modelfile *os.File

//
// This is to write an LP Format model readable by lp_solve
//

// TODO: FIXME: Create UNIT tests: last two parameters need s vector (s is output from simplex run)

// PrintModelMatrix prints to object function (cx) and constraint matrix (Ax<=b)
func (ms ModelSpecs) WriteLPFormatModel(c []float64, A [][]float64, b []float64, notes []ModelNote, filename string, row []float64, cmdline string) error {
	modelfile, err := os.Create(filename)
	if err != nil {
		e := fmt.Errorf("Could not create new model file %s", filename)
		return e
	}
	fmt.Fprintf(modelfile, "\n// %s \n", cmdline)
	note := ""
	notesIndex := 0
	nextModelIndex := len(A) + 1 // beyond the end of A
	if notes != nil {
		nextModelIndex = notes[notesIndex].index
		note = notes[notesIndex].note
		notesIndex++
	}
	if nextModelIndex < 0 { // Object function index -1
		from := nextModelIndex
		nextModelIndex = notes[notesIndex].index
		to := nextModelIndex - 1
		fmt.Fprintf(modelfile, "\n// ##== [%d-%d]: %s ==##\n", from, to, note)
		note = notes[notesIndex].note
		notesIndex++
	}
	fmt.Fprintf(modelfile, "min: ")
	ms.writeModelRow(c, false, modelfile)
	fmt.Fprintf(modelfile, "\n")
	for constraint := 0; constraint < len(A); constraint++ {
		if nextModelIndex == constraint {
			from := nextModelIndex
			nextModelIndex = notes[notesIndex].index
			to := nextModelIndex - 1
			for to < from {
				fmt.Fprintf(modelfile, "\n// ##== [%d-%d]: %s ==##\n", from, to, note)
				note = notes[notesIndex].note
				notesIndex++
				from = nextModelIndex
				nextModelIndex = notes[notesIndex].index
				to = nextModelIndex - 1
			}
			fmt.Fprintf(modelfile, "\n// ##== [%d-%d]: %s ==##\n", from, to, note)
			note = notes[notesIndex].note
			notesIndex++
		}
		//fmt.Fprintf(modelfile, "%3d: ", constraint)
		ms.writeConstraint(A[constraint], b[constraint], modelfile)
	}
	fmt.Fprintf(modelfile, "\n")
	ms.WriteObjectFunctionSolution(c, row, modelfile)
	return nil
}

func (ms ModelSpecs) writeConstraint(row []float64, b float64, modelfile *os.File) {
	ms.writeModelRow(row, true, modelfile)
	fmt.Fprintf(modelfile, "<= %10.4f;\n", b)
}

func (ms ModelSpecs) writeModelRow(row []float64, suppressNewline bool, modelfile *os.File) {
	if ms.Ip.Numacc < 0 || ms.Ip.Numacc > 5 {
		e := fmt.Errorf("PrintModelRow: number of accounts is out of bounds, should be between [0, 5] but is %d", ms.Ip.Numacc)
		fmt.Fprintf(modelfile, "%s\n", e)
		return
	}
	for i := 0; i < ms.Ip.Numyr; i++ { // x[]
		for k := 0; k < len(*ms.Ti.Taxtable); k++ {
			if row[ms.Vindx.X(i, k)] != 0 {
				fmt.Fprintf(modelfile, "+ %10.4f x.%d.%d ", row[ms.Vindx.X(i, k)], i, k)
			}
		}
	}
	if ms.Ip.Accmap[Aftertax] > 0 {
		for i := 0; i < ms.Ip.Numyr; i++ { // sy[]
			for l := 0; l < len(*ms.Ti.Capgainstable); l++ {
				if row[ms.Vindx.Sy(i, l)] != 0 {
					fmt.Fprintf(modelfile, "+ %10.4f sy.%d.%d ", row[ms.Vindx.Sy(i, l)], i, l)
				}
			}
		}
		for i := 0; i < ms.Ip.Numyr; i++ { // y[]
			for l := 0; l < len(*ms.Ti.Capgainstable); l++ {
				if row[ms.Vindx.Y(i, l)] != 0 {
					fmt.Fprintf(modelfile, "+ %10.4f y.%d.%d ", row[ms.Vindx.Y(i, l)], i, l)
				}
			}
		}
	}
	for i := 0; i < ms.Ip.Numyr; i++ { // w[]
		for j := 0; j < ms.Ip.Numacc; j++ {
			if row[ms.Vindx.W(i, j)] != 0 {
				fmt.Fprintf(modelfile, "+ %10.4f w.%d.%d ", row[ms.Vindx.W(i, j)], i, j)
			}
		}
	}
	for i := 0; i < ms.Ip.Numyr+1; i++ { // b[] has an extra year
		for j := 0; j < ms.Ip.Numacc; j++ {
			if row[ms.Vindx.B(i, j)] != 0 {
				fmt.Fprintf(modelfile, "+ %10.4f b.%d.%d ", row[ms.Vindx.B(i, j)], i, j)
			}
		}
	}
	for i := 0; i < ms.Ip.Numyr; i++ { // s[]
		if row[ms.Vindx.S(i)] != 0 {
			fmt.Fprintf(modelfile, "+ %10.4f s.%d ", row[ms.Vindx.S(i)], i)
		}
	}
	for i := 0; i < ms.Ip.Numyr; i++ { // D[]
		if row[ms.Vindx.D(i)] != 0 {
			fmt.Fprintf(modelfile, "+ %10.4f D.%d ", row[ms.Vindx.D(i)], i)
		}
	}
	if !suppressNewline {
		fmt.Fprintf(modelfile, ";\n")
	}
}

func (ms ModelSpecs) WriteObjectFunctionSolution(c []float64, row []float64, modelfile *os.File) {
	//
	// if modelfile is nil use ms.Ao.Output to write to wherever it points
	// if it is not nil create an AppOutput() pointing to it and use that
	//
	objOut := ms.Ao
	writeAll := true // write the complex c and res.X vectors
	if modelfile != nil {
		objOut = NewAppOutput(nil, modelfile)
		writeAll = false
	}
	if ms.Ip.Numacc < 0 || ms.Ip.Numacc > 5 {
		e := fmt.Errorf("PrintObjectFunc: number of accounts is out of bounds, should be between [0, 5] but is %d", ms.Ip.Numacc)
		objOut.Output(fmt.Sprintf("%s\n", e))
		return
	}
	objOut.Output(fmt.Sprintf("/* LPSimplex solution: \n"))
	localSum := 0.0
	globalSum := 0.0
	for i := 0; i < ms.Ip.Numyr; i++ { // x[]
		for k := 0; k < len(*ms.Ti.Taxtable); k++ {
			cIndx := ms.Vindx.X(i, k)
			if c[cIndx] != 0 || writeAll { // should all values for c and vars
				cXrow := c[cIndx] * row[cIndx]
				localSum += cXrow
				objOut.Output(fmt.Sprintf("C[%d]=@%6.3f@&*@&x[%d,%d]=@%6.3f@&== @&%6.3f\n", cIndx, c[cIndx], i, k, row[cIndx], cXrow))
			}
		}
	}
	objOut.Output(fmt.Sprintf("\tSum Ci*Xi == %6.3f\n", localSum))
	globalSum += localSum
	localSum = 0.0
	if ms.Ip.Accmap[Aftertax] > 0 {
		for i := 0; i < ms.Ip.Numyr; i++ { // sy[]
			for l := 0; l < len(*ms.Ti.Capgainstable); l++ {
				cIndx := ms.Vindx.Sy(i, l)
				if c[cIndx] != 0 || writeAll {
					cXrow := c[cIndx] * row[cIndx]
					localSum += cXrow
					objOut.Output(fmt.Sprintf("C[%d]=@%6.3f@&*@&Sy[%d,%d]=@%6.3f@&==@&%6.3f\n", cIndx, c[cIndx], i, l, row[cIndx], cXrow))
				}
			}
		}
		objOut.Output(fmt.Sprintf("\tSum Ci*Syi@&==@&%6.3f\n", localSum))
		globalSum += localSum
		localSum = 0.0
		for i := 0; i < ms.Ip.Numyr; i++ { // y[]
			for l := 0; l < len(*ms.Ti.Capgainstable); l++ {
				cIndx := ms.Vindx.Y(i, l)
				if c[cIndx] != 0 || writeAll {
					cXrow := c[cIndx] * row[cIndx]
					localSum += cXrow
					objOut.Output(fmt.Sprintf("C[%d]=@%6.3f@&*@&Y[%d,%d]=@%6.3f@&==@&%6.3f\n", cIndx, c[cIndx], i, l, row[cIndx], cXrow))
				}
			}
		}
		objOut.Output(fmt.Sprintf("\tSum Ci*Yi@&==@&%6.3f\n", localSum))
		globalSum += localSum
		localSum = 0.0
	}
	for i := 0; i < ms.Ip.Numyr; i++ { // w[]
		for j := 0; j < ms.Ip.Numacc; j++ {
			cIndx := ms.Vindx.W(i, j)
			if c[cIndx] != 0 || writeAll {
				cXrow := c[cIndx] * row[cIndx]
				localSum += cXrow
				objOut.Output(fmt.Sprintf("C[%d]=@%6.3f@&*@&w[%d,%d]=@%6.3f@&==@&%6.3f\n", cIndx, c[cIndx], i, j, row[cIndx], cXrow))
			}
		}
	}
	objOut.Output(fmt.Sprintf("\tSum Ci*wi@&==@&%6.3f\n", localSum))
	globalSum += localSum
	localSum = 0.0
	for i := 0; i < ms.Ip.Numyr+1; i++ { // b[] has an extra year
		for j := 0; j < ms.Ip.Numacc; j++ {
			cIndx := ms.Vindx.B(i, j)
			if c[cIndx] != 0 || writeAll {
				cXrow := c[cIndx] * row[cIndx]
				localSum += cXrow
				objOut.Output(fmt.Sprintf("C[%d]=@%6.3f@&*@&b[%d,%d]=@%6.3f@&==@&%6.3f\n", cIndx, c[cIndx], i, j, row[cIndx], cXrow))
			}
		}
	}
	objOut.Output(fmt.Sprintf("\tSum Ci*bi@&==@&%6.3f\n", localSum))
	globalSum += localSum
	localSum = 0.0
	for i := 0; i < ms.Ip.Numyr; i++ { // s[]
		cIndx := ms.Vindx.S(i)
		if c[cIndx] != 0 || writeAll {
			cXrow := c[cIndx] * row[cIndx]
			localSum += cXrow
			objOut.Output(fmt.Sprintf("C[%d]=@%6.3f@&*@&S[%d]=@%6.3f@&==@&%6.3f\n", cIndx, c[cIndx], i, row[cIndx], cXrow))
		}
	}
	objOut.Output(fmt.Sprintf("\tSum Ci*Si@&==@&%6.3f\n", localSum))
	globalSum += localSum
	localSum = 0.0
	for i := 0; i < ms.Ip.Numyr; i++ { // D[]
		cIndx := ms.Vindx.D(i)
		if c[cIndx] != 0 || writeAll {
			cXrow := c[cIndx] * row[cIndx]
			localSum += cXrow
			objOut.Output(fmt.Sprintf("C[%d]=@%6.3f@&*@&D[%d]=@%6.3f@&==@&%6.3f\n", cIndx, c[cIndx], i, row[cIndx], cXrow))
		}
	}
	objOut.Output(fmt.Sprintf("\tSum Ci*Di@&==@&%6.3f\n", localSum))
	globalSum += localSum
	objOut.Output(fmt.Sprintf("\t\tSum overall@&==@&%6.3f\n", globalSum))
	localSum = 0.0
	objOut.Output(fmt.Sprintf(" End LPSimplex solution */\n"))
}
