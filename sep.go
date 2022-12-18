package sep

import (
	"regexp"
	"bufio"
	"fmt"
	"strconv"
	"math"
	"sort"
	"os"
)

type Point struct { X, Y float64 }
type Segment struct { P, Q Point }
type Line struct { A, B float64}
type Horizontal float64
type Vertical float64
type Stream struct { 
	m,c float64 //c is the composition of C
	x Point		//x represents the compositions of A and B
}
type Diagram struct { EqLines, Alders []Segment }

const sq2 = 1.41421356237

func main() {
	F, S := readFS()
	F.c, S.c = to3(F.x)[2], to3(S.x)[2]
	diagram := readDiagram()
	_, R, E := simple(F, S, diagram)

	F.Print()
	S.Print()
	R.Print()
	E.Print()
	fmt.Printf("%.2f%%\n", efficiency(F, E)*100)
	fmt.Printf("%.3f\t%.3f\t%.3f\n", partition1(E,R), partition2(E,R), selectivity(E,R))


//	minF, maxF := 100.0, 1000.0
//	minS, maxS := 100.0, 1000.0
//	pasoF := (maxF - minF) / 100.0
//	pasoS := (maxS - minS) / 100.0
//	var usage, eff, cost float64
//	var usageScore, effScore, costScore, score float64
//	maxusage, maxeff, maxcost := 2000.0, 1.0, 10100.0
//	for F.m = minF; F.m < maxF; F.m += pasoF {
//		for S.m = minS; S.m < maxS; S.m += pasoS {
//			_, _, E := simple(F, S, diagram)
//			usage = F.m + S.m
//			cost = 0.1*F.m + 10*S.m 
//			eff = efficiency(F, E)
//			usageScore = 1 - (usage/maxusage)
//			effScore = 1 - (eff/maxeff)
//			costScore = cost/maxcost
//			score = (usageScore + costScore + effScore)/3
//			fmt.Printf("%f, ", score)
//		}
//		fmt.Printf("\n")
//	}
}

func efficiency(F, E Stream) float64 {
	return ((E.c * E.m) / (F.c * F.m))
}

func partition1(E, R Stream) float64 {
	return E.c/R.c
}

func partition2(E, R Stream) float64 {
	return E.x.X/R.x.X
}

func selectivity(E, R Stream) float64 {
	return partition1(E,R)/partition2(E,R)
}

func simple(F, S Stream, d Diagram) (M, R, E Stream){
	M, R, E = mixBal(F, S, d.EqLines)	
	//interpSeg := Segment{ R.x, E.x }
	//no need to print interpseg; it's R and E
	//println("> M, R, E")
	//M.Print()
	//R.Print()
	//E.Print()
	return M, R, E
}

func direct(F0, S Stream, d Diagram, finalExtC float64) {
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
	//F.c, S.c = to3(F.x)[2], to3(S.x)[2]
}

//func counter(F, S Stream, eqLines Segment[], Rn Stream) {
//	var E, R, E0 []Stream
//	Rn.x := schnittRaf(Line{-1, 1-c})
//	E0.x = schnittExt(segLine(Rn.x,M.x))	
//	E = append(E, E0)
//	P := segLine(F.x,E0).schnitt(segLine(Rn.x, S.x))
//	var nextRx, nextEx Point
//	nextRx = Point{0, 0}
//	for i := 0; to3(nextRx)[2] > Rn.c; i++ {
//		nextRx = alders(E[i].x, d.Alders)
//		R = append(Rcomps, nextRcomp)
//		nextEcomp = schnittExt(line(nextRcomp, P))
//		Ecomps = append(Ecomps, nextEcomp)
//	}
//	//TODO: Discard the last Ecomp, or return Ecomps[i]
//}

//func schnittRaf (l Line, eqLines []Segment) Point {
//	var sch Point
//	for i := 0; i < len(eqLines) - 1; i++ {
//		rafSeg := Segment{eqLines[i].P, eqLines[i+1].P}	
//		sch = l.schnitt(segLine(rafSeg))	
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
//		sch = l.schnitt(segLine(extSeg))	
//		//check orientation of P and Q
//		if twopBound(sch, extSeg.P, extSeg.Q) == true {
//			break
//		}
//	}
//	return sch
//}

