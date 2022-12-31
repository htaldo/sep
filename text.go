package sep

import (
	"regexp"
	"bufio"
	"fmt"
	"strconv"
	"os"
)

func (s Segment) Print() {
	//this loop could be put in a func of its own
	var str, sep string
	p, q := To3(s.P), To3(s.Q)
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

func (point Point) Print() {
	//this loop could be put in a func of its own
	var str, sep string
	p := To3(point)
	for i := 0; i < 3; i++ {
		str += sep + fmt.Sprintf("%.4f", p[i])
		sep = "\t"
	}
	fmt.Println(str)
}

func ReadSegments(filename string) []Segment {
	var segments []Segment
	f, _ := os.Open(filename)
	input := bufio.NewScanner(f)
	var fields []float64
	for input.Scan() {
		fields = splitRow(input.Text(), 6)	
		P, Q := To2(fields[0:2]), To2(fields[3:5])
		segments = append(segments, Segment{P, Q})
	}
	f.Close()
	return segments
}

func SplitRow(row string, length int) []float64 {
	sarray := regexp.MustCompile("\t+").Split(row, -1)
	var array []float64
	var newnumber float64
	for i := 0; i < length; i++ {
		newnumber, _ = strconv.ParseFloat(sarray[i], 64)
		array = append(array, newnumber)
	}
	return array
}

