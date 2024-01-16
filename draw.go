package main

import (
	"fmt"
	"math"
	"os"
	"errors"
)

//Color struct
type Color struct {
	R, G, B int
}

//Point struct 
type Point struct {
	x, y int
}

//Display struct 
type Display struct {
	maxX, maxY int
	matrix     [][]Color
}


//screen interface calling all the intializing functions
type screen interface {
	initialize(maxX, maxY int)
	getMaxXY() (int, int)
	drawPixel(x, y int, c Color) error
	getPixel(x, y int) (Color, error)
	clearScreen()
	screenShot(f string) error
}


//makes the screen and then clears it using the clear function
func (d *Display) initialize(maxX, maxY int) {
	d.maxX = maxX
	d.maxY = maxY
	d.matrix = make([][]Color, maxY)
	for i := range d.matrix {
		d.matrix[i] = make([]Color, maxX)
	}
	d.clearScreen()
}

//gets the pixel to see if its in the bounds of the display
func (d *Display) getPixel(x, y int) (Color, error) {
	if x < 0 || x >= d.maxX || y < 0 || y >= d.maxY {
		return Color{}, errors.New("geometry out of bounds")
	}
	return d.matrix[y][x], nil
}

//makes screen just fully white
func (d *Display) clearScreen() {
	for i := range d.matrix {
		for j := range d.matrix[i] {
			d.matrix[i][j] = Color{255, 255, 255}
		}
	}
}

//function to get the screenshot of the drawing and make it into a pmm file
func (d *Display) screenShot(f string) error {
	file, err := os.Create(f + ".ppm")
	if err != nil {
		return err
	}
	defer file.Close()

	fmt.Fprintf(file, "P3\n%d %d\n255\n", d.maxX, d.maxY)

	for _, row := range d.matrix {
		for _, pixel := range row {
			fmt.Fprintf(file, "%d %d %d ", pixel.R, pixel.G, pixel.B)
		}
		fmt.Fprintln(file)
	}

	return nil
}

//thgis function draws the pixels and returns a error if its outside the desired bounds and will not print it
func (d *Display) drawPixel(x, y int, c Color) error {
	if x < 0 || x >= d.maxX || y < 0 || y >= d.maxY {
		return errors.New("geometry out of bounds")
	}
	d.matrix[y][x] = c
	return nil
}

//function to get the max X and Y values
func (d *Display) getMaxXY() (int, int) {
	return d.maxX, d.maxY
}

//color map to assign the name of the color to its desired RBG value using the Color class
var (
	red    = Color{255, 0, 0}
	green  = Color{0, 255, 0}
	blue   = Color{0, 0, 255}
	yellow = Color{255, 255, 0}
	orange = Color{255, 165, 0}
	purple = Color{128, 0, 128}
	brown  = Color{165, 42, 42}
	black  = Color{0, 0, 0}
	white  = Color{255, 255, 255}
)

//ERRORS for output
var outOfBoundsErr = errors.New("geometry out of bounds")
var colorUnknownErr = errors.New("color unknown")

//Structs for Triangles
type Triangle struct {
	pt0, pt1, pt2 Point
	c             interface{}
}

//Structs for Rectangles
type Rectangle struct {
	ll, ur Point
	c      interface{}
}

//Structs for Circles
type Circle struct {
	cp Point
	r  int
	c  interface{}
}

//geometry interface
type geometry interface {
	draw(scn screen) error
	shape() string
}


//  https://gabrielgambetta.com/computer-graphics-from-scratch/07-filled-triangles.html
func interpolate (l0, d0, l1, d1 int) (values []int) {
	a := float64(d1 - d0) / float64(l1 - l0)
	d  := float64(d0)

	count := l1-l0+1
	for ; count>0; count-- {
		values = append(values, int(d))
		d = d+a
	}
	return
}

// function to check if the color is within valid RGB range
func colorUnknown(c Color) bool {
	return c.R < 0 || c.G < 0 || c.B < 0 || c.R > 255 || c.G > 255 || c.B > 255
}

