package models

type CommonRes struct {
	Code    string      `json:"code"`
	Message string      `json:"message"`
	Content interface{} `json:"content"`
}
