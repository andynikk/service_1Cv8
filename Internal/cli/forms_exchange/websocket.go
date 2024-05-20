package forms_exchange

import (
	"context"
	"fmt"
	"gopkg.in/yaml.v3"
	"net/http"
	"os"
	"sort"
	"time"

	"github.com/charmbracelet/bubbles/table"
	"github.com/gorilla/websocket"

	"Service_1Cv8/internal/constants"
	"Service_1Cv8/internal/repository"
)

func (m *FormListQueues) wsQueuesInfo(ctx context.Context) {

	ticker := time.NewTicker(time.Duration(30) * time.Millisecond)

	hostname, err := os.Hostname()
	if err != nil {
		hostname = "localhost"
	}

	port := constants.Port

	data := repository.DataJSON{}
	_ = repository.GetPudgeSetting(&data.Settings)
	if data.Settings.HTTPPort != "" {
		port = data.Settings.HTTPPort
	}

	m.message = ""
	m.errorConnet = false
	socketUrl := fmt.Sprintf("ws://%s:%s/socket/queues", hostname, port)
	conn, _, err := websocket.DefaultDialer.Dial(socketUrl, nil)
	if err != nil {
		m.message = "Connect error"

		_, ok := <-m.chanOut
		_ = ok
		m.errorConnet = true

		return
	}
	defer conn.Close()

	for {
		select {
		case <-ticker.C:

			res, ok := <-m.chanOut
			if !ok || !res {
				//m.message = err.Error()
				return
			}

			_, messageContent, err := conn.ReadMessage()
			if err != nil {

				_, ok = <-m.chanOut
				_ = ok

				m.message = err.Error()
				m.errorConnet = true

				return
			}
			err = yaml.Unmarshal(messageContent, &m.itemsDB)
			if err != nil {
				m.message = err.Error()
				continue
			}
			sort.Slice(m.itemsDB, func(i, j int) bool {
				return m.itemsDB[i].Name < m.itemsDB[j].Name
			})

			m.rows = []table.Row{}

			pp := 0
			for _, v := range m.itemsDB {
				pp++
				rowT := table.Row{
					fmt.Sprintf("%d", pp),
					v.Name,
					v.TypeStorage,
					fmt.Sprintf("%v", v.DownloadAt),
					fmt.Sprintf("%v", v.UploadAt),
					v.Messages,
					v.Size,
				}
				m.rows = append(m.rows, rowT)
			}
		case <-ctx.Done():
			return
		}
	}
}

func (m *FormQueue) wsQueueInfo(ctx context.Context) {

	ticker := time.NewTicker(time.Duration(30) * time.Millisecond)

	hostname, err := os.Hostname()
	if err != nil {
		hostname = "localhost"
	}

	port := constants.Port

	data := repository.DataJSON{}
	_ = repository.GetPudgeSetting(&data.Settings)
	if data.Settings.HTTPPort != "" {
		port = data.Settings.HTTPPort
	}

	h := http.Header{}
	h.Add("Queue", m.inputs.name.Value())

	m.message = ""
	m.errorConnet = false
	socketUrl := fmt.Sprintf("ws://%s:%s/socket/queue", hostname, port)
	conn, _, err := websocket.DefaultDialer.Dial(socketUrl, h)
	if err != nil {
		m.message = "Connect error"

		_, ok := <-m.chanOut
		_ = ok
		m.errorConnet = true

		return
	}
	defer conn.Close()

	for {
		select {
		case <-ticker.C:

			res, ok := <-m.chanOut
			if !ok || !res {
				//m.message = err.Error()
				return
			}

			_, messageContent, err := conn.ReadMessage()
			if err != nil {

				_, ok = <-m.chanOut
				_ = ok

				m.message = err.Error()
				m.errorConnet = true

				return
			}
			err = yaml.Unmarshal(messageContent, &m.exchangeQueueInfo)
			if err != nil {
				m.message = err.Error()
				continue
			}

			m.refreshData()

		case <-ctx.Done():
			return
		}
	}
}
