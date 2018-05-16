package main

import (
	"log"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
)

func main() {
	// Serve up client side
	fs := http.FileServer(http.Dir("html"))
	http.Handle("/", fs)

	// Websocket side
	http.HandleFunc("/websocket", func(w http.ResponseWriter, r *http.Request) {

		// the upgrader takes an http (tcp) connection and upgrades
		// it to a websocket connection with the following parameters
		upgrader := websocket.Upgrader{
			ReadBufferSize:  1024,
			WriteBufferSize: 1024,
		}

		// Upgrade connection to websocket
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			// Do not want to bring server down if client cannot upgrade
			log.Println(err)
			return
		}

		log.Println("Client subscribed")

		// This loop blocks on conn.ReadMessage()
		for {
			msgType, msg, err := conn.ReadMessage()
			if err != nil {
				log.Println(err)
				break
			}

			// If we get ping, wait 2 seconds then reply with pong.
			if string(msg) == "ping" {
				log.Println("ping received")
				time.Sleep(2 * time.Second)
				err = conn.WriteMessage(msgType, []byte("pong"))
				if err != nil {
					log.Println(err)
					break
				}
				log.Println("pong sent")

			} else {
				if string(msg) == "exit" {
					conn.WriteMessage(msgType, []byte("bye!"))
				} else {
					log.Println("Unexpected message received: " + string(msg) + ". terminating connection.")
				}
				conn.Close()
				break
			}

		}
		log.Println("Client unsubscribed")
	})

	// Start server with default handler
	http.ListenAndServe(":3000", nil)
}
