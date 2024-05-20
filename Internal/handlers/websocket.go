package handlers

import (
	"fmt"
	"github.com/gorilla/websocket"
	"gopkg.in/yaml.v3"
	"log"
	"net/http"
	"time"

	"Service_1Cv8/internal/exchange"
	"Service_1Cv8/internal/files"
)

// wsGetInfoQueues вебсокет сервера, отправляет информацию по очередям
func (srv *Server) wsGetInfoQueues(conn *websocket.Conn, w http.ResponseWriter) {
	for {
		time.Sleep(1 * time.Second)

		arrEQI := []exchange.ExchangeQueueInfo{}
		for k, v := range srv.ExchangeStorage {

			msgCounts := v.Len()
			msgSize := v.SizeMemory()
			//msgSize := 0

			arrEQI = append(arrEQI, exchange.ExchangeQueueInfo{
				Name:        k,
				DownloadAt:  v.DownloadAt,
				UploadAt:    v.UploadAt,
				TypeStorage: v.TypeStorage.String(),
				Size:        fmt.Sprintf("%s KBt", files.GroupSeparator(fmt.Sprintf("%d", msgSize/1024))),
				Messages:    fmt.Sprintf("%s", files.GroupSeparator(fmt.Sprintf("%d", msgCounts))),
			})
		}

		out, err := yaml.Marshal(arrEQI)
		if err != nil {
			log.Println("wsGetInfoQueues Marshal() error:", err)
			continue
		}

		err = conn.WriteMessage(1, out)
		if err != nil {
			log.Println("ws closed")
			return
		}
	}
}

// wsGetInfoQueues вебсокет сервера, отправляет информацию по очередям
func (srv *Server) wsGetInfoQueue(conn *websocket.Conn, r *http.Request) {
	qu := r.Header.Get("Queue")
	for {
		time.Sleep(1 * time.Second)

		eq, ok := srv.ExchangeStorage[qu]
		if !ok {
			eq = exchange.ExchangeQueue{}
		}
		msgCounts := eq.Len()
		msgSize := eq.SizeMemory()

		eqi := exchange.ExchangeQueueInfo{
			Name:        qu,
			DownloadAt:  eq.DownloadAt,
			UploadAt:    eq.UploadAt,
			TypeStorage: eq.TypeStorage.String(),
			Size:        fmt.Sprintf("%s KBt", files.GroupSeparator(fmt.Sprintf("%d", msgSize/1024))),
			Messages:    fmt.Sprintf("%s", files.GroupSeparator(fmt.Sprintf("%d", msgCounts))),
		}

		out, err := yaml.Marshal(eqi)
		if err != nil {
			log.Println("wsGetInfoQueue Marshal() error:", err)
			continue
		}

		err = conn.WriteMessage(1, out)
		if err != nil {
			log.Println("ws closed")
			return
		}
	}
}