// function to check if the point is outside the screen
func outOfBounds(pt Point, scn screen) bool {
	maxX, maxY := scn.getMaxXY()
	return pt.x < 0 || pt.x >= maxX || pt.y < 0 || pt.y >= maxY
}


//https://gabrielgambetta.com/computer-graphics-from-scratch/
//function to draw the circle using the source above
// func insideCircle(center, tile Point, radius int) bool {
// 	dx := center.x - tile.x
// 	dy := center.y - tile.y
// 	distanceSquared := dx*dx + dy*dy
// 	return 4*distanceSquared <= radius*radius
// }

//https://gabrielgambetta.com/computer-graphics-from-scratch/
//function to draw the circle using the source above
func (circ Circle) draw(scn screen) error {

		top := circ.cp.y - circ.r
		bottom := circ.cp.y + circ.r

	if outOfBounds(circ.cp,scn){ //checks if in bounds
		return outOfBoundsErr
	}

		for y := top; y <= bottom; y++ {
				dy := y - circ.cp.y
				dx := int(math.Sqrt(float64(circ.r*circ.r - dy*dy)))

				left := circ.cp.x - dx
				right := circ.cp.x + dx

				for x := left; x <= right; x++ {
						color, ok := circ.c.(Color) //set color 

						if !ok { //check if color is valid 
							return colorUnknownErr
						}
						if err := scn.drawPixel(x, y, color); err != nil { //draw if not error
								fmt.Println("drawPixel error:", err)
								return outOfBoundsErr
						}
				}
		}
		return nil
}

//function to draw Rectangles
func (rect Rectangle) draw(scn screen) error {
	color, ok := rect.c.(Color) //set color

	if !ok { //check if color is in the map
		return colorUnknownErr
	}
	if outOfBounds(rect.ll,scn) || outOfBounds(rect.ur,scn){ //check if in bounds or not
		return outOfBoundsErr
	}

		for y := rect.ll.y; y <= rect.ur.y; y++ { //draw box
				for x := rect.ll.x; x <= rect.ur.x; x++ {
						if err := scn.drawPixel(x, y, color); err != nil {
								fmt.Println("drawPixel error:", err)
								return outOfBoundsErr
						}
				}
		}
		return nil
}


//  https://gabrielgambetta.com/computer-graphics-from-scratch/07-filled-triangles.html
func (tri Triangle) draw(scn screen) (err error) {
	color := tri.c.(Color) //set color
	if outOfBounds(tri.pt0,scn) || outOfBounds(tri.pt1,scn)  || outOfBounds(tri.pt2,scn){
		return outOfBoundsErr
	}
	if colorUnknown(color) {
		return colorUnknownErr
	}

	y0 := tri.pt0.y
	y1 := tri.pt1.y
	y2 := tri.pt2.y

	// Sort the points so that y0 <= y1 <= y2
	if y1 < y0 { tri.pt1, tri.pt0 = tri.pt0, tri.pt1 }
	if y2 < y0 { tri.pt2, tri.pt0 = tri.pt0, tri.pt2 }
	if y2 < y1 { tri.pt2, tri.pt1 = tri.pt1, tri.pt2 }

	x0,y0,x1,y1,x2,y2 := tri.pt0.x, tri.pt0.y, tri.pt1.x, tri.pt1.y, tri.pt2.x, tri.pt2.y

	x01 := interpolate(y0, x0, y1, x1)
	x12 := interpolate(y1, x1, y2, x2)
	x02 := interpolate(y0, x0, y2, x2)

	// Concatenate the short sides

	x012 := append(x01[:len(x01)-1],  x12...)

	// Determine which is left and which is right
	var x_left, x_right []int
	m := len(x012) / 2
	if x02[m] < x012[m] {
		x_left = x02
		x_right = x012
	} else {
		x_left = x012
		x_right = x02
	}

	// Draw the horizontal segments
	for y := y0; y<= y2; y++  {
		for x := x_left[y - y0]; x <=x_right[y - y0]; x++ {
			scn.drawPixel(x, y, color)
		}
	}
	return
}

// display 
// TODO: you must implement the struct for this variable, and the interface it implements (screen)
var display Display

