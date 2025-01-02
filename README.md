# Subpixel

Replace each horizontal pixal in an image with an RGB sub-pixel. Reduces the width of an image to one third of its original width. Converts an image to grayscale before scaling.

Aspect ratio can be maintained using the `-a` flag, but it reduces the images height to 1/3 its original height when converting to subpixel to match the new width, and stretches the image's height 3x when converting back to full size grayscale or expanded pixel format.

## Why?

I thought it would be interesting

## Usage

```
Usage of subpixel:
  -a    maintain aspect ratio when converting, makes the image 1/3 as tall when shrinking and 3x taller when reversing
  -e    expand each pixel to three full color pixels
  -i string
        input png
  -o string
        output png
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
