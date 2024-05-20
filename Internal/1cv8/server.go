package OneCv8

import (
	"fmt"
	"log"
	"sort"
	"strings"
	"time"

	"Service_1Cv8/internal/repository"

	"github.com/go-ole/go-ole"
	"github.com/go-ole/go-ole/oleutil"
)

type SettingCloseDoubleConnection struct {
	IntervalWebClient int
	OutputMessages    bool
}

type MassageJSON struct {
	NameServer     string `yaml:"name_server"`
	NameDB         string `yaml:"name_db"`
	NameUser       string `yaml:"name_user"`
	PasswordUser   string `yaml:"password_user"`
	Block          bool   `yaml:"block"`
	PermissionCode string `yaml:"permission_code"`
	DeniedMessage  string `yaml:"denied_message"`
	DeniedFrom     string `yaml:"denied_from"`
	DeniedTo       string `yaml:"denied_to"`
}

type Clasters struct {
	HostName    string
	ClusterName string
	MainPort    int
}

type InfoBases struct {
	ConnectDenied  string
	SessionsDenied string
	PermissionCode string
	DeniedMessage  string
	DeniedFrom     string
	DeniedTo       string
	UpdateInfoBase string
	DBMS           string
	DBName         string
	DBServerName   string
	DBUser         string
	Name           string
}

type ItemDB struct {
	Descr    string
	Name     string
	MainPort string
	User     string
	Password string
	UID      string
}

func PropertyDB(massageJSON MassageJSON) error {

	ole.CoInitialize(0)

	Server1С, err := oleutil.CreateObject("V83.COMConnector")
	if err != nil {
		return err
	}
	wmi, err := Server1С.QueryInterface(ole.IID_IDispatch)
	if err != nil {
		return err
	}
	arg := []interface{}{massageJSON.NameServer}
	agent, err := oleutil.CallMethod(wmi, "ConnectAgent", arg...)
	if err != nil {
		return err
	}
	ServerAgent := agent.ToIDispatch()
	defer ServerAgent.Release()
	arg = []interface{}{}
	Clasters, err := oleutil.CallMethod(ServerAgent, "GetClusters", arg...)
	if err != nil {
		return err
	}
	resClasters := Clasters.ToArray()
	defer resClasters.Release()

	for _, Claster := range resClasters.ToValueArray() {

		arg = []interface{}{Claster, "", ""}
		_, err = oleutil.CallMethod(ServerAgent, "Authenticate", arg...)
		if err != nil {
			continue
		}

		arg = []interface{}{Claster}
		wps, err := oleutil.CallMethod(ServerAgent, "GetWorkingProcesses", arg...)
		if err != nil {
			continue
		}

		WorkingProcesses := wps.ToArray()
		defer WorkingProcesses.Release()
		for _, wp := range WorkingProcesses.ToValueArray() {
			WorkingProcess := wp.(*ole.IDispatch)
			r, _ := oleutil.GetProperty(WorkingProcess, "Running")
			running := r.Value()
			if running == 0 {
				continue
			}

			hn, _ := oleutil.GetProperty(WorkingProcess, "HostName")
			mp, _ := oleutil.GetProperty(WorkingProcess, "MainPort")

			HostName := hn.Value()
			MainPort := mp.Value()

			CWPAddr := fmt.Sprintf("tcp://%s:%d", HostName, MainPort)

			arg = []interface{}{CWPAddr}
			CWP, err := oleutil.CallMethod(wmi, "ConnectWorkingProcess", arg...)
			if err != nil {
				return err
			}

			oleCWP := CWP.ToIDispatch()
			arg = []interface{}{}
			ibs, err := oleutil.CallMethod(oleCWP, "GetInfoBases", arg...)

			InfoBases := ibs.ToArray()
			for _, ib := range InfoBases.ToValueArray() {
				InfoBase := ib.(*ole.IDispatch)

				ndb, _ := oleutil.GetProperty(InfoBase, "Name")
				nameDB := ndb.Value()
				if nameDB != massageJSON.NameDB {
					continue
				}

				arg = []interface{}{massageJSON.NameUser, massageJSON.PasswordUser}
				_, err = oleutil.CallMethod(oleCWP, "AddAuthentication", arg...)
				if err != nil {
					return err
				}

				arg = []interface{}{massageJSON.Block}
				_, err = oleutil.PutProperty(InfoBase, "ConnectDenied", arg...)
				if err != nil {
					return err
				}
				arg = []interface{}{massageJSON.Block}
				_, err = oleutil.PutProperty(InfoBase, "SessionsDenied", arg...)
				if err != nil {
					return err
				}
				arg = []interface{}{massageJSON.PermissionCode}
				_, err = oleutil.PutProperty(InfoBase, "PermissionCode", arg...)
				if err != nil {
					return err
				}
				arg = []interface{}{massageJSON.DeniedMessage}
				_, err = oleutil.PutProperty(InfoBase, "DeniedMessage", arg...)
				if err != nil {
					return err
				}
				arg = []interface{}{massageJSON.DeniedFrom}
				_, err = oleutil.PutProperty(InfoBase, "DeniedFrom", arg...)
				if err != nil {
					return err
				}
				arg = []interface{}{massageJSON.DeniedTo}
				_, err = oleutil.PutProperty(InfoBase, "DeniedTo", arg...)
				if err != nil {
					return err
				}

				arg = []interface{}{InfoBase}
				_, err = oleutil.CallMethod(oleCWP, "UpdateInfoBase", arg...)
				if err != nil {
					return err
				}
				break
			}
			break
		}
	}

	return nil
}

