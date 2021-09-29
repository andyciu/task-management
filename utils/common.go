package utils

import (
	"github.com/pc01pc013/task-management/enums/response"
	"github.com/pc01pc013/task-management/models"
)

func MakeResponseResultSuccess(content interface{}) models.CommonRes {
	return models.CommonRes{
		Code:    string(response.Success),
		Message: "",
		Content: content,
	}
}

func MakeResponseResult(code response.CommonResCode, message string, content interface{}) models.CommonRes {
	return models.CommonRes{
		Code:    string(code),
		Message: message,
		Content: content,
	}
}
