# Subpixel

Replace each horizontal pixal in an image with an RGB sub-pixel. Reduces the width of an image to one third of its original width. Converts an image to grayscale before scaling.

If you want the image to maintain the same aspect ratio after conversion, resize the image to be 33% its original height while maintaining the same width before converting it.

## Why?

I thought it would be interesting

## Usage

```
Usage of subpixel:
  -a    maintain aspect ratio when converting, makes the image 1/3 as tall when shrinking and 3x taller when reversing
  -e    expand each pixel to three full color pixels
  -i string
        input PNG
  -o string
        output PNG
  -r    reverse process
```

## Building

```
go build
```

## Example

### Reduce width to subpixel rendering

```
subpixel -i smiley.png -o smiley-subpixel.png
```

![input](smiley.png) -> ![output](smiley-subpixel.png)

### Convert back to full width from subpixel rendering, grayscale

```
subpixel -r -i smiley-subpixel.png -o smiley-restored.png
```

![input](smiley-subpixel.png) -> ![output](smiley-restored.png)

### Expand subpixel image to full width, maintaining subpixel colors

```
subpixel -e -i smiley-subpixel.png -o smiley-expanded.png
```

![input](smiley-subpixel.png) -> ![output](smiley-expanded.png)