func CheckPropertyDB(massageJSON MassageJSON) (bool, error) {

	ole.CoInitialize(0)

	Server1С, err := oleutil.CreateObject("V83.COMConnector")
	if err != nil {
		return false, err
	}
	defer Server1С.Release()

	wmi, err := Server1С.QueryInterface(ole.IID_IDispatch)
	if err != nil {
		return false, err
	}

	//arg1 := []interface{}{"Srvr='1c.telematika.local:1741';Ref='tsk_buh';Usr='Михайлов_АН';Pwd='101650'"}
	//connDB, err := oleutil.CallMethod(wmi, "Connect", arg1...)
	//_ = connDB
	//massageJSON.NameServer = "1c1.telematika.local"
	arg := []interface{}{massageJSON.NameServer}
	agent, err := oleutil.CallMethod(wmi, "ConnectAgent", arg...)

	if err != nil {
		return false, err
	}
	ServerAgent := agent.ToIDispatch()
	defer ServerAgent.Release()

	arg = []interface{}{}
	Clasters, err := oleutil.CallMethod(ServerAgent, "GetClusters", arg...)
	if err != nil {
		return false, err
	}

	resClasters := Clasters.ToArray()
	defer resClasters.Release()

	for _, Claster := range resClasters.ToValueArray() {

		arg = []interface{}{Claster, "", ""}
		_, err = oleutil.CallMethod(ServerAgent, "Authenticate", arg...)
		if err != nil {
			continue
		}

		arg = []interface{}{Claster}
		wps, err := oleutil.CallMethod(ServerAgent, "GetWorkingProcesses", arg...)
		if err != nil {
			continue
		}

		WorkingProcesses := wps.ToArray()
		defer WorkingProcesses.Release()

		result := false
		for _, wp := range WorkingProcesses.ToValueArray() {
			WorkingProcess := wp.(*ole.IDispatch)
			r, _ := oleutil.GetProperty(WorkingProcess, "Running")
			running := r.Value()
			if running == 0 {
				continue
			}

			hn, _ := oleutil.GetProperty(WorkingProcess, "HostName")
			mp, _ := oleutil.GetProperty(WorkingProcess, "MainPort")

			HostName := hn.Value()
			MainPort := mp.Value()

			CWPAddr := fmt.Sprintf("tcp://%s:%d", HostName, MainPort)

			arg = []interface{}{CWPAddr}
			CWP, err := oleutil.CallMethod(wmi, "ConnectWorkingProcess", arg...)
			if err != nil {
				return false, err
			}

			oleCWP := CWP.ToIDispatch()

			arg = []interface{}{massageJSON.NameUser, massageJSON.PasswordUser}
			_, err = oleutil.CallMethod(oleCWP, "AddAuthentication", arg...)
			if err != nil {
				return false, err
			}

			arg = []interface{}{}
			ibs, err := oleutil.CallMethod(oleCWP, "GetInfoBases", arg...)

			InfoBases := ibs.ToArray()
			for _, ib := range InfoBases.ToValueArray() {
				InfoBase := ib.(*ole.IDispatch)

				ndb, _ := oleutil.GetProperty(InfoBase, "Name")
				if ndb == nil {
					continue
				}
				nameDB := ndb.Value()
				if nameDB != massageJSON.NameDB {
					continue
				}

				cdDB, _ := oleutil.GetProperty(InfoBase, "ConnectDenied")
				connectDeniedDB := cdDB.Value()

				oleCWP.Release()

				result = connectDeniedDB.(bool)
				break
			}

			arg = []interface{}{Claster}
			ibs, err = oleutil.CallMethod(ServerAgent, "GetInfoBases", arg...)
			if err != nil {
				return false, err
			}
			InfoBasesDrop := ibs.ToArray()
			for _, ib := range InfoBasesDrop.ToValueArray() {
				InfoBase := ib.(*ole.IDispatch)

				ndb, _ := oleutil.GetProperty(InfoBase, "Name")
				nameDB := ndb.Value()
				if nameDB != massageJSON.NameDB {
					continue
				}

				arg = []interface{}{Claster, InfoBase}
				csdb, err := oleutil.CallMethod(ServerAgent, "GetInfoBaseSessions", arg...)
				if err != nil {
					return false, err
				}

				ConnectsDB := csdb.ToArray()
				for _, cdb := range ConnectsDB.ToValueArray() {
					connect := cdb.(*ole.IDispatch)

					aID, _ := oleutil.GetProperty(connect, "AppID")
					appID := aID.Value().(string)

					if strings.ToLower(appID) != "comconsole" {
						continue
					}

					arg = []interface{}{Claster, connect}
					_, err = oleutil.CallMethod(ServerAgent, "TerminateSession", arg...)
					if err != nil {
						return false, err
					}

				}
				return result, nil
			}

		}
	}

	return false, nil
}

