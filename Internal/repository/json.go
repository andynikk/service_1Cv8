package repository

import (
	"encoding/json"
	"os"
	"strings"
)

var patchDB = "./db.json"

type DataDBJSON struct {
	DB         []DataBases  `yaml:"db,omitempty"`
	PropertyDB []PropertyDB `yaml:"property_db,omitempty"`
}

type Settings struct {
	PathExe1C    string
	NameUser     string
	PasswordUser string
	StartBlock   string
	FinishBlock  string
	KeyUnlock    string
	Massage      string

	PathStorage string
	PathCopy    string

	//HTTPServer        string
	HTTPPort          string
	IntervalService   string
	IntervalWebClient string
	KillProcessKb     string
	ControlServer     string

	TgAPI string
	TgID  string
}

func (d *DataJSON) MarshalData() (val []byte, err error) {
	arrJSON, err := json.MarshalIndent(d, "", " ")
	if err != nil {
		return nil, err
	}

	return arrJSON, nil
}

// WriteData Запись данных в файл
func (d DataJSON) WriteData() error {
	arrJSON, err := json.MarshalIndent(d, "", " ")
	if err != nil {
		return err
	}

	if IsOSWindows() {
		patchDB = strings.Replace(patchDB, "/", "\\", -1)
	}

	if err = os.WriteFile(patchDB, arrJSON, 0664); err != nil {
		return err
	}

	return nil
}

// GetData Получение баз и серверов из файла
func (d *DataJSON) GetData() error {

	if IsOSWindows() {
		patchDB = strings.Replace(patchDB, "/", "\\", -1)
	}

	res, err := os.ReadFile(patchDB)
	if err != nil {
		return err
	}
	if err = json.Unmarshal(res, d); err != nil {
		return err
	}

	return nil
}
