// Algorithm: https://www.cise.ufl.edu/~ungor/delaunay/delaunay/node5.html
// Plotting: Export to list, plot using PyPlot
// Export Lines of the form [(x0, y0), (x1, y1)]

package main

import ("fmt"
		"math/rand"
		"sort"
		"math"
)

type point struct {
	x float64
	y float64
}

// Define the collection
type points []point

// Implementing the interface
//
func (slice points) Len() int {
	return len(slice)
}

func (slice points) Less(i, j int) bool {
	return slice[i].x < slice[j].x
}

func (slice points) Swap(i, j int) {
	slice[i], slice[j] = slice[j], slice[i]
}

// Line between two points.
// Contains a pointer to the next line
//
type line struct {
	p1 point
	p2 point
}

// Three Lines that make up the triangle
//
type triangle struct {
	A point
	B point
	C point
}

// Triangle method
// Determine if one of the Triangle's vertices is the Point p
// Return: True if the Point is equal to one of the vertices
func (t triangle) ContainsPoint(p point) bool {
	return t.A == p || t.B == p || t.C == p
}

// Triangle method
// Determines if a given Point is contained within the circumcircle of the triangle
// A circumcircle is the circle whose circumference contains all 3 vertices of a triangle
// Return: True if point is contained
func (t triangle) CircumcircleContains(p point) bool {

	var ab = math.Pow(t.A.x, 2) + math.Pow(t.A.y, 2)
	var cd = math.Pow(t.B.x, 2) + math.Pow(t.B.y, 2)
	var ef = math.Pow(t.C.x, 2) + math.Pow(t.C.y, 2)

	var circum_x = (ab * (t.C.y - t.B.y) + cd * (t.A.y - t.C.y) + ef * (t.B.y - t.A.y)) / (t.A.x * (t.C.y - t.B.y) + t.B.x * (t.A.y - t.C.y) + t.C.x * (t.B.y - t.A.y)) / 2
	var circum_y = (ab * (t.C.x - t.B.x) + cd * (t.A.x - t.C.x) + ef * (t.B.x - t.A.x)) / (t.A.y * (t.C.x - t.B.x) + t.B.y * (t.A.x - t.C.x) + t.C.y * (t.B.x - t.A.x)) / 2
	var circum_radius = math.Sqrt(math.Pow(t.A.x - circum_x, 2) + math.Pow(t.A.y - circum_y, 2))

	var dist = math.Sqrt(math.Pow(p.x - circum_x, 2) + math.Pow(p.y - circum_y, 2))
	return dist <= circum_radius
}

func generatePoints(numPoints int) points{
	r := rand.New(rand.NewSource(99))
	var pointList = make([]point, numPoints)

	for i:= 0; i < numPoints; i++ {
		pointList[i] = point{r.Float64(), r.Float64()}
	}
	return pointList
}

func (e1 line) isEqual(e2 line) bool {
	return (e1.p1 == e2.p1 && e1.p2 == e2.p2 || e1.p1 == e2.p2 && e1.p2 == e2.p1)
}

// Generate a random triangulation
//
func generateArbitraryTriangulation(pointList points) []line {
	// Sort the list by its x-coordinates
	// 
	sort.Sort(pointList)
	var lineList []line

	// Connect p0->p1->p2, then p1->p2->p3
	//
	for i:=2; i < len(pointList); i++ {
		firstPoint:=pointList[i - 2]
		secondPoint:=pointList[i - 1]
		thirdPoint:=pointList[i]
		
		firstLine := line{firstPoint, secondPoint}
		secondLine := line{firstPoint, thirdPoint}

		lineList = append(lineList, firstLine)
		lineList = append(lineList, secondLine)
	}
	return lineList
}

// TODO Restructure to use slices
//
func DelaunayTriangle(pointList []point, superT triangle) []triangle {
	var tList []triangle
	tList = append(tList, superT)
	fmt.Printf("The length of the triangle list is %d\n", len(tList))


	fmt.Printf("Super Triangle point: %1.3f\n", superT.C.x)
	
	for _, p := range pointList {
		fmt.Printf("Iterator point: %1.3f\n", p.x)
		var lineList []line

		// TODO: I think we do potentially want to iterate through
		// and populate a list of "triangles to remove"
		//
		fmt.Printf("The length of the triangle list is %d\n", len(tList))
		for index, t := range tList {
			fmt.Printf("Index is %d\n", index)
			if t.CircumcircleContains(p) {

				lineList = append(lineList, line{t.A, t.B})
				lineList = append(lineList, line{t.A, t.C})
				lineList = append(lineList, line{t.B, t.C})


				// Remove the triangle using slice witchcraft
				//
				// 
				if index == len(tList) - 1 {
					tList = tList[:len(tList) - 1]
				} else {
					tList = tList[:index+copy(tList[index:], tList[index+1:])]
				}
			}
		}

		for index, l := range lineList {
			left := l
			if index == len(lineList) - 2 {
				break
			}
			right := lineList[index + 1]
			
			if left.isEqual(right) {

				// Delete both left and right if they are the same
				//
				lineList = lineList[:index+copy(lineList[index:], lineList[index+2:])]
			}
		}

		// Make the new lines from this point.
		//
		for _, l := range lineList {
			new_triangle := triangle{l.p1,l.p2, p}
			tList = append(tList, new_triangle)
			
		}
	}

	//Remove any triangles that contain supertriangle points. 
	for index, t := range tList {
		if t.ContainsPoint(superT.A) ||
		   t.ContainsPoint(superT.B) ||
		   t.ContainsPoint(superT.C) {	
	   		
	   	    // Remove the triangle from the slice
		    //
			tList = tList[:index+copy(tList[index:], tList[index+1:])]
		} 
	}


	return tList
}


// Push all non-locally interior edges onto the stack
//

func main() {
	pointList := generatePoints(10)
	//lineList := generateArbitraryTriangulation(pointList)

	superT := &triangle{
		A: point{
			x: 0,
			y: 0,
		},
		B: point{
			x: 1,
			y: 0,
		},
		C: point{
			x: 1,
			y: 1,
		},
		}

	triangleList := DelaunayTriangle(pointList, *superT)

	for _, triangle := range triangleList {
		fmt.Printf("[(%1.3f, %1.3f):(%1.3f, %1.3f)]\n", triangle.A.x, triangle.A.y, triangle.B.x, triangle.B.y)
	}

	for _, element := range pointList {
		fmt.Printf("(%1.3f, %1.3f)\n", element.x, element.y)
	}

	fmt.Printf("*****\n");
	//for _, element := range lineList {
	//	fmt.Printf("[(%1.3f, %1.3f):(%1.3f, %1.3f)]\n", element.p1.x, element.p1.y, element.p2.x, element.p2.y)
	//}
}
