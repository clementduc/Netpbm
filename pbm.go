package Netpbm

import (
	"bufio"
	"fmt"
	"math"
	"os"
	"sort"
	"strconv"
	"strings"
)
 
type Pixel struct{
    R, G, B uint8
}

type PPM struct{
    data [][]Pixel
    width, height int
    magicNumber string
    max int
}



func ReadPPM(filename string) (*PPM, error) {
	var err error
	var magicNumber string = ""
	var width int
	var height int
	var maxval int
	var counter int
	var headersize int

	// Read content of the file
	file, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	// Split the file content into lines
	splitfile := strings.SplitN(string(file), "\r\n", -1)

	// Iterate through each line of the file
	for i, _ := range splitfile {
		// Identify the PPM magic number
		if strings.Contains(splitfile[i], "P3") {
			magicNumber = "P3"
		} else if strings.Contains(splitfile[i], "P6") {
			magicNumber = "P6"
		}

		// Determine the header size 
		if strings.HasPrefix(splitfile[i], "#") && maxval != 0 {
			headersize = counter
		}

		// Extract width and height from the first line
		splitl := strings.SplitN(splitfile[i], " ", -1)
		if width == 0 && height == 0 && len(splitl) >= 2 {
			width, err = strconv.Atoi(splitl[0])
			height, err = strconv.Atoi(splitl[1])
			headersize = counter
		}

		// Extract maxval from the line
		if maxval == 0 && width != 0 {
			maxval, err = strconv.Atoi(splitfile[i])
			headersize = counter
		}

		counter++
	}

	// Create a slice to store pixel data
	data := make([][]Pixel, height)
	for j := 0; j < height; j++ {
		data[j] = make([]Pixel, width)
	}

	// Process and store pixel data if the counter is greater than the header size
	var splitdata []string
	if counter > headersize {
		for i := 0; i < height; i++ {
			splitdata = strings.SplitN(splitfile[headersize+1+i], " ", -1)
			for j := 0; j < width*3; j += 3 {
				r, _ := strconv.Atoi(splitdata[j])
				g, _ := strconv.Atoi(splitdata[j+1])
				b, _ := strconv.Atoi(splitdata[j+2])
				data[i][j/3] = Pixel{R: uint8(r), G: uint8(g), B: uint8(b)}
			}
		}
	}
	return &PPM{data: data, width: width, height: height, magicNumber: magicNumber, max: int(maxval)}, err
}



func (ppm *PPM) Size() (int, int) {
    return ppm.width, ppm.height
}




func (ppm *PPM) At(x, y int) Pixel {
	return ppm.data[y][x]
}



func (ppm *PPM) Set(x, y int, value Pixel) {
	ppm.data[y][x] = value
}




func (ppm *PPM) Save(filename string) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	writer := bufio.NewWriter(file)

	// Write header
	fmt.Fprintf(writer, "%s\n", ppm.magicNumber)
	fmt.Fprintf(writer, "%d %d\n", ppm.width, ppm.height)
	fmt.Fprintf(writer, "%d\n", ppm.max)

	// Write pixel data
	for i := 0; i < ppm.height; i++ {
		for j := 0; j < ppm.width; j++ {
			fmt.Fprintf(writer, "%d %d %d\n", ppm.data[i][j].R, ppm.data[i][j].G, ppm.data[i][j].B)
		}
	}
	return writer.Flush()
}



func (ppm *PPM) Invert() {
	for i := 0; i < ppm.height; i++ {
		for j := 0; j < ppm.width; j++ {
			ppm.data[i][j].R = uint8(ppm.max) - ppm.data[i][j].R
			ppm.data[i][j].G = uint8(ppm.max) - ppm.data[i][j].G
			ppm.data[i][j].B = uint8(ppm.max) - ppm.data[i][j].B
		}
	}
}



func (ppm *PPM) Flip() {
	for i := 0; i < len(ppm.data); i++ {
		for j := 0; j < len(ppm.data[i])/2; j++ {
			sdata := ppm.data[i][j]
			ppm.data[i][j] = ppm.data[i][len(ppm.data[i])-1-j]
			ppm.data[i][len(ppm.data[i])-1-j] = sdata
		}
	}
}


