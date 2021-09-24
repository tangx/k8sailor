package httpresponse

import (
	"github.com/gin-gonic/gin"
)

type Response struct {
	Code  int         `json:"code"`
	Data  interface{} `json:"data"`
	Error string      `json:"error"`
}

func OK(c *gin.Context, data interface{}) {
	Common(c, 0, data, nil)
}

func Error(c *gin.Context, code int, err error) {
	Common(c, code, nil, err)
}

func Common(c *gin.Context, code int, data interface{}, err error) {
	_err := ""
	if err != nil {
		_err = err.Error()
	}

	// 强制设置
	if code == 200 {
		code = 0
	}

	resp := Response{
		Code:  code,
		Data:  data,
		Error: _err,
	}

	c.JSON(code, resp)
}
