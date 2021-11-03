package blockchain

type Block struct {
	Index     int
	Timestamp string
	Data      BlockData
	Hash      string
	PrevHash  string
}

type Blockchain []Block

type BlockData struct {
	BMP int
}
