package services

import (
	"plum/models"
	"io/ioutil"
	"os"
	"encoding/json"
	"plum/utils"
)

func WriteBlockChain(block models.Block) {
	dat, err := ioutil.ReadFile(os.Getenv("CHAIN_FILE"))
	utils.Check(err)

	var bc []models.Block
	err = json.Unmarshal(dat, &bc)
	utils.Check(err)

	bc = append(bc, block)

	result, err := json.Marshal(bc)
	utils.Check(err)
	err = ioutil.WriteFile("blockchain.json", result, 0644)
	utils.Check(err)
}

func GetBlockChain() []models.Block {
	dat, err := ioutil.ReadFile(os.Getenv("CHAIN_FILE"))
	utils.Check(err)

	var bc []models.Block
	err = json.Unmarshal(dat, &bc)
	utils.Check(err)
	return bc
}

func GetBlockChainString() string {
	dat, err := ioutil.ReadFile(os.Getenv("CHAIN_FILE"))
	utils.Check(err)

	return string(dat)
}
