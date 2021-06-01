package dataBase

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"github.com/ruraomsk/TLServer/logger"
	"golang.org/x/crypto/bcrypt"
	"io"
	"strings"
)

func GetHasPassword(password string) string {
	b, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(b)
}

var key = []byte("TZPtSIacEJG18IrUrAkTE6luYmnCNKgR")

func CodeString(message string) string {
	plaintext := []byte(message)
	for len(plaintext)%aes.BlockSize != 0 {
		plaintext = append(plaintext, ' ')
	}
	block, err := aes.NewCipher(key)
	if err != nil {
		logger.Error.Printf("AES %s", err.Error())
		return ""
	}
	chipText := make([]byte, aes.BlockSize+len(plaintext))
	iv := chipText[:aes.BlockSize]
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		logger.Error.Printf("rand %s", err.Error())
		return ""
	}
	//fmt.Printf("iv %v\n",iv)
	cbc := cipher.NewCBCEncrypter(block, iv)
	//if err != nil {
	//	logger.Error.Printf("GCM %s",err.Error())
	//	return ""
	//}
	cbc.CryptBlocks(chipText[aes.BlockSize:], plaintext)
	//outbuff=append(outbuff,ciphertext...)

	//fmt.Printf("cipher out %v\n",outbuff)

	return base64.StdEncoding.EncodeToString(chipText)
}
func DecodeString(message string) string {
	inbuffer, _ := base64.StdEncoding.DecodeString(message)
	//fmt.Printf("cipher  in %v\n",inbuffer)

	block, err := aes.NewCipher(key)
	if err != nil {
		logger.Error.Print(err.Error())
		return ""
	}

	if len(inbuffer) < aes.BlockSize {
		logger.Error.Print("ciphertext too short")
		return ""
	}
	iv := inbuffer[:aes.BlockSize]
	cbc := cipher.NewCBCDecrypter(block, iv)
	inbuffer = inbuffer[aes.BlockSize:]
	if len(inbuffer)%aes.BlockSize != 0 {
		fmt.Printf(" len %d not % %d", len(inbuffer), aes.BlockSize)
		return ""
	}
	cbc.CryptBlocks(inbuffer, inbuffer)
	return strings.TrimRight(string(inbuffer), " ")
}
