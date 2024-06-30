package main

import (
	"encoding/json"
	"fmt"
	"log"
	"lru/cache"
	"net/http"
	"strconv"
	"time"
)

var c *cache.Cache

func deleteKey(w http.ResponseWriter, r *http.Request) {
	key := r.URL.Query().Get("key")
	if key == "" {
		http.Error(w, "Missing key parameter", http.StatusBadRequest)
		return
	}

	c.Delete(key)

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Key deleted successfully"))
}

func setKey(w http.ResponseWriter, r *http.Request) {
	var data struct {
		Key   string `json:"key"`
		Value string `json:"value"`
		TTL   string `json:"ttl"`
	}
	if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	ttl, err := strconv.Atoi(data.TTL)
	if err != nil {
		http.Error(w, "TTl should be a valid number", http.StatusBadRequest)
		return
	}
	c.Set(data.Key, data.Value, ttl)
	fmt.Println("set key", data.Key, data.Value, data.TTL)

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Key(s) added successfully"))
}

func getKey(w http.ResponseWriter, r *http.Request) {
	key := r.URL.Query().Get("key")
	if key == "" {
		http.Error(w, "Key missing", http.StatusBadRequest)
		return
	}

	value, exists := c.Get(key)
	if !exists {
		http.Error(w, "Key does not exist", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(value.(string)))
}

func getAllKeys(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	allKeys := c.GetAll()
	if err := json.NewEncoder(w).Encode(allKeys); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
	}
}

func mainHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		getKey(w, r)
	case http.MethodPost:
		setKey(w, r)
	case http.MethodDelete:
		deleteKey(w, r)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func ping(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	resp := map[string]string{
		"status": "OK",
	}
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
	}
}

func main() {
	c = cache.NewCache(50)
	c.StartCleanup(10 * time.Second)

	http.HandleFunc("/ping", ping)
	http.HandleFunc("/mycache", mainHandler)
	http.HandleFunc("/mycache/all", getAllKeys)

	fmt.Printf("Starting server at http://localhost:3001\n")
	if err := http.ListenAndServe(":3001", enableCors(http.DefaultServeMux)); err != nil {
		log.Fatal(err)
	}
}

func enableCors(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS, PUT, DELETE")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

		if r.Method == "OPTIONS" {
			return
		}

		next.ServeHTTP(w, r)
	})
}
