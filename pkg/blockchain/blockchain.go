package blockchain

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"time"

	"github.com/sirupsen/logrus"
)

type Service interface {
	InsertDataBlock(BlockData) error
	GetBlockchain() Blockchain
}

type blockchain struct {
	log        *logrus.Logger
	version    string
	blockchain Blockchain
}

func NewBlockChain(log *logrus.Logger, version string) Service {
	genesisBlock := Block{
		Index:     0,
		Timestamp: time.Now().String(),
		Data: BlockData{
			BMP: 0,
		},
		Hash:     "",
		PrevHash: "",
	}
	chain := append(Blockchain{}, genesisBlock)

	return &blockchain{
		log:        log,
		version:    version,
		blockchain: chain,
	}
}

func (b *blockchain) GetBlockchain() Blockchain {
	return b.blockchain
}

func (b *blockchain) InsertDataBlock(data BlockData) error {
	chainlen := len(b.blockchain)
	newBlock, err := b.generateBlock(b.blockchain[chainlen-1], data)
	if err != nil {
		return err
	}

	if b.isBlockValid(newBlock, b.blockchain[chainlen-1]) {
		newBlockchain := append(b.blockchain, newBlock)
		b.replaceChain(newBlockchain)
	}

	return nil
}

//calculateHash calculates our Block's hash
func (b *blockchain) calculateHash(block Block) (string, error) {
	dataBytes, err := json.Marshal(block.Data)
	if err != nil {
		return "", err
	}
	record := fmt.Sprint(block.Index) + block.Timestamp + string(dataBytes) + block.PrevHash
	h := sha256.New()
	h.Write([]byte(record))
	hashed := h.Sum(nil)
	return hex.EncodeToString(hashed), nil
}

//generateBlock will generate a new block using the old one and the new data (BPM) inside it
func (b *blockchain) generateBlock(oldBlock Block, data BlockData) (Block, error) {

	newBlock := Block{
		Index:     oldBlock.Index + 1,
		Timestamp: time.Now().String(),
		Data:      data,
		PrevHash:  oldBlock.Hash,
	}

	newhash, err := b.calculateHash(newBlock)
	if err != nil {
		return Block{}, err
	}
	newBlock.Hash = newhash

	return newBlock, nil
}

//isBlockValid validates a single block based on the old block and the hash calculation
func (b *blockchain) isBlockValid(newBlock, oldBlock Block) bool {
	if oldBlock.Index+1 != newBlock.Index {
		return false
	}

	if oldBlock.Hash != newBlock.PrevHash {
		return false
	}
	calcHash, err := b.calculateHash(newBlock)
	if err != nil {
		return false
	}
	if calcHash != newBlock.Hash {
		return false
	}

	return true
}

//replaceChain will be of use if we have two sets of blocks,
// we'll replace our blockchain with the most up to date chain which is the longest
func (b *blockchain) replaceChain(newBlocks []Block) {
	if len(newBlocks) > len(b.blockchain) {
		b.blockchain = newBlocks
	}
}
