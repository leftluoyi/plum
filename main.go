package main

import (
	"plum/models"
	"plum/utils"

	"fmt"
	"github.com/joho/godotenv"
	"io/ioutil"
	"net"
	"os"
	"plum/services"
	"strconv"
	"time"
)

var bcServer chan []models.Block

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
	server, err := net.Listen("tcp", ":"+httpPort)
	utils.Check(err)
	fmt.Println("HTTP Server Listening on port:", httpPort)

	defer server.Close()

	for {
		conn, err := server.Accept()
		utils.Check(err)
		go services.HandleTcpConn(conn, bcServer)
	}

	//log.Fatal(services.Run())		# the HTTP server
}
