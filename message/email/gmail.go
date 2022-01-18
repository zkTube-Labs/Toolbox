package Email

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"os"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/gmail/v1"
	"google.golang.org/api/option"
)

var Gmail *GMail

type GMail struct {
	srv         *gmail.Service
	SendUser    string
	SendAddress string
}

func InitGmail(credentials string, tokFile string) (err error) {
	Gmail = &GMail{}
	config, err := google.ConfigFromJSON([]byte(credentials), gmail.GmailReadonlyScope)
	if err != nil {
		return
	}
	tok, err := tokenFromFile(tokFile)
	if err != nil {
		return
	}
	client := config.Client(context.Background(), tok)
	Gmail.srv, err = gmail.NewService(context.Background(), option.WithHTTPClient(client))
	return
}

// Retrieves a token from a local file.
func tokenFromFile(file string) (*oauth2.Token, error) {
	f, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	tok := &oauth2.Token{}
	err = json.NewDecoder(f).Decode(tok)
	return tok, err
}

func NewGmail() *GMail {
	return Gmail
}

type GmailMsg interface {
	GetTo() string
	GetSubject() string
	GetBody() string
}

func (G *GMail) SendEmail(sendmsg GmailMsg) (*gmail.Message, error) {
	header := make(map[string]string)
	header["From"] = G.SendAddress
	header["To"] = sendmsg.GetTo()
	header["Subject"] = sendmsg.GetSubject()
	header["MIME-Version"] = "1.0"
	header["Content-Type"] = "text/html; charset=\"utf-8\""
	header["Content-Transfer-Encoding"] = "base64"

	var msg string
	for k, v := range header {
		msg += fmt.Sprintf("%s: %s\r\n", k, v)
	}
	msg += "\r\n" + sendmsg.GetBody()

	gmsg := gmail.Message{
		Raw: base64.RawURLEncoding.EncodeToString([]byte(msg)),
	}

	return G.srv.Users.Messages.Send(G.SendAddress, &gmsg).Do()
}
