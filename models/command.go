package models

import "encoding/json"

type Command struct {
	Code string      `json:"command"`
	Data interface{} `json:"data"`
}

func (c Command) Json() []byte {
	r, _ := json.Marshal(c)
	return r
}
