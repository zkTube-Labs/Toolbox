package middlewares

import (
	"github.com/gin-gonic/gin"
	"github.com/zkTube-Labs/Toolbox/crypto/jwt"
	"github.com/zkTube-Labs/Toolbox/message/responses"
)

func AuthUser() gin.HandlerFunc {
	return func(c *gin.Context) {
		Token := c.GetHeader("token")
		if Token == "" {
			responses.ParamErrRep(c, &responses.Responses{
				Msg: "missing login information",
			})
			c.Abort()
			return
		}
		RJ := jwt.NewJwt()
		User := &jwt.RJMsg{}
		err := RJ.ParseToken(Token, User)
		if err != nil {
			responses.ParamErrRep(c, &responses.Responses{
				Msg: err.Error(),
			})
			c.Abort()
			return
		}
		c.Set("UUID", User.UUID)
	}
}

func AuthOperate() gin.HandlerFunc {
	return func(c *gin.Context) {
		Token := c.GetHeader("operate")
		if Token == "" {
			responses.ParamErrRep(c, &responses.Responses{
				Msg: "missing operate information",
			})
			c.Abort()
			return
		}
		RJ := jwt.NewJwt()
		Operate := &jwt.AuthMsg{}
		err := RJ.ParseToken(Token, Operate)
		if err != nil {
			responses.ParamErrRep(c, &responses.Responses{
				Msg: err.Error(),
			})
			c.Abort()
			return
		}
		UUID := c.GetString("UUID")
		if UUID != Operate.UUID {
			responses.ParamErrRep(c, &responses.Responses{
				Msg: "the user does not have permission to operate",
			})
			c.Abort()
			return
		}
		c.Set("OperateKey", Operate.Key)
		c.Set("Operate", Operate.Action)
	}
}
