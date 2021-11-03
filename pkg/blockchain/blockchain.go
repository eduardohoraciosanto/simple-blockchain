package blockchain

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"time"

	"github.com/eduardohoraciosanto/simple-blockchain/config"
	"github.com/sirupsen/logrus"
)

type Service interface {
	InsertDataBlock(BlockData) error
	GetBlockchain() Blockchain
}

const (
	blocksToSetDifficulty     = 5
	timeToDetermineDifficulty = 30 * time.Second
)

type blockchain struct {
	log                      *logrus.Entry
	version                  string
	blockchain               Blockchain
	targetLeadingZeroes      int
	successfullyHashedBlocks int
	lastDifficultySet        time.Time
}

func NewBlockchain(log *logrus.Entry, version string) (Service, error) {
	log.Info("Creating new BlockchainInstance")
	tlz, err := config.GetEnvInt(config.STARTING_LEADING_ZEROES)
	if err != nil {
		log.WithError(err).Error("Unable to get starting target leading zeroes from ENV")
		return nil, err
	}
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
		log:                      log,
		version:                  version,
		blockchain:               chain,
		targetLeadingZeroes:      tlz,
		successfullyHashedBlocks: 0,
		lastDifficultySet:        time.Now(),
	}, nil
}

func (b *blockchain) GetBlockchain() Blockchain {
	b.log.Info("Returning full Blockchain")
	return b.blockchain
}

func (b *blockchain) InsertDataBlock(data BlockData) error {
	b.log.Info("Inserting Block on chain")
	chainlen := len(b.blockchain)
	newBlock, err := b.generateBlock(b.blockchain[chainlen-1], data)
	if err != nil {
		b.log.WithError(err).Error("Error generating block")
		return err
	}
	b.log.Info("Block generated")
	if !b.isBlockValid(newBlock, b.blockchain[chainlen-1]) {
		return fmt.Errorf("Block Not Valid")
	}

	//add the block to the ones that were successfully hashed
	b.successfullyHashedBlocks = b.successfullyHashedBlocks + 1
	b.log.WithField("generated_blocks", b.successfullyHashedBlocks).Info("Block is Valid")
	newBlockchain := append(b.blockchain, newBlock)
	b.replaceChain(newBlockchain)
	b.setDifficulty()

	return nil
}

//calculateHash calculates our Block's hash
func (b *blockchain) calculateHash(block Block) (string, error) {
	b.log.Info("Calculating Hash for Block")
	dataBytes, err := json.Marshal(block.Data)
	if err != nil {
		b.log.WithError(err).Error("Unable to Marshal Block Data")
		return "", err
	}
	record := fmt.Sprint(block.Index) + block.Timestamp + string(dataBytes) + block.PrevHash + fmt.Sprint(block.Nuance)
	h := sha256.New()
	h.Write([]byte(record))
	hashed := h.Sum(nil)
	return hex.EncodeToString(hashed), nil
}

//generateBlock will generate a new block using the old one and the new data (BPM) inside it
func (b *blockchain) generateBlock(oldBlock Block, data BlockData) (Block, error) {
	ts := time.Now().String()
	b.log.WithField("timestamp", ts).Info("Generating new Block with timestamp")
	// A block will not be generated until the hash complies with the LEADING_ZEROES value
	newBlock := Block{
		Index:     oldBlock.Index + 1,
		Timestamp: ts,
		Data:      data,
		PrevHash:  oldBlock.Hash,
	}
	var newhash string
	var err error
	nuance := 0
	to := time.Now()
	b.log.Info("Beginning Hash calculation")
	for {
		newBlock.Nuance = nuance
		//we'll calculate until we successfully find a proper hash
		newhash, err = b.calculateHash(newBlock)
		if err != nil {
			return Block{}, err
		}
		if b.hashCompliesWithRule(newhash) {
			break
		}
		nuance += 1
		if nuance < 0 {
			b.log.Warn("All nuance values were used, changing timestamp")
			newBlock.Timestamp = time.Now().String()
		}
	}
	b.log.WithField("operation_took", time.Since(to).String()).WithField("new_hash", newhash).Info("Hash obtained successfully!")
	newBlock.Hash = newhash

	return newBlock, nil
}

//isBlockValid validates a single block based on the old block and the hash calculation
func (b *blockchain) isBlockValid(newBlock, oldBlock Block) bool {
	b.log.Info("Verifying Block")
	if oldBlock.Index+1 != newBlock.Index {
		b.log.Error("New block is not directly next to the old block")
		return false
	}

	if oldBlock.Hash != newBlock.PrevHash {
		b.log.Error("Old block hash does not match with new block prev hash")
		return false
	}
	calcHash, err := b.calculateHash(newBlock)
	if err != nil {
		b.log.WithError(err).Error("hash calculation failed")
		return false
	}
	if calcHash != newBlock.Hash {
		b.log.WithField("calculated_hash", calcHash).
			WithField("new_block_hash", newBlock.Hash).
			Error("calculated hash differs from block hash")
		return false
	}

	return true
}

//replaceChain will be of use if we have two sets of blocks,
// we'll replace our blockchain with the most up to date chain which is the longest
func (b *blockchain) replaceChain(newBlocks []Block) {
	if len(newBlocks) > len(b.blockchain) {
		b.log.Info("Replacing chain as newone has more blocks")
		b.blockchain = newBlocks
	}
}

func (b *blockchain) hashCompliesWithRule(hash string) bool {
	le := b.log.WithField("target_leading_zeroes", b.targetLeadingZeroes).WithField("hash", hash)
	for i := 0; i < b.targetLeadingZeroes; i++ {
		if hash[i] != '0' {
			le.Info("Hash does not comply with target leading zeroes")
			return false
		}
	}
	le.Info("Hash successfully complies with target leading zeroes!")
	return true
}

func (b *blockchain) setDifficulty() {
	if b.successfullyHashedBlocks >= blocksToSetDifficulty {
		b.successfullyHashedBlocks = 0
		//we'll set the new difficulty.
		b.log.WithField("block_limit_took_ms", time.Since(b.lastDifficultySet).Milliseconds()).
			Info("Blocks limit reached, evaluating difficulty adjustment")
		if time.Since(b.lastDifficultySet) >= timeToDetermineDifficulty {
			//it was hard? decrease the difficulty
			if b.targetLeadingZeroes > 1 {
				b.targetLeadingZeroes = b.targetLeadingZeroes - 1
				b.log.WithField("new_target_leading_zeroes", b.targetLeadingZeroes).Info("Difficulty Decreased")
			} else {
				b.log.Info("Difficulty already at minimum")
			}

		} else {
			//it was easy? increase the difficulty
			b.targetLeadingZeroes = b.targetLeadingZeroes + 1
			b.log.WithField("new_target_leading_zeroes", b.targetLeadingZeroes).Info("Difficulty Increased")
		}
		b.lastDifficultySet = time.Now()
	}
}
