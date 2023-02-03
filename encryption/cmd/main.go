package main

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func decrypt(c echo.Context) error {
	message := c.FormValue("message")
	msg, err := hex.DecodeString(message)
	if err != nil {
		fmt.Println(err)
	}
	key, _ := os.ReadFile("key.txt")

	block, err := aes.NewCipher(key)
	if err != nil {
		return err
	}

	aesgcm, err := cipher.NewGCM(block)
	if err != nil {
		return err
	}
	nonceSize := aesgcm.NonceSize()
	if len(msg) < nonceSize {
		return err
	}
	nonce, ciphertext := msg[:nonceSize], msg[nonceSize:]
	plaintext, err := aesgcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return err
	}

	return c.String(http.StatusOK, string(plaintext))

}

func encrypt(c echo.Context) error {
	message := c.FormValue("message")
	msg := []byte(message)
	key, _ := os.ReadFile("key.txt")
	block, err := aes.NewCipher(key)
	if err != nil {
		return err
	}

	aesgcm, err := cipher.NewGCM(block)
	if err != nil {
		return err
	}

	nonce := make([]byte, aesgcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return err
	}

	ciphertext := aesgcm.Seal(nonce, nonce, msg, nil)

	return c.String(http.StatusOK, hex.EncodeToString(ciphertext))

}

func main() {
	e := echo.New()

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.Secure())

	e.Static("/", "")
	e.POST("/encrypt", encrypt)
	e.POST("/decrypt", decrypt)

	e.Logger.Fatal(e.Start(":1324"))
}