func (ppm *PPM) Flop() {
	for i := 0; i < len(ppm.data)/2; i++ {
		sdata := ppm.data[i]
		ppm.data[i] = ppm.data[len(ppm.data)-1-i]
		ppm.data[len(ppm.data)-1-i] = sdata
	}
}




func (ppm *PPM) SetMagicNumber(magicNumber string){
    fmt.Println(ppm.data[0][0], ppm.data[0][1])
}



func (ppm *PPM) SetMaxValue(maxValue uint8){
    ppm.max = int(maxValue)
}



func (ppm *PPM) Rotate90CW() {
	// Transpose the image
	for i := 0; i < ppm.height; i++ {
		for j := i + 1; j < ppm.width; j++ {
			ppm.data[i][j], ppm.data[j][i] = ppm.data[j][i], ppm.data[i][j]
		}
	}

	// Reverse each row
	for i := 0; i < ppm.height; i++ {
		for j := 0; j < ppm.width/2; j++ {
			ppm.data[i][j], ppm.data[i][ppm.width-j-1] = ppm.data[i][ppm.width-j-1], ppm.data[i][j]
		}
	}
}





func (ppm *PPM) ToPGM() *PGM {
	pgm := &PGM{
		data:        make([][]uint8, ppm.height),
		width:       ppm.width,
		height:      ppm.height,
		magicNumber: "P2",
		max:         int(ppm.max), 
	}

	for i := 0; i < ppm.height; i++ {
		pgm.data[i] = make([]uint8, ppm.width)
		for j := 0; j < ppm.width; j++ {
			a := uint8(0.299*float64(ppm.data[i][j].R) + 0.587*float64(ppm.data[i][j].G) + 0.114*float64(ppm.data[i][j].B))
			pgm.data[i][j] = a
		}
	}
	return pgm
}




func (ppm *PPM) ToPBM() *PBM {
	pbm := &PBM{
		data:        make([][]bool, ppm.height),
		width:       ppm.width,
		height:      ppm.height,
		magicNumber: "P1",
	}

	for i := 0; i < ppm.height; i++ {
		pbm.data[i] = make([]bool, ppm.width)
		for j := 0; j < ppm.width; j++ {
			a := 0.299*float64(ppm.data[i][j].R) + 0.587*float64(ppm.data[i][j].G) + 0.114*float64(ppm.data[i][j].B)
			pbm.data[i][j] = a > float64(ppm.max)/2
		}
	}
	return pbm
}




type Point struct{
    X, Y int
}



func (ppm *PPM) DrawLine(p1, p2 Point, color Pixel) {
	dx := float64(p2.X - p1.X)
	dy := float64(p2.Y - p1.Y)
	steps := int(math.Max(math.Abs(dx), math.Abs(dy)))

	xi := dx / float64(steps)
	yi := dy / float64(steps)

	x, y := float64(p1.X), float64(p1.Y)

	for i := 0; i <= steps; i++ {
		ppm.Set(int(x), int(y), color)
		x += xi
		y += yi
	}
}


func (ppm *PPM) DrawRectangle(p1 Point, width, height int, color Pixel) {
	p2 := Point{p1.X + width, p1.Y}
	p3 := Point{p1.X, p1.Y + height}
	p4 := Point{p1.X + width, p1.Y + height}

	ppm.DrawLine(p1, p2, color)
	ppm.DrawLine(p2, p4, color)
	ppm.DrawLine(p4, p3, color)
	ppm.DrawLine(p3, p1, color)
}



func (ppm *PPM) DrawFilledRectangle(p1 Point, width, height int, color Pixel) {
    for i := p1.Y; i < p1.Y+height; i++ {
        for j := p1.X; j < p1.X+width; j++ {
            ppm.Set(j, i, color)
        }
    }
}



func (ppm *PPM) DrawCircle(center Point, radius int, color Pixel) {
	for i := 0; i <= radius*2; i++ {
		for j := 0; j <= radius*2; j++ {
			x := i - radius
			y := j - radius

			if x*x+y*y <= radius*radius {
				ppm.Set(center.X+x, center.Y+y, color)
			}
		}
	}
}




