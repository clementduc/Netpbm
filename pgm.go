package Netpbm

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"
)
 
type PGM struct{
    data [][]uint8
    width, height int
    magicNumber string
    max int
}



func ReadPGM(filename string) (*PGM, error) {
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
    if magicNumber != "P2" && magicNumber != "P5" {
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
        // just for p2
        if magicNumber == "P2" {
            // placer les données de la ligne découper en int
            row := make([]uint8, 0)
            // met les donnée de l'image dans une matrice 
            for _, char := range strings.Fields(line) {
                pixel, err := strconv.Atoi(char)
                if err != nil {
                    fmt.Println("Erreur de conversion en entier :", err)
                }
                // max = 
                if pixel >= 0 && pixel <= max {
                    row = append(row, uint8(pixel))
                } else {
                    fmt.Println("Valeur de pixel invalide :", pixel)
                }
            }
            data = append(data, row)
        }
    }

    // structure qui contient l'image
    return &PGM{
        data:        data,
        width:       width,
        height:      height,
        magicNumber: magicNumber,
        max:         max,
    }, nil
}




func (pgm *PGM) Size() (int, int) {
    return pgm.width, pgm.height
}




func (pgm *PGM) At(x, y int) uint8{
    return pgm.data[y][x]
}


func (pgm *PGM) Set(x, y int, value uint8) {
	pgm.data[y][x] = value
}




func (pgm *PGM) Save(filename string) error {
    file, err := os.Create(filename)
    if err != nil {
        return err
    }
    defer file.Close()

    // Writing header
    fmt.Fprintf(file, "%s\n%d %d\n%d\n", pgm.magicNumber, pgm.width, pgm.height, pgm.max)

    // Writing image's data
    for _, row := range pgm.data {
        for _, pixel := range row {
            fmt.Fprintf(file, "%d ", pixel)
        }
        fmt.Fprintln(file) // new line after wich lines of the image
    }

    return nil
}



func (pgm *PGM) Invert() {
	for i := 0; i < len(pgm.data); i++ {
		for j := 0; j < len(pgm.data[i]); j++ {
			pgm.data[i][j] = uint8(pgm.max) - pgm.data[i][j]
		}
	}
}



func (pgm *PGM) Flip() {
	row := pgm.width
	colums := pgm.height
	for i := 0; i < row; i++ {
		for j := 0; j < colums/2; j++ {
			pgm.data[i][j], pgm.data[i][colums-j-1] = pgm.data[i][colums-j-1], pgm.data[i][j]
		}
	}
}



func (pgm *PGM) Flop() {
	row := len(pgm.data)
	if row == 0 {
		return
	}
	for i := 0; i < row/2; i++ {
		pgm.data[i], pgm.data[row-i-1] = pgm.data[row-i-1], pgm.data[i]
	}
}


func (pgm *PGM) SetMagicNumber(magicNumber string) {
    fmt.Println(pgm.data[0][0], pgm.data[0][1])
}


func (pgm *PGM) SetMaxValue(maxValue uint8) {
	pgm.max = int(maxValue) // convert maxValue to int 
}




func (pgm *PGM) Rotate90CW(){
    //reverse rows and colums
     pgm.width, pgm.height = pgm.height, pgm.height
     //create a new data
     newdata := make([][]uint8, pgm.width)
     for i:= range newdata{
         newdata[i] = make([]uint8, pgm.height)
     }
     //
     for y:=0;y<pgm.height;y++ {
         for x:=0;x<pgm.width;x++ {
             newdata[x][pgm.width-1-y] = pgm.data[y][x]
         }
     }
     pgm.data = newdata
}



func (pgm *PGM) ToPBM() *PBM {
	data := make([][]bool, pgm.height)
	for y := 0; y < pgm.height; y++ {
		data[y] = make([]bool, pgm.width)
		for x := 0; x < pgm.width; x++ {
			data[y][x] = pgm.data[y][x] > uint8(pgm.max/2)
		}
	}
	return &PBM{
		data:        data,
		width:       pgm.width,
		height:      pgm.height,
		magicNumber: "P4",
	}
}
