package main

import (
	"distr-cache/internal/cache"
	"distr-cache/internal/server"
	"flag"
	"log/slog"
	"net/http"
	"strings"
	"time"
)

var (
	port  string
	peers string
)

func main() {
	flag.StringVar(&port, "port", ":8080", "HTTP server port")
	flag.StringVar(&peers, "peers", "", "Comma-separated list of peer addresses")
	flag.Parse()

	peerList := strings.Split(peers, ",")
	cache := cache.NewCache(5) // Setting capacity to 5 for LRU
	cache.StartEvictionTicker(1 * time.Minute)
	cs := server.NewCacheServer(cache, peerList)
	http.HandleFunc("/set", cs.SetHandler)
	http.HandleFunc("/get", cs.GetHandler)
	slog.Info("Starting server on port", "port", port)
	err := http.ListenAndServe(port, nil)
	if err != nil {
		slog.Error(err.Error())
	}
}