func (ppm *PPM) DrawFilledCircle(center Point, radius int, color Pixel) {
    for y := 0; y < ppm.height; y++ {
        for x := 0; x < ppm.width; x++ {
            dx := float64(x - center.X)
            dy := float64(y - center.Y)
            distance := math.Sqrt(dx*dx + dy*dy)

            if int(distance) <= radius {
                ppm.data[y][x] = color
            }
        }
    }
}



func (ppm *PPM) DrawTriangle(p1, p2, p3 Point, color Pixel) {
	ppm.DrawLine(p1, p2, color)
	ppm.DrawLine(p2, p3, color)
	ppm.DrawLine(p3, p1, color)
}


func (ppm *PPM) DrawFilledTriangle(p1, p2, p3 Point, color Pixel) {
    // Sort the points by Y coordinate
    var points [3]Point
    if p1.Y > p2.Y {
        p1, p2 = p2, p1
    }
    if p1.Y > p3.Y {
        p1, p3 = p3, p1
    }
    if p2.Y > p3.Y {
        p2, p3 = p3, p2
    }

    // Calculate slopes for the two edges of the triangle
    slope1 := float64(p2.X-p1.X) / float64(p2.Y-p1.Y)
    slope2 := float64(p3.X-p1.X) / float64(p3.Y-p1.Y)

    // Initialize starting and ending X coordinates for each row
    x1 := float64(p1.X)
    x2 := float64(p1.X)

    // Process each scanline
    for y := p1.Y; y <= p3.Y; y++ {
        // Draw the current row
        for x := int(x1); x <= int(x2); x++ {
            ppm.data[y][x] = color
        }

        // Update starting and ending X coordinates for the next row
        x1 += slope1
        x2 += slope2
    }
}


func (ppm *PPM) DrawPolygon(points []Point, color Pixel) {
	for i := 0; i < len(points)-1; i++ {
		ppm.DrawLine(points[i], points[i+1], color)
	}
	ppm.DrawLine(points[len(points)-1], points[0], color)
}


func (ppm *PPM) DrawFilledPolygon(points []Point, color Pixel) {
    // Check if there are enough points to form a polygon
    if len(points) < 3 {
        return
    }
    // Create a list to store the edges of the polygon
    edges := make([]Point, 0)
    minY, maxY := points[0].Y, points[0].Y
    for _, point := range points {
        if point.Y < minY {
            minY = point.Y
        }
        if point.Y > maxY {
            maxY = point.Y
        }
    }
    // Populate the edges list with the edges of the polygon
    for i, point := range points {
        nextIndex := (i + 1) % len(points)
        nextPoint := points[nextIndex]

        // Add only the edges that are not horizontal
        if point.Y != nextPoint.Y {
            edges = append(edges, point)
        }
    }

    // Sort the edges list based on the X-coordinate
    sort.Slice(edges, func(i, j int) bool {
        return edges[i].X < edges[j].X
    })

    // Create a scanline buffer to keep track of the active edges
    activeEdges := make(map[int]int)

    // Loop through each scanline in the bounding box of the polygon
    for scanlineY := minY; scanlineY <= maxY; scanlineY++ {
        // Update the active edges list by removing edges that end at the current scanline
        for x, endY := range activeEdges {
            if endY <= scanlineY {
                delete(activeEdges, x)
            }
        }

        // Add new edges to the active edges list
        for _, edge := range edges {
            nextIndex := (i + 1) % len(points)
            nextPoint := points[nextIndex]

            if edge.Y <= scanlineY && nextPoint.Y > scanlineY || nextPoint.Y <= scanlineY && edge.Y > scanlineY {
                // Calculate the intersection point with the scanline
                xIntersection := int(float64(edge.X) + float64(scanlineY-edge.Y)/float64(nextPoint.Y-edge.Y)*(float64(nextPoint.X)-float64(edge.X)))

                // Add the intersection point to the active edges list
                activeEdges[xIntersection] = nextPoint.Y
            }
        }

        // Fill the pixels between pairs of active edges
        sortedX := make([]int, 0, len(activeEdges))
        for x := range activeEdges {
            sortedX = append(sortedX, x)
        }
        sort.Ints(sortedX)

        for i := 0; i < len(sortedX); i += 2 {
            startX, endX := sortedX[i], sortedX[i+1]

            // Ensure startX is within bounds
            if startX < 0 {
                startX = 0
            }

            // Ensure endX is within bounds
            if endX >= ppm.width {
                endX = ppm.width - 1
            }

            // Fill the pixels between startX and endX on the current scanline
            for x := startX; x <= endX; x++ {
                ppm.data[scanlineY][x] = color
            }
        }
    }
}



