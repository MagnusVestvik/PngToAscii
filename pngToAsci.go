package main

import (
	"fmt"
	"image"
	_ "image/png"
	_ "image/jpeg"
	"os"
	"strings"
	"log"
	"github.com/nfnt/resize"
)
/*
var asciiChars = []string{
	"@","%","#","*","+","=","-",";",":","."," ",
}
*/
var asciiChars = []string{
	"@", "@", "@", "%", "%", "#", "*", "*", "+", "+", "=", "-", "-", ":", ".", " ", " ",
}
func mapIntensityToASCII(intensity int) string {
	asciiIndex := (intensity * len(asciiChars)) / 256 % len(asciiChars)
	return asciiChars[asciiIndex]
}

func extractRgbValues(image image.Image) [][]int {
	bounds := image.Bounds()
	width, height := bounds.Max.X, bounds.Max.Y
	rgbValues := make([][]int, height)

	for i := range rgbValues {
		rgbValues[i] = make([]int, width)
	}

	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			r, g, b, _ := image.At(x, y).RGBA()
			gray := 0.299*float64(r) + 0.587*float64(g) + 0.114*float64(b) // convert to grayscale
			rgbValues[y][x] = int(gray)>>8 // scaling down the colors to 255
		}
	}
	return rgbValues
}

func mapRgbToASCII(image [][]int) [][]string {
	height, width := len(image), len(image[0])
	asciiChars := make([][]string, len(image))

	for i := 0; i < height; i++ {
		asciiChars[i] = make([]string, width)
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
	height := len(image)
	width := len(image[0])
	for i := 0; i < height; i++ {
		for j := 0; j < width; j++ {
			drawnImage.WriteString(image[i][j])
		}
		drawnImage.WriteString("\n")
	}
	fmt.Println(drawnImage.String())
	file, err := os.OpenFile("asciArt.txt", os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	data := []byte(drawnImage.String())
	_, err = file.Write(data)
	if err != nil {
		log.Fatal(err)
	}

	err = file.Sync()
	if err != nil {
		log.Fatal(err)
	}
}

func resizeImg(img image.Image, width uint) image.Image{
	oldWidth := img.Bounds().Dx()
	oldHeight := img.Bounds().Dy()
	newHeight := uint((float64(width)/float64(oldWidth))*float64(oldHeight))
	img = resize.Resize(width, newHeight, img, resize.Lanczos3)
	return img
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
	
	//const width uint = 430
	//img = resizeImg(img, width)
	
	rgbValues := extractRgbValues(img)

	asciiImage := mapRgbToASCII(rgbValues)

	drawImage(asciiImage)
}
