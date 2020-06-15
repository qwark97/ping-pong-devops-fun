package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"runtime"
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

func pinger() {
	var httpClient = &http.Client{
		Timeout: time.Second * 5,
	}
	msg := "Hello Goopher, I'm " + config.Name
	for {
		func() {
			greetings := bytes.NewBufferString(msg)
			resp, err := httpClient.Post(config.Address, "application/json", greetings)
			if err != nil {
				log.Print("Something went wrong during request")
				return
			}
			defer resp.Body.Close()
			body, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				log.Print("Something went wrong during reading response")
				return
			}
			log.Print(string(body))
		}()
		time.Sleep(time.Second * 5)
	}
}
func main() {

	config.loadConf("/home/marcin/go-code/src/pinger/conf.json")

	go func() {
		http.HandleFunc("/ping", serverHandler)
		log.Fatal(http.ListenAndServe(":8080", nil))
	}()

	go pinger()
	runtime.Goexit()
	fmt.Println("Exit")

}
