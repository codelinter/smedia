package main

import (
	"bytes"
	"image"
	"image/jpeg"
	"io"
	"io/ioutil"
	"log"
	"os"

	"github.com/muesli/smartcrop"
)

func analyzerAndCrop(original string) (image.Image, error) {
	f, err := os.Open(original)
	if err != nil {
		//fmt.Printf("\nCould not find image %v\n", original)
		//os.Exit(1)
		return nil, err
	}

	ratio := getImageDimRatio(original)

	if ratio > 0.8 && ratio < 1.91 {
		//fmt.Println(getImageDimRatio(original))
		return nil, nil
	}
	var x int
	if ratio < 0.8 {
		x = 4
	} else {
		x = 9
	}
	img, _, err := image.Decode(f)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	analyzer := smartcrop.NewAnalyzer()
	topCrop, err := analyzer.FindBestCrop(img, x, 5)
	if err != nil {
		log.Fatal(err)
	}

	type SubImager interface {
		SubImage(r image.Rectangle) image.Image
	}
	croppedimg := img.(SubImager).SubImage(topCrop)
	return croppedimg, nil
}

func makeReaderImage(img image.Image) (io.ReadCloser, error) {
	var buf *bytes.Buffer
	buf = new(bytes.Buffer)
	err := jpeg.Encode(buf, img, nil)
	if err != nil {
		return nil, err
	}
	return ioutil.NopCloser(bytes.NewReader(buf.Bytes())), nil
}

func getCroppedReader(file string) (io.ReadCloser, error) {
	img, err := analyzerAndCrop(file)
	if err != nil {
		//fmt.Println("Analy", err)
		return nil, err
	}
	if img == nil {
		//fmt.Println("Nil Image no cropping")
		return nil, nil
	}
	rc, err := makeReaderImage(img)
	if err != nil {
		//fmt.Println("Make Reader", err)
		return nil, err
	}
	return rc, nil
}
func getImageDimRatio(imagePath string) float64 {
	file, err := os.Open(imagePath)
	if err != nil {
		return 0
		//fmt.Fprintf(os.Stderr, "%v\n", err)
	}
	defer file.Close()
	image, _, err := image.DecodeConfig(file)
	if err != nil {
		//fmt.Fprintf(os.Stderr, "Err %v\n", err)
		return 0
	}
	return float64(image.Width) / float64(image.Height)
}
