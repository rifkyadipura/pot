package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"
)

type EchoRequest struct {
	Message string `json:"message"`
}
type EchoResponse struct {
	Message   string    `json:"message"`
	Timestamp time.Time `json:"timestamp"`
}

type Team struct {
	ID     int     `json:"id"`
	Name   string  `json:"name"`
	Score  int     `json:"score"`
	Lat    float64 `json:"lat"`
	Lng    float64 `json:"lng"`
	Status string  `json:"status"`
}

var teams = []Team{
	{1, "GreenSnake", 420, -6.2088, 106.8456, "online"},
	{2, "RedFox", 350, -6.9147, 107.6098, "online"},
	{3, "BlueTiger", 280, -7.2504, 112.7688, "under_attack"},
}

func main() {
	// GET /ping → health check sederhana (JSON)
	http.HandleFunc("/ping", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprintf(w, `{"status":"meriang","time":"%s"}`, time.Now().Format(time.RFC3339))
	})

	http.HandleFunc("/api/teams", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("Access-Control-Allow-Origin", "*")

		json.NewEncoder(w).Encode(teams)
	})

	// GET /hello?name=Rolly → contoh baca query parameter
	http.HandleFunc("/hello", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}
		name := r.URL.Query().Get("name")
		if name == "" {
			name = "Anon"
		}
		w.Header().Set("Content-Type", "text/plain")
		fmt.Fprintf(w, "Hello, %s!\n", name)
	})

	// POST /echo (JSON) → baca body JSON dan kembalikan lagi (echo) + timestamp
	http.HandleFunc("/echo", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}
		if r.Header.Get("Content-Type") != "" && r.Header.Get("Content-Type")[:16] != "application/json" {
			// tidak wajib keras, hanya contoh validasi sederhana
		}

		var req EchoRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "invalid JSON: "+err.Error(), http.StatusBadRequest)
			return
		}
		resp := EchoResponse{Message: req.Message, Timestamp: time.Now()}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(resp)
	})

	// POST /form (application/x-www-form-urlencoded) → baca form field
	http.HandleFunc("/form", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}
		if err := r.ParseForm(); err != nil {
			http.Error(w, "cannot parse form: "+err.Error(), http.StatusBadRequest)
			return
		}
		username := r.Form.Get("username")
		ageStr := r.Form.Get("age")
		age, _ := strconv.Atoi(ageStr)

		w.Header().Set("Content-Type", "application/json")
		fmt.Fprintf(w, `{"username":%q,"age":%d,"received_at":%q}`, username, age, time.Now().Format(time.RFC3339))
	})

	addr := ":8080"
	log.Println("Server jalan di http://localhost" + addr)
	log.Fatal(http.ListenAndServe(addr, nil))
}
