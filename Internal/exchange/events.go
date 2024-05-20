package exchange

import (
	"encoding/binary"
	"fmt"
	"github.com/google/uuid"

	"github.com/recoilme/pudge"

	"Service_1Cv8/internal/constants"
)

func (m *MsgWithSorter) PutPudge(n string, p int) error {

	cfg := &pudge.Config{StoreMode: 1}
	db := &pudge.Db{}
	db.Lock()

	defer db.Close()
	defer db.Unlock()

	db, err := pudge.Open(fmt.Sprintf("%s/%s/%d", constants.Pudge, n, p), cfg)
	if err != nil {
		return err
	}

	key := m.Sorter
	if key == "" {
		key = uuid.New().String()
	}

	err = db.Set(key, m.Message)
	if err != nil {
		return err
	}

	return nil
}

func (eq *ExchangeQueue) TransferToDisk() error {

	//Вернуть
	return nil
	//cfg := &pudge.Config{StoreMode: 1}
	//db := &pudge.Db{}
	//db.Lock()
	//defer db.Close()
	//defer db.Unlock()
	//
	//for i := 0; i <= constants.NumberPriorities; i++ {
	//	pm, ok := eq.PriorityMessages[i]
	//	if !ok || len(pm) == 0 {
	//		continue
	//	}
	//
	//	db, err := pudge.Open(fmt.Sprintf("%s/%s/%d", constants.Pudge, eq.Name, i), cfg)
	//	if err != nil {
	//		return err
	//	}
	//
	//	for _, val := range pm {
	//
	//		key := val.Sorter
	//		if key == "" {
	//			key = uuid.New().String()
	//		}
	//
	//		err = db.Set(key, val.Message)
	//		if err != nil {
	//			return err
	//		}
	//
	//		pm = append(pm[1:])
	//
	//		eq.Size = eq.Size - int64(binary.Size(val.Message))
	//	}
	//	eq.PriorityMessages[i] = pm
	//
	//	c, err := db.Count()
	//	if err == nil && c == 0 {
	//		_ = db.DeleteFile()
	//	}
	//}
	//
	//if eq.TypeStorage != constants.Hard {
	//	eq.TypeStorage = constants.Hard
	//}
	//
	//return nil
}

func (eq *ExchangeQueue) TransferToMemory() error {

	cfg := &pudge.Config{StoreMode: 1}

	totalByte := 0
	totalCount := 0

	db := &pudge.Db{}
	db.Lock()

	defer db.Close()
	defer db.Unlock()

	for i := 0; i <= constants.NumberPriorities; i++ {

		db, err := pudge.Open(fmt.Sprintf("%s/%s/%d", constants.Pudge, eq.Name, i), cfg)
		if err != nil {
			return err
		}

		pm, ok := eq.PriorityMessages[i]
		if !ok {
			pm = []MsgWithSorter{}
		}

		keys, _ := db.Keys(0, 0, 0, true)
		for _, key := range keys {
			b := []byte{}

			err = db.Get(key, &b)
			if err != nil {
				continue
			}
			err = db.Delete(key)
			if err != nil {
				continue
			}

			eq.Size = eq.Size + int64(binary.Size(b))
			pm = append(pm, MsgWithSorter{Sorter: string(key), Message: b})

			totalByte = totalByte + len(b)
			totalCount++

			if totalByte > eq.MaxSizeSoftMode || totalCount > eq.MaxLenSoftMode {
				eq.PriorityMessages[i] = pm
				_ = db.Close()

				return nil

			}

		}
		if len(pm) != 0 {
			eq.PriorityMessages[i] = pm
		}

		c, err := db.Count()
		if err == nil && c == 0 {
			_ = db.DeleteFile()
		}
	}

	if eq.TypeStorage != constants.Soft {
		eq.TypeStorage = constants.Soft
	}

	return nil
}

func (eq *ExchangeQueue) LenSoft() int {

	totalLen := 0

	for _, v := range eq.PriorityMessages {
		totalLen = totalLen + len(v)
	}

	return totalLen
}

func (eq *ExchangeQueue) Len() int {

	totalLen := 0

	for _, v := range eq.PriorityMessages {
		totalLen = totalLen + len(v)
	}

	cfg := &pudge.Config{StoreMode: 1}
	db := &pudge.Db{}
	db.RLock()

	defer db.Close()
	defer db.RUnlock()

	for i := 0; i <= constants.NumberPriorities; i++ {

		db, err := pudge.Open(fmt.Sprintf("%s/%s/%d", constants.Pudge, eq.Name, i), cfg)
		c, err := db.Count()
		if err != nil {
			c = 0
		}

		totalLen = totalLen + c
	}

	return totalLen
}

func (eq *ExchangeQueue) SizeMemory() int64 {

	var totalSize int64

	if eq.TypeStorage == constants.Soft {
		return eq.Size
	}

	cfg := &pudge.Config{StoreMode: 1}
	db := &pudge.Db{}
	db.RLock()

	defer db.Close()
	defer db.RUnlock()

	for i := 0; i <= constants.NumberPriorities; i++ {
		db, err := pudge.Open(fmt.Sprintf("%s/%s/%d", constants.Pudge, eq.Name, i), cfg)

		if err != nil {
			totalSize = totalSize + 0
			continue
		}

		c, err := db.FileSize()
		if err != nil {
			c = 0
		}

		totalSize = totalSize + c
	}

	return totalSize
}