//func twopBound(o, P, Q Point) bool {
//	//write between()
//	xcoords := [2]float64{P.X, Q.X}
//	ycoords := [2]float64{P.Y, Q.Y}
//	sortedx := sort.Ints(xcoords)
//	sortedy := sort.Ints(ycoords)
//	if o.X > sortedx[0] && o.X < sortedx[1] {
//		if o.Y > sortedy[0] && o.Y < sortedy[1] {
//			return true
//		}
//	}
//}

func mixBal (F, S Stream, eqlines []Segment) (Stream, Stream, Stream) {
	//it's probably fine to pass the whole diagram (or a pointer to it),
	//instead of just eqlines
	var M Stream
	M.m = F.m + S.m
	M.c = (F.m*F.c + S.m*S.c)/(M.m)
	Mline := Line{ -1, 1 - M.c }
	FSline := ppline(S.x, F.x)
	M.x = Mline.schnitt(FSline)
	//find M.x between two equilibrium lines
	//this assumes the lines are ordered by decreasing c
	var quadpoints [4]Point
	var s1, s2 Segment
	for i, j := 0, 1; j < len(eqlines); i, j = i+1, j+1 {
		quadpoints[0], quadpoints[1] = eqlines[i].Points()
		quadpoints[2], quadpoints[3] = eqlines[j].Points()
		if M.x.inQuad(quadpoints) == true {
			s1, s2 = eqlines[i], eqlines[j]
			break
		}
		//holi c:
	}
	l1, l2 := segLine(s1), segLine(s2)
	interpPoint := l1.schnitt(l2)
	interpLine := ppline(interpPoint, M.x)
	//use the first 2 points of the quadrilateral
	//(the segment of the extraction curve)
	//to intersect with interpLine
	//this assumes both segments have slope >= 0
	rafSeg := Segment{ s1.P, s2.P }
	extSeg := Segment{ s1.Q, s2.Q }
	var R, E Stream
	//this assumes interpLine will intersect with extSeg & rafSeg
	R.x = interpLine.schnitt(segLine(rafSeg))
	E.x = interpLine.schnitt(segLine(extSeg))
	R.c, E.c = to3(R.x)[2], to3(E.x)[2]
	R.m = M.m * (M.c - E.c)/(R.c - E.c)
	E.m = M.m - R.m
	return M, R, E
}

func (p Point) inQuad (Ps [4]Point) bool {
	var totalArea float64
	sorted := sortPolPoints(Ps[:])
	n := len(Ps)
	j := n - 1
	for i := 0; i < n; i++ {
		totalArea += area([]Point{sorted[j], sorted[i], p})
		j = i
	}
	quadArea := area(sorted)
	if (totalArea - quadArea) < 0.001 {
		return true
	} else {
		return false
	}
}

func (s Segment) Print() {
	//this loop could be put in a func of its own
	var str, sep string
	p, q := to3(s.P), to3(s.Q)
	for i := 0; i < 3; i++ {
		str += sep + fmt.Sprintf("%.4f", p[i])
		sep = "\t"
	}
	str += sep
	for i := 0; i < 3; i++ {
		str += fmt.Sprintf("%.4f", q[i]) + sep
	}
	fmt.Println(str)
}

func (line Line) Print() {
	fmt.Printf("%.4f\t%.4f\n", line.A, line.B)
}

func (s Stream) Print() {
	fmt.Printf("%.1f\t%.4f\t%.4f\t%.4f\n", s.m, s.x.X, s.x.Y, s.c)
}

func (point Point) Print() {
	//this loop could be put in a func of its own
	var str, sep string
	p := to3(point)
	for i := 0; i < 3; i++ {
		str += sep + fmt.Sprintf("%.4f", p[i])
		sep = "\t"
	}
	fmt.Println(str)
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
		s.x, s.m = to2(fields[0:2]), fields[3]
	}
	f.Close()
	return s
}

