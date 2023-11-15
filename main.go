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
	"strconv"
	"strings"
	"time"
)

var CHAT_HOST string
var CHAT_PORT int

func init() {
	err := godotenv.Load(".env")

	if err != nil {
		log.Println("Error while loading .env file:", err)
	}

	CHAT_HOST = os.Getenv("CHAT_HOST")
	CHAT_PORT, _ = strconv.Atoi(os.Getenv("CHAT_PORT"))
}

func setHttpHeaders(w *http.ResponseWriter) {
	(*w).Header().Set("Access-Control-Allow-Origin", "*")
	(*w).Header().Set("Cache-Control", "no-cache")
	(*w).Header().Set("Connection", "keep-alive")
	(*w).Header().Set("Content-Type", "text/event-stream")
}

func sseResponse(w *http.ResponseWriter, message string) {
	message = fmt.Sprintf("data: %s\n", message)
	message += fmt.Sprintf("id: %s\n\n", uuid.New().String())
	log.Println("MESSAGE:", message)
	_, err := (*w).Write([]byte(message))

	if err != nil {
		log.Println("Error response write")
	}

	(*w).(http.Flusher).Flush()
}

func sse(w http.ResponseWriter, r *http.Request) {
	setHttpHeaders(&w)

	var m, _ = json.Marshal(JsonMessage{
		Id:      "ID-123",
		Type:    "chat/message",
		Service: "twitch",
		Text:    "Hello world from SERVER",
		TextSrc: "Hello world from SERVER",
		User: User{
			Id:       "ID-user-1",
			Nickname: "E.Wolf",
			Login:    "EWolf",
			Meta:     Meta{},
		},
	})

	for {
		log.Println("Send client data:", "Hello 101")
		sseResponse(&w, string(m))
		time.Sleep(2 * time.Second)
	}
}

func main() {
	go runTCPServer()

	http.HandleFunc("/chat", sse)
	http.Handle("/", http.FileServer(http.Dir("./public")))
	addr := fmt.Sprintf("%s:%d", CHAT_HOST, CHAT_PORT)

	log.Println("Server started at", addr)
	log.Fatal(http.ListenAndServe(addr, nil))
}

func runTCPServer() {
	server, err := net.Listen("tcp", "0.0.0.0:8889")
	if err != nil {
		fmt.Println("Error listening:", err.Error())
		os.Exit(1)
	} else {
		fmt.Println("OK!")
	}

	defer server.Close()

	for {
		conn, err := server.Accept()
		if err != nil {
			fmt.Println("Error accepting: ", err.Error())
			os.Exit(1)
		}

		go tcpAccept(conn)
	}
}

func tcpAccept(conn net.Conn) {
	msg := json.NewDecoder(conn)

	for {
		var message JsonMessage
		err := msg.Decode(&message)

		if err != nil {
			log.Println("Error while reading JSON message from TCP server", err)
			break
		}

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
