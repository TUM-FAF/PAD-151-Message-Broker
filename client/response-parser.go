package client

import (
	"PAD-151-Message-Broker/model"
)

type Emptyness interface {
}

// ResponseParser ...
type ResponseParser struct {
}

// Parse ..
func (rp *ResponseParser) Parse(response string) (interface{}, error) {
	m := model.ResponseMessageModel{}
	i := helper(m)
	err := model.DecodeYamlMessage([]byte(response), &i)
	return i, err
}

func helper(in model.ResponseMessageModel) interface{} {
	return in
}
