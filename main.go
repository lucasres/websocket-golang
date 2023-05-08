package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gobwas/ws"
	"github.com/gobwas/ws/wsutil"
)

type Message struct {
	Action string    `json:"action"`
	Params CarDetail `json:"params"`
}

type CarDetail struct {
	Id   string `json:"id"`
	Name string `json:"name"`
	Year string `json:"year"`
}

var data map[string]string = map[string]string{
	"1": "{\"id\": \"1\", \"name\": \"Monza\", \"year\": \"1999\"}",
	"2": "{\"id\": \"2\", \"name\": \"Marea\", \"year\": \"2000\"}",
	"3": "{\"id\": \"3\", \"name\": \"Kaddet\", \"year\": \"1996\"}",
	"4": "{\"id\": \"4\", \"name\": \"Uno Mile\", \"year\": \"2001\"}",
	"5": "{\"id\": \"5\", \"name\": \"Clio\", \"year\": \"2001\"}",
	"6": "{\"id\": \"6\", \"name\": \"Pegout\", \"year\": \"2001\"}",
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
					log.Printf("err when decode client message: %v\n", err)
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
				} else if decoded.Action == "update" {
					newValue := fmt.Sprintf(
						"{\"id\": \"%s\", \"name\": \"%s\", \"year\": \"%s\"}",
						decoded.Params.Id,
						decoded.Params.Name,
						decoded.Params.Year,
					)

					log.Printf("atualizando para %s", newValue)
					data[decoded.Params.Id] = newValue
				} else {
					return
				}
			}
		}()
	})

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "./index.html")
	})

	http.ListenAndServe(":8000", nil)
}