func DropUsersDB(massageJSON MassageJSON) (bool, error) {

	res := false
	ole.CoInitialize(0)

	Server1С, err := oleutil.CreateObject("V83.COMConnector")
	if err != nil {
		return false, err
	}
	wmi, err := Server1С.QueryInterface(ole.IID_IDispatch)
	if err != nil {
		return false, err
	}
	arg := []interface{}{massageJSON.NameServer}
	agent, err := oleutil.CallMethod(wmi, "ConnectAgent", arg...)
	if err != nil {
		return false, err
	}
	ServerAgent := agent.ToIDispatch()
	defer ServerAgent.Release()
	arg = []interface{}{}
	Clasters, err := oleutil.CallMethod(ServerAgent, "GetClusters", arg...)
	if err != nil {
		return false, err
	}
	resClasters := Clasters.ToArray()
	defer resClasters.Release()

	for _, Claster := range resClasters.ToValueArray() {

		arg = []interface{}{Claster, "", ""}
		_, err = oleutil.CallMethod(ServerAgent, "Authenticate", arg...)
		if err != nil {
			continue
		}

		arg = []interface{}{Claster}
		wps, err := oleutil.CallMethod(ServerAgent, "GetWorkingProcesses", arg...)
		if err != nil {
			continue
		}

		WorkingProcesses := wps.ToArray()
		defer WorkingProcesses.Release()
		for _, wp := range WorkingProcesses.ToValueArray() {
			WorkingProcess := wp.(*ole.IDispatch)
			r, _ := oleutil.GetProperty(WorkingProcess, "Running")
			running := r.Value()
			if running == 0 {
				continue
			}

			///////////////////////////////////////////////////////////////////

			arg = []interface{}{Claster}
			ibs, err := oleutil.CallMethod(ServerAgent, "GetInfoBases", arg...)
			if err != nil {
				return false, err
			}
			InfoBases := ibs.ToArray()
			for _, ib := range InfoBases.ToValueArray() {
				InfoBase := ib.(*ole.IDispatch)

				ndb, _ := oleutil.GetProperty(InfoBase, "Name")
				nameDB := ndb.Value()
				if nameDB != massageJSON.NameDB {
					continue
				}

				arg = []interface{}{Claster, InfoBase}
				csdb, err := oleutil.CallMethod(ServerAgent, "GetInfoBaseSessions", arg...)
				if err != nil {
					return false, err
				}

				ConnectsDB := csdb.ToArray()
				for _, cdb := range ConnectsDB.ToValueArray() {
					connect := cdb.(*ole.IDispatch)

					aID, _ := oleutil.GetProperty(connect, "AppID")
					appID := aID.Value().(string)

					uName, _ := oleutil.GetProperty(connect, "userName")
					userName := uName.Value().(string)

					if strings.ToLower(appID) == "designer" ||
						(strings.ToLower(appID) == "httpserviceconnection" && userName == massageJSON.NameUser) {
						continue
					}

					arg = []interface{}{Claster, connect}
					_, err = oleutil.CallMethod(ServerAgent, "TerminateSession", arg...)
					if err != nil {
						return false, err
					}
				}
				res = true
				break
			}
			break
		}
	}

	return res, nil

	//msgJSON, err := json.MarshalIndent(massageJSON, "", " ")
	//
	//body := bytes.NewReader(msgJSON)
	//addressPost := "http://tk-test-app/mikhailov_uh/hs/DataExchangeUH/Exchange/drobusers"
	//req, err := http.NewRequest("POST", addressPost, body)
	//if err != nil {
	//	return false, errors.New("-- ошибка отправки данных на сервер")
	//}
	//req.SetBasicAuth(massageJSON.NameUser, massageJSON.PasswordUser)
	////u, p, ok := req.BasicAuth()
	//
	//req.Header.Set("Content-Type", "application/json")
	//defer req.Body.Close()
	//
	//client := &http.Client{}
	//resp, err := client.Do(req)
	//if err != nil {
	//	return false, errors.New("-- ошибка отправки данных на сервер")
	//}
	//defer resp.Body.Close()
	//
	//if resp.StatusCode != http.StatusOK {
	//	return false, errors.New("-- ошибка обновления данных на сервере")
	//}
	//
	//b := result{}
	//bodyJSON, err := io.ReadAll(resp.Body)
	//
	//err = json.Unmarshal(bodyJSON, &b)
	//if err != nil {
	//	return false, err
	//}
	//
	//return b.Result, nil
}

