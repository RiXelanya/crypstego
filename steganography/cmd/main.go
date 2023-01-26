package main

import (
	"bufio"
	"bytes"
	"fmt"
	"image"
	"net/http"
	"os"
	"github.com/auyer/steganography"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func decode(c echo.Context) error {
	file, err := c.FormFile("file")
	if err != nil {
		return err
	}
	src, err := file.Open()
	if err != nil {
		return err
	}
	defer src.Close()
	reader := bufio.NewReader(src)
	img, _, err := image.Decode(reader)
	if err != nil {
		return err
	}

	sizeOfMessage := steganography.GetMessageSizeFromImage(img) 

	msg := steganography.Decode(sizeOfMessage, img) 
	message := string(msg[:])

	return c.String(http.StatusOK, message)

}

func encode(c echo.Context) error {
	message := c.FormValue("message")
	msg := []byte(message)
	file, err := c.FormFile("file")
	if err != nil {
		return err
	}
	src, err := file.Open()
	if err != nil {
		return err
	}
	defer src.Close()
	reader := bufio.NewReader(src)
	img, _, err := image.Decode(reader)
	if err != nil {
		return err
	}
	encodedImg := new(bytes.Buffer)
	err = steganography.Encode(encodedImg, img, msg) // Calls library and Encodes the message into a new buffer
	if err != nil {
		return err
	}
	outFile, err := os.Create("encoded.png") // Creates file to write the message into
	if err != nil {
		return err
	}
	bufio.NewWriter(outFile).Write(encodedImg.Bytes())
	defer outFile.Close()
	c.Response().After(func() { os.Remove("encoded.png") })

	return c.Attachment("encoded.png", "encoded.png")

}

func main() {
	e := echo.New()

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.Secure())

	e.Static("/", "")
	e.POST("/decode", decode)
	e.POST("/encode", encode)

	e.Logger.Fatal(e.Start(":1323"))
}
