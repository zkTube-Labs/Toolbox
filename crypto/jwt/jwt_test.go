package jwt

import (
	"testing"
)

func TestJwt(t *testing.T) {
	// d, _ := time.ParseDuration("1s")
	// InitRsaJWT("public.key", "private.key")
	// m := &RJMsg{
	// 	Expiration: time.Now().Add(d),
	// 	NotBefore:  time.Now(),
	// 	UUID:       helper.GetUuid(),
	// }
	// jstr, err := RJ.CreateRsaJWT(m)
	// time.Sleep(time.Duration(2) * time.Second)
	// if err != nil {
	// 	panic(err)
	// }
	// fmt.Println("str", jstr)

	// rm := &RJMsg{}
	// // err = RJ.ParseToken(jstr, rm)
	// // if err != nil {
	// // 	panic(err)
	// // }
	// fmt.Println("UUID", rm.UUID)
}
