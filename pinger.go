package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

type data struct {
	Name    string `json:"name"`
	Address string `json:"address"`
}

func (d *data) loadConf(path string) {
	fileContent, err := ioutil.ReadFile(path)
	if err != nil {
		panic("Provided configuration path is not valid!")
	}

	err = json.Unmarshal(fileContent, d)
	if err != nil {
		panic("Configuration file is not valid JSON file!")
	}
}

var config data

func serverHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("PING from %s", r.Host)
	fmt.Fprintf(w, "Hi there, I'm PINGER. My name is %s", config.Name)
}
func main() {

	config.loadConf("/home/marcin/go-code/src/pinger/conf.json")

	go func() {
		http.HandleFunc("/ping", serverHandler)
		log.Fatal(http.ListenAndServe(":8080", nil))
	}()

	fmt.Println("Am here now")
	time.Sleep(time.Second * 10)
}
