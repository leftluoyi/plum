package main

import (
	"log"
	"plum/models"

	"github.com/joho/godotenv"
	"io/ioutil"
	"plum/utils"
	"time"
)


func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal(err)
	}

	err = ioutil.WriteFile("blockchain.json", []byte("[]"), 0644)
	utils.Check(err)

	go func() {
		t := time.Now()
		genesisBlock := models.Block{0, t.String(), 0, "", "", models.Content{""}}
		utils.WriteBlockChain(genesisBlock)
	}()

	log.Fatal(utils.Run())
}