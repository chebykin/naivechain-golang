package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"time"

	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
)

var (
	upgrader = websocket.Upgrader{}
	sockets  = make(map[*websocket.Conn]bool, 0)
)

type MessageType int

const (
	QUERY_LATEST MessageType = iota
	QUERY_ALL
	RESPONSE_BLOCKCHAIN
)

func write(conn *websocket.Conn, message interface{}) {
	fmt.Printf(">>> %#v\n", message)
	// err := conn.WriteJSON(message)
	// msg, _ := json.Marshal(message)
	err := conn.WriteMessage(websocket.TextMessage, []byte(`{"type": 0}`))
	if err != nil {
		log.Println("write:", err)
	}
}

func peerHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("New Peer")
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Upgrade:", err)
		return
	}

	// conn.SetPongHandler(func(string) error {
	// 	conn.SetReadDeadline(time.Now().Add(10 * time.Second))
	// 	return nil
	// })

	sockets[conn] = true

	defer func() {
		log.Println("xxx: closing connection")
		conn.Close()
		delete(sockets, conn)
	}()

	peerReader(conn)
	// peerWriter(conn)
}

type PeerMessage struct {
	Type MessageType
	Data interface{}
}

func peerReader(conn *websocket.Conn) {
	for {
		var msg PeerMessage
		// err := conn.ReadJSON(&msg)
		_, message, err := conn.ReadMessage()
		fmt.Println("peer read <<<", msg)
		if err != nil {
			fmt.Println("p2p reader:", err)
			return
		}

		err = json.Unmarshal(message, &msg)
		if err != nil {
			fmt.Println("p2p unmarshaler:", err)
			return
		}

		messageHandler(conn, msg)
	}
	fmt.Println("out of loop")
}

func peerWriter(ws *websocket.Conn) {
	for {
		time.Sleep(1 * time.Minute)

	}
}

func messageHandler(conn *websocket.Conn, msg PeerMessage) {
	switch msg.Type {
	case QUERY_LATEST:
		fmt.Println("Request for a latest")
		write(conn, struct {
			Type MessageType
			Data []*Block
		}{RESPONSE_BLOCKCHAIN, []*Block{latestBlock()}})
	case QUERY_ALL:
		fmt.Println("Request for all")
		write(conn, struct {
			Type MessageType
			Data []*Block
		}{RESPONSE_BLOCKCHAIN, blockchain})
	case RESPONSE_BLOCKCHAIN:
		fmt.Println("Request for latest")
	}
}

func connectToPeers(peers []string) {
	for _, host := range peers {
		u := url.URL{Scheme: "ws", Host: host, Path: "/p2p"}
		conn, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
		log.Printf("p2p to %s", u.String())
		if err != nil {
			log.Println("p2p dial:", err)
			return
		}

		defer conn.Close()

		go peerReader(conn)

		write(conn, struct {
			Type MessageType
		}{QUERY_LATEST})
	}
}

func initP2PServer() {
	r := mux.NewRouter()

	r.HandleFunc("/p2p", peerHandler)

	s := &http.Server{
		Addr:    fmt.Sprintf(":%s", p2p_port),
		Handler: r,
	}

	fmt.Printf("Listening websocket p2p on port http://localhost:%s\n", p2p_port)
	s.ListenAndServe()
}
