package cli

import (
	OneCv8 "Service_1Cv8/internal/1cv8"
	"Service_1Cv8/internal/constants"
	"fmt"
	"log"
	"sync"

	"Service_1Cv8/internal/environment"
	"Service_1Cv8/internal/repository"

	tea "github.com/charmbracelet/bubbletea"
)

var client = Client{}

type KeyContext string

type CustomModel interface {
	SetParameters([]interface{})
}

type EventsService struct {
	Stop   bool
	Del    bool
	Start  bool
	Reboot bool
}

type StringArray []string

type PerformedActions struct {
	UpdateDB       StringArray
	ClearCash      StringArray
	RestartService StringArray
	RebutServer    StringArray
}

type CM CustomModel

type Client struct {
	Config        *environment.Config
	Models        []CM
	Tea           *tea.Program
	EventsService map[string]EventsService
	PerformedActions
	sync.RWMutex
	Storage *repository.DataJSON
}

func NewClient() Client {

	fmt.Print("loading...")
	client.Storage = repository.NewStorage()
	for i, pDB := range client.Storage.PropertyDB {

		db, k := repository.GetDB(client.Storage.DB, pDB.UID)
		if k == -1 {
			continue
		}

		massageJSON := OneCv8.MassageJSON{
			NameServer:   db.Server,
			NameDB:       db.NameOnServer,
			NameUser:     pDB.NameUser,
			PasswordUser: pDB.PasswordUser,
		}

		r, _ := OneCv8.CheckPropertyDB(massageJSON)
		client.Storage.PropertyDB[i].Block = r
	}

	client.Models = []CM{
		NewFormMain(),
		NewFormAddDB(),
		NewFormPropertyDB(repository.DataBases{}, repository.PropertyDB{}),
		NewFormServer(),
		NewFormLogin(),
		NewFormSetting(),
		NewFormAddService(),
		NewFormServerDB(),
		NewFormSQLDB(),
		NewFormAddSQLSrv(),
		NewFormCopyFile(),
		NewFormResTlg(),
		NewFormTgMsg(),
	}

	//frm := client.Models[constants.FormLogin].(*FormLogin)
	//frm.SetParameters(nil)

	frm := client.Models[constants.FormMain].(*FormMain)
	frm.SetParameters(nil)

	client.Tea = tea.NewProgram(frm)
	config, err := environment.NewConfig()
	if err != nil {
		log.Fatal(err.Error())
	}

	client.Config = config
	return client
}

func (c *Client) Run() error {
	_, err := c.Tea.Run()
	return err
}

func (c *Client) Shutdown() error {
	return c.Storage.SetPudgelData()
}

func (c *Client) EditDB(uid string, db repository.DataBases, propertyDB repository.PropertyDB) error {

	editDB, keyDB := repository.GetDB(c.Storage.DB, uid)
	if keyDB == -1 {
		editDB = repository.DataBases{}
	}
	editDB.Name = db.Name
	editDB.Server = db.Server
	editDB.Port = db.Port
	editDB.NameOnServer = db.NameOnServer
	editDB.UID = uid

	editPropertyDB, keyPDB := repository.GetPropertiesDB(client.Storage.PropertyDB, uid)
	if keyPDB == -1 {
		editPropertyDB = repository.PropertyDB{}
	}
	editPropertyDB.StartBlock = propertyDB.StartBlock
	editPropertyDB.FinishBlock = propertyDB.FinishBlock
	editPropertyDB.KeyUnlock = propertyDB.KeyUnlock
	editPropertyDB.NameUser = propertyDB.NameUser
	editPropertyDB.PasswordUser = propertyDB.PasswordUser
	editPropertyDB.Massage = propertyDB.Massage
	editPropertyDB.UID = uid

	dataDBJSON := &client.Storage.DataDBJSON
	if keyDB != -1 {
		dataDBJSON.DB[keyDB] = dataDBJSON.DB[len(dataDBJSON.DB)-1]
		dataDBJSON.DB[len(dataDBJSON.DB)-1] = repository.DataBases{}
		dataDBJSON.DB = dataDBJSON.DB[:len(dataDBJSON.DB)-1]
	}
	dataDBJSON.DB = append(dataDBJSON.DB, db)

	if keyPDB != -1 {
		dataDBJSON.PropertyDB[keyPDB] = dataDBJSON.PropertyDB[len(dataDBJSON.PropertyDB)-1]
		dataDBJSON.PropertyDB[len(dataDBJSON.PropertyDB)-1] = repository.PropertyDB{}
		dataDBJSON.PropertyDB = dataDBJSON.PropertyDB[:len(dataDBJSON.PropertyDB)-1]
	}
	dataDBJSON.PropertyDB = append(dataDBJSON.PropertyDB, propertyDB)

	err := client.Storage.SetPudgelData()
	if err != nil {
		return err
	}

	return nil
}

func (c *Client) EditControlDoubleConDB(uid string, db repository.BasesDoubleControl) error {

	editDB, keyDB := repository.GetBasesDoubleControlDB(c.Storage.BasesDoubleControl, uid)
	if keyDB == -1 {
		editDB = repository.BasesDoubleControl{}
	}

	editDB.Name = db.Name
	editDB.Server = db.Server
	editDB.User = db.User
	editDB.Password = db.Password
	editDB.UID = uid

	dataDBJSON := client.Storage
	if keyDB != -1 {
		dataDBJSON.BasesDoubleControl[keyDB] = dataDBJSON.BasesDoubleControl[len(dataDBJSON.BasesDoubleControl)-1]
		dataDBJSON.BasesDoubleControl[len(dataDBJSON.BasesDoubleControl)-1] = repository.BasesDoubleControl{}
		dataDBJSON.BasesDoubleControl = dataDBJSON.BasesDoubleControl[:len(dataDBJSON.BasesDoubleControl)-1]
	}
	dataDBJSON.BasesDoubleControl = append(dataDBJSON.BasesDoubleControl, db)

	err := client.Storage.SetPudgelData()
	if err != nil {
		return err
	}

	return nil
}

func (c *Client) ClearPerformedActions() {
	c.PerformedActions = PerformedActions{}
}

func (s *StringArray) Find(value string) int {

	for k, v := range *s {
		if v == value {
			return k
		}
	}
	return -1
}
