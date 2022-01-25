package middlewares

import (
	"github.com/gin-gonic/gin"
	"github.com/zkTube-Labs/Toolbox/crypto/jwt"
	"github.com/zkTube-Labs/Toolbox/message/responses"
)

func AuthUser() gin.HandlerFunc {
	return func(c *gin.Context) {
		Token := c.GetHeader("token")
		RJ := jwt.NewJwt()
		User := &jwt.RJMsg{}
		err := RJ.ParseToken(Token, User)
		if err != nil {
			responses.ParamErrRep(c, &responses.Responses{
				Msg: err.Error(),
			})
			return
		}
		c.Set("UUID", User.UUID)
	}
}
