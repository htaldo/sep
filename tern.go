package sep

import (
	"bufio"
	"fmt"
	"os"
)

type Stream struct { 
	m,c float64 //c is the composition of C
	x Point		//x represents the compositions of A and B
}
type Diagram struct { EqLines, Alders []Segment }

func Efficiency(F, E Stream) float64 {
	return ((E.c * E.m) / (F.c * F.m))
}

func Partition1(E, R Stream) float64 {
	return E.c/R.c
}

func Partition2(E, R Stream) float64 {
	return E.x.X/R.x.X
}

func Selectivity(E, R Stream) float64 {
	return Partition1(E,R)/Partition2(E,R)
}

func LleSimple(F, S Stream, d Diagram) (M, R, E Stream){
	M, R, E = mixBal(F, S, d.EqLines)	
	//interpSeg := Segment{ R.x, E.x }
	//no need to print interpseg; it's R and E
	//println("> M, R, E")
	//M.Print()
	//R.Print()
	//E.Print()
	return M, R, E
}

func LleDirect(F0, S Stream, d Diagram, finalExtC float64) {
	//append F0 and S0 to F and S
	var F, M, R, E []Stream
	var nextM, nextR, nextE Stream
	var interpSeg []Segment
	var currentEc float64 = 1.0
	F = append(F, F0)
	for i := 0; currentEc > finalExtC; i++ {
		fmt.Printf("%d\n", i)
		nextM, nextR, nextE = mixBal(F[i], S, d.EqLines)	
		M, R, E = append(M, nextM), append(R, nextR), append(E, nextE)
		interpSeg = append(interpSeg, Segment{ R[i].x, E[i].x })
		interpSeg[i].Print()
		println("> M, R, E")
		fmt.Printf("%v\n", M)
		fmt.Printf("%v\n", R)
		fmt.Printf("%v\n", E)
		currentEc = E[i].c 
		F = append(F, R[i])
	}
	//F.c, S.c = To3(F.x)[2], To3(S.x)[2]
}

//func LleCounter(F, S Stream, eqLines Segment[], Rn Stream) {
//	var E, R, E0 []Stream
//	Rn.x := schnittRaf(Line{-1, 1-c})
//	E0.x = schnittExt(PPLine(Rn.x,M.x))	
//	E = append(E, E0)
//	P := PPLine(F.x,E0).SchnittLine(PPLine(Rn.x, S.x))
//	var nextRx, nextEx Point
//	nextRx = Point{0, 0}
//	for i := 0; To3(nextRx)[2] > Rn.c; i++ {
//		nextRx = alders(E[i].x, d.Alders)
//		R = append(Rcomps, nextRcomp)
//		nextEcomp = SchnittExt(line(nextRcomp, P))
//		Ecomps = append(Ecomps, nextEcomp)
//	}
//	//TODO: Discard the last Ecomp, or return Ecomps[i]
//}

//func schnittRaf (l Line, eqLines []Segment) Point {
//	var sch Point
//	for i := 0; i < len(eqLines) - 1; i++ {
//		rafSeg := Segment{eqLines[i].P, eqLines[i+1].P}	
//		sch = l.SchnittLine(rafSeg.Line())	
//		//check orientation of P and Q
//		if twopBound(sch, rafSeg.P, rafSeg.Q) == true {
//			break
//		}
//	}
//	return sch
//}

//func schnittExt (l Line, eqLines []Segment) Point {
//	var sch Point
//	for i := 0; i < len(eqLines) - 1; i++ {
//		extSeg := Segment{eqLines[i].Q, eqLines[i+1].Q}	
//		sch = l.SchnittLine(extSeg.Line())	
//		//check orientation of P and Q
//		if twopBound(sch, extSeg.P, extSeg.Q) == true {
//			break
//		}
//	}
//	return sch
//}

func mixBal (F, S Stream, eqlines []Segment) (Stream, Stream, Stream) {
	//it's probably fine to pass the whole diagram (or a pointer to it),
	//instead of just eqlines
	var M Stream
	M.m = F.m + S.m
	M.c = (F.m*F.c + S.m*S.c)/(M.m)
	Mline := Line{ -1, 1 - M.c }
	FSline := PPLine(S.x, F.x)
	M.x = Mline.SchnittLine(FSline)
	//find M.x between two equilibrium lines
	//this assumes the lines are ordered by decreasing c
	var quadpoints [4]Point
	var s1, s2 Segment
	for i, j := 0, 1; j < len(eqlines); i, j = i+1, j+1 {
		quadpoints[0], quadpoints[1] = eqlines[i].Points()
		quadpoints[2], quadpoints[3] = eqlines[j].Points()
		if M.x.InQuad(quadpoints) == true {
			s1, s2 = eqlines[i], eqlines[j]
			break
		}
		//holi c:
	}
	l1, l2 := s1.Line(), s2.Line()
	interpPoint := l1.SchnittLine(l2)
	interpLine := PPLine(interpPoint, M.x)
	//use the first 2 points of the quadrilateral
	//(the segment of the extraction curve)
	//to intersect with interpLine
	//this assumes both segments have slope >= 0
	rafSeg := Segment{ s1.P, s2.P }
	extSeg := Segment{ s1.Q, s2.Q }
	var R, E Stream
	//this assumes interpLine will intersect with extSeg & rafSeg
	R.x = interpLine.SchnittLine(rafSeg.Line())
	E.x = interpLine.SchnittLine(extSeg.Line())
	R.c, E.c = To3(R.x)[2], To3(E.x)[2]
	R.m = M.m * (M.c - E.c)/(R.c - E.c)
	E.m = M.m - R.m
	return M, R, E
}

// beg text

func (s Stream) Print() {
	fmt.Printf("%.1f\t%.4f\t%.4f\t%.4f\n", s.m, s.x.X, s.x.Y, s.c)
}

func readDiagram() Diagram {
	var d Diagram
	d.EqLines = readSegments("reparto")	
	d.Alders = readSegments("alders")	
	return d
}

func readFS() (F, S Stream) {
	F = readStream("feed")	
	S = readStream("solvent")	
	return F, S
}

func readStream(filename string) Stream {
	var s Stream
	f, _ := os.Open(filename)
	input := bufio.NewScanner(f)
	var fields []float64
	for input.Scan() {
		fields = splitRow(input.Text(), 4)	
		s.x, s.m = To2(fields[0:2]), fields[3]
	}
	f.Close()
	return s
}

// end text
