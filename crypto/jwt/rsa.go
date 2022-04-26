package jwt

import (
	"crypto/rsa"
	"encoding/json"
	"errors"
	"io/ioutil"
	"time"

	"github.com/golang-jwt/jwt"
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
	GetExpiration() time.Time
	GetNotBefore() time.Time
}

func (J *RsaJwt) CreateRsaJWT(M JMsg) (tokenStr string, err error) {
	exp := M.GetExpiration().Unix()
	nbf := M.GetNotBefore().Unix()
	if exp < 0 || nbf < 0 || exp < nbf {
		err = errors.New("time is invalid")
		return
	}
	token := jwt.NewWithClaims(jwt.SigningMethodRS256, M.GetMapClaims())
	return token.SignedString(J.privateKey)
}

// Parses the Token and converts the result to struct return
// @params rMsg must pass in the struct pointer type
func (J *RsaJwt) ParseToken(tokenStr string, rMsg JMsg) (err error) {
	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, errors.New("verify that the encryption type of the token is incorrect")
		}
		return J.publicKey, nil
	})
	if err != nil {
		return
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !token.Valid || !ok {
		return errors.New("token is invalid or has no corresponding value")
	}
	data, _ := json.Marshal(claims)
	err = json.Unmarshal(data, rMsg)
	if err != nil {
		return
	}
	now := time.Now().Unix()
	exp := rMsg.GetExpiration().Unix()
	nbf := rMsg.GetNotBefore().Unix()
	if exp < 0 || nbf < 0 || exp < now || now < nbf || exp < nbf {
		return errors.New("token is invalid")
	}
	return
}