type Connect struct {
	Host     string
	AppID    string
	UserName string

	SesConnects []SesConnect
}

type SesConnect struct {
	SessionID int32
	StartedAt time.Time
	Cnn       *ole.IDispatch
}

func DropDoubleUsersDB(massagesJSON []MassageJSON, scdc SettingCloseDoubleConnection, c chan repository.ClosedConnect) {
	defer close(c)

	if scdc.OutputMessages {
		log.Println("start kill DD")
	}
	if len(massagesJSON) == 0 {
		if scdc.OutputMessages {
			log.Println("Control BD 0. finish kill DD")
		}
		return
	}

	Server1С, err := oleutil.CreateObject("V83.COMConnector")
	if err != nil {
		if scdc.OutputMessages {
			log.Println(2, err.Error())
		}
		return
	}
	wmi, err := Server1С.QueryInterface(ole.IID_IDispatch)
	if err != nil {
		if scdc.OutputMessages {
			log.Println(3, err.Error())
		}
		return
	}

	allClosedDD := 0
	for _, massageJSON := range massagesJSON {

		arg := []interface{}{massageJSON.NameServer}
		agent, err := oleutil.CallMethod(wmi, "ConnectAgent", arg...)
		if err != nil {
			if scdc.OutputMessages {
				log.Println(4, err.Error())
			}
			return
		}
		ServerAgent := agent.ToIDispatch()
		defer ServerAgent.Release()
		arg = []interface{}{}
		Clasters, err := oleutil.CallMethod(ServerAgent, "GetClusters", arg...)
		if err != nil {
			if scdc.OutputMessages {
				log.Println(5, err.Error())
			}
			return
		}
		resClasters := Clasters.ToArray()
		defer resClasters.Release()

		var Connects []Connect
		var ConnectsWC []Connect

		for _, Claster := range resClasters.ToValueArray() {
			arg = []interface{}{Claster, "", ""}
			_, err = oleutil.CallMethod(ServerAgent, "Authenticate", arg...)
			if err != nil {
				if scdc.OutputMessages {
					log.Println(6, err.Error())
				}
				continue
			}

			arg = []interface{}{Claster}
			wps, err := oleutil.CallMethod(ServerAgent, "GetWorkingProcesses", arg...)
			if err != nil {
				if scdc.OutputMessages {
					log.Println(7, err.Error())
				}
				continue
			}

			WorkingProcesses := wps.ToArray()
			defer WorkingProcesses.Release()

			for _, wp := range WorkingProcesses.ToValueArray() {
				WorkingProcess := wp.(*ole.IDispatch)
				r, err := oleutil.GetProperty(WorkingProcess, "Running")
				if err != nil {
					if scdc.OutputMessages {
						log.Println(8, err.Error())
					}
					continue
				}
				running := r.Value()
				if running == 0 {
					continue
				}

				///////////////////////////////////////////////////////////////////

				arg = []interface{}{Claster}
				ibs, err := oleutil.CallMethod(ServerAgent, "GetInfoBases", arg...)
				if err != nil {
					if scdc.OutputMessages {
						log.Println(9, err.Error())
					}
					continue
				}
				InfoBases := ibs.ToArray()
				for _, ib := range InfoBases.ToValueArray() {
					InfoBase := ib.(*ole.IDispatch)

					ndb, err := oleutil.GetProperty(InfoBase, "Name")
					if err != nil {
						if scdc.OutputMessages {
							log.Println(10, err.Error())
						}
						continue
					}
					nameDB := ndb.Value()
					if nameDB != massageJSON.NameDB {
						continue
					}

					arg = []interface{}{Claster, InfoBase}
					csdb, err := oleutil.CallMethod(ServerAgent, "GetInfoBaseSessions", arg...)
					if err != nil {
						if scdc.OutputMessages {
							log.Println(11, err.Error())
						}
						return
					}

					ConnectsDB := csdb.ToArray()
					for _, cdb := range ConnectsDB.ToValueArray() {
						connect := cdb.(*ole.IDispatch)

						h, err := oleutil.GetProperty(connect, "Host")
						if err != nil {
							if scdc.OutputMessages {
								log.Println(12, err.Error())
							}
							continue
						}
						host := h.Value().(string)

						aID, err := oleutil.GetProperty(connect, "AppID")
						if err != nil {
							if scdc.OutputMessages {
								log.Println(13, err.Error())
							}
							continue
						}
						appID := aID.Value().(string)

						if appID == "BackgroundJob" {
							continue
						}

						uName, err := oleutil.GetProperty(connect, "userName")
						if err != nil {
							continue
						}
						userName := uName.Value().(string)

						sAt, err := oleutil.GetProperty(connect, "StartedAt")
						if err != nil {
							if scdc.OutputMessages {
								log.Println(14, err.Error())
							}
							continue
						}
						startedAt := sAt.Value().(time.Time)

						//sLastAt, err := oleutil.GetProperty(connect, "LastActiveAt")
						//if err != nil {
						//	if scdc.OutputMessages {
						//		log.Println(14, err.Error())
						//	}
						//	continue
						//}
						//lastStartedAt := sLastAt.Value().(time.Time)

						if appID == "WebClient" && scdc.IntervalWebClient != 0 {

							var sc []SesConnect
							sc = append(sc, SesConnect{SessionID: 0, StartedAt: startedAt, Cnn: connect})
							c := Connect{
								Host:        host,
								AppID:       appID,
								UserName:    userName,
								SesConnects: sc,
							}
							ConnectsWC = append(ConnectsWC, c)
						}

						k := findAUserID(Connects, userName, host, appID)
						if k == -1 {
							c := Connect{
								Host:        host,
								AppID:       appID,
								UserName:    userName,
								SesConnects: []SesConnect{},
							}
							Connects = append(Connects, c)
							k = len(Connects) - 1
						}

						Connects[k].SesConnects = append(Connects[k].SesConnects,
							SesConnect{SessionID: 0, StartedAt: startedAt, Cnn: connect})
					}

					for _, v := range Connects {
						if len(v.SesConnects) <= 1 {
							continue
						}

						sort.Slice(v.SesConnects, func(i, j int) bool {
							return v.SesConnects[i].StartedAt.Before(v.SesConnects[j].StartedAt)
						})

						for k, val := range v.SesConnects {
							if k == len(v.SesConnects)-1 {
								break
							}

							arg = []interface{}{Claster, val.Cnn}
							_, err = oleutil.CallMethod(ServerAgent, "TerminateSession", arg...)
							if err != nil {
								if scdc.OutputMessages {
									log.Println("15.2", err.Error())
								}
								continue
							}

							cc := repository.ClosedConnect{
								Time:         time.Now(),
								TimeStartCon: val.StartedAt,
								Host:         v.Host,
								DB:           massageJSON.NameDB,
								User:         v.UserName,
								AppID:        v.AppID,
							}

							c <- cc
							allClosedDD++
						}
					}

					for _, v := range ConnectsWC {
						for _, val := range v.SesConnects {

							StartedAtPlus := val.StartedAt.Add(time.Minute * time.Duration(scdc.IntervalWebClient))
							StartedAtPlusLocal := time.Date(StartedAtPlus.Year(), StartedAtPlus.Month(), StartedAtPlus.Day(),
								StartedAtPlus.Hour(), StartedAtPlus.Minute(), StartedAtPlus.Second(),
								StartedAtPlus.Nanosecond(), time.Local)

							now := time.Now()
							if StartedAtPlusLocal.After(now) {
								continue
							}

							arg = []interface{}{Claster, val.Cnn}
							_, err = oleutil.CallMethod(ServerAgent, "TerminateSession", arg...)
							if err != nil {
								if scdc.OutputMessages {
									log.Println("15.1", err.Error())
								}
								continue
							}

							cc := repository.ClosedConnect{
								Time:         time.Now(),
								TimeStartCon: val.StartedAt,
								Host:         v.Host,
								DB:           massageJSON.NameDB,
								User:         v.UserName,
								AppID:        v.AppID,
							}

							c <- cc
							allClosedDD++
						}
					}

					break
				}
				break
			}
		}
	}

	if scdc.OutputMessages {
		log.Println(allClosedDD, "finish kill DD")
	}
}

