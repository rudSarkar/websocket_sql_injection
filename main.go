package main

import (
	"flag"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/fatih/color"
	"github.com/gorilla/websocket"
)

var wsServer string

func init() {
	flag.StringVar(&wsServer, "ws", "ws://localhost:9091/ws", "websocket server address")
	flag.Usage = func() {
		fmt.Println("Usage: middleware [options]")
		fmt.Println("Options:")
		flag.PrintDefaults()
	}
}

func sendWS(payload map[string]string) (string, error) {
	u, err := url.Parse(wsServer)
	if err != nil {
		return "", err
	}

	c, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		return "", err
	}
	defer c.Close()

	// For our case, format the payload in JSON
	var data []string
	for key, value := range payload {
		data = append(data, fmt.Sprintf(`"%s":"%s"`, key, strings.Replace(value, `"`, `'`, -1)))
	}
	dataString := fmt.Sprintf("{%s}", strings.Join(data, ","))

	err = c.WriteMessage(websocket.TextMessage, []byte(dataString))
	if err != nil {
		return "", err
	}

	_, resp, err := c.ReadMessage()
	if err != nil {
		return "", err
	}

	return string(resp), nil
}

func middlewareServer(host, port string) error {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		query := r.URL.Query()
		if len(query) == 0 {
			w.Write([]byte("No parameters specified!"))
			return
		}

		now := time.Now().Format("2006-01-02 15:04:05")
		color.Green("[%s] Received request from client: %s\n", now, r.URL)

		payload := make(map[string]string, len(query))
		for key, values := range query {
			payload[key] = values[0]
		}

		resp, err := sendWS(payload)
		if err != nil {
			w.Write([]byte(err.Error()))
			return
		}

		w.Write([]byte(resp))
	})

	return http.ListenAndServe(fmt.Sprintf("%s:%s", host, port), nil)
}

func main() {
	flag.Parse()

	_, err := sendWS(map[string]string{})
	if err != nil {
		flag.Usage()
		os.Exit(1)
	} else {
		fmt.Println("[i] Starting Server")
		fmt.Println("[i] Send payloads in http://localhost:8000/?key1=value1&key2=value2...")
	}

	if err := middlewareServer("0.0.0.0", "8000"); err != nil {
		flag.Usage()
		os.Exit(1)
	}
}
