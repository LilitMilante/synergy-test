package main

import (
	"bufio"
	"encoding/json"
	"log"
	"net"
	"net/http"
	"sync"
)

var messages = make(map[string]interface{})
var mu = sync.Mutex{}

func main() {
	http.HandleFunc("/test", Handler)

	go http.ListenAndServe(":8081", nil)

	conn, err := net.Dial("tcp", ":8080")
	if err != nil {
		log.Fatal(err)
	}

	s := bufio.NewScanner(conn)

	for s.Scan() {
		mu.Lock()
		err := json.Unmarshal(s.Bytes(), &messages)
		mu.Unlock()

		if err != nil {
			log.Println(err)
			continue
		}
	}
}

func Handler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte(GetAction()))
}

func GetAction() string {
	mu.Lock()
	defer mu.Unlock()

	action := messages["action"]

	actionStr, ok := action.(string)

	if !ok {
		return GetAction()
	}

	return actionStr
}
