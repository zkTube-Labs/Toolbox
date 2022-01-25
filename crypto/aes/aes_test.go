package aes

import (
	"crypto/aes"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"testing"
)

func getKey() []byte {
	str := "BpLnfgDsc2WD8F2qNfHK5a84jjJkwzDk"
	fmt.Println("salt", str)
	byt := []byte(str)
	str = hex.EncodeToString(byt)
	key, err := hex.DecodeString(str)
	if err != nil {
		panic(err)
	}
	return key
}

func TestCipher(t *testing.T) {
	key := getKey()

	plaintext := []byte("0xB74693f2DAbdb84570755E536e016d3CBDEB0810")

	c := make([]byte, aes.BlockSize+len(plaintext))
	iv := c[:aes.BlockSize]

	ciphertext, err := AesEncrypt(plaintext, key, iv)
	if err != nil {
		panic(err)
	}
	fmt.Println("ciphertext", ciphertext)
	fmt.Println("base64 ciphertext", base64.StdEncoding.EncodeToString(ciphertext))

	plaintext, err = AesDecrypt(ciphertext, key, iv)
	if err != nil {
		panic(err)
	}

	fmt.Println("decrypt", string(plaintext))
}

func TestEcb(t *testing.T) {
	key := getKey()

	plaintext := []byte("0xB74693f2DAbdb84570755E536e016d3CBDEB0810")
	ciphertext := EcbEncrypt(plaintext, key)
	fmt.Println("ciphertext", base64.StdEncoding.EncodeToString(ciphertext))

	plaintext = EcbDecrypt(ciphertext, key)
	fmt.Println("plaintext", string(plaintext))
}
