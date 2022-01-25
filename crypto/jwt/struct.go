package jwt

import (
	"encoding/json"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
)

type RJMsg struct {
	IssuanceAt time.Time   `json:"iat"`
	NotBefore  time.Time   `json:"nbf"`
	Expiration time.Time   `json:"exp"`
	Issuer     string      `json:"iss"`
	Subject    string      `json:"sub"`
	Role       string      `json:"role"`
	UUID       string      `json:"uuid"`
	UserInfo   interface{} `json:"user_info"`
}

func (m *RJMsg) GetMapClaims() (result jwt.MapClaims) {
	result = jwt.MapClaims{}
	json_str, _ := json.Marshal(m)
	_ = json.Unmarshal([]byte(json_str), &result)
	return
}

func (m *RJMsg) GetExpiration() time.Time {
	return m.Expiration
}

func (m *RJMsg) GetNotBefore() time.Time {
	return m.NotBefore
}

type AuthMsg struct {
	IssuanceAt time.Time `json:"iat"`
	NotBefore  time.Time `json:"nbf"`
	Expiration time.Time `json:"exp"`
	Issuer     string    `json:"iss"`
	Subject    string    `json:"sub"`
	UUID       string    `json:"uuid"`
	Action     string    `json:"action"`
	Key        string    `json:"key"`
}

func (m *AuthMsg) GetMapClaims() (result jwt.MapClaims) {
	result = jwt.MapClaims{}
	json_str, _ := json.Marshal(m)
	_ = json.Unmarshal([]byte(json_str), &result)
	return
}

func (m *AuthMsg) GetExpiration() time.Time {
	return m.Expiration
}

func (m *AuthMsg) GetNotBefore() time.Time {
	return m.NotBefore
}
