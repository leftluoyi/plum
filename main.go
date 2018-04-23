package main

import (
	"log"
	"plum/models"
	"plum/utils"

	"github.com/joho/godotenv"
	"io/ioutil"
	"time"
	"os"
	"strconv"
	"net"
	"bufio"
	"io"
	"fmt"
	"encoding/json"
	"github.com/davecgh/go-spew/spew"
	"sync"
)

var bcServer chan []models.Block
var mutex = &sync.Mutex{}

func main() {
	err := godotenv.Load("config.env")
	utils.Check(err)

	bcServer = make(chan []models.Block)

	rewrite, err := strconv.ParseBool(os.Getenv("REWRITE_BLOCKCHAIN"))
	utils.Check(err)

	if rewrite {
		err = ioutil.WriteFile("blockchain.json", []byte("[]"), 0644)
		utils.Check(err)
		t := time.Now()
		genesisBlock := models.Block{0, t.String(), 0, "", "", 0, "", models.Content{""}}
		models.AppendToBlockChain(genesisBlock)
	}

	httpPort := os.Getenv("TCP_PORT")
	server, err := net.Listen("tcp", ":" + httpPort)
	utils.Check(err)
	fmt.Println("HTTP Server Listening on port:", httpPort)

	defer server.Close()

	for {
		conn, err := server.Accept()
		utils.Check(err)
		go handleConn(conn)
	}

	//log.Fatal(services.Run())		# the HTTP server
}


func handleConn(conn net.Conn) {
	defer conn.Close()

	io.WriteString(conn, "Enter a new BPM:")

	scanner := bufio.NewScanner(conn)

	// take in BPM from stdin and add it to blockchain after conducting necessary validation
	go func() {
		for scanner.Scan() {
			blockchain := models.GetBlockChain()
			bpm, err := strconv.Atoi(scanner.Text())
			if err != nil {
				log.Printf("%v not a number: %v", scanner.Text(), err)
				continue
			}
			content := models.Content{"empty"}
			newBlock, err := blockchain[len(blockchain)-1].GenerateNextBlock(bpm, content)
			utils.Check(err)
			if newBlock.IsBlockValid(blockchain[len(blockchain)-1]) {
				fmt.Println("Valid")
				newBlockchain := append(blockchain, newBlock)
				models.ReplaceChain(newBlockchain)
			}

			bcServer <- models.GetBlockChain()
			io.WriteString(conn, "\nEnter a new BPM:")
		}
	}()

	go func() {
		for {
			blockchain := models.GetBlockChain()
			time.Sleep(30 * time.Second)
			mutex.Lock()
			output, err := json.Marshal(blockchain)
			if err != nil {
				log.Fatal(err)
			}
			mutex.Unlock()
			io.WriteString(conn, string(output))
		}
	}()

	for _ = range bcServer {
		spew.Dump(models.GetBlockChain())
	}
}