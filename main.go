package main

import (
	"plum/models"
	"plum/utils"

	"fmt"
	"github.com/joho/godotenv"
	"io/ioutil"
	"os"
	golog "gx/ipfs/QmTG23dvpBCBjqQwyDxV8CQT6jmS4PSftNr1VqHhE3MLy7/go-log"
	gologging "gx/ipfs/QmQvJiADDe7JR4m968MwXobTCCzUqQkP87aRHe29MEBGHV/go-logging"
	ma "gx/ipfs/QmWWQ2Txc2c6tqjsBpzg5Ar652cHPGNsQQp2SejkNmkUMb/go-multiaddr"
	peer "gx/ipfs/QmcJukH2sAFjY3HdBKq35WDzWoL3UUu2gt9wdfqZTUyM74/go-libp2p-peer"
	pstore "gx/ipfs/QmdeiKhUy1TVGBaKxt7y1QmBDLBdisSrLJ1x58Eoj4PXUh/go-libp2p-peerstore"

	"strconv"
	"time"
	"bufio"
	"flag"
	"log"
	"plum/services"
	"context"
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

	golog.SetAllLoggers(gologging.INFO) // Change to DEBUG for extra info

	listenF := flag.Int("l", 0, "wait for incoming connections")
	target := flag.String("d", "", "target peer to dial")
	secio := flag.Bool("secio", false, "enable secio")
	seed := flag.Int64("seed", 0, "set random seed for id generation")
	flag.Parse()

	if *listenF == 0 {
		log.Fatal("Please provide a port to bind on with -l")
	}

	// Make a host that listens on the given multiaddress
	ha, err := services.MakeBasicHost(*listenF, *secio, *seed)
	if err != nil {
		log.Fatal(err)
	}

	if *target == "" {
		log.Println("listening for connections")
		//ha.SetStreamHandler("/p2p/1.0.0", handleStream)

		select {} // hang forever
	} else {
		ipfsaddr, err := ma.NewMultiaddr(*target)
		if err != nil {
			log.Fatalln(err)
		}

		pid, err := ipfsaddr.ValueForProtocol(ma.P_IPFS)
		if err != nil {
			log.Fatalln(err)
		}

		peerid, err := peer.IDB58Decode(pid)
		if err != nil {
			log.Fatalln(err)
		}

		targetPeerAddr, _ := ma.NewMultiaddr(
			fmt.Sprintf("/ipfs/%s", peer.IDB58Encode(peerid)))
		targetAddr := ipfsaddr.Decapsulate(targetPeerAddr)

		ha.Peerstore().AddAddr(peerid, targetAddr, pstore.PermanentAddrTTL)

		log.Println("opening stream")
		s, err := ha.NewStream(context.Background(), peerid, "/p2p/1.0.0")
		if err != nil {
			log.Fatalln(err)
		}

		rw := bufio.NewReadWriter(bufio.NewReader(s), bufio.NewWriter(s))

		go services.WriteData(rw)
		go services.ReadData(rw)

		select {}

	}
}
