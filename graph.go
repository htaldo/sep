package sep 
//TODO: change 1000 by width and height
//define functions instead of main
//remove canvas from draw calls
//equilateral based on a new draw function (which does the tf)
//convert functions to methods where applicable
import (
	"os"
	svg "github.com/ajstarks/svgo"
)

func DrawSegment(s Segment, canvas *svg.SVG, style string) {
	x1, x2 := int(1000*s.P.X), int(1000*s.Q.X)
	y1, y2 := int(1000 - s.P.Y*1000), int(1000 - s.Q.Y*1000)
	canvas.Line(x1, y1, x2, y2, style)	
}

func DrawSegments(segs []Segment, canvas *svg.SVG, style string) {
	for _, segment := range segs {
		DrawSegment(segment, canvas, style)
	}
}

func LoadGrid(canvas *svg.SVG) {
	var P, Q Point
	var segment Segment
	var i float64
	for i = 0; i < 1; i += 0.1 {
		//horizontal lines
		P = Point{0, i}
		Q = Point{1, i}
		segment = Segment{P,Q}.Transform()
		DrawSegment(segment, canvas, "stroke:gray")
		//vertical lines
		P = Point{i, 0}
		Q = Point{i, 1}
		segment = Segment{P,Q}.Transform()
		DrawSegment(segment, canvas, "stroke:gray")
		//diagonals
		P = Point{i, 0}
		Q = Point{0, i}
		segment = Segment{P,Q}.Transform()
		DrawSegment(segment, canvas, "stroke:gray")
	}
}

func Transform(oldP Point) {
	var newP Point	
	newP.X = oldP.X + 0.5*oldP.Y
	newP.Y = 0 + oldP.Y
	return newP
}

func Transform(oldS Segment) {
	var newS Segment
	newS.P = oldS.P.Transform()
	newS.Q = oldS.Q.Transform()
	return newS
}
