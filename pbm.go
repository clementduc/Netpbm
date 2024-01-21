package Netpbm

import (
	"bufio"
    "fmt"
    "os"
    "strings"
    "errors"
    "strconv"
)

type PBM struct{
    data [][]bool
    width, height int
    magicNumber string
}



func ReadPBM(filename string) (*PBM, error) {
    // Open the file and verify if it has been opened 
    file, err := os.Open(filename)
    if err != nil {
        return nil, err
    } 
    defer file.Close()

    // Initialise 
    var width, height, max int
    var data [][]uint8

    // Verify the magic number 
    scanner := bufio.NewScanner(file)
    scanner.Scan()
    magicNumber := scanner.Text()
    if magicNumber != "P1" && magicNumber != "P4" {
        return nil, errors.New("type de fichier non pris en charge")
    }

    // 
    for scanner.Scan() {
        line := scanner.Text()
        if !strings.HasPrefix(line, "#") {
            _, err := fmt.Sscanf(line, "%d %d", &width, &height)
            if err == nil {
                break
            } else {
                fmt.Println("Largeur ou hauteur invalide :", err)
            }
        }
    }

    // ASCII in integer 
    scanner.Scan()
    max, err = strconv.Atoi(scanner.Text())
    if err != nil {
        return nil, errors.New("valeur maximale de pixel invalide")
    }

    // run all the lines
    for scanner.Scan() {
        // line by line
        line := scanner.Text()
        if magicNumber == "P1" {
            row := make([]uint8, 0)
            // data in matrix
            for _, char := range strings.Fields(line) {
                pixel, err := strconv.Atoi(char)
                if err != nil {
                    fmt.Println("Erreur de conversion en entier :", err)
                }
                // 
                if pixel >= 0 && pixel <= max {
                    row = append(row, uint8(pixel))
                } else {
                    fmt.Println("Valeur de pixel invalide :", pixel)
                }
            }
			// add on data
            data = append(data, row)
        }
    }

    // structure who contain the image
    return &PBM{
        data:        data,
        width:       width,
        height:      height,
        magicNumber: magicNumber,
    }, nil
}



func (pbm *PBM) Size() (int, int) {
    return pbm.width, pbm.height
}




func (pbm *PBM) At(x, y int) bool{
	if x >= 0 && x < pbm.width && y >= 0 && y < pbm.height {
		return pbm.data[y][x]
	}
	return false
}



func (pbm *PBM) Set(x, y int, value bool) {
	pbm.data[y][x] = value
}




func (pbm *PBM) Save(filename string) error {
    file, err := os.Create(filename)
    if err != nil {
        return err
    }
    defer file.Close()

    writer := bufio.NewWriter(file)

    // Write the magic number
    _, err = fmt.Fprintln(writer, pbm.magicNumber)
    if err != nil {
        return err
    }

    // Write the dimensions
    _, err = fmt.Fprintf(writer, "%d %d\n", pbm.width, pbm.height)
    if err != nil {
        return err
    }

    // Write pixel data based on the magic number
    if pbm.magicNumber == "P1" {
        for i := 0; i < pbm.height; i++ {
            for j := 0; j < pbm.width; j++ {
                if pbm.data[i][j] {
                    _, err = fmt.Fprint(writer, "1 ")
                } else {
                    _, err = fmt.Fprint(writer, "0 ")
                }
                if err != nil {
                    return err
                }
            }
            
            _, err = fmt.Fprintln(writer)
            if err != nil {
                return err
            }
        }
    }
return writer.Flush() // write in the file
}




func (pbm *PBM) Invert() {
    for y := 0; y < pbm.height; y++ {
        for x := 0; x < pbm.width; x++ {
			// reverse 2 pixels
            pbm.data[y][x] = !pbm.data[y][x]
        }
    }
}




// flip the image horizontaly
func (pbm *PBM) Flip() {
	for i := 0; i < pbm.height; i++ {
		for j := 0; j < pbm.width/2; j++ {
			pbm.data[i][j], pbm.data[i][pbm.width-j-1] = pbm.data[i][pbm.width-j-1], pbm.data[i][j]
		}
	}
}


// flop the image verticaly
func (pbm *PBM) Flop() {
	for i := 0; i < pbm.height/2; i++ {
		pbm.data[i], pbm.data[pbm.height-i-1] = pbm.data[pbm.height-i-1], pbm.data[i]
	}
}



// define the magic number
func (pbm *PBM) SetMagicNumber(magicNumber string) {
    fmt.Println(pbm.data[0][0], pbm.data[0][1])
}

