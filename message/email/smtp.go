package Email

import (
	"crypto/tls"
	"net/smtp"
	"strings"

	"github.com/jordan-wright/email"
)

var SMTP *Smtp

type Smtp struct {
	Host      string
	SendUser  string
	SendPwd   string
	ReplyUser string
	UserName  string
}

func New() *Smtp {
	return SMTP
}

func Init(Host, User, Pwd, Reply string) {
	SMTP = &Smtp{
		Host:      Host,
		SendUser:  User,
		SendPwd:   Pwd,
		ReplyUser: Reply,
	}
}

type SendData interface {
	GetTo() string
	GetSubject() string
	GetBody() string
	GetMailType() string
}

func (E *Smtp) SendEmailWithTLS(Data SendData) error {
	hp := strings.Split(E.Host, ":")
	auth := smtp.PlainAuth("", E.SendUser, E.SendPwd, hp[0])

	e := email.NewEmail()
	e.From = E.SendUser
	e.To = []string{Data.GetTo()}
	e.Subject = Data.GetSubject()
	e.HTML = []byte(Data.GetBody())

	return e.SendWithTLS(E.Host, auth, &tls.Config{ServerName: hp[0]})
}
