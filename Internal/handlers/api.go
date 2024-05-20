package handlers

import (
	"Service_1Cv8/internal/answers"
	"Service_1Cv8/internal/compression"
	"Service_1Cv8/internal/constants"
	"Service_1Cv8/internal/exchange"
	"Service_1Cv8/internal/files"
	"Service_1Cv8/internal/repository"
	"Service_1Cv8/internal/telegram"
	"Service_1Cv8/internal/token"
	"bytes"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"github.com/boltdb/bolt"
	"github.com/gorilla/mux"
	"github.com/jinzhu/copier"
	"github.com/recoilme/pudge"
	"gopkg.in/yaml.v3"
	"io"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"
)

func (srv *Server) hidePassword() (*Setting, error) {
	setting := Setting{}
	err := copier.Copy(&setting, srv.Setting)
	if err != nil {
		return nil, err
	}

	var messageCom []repository.BasesDoubleControl
	err = copier.Copy(&messageCom, srv.MessageCom)
	if err != nil {
		return nil, err
	}

	for kMessageCom := range messageCom {
		messageCom[kMessageCom].Password = "•"
	}

	setting.MessageCom = messageCom

	return &setting, nil
}

// handlerNotFound, хендлер начальной страницы сервера
func (srv *Server) handleFunc(rw http.ResponseWriter, rq *http.Request) {

	setting, err := srv.hidePassword()
	if err != nil {
		return
	}

	out, _ := yaml.Marshal(setting)

	if _, err := rw.Write(out); err != nil {
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}

	rw.WriteHeader(http.StatusOK)
}