// DrawSierpinskiTriangle dessine un triangle de Sierpinski.
func (ppm *PPM) DrawSierpinskiTriangle(n int, start Point, width int, color Pixel) {
	if n == 0 {
		// Cas de base : dessiner un triangle simple
		ppm.DrawFilledTriangle(
			start,
			Point{start.X + width, start.Y},
			Point{start.X + width/2, start.Y - int(float64(width)*math.Sqrt(3)/2)},
			color,
		)
	} else {
		// Limitez la largeur à une valeur raisonnable
		if width > ppm.width {
			width = ppm.width
		}

		// Calculer les points pour les trois triangles récursifs
		p1 := start
		p2 := Point{start.X + width, start.Y}
		p3 := Point{start.X + width/2, start.Y - int(float64(width)*math.Sqrt(3)/2)}

		// Dessiner le triangle central
		ppm.DrawSierpinskiTriangle(n-1, Point{(p1.X + p2.X) / 2, (p1.Y + p2.Y) / 2}, width/2, color)

		// Dessiner le triangle supérieur
		ppm.DrawSierpinskiTriangle(n-1, p1, width/2, color)

		// Dessiner le triangle supérieur droit
		ppm.DrawSierpinskiTriangle(n-1, Point{(p2.X + p3.X) / 2, (p2.Y + p3.Y) / 2}, width/2, color)

		// Dessiner le triangle supérieur gauche
		ppm.DrawSierpinskiTriangle(n-1, Point{(p1.X + p3.X) / 2, (p1.Y + p3.Y) / 2}, width/2, color)
	}
}


func (ppm *PPM) DrawPerlinNoise(color1 Pixel , color2 Pixel){
}



// KNearestNeighbors resizes the PPM image using the k-nearest neighbors algorithm.
func (ppm *PPM) KNearestNeighbors(newWidth, newHeight, k int) {
    // Check if the new dimensions are valid
    if newWidth <= 0 || newHeight <= 0 {
        fmt.Println("Invalid dimensions for resizing.")
        return
    }

    // Calculate the scaling factors
    xScale := float64(ppm.width) / float64(newWidth)
    yScale := float64(ppm.height) / float64(newHeight)

    // Create a new PPM image with the new dimensions
    resizePPM := PPM{
        data:        make([][]Pixel, newHeight),
        width:       newWidth,
        height:      newHeight,
        magicNumber: ppm.magicNumber,
        max:         ppm.max,
    }

    for i := range resizePPM.data {
        resizePPM.data[i] = make([]Pixel, newWidth)
    }

    // Resize the image using k-nearest neighbors
    for y := 0; y < newHeight; y++ {
        for x := 0; x < newWidth; x++ {
            // Calculate the corresponding point in the original image
            originalX := int(float64(x) * xScale)
            originalY := int(float64(y) * yScale)

            // Collect the k-nearest neighbors
            neighbors := make([]Pixel, 0, k)
            for i := -k / 2; i <= k/2; i++ {
                for j := -k / 2; j <= k/2; j++ {
                    nx := originalX + i
                    ny := originalY + j

                    // Ensure the neighbor is within bounds
                    if nx >= 0 && nx < ppm.width && ny >= 0 && ny < ppm.height {
                        neighbors = append(neighbors, ppm.data[ny][nx])
                    }
                }
            }

            // Calculate the average color of the neighbors
            var avgR, avgG, avgB uint
            for _, neighbor := range neighbors {
                avgR += uint(neighbor.R)
                avgG += uint(neighbor.G)
                avgB += uint(neighbor.B)
            }

            avgR /= uint(len(neighbors))
            avgG /= uint(len(neighbors))
            avgB /= uint(len(neighbors))

            // Set the pixel color in the resized image
            resizePPM.data[y][x] = Pixel{uint8(avgR), uint8(avgG), uint8(avgB)}
        }
    }

    // Update the original PPM image with the resized image
    ppm.data = resizePPM.data
    ppm.width = newWidth
    ppm.height = newHeight
}
