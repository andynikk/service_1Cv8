package repository

import (
	"Service_1Cv8/internal/repository/datainterface"
	"bytes"
	"os/exec"
	"strings"
	"time"
)

type MSettings map[string]datainterface.Setting

type DataJSON struct {
	DataDBJSON
	Settings
	Services           []Services
	BasesDoubleControl []BasesDoubleControl
	DatabaseTokens     []DatabaseTokens
}

type BasesDoubleControl struct {
	Server   string `yaml:"server"`
	Name     string `yaml:"name"`
	User     string `yaml:"user"`
	Password string `yaml:"password"`
	UID      string `yaml:"uid"`
}

type DatabaseTokens struct {
	UID          string `yaml:"uid"`
	Name         string `yaml:"name"`
	Server       string `yaml:"server"`
	NameOnServer string `yaml:"name_on_server"`
	Port         string `yaml:"port"`
	PublicKey    string `yaml:"public_key"`
	PrivateKey   string `yaml:"private_key"`
}

type DataBases struct {
	Name         string `yaml:"name"`
	Server       string `yaml:"server"`
	NameOnServer string `yaml:"name_on_server"`
	Port         string `yaml:"port"`
	UID          string `yaml:"uid"`
}

type PropertyDB struct {
	NameUser     string `yaml:"name_user"`
	PasswordUser string `yaml:"password_user"`
	StartBlock   string `yaml:"start_block"`
	FinishBlock  string `yaml:"finish_block"`
	KeyUnlock    string `yaml:"key_unlock"`
	Massage      string `yaml:"massage"`
	Block        bool   `yaml:"block"`
	UID          string `yaml:"uid"`
}

type SQLServer struct {
	Description string `yaml:"description"`
	Name        string `yaml:"name"`
	User        string `yaml:"user"`
	Password    string `yaml:"password"`
}

type Services struct {
	NameServer  string      `yaml:"name_server"`
	IP          string      `yaml:"ip"`
	NameService string      `yaml:"name_service"`
	UID         string      `yaml:"uid"`
	User        string      `yaml:"user"`
	Password    string      `yaml:"password"`
	SQLServers  []SQLServer `yaml:"sql_servers"`
}

type ClosedConnect struct {
	Time         time.Time `yaml:"time"`
	TimeStartCon time.Time `yaml:"time_start_con"`
	Host         string    `yaml:"host"`
	DB           string    `yaml:"db"`
	AppID        string    `yaml:"app_id"`
	User         string    `yaml:"user"`
}

type ClosedTask struct {
	Time time.Time `yaml:"time"`
	Host string    `yaml:"host"`
	Size string    `yaml:"size"`
}

func NewStorage() *DataJSON {

	dataBases := []DataBases{}
	propertyDB := []PropertyDB{}
	services := []Services{}

	dataDBJSON := DataDBJSON{
		DB:         dataBases,
		PropertyDB: propertyDB,
	}

	d := DataJSON{
		DataDBJSON: dataDBJSON,
		Services:   services,
		Settings:   Settings{},
	}
	d.GetPudgeData()
	//if err != nil {
	//	return &d
	//}

	return &d
}

func IsOSWindows() bool {

	var stderr bytes.Buffer
	defer stderr.Reset()

	var out bytes.Buffer
	defer out.Reset()

	cmd := exec.Command("cmd", "ver")
	cmd.Stdin = strings.NewReader("some input")
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	err := cmd.Run()
	if err != nil {
		return false
	}
	myOS := out.String()
	return strings.Contains(myOS, "Microsoft Windows")
}

func GetPropertiesDB(pryDB []PropertyDB, uid string) (PropertyDB, int) {
	for k, v := range pryDB {
		if v.UID == uid {
			return v, k
		}
	}

	return PropertyDB{}, -1
}

func GetDB(db []DataBases, uid string) (DataBases, int) {
	for k, v := range db {
		if v.UID == uid {
			return v, k
		}
	}

	return DataBases{}, -1
}

func GetService(services []Services, uid string) (Services, int) {
	for k, v := range services {
		if v.UID == uid {

			return v, k

		}
	}

	return Services{}, -1
}

func GetBasesDoubleControlDB(db []BasesDoubleControl, uid string) (BasesDoubleControl, int) {
	for k, v := range db {
		if v.UID == uid {
			return v, k
		}
	}

	return BasesDoubleControl{}, -1
}
