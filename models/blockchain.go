package models

import (
	"io/ioutil"
	"os"
	"encoding/json"
	"crypto/sha256"
	"encoding/hex"
	"strconv"
	"time"
	"fmt"
	"strings"
	"plum/utils"
	"log"
)

type Content struct {
	Affiliation		string
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

func ReplaceChain(newBlocks []Block) {
	blockchain := GetBlockChain()
	if len(newBlocks) > len(blockchain) {
		blockchain = newBlocks
	}
	WriteBlockChain(blockchain)
}

func AppendToBlockChain(block Block) {
	dat, err := ioutil.ReadFile(os.Getenv("CHAIN_FILE"))
	utils.Check(err)

	var bc []Block
	err = json.Unmarshal(dat, &bc)
	utils.Check(err)

	bc = append(bc, block)

	result, err := json.Marshal(bc)
	utils.Check(err)
	err = ioutil.WriteFile("blockchain.json", result, 0644)
	utils.Check(err)
}

func WriteBlockChain(blocks []Block) {
	result, err := json.Marshal(blocks)
	utils.Check(err)
	err = ioutil.WriteFile("blockchain.json", result, 0644)
}

func GetBlockChain() []Block {
	dat, err := ioutil.ReadFile(os.Getenv("CHAIN_FILE"))
	utils.Check(err)

	var bc []Block
	err = json.Unmarshal(dat, &bc)
	utils.Check(err)
	return bc
}

func GetBlockChainString() string {
	dat, err := ioutil.ReadFile(os.Getenv("CHAIN_FILE"))
	utils.Check(err)

	return string(dat)
}
