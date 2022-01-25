package authenticator

import (
	"fmt"
	"testing"
)

func Test_InitAuth(t *testing.T) {
	user := "cody@zktube.io"
	issuer := "zkTube"

	secret := NewGoogleAuth().GetSecret()
	fmt.Println("Secret:", secret, len(secret))

	code, err := NewGoogleAuth().GetCode(secret)
	fmt.Println("Code:", code, err)

	qrCode := NewGoogleAuth().GetQrcode(user, secret, issuer)
	fmt.Println("Qrcode", qrCode)

	qrCodeUrl := NewGoogleAuth().GetQrcodeUrl(user, secret, issuer)
	fmt.Println("QrcodeUrl", qrCodeUrl)
}
