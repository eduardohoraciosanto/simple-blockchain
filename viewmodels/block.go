package viewmodels

type BlockData struct {
	BMP int `json:"bmp"`
}

type InsertDataRequest struct {
	BlockData BlockData `json:"block_data"`
}
