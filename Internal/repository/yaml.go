package repository

import (
	"Service_1Cv8/internal/constants"
	"Service_1Cv8/internal/cryptography"
	"gopkg.in/yaml.v3"
	"os"
	"strings"
)

var patchYAMLDB = "./db.yaml"

func (d *DataJSON) MarshalYamlData() (val []byte, err error) {
	arrJSON, err := yaml.Marshal(d)
	if err != nil {
		return nil, err
	}

	return arrJSON, nil
}

// WriteYamlData Запись данных в файл
func (d DataJSON) WriteYamlData() error {

	hashKey := constants.HashKey
	for kServices, vServices := range d.Services {
		d.Services[kServices].Password = cryptography.EncryptString(vServices.Password, hashKey)
		for kSQLServers, vSQLServers := range vServices.SQLServers {
			d.Services[kServices].SQLServers[kSQLServers].Password =
				cryptography.EncryptString(vSQLServers.Password, hashKey)
		}
	}

	for kBasesDoubleControl, vBasesDoubleControl := range d.BasesDoubleControl {
		d.BasesDoubleControl[kBasesDoubleControl].Password =
			cryptography.EncryptString(vBasesDoubleControl.Password, hashKey)
	}

	d.Settings.PasswordUser = cryptography.EncryptString(d.Settings.PasswordUser, hashKey)

	arrJSON, err := yaml.Marshal(d)
	if err != nil {
		return err
	}

	if IsOSWindows() {
		patchYAMLDB = strings.Replace(patchYAMLDB, "/", "\\", -1)
	}

	if err = os.WriteFile(patchYAMLDB, arrJSON, 0664); err != nil {
		return err
	}

	return nil
}

// GetYamlData Получение баз и серверов из файла
func (d *DataJSON) GetYamlData() error {

	if IsOSWindows() {
		patchYAMLDB = strings.Replace(patchYAMLDB, "/", "\\", -1)
	}

	res, err := os.ReadFile(patchYAMLDB)
	if err != nil {
		return err
	}
	if err = yaml.Unmarshal(res, d); err != nil {
		return err
	}

	hashKey := constants.HashKey
	for kServices, vServices := range d.Services {
		d.Services[kServices].Password = cryptography.DecryptString(vServices.Password, hashKey)
		for kSQLServers, vSQLServers := range vServices.SQLServers {
			d.Services[kServices].SQLServers[kSQLServers].Password = cryptography.DecryptString(vSQLServers.Password, hashKey)
		}
	}

	for kBasesDoubleControl, vBasesDoubleControl := range d.BasesDoubleControl {
		d.BasesDoubleControl[kBasesDoubleControl].Password =
			cryptography.DecryptString(vBasesDoubleControl.Password, hashKey)
	}

	d.Settings.PasswordUser = cryptography.DecryptString(d.Settings.PasswordUser, hashKey)

	return nil
}
