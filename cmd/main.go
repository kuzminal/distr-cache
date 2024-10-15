package main

import (
	"distr-cache/internal/cache"
	"distr-cache/internal/server"
	"fmt"
	"net/http"
	"time"
)

func main() {
	cache := cache.NewCache(5) // Setting capacity to 5 for LRU
	cache.StartEvictionTicker(1 * time.Minute)
	cs := server.NewCacheServer(cache)
	http.HandleFunc("/set", cs.SetHandler)
	http.HandleFunc("/get", cs.GetHandler)
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		fmt.Println(err)
		return
	}
}
