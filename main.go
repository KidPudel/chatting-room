package main

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/coder/websocket"
)

func client(ctx context.Context) {
	fmt.Println("client initiation")
	// handshake
	// send a message
	connection, r, err := websocket.Dial(ctx, "http://localhost:8000/connect", nil)
	if err != nil {
		log.Fatalf(fmt.Sprintf("failed to handshake: %s", err))
	}
	defer connection.Close(websocket.StatusNormalClosure, "end of session")
	fmt.Println(r)

	for {
		var message string
		read, err := fmt.Scanln(&message)
		if err != nil {
			fmt.Printf("error while listening for input: %s", err)
			continue
		}
		fmt.Println(read)

		connection.Write(ctx, websocket.MessageText, []byte(message))
	}
}

func server(ctx context.Context) (err error) {
	fmt.Println("server initiation")
	// accept
	// router to help find endpoint for start
	mux := http.NewServeMux()
	// connect to the socket
	mux.HandleFunc("GET /connect", func(w http.ResponseWriter, r *http.Request) {
		// listen for tcp handshake
		connection, err := websocket.Accept(w, r, nil)
		if err != nil {
			log.Fatalf(fmt.Sprintf("failed to accept a handshake from the client: %s", err))
			return
		}
		defer connection.Close(websocket.StatusNormalClosure, "end of session")

		for {
			msgType, message, err := connection.Read(ctx)
			if err != nil {
				fmt.Printf("failed to read the message: %s\n", err)
			}
			fmt.Printf("message: %s, type of %s", message, msgType)
			if string(message) == "end" {
				break
			}
		}
	})
	// listen
	http.ListenAndServe(":8000", mux)

	return nil
}

func main() {
	fmt.Print("Choose client or server (1/2)")
	var choice string
	read, err := fmt.Scanln(&choice)
	if err != nil {
		fmt.Printf("error during input: %s", err)
		return
	}
	if read != 1 {
		fmt.Printf("input must be 1 or 2")
		return
	}

	if choice == "1" {
		server(context.Background())
	} else {
		client(context.Background())
	}
}
