package crypto

import (
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
	"errors"
	"strconv"
	"strings"
)

// Decrypt function is used to decrypt the string
func Decrypt(securemess string) (decodedmess string, err error) {
	if len(strings.TrimSpace(securemess)) == 0 {
		return "", errors.New("string is empty")
	}
	decodedStr := strings.Split(securemess, "||")
	if len(decodedStr) == 2 {
		ver, err := strconv.Atoi(decodedStr[0])
		if err != nil {
			return "", err
		}
		switch ver {
		case encV1:
			decodedmess, err = decrypt1(decodedStr[1])
			if err != nil {
				return "", err
			}
		default:
			return "", errors.New("invalid encryption")
		}
	}

	return decodedmess, nil
}

func decrypt1(securemess string) (string, error) {
	cipherText, err := base64.URLEncoding.DecodeString(securemess)
	if err != nil {
		return "", err
	}

	block, err := aes.NewCipher(cipherKey)
	if err != nil {
		return "", err
	}

	if len(cipherText) < aes.BlockSize {
		err = errors.New("ciphertext block size is too short")
		return "", err
	}

	//IV needs to be unique, but doesn't have to be secure.
	//It's common to put it at the beginning of the ciphertext.
	iv := cipherText[:aes.BlockSize]
	cipherText = cipherText[aes.BlockSize:]

	stream := cipher.NewCFBDecrypter(block, iv)
	// XORKeyStream can work in-place if the two arguments are the same.
	stream.XORKeyStream(cipherText, cipherText)

	decodedmess := string(cipherText)
	return decodedmess, nil
}