// handlerNotFound, хендлер начальной страницы сервера
func (srv *Server) handleFuncQueues(w http.ResponseWriter, r *http.Request) {

	arrEQI := []exchange.ExchangeQueueInfo{}
	for k, v := range srv.ExchangeStorage {
		arrEQI = append(arrEQI, exchange.ExchangeQueueInfo{
			Name:        k,
			DownloadAt:  v.DownloadAt,
			UploadAt:    v.UploadAt,
			TypeStorage: v.TypeStorage.String(),
			Size:        fmt.Sprintf("%s KBt", files.GroupSeparator(fmt.Sprintf("%d", v.Size/1024))),
			Messages:    fmt.Sprintf("%s", files.GroupSeparator(fmt.Sprintf("%d", v.Len()))),
		})
	}

	out, _ := yaml.Marshal(arrEQI)
	if _, err := w.Write(out); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

// handlerNotFound, хендлер начальной страницы сервера
func (srv *Server) handleFuncUpdateDB(rw http.ResponseWriter, rq *http.Request) {

	srv.InitConfig()
	out, _ := yaml.Marshal(srv.Setting)

	if _, err := rw.Write(out); err != nil {
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}

	rw.WriteHeader(http.StatusOK)
}

func (srv *Server) handleFuncCheckQueue(w http.ResponseWriter, r *http.Request) {

	muxVars := mux.Vars(r)
	nameQ := muxVars["nameQ"]

	eq, ok := srv.ExchangeStorage[nameQ]
	if !ok {
		w.Write([]byte(fmt.Sprintf("Queue %s not found", nameQ)))
		w.WriteHeader(http.StatusNotFound)
		return
	}

	if eq.Len() != 0 {
		nErr := answers.AnPreviousMessageNotRead
		http.Error(w, nErr.Error(), answers.StatusHTTP(nErr))
		w.WriteHeader(http.StatusConflict)
		return
	}

	nErr := answers.AnEmpty
	http.Error(w, nErr.Error(), answers.StatusHTTP(nErr))
	w.WriteHeader(http.StatusNoContent)
}

/////////////////////////////////////////////////////////////////////////////////////////////////////

// handleFuncPutQueue кладет в хранилище
func (srv *Server) handleFuncPutQueue(w http.ResponseWriter, r *http.Request) {
	k := r.Header.Get("Authentication-Key")

	cs, ok := srv.findClaim(k)
	switch ok {
	case true:
		if _, ok = token.ExtractClaims(string(cs.Value), cs.Secret); !ok {
			cs, ok = repository.CheckToken(k)
			if !ok {
				http.Error(w, fmt.Sprintf("not valid key %s", k), http.StatusNotAcceptable)
				return
			}
			srv.ClaimsStore = append(srv.ClaimsStore, cs)
		}
	default:
		cs, ok = repository.CheckToken(k)
		if !ok {
			http.Error(w, fmt.Sprintf("not valid key %s", k), http.StatusNotAcceptable)
			return
		}
		srv.ClaimsStore = append(srv.ClaimsStore, cs)
	}

	bytBody, err := io.ReadAll(r.Body)
	if err != nil {
		log.Println(fmt.Sprintf("$$ 1 %s", err.Error()))
		http.Error(w, "Ошибка получения Content-Encoding", http.StatusInternalServerError)
		return
	}

	bodyJSON := []byte{}
	contentEncoding := r.Header.Get("Content-Encoding")
	if !strings.Contains(contentEncoding, "gzip") {
		bodyJSON, err = compression.Compress(bytBody)
		if err != nil {
			log.Println(fmt.Sprintf("$$ 3 %s", err.Error()))
			http.Error(w, "Ошибка упаковки сообщения GZIP", http.StatusInternalServerError)
			return
		}
	}

	srv.Lock()
	defer srv.Unlock()

	name := r.Header.Get("Pudge")
	priority := r.Header.Get("Priority")
	key := r.Header.Get("Key")
	timeTransfer := r.Header.Get("Time-Transfer-To-Disk")
	maxSizeSoftMode := r.Header.Get("Max-Size-Soft-Mode")
	maxLenSoftMode := r.Header.Get("Max-Len-Soft-Mode")

	fltTimeTransfer, err := strconv.ParseFloat(timeTransfer, 64)
	if err != nil {
		fltTimeTransfer = constants.TimeTransferToDisk
	}

	intPriority, err := strconv.Atoi(priority)
	if err != nil {
		intPriority = constants.NumberPriorities
	}

	intMaxSizeSoftMode, err := strconv.Atoi(maxSizeSoftMode)
	if err != nil {
		intMaxSizeSoftMode = constants.TotalByteInPudge
	}

	intMaxLenSoftMode, err := strconv.Atoi(maxLenSoftMode)
	if err != nil {
		intMaxLenSoftMode = constants.TotalCount
	}

	eq, ok := srv.ExchangeStorage[name]
	if !ok {
		eq = exchange.ExchangeQueue{
			Name:             name,
			UploadAt:         time.Now(),
			DownloadAt:       time.Now(),
			TypeStorage:      constants.Soft,
			PriorityMessages: exchange.PriorityMessages{},
		}
	}

	eq.TimeTransferToDisk = fltTimeTransfer
	eq.MaxSizeSoftMode = intMaxSizeSoftMode
	eq.MaxLenSoftMode = intMaxLenSoftMode

	if eq.TypeStorage == constants.Hard {
		msg := exchange.MsgWithSorter{
			Sorter:  key,
			Message: bodyJSON,
		}
		err = msg.PutPudge(name, intPriority)
		if err != nil {
			log.Println(err)
		}
		return
	}

	eq.Size = eq.Size + int64(binary.Size(bodyJSON))

	eq.UploadAt = time.Now()
	pm, ok := eq.PriorityMessages[intPriority]
	if !ok {
		pm = []exchange.MsgWithSorter{}
	}

	pm = append(pm, exchange.MsgWithSorter{Sorter: key, Message: bodyJSON})

	eq.PriorityMessages[intPriority] = pm
	srv.ExchangeStorage[name] = eq

	w.WriteHeader(http.StatusOK)
}

// handleFuncPickQueue забирает из хранилища
func (srv *Server) handleFuncPickQueue(w http.ResponseWriter, r *http.Request) {
	rVars := mux.Vars(r)
	nameP := rVars["nameQ"]

	k := r.Header.Get("Authentication-Key")
	cs, ok := srv.findClaim(k)
	switch ok {
	case true:
		if _, ok = token.ExtractClaims(string(cs.Value), cs.Secret); !ok {
			cs, ok = repository.CheckToken(k)
			if !ok {
				http.Error(w, fmt.Sprintf("not valid key %s", k), http.StatusNotAcceptable)
				return
			}
			//srv.ClaimsStore = append(srv.ClaimsStore, cs)
		}
	default:
		cs, ok = repository.CheckToken(k)
		if !ok {
			http.Error(w, fmt.Sprintf("not valid key %s", k), http.StatusNotAcceptable)
			return
		}
		//srv.ClaimsStore = append(srv.ClaimsStore, cs)
	}

	srv.Lock()
	defer srv.Unlock()

	eq, ok := srv.ExchangeStorage[nameP]

	if !ok {
		http.Error(w, fmt.Sprintf("Queue %s no content", nameP), http.StatusNoContent)
		return
	}

	if eq.LenSoft() == 0 && eq.TypeStorage == constants.Hard {
		err := eq.TransferToMemory()
		if err != nil {
			http.Error(w, fmt.Sprintf("Error transfer to memory %s", nameP), http.StatusInternalServerError)
			return
		}
	}

	eq.DownloadAt = time.Now()

	w.Header().Set("content-encoding", "gzip")
	if len(eq.PriorityMessages) == 0 {
		http.Error(w, fmt.Sprintf("The queue %s is empty", nameP), http.StatusNoContent)
		return
	}

	msg := exchange.MsgWithSorter{}
	for i := 0; i <= constants.NumberPriorities; i++ {
		pm, ok := eq.PriorityMessages[i]
		if !ok || len(pm) == 0 {
			continue
		}

		msg = pm[0]
		pm = append(pm[1:])

		eq.PriorityMessages[i] = pm
		eq.Size = eq.Size - int64(binary.Size(msg.Message))

		break
	}

	if len(msg.Message) == 0 {
		http.Error(w, fmt.Sprintf("The queue %s is empty", nameP), http.StatusNoContent)
		return
	}

	srv.ExchangeStorage[nameP] = eq

	w.Write(msg.Message)
	w.Header().Set("Key", msg.Sorter)
	w.WriteHeader(http.StatusOK)
}

// handleFuncAddQueue добавляет очередь
func (srv *Server) handleFuncAddQueue(w http.ResponseWriter, r *http.Request) {
	rVars := mux.Vars(r)
	nameQ := rVars["nameQ"]

	k := r.Header.Get("Authentication-Key")
	cs, ok := srv.findClaim(k)
	switch ok {
	case true:
		if _, ok = token.ExtractClaims(string(cs.Value), cs.Secret); !ok {
			cs, ok = repository.CheckToken(k)
			if !ok {
				http.Error(w, fmt.Sprintf("not valid key %s", k), http.StatusNotAcceptable)
				return
			}
			//srv.ClaimsStore = append(srv.ClaimsStore, cs)
		}
	default:
		cs, ok = repository.CheckToken(k)
		if !ok {
			http.Error(w, fmt.Sprintf("not valid key %s", k), http.StatusNotAcceptable)
			return
		}
		//srv.ClaimsStore = append(srv.ClaimsStore, cs)
	}

	srv.Lock()
	defer srv.Unlock()

	eq, ok := srv.ExchangeStorage[nameQ]

	if !ok {
		eq = exchange.ExchangeQueue{
			Name:             nameQ,
			TypeStorage:      constants.Soft,
			UploadAt:         time.Now(),
			DownloadAt:       time.Now(),
			PriorityMessages: map[int][]exchange.MsgWithSorter{},
		}
		srv.ExchangeStorage[nameQ] = eq
	}

	w.Write([]byte{})
	w.WriteHeader(http.StatusOK)
}

// handleFuncDelQueue удаляет очередь
func (srv *Server) handleFuncDelQueue(w http.ResponseWriter, r *http.Request) {
	rVars := mux.Vars(r)
	nameQ := rVars["nameQ"]

	k := r.Header.Get("Authentication-Key")
	cs, ok := srv.findClaim(k)
	switch ok {
	case true:
		if _, ok = token.ExtractClaims(string(cs.Value), cs.Secret); !ok {
			cs, ok = repository.CheckToken(k)
			if !ok {
				http.Error(w, fmt.Sprintf("not valid key %s", k), http.StatusNotAcceptable)
				return
			}
			//srv.ClaimsStore = append(srv.ClaimsStore, cs)
		}
	default:
		cs, ok = repository.CheckToken(k)
		if !ok {
			http.Error(w, fmt.Sprintf("not valid key %s", k), http.StatusNotAcceptable)
			return
		}
		//srv.ClaimsStore = append(srv.ClaimsStore, cs)
	}

	srv.Lock()
	defer srv.Unlock()

	delete(srv.ExchangeStorage, nameQ)
	//_ = files.DelFolder(nameQ)

	w.WriteHeader(http.StatusOK)
}

// handleFuncTgSend очищает очередь
func (srv *Server) handleFuncClearQueue(w http.ResponseWriter, r *http.Request) {
	rVars := mux.Vars(r)
	nameQ := rVars["nameQ"]

	k := r.Header.Get("Authentication-Key")
	cs, ok := srv.findClaim(k)
	switch ok {
	case true:
		if _, ok = token.ExtractClaims(string(cs.Value), cs.Secret); !ok {
			cs, ok = repository.CheckToken(k)
			if !ok {
				http.Error(w, fmt.Sprintf("not valid key %s", k), http.StatusNotAcceptable)
				return
			}
			//srv.ClaimsStore = append(srv.ClaimsStore, cs)
		}
	default:
		cs, ok = repository.CheckToken(k)
		if !ok {
			http.Error(w, fmt.Sprintf("not valid key %s", k), http.StatusNotAcceptable)
			//return
		}
		//srv.ClaimsStore = append(srv.ClaimsStore, cs)
	}

	srv.Lock()
	defer srv.Unlock()

	eq, ok := srv.ExchangeStorage[nameQ]

	if ok {
		eq = exchange.ExchangeQueue{
			Name:             nameQ,
			TypeStorage:      constants.Soft,
			UploadAt:         time.Now(),
			DownloadAt:       time.Now(),
			PriorityMessages: map[int][]exchange.MsgWithSorter{},
		}
		srv.ExchangeStorage[nameQ] = eq
	}

	_ = pudge.Delete(fmt.Sprintf("%s", constants.Pudge), nameQ)

	w.WriteHeader(http.StatusOK)
}

//////////////////////////////////////////////////////////////////////////////////////////

// handleFuncTgSend очищает очередь
func (srv *Server) handleFuncTgSend(w http.ResponseWriter, r *http.Request) {
	k := r.Header.Get("Authentication-Key")

	cs, ok := srv.findClaim(k)
	switch ok {
	case true:
		if _, ok = token.ExtractClaims(string(cs.Value), cs.Secret); !ok {
			cs, ok = repository.CheckToken(k)
			if !ok {
				http.Error(w, fmt.Sprintf("not valid key %s", k), http.StatusNotAcceptable)
				return
			}
			srv.ClaimsStore = append(srv.ClaimsStore, cs)
		}
	default:
		cs, ok = repository.CheckToken(k)
		if !ok {
			http.Error(w, fmt.Sprintf("not valid key %s", k), http.StatusNotAcceptable)
			return
		}
		srv.ClaimsStore = append(srv.ClaimsStore, cs)
	}

	bytBody, err := io.ReadAll(r.Body)
	if err != nil {
		log.Println(fmt.Sprintf("$$ 1 %s", err.Error()))
		http.Error(w, "Ошибка получения", http.StatusInternalServerError)
		return
	}

	if r.Header.Get("Content-Encoding") == "gzip" {
		bytBody, err = compression.Decompress(bytBody)
		if err != nil {
			log.Println(fmt.Sprintf("$$ 1 %s", err.Error()))
			http.Error(w, "Ошибка получения Content-Encoding", http.StatusInternalServerError)
			return
		}
	}

	tMsg := telegram.TgMsg{}
	err = json.Unmarshal(bytBody, &tMsg)
	if err != nil {
		http.Error(w, "Ошибка получения тела тг сообщения", http.StatusInternalServerError)
		return
	}

	msg := tMsg.Message
	emoji := ""
	if tMsg.Emoji.ID != "" {
		emoji = fmt.Sprintf("%s%s", tMsg.Emoji.ID, tMsg.Emoji.Сaption)
	}
	c := telegram.New(tMsg.API)
	err = c.SendMessage(fmt.Sprintf("*%s*\n%s", msg, emoji), tMsg.ID)
	if err != nil {
		http.Error(w, "Ошибка отправки тг сообщения", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

///////////////////////////////////////////////////

func (srv *Server) handleFuncPutBucket(w http.ResponseWriter, r *http.Request) {

	bytBody, err := io.ReadAll(r.Body)
	if err != nil {
		log.Println(fmt.Sprintf("$$ 1 %s", err.Error()))
		http.Error(w, "Ошибка получения Content-Encoding", http.StatusInternalServerError)
		return
	}

	contentEncoding := r.Header.Get("Content-Encoding")
	if strings.Contains(contentEncoding, "gzip") {
		bytBody, err = compression.Decompress(bytBody)
		log.Println(fmt.Sprintf("$$ 2 %s", err.Error()))
		if err != nil {
			http.Error(w, "Ошибка распаковки", http.StatusInternalServerError)
			return
		}
	}

	bodyJSON := bytes.NewReader(bytBody)

	var m exchange.MessageBolt
	err = json.NewDecoder(bodyJSON).Decode(&m)
	if err != nil {
		log.Println(fmt.Sprintf("$$ 3 %s", err.Error()))
		http.Error(w, "Ошибка получения JSON", http.StatusInternalServerError)
		return
	}

	gzipBody, err := compression.Compress(bytBody)

	tx, err := srv.BoltDB.Begin(true)
	if err != nil {
		w.WriteHeader(http.StatusLocked)
		return
	}
	defer tx.Rollback()

	bucketQ, err := tx.CreateBucketIfNotExists([]byte(m.Bucket))
	if err != nil {
		tx.Rollback()
		w.WriteHeader(http.StatusLocked)
		return
	}

	bucketP, err := bucketQ.CreateBucketIfNotExists([]byte(m.Priority))
	if err != nil {
		tx.Rollback()
		w.WriteHeader(http.StatusLocked)
		return
	}

	err = bucketP.Put([]byte(m.UID), gzipBody)
	if err != nil {
		tx.Rollback()
		w.WriteHeader(http.StatusLocked)
		return
	}

	if err = tx.Commit(); err != nil {
		tx.Rollback()
		w.WriteHeader(http.StatusLocked)
		return
	}

	//srv.BoltDB.View(func(tx *bolt.Tx) error {
	//	gzipBody, err := compression.Compress(bytBody)
	//
	//	bucketQ, err := tx.CreateBucketIfNotExists([]byte(m.Bucket))
	//	if err != nil {
	//		w.WriteHeader(http.StatusLocked)
	//		return err
	//	}
	//
	//	bucketP, err := bucketQ.CreateBucketIfNotExists([]byte(m.Priority))
	//	if err != nil {
	//		w.WriteHeader(http.StatusLocked)
	//		return err
	//	}
	//
	//	err = bucketP.Put([]byte(m.UID), gzipBody)
	//	if err != nil {
	//		w.WriteHeader(http.StatusLocked)
	//		return err
	//	}
	//
	//	return nil
	//})

	w.WriteHeader(http.StatusOK)
}

func (srv *Server) handleFuncDelBucket(w http.ResponseWriter, r *http.Request) {
	rVars := mux.Vars(r)
	nameB := rVars["nameB"]

	srv.BoltDB.Update(func(tx *bolt.Tx) error {
		err := tx.DeleteBucket([]byte(nameB))
		if err != nil {
			w.WriteHeader(http.StatusLocked)
			return err
		}

		return nil
	})

	w.WriteHeader(answers.StatusHTTP(answers.AnEmpty))
}

func (srv *Server) handleFuncPickBucket(w http.ResponseWriter, r *http.Request) {
	rVars := mux.Vars(r)
	nameB := rVars["nameB"]
	srv.BoltDB.Stats()
	messageReceived := false

	srv.BoltDB.View(func(tx *bolt.Tx) error {

		bucketDb := tx.Bucket([]byte(nameB))
		if bucketDb == nil {
			return fmt.Errorf("bucket not found")
		}

		c := bucketDb.Cursor()
		for k, _ := c.First(); k != nil; k, _ = c.Next() {

			bucketP := bucketDb.Bucket(k)
			cMsg := bucketP.Cursor()

			kMsg, vMsg := cMsg.First()
			if kMsg == nil {
				continue
			}

			////for kMsg, vMsg := cMsg.First(); kMsg != nil; kMsg, vMsg = cMsg.Next() {
			//for kMsg, vMsg := cMsg.First(); kMsg != nil; {

			b, _ := compression.Decompress(vMsg)
			_, err := w.Write(b)
			if err != nil {
				http.Error(w, "Ошибка передачи в тело", http.StatusInternalServerError)
				w.WriteHeader(http.StatusInternalServerError)
			}

			txUpdate, err := srv.BoltDB.Begin(true)
			if err != nil {
				break
			}
			defer txUpdate.Rollback()

			bUpdate := txUpdate.Bucket([]byte(nameB))
			pUpdate := bUpdate.Bucket(k)
			err = pUpdate.Delete(kMsg)
			if err != nil {
				break
			}

			if err = txUpdate.Commit(); err != nil {
				break
			}

			w.WriteHeader(http.StatusOK)
			messageReceived = true
			return nil
		}
		return nil
	})

	if !messageReceived {
		w.WriteHeader(http.StatusNoContent)
	}
}

//////////////////////////////////////////////////////////////////////////////////////////

func (srv *Server) handleFuncPutPudge(w http.ResponseWriter, r *http.Request) {
	bytBody, err := io.ReadAll(r.Body)
	if err != nil {
		log.Println(fmt.Sprintf("$$ 1 %s", err.Error()))
		http.Error(w, "Ошибка получения Content-Encoding", http.StatusInternalServerError)
		return
	}

	contentEncoding := r.Header.Get("Content-Encoding")
	if strings.Contains(contentEncoding, "gzip") {
		bytBody, err = compression.Decompress(bytBody)
		log.Println(fmt.Sprintf("$$ 2 %s", err.Error()))
		if err != nil {
			http.Error(w, "Ошибка распаковки", http.StatusInternalServerError)
			return
		}
	}

	bodyJSON := bytes.NewReader(bytBody)

	arrM := []exchange.MessagePudge{}
	err = json.NewDecoder(bodyJSON).Decode(&arrM)
	if err != nil {
		log.Println(fmt.Sprintf("$$ 3 %s", err.Error()))
		http.Error(w, "Ошибка получения JSON", http.StatusInternalServerError)
		return
	}

	//gzipBody, err := compression.Compress(bytBody)

	cfg := &pudge.Config{StoreMode: 1}
	db, err := pudge.Open(fmt.Sprintf("%s/%s", constants.Pudge, r.Header.Get("Pudge")), cfg)
	if err != nil {
		w.WriteHeader(http.StatusConflict)
		return
	}
	defer db.Close()

	for _, m := range arrM {
		err = db.Set(m.Key, m)
		if err != nil {
			w.WriteHeader(http.StatusConflict)
			return
		}
	}

	w.WriteHeader(http.StatusOK)
}

func (srv *Server) handleFuncPickPudge(w http.ResponseWriter, r *http.Request) {
	rVars := mux.Vars(r)
	nameP := rVars["nameP"]

	cfg := &pudge.Config{StoreMode: 1}
	db, err := pudge.Open(fmt.Sprintf("%s/%s", constants.Pudge, nameP), cfg)
	if err != nil {
		return
	}
	defer db.Close()

	//db.RLock()
	//defer db.RUnlock()

	keys, _ := db.Keys(0, 0, 0, true)
	if len(keys) == 0 {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	//arrKey := []string{}
	arrM := []exchange.MessagePudge{}

	for _, key := range keys {
		var m exchange.MessagePudge

		db.Get(key, &m)
		//db.Delete(key)

		//arrKey = append(arrKey, string(key))
		arrM = append(arrM, m)
	}

	mJSON, _ := json.Marshal(arrM)
	_, err = w.Write(mJSON)

	//for _, v := range arrKey {
	//	db.Delete(v)
	//}
	//
	//count, err := db.Count()
	//if err != nil {
	//	w.WriteHeader(http.StatusConflict)
	//	return
	//}
	//
	//if count == 0 {
	//	db.DeleteFile()
	//}

	w.WriteHeader(http.StatusOK)
}

func (srv *Server) handleFuncClearPudge(w http.ResponseWriter, r *http.Request) {
	bytBody, err := io.ReadAll(r.Body)
	if err != nil {
		log.Println(fmt.Sprintf("$$ 1 %s", err.Error()))
		http.Error(w, "Ошибка получения Content-Encoding", http.StatusInternalServerError)
		return
	}

	contentEncoding := r.Header.Get("Content-Encoding")
	if strings.Contains(contentEncoding, "gzip") {
		bytBody, err = compression.Decompress(bytBody)
		log.Println(fmt.Sprintf("$$ 2 %s", err.Error()))
		if err != nil {
			http.Error(w, "Ошибка распаковки", http.StatusInternalServerError)
			return
		}
	}

	bodyJSON := bytes.NewReader(bytBody)

	arrK := []exchange.MessageKey{}
	err = json.NewDecoder(bodyJSON).Decode(&arrK)
	if err != nil {
		log.Println(fmt.Sprintf("$$ 3 %s", err.Error()))
		http.Error(w, "Ошибка получения JSON", http.StatusInternalServerError)
		return
	}

	//gzipBody, err := compression.Compress(bytBody)

	cfg := &pudge.Config{StoreMode: 1}
	db, err := pudge.Open(fmt.Sprintf("%s/%s", constants.Pudge, r.Header.Get("Pudge")), cfg)
	if err != nil {
		w.WriteHeader(http.StatusConflict)
		return
	}
	defer db.Close()

	for _, k := range arrK {
		err = db.Delete(k.Key)
		if err != nil {
			w.WriteHeader(http.StatusConflict)
			return
		}
	}

	count, err := db.Count()
	if err != nil {
		w.WriteHeader(http.StatusConflict)
		return
	}

	if count == 0 {
		db.DeleteFile()
	}

	w.WriteHeader(http.StatusOK)
}

/////////////////////////////////////////////////////////////////////////////////////////
