package handlers

import (
	"Service_1Cv8/internal/encryption"
	"Service_1Cv8/internal/exchange"
	"Service_1Cv8/internal/token"
	"context"
	"fmt"
	"github.com/boltdb/bolt"
	"github.com/go-ole/go-ole"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"sync"
	"syscall"
	"time"

	"Service_1Cv8/internal/constants"
	"Service_1Cv8/internal/repository"
	"github.com/gorilla/mux"
)

type Setting struct {
	HTTPPort          string    `yaml:"http_port"`
	IntervalService   int       `yaml:"interval_service"`
	IntervalWebClient int       `yaml:"interval_web_client"`
	KillProcessKb     int       `yaml:"kill_process_kb"`
	ExtremeStartKill  time.Time `yaml:"extreme_start_kill_process"`
	ExtremeStartDD    time.Time `yaml:"extreme_start_drop_double"`
	ControlServer     string    `yaml:"control_server"`
	PrivateKey        *encryption.KeyEncryption
	MessageCom        []repository.BasesDoubleControl
	ClosedConnects    []repository.ClosedConnect
	ClosedTasks       []repository.ClosedTask
}

//var ClaimsStore []token.ClaimStore

type Server struct {
	*mux.Router
	Setting
	BoltDB      *bolt.DB
	ClaimsStore []token.ClaimStore

	sync.RWMutex
	ExchangeStorage exchange.ExchangeStorage
}

// NewServer создание сервера
func NewServer() *Server {
	srv := &Server{}

	srv.InitRouters()
	srv.InitConfig()
	srv.InitData()

	return srv
}

func (srv *Server) InitRouters() {
	r := mux.NewRouter()

	r.HandleFunc("/", srv.handleFunc).Methods("GET")
	r.HandleFunc("/update", srv.handleFuncUpdateDB).Methods("GET")
	r.HandleFunc("/queues", srv.handleFuncQueues).Methods("GET")

	r.HandleFunc("/api/checkqueue/{nameQ}", srv.handleFuncCheckQueue).Methods("GET")
	r.HandleFunc("/api/putqueue", srv.handleFuncPutQueue).Methods("POST")
	r.HandleFunc("/api/pickqueue/{nameQ}", srv.handleFuncPickQueue).Methods("GET")

	r.HandleFunc("/api/addqueue/{nameQ}", srv.handleFuncAddQueue).Methods("GET")
	r.HandleFunc("/api/delqueue/{nameQ}", srv.handleFuncDelQueue).Methods("GET")
	r.HandleFunc("/api/clearqueue/{nameQ}", srv.handleFuncClearQueue).Methods("GET")

	r.HandleFunc("/api/tg/send", srv.handleFuncTgSend).Methods("POST")

	var upgrader = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
	}
	r.HandleFunc("/socket/queues", func(w http.ResponseWriter, r *http.Request) {
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			log.Println(err)
			return
		}

		srv.wsGetInfoQueues(conn, w)
	})

	r.HandleFunc("/socket/queue", func(w http.ResponseWriter, r *http.Request) {
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			log.Println(err)
			return
		}

		srv.wsGetInfoQueue(conn, r)
	})

	//r.Handle("/api/putqueue", midware.IsAuthorized(srv.handleFuncPutQueue)).Methods("POST")

	//r.HandleFunc("/api/pickqueue1/{nameQ}", srv.handleFuncPickQueue1).Methods("GET")

	//r.HandleFunc("/api/pickbucket/{nameB}", srv.handleFuncPickBucket).Methods("GET")
	//r.HandleFunc("/api/putbucket", srv.handleFuncPutBucket).Methods("POST")
	//r.HandleFunc("/api/delbucket", srv.handleFuncDelBucket).Methods("POST")
	//
	//r.HandleFunc("/api/pickpudge/{nameP}", srv.handleFuncPickPudge).Methods("GET")
	//r.HandleFunc("/api/putpudge", srv.handleFuncPutPudge).Methods("POST")
	//r.HandleFunc("/api/clearpudge", srv.handleFuncClearPudge).Methods("POST")
	//
	//r.HandleFunc("/api/pickpudge/{nameP}", srv.handleFuncPickPudge).Methods("GET")
	//r.HandleFunc("/api/putpudge", srv.handleFuncPutPudge).Methods("POST")
	//r.HandleFunc("/api/clearpudge", srv.handleFuncClearPudge).Methods("POST")

	srv.Router = r
}

