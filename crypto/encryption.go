package crypto

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"fmt"
	"io"
	"strings"
)

var cipherKey = []byte("0123456789012345")

const (
	encV1 = 1
)

//Encrypt function is used to encrypt the string
func Encrypt(message string) (encmess string, err error) {
	if len(strings.TrimSpace(message)) == 0 {
		return "", errors.New("string is empty")
	}
	plainText := []byte(message)

	block, err := aes.NewCipher(cipherKey)
	if err != nil {
		return "", err
	}

	//IV needs to be unique, but doesn't have to be secure.
	//It's common to put it at the beginning of the ciphertext.
	cipherText := make([]byte, aes.BlockSize+len(plainText))
	iv := cipherText[:aes.BlockSize]
	if _, err = io.ReadFull(rand.Reader, iv); err != nil {
		return "", err
	}

	stream := cipher.NewCFBEncrypter(block, iv)
	stream.XORKeyStream(cipherText[aes.BlockSize:], plainText)

	//returns to base64 encoded string
	encmess = base64.URLEncoding.EncodeToString(cipherText)
	finalEnc := fmt.Sprintf("%d%s%s", encV1, "||", encmess)
	return finalEnc, nil
}
