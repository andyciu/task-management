package utils

import (
	"encoding/json"

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

func JsonNumberPointToIntPoint(jnum *json.Number) *int {
	if jnum == nil {
		return nil
	}

	tempnum, _ := jnum.Int64()
	result := int(tempnum)

	return &result
}

func QueryCondLikeString(str string) string {
	return "%" + str + "%"
}
