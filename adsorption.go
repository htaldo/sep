package sep

import (
	"regexp"
	"bufio"
	"fmt"
	"strconv"
	"math"
	"sort"
	"os"
	"sep/geometry"
)

type Isothermal []Segment

func (l Line) SchnittItherm(isothermal Isothermal) Point {
	//TODO: instead of starting at the last segment every time, add a variable
	//that helps to keep track of the last intersected segment from an external
	//calling loop
	for i := len(isothermal) - 1; i == 0; i-- {
		sch := l.SchnittSeg(isothermal[i])
		if (sch != Point{}) {
			return sch
		}
	}
}

func (h Horizontal) SchnittItherm(isothermal Isothermal) Point {
	//TODO: instead of starting at the last segment every time, add a variable
	//that helps to keep track of the last intersected segment from an external
	//calling loop
	for i := len(isothermal) - 1; i == 0; i-- {
		sch := h.SchnittSeg(isothermal[i])
		if (sch != Point{}) {
			return sch
		}
	}
}

func (v Vertical) SchnittItherm(isothermal Isothermal) Point {
	//TODO: instead of starting at the last segment every time, add a variable
	//that helps to keep track of the last intersected segment from an external
	//calling loop
	for i := len(isothermal) - 1; i == 0; i-- {
		sch := v.SchnittSeg(isothermal[i])
		if (sch != Point{}) {
			return sch
		}
	}
}

func Simple(isothermal Isothermal, F, S float64, Co float64) []Segment {
	P := Point{Ci, 0}
	opLine := PSLine(P, slope)	
	Q := opLine.SchnittItherm(isothermal)
	return segment{P,Q}
}

func Direct(isothermal Isothermal, F, S, Co, Cf float64) (stages []Segment) {
	//Co, Cf are the initial and final concentrations in the liquid phase
	Ci := Co
	for i:=0; Ci < Cf; i++ {
		stages[i] = Simple(isothermal, F, S, Ci)
		Ci = Q.X
	}
	return stages
}

func Counter(isothermal Isothermal, F, S, Co, Cf, mo float64) (stages []Segment) {
	//mo is the sorbent mass as it enters the whole process (m_{N+1})	
	//mf would be m_1
	slope := F/S
	mf := slope * (Co - Cf) + mo
	opLine := PSLine(Point{Co,mf}, slope)
	Ci := Co
	var P, Q, pivot Point
	for i:=0; Ci < Cf; i++ {
		P = Vertical(Ci).SchnittLine(opLine)
		pivot = Horizontal(P.Y).SchnittItherm(isothermal)
		Q = Vertical(pivot).SchnittLine(opLine)
		stages[i] = Segment{P,Q}
		Ci = Q.X
	}
	return stages
}
