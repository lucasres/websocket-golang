package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/gobwas/ws"
	"github.com/gobwas/ws/wsutil"
)

type Message struct {
	Action string `json:"action"`
}

var data []string = []string{
	"{\"name\": \"Monza\", \"year\": \"1999\"}",
	"{\"name\": \"Marea\", \"year\": \"2000\"}",
	"{\"name\": \"Kaddet\", \"year\": \"1996\"}",
	"{\"name\": \"Uno Mile\", \"year\": \"2001\"}",
	"{\"name\": \"Clio\", \"year\": \"2001\"}",
	"{\"name\": \"Pegout\", \"year\": \"2001\"}",
}

func main() {
	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		conn, _, _, err := ws.UpgradeHTTP(r, w)
		if err != nil {
			log.Printf("err when update http connection: %e\n", err)
			return
		}

		go func() {
			defer conn.Close()
			log.Println("client connected")

			for {
				msg, op, err := wsutil.ReadClientData(conn)
				if err != nil {
					log.Printf("error when read client data: %e\n", err)
					return
				}

				if op == ws.OpClose {
					return
				}

				decoded := &Message{}
				err = json.Unmarshal(msg, decoded)
				if err != nil {
					log.Printf("err when decode client message: %e\n", err)
					return
				}

				if decoded.Action == "getProducts" {
					for _, v := range data {
						err := wsutil.WriteServerMessage(conn, op, []byte(v))
						if err != nil {
							log.Printf("err when write client message: %e\n", err)
						}
						continue
					}
				} else {
					return
				}
			}
		}()
	})

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "../front/index.html")
	})

	http.ListenAndServe(":8000", nil)
}
