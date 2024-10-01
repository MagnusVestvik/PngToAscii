package main

import (
	"fmt"
	"github.com/gdamore/tcell/v2"
	"github.com/nfnt/resize"
	"image"
	_ "image/jpeg"
	_ "image/png"
	"os"
	"strings"
)

/*
	var asciiChars = []string{
		"@","%","#","*","+","=","-",";",":","."," ",
	}
*/
var asciiChars = []string{
	"@", "@", "@", "%", "%", "#", "*", "*", "+", "+", "=", "-", "-", ":", ".", " ", " ",
}

func getTerminalSize() (int, int) {
	screen, err := tcell.NewScreen()
	if err != nil {
		panic(err)
	}
	defer screen.Fini()

	screen.Init()
	width, height := screen.Size()
	return width, height
}

func mapIntensityToASCII(intensity int) string {
	asciiIndex := (intensity * len(asciiChars)) / 256 % len(asciiChars)
	return asciiChars[asciiIndex]
}

func extractRgbValues(image image.Image) [][]int {
	bounds := image.Bounds()
	height, width := bounds.Max.Y, bounds.Max.X
	rgbValues := make([][]int, height)

	for i := range rgbValues {
		rgbValues[i] = make([]int, width)
	}

	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			r, g, b, _ := image.At(x, y).RGBA()
			gray := 0.299*float64(r) + 0.587*float64(g) + 0.114*float64(b) // convert to grayscale
			rgbValues[y][x] = int(gray) >> 8                               // scaling down the colors to 255
		}
	}
	return rgbValues
}

func mapRgbToASCII(image [][]int) [][]string {
	height, width := len(image), len(image[0])
	asciiChars := make([][]string, height)

	for x := 0; x < height; x++ {
		asciiChars[x] = make([]string, width)
	}

	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			asciiChars[y][x] = mapIntensityToASCII(image[y][x])
		}
	}
	return asciiChars
}

func drawImage(image [][]string) {
	drawnImage := strings.Builder{}
	height, width := len(image), len(image[0])
	for i := 0; i < height; i++ {
		for j := 0; j < width; j++ {
			drawnImage.WriteString(image[i][j])
		}
		drawnImage.WriteString("\n")
	}
	fmt.Println(drawnImage.String())
	file, err := os.OpenFile("asciArt.txt", os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		fmt.Println("Could create file... ", err)
	}
	defer file.Close()

	data := []byte(drawnImage.String())
	_, err = file.Write(data)
	if err != nil {
		fmt.Println("Somthing went wrong, could not write to file... ", err)
	}

	err = file.Sync()
	if err != nil {
		fmt.Println(err)
	}
}

func resizeImg(img image.Image) image.Image {
	width, height := getTerminalSize()
	newWidth := uint(width)
	newHeight := uint(height)
	resizedImage := resize.Resize(newWidth, newHeight, img, resize.Lanczos3)
	return resizedImage
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Please provide a path to the image file you want to convert")
		return
	}
	pathToImage := os.Args[1]

	file, err := os.Open(pathToImage)
	if err != nil {
		fmt.Println("Woops an error occurred when trying to open the file", err)
		return
	}
	defer file.Close()

	img, _, err := image.Decode(file)
	if err != nil {
		fmt.Println("Whoops an error occurred when trying to decode the image", err)
		return
	}
	img = resizeImg(img)

	rgbValues := extractRgbValues(img)

	asciiImage := mapRgbToASCII(rgbValues)

	drawImage(asciiImage)
}
