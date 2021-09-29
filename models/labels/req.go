package labels

type LabelCreateReq struct {
	Name string
}

type LabelUpdateReq struct {
	ID   int
	Name string
}

type LabelDeleteReq struct {
	ID int
}
