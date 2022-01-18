package responses

import (
	"encoding/json"
	"net/http"

	"github.com/gin-gonic/gin"
)

type Msg interface {
	toJson() []byte
	setCode(int)
	setMsg(string)
}

type Responses struct {
	Code int         `json:"code"`
	Data interface{} `json:"data"`
	Msg  string      `json:"message"`
}

func (R *Responses) setCode(code int) {
	R.Code = code
}

func (R *Responses) setMsg(msg string) {
	R.Msg = msg
}

func (R *Responses) toJson() (msg []byte) {
	msg, _ = json.Marshal(R)
	return
}

func SuccessRep(c *gin.Context, m Msg) {
	m.setCode(SUCCESS)
	m.setMsg("Success")
	c.JSON(http.StatusOK, m)
}

func ParamErrRep(c *gin.Context, m Msg) {
	if m == nil {
		m = &Responses{
			Code: ParameterError,
			Msg:  "Parameter error",
		}
	}
	c.JSON(http.StatusBadRequest, m)
}

func ServerErrRep(c *gin.Context, m Msg) {
	c.JSON(http.StatusInternalServerError, m)
}