func (srv *Server) InitConfig() {
	var bdc []repository.BasesDoubleControl

	data := repository.DataJSON{}
	data.GetPudgeData()

	srv.HTTPPort = data.Settings.HTTPPort
	i, err := strconv.Atoi(data.Settings.IntervalService)
	if err != nil || i == 0 {
		log.Println(err)
		i = constants.IntervalService
	}
	srv.IntervalService = i

	i, err = strconv.Atoi(data.Settings.IntervalWebClient)
	if err != nil || i == 0 {
		log.Println(err)
		i = 0
	}
	srv.IntervalWebClient = i

	srv.ControlServer = data.Settings.ControlServer

	i, err = strconv.Atoi(data.Settings.KillProcessKb)
	if err != nil || i == 0 {
		log.Println(err)
		i = constants.KillProcessKb
	}
	srv.KillProcessKb = i

	for _, v := range data.BasesDoubleControl {
		c := repository.BasesDoubleControl{
			Server:   v.Server,
			Name:     v.Name,
			User:     v.User,
			Password: v.Password,
		}

		bdc = append(bdc, c)
	}

	srv.MessageCom = bdc
	srv.ClosedTasks = []repository.ClosedTask{}
	srv.ClosedConnects = []repository.ClosedConnect{}
	srv.ExchangeStorage = exchange.ExchangeStorage{}

	srv.ClaimsStore = []token.ClaimStore{}

	//srv.BoltDB = dbbolt.GetDB()
}

func (srv *Server) InitData() {
	path := constants.Pudge
	entries, err := os.ReadDir(path)
	if err != nil {
		return
	}
	for _, e := range entries {
		if !e.IsDir() {
			continue
		}

		eq := exchange.ExchangeQueue{
			Name:               e.Name(),
			UploadAt:           time.Now(),
			DownloadAt:         time.Now(),
			TypeStorage:        constants.Hard,
			TimeTransferToDisk: constants.TimeTransferToDisk,
			MaxSizeSoftMode:    constants.TotalByteInPudge,
			MaxLenSoftMode:     constants.TotalCount,
			PriorityMessages:   exchange.PriorityMessages{},
		}

		srv.ExchangeStorage[e.Name()] = eq
	}
}

// Run Запуск сервера
func (srv *Server) Run() {

	ctx, cancelFunc := context.WithCancel(context.Background())
	if srv.IntervalService != 0 || srv.KillProcessKb != 0 {
		go srv.serviceKillWinProc(ctx, cancelFunc)
	}
	go srv.serviceResetDataDisk(ctx, cancelFunc)

	//srv.PrivateKey = ServicePrivateKey()

	if srv.IntervalService != 0 || srv.IntervalWebClient != 0 {
		var p uintptr
		const coin uint32 = 0
		err := ole.CoInitializeEx(p, coin)
		if err == nil {
			log.Println("start control DD ok!")
			go srv.serviceDropDouble(ctx, cancelFunc)
		}
	}

	hostname, err := os.Hostname()
	if err != nil {
		hostname = "localhost"
	}

	port := constants.Port
	if srv.HTTPPort != "" {
		port = srv.HTTPPort
	}
	addressString := fmt.Sprintf("%s:%s", hostname, port)

	go func() {
		s := &http.Server{
			Addr:    addressString,
			Handler: srv.Router}

		log.Println(addressString)
		if err := s.ListenAndServe(); err != nil {
			log.Fatalln(err)
		}
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGTERM, syscall.SIGINT, syscall.SIGQUIT)
	<-stop

	srv.Stop()
}

func (srv *Server) Stop() {
	//for k, v := range srv.ExchangeStorage {
	//
	//	srv.Lock()
	//
	//	_ = v.TransferToDisk()
	//	srv.ExchangeStorage[k] = v
	//
	//	srv.Unlock()
	//
	//}
}
