package main

import (
	"log"
	"net/http"
	"sync"
)

var (
	sum        = 0
	productNum = 10000
	mutex      = sync.Mutex{}
)

func GetProduct(w http.ResponseWriter, r *http.Request) {
	if GetOneProduct() {
		w.Write([]byte("true"))
	}
	w.Write([]byte("false"))
}

func GetOneProduct() bool {
	mutex.Lock()
	defer mutex.Unlock()
	if sum < productNum {
		sum++
		return true
	}
	return false
}

func main() {
	http.HandleFunc("/getOne", GetProduct)
	err := http.ListenAndServe(":12345", nil)
	if err != nil {
		log.Fatal(err)
	}
}
