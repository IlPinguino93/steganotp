package steganotp

import (
	"bytes"
	"errors"
	"fmt"
	"image"
	"image/color"
	"image/png"
	"os"
	"strconv"
)

// FileToBytes reads a file into byte array
func FileToBytes(filename string) ([]byte, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	meta, err := file.Stat()
	if err != nil {
		return nil, err
	}

	size := meta.Size()
	result := make([]byte, size)

	_, err = file.Read(result)
	if err != nil {
		return nil, err
	}

	return result, nil
}

// OpenImage opens the key image
func OpenImage(filename string) (image.Image, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	img, err := png.Decode(file) // Force PNG format. We need lossless, so JPG is not an option anyway, and we need pixels.

	if err != nil {
		return nil, err
	}

	return img, nil
}

func addOrSub(RGB uint32, number uint32) uint32 { // If RGB is smaller than number, add. Otherwise, subtract. Unsigned variables.
	if RGB < number {
		return RGB + number
	}
	return RGB - number
}

func colordiff(lock uint32, key uint32) byte { // Produce the difference between two values in absolutes.
	if lock > key {
		return (byte)(lock - key)
	}
	return (byte)(key - lock)
}

// Size - how many bytes can be writteng into img?
func Size(img image.Image) int {
	return img.Bounds().Max.X * img.Bounds().Max.Y
}

// Encode - write data into img, producing a new image.Image with the concealed information on it.
func Encode(img image.Image, data []byte) (image.Image, error) {
	maxX := img.Bounds().Max.X
	maxY := img.Bounds().Max.Y

	if len(data) > maxX*maxY { // Data too large for the image? (1 Byte = 1 Pixel)
		return nil, errors.New("Image too small for information")
	}

	//println("R - G - B")

	newImage := image.NewRGBA(img.Bounds())
	for y := 0; y < maxY; y++ { // Iterate over all the pixels in the key image.
		for x := 0; x < maxX; x++ {
			oldR, oldG, oldB, oldA := img.At(x, y).RGBA()
			oldR = oldR / 257 // Go's img library uses 0-65536, but PNG only has 0-255 so we'd lose information here. Convert.
			oldG = oldG / 257
			oldB = oldB / 257
			oldA = oldA / 257
			octal := []rune{'0', '0', '0'} // Octal number. If there is no data to write, this will be written.

			if len(data) > (x + y*maxX) { // If there is data to write, convert it to octal, 3 digits.
				octal = []rune(fmt.Sprintf("%03o", data[x+y*maxX]))
			}

			dataR, _ := strconv.Atoi(string(octal[0])) // Three digits to three color values.
			dataG, _ := strconv.Atoi(string(octal[1]))
			dataB, _ := strconv.Atoi(string(octal[2]))

			newR := addOrSub(oldR, uint32(dataR)) // Add or subtract the data values from the colors.
			newG := addOrSub(oldG, uint32(dataG))
			newB := addOrSub(oldB, uint32(dataB))
			newA := oldA // We'll carry over Alpha, if it's there.

			newColor := color.RGBA{uint8(newR), uint8(newG), uint8(newB), uint8(newA)} // The new pixel's color.

			newImage.SetRGBA(x, y, newColor) // Write new pixel to the lock image.

		}
	}
	return newImage, nil
}

// Decode - decode the information encoded in lock using key as key.
func Decode(lock image.Image, key image.Image) (data []byte, err error) {
	var result []byte // Our decoded data.
	// The lock and key image sizes must match. If they don't, cancel.
	if lock.Bounds().Max.X != key.Bounds().Max.X || lock.Bounds().Max.Y != key.Bounds().Max.Y {
		return nil, errors.New("Image sizes mismatch. This key cannot be for this lock")
	}
	// Iterate over each pixel of the lock image.
	for y := 0; y < lock.Bounds().Max.Y; y++ {
		for x := 0; x < lock.Bounds().Max.X; x++ {

			// Get lock and key pixel values
			lockR, lockG, lockB, _ := lock.At(x, y).RGBA()
			keyR, keyG, keyB, _ := key.At(x, y).RGBA()

			diffR := colordiff(lockR/257, keyR/257) // Difference between these, in absolute, is the data value of this
			diffG := colordiff(lockG/257, keyG/257)
			diffB := colordiff(lockB/257, keyB/257)

			num := fmt.Sprintf("%d%d%d", diffR, diffG, diffB) // Join the three digits to make an "octal" string.
			n, err := strconv.ParseInt(num, 8, 16)            // Parse octal string into an int.
			if err != nil {                                   // Parsing went wrong
				return nil, err
			}
			decoded := byte(n) // Cast to byte, because that's what the result is.
			result = append(result, decoded)
		}
	}
	return bytes.TrimRight(result, "\x00"), nil
}
