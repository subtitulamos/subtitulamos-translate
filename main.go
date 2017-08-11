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

	"github.com/garyburd/redigo/redis"
)

var redisConn redis.Conn
var psConn redis.PubSubConn
var redisEnvPrefix *string

func main() {
	addr := flag.String("addr", ":8080", "http servicing address")
	redisPort := flag.String("redis-port", ":6379", "redis service address")
	redisEnvPrefix = flag.String("redis-pubsub-env", "dev", "redis pub/sub environment prefix")
	flag.Parse()

	var err error
	redisConn, err = redis.Dial("tcp", *redisPort)
	if err != nil {
		log.Fatal("redis connection failed: ", err)
	}
	psConn = redis.PubSubConn{Conn: redisConn}
	go redisListener()

	http.HandleFunc("/", serveWs)
	err = http.ListenAndServe(*addr, nil)
	if err != nil {
		log.Fatal("ListenAndServe failed: ", err)
	}
}
