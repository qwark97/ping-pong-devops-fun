package main

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"log"
	"math/rand"
	"time"

	"github.com/go-redis/redis/v8"
)

var config data
var ctx = context.Background()

type data struct {
	Names []string `json:"names"`
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

func redisNewClient(port string) *redis.Client {
	addr := "redis:" + port
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
	config.loadConf("conf.json")
	rdb := redisNewClient("6379")
	var randIdx int
	var name string

	for {
		randIdx = rand.Intn(3)
		name = config.Names[randIdx]
		err := rdb.Set(ctx, "name", name, 0).Err()
		if err != nil {
			log.Print(err)
		}
		log.Printf("Name \"%s\" set successfully", name)
		time.Sleep(5 * time.Second)
	}
}
