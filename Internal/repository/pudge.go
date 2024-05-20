package repository

import (
	"fmt"
	"log"
	"sync"

	"github.com/jinzhu/copier"
	"github.com/recoilme/pudge"

	"Service_1Cv8/internal/constants"
	"Service_1Cv8/internal/cryptography"
)

// SetPudgelData Сохранение баз и серверов из файла
func (d *DataJSON) SetPudgelData() error {

	wg := sync.WaitGroup{}

	wg.Add(1)
	go func() {
		defer wg.Done()
		_ = SetPudgeSetting(d.Settings)
	}()
	wg.Add(1)
	go func() {
		defer wg.Done()
		_ = SetPudgeDataBases(d.DB)
	}()
	wg.Add(1)
	go func() {
		defer wg.Done()
		_ = SetPudgePropertyDB(d.PropertyDB)
	}()
	wg.Add(1)
	go func() {
		defer wg.Done()
		_ = SetPudgeServices(d.Services)
	}()
	wg.Add(1)
	go func() {
		defer wg.Done()
		_ = SetPudgeBasesDoubleControl(d.BasesDoubleControl)
	}()

	wg.Wait()

	return nil
}

// GetPudgeData Получение баз и серверов из файла
func (d *DataJSON) GetPudgeData() {

	wg := sync.WaitGroup{}

	wg.Add(1)
	go func() {
		defer wg.Done()
		err := GetPudgeSetting(&d.Settings)
		if err != nil {
			log.Println("GetPudgeSetting", err)
		}
	}()
	wg.Add(1)
	go func() {
		defer wg.Done()
		err := GetPudgeDataBases(&d.DB)
		if err != nil {
			log.Println("GetPudgeDataBases", err)
		}
	}()
	wg.Add(1)
	go func() {
		defer wg.Done()
		err := GetPudgePropertyDB(&d.PropertyDB)
		if err != nil {
			log.Println("GetPudgePropertyDB", err)
		}
	}()
	wg.Add(1)
	go func() {
		defer wg.Done()
		err := GetPudgeServices(&d.Services)
		if err != nil {
			log.Println("GetPudgeServices", err)
		}
	}()
	wg.Add(1)
	go func() {
		defer wg.Done()
		err := GetPudgeBasesDoubleControl(&d.BasesDoubleControl)
		if err != nil {
			log.Println("GetPudgeBasesDoubleControl", err)
		}
	}()

	wg.Wait()
}

func SetPudgeSetting(s Settings) error {

	var settings Settings
	err := copier.Copy(&settings, s)
	if err != nil {
		return err
	}

	settings.PasswordUser = cryptography.EncryptString(settings.PasswordUser, constants.HashKey)
	err = pudge.Set(constants.PudgeSetting, 0, settings)
	if err != nil {
		log.Println(fmt.Sprintf("error save setting value `%s`", s))
		return err
	}

	return nil
}

func SetPudgeDataBases(dbs []DataBases) error {
	db, err := pudge.Open(constants.PudgeDataBases, &pudge.Config{SyncInterval: 1})
	defer db.Close()

	if err != nil {
		log.Println(fmt.Sprintf("error del table `%s`", constants.PudgeDataBases))
		return err
	}

	for key, val := range dbs {
		err = db.Set(key, val)
		if err != nil {
			log.Println(fmt.Sprintf("error save setting `%s` value `%s`", "databases", val))
			return err
		}
	}

	return nil
}

func SetPudgePropertyDB(pdb []PropertyDB) error {
	db, err := pudge.Open(constants.PudgePropertyDB, &pudge.Config{SyncInterval: 1})
	defer db.Close()
	if err != nil {
		log.Println(fmt.Sprintf("error del table `%s`", constants.PudgePropertyDB))
		return err
	}

	var propertyDB []PropertyDB
	err = copier.Copy(&propertyDB, pdb)
	if err != nil {
		return err
	}

	hashKey := constants.HashKey
	for k, v := range propertyDB {
		propertyDB[k].PasswordUser =
			cryptography.EncryptString(v.PasswordUser, hashKey)
	}

	for key, val := range propertyDB {
		err = db.Set(key, val)
		if err != nil {
			log.Println(fmt.Sprintf("error save setting `%s` value `%s`", "property_db", val.UID))
			return err
		}
	}

	return nil
}

