package sep

type Stream struct { 
	m,c float64 //c is the composition of C
	x Point		//x represents the compositions of A and B
}
type Diagram struct { EqLines, Alders []Segment }

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
//	E0.x = schnittExt(ppLine(Rn.x,M.x))	
//	E = append(E, E0)
//	P := ppLine(F.x,E0).SchnittLine(ppLine(Rn.x, S.x))
//	var nextRx, nextEx Point
//	nextRx = Point{0, 0}
//	for i := 0; to3(nextRx)[2] > Rn.c; i++ {
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
	FSline := ppLine(S.x, F.x)
	M.x = Mline.SchnittLine(FSline)
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
	l1, l2 := s1.Line(), s2.Line()
	interpPoint := l1.SchnittLine(l2)
	interpLine := ppLine(interpPoint, M.x)
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
	R.c, E.c = to3(R.x)[2], to3(E.x)[2]
	R.m = M.m * (M.c - E.c)/(R.c - E.c)
	E.m = M.m - R.m
	return M, R, E
}

// end lle

// beg text

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

// end text
