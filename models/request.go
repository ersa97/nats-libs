package models

type Request struct {
	Route string      `json:"route"`
	Param interface{} `json:"param"`
	Data  interface{} `json:"data"`
}
