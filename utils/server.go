package utils

import (
	"net/http"
	"github.com/gorilla/mux"
	"os"
	"log"
	"io"
	"time"
	"encoding/json"
	"plum/models"
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
	io.WriteString(w, GetBlockChainString())
}

type Message struct {
	BPM int
	Content	models.Content
}

func handleWriteBlock(w http.ResponseWriter, r *http.Request) {
	var m Message
	var blockchain = GetBlockChain()

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
		WriteBlockChain(newBlock)
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
