package handlers

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	OneCv8 "Service_1Cv8/internal/1cv8"
	"Service_1Cv8/internal/constants"
	"Service_1Cv8/internal/repository"
	"Service_1Cv8/internal/winsys"

	"github.com/gorilla/mux"
)

type Setting struct {
	HTTPPort        string    `yaml:"http_port"`
	IntervalService int       `yaml:"interval_service"`
	KillProcessKb   int       `yaml:"kill_process_kb"`
	ExtremeStart    time.Time `yaml:"extreme_start"`
	ControlServer   string    `yaml:"control_server"`
	MessageCom      []repository.BasesDoubleControl
	ClosedConnects  []repository.ClosedConnect
	ClosedTasks     []repository.ClosedTask
}

type Server struct {
	*mux.Router
	Setting
}

// NewServer создание сервера
func NewServer() *Server {
	srv := &Server{}

	srv.InitRouters()
	srv.InitDB()

	return srv
}

func (srv *Server) InitRouters() {
	r := mux.NewRouter()

	r.HandleFunc("/", srv.handleFunc).Methods("GET")
	r.HandleFunc("/u", srv.handleFuncUpdateDB).Methods("GET")

	srv.Router = r
}

func (srv *Server) InitDB() {
	var bdc []repository.BasesDoubleControl

	data := repository.DataJSON{}
	if err := data.GetYamlData(); err != nil {
		log.Println(err)
	}
	srv.HTTPPort = data.Settings.HTTPPort

	i, err := strconv.Atoi(data.Settings.IntervalService)
	if err != nil || i == 0 {
		log.Println(err)
		i = constants.IntervalService
	}
	srv.IntervalService = i

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
}

// Run Запуск сервера
func (srv *Server) Run() {
	ctx, cancelFunc := context.WithCancel(context.Background())
	go srv.Service(ctx, cancelFunc)

	port := constants.Port
	if srv.HTTPPort != "" {
		port = srv.HTTPPort
	}
	go func() {
		s := &http.Server{
			Addr:    fmt.Sprintf("localhost:%s", port),
			Handler: srv.Router}

		log.Println(port)
		if err := s.ListenAndServe(); err != nil {
			log.Fatalln(err)
		}
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGTERM, syscall.SIGINT, syscall.SIGQUIT)
	<-stop
}

func (srv *Server) Service(ctx context.Context, cancelFunc context.CancelFunc) {
	ticker := time.NewTicker(time.Duration(srv.IntervalService) * time.Minute)
	for {
		select {
		case <-ticker.C:

			srv.ExtremeStart = time.Now()

			closedConnects, err := winsys.KillWinProc(srv.ControlServer, "rphost.exe", srv.KillProcessKb)
			if err != nil {
				for _, v := range closedConnects {
					srv.Setting.ClosedTasks = append(srv.Setting.ClosedTasks, v)
				}
			}

			for _, v := range srv.MessageCom {
				if v.Server == "" || v.Name == "" {
					continue
				}

				m := OneCv8.MassageJSON{
					NameServer:   v.Server,
					NameDB:       v.Name,
					NameUser:     v.User,
					PasswordUser: v.Password,
				}

				closedConnects, err := OneCv8.DropDoubleUsersDB(m)
				if err != nil {
					log.Println(err)
					continue
				}
				for _, v := range closedConnects {
					srv.Setting.ClosedConnects = append(srv.Setting.ClosedConnects, v)
				}
			}

		case <-ctx.Done():
			cancelFunc()
			return
		}
	}
}
