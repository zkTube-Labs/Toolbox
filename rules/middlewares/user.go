package middlewares

import (
	"errors"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis"
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

func AuthOperate(op string, r *redis.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		Token := c.GetHeader("opToken")
		if Token == "" {
			responses.ParamErrRep(c, &responses.Responses{
				Msg: "missing operate information",
			})
			c.Abort()
			return
		}
		RJ := jwt.NewJwt()
		Operate := &jwt.AuthMsg{}
		if err := RJ.ParseToken(Token, Operate); err != nil {
			responses.ParamErrRep(c, &responses.Responses{
				Msg: err.Error(),
			})
			c.Abort()
			return
		}
		UUID := c.GetString("UUID")
		err := errors.New("you do not have permission to do this")
		if UUID != Operate.UUID {
			responses.ParamErrRep(c, &responses.Responses{
				Msg: err.Error(),
			})
			c.Abort()
			return
		}
		if op != Operate.Action {
			responses.ParamErrRep(c, &responses.Responses{
				Msg: err.Error(),
			})
			c.Abort()
			return
		}
		if str := r.Get(Operate.Key).Val(); str != "" {
			responses.ParamErrRep(c, &responses.Responses{
				Msg: "repeat the operation",
			})
			c.Abort()
			return
		}
		r.Del(Operate.Key)
	}
}
