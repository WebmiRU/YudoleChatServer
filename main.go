package main

import (
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"github.com/joho/godotenv"
	"log"
	"net"
	"net/http"
	"os"
	"slices"
	"strings"
	"sync"
)

var config Config
var sseClientsMutex sync.Mutex
var sseClietns []chan string
var out = make(chan TypeMessage, 9999)

func init() {
	err := godotenv.Load(".env")

	if err != nil {
		log.Println("Error while loading .env file:", err)
	}
}

func configLoad() {
	configFile, _ := os.ReadFile("config.json")
	err := json.Unmarshal(configFile, &config)

	if err != nil {
		log.Println("Error while read config file", err)
		os.Exit(1)
	}
}

func setHttpHeaders(w *http.ResponseWriter) {
	(*w).Header().Set("Access-Control-Allow-Origin", "*")
	(*w).Header().Set("Cache-Control", "no-cache")
	(*w).Header().Set("Connection", "keep-alive")
	(*w).Header().Set("Content-Type", "text/event-stream")
	(*w).(http.Flusher).Flush()
}

func sseResponse(w *http.ResponseWriter, message string) {
	message = fmt.Sprintf("data: %s\n", message)
	message += fmt.Sprintf("id: %s\n\n", uuid.New().String())
	log.Println("MESSAGE:", message)
	_, err := (*w).Write([]byte(message))

	if err != nil {
		log.Println("Error response write", err)
	}

	(*w).(http.Flusher).Flush()
}

func sseHandler(w http.ResponseWriter, r *http.Request) {
	setHttpHeaders(&w)

	c := make(chan string, 999)

	sseClientsMutex.Lock()
	sseClietns = append(sseClietns, c)
	sseClientsMutex.Unlock()

	for {
		select {
		case message := <-c:
			sseResponse(&w, message)
		case <-r.Context().Done():
			log.Println("SSE client closed connection")

			sseClientsMutex.Lock()
			idx := slices.Index(sseClietns, c)
			sseClietns = slices.Delete(sseClietns, idx, idx+1)
			sseClientsMutex.Unlock()
			close(c)

			return
		}
	}
}

func sseBroadcast() {
	for {
		var msg, _ = <-out
		log.Println("INCOME MESSAGE")
		message, _ := json.Marshal(msg)

		sseClientsMutex.Lock()
		for _, c := range sseClietns {
			c <- string(message)
		}
		sseClientsMutex.Unlock()
	}
}

func main() {
	configLoad()

	go runTCPServer()
	go sseBroadcast()

	http.HandleFunc("/chat", sseHandler)
	http.Handle("/", http.FileServer(http.Dir("./public")))
	addr := fmt.Sprintf("%s:%d", config.Servers.Host.Http.Address, config.Servers.Host.Http.Port)

	log.Println("Server started at", addr)
	log.Fatal(http.ListenAndServe(addr, nil))
}

func runTCPServer() {
	socket, err := net.Listen("tcp", fmt.Sprintf("%s:%d", config.Servers.Host.Server.Address, config.Servers.Host.Server.Port))
	if err != nil {
		fmt.Println("Error listening:", err.Error())
		os.Exit(1)
	}

	defer socket.Close()

	for {
		conn, err := socket.Accept()
		if err != nil {
			log.Println("Error accepting:", err.Error())
			os.Exit(1)
		}

		go tcpAccept(conn)
	}
}

func tcpAccept(conn net.Conn) {
	msg := json.NewDecoder(conn)

	for {
		var message TypeMessage
		err := msg.Decode(&message)

		if err != nil {
			log.Println("Error while reading JSON message from TCP server", err)
			break
		}

		log.Println("SEND OUT MESSAGE", message)
		out <- message

		switch strings.ToLower(message.Type) {
		case "chat/message":

			break

		case "chat/private_message":

			break

		case "system/message":

			break
		}
	}

	conn.Close()
}
