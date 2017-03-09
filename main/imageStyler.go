package main

import (
	"fmt"
	"os"
	"image"
	"github.com/nfnt/resize"
	"image/jpeg"
	"strconv"
	"strings"
	"mime/multipart"
)

func SaveOriginalImage(originalImage multipart.File, originalImageName string)  {
	err := os.Mkdir("pictures/original", 0700)
	if err != nil {
		fmt.Println("Couldn't create a directory or directory alrady exists")
	}
	file, err := os.Create("pictures/original/" + originalImageName + ".jpg")
	defer file.Close()
	if err != nil {
		fmt.Println("Couldn't create a file")
	}
	_, err = file.Write([]byte(originalImage))
	if err != nil {
		fmt.Println("Couldn't write to file")
	}
}

func RestyleImage(imageName, neededImageSize string) multipart.File {
	err := os.Mkdir("pictures/styled/"+neededImageSize, 0700)
	if err != nil {
		fmt.Println("Couldn't create a directory or directory alrady exists")
	}
	neededImage := searchImage(imageName, neededImageSize)
	return neededImage

}

func searchImage(fileName, neededImageSize string) *os.File {
	file, err := os.Open("pictures/styled/" + neededImageSize + "/" + fileName + ".jpg")
	if err != nil || file == nil {
		fmt.Println("File not exists, crete new one")
		return createImage(fileName, neededImageSize)
	}
	return file
}

func createImage(fileName, neededImageSize string) *os.File {
	file, err := os.Create("pictures/styled/" + neededImageSize + "/" + fileName + ".jpg")
	defer file.Close()
	if err != nil {
		fmt.Println("Couldn't create a file")
	} else {
		originalImage, err := os.Open("pictures/original/" + fileName + ".jpg")
		if err != nil {
			fmt.Println("Original file isn't exists")
		} else {
			width, height := parseSize(neededImageSize)
			decodedImage, _, err := image.Decode(originalImage)
			if err != nil {
				fmt.Println("Couldn't decode a file")
			}
			decodedImage = resize.Resize(uint(width), uint(height), decodedImage, resize.Lanczos3)
			jpeg.Encode(file, decodedImage, &jpeg.Options{100})
			return file
		}
		return nil
	}
	return nil
}

func parseSize(size string) (uint, uint) {
	width, err := strconv.ParseUint(strings.Split(size, "x")[0], 10, 32)
	if err != nil {
		fmt.Println("Couldn't parse a width")
	}
	height, err := strconv.ParseUint(strings.Split(size, "x")[1], 10, 32)
	if err != nil {
		fmt.Println("Couldn't parse a height")
	}
	return uint(width), uint(height)
}
