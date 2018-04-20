package models

import (
	"crypto/sha256"
	"encoding/hex"
	"time"
)

type Component interface {
	calculateHash() string
}

type Block struct {
	Index     int
	Timestamp string
	BPM       int
	Hash      string
	PrevHash  string
	Content		Content
}

type Content struct {
	Affiliation		string
}

func(block Block) CalculateHash() string {
	record := string(block.Index) + block.Timestamp + string(block.BPM) + block.PrevHash
	h := sha256.New()
	h.Write([]byte(record))
	hashed := h.Sum(nil)
	return hex.EncodeToString(hashed)
}

func (oldBlock Block) GenerateNextBlock(bpm int, content Content) (Block, error) {
	t := time.Now()

	newBlock := Block{oldBlock.Index + 1, t.String(), bpm, "", oldBlock.Hash, content}
	newBlock.Hash = newBlock.CalculateHash()

	return newBlock, nil
}

func (block Block) IsBlockValid(oldBlock Block) bool {
	if(oldBlock.Hash != block.PrevHash) {
		return false
	}

	if(oldBlock.Index != block.Index - 1) {
		return false
	}

	if(block.CalculateHash() != block.Hash) {
		return false
	}

	return true
}