package controller

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func AbortWithCode(c *gin.Context, code int) {
	AbortWithErrorResponse(c, code, http.StatusText(code))
}

func AbortWithErrorResponse(c *gin.Context, code int, err string) {
	response := ErrorResponse{
		Code:  code,
		Text:  http.StatusText(code),
		Error: err,
	}
	if accept := c.NegotiateFormat(OfferedJsonXmlYaml...); accept == "" {
		c.AbortWithStatusJSON(code, response)
	} else {
		c.Negotiate(code, gin.Negotiate{Offered: OfferedJsonXmlYaml, Data: response})
	}
}

type ErrorResponse struct {
	Code  int    `json:"code" xml:"code" yaml:"code"`
	Text  string `json:"text" xml:"text" yaml:"text"`
	Error string `json:"error" xml:"error" yaml:"error"`
}

type OfferedOption func() []string

//goland:noinspection GoUnusedGlobalVariable
var (
	OfferedJson        = func() []string { return []string{"application/json", "application/json; charset=utf-8"} }
	OfferedXml         = func() []string { return []string{"application/xml", "application/xml; charset=utf-8"} }
	OfferedYaml        = func() []string { return []string{"application/yaml", "application/yaml; charset=utf-8"} }
	OfferedProtobuf    = func() []string { return []string{"application/x-protobuf"} }
	OfferedJsonXmlYaml = Offered(OfferedJson, OfferedXml, OfferedYaml)
)

func Offered(options ...OfferedOption) []string {
	offered := make([]string, 0, 6)
	for _, option := range options {
		offered = append(offered, option()...)
	}
	return offered
}
