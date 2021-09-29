package labels

type LabelCreateReq struct {
	Name string `json:"name" binding:"required"`
}

type LabelUpdateReq struct {
	ID   int    `json:"id" binding:"required"`
	Name string `json:"name" binding:"required"`
}

type LabelDeleteReq struct {
	ID int `json:"id" binding:"required"`
}
