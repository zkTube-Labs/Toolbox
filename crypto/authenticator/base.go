package authenticator

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha1"
	"encoding/base32"
	"encoding/binary"
	"fmt"
	"net/url"
	"strings"
	"time"
)

type GoogleAuth struct {
}

func NewGoogleAuth() *GoogleAuth {
	return &GoogleAuth{}
}

func (G *GoogleAuth) un() int64 {
	return time.Now().UnixNano() / 1000 / 30
}

func (G *GoogleAuth) hmacSha1(key, data []byte) []byte {
	h := hmac.New(sha1.New, key)
	if total := len(data); total > 0 {
		h.Write(data)
	}
	return h.Sum(nil)
}

func (G *GoogleAuth) base32encode(src []byte) string {
	return base32.StdEncoding.EncodeToString(src)
}

func (G *GoogleAuth) base32decode(s string) ([]byte, error) {
	return base32.StdEncoding.DecodeString(s)
}

func (G *GoogleAuth) toBytes(value int64) []byte {
	var result []byte
	mask := int64(0xFF)
	shifts := [8]uint16{56, 48, 40, 32, 24, 16, 8, 0}
	for _, shift := range shifts {
		result = append(result, byte((value>>shift)&mask))
	}
	return result
}

func (G *GoogleAuth) toUint32(bts []byte) uint32 {
	return (uint32(bts[0]) << 24) + (uint32(bts[1]) << 16) +
		(uint32(bts[2]) << 8) + uint32(bts[3])
}

func (G *GoogleAuth) oneTimePassword(key []byte, data []byte) uint32 {
	hash := G.hmacSha1(key, data)
	offset := hash[len(hash)-1] & 0x0F
	hashParts := hash[offset : offset+4]
	hashParts[0] = hashParts[0] & 0x7F
	number := G.toUint32(hashParts)
	return number % 1000000
}

func (G *GoogleAuth) GetSecret() string {
	var buf bytes.Buffer
	_ = binary.Write(&buf, binary.BigEndian, G.un())
	return strings.ToUpper(G.base32encode(G.hmacSha1(buf.Bytes(), nil)))
}

func (G *GoogleAuth) GetCode(secret string) (string, error) {
	secretUpper := strings.ToUpper(secret)
	secretKey, err := G.base32decode(secretUpper)
	if err != nil {
		return "", err
	}
	number := G.oneTimePassword(secretKey, G.toBytes(time.Now().Unix()/30))
	return fmt.Sprintf("%06d", number), nil
}

func (G *GoogleAuth) GetQrcode(user, secret, issuer string) string {
	return fmt.Sprintf("otpauth://totp/%s?secret=%s&issuer=%s", user, secret, issuer)
}

func (G *GoogleAuth) GetQrcodeUrl(user, secret, issuer string) string {
	qrcode := G.GetQrcode(user, secret, issuer)
	width := "200"
	height := "200"
	data := url.Values{}
	data.Set("data", qrcode)
	return "https://api.qrserver.com/v1/create-qr-code/?" + data.Encode() + "&size=" + width + "x" + height + "&ecc=M"
}

func (G *GoogleAuth) VerifyCode(secret, code string) (bool, error) {
	_code, err := G.GetCode(secret)
	fmt.Println(_code, code, err)
	if err != nil {
		return false, err
	}
	return _code == code, nil
}
