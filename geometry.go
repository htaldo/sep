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

const sq2 = 1.41421356237

func PPLine(p1, p2 Point) Line {
	A := Slope(p1, p2)
	B := p1.Y - A*p1.X
	//this satisfies p1.Y = A*p1.X + B
	return Line{A, B}
}

func Slope(p1, p2 Point) float64 {
	return (p2.Y - p1.Y)/(p2.X - p1.X)
}

func PSLine(p Point, slope float64) Line {
	A := slope
	B := p.Y - m*p.X
	return Line{A, B}
}

func Area(Ps []Point) float64 {
	var Area float64	
	sorted := SortPolPoints(Ps)
	n := len(sorted)
	j := n - 1
	for i := 0; i < n; i++ {
		Area += (sorted[j].X + sorted[i].X) * (sorted[j].Y - sorted[i].Y)
		j = i
	}
	return math.Abs(Area/2.0)
}

func (s Segment) Points() (P Point, Q Point) {
	return s.P, s.Q
}

func (s Segment) Line() Line {
	return PPLine(s.P, s.Q)	
}

func (l1 Line) SchnittLine(l2 Line) Point {
	//TODO: break if A1 and A2 are the same
	x := (l2.B - l1.B) / (l1.A - l2.A)
	y := l1.A * x + l1.B 
	return Point{x, y}
}

func (v Vertical) SchnittLine(l line) Point {
	var x float64 = v  
	y := l.A * x + l.B
	return Point{x, y}
}

func (h Horizontal) SchnittLine(l line) Point {
	var y float64 = h 
	x := (y - l.B)/l.A
	return Point{x, y}
}

//TODO: generalize l, h and v
//maybe a type assertion would help (check if type is vertical
//ie if m is infinite

func (l Line) SchnittSeg(s Segment) Point {
	sch = l.SchnittLine(s.line())
		if sch.TwopBound(s.P, s.Q) == true {
			return sch
		}
}

func (h Horizontal) SchnittSeg(s Segment) Point {
	sch = h.SchnittLine(s.Line())
		if sch.TwopBound(s.P, s.Q) == true {
			return sch
		}
}

func (v Vertical) SchnittSeg(s Segment) Point {
	sch = v.SchnittLine(s.Line())
		if sch.TwopBound(s.P, s.Q) == true {
			return sch
		}
}

func (s1 Segment) SchnittSeg(s2 Segment) Point {
	l1 := PPLine(s1.P, s1.Q)
	l2 := PPLine(s2.P, s2.Q)
	sch := l1.SchnittLine(l2)
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

func SortPolPoints (oldpoints []Point) []Point {
	Centroid := Centroid(oldpoints)
	//var normPoints []Point
	n := len(oldpoints)
	points := make(map[float64]Point)
	for _, oldpoint := range oldpoints { 
		//angle := math.Atan(Slope(Centroid, oldpoint))
		angle := Satan(Centroid, oldpoint)
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

func (p Point) InQuad (Ps [4]Point) bool {
	var totalArea float64
	sorted := SortPolPoints(Ps[:])
	n := len(Ps)
	j := n - 1
	for i := 0; i < n; i++ {
		totalArea += Area([]Point{sorted[j], sorted[i], p})
		j = i
	}
	quadArea := Area(sorted)
	if (totalArea - quadArea) < 0.001 {
		return true
	} else {
		return false
	}
}

func Centroid (points []Point) Point {
	var sumX, sumY float64
	for _, point := range points {
		sumX += point.X
		sumY += point.Y
	}
	avgX := sumX/float64(len(points))
	avgY := sumY/float64(len(points))
	return Point{avgX, avgY}
}

func Coords (points []Point) ([]float64, []float64) {
	var xcoords, ycoords []float64
	for i := 0; i < len(points); i++ {
		xcoords = append(xcoords, points[i].X)			
		ycoords = append(ycoords, points[i].Y)			
	}
	return xcoords, ycoords
}

func To3(p Point) [3]float64 {
	return [3]float64{p.X, p.Y, 1-p.X-p.Y}
}

func To2(t []float64) Point {
		return Point{t[0], t[1]}
}

func Satan (O, P Point) float64 {
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

func (o Point) TwopBound(P, Q Point) bool {
	//break if two points are the same
	//write between()
	xcoords := [2]float64{P.X, Q.X}
	ycoords := [2]float64{P.Y, Q.Y}
	sortedx := sort.Ints(xcoords)
	sortedy := sort.Ints(ycoords)
	if o.X >= sortedx[0] && o.X <= sortedx[1] {
		if o.Y >= sortedy[0] && o.Y <= sortedy[1] {
			return true
		}
	}
}
