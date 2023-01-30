package main

import (
	"bufio"
	"bytes"
	"image"
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
	outFile, err := os.Create("message.txt")
	if err != nil {
		return err
	}
	os.WriteFile("message.txt", []byte(message), 0777)
	defer outFile.Close()
	c.Response().After(func() { os.Remove("message.txt") })
	return c.Attachment("message.txt", "message.txt")
	// return c.String(http.StatusOK, message)

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
	err = steganography.Encode(encodedImg, img, msg)
	if err != nil {
		return err
	}
	outFile, err := os.Create("encoded.png")
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
