package logic

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/sha256"
	"encoding/base64"
)

var bytes = []byte{35, 46, 57, 24, 85, 35, 24, 74, 87, 35, 88, 98, 66, 32, 14, 05}

const MySecret string = "top secret key"

//Encrypt шифрует идентификатор
func Encrypt(text string) (string, error) {
	key := sha256.Sum256([]byte(MySecret))
	block, err := aes.NewCipher(key[:])
	if err != nil {
		return "", err
	}
	plainText := []byte(text)
	cfb := cipher.NewCFBEncrypter(block, bytes)
	cipherText := make([]byte, len(plainText))
	cfb.XORKeyStream(cipherText, plainText)
	return Encode(cipherText), nil
}

//Decrypt расшифровывает идентификатор
func Decrypt(text string) (string, error) {
	key := sha256.Sum256([]byte(MySecret))
	_, err := base64.StdEncoding.DecodeString(text)
	if err != nil {
		return "", err
	}
	block, err := aes.NewCipher([]byte(key[:]))
	if err != nil {
		return "", err
	}

	cipherText := Decode(text)
	cfb := cipher.NewCFBDecrypter(block, bytes)
	plainText := make([]byte, len(cipherText))
	cfb.XORKeyStream(plainText, cipherText)
	return string(plainText), nil
}

func Encode(b []byte) string {
	return base64.StdEncoding.EncodeToString(b)
}

func Decode(s string) []byte {
	data, err := base64.StdEncoding.DecodeString(s)
	if err != nil {
		panic(err)
	}
	return data
}
