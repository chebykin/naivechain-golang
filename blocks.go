package main

import (
	"bytes"
	"crypto/sha256"
	"fmt"
	"time"
)

var (
	blockchain []*Block
)

type Block struct {
	Index        int64
	PreviousHash []byte
	Timestamp    int64
	Data         []byte
	Hash         []byte
}

func (b *Block) String() string {
	return fmt.Sprintf("#%d, %d, %s %x >>> <<<%x", b.Index, b.Timestamp, b.Data,
		b.PreviousHash, b.Hash)
}

func getGenesisBlock() *Block {
	return &Block{
		Index:        0,
		PreviousHash: []byte("0"),
		Timestamp:    1498381610,
		Data:         []byte("my genesis block!!!"),
		Hash:         []byte("816534932c2b7154836da6afc367695e6337db8a921823784c14378abed4f7d7"),
	}
}

func generateNextBlock(data []byte) *Block {
	previousBlock := latestBlock()
	nextIndex := previousBlock.Index + 1
	nextTimestamp := time.Now().Unix()
	nextHash := calculateHash([]byte(string(nextIndex)), previousBlock.Hash,
		[]byte(string(nextTimestamp)), data)

	return &Block{
		Index:        nextIndex,
		PreviousHash: previousBlock.Hash,
		Timestamp:    nextTimestamp,
		Data:         data,
		Hash:         nextHash,
	}
}

func latestBlock() *Block {
	return blockchain[len(blockchain)-1]
}

func addBlock(block *Block) {
	// TODO: add validation
	blockchain = append(blockchain, block)
}

func calculateHash(elements ...[]byte) []byte {
	h := sha256.New()
	h.Write(bytes.Join(elements, nil))

	return h.Sum(nil)
}
