package sep

//F and S are just mass flow rates, not Streams

type Isothermal []Segment

func (l Line) SchnittItherm(isothermal Isothermal) Point {
	//TODO: instead of starting at the last segment every time, add a variable
	//that helps to keep track of the last intersected segment from an external
	//calling loop
	var sch Point
	for i := len(isothermal) - 1; i >= 0; i-- {
		sch = l.SchnittSeg(isothermal[i])
		if (sch != Point{}) {
			break
		}
	}
	return sch
}

func (h Horizontal) SchnittItherm(isothermal Isothermal) Point {
	//TODO: instead of starting at the last segment every time, add a variable
	//that helps to keep track of the last intersected segment from an external
	//calling loop
	var sch Point
	for i := len(isothermal) - 1; i >= 0; i-- {
		sch = h.SchnittSeg(isothermal[i])
		if (sch != Point{}) {
			break
		}
	}
	return sch
}

func (v Vertical) SchnittItherm(isothermal Isothermal) Point {
	//TODO: instead of starting at the last segment every time, add a variable
	//that helps to keep track of the last intersected segment from an external
	//calling loop
	var sch Point
	for i := len(isothermal) - 1; i >= 0; i-- {
		sch = v.SchnittSeg(isothermal[i])
		if (sch != Point{}) {
			break
		}
	}
	return sch
}

func AdsSimple(isothermal Isothermal, F, S float64, Co float64) Segment {
	P := Point{Co, 0}
	slope := -F/S
	opLine := PSLine(P, slope)	
	Q := opLine.SchnittItherm(isothermal)
	return Segment{P,Q}
}

func AdsDirect(isothermal Isothermal, F, S, Co, Cf float64) (stages []Segment) {
	//Co, Cf are the initial and final concentrations in the liquid phase
	Ci := Co
	for i:=0; Ci < Cf; i++ {
		stages[i] = AdsSimple(isothermal, F, S, Ci)
		Ci = stages[i].Q.X
	}
	return stages
}

func AdsCounter(isothermal Isothermal, F, S, Co, Cf, mo float64) (stages []Segment) {
	//mo is the sorbent mass as it enters the whole process (m_{N+1})	
	//mf would be m_1
	slope := F/S
	mf := slope * (Co - Cf) + mo
	opLine := PSLine(Point{Co,mf}, slope)
	Ci := Co
	var P, Q, pivot Point
	for i:=0; Ci > Cf; i++ {
		P = Vertical(Ci).SchnittLine(opLine)
		pivot = Horizontal(P.Y).SchnittItherm(isothermal)
		Q = Vertical(pivot.X).SchnittLine(opLine)
		stages = append(stages, Segment{P,Q})
		Ci = Q.X
	}
	return stages
}
