package main

import (
	"flag"
	"fmt"
	"image"
	"image/color"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"os"
	"strconv"
	"time"
	"math/rand"
)

type Coordinates struct {
	x int
	y int
	char string
}

func main() {
	// Process flag
	path := flag.String("fpath", "fpath", "help message for flagname")
	flag.Parse()
	if path != nil {
		// Do we have a file to play with
		if _, err := os.Stat(*path); err == nil {
			fmt.Printf("file exists; processing...\n")
			process(*path)
		} else {
			fmt.Println("File doesnt exist")
		}
	}
}

func process(path string) {

	infile, err := os.Open(path)
	if err != nil {
		// We can't open the file so dont go any further
		fmt.Println("Cant open file")
		return
	}
	defer infile.Close()

	img, _, err := image.Decode(infile)

	if err != nil {
		fmt.Println(err)
		return
	}

	// Get the width and height of the image
	bounds := img.Bounds()
	w, h := bounds.Max.X, bounds.Max.Y
	
	// Instantiate slices
	var cords = make([][]string, h)
	for y := 0; y < h; y++ {
		cords[y] = make([]string, w)
	} 

	// Creat channel for the char comunication
	pictureChan := make(chan Coordinates)

	// Allocate chars as seperate threads
	for y := 0; y < h; y++ {
		for x := 0; x < w; x++ {
			go pickCharForPixel(img, x, y, pictureChan)
		}
	}

	// Clear the screen
	fmt.Printf("\033[2J")
	// Reset cursor to top left
	fmt.Printf("\033[0;0H")

	for p := 0; p < h * w; p++ {
		// Recieve the data from the channel
		row := <- pictureChan
		// Bulid up string to move cursor
		com := "\033[" + strconv.Itoa(row.y) + ";" + strconv.Itoa(row.x) + "H"
		// Move cursor
		fmt.Printf(com)
		// Output char
		fmt.Printf(row.char)
	}
}

func pickCharForPixel(img image.Image, x int, y int, pictureChan chan Coordinates) {
	// Set a random time delay so that the chars dont all load at once
	amt := time.Duration(rand.Intn(2000))
    time.Sleep(time.Millisecond * amt)

	var char string

	// Get the color value between 0 - 255
	oldColor := img.At(x, y)
	value := color.GrayModel.Convert(oldColor).(color.Gray).Y

	// Choose different chars based on the pixel brightness
	// Could improve selcetion of this and put in more chars
	if value > 0 && value < 50 {
		char = ":"
	} else if value > 50 && value < 100 {
		char = "+"
	} else if value > 100 && value < 150 {
		char = "&"
	} else if value > 150 {
		char = "@"
	} else {
		char = "#"
	}

	// Send this data back to the channel
	pictureChan <- Coordinates{x : x, y : y, char : char}
}
