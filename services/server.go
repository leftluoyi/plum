package services

import (
	"net/http"
	"github.com/gorilla/mux"
	"os"
	"log"
	"io"
	"time"
	"encoding/json"
	"plum/models"
	"strconv"
	"fmt"
	"net"
	"bufio"
	"plum/utils"
	"sync"
)

func Run() error {
	muxRouter := mux.NewRouter()
	muxRouter.HandleFunc("/", handleGetBlockchain).Methods("GET")
	muxRouter.HandleFunc("/", handleWriteBlock).Methods("POST")

	httpAddr := os.Getenv("ADDR")
	log.Println("Listening on ", os.Getenv("ADDR"))
	s := &http.Server{
		Addr:           ":" + httpAddr,
		Handler:        muxRouter,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	if err := s.ListenAndServe(); err != nil {
		return err
	}

	return nil
}

func handleGetBlockchain(w http.ResponseWriter, r *http.Request) {
	io.WriteString(w, models.GetBlockChainString())
}

type Message struct {
	BPM int
	Content	models.Content
}

func handleWriteBlock(w http.ResponseWriter, r *http.Request) {
	var m Message
	var blockchain = models.GetBlockChain()

	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&m); err != nil {
		respondWithJSON(w, r, http.StatusBadRequest, r.Body)
		return
	}
	defer r.Body.Close()

	newBlock, err := blockchain[len(blockchain)-1].GenerateNextBlock(m.BPM, m.Content)
	if err != nil {
		respondWithJSON(w, r, http.StatusInternalServerError, m)
		return
	}
	if newBlock.IsBlockValid(blockchain[len(blockchain)-1]) {
		models.AppendToBlockChain(newBlock)
		//replaceChain(newBlockchain, blockchain)
		//spew.Dump(Blockchain)
	}

	respondWithJSON(w, r, http.StatusCreated, newBlock)
}

func replaceChain(newBlocks []models.Block, blockchain []models.Block) {
	if len(newBlocks) > len(blockchain) {
		blockchain = newBlocks
	}
}

func respondWithJSON(w http.ResponseWriter, r *http.Request, code int, payload interface{}) {
	response, err := json.MarshalIndent(payload, "", "  ")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("HTTP 500: Internal Server Error"))
		return
	}
	w.WriteHeader(code)
	w.Write(response)
}

func HandleTcpConn(conn net.Conn, bcServer chan []models.Block) {
	defer conn.Close()

	var mutex = &sync.Mutex{}
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

	for true {
		<- bcServer
	}
}

