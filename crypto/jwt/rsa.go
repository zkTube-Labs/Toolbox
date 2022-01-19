package jwt

import (
	"crypto/rsa"
	"encoding/json"
	"errors"
	"io/ioutil"

	jwt "github.com/dgrijalva/jwt-go"
)

var RJ *RsaJwt

type RsaJwt struct {
	publicKey  *rsa.PublicKey
	privateKey *rsa.PrivateKey
}

func NewJwt() *RsaJwt {
	return RJ
}

func InitRsaJWT(pubFile, priFile string) (err error) {
	RJ = &RsaJwt{}
	publicKeyByte, err := ioutil.ReadFile(pubFile)
	if err != nil {
		return
	}
	RJ.publicKey, err = jwt.ParseRSAPublicKeyFromPEM(publicKeyByte)
	if err != nil {
		return
	}
	privateKeyByte, err := ioutil.ReadFile(priFile)
	if err != nil {
		return
	}
	RJ.privateKey, err = jwt.ParseRSAPrivateKeyFromPEM(privateKeyByte)
	if err != nil {
		return
	}
	return
}

func InitIssuer(priFile string) (err error) {
	RJ = &RsaJwt{}
	privateKeyByte, err := ioutil.ReadFile(priFile)
	if err != nil {
		return
	}
	RJ.privateKey, err = jwt.ParseRSAPrivateKeyFromPEM(privateKeyByte)
	if err != nil {
		return
	}
	return
}

func InitUser(pubFile string) (err error) {
	RJ = &RsaJwt{}
	publicKeyByte, err := ioutil.ReadFile(pubFile)
	if err != nil {
		return
	}
	RJ.publicKey, err = jwt.ParseRSAPublicKeyFromPEM(publicKeyByte)
	if err != nil {
		return
	}
	return
}

type JMsg interface {
	GetMapClaims() jwt.MapClaims
}

func (J *RsaJwt) CreateRsaJWT(M JMsg) (tokenStr string, err error) {
	token := jwt.NewWithClaims(jwt.SigningMethodRS256, M.GetMapClaims())
	return token.SignedString(J.privateKey)
}

// Parses the Token and converts the result to struct return
// @params rMsg must pass in the struct pointer type
func (J *RsaJwt) ParseToken(tokenStr string, rMsg interface{}) (err error) {
	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, errors.New("verify that the encryption type of the token is incorrect")
		}
		return J.publicKey, nil
	})
	if err != nil {
		return
	}
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		data, _ := json.Marshal(claims)
		err = json.Unmarshal(data, rMsg)
		return
	}
	err = errors.New("token is invalid or has no corresponding value")
	return
}
