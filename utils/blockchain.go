package utils

import (
	"plum/models"
	"io/ioutil"
	"os"
	"encoding/json"
)

func WriteBlockChain(block models.Block) {
	dat, err := ioutil.ReadFile(os.Getenv("CHAIN_FILE"))
	Check(err)

	var bc []models.Block
	err = json.Unmarshal(dat, &bc)
	Check(err)

	bc = append(bc, block)

	result, err := json.Marshal(bc)
	Check(err)
	err = ioutil.WriteFile("blockchain.json", result, 0644)
	Check(err)
}

func GetBlockChain() []models.Block {
	dat, err := ioutil.ReadFile(os.Getenv("CHAIN_FILE"))
	Check(err)

	var bc []models.Block
	err = json.Unmarshal(dat, &bc)
	Check(err)
	return bc
}

func GetBlockChainString() string {
	dat, err := ioutil.ReadFile(os.Getenv("CHAIN_FILE"))
	Check(err)

	return string(dat)
}
