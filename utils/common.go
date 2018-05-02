package utils

import (
	"log"
	"sync"
)

var mutex = &sync.Mutex{}

func GetMutex() *sync.Mutex {
	return mutex
}

func Check(e error) {
	if e != nil {
		log.Fatal(e)
	}
}
