package labels

import "encoding/json"

type LabelCreateReq struct {
	Name string `json:"name" binding:"required"`
}

type LabelUpdateReq struct {
	ID   json.Number `json:"id" binding:"required"`
	Name string      `json:"name" binding:"required"`
}

type LabelDeleteReq struct {
	ID json.Number `json:"id" binding:"required"`
}
