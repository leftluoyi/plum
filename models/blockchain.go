package models

import (
	"crypto/sha256"
	"encoding/hex"
	"time"
	"strings"
	"os"
	"strconv"
	"fmt"
	"encoding/json"
	"log"
	"plum/utils"
)

type Component interface {
	calculateHash() string
}

type Block struct {
	Index     	int
	Timestamp 	string
	BPM       	int
	Hash      	string
	PrevHash  	string
	Difficulty 	int
	Nonce	  	string
	Content		Content
}

type Content struct {
	Affiliation		string
}

func(block Block) CalculateHash() string {
	content, err := json.Marshal(block.Content)
	utils.Check(err)
	if err != nil {
		log.Fatal(err)
	}
	record := string(block.Index) + block.Timestamp + string(block.BPM) + block.PrevHash + block.Nonce + string(content)
	h := sha256.New()
	h.Write([]byte(record))
	hashed := h.Sum(nil)
	return hex.EncodeToString(hashed)
}

func (oldBlock Block) GenerateNextBlock(bpm int, content Content) (Block, error) {
	difficulty, _ := strconv.Atoi(os.Getenv("DIFFICULTY"))
	t := time.Now()

	newBlock := Block{oldBlock.Index + 1, t.String(), bpm, "", oldBlock.Hash, difficulty, "", content}
	for i := 0; ; i++ {
		hex := fmt.Sprintf("%x", i)
		newBlock.Nonce = hex
		if !isHashValid(newBlock.CalculateHash(), newBlock.Difficulty) {
			fmt.Println(newBlock.CalculateHash(), " do more work!")
			time.Sleep(time.Second)
			continue
		} else {
			fmt.Println(newBlock.CalculateHash(), " work done!")
			newBlock.Hash = newBlock.CalculateHash()
			break
		}
	}

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

func isHashValid(hash string, difficulty int) bool {
	prefix := strings.Repeat("0", difficulty)
	return strings.HasPrefix(hash, prefix)
}