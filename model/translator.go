package model

import (
	"encoding/json"

	yaml "gopkg.in/yaml.v2"
)

// DecodeYamlMessage ...
func DecodeYamlMessage(data []byte, dataType interface{}) error {
	err := yaml.Unmarshal(data, dataType)
	return err
}

// EncodeYamlMessage ..
func EncodeYamlMessage(message interface{}) ([]byte, error) {
	data, err := yaml.Marshal(message)
	if err != nil {
	}
	return data, err
}

// DecodeJsonMessage ...
func DecodeJsonMessage(data []byte, dataType interface{}) error {
	err := json.Unmarshal(data, dataType)
	return err
}

// EncodeJsonMessage ..
func EncodeJsonMessage(message interface{}) ([]byte, error) {
	data, err := json.Marshal(message)
	if err != nil {
	}
	return data, err
}
