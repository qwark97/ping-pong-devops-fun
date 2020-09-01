package main

import (
	"bytes"
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"runtime"
	"time"

	"github.com/go-redis/redis/v8"
)

type data struct {
	Name    string `json:"name"`
	Address string `json:"address"`
}

var config data
var ctx = context.Background()

func serverHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("PING from %s", r.Host)
	fmt.Fprintf(w, "Hi there, I'm PINGER. My name is %s", config.Name)
}

func namesGiver(rdb *redis.Client) {
	for {
		func() {
			val, err := rdb.Get(ctx, "name").Result()
			if err != nil {
				if err == redis.Nil {
					fmt.Println("There is no name")
					return
				}
			}
			if val == config.Name {
				log.Print("There is no new name")
				log.Print("Let's wait and try again")
				time.Sleep(time.Second * 3)
			} else {
				config.Name = val
				log.Print("New given name: ", string(config.Name))
			}
		}()
		time.Sleep(time.Second * 3)
	}
}

func pinger() {
	var httpClient = &http.Client{
		Timeout: time.Second * 5,
	}

	for {
		func() {
			msg := "Hello Goopher, I'm " + config.Name
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
			log.Print("Pinged response: ", string(body))
		}()
		time.Sleep(time.Second * 5)
	}
}

func redisNewClient(port string) *redis.Client {
	addr := "redis-host:" + port
	rdb := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	err := rdb.Ping(ctx).Err()
	if err != nil {
		panic(err)
	}
	return rdb
}

func main() {

	rdb := redisNewClient("6379")
	config.Address = "http://pinged-host:5050/ping"

	go func() {
		http.HandleFunc("/ping", serverHandler)
		log.Fatal(http.ListenAndServe(":8080", nil))
	}()

	go namesGiver(rdb)
	time.Sleep(time.Millisecond * 500)
	go pinger()
	runtime.Goexit()
	fmt.Println("Exit")

}
