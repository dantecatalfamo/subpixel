package main

import (
	"flag"
	"image"
	"image/color"
	"image/png"
	"log"
	"os"
)

func main() {
	inPath := flag.String("in", "", "input PNG")
	outPath := flag.String("out", "", "output PNG")
	reverse := flag.Bool("r", false, "reverse process")
	expand := flag.Bool("e", false, "expand each pixel to three full color pixels")
	flag.Parse()

	if *inPath == "" {
		log.Fatal("no input file provided")
	}

	if *outPath == "" {
		log.Fatal("no output file provided")
	}

	if *reverse && *expand {
		log.Fatal("cannot use reverse and expand at the same time")
	}

	inFile, err := os.Open(*inPath)
	if err != nil {
		log.Fatalf("failed to open file %s: %s", *inPath, err)
	}

	outFile, err := os.Create(*outPath)
	if err != nil {
		log.Fatalf("failed to create output file: %s: %s", *outPath, err)
	}
	_ = outFile

	inImage, err := png.Decode(inFile)
	if err != nil {
		log.Fatalf("failed to decode png: %s", err)
	}

	var outImage image.Image

	if *reverse {
		outImage = subpixelToFull(inImage)
	} else if *expand {
		outImage = expandPixels(inImage)
	} else {
		outImage = fullToSubpixel(inImage)
	}

	err = png.Encode(outFile, outImage)
	if err != nil {
		log.Fatalf("failed to write png: %s", err)
	}

	log.Print("done")
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
			lum := uint8((0.299*float64(r) + 0.587*float64(g) + 0.114*float64(b)) / 256)
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
			r, g, b, a := inImage.At(x, y).RGBA()
			outImage.Set(x*3, y, color.RGBA{R: uint8(r), A: uint8(a)})
			outImage.Set(x*3+1, y, color.RGBA{G: uint8(g), A: uint8(a)})
			outImage.Set(x*3+2, y, color.RGBA{B: uint8(b), A: uint8(a)})
		}
	}

	return outImage
}
