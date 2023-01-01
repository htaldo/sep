package sep

import (
	"bufio"
	"fmt"
	"os"
)

type Stream struct { 
	Mass,SComp float64 //SComp is the solute composition 
	Comp Point		//Comp represents the compositions of A and B
}
type Diagram struct { EqLines, Alders []Segment }

func Efficiency(F, E Stream) float64 {
	return ((E.SComp * E.Mass) / (F.SComp * F.Mass))
}

func Partition1(E, R Stream) float64 {
	return E.SComp/R.SComp
}

func Partition2(E, R Stream) float64 {
	return E.Comp.X/R.Comp.X
}

func Selectivity(E, R Stream) float64 {
	return Partition1(E,R)/Partition2(E,R)
}

func LleSimple(F, S Stream, d Diagram) (M, R, E Stream){
	M, R, E = mixBal(F, S, d.EqLines)	
	//interpSeg := Segment{ R.Comp, E.Comp }
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
		interpSeg = append(interpSeg, Segment{ R[i].Comp, E[i].Comp })
		interpSeg[i].Print()
		println("> M, R, E")
		fmt.Printf("%v\n", M)
		fmt.Printf("%v\n", R)
		fmt.Printf("%v\n", E)
		currentEc = E[i].SComp 
		F = append(F, R[i])
	}
	//F.SComp, S.SComp = To3(F.Comp)[2], To3(S.Comp)[2]
}

//func LleCounter(F, S Stream, eqLines Segment[], Rn Stream) {
//	var E, R, E0 []Stream
//	Rn.Comp := schnittRaf(Line{-1, 1-c})
//	E0.Comp = schnittExt(PPLine(Rn.Comp,M.Comp))	
//	E = append(E, E0)
//	P := PPLine(F.Comp,E0).SchnittLine(PPLine(Rn.Comp, S.Comp))
//	var nextRx, nextEx Point
//	nextRx = Point{0, 0}
//	for i := 0; To3(nextRx)[2] > Rn.SComp; i++ {
//		nextRx = alders(E[i].Comp, d.Alders)
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
	M.Mass = F.Mass + S.Mass
	M.SComp = (F.Mass*F.SComp + S.Mass*S.SComp)/(M.Mass)
	Mline := Line{ -1, 1 - M.SComp }
	FSline := PPLine(S.Comp, F.Comp)
	M.Comp = Mline.SchnittLine(FSline)
	//find M.Comp between two equilibrium lines
	//this assumes the lines are ordered by decreasing c
	var quadpoints [4]Point
	var s1, s2 Segment
	for i, j := 0, 1; j < len(eqlines); i, j = i+1, j+1 {
		quadpoints[0], quadpoints[1] = eqlines[i].Points()
		quadpoints[2], quadpoints[3] = eqlines[j].Points()
		if M.Comp.InQuad(quadpoints) == true {
			s1, s2 = eqlines[i], eqlines[j]
			break
		}
		//holi c:
	}
	l1, l2 := s1.Line(), s2.Line()
	interpPoint := l1.SchnittLine(l2)
	interpLine := PPLine(interpPoint, M.Comp)
	//use the first 2 points of the quadrilateral
	//(the segment of the extraction curve)
	//to intersect with interpLine
	//this assumes both segments have slope >= 0
	rafSeg := Segment{ s1.P, s2.P }
	extSeg := Segment{ s1.Q, s2.Q }
	var R, E Stream
	//this assumes interpLine will intersect with extSeg & rafSeg
	R.Comp = interpLine.SchnittLine(rafSeg.Line())
	E.Comp = interpLine.SchnittLine(extSeg.Line())
	R.SComp, E.SComp = To3(R.Comp)[2], To3(E.Comp)[2]
	R.Mass = M.Mass * (M.SComp - E.SComp)/(R.SComp - E.SComp)
	E.Mass = M.Mass - R.Mass
	return M, R, E
}

// beg text

func (s Stream) Print() {
	fmt.Printf("%.1f\t%.4f\t%.4f\t%.4f\n", s.Mass, s.Comp.X, s.Comp.Y, s.SComp)
}

func ReadDiagram() Diagram {
	var d Diagram
	d.EqLines = ReadSegments("reparto")	
	d.Alders = ReadSegments("alders")	
	return d
}

func ReadFS() (F, S Stream) {
	F = ReadStream("feed")	
	S = ReadStream("solvent")	
	return F, S
}

func ReadStream(filename string) Stream {
	var s Stream
	f, _ := os.Open(filename)
	input := bufio.NewScanner(f)
	var fields []float64
	for input.Scan() {
		fields = SplitRow(input.Text(), 4)	
		s.Comp, s.Mass = To2(fields[0:2]), fields[3]
	}
	f.Close()
	return s
}

// end text
