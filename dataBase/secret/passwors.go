package secret

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"github.com/ruraomsk/TLServer/logger"
	"golang.org/x/crypto/bcrypt"
	"io"
)

func GetHasPassword(password string) string {
	b, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(b)
}

var key = []byte("TZPtSIacEJG18IrUrAkTE6luYmnCNKgR")

func CodeString(message string) string {
	plaintext := []byte(message)
	block, err := aes.NewCipher(key)
	if err != nil {
		return ""
	}
	ciphertext := make([]byte, aes.BlockSize+len(plaintext))
	iv := ciphertext[:aes.BlockSize]
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return ""
	}
	steam := cipher.NewCFBEncrypter(block, iv)
	steam.XORKeyStream(ciphertext[aes.BlockSize:], plaintext)
	return base64.URLEncoding.EncodeToString(ciphertext)
}
func DecodeString(message string) string {
	ciphertext, _ := base64.URLEncoding.DecodeString(message)

	block, err := aes.NewCipher(key)
	if err != nil {
		logger.Error.Print(err.Error())
		return ""
	}

	if len(ciphertext) < aes.BlockSize {
		logger.Error.Print("ciphertext too short")
		return ""
	}
	iv := ciphertext[:aes.BlockSize]
	ciphertext = ciphertext[aes.BlockSize:]

	stream := cipher.NewCFBDecrypter(block, iv)

	// XORKeyStream can work in-place if the two arguments are the same.
	stream.XORKeyStream(ciphertext, ciphertext)

	return fmt.Sprintf("%s", ciphertext)
}
