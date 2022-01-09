package main

import (
	"errors"
	"io/ioutil"
	"log"
	"math/rand"
	"net"
	"net/http"
	"sync"
	"time"
)

var clients = make(map[net.Conn]struct{}, 0)
var mu = sync.Mutex{}

func main() {
	c := make(chan []byte)

	go Connection()
	go RandMakeRequest("https://novasite.su/test1.php", c)
	go RandMakeRequest("https://novasite.su/test2.php", c)
	for v := range c {
		mu.Lock()
		for client := range clients {
			v = append(v, '\n')
			_, err := client.Write(v)
			if err != nil {
				delete(clients, client)
			}
		}
		mu.Unlock()
	}
}

func Connection() {
	listener, err := net.Listen("tcp", ":8080")
	if err != nil {
		log.Fatal(err)
	}

	defer listener.Close()

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Fatal(err)
		}

		mu.Lock()

		clients[conn] = struct{}{}

		mu.Unlock()
	}
}

func MakeRequest(adr string) ([]byte, error) {
	respOne, err := http.Get(adr)
	if err != nil {
		return nil, err
	}

	body, err := ioutil.ReadAll(respOne.Body)
	if err != nil {
		return nil, err
	}

	if respOne.StatusCode != http.StatusOK {
		return nil, errors.New("status code not 200")
	}

	return body, nil
}

func RandMakeRequest(adr string, c chan []byte) {
	for {
		time.Sleep(time.Second * time.Duration(rand.Intn(3)+1))
		body, err := MakeRequest(adr)
		if err != nil {
			log.Println(err)
			continue
		}
		c <- body
	}
}
