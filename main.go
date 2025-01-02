package main

import (
	"flag"
	"fmt"
	"image"
	"image/color"
	"image/png"
	"os"

	"golang.org/x/image/draw"
)

func main() {
	inPath := flag.String("i", "", "input PNG")
	outPath := flag.String("o", "", "output PNG")
	reverse := flag.Bool("r", false, "reverse process")
	expand := flag.Bool("e", false, "expand each pixel to three full color pixels")
	aspect := flag.Bool("a", false, "maintain aspect ratio when converting, makes the image 1/3 as tall when shrinking and 3x taller when reversing")
	flag.Parse()

	if *inPath == "" {
		fmt.Fprintln(os.Stderr, "no input file provided")
		os.Exit(1)
	}

	if *outPath == "" {
		fmt.Fprintln(os.Stderr, "no output file provided")
		os.Exit(1)
	}

	if *reverse && *expand {
		fmt.Fprintln(os.Stderr, "cannot use reverse and expand at the same time")
		os.Exit(1)
	}

	inFile, err := os.Open(*inPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to open file %s: %s\n", *inPath, err)
		os.Exit(1)
	}

	inImage, err := png.Decode(inFile)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to decode png: %s\n", err)
		os.Exit(1)
	}

	if *aspect && !*reverse && !*expand {
		scaledImage := image.NewRGBA(image.Rect(0, 0, inImage.Bounds().Dx(), inImage.Bounds().Dy()/3))
		draw.BiLinear.Scale(scaledImage, scaledImage.Bounds(), inImage, inImage.Bounds(), draw.Over, nil)
		inImage = scaledImage
	}

	outFile, err := os.Create(*outPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to create output file: %s: %s\n", *outPath, err)
		os.Exit(1)
	}
	_ = outFile

	var outImage image.Image

	if *reverse {
		outImage = subpixelToFull(inImage)
	} else if *expand {
		outImage = expandPixels(inImage)
	} else {
		outImage = fullToSubpixel(inImage)
	}

	if *aspect && (*reverse || *expand) {
		scaledImage := image.NewRGBA(image.Rect(0, 0, outImage.Bounds().Dx(), outImage.Bounds().Dy()*3))
		draw.BiLinear.Scale(scaledImage, scaledImage.Bounds(), outImage, outImage.Bounds(), draw.Over, nil)
		outImage = scaledImage
	}

	if err := png.Encode(outFile, outImage); err != nil {
		fmt.Fprintf(os.Stderr, "failed to write png: %s\n", err)
		os.Exit(1)
	}
}

func fullToSubpixel(inImage image.Image) image.Image {
	width := inImage.Bounds().Dx() / 3
	if inImage.Bounds().Dx()%3 != 0 {
		width += 1
	}

	outImage := image.NewRGBA(image.Rect(0, 0, width, inImage.Bounds().Dy()))
	inImgHeight := inImage.Bounds().Dy()
	inImgWidth := inImage.Bounds().Dx()
	for y := 0; y < inImgHeight; y += 1 {
		for x := 0; x < inImgWidth; x += 1 {
			outPx := x / 3
			outSubPx := x % 3
			inPx := inImage.At(x, y)
			r, g, b, _ := inPx.RGBA()
			// https://en.wikipedia.org/wiki/Luma_(video)#Rec._601_luma_versus_Rec._709_luma_coefficients
			lum := uint8((0.2126*float64(r) + 0.7152*float64(g) + 0.0722*float64(b)) / 256)
			outR, outG, outB, _ := outImage.At(outPx, y).RGBA()
			outPxColor := color.RGBA{R: uint8(outR), G: uint8(outG), B: uint8(outB), A: 255}
			switch outSubPx {
			case 0:
				outPxColor.R = lum
			case 1:
				outPxColor.G = lum
			case 2:
				outPxColor.B = lum
			}
			outImage.Set(outPx, y, outPxColor)
		}
	}

	return outImage
}

func subpixelToFull(inImage image.Image) image.Image {
	width := inImage.Bounds().Dx() * 3

	outImage := image.NewRGBA(image.Rect(0, 0, width, inImage.Bounds().Dy()))
	inImgHeight := inImage.Bounds().Dy()
	inImgWidth := inImage.Bounds().Dx()
	for y := 0; y < inImgHeight; y += 1 {
		for x := 0; x < inImgWidth; x += 1 {
			r, g, b, _ := inImage.At(x, y).RGBA()
			outImage.Set(x*3, y, color.Gray{Y: uint8(r)})
			outImage.Set(x*3+1, y, color.Gray{Y: uint8(g)})
			outImage.Set(x*3+2, y, color.Gray{Y: uint8(b)})
		}
	}

	return outImage
}

func expandPixels(inImage image.Image) image.Image {
	width := inImage.Bounds().Dx() * 3

	outImage := image.NewRGBA(image.Rect(0, 0, width, inImage.Bounds().Dy()))
	inImgHeight := inImage.Bounds().Dy()
	inImgWidth := inImage.Bounds().Dx()
	for y := 0; y < inImgHeight; y += 1 {
		for x := 0; x < inImgWidth; x += 1 {
			r, g, b, _ := inImage.At(x, y).RGBA()
			outImage.Set(x*3, y, color.RGBA{R: uint8(r), A: 255})
			outImage.Set(x*3+1, y, color.RGBA{G: uint8(g), A: 255})
			outImage.Set(x*3+2, y, color.RGBA{B: uint8(b), A: 255})
		}
	}

	return outImage
}
