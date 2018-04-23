package main

import (
	"log"
	"plum/models"
	"plum/utils"
	"plum/services"


	"github.com/joho/godotenv"
	"io/ioutil"
	"time"
	"os"
	"strconv"
)


func main() {
	err := godotenv.Load("config.env")
	if err != nil {
		log.Fatal(err)
	}

	rewrite, err := strconv.ParseBool(os.Getenv("REWRITE_BLOCKCHAIN"))
	if err != nil {
		log.Fatal(err)
	}

	if rewrite {
		err = ioutil.WriteFile("blockchain.json", []byte("[]"), 0644)
		utils.Check(err)
	} else {
		t := time.Now()
		genesisBlock := models.Block{0, t.String(), 0, "", "", 0, "", models.Content{""}}
		services.WriteBlockChain(genesisBlock)
	}


	log.Fatal(services.Run())
}