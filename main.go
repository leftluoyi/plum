package main

import (
	"log"
	"plum/models"

	"github.com/joho/godotenv"
	"io/ioutil"
	"plum/utils"
	"time"
	"os"
)


func main() {
	err := godotenv.Load("config.env")
	if err != nil {
		log.Fatal(err)
	}

	if(os.Getenv("REWRITE_BLOCKCHAIN") == "true") {
		err = ioutil.WriteFile("blockchain.json", []byte("[]"), 0644)
		utils.Check(err)
	}

	go func() {
		t := time.Now()
		genesisBlock := models.Block{0, t.String(), 0, "", "", models.Content{""}}
		utils.WriteBlockChain(genesisBlock)
	}()

	log.Fatal(utils.Run())
}