func readSegments(filename string) []Segment {
	var segments []Segment
	f, _ := os.Open(filename)
	input := bufio.NewScanner(f)
	var fields []float64
	for input.Scan() {
		fields = splitRow(input.Text(), 6)	
		P, Q := to2(fields[0:2]), to2(fields[3:5])
		segments = append(segments, Segment{P, Q})
	}
	f.Close()
	return segments
}

func splitRow(row string, length int) []float64 {
	sarray := regexp.MustCompile("\t+").Split(row, -1)
	var array []float64
	var newnumber float64
	for i := 0; i < length; i++ {
		newnumber, _ = strconv.ParseFloat(sarray[i], 64)
		array = append(array, newnumber)
	}
	return array
}

func (s Segment) Points() (P Point, Q Point) {
	return s.P, s.Q
}

func (s1 Segment) schnitt(s2 Segment) Point {
	l1 := ppline(s1.P, s1.Q)
	l2 := ppline(s2.P, s2.Q)
	sch := schnitt (l1, l2)
	xvalues := []float64{sch.X, s1.P.X, s1.Q.X, s2.P.X, s2.Q.X}
	yvalues := []float64{sch.Y, s1.P.Y, s1.Q.Y, s2.P.Y, s2.Q.Y}
	sort.Float64s(xvalues)
	sort.Float64s(yvalues)
	//check if x and y sit in the middle of all coordinates
	if sch.X == xvalues[2] && sch.Y == yvalues[2] {
		return sch
	} else {
		//TODO: return error, which will be checked by the calling func
		return Point{0, 0}
	}
}

func (l1 line) schnitt(l2 Line) Point {
	//TODO: break if A1 and A2 are the same
	x := (l2.B - l1.B) / (l1.A - l2.A)
	y := l1.A * x + l1.B 
	return Point{x, y}
}

func (l line) CurSch()

func segLine(s Segment) Line {
	return ppline(s.P, s.Q)	
}

func ppline(p1, p2 Point) Line {
	A := slope(p1, p2)
	B := p1.Y - A*p1.X
	//p1.Y = A*p1.X + B
	return Line{A, B}
}

func slope(p1, p2 Point) float64 {
	return (p2.Y - p1.Y)/(p2.X - p1.X)
}

func area(Ps []Point) float64 {
	var area float64	
	sorted := sortPolPoints(Ps)
	n := len(sorted)
	j := n - 1
	for i := 0; i < n; i++ {
		area += (sorted[j].X + sorted[i].X) * (sorted[j].Y - sorted[i].Y)
		j = i
	}
	return math.Abs(area/2.0)
}

func sortPolPoints (oldpoints []Point) []Point {
	centroid := centroid(oldpoints)
	//var normPoints []Point
	n := len(oldpoints)
	points := make(map[float64]Point)
	for _, oldpoint := range oldpoints { 
		//angle := math.Atan(slope(centroid, oldpoint))
		angle := satan(centroid, oldpoint)
		points[angle] = oldpoint
	}
	angles := make([]float64, 0, n)
	for angle := range points{
		angles = append(angles, angle)
	}
	sort.Float64s(angles)
	var newpoints []Point
	for _, angle := range angles {
		newpoints = append (newpoints, points[angle])
	}
	return newpoints
}

func centroid (points []Point) Point {
	var sumX, sumY float64
	for _, point := range points {
		sumX += point.X
		sumY += point.Y
	}
	avgX := sumX/float64(len(points))
	avgY := sumY/float64(len(points))
	return Point{avgX, avgY}
}

func coords (points []Point) ([]float64, []float64) {
	var xcoords, ycoords []float64
	for i := 0; i < len(points); i++ {
		xcoords = append(xcoords, points[i].X)			
		ycoords = append(ycoords, points[i].Y)			
	}
	return xcoords, ycoords
}

func to3(p Point) [3]float64 {
	return [3]float64{p.X, p.Y, 1-p.X-p.Y}
}

func to2(t []float64) Point {
		return Point{t[0], t[1]}
}

func satan (O, P Point) float64 {
	pi := math.Pi
	y := P.Y - O.Y
	x := P.X - O.X
	atan2angle := math.Atan2(y,x)
	if atan2angle > 0 {
		return atan2angle
	} else {
		return 2*pi + atan2angle
	}
}