func SetPudgeServices(s []Services) error {
	db, err := pudge.Open(constants.PudgeServices, &pudge.Config{SyncInterval: 1})
	defer db.Close()

	if err != nil {
		log.Println(fmt.Sprintf("error del table `%s`", constants.PudgeServices))
		return err
	}

	var services []Services
	err = copier.Copy(&services, s)
	if err != nil {
		return err
	}

	hashKey := constants.HashKey
	for k, v := range s {
		services[k].Password = cryptography.EncryptString(v.Password, hashKey)

		var sqlServer []SQLServer
		err = copier.Copy(&sqlServer, services[k].SQLServers)
		if err != nil {
			return err
		}
		services[k].SQLServers = sqlServer
		for kSQLServers, vSQLServers := range v.SQLServers {
			services[k].SQLServers[kSQLServers].Password =
				cryptography.EncryptString(vSQLServers.Password, hashKey)
		}
	}

	for key, val := range services {
		err = db.Set(key, val)
		if err != nil {
			log.Println(fmt.Sprintf("error save setting `%s` value `%s`", "services", val))
			return err
		}
	}

	return nil
}

func SetPudgeBasesDoubleControl(bdc []BasesDoubleControl) error {
	db, err := pudge.Open(constants.PudgeBasesDoubleControl, &pudge.Config{SyncInterval: 1})
	defer db.Close()
	if err != nil {
		log.Println(fmt.Sprintf("error del table `%s`", constants.PudgeBasesDoubleControl))
		return err
	}

	var basesDoubleControl []BasesDoubleControl
	err = copier.Copy(&basesDoubleControl, bdc)
	if err != nil {
		return err
	}

	hashKey := constants.HashKey
	for k, v := range basesDoubleControl {
		basesDoubleControl[k].Password =
			cryptography.EncryptString(v.Password, hashKey)
	}

	for key, val := range basesDoubleControl {
		err = db.Set(key, val)
		if err != nil {
			log.Println(fmt.Sprintf("error save setting `%s` value `%s`", "bases_double_control", val))
			return err
		}
	}

	return nil
}

func GetPudgeSetting(s *Settings) error {

	db, err := pudge.Open(constants.PudgeSetting, &pudge.Config{SyncInterval: 1})
	defer db.Close()

	if err != nil {
		return err
	}
	err = db.Get(0, s)
	if err != nil {
		return err
	}

	s.PasswordUser = cryptography.DecryptString(s.PasswordUser, constants.HashKey)

	return nil
}

func GetPudgeDataBases(dbs *[]DataBases) error {
	db, err := pudge.Open(constants.PudgeDataBases, &pudge.Config{SyncInterval: 1})
	defer db.Close()

	count, err := db.Count()
	if err != nil {
		return err
	}

	for i := 0; i < count; i++ {
		var newDbs = new(DataBases)

		err = db.Get(i, newDbs)
		if err != nil {
			return err
		}
		*dbs = append(*dbs, *newDbs)
	}

	return nil
}

func GetPudgePropertyDB(pdb *[]PropertyDB) error {
	db, err := pudge.Open(constants.PudgePropertyDB, &pudge.Config{SyncInterval: 1})
	defer db.Close()

	count, err := db.Count()
	if err != nil {
		return err
	}

	for i := 0; i < count; i++ {
		var newPdb = new(PropertyDB)

		err = db.Get(i, newPdb)
		if err != nil {
			return err
		}
		newPdb.PasswordUser = cryptography.DecryptString(newPdb.PasswordUser, constants.HashKey)

		*pdb = append(*pdb, *newPdb)
	}

	return nil
}

func GetPudgeServices(s *[]Services) error {
	db, err := pudge.Open(constants.PudgeServices, &pudge.Config{SyncInterval: 1})
	defer db.Close()

	count, err := db.Count()
	if err != nil {
		return err
	}

	for i := 0; i < count; i++ {
		var newS = new(Services)

		err = db.Get(i, newS)
		if err != nil {
			return err
		}
		newS.Password = cryptography.DecryptString(newS.Password, constants.HashKey)
		for k, v := range newS.SQLServers {
			newS.SQLServers[k].Password = cryptography.DecryptString(v.Password, constants.HashKey)
		}

		*s = append(*s, *newS)
	}

	return nil
}

func GetPudgeBasesDoubleControl(bdc *[]BasesDoubleControl) error {
	db, err := pudge.Open(constants.PudgeBasesDoubleControl, &pudge.Config{SyncInterval: 1})
	if err != nil {
		return err
	}

	defer db.Close()

	count, err := db.Count()
	if err != nil {
		return err
	}

	for i := 0; i < count; i++ {
		var newBdc = new(BasesDoubleControl)

		err = db.Get(i, newBdc)
		if err != nil {
			return err
		}
		newBdc.Password = cryptography.DecryptString(newBdc.Password, constants.HashKey)

		*bdc = append(*bdc, *newBdc)
	}

	return nil
}

func DelPudge(n string) error {
	db, err := pudge.Open(fmt.Sprintf("%s/%s", constants.Pudge, n), &pudge.Config{SyncInterval: 1})
	if err != nil {
		return err
	}
	defer db.Close()

	err = db.DeleteFile()
	if err != nil {
		return err
	}

	return nil
}
