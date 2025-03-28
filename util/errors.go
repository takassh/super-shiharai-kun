package util

import (
	"encoding/json"
	"os"
	"sync"
)

type Response struct {
	Code    string `json:"code"`
	Message string `json:"message"`
	Error   string `json:"error"`
}

// Const constants defined in const.json
type Const struct {
	Errors `json:"errors"`
}

// Errors errors config
type Errors map[string]string

var once sync.Once
var ErrorCodes *Const

func LoadErrors() *Const {
	once.Do(func() {
		var constant *Const
		constFile, err := os.Open("errors.json")
		if err != nil {
			panic(err)
		}
		defer constFile.Close()

		jsonParser := json.NewDecoder(constFile)
		err = jsonParser.Decode(&constant)
		if err != nil {
			panic(err)
		}

		ErrorCodes = constant
	})

	return ErrorCodes
}

// GetError get error by error code, if error key not found, return "undefined"
func (c *Const) GetError(code string, err error) *Response {
	value, found := c.Errors[code]
	if !found {
		return &Response{
			Code:    code,
			Message: "undefined",
			Error:   err.Error(),
		}
	}

	return &Response{
		Code:    code,
		Message: value,
		Error:   err.Error(),
	}
}
