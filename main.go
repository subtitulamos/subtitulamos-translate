/**
 * This file is covered by the AGPLv3 license, which can be found at the LICENSE file in the root of this project.
 * @copyright 2013 The Gorilla WebSocket Authors
 * @copyright 2017 subtitulamos.tv
 */

package main

import (
	"flag"
	"log"
	"net/http"

	"github.com/go-redis/redis"
)

var redisClient *redis.Client
var psConn *redis.PubSub
var redisEnvPrefix *string

func main() {
	addr := flag.String("http-addr", ":8080", "http servicing address")
	redisAddr := flag.String("redis-addr", ":6379", "redis service address")
	redisEnvPrefix = flag.String("redis-pubsub-env", "dev", "redis pub/sub environment prefix")
	flag.Parse()

	redisClient = redis.NewClient(&redis.Options{
		Addr: *redisAddr,
		DB:   0,
	})

	_, err := redisClient.Ping().Result()
	if err != nil {
		log.Fatal("pinging redis failed: ", err)
	}

	psConn = redisClient.Subscribe("") // Empty subscription
	go redisListener()

	http.HandleFunc("/", serveWs)
	err = http.ListenAndServe(*addr, nil)
	if err != nil {
		log.Fatal("ListenAndServe failed: ", err)
	}
}
