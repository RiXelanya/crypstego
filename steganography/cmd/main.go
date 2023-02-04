package main

import (
	"bufio"
	"bytes"
	"fmt"
	"image"
	"io"
	"mime/multipart"
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
	var routemap string = "decrypt"

	return c.String(http.StatusOK, string(comm(message, routemap)))

}

func encode(c echo.Context) error {
	var routemap string = "encrypt" // the url of the encryption service
	message := c.FormValue("message")
	msg := comm(message, routemap)
	file, err := c.FormFile("file") // read the file
	if err != nil {
		return err
	}
	src, err := file.Open()
	if err != nil {
		return err
	}
	defer src.Close()
	reader := bufio.NewReader(src)
	img, _, err := image.Decode(reader) // decode the image to its nrgba form
	if err != nil {
		return err
	}
	encodedImg := new(bytes.Buffer)
	err = steganography.Encode(encodedImg, img, msg) // steganographically encode the msg and img into encodedimg
	if err != nil {
		return err
	}
	outFile, err := os.Create("encoded.png") // Creates the encoded image file
	if err != nil {
		return err
	}
	bufio.NewWriter(outFile).Write(encodedImg.Bytes()) // write the encoded image into file
	defer outFile.Close()
	c.Response().After(func() { os.Remove("encoded.png") })

	return c.Attachment("encoded.png", "encoded.png") // return the encoded image file as download

}

func comm(message string, routemap string) []byte {
	var url string = "http://127.0.0.1:1324/"
	url = url + routemap
	var r bytes.Buffer
	w := multipart.NewWriter(&r)
	w.WriteField("message", message)
	w.Close()                                                // create form data to send to the encryption service
	resp, err := http.Post(url, w.FormDataContentType(), &r) // send it to the encryption service
	if err != nil {
		fmt.Println(err)
	}

	msgr := bufio.NewReader(resp.Body) // read the message response
	msg, err := io.ReadAll(msgr)       // store reading of msgr inside the msg bytes array
	if err != nil {
		fmt.Println(err)
	}
	resp.Body.Close() // close the response body
	return msg        // return the msg bytes
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
