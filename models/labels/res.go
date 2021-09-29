package labels

type LabelListRes struct {
	ID   int    `json:"id" binding:"required"`
	Name string `json:"name" binding:"required"`
}