func findAUserID(cnn []Connect, userName, host, appID string) int {

	for k, v := range cnn {
		if v.UserName == userName && v.Host == host && v.AppID == appID {
			return k
		}
	}

	return -1
}

func ListDB(massageJSON MassageJSON) ([]ItemDB, error) {

	var itemsDB []ItemDB

	ole.CoInitialize(0)

	Server1С, err := oleutil.CreateObject("V83.COMConnector")
	if err != nil {
		return nil, err
	}
	wmi, err := Server1С.QueryInterface(ole.IID_IDispatch)
	if err != nil {
		return nil, err
	}
	arg := []interface{}{massageJSON.NameServer}
	agent, err := oleutil.CallMethod(wmi, "ConnectAgent", arg...)
	if err != nil {
		return nil, err
	}
	ServerAgent := agent.ToIDispatch()
	defer ServerAgent.Release()
	arg = []interface{}{}
	Clasters, err := oleutil.CallMethod(ServerAgent, "GetClusters", arg...)
	if err != nil {
		return nil, err
	}
	resClasters := Clasters.ToArray()
	defer resClasters.Release()

	for _, Claster := range resClasters.ToValueArray() {
		cl := Claster.(*ole.IDispatch)

		arg = []interface{}{Claster, "", ""}
		_, err = oleutil.CallMethod(ServerAgent, "Authenticate", arg...)
		if err != nil {
			continue
		}

		mainPort, _ := oleutil.GetProperty(cl, "MainPort")

		arg = []interface{}{Claster}
		ibs, err := oleutil.CallMethod(ServerAgent, "GetInfoBases", arg...)
		if err != nil {
			continue
		}
		infoBases := ibs.ToArray()
		for _, v := range infoBases.ToValueArray() {
			infoBases := v.(*ole.IDispatch)

			descr, _ := oleutil.GetProperty(infoBases, "Descr")
			name, _ := oleutil.GetProperty(infoBases, "Name")

			itemDB := ItemDB{
				Descr:    descr.ToString(),
				Name:     name.ToString(),
				MainPort: fmt.Sprintf("%d", mainPort.Value().(int32)),
			}

			itemsDB = append(itemsDB, itemDB)
		}
	}

	return itemsDB, nil
}
