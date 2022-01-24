package aes

import (
	"crypto/aes"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"testing"

	"github.com/zkTube-Labs/Toolbox/helper"
)

func TestCipher(t *testing.T) {
	str := helper.RandomString(32)
	fmt.Println(str)
	byt := []byte(str)
	str = hex.EncodeToString(byt)
	fmt.Println(str)
	key, err := hex.DecodeString(str)
	if err != nil {
		panic(err)
	}
	fmt.Println("key", key)
	plaintext := []byte("hello ming")

	c := make([]byte, aes.BlockSize+len(plaintext))
	iv := c[:aes.BlockSize]

	ciphertext, err := AesEncrypt(plaintext, key, iv)
	if err != nil {
		panic(err)
	}

	fmt.Println(base64.StdEncoding.EncodeToString(ciphertext))

	plaintext, err = AesDecrypt(ciphertext, key, iv)
	if err != nil {
		panic(err)
	}

	fmt.Println(string(plaintext))
}
