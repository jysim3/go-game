package models

type Command struct {
	Code string      `json:"command"`
	Data interface{} `json:"data"`
}
