package handlers

import (
	"Service_1Cv8/internal/token"
	"context"
	"log"
	"time"

	OneCv8 "Service_1Cv8/internal/1cv8"
	"Service_1Cv8/internal/repository"
	"Service_1Cv8/internal/winsys"
)

func (srv *Server) serviceKillWinProc(ctx context.Context, cancelFunc context.CancelFunc) {
	ticker := time.NewTicker(time.Duration(srv.IntervalService) * time.Minute)
	//ticker := time.NewTicker(2 * time.Minute)

	for {
		select {
		case <-ticker.C:
			closedConnects, _ := winsys.KillWinProc(srv.ControlServer, "rphost.exe", srv.KillProcessKb)
			for _, v := range closedConnects {
				srv.Setting.ClosedTasks = append(srv.Setting.ClosedTasks, v)
			}

		case <-ctx.Done():
			cancelFunc()
			return
		}
	}
}

func (srv *Server) serviceDropDouble(ctx context.Context, cancelFunc context.CancelFunc) {
	ticker := time.NewTicker(time.Duration(srv.IntervalService) * time.Minute)

	scdc := OneCv8.SettingCloseDoubleConnection{
		IntervalWebClient: srv.IntervalWebClient,
		OutputMessages:    true,
	}
	for {
		select {
		case <-ticker.C:

			srv.ExtremeStartDD = time.Now()

			var arrM []OneCv8.MassageJSON

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

				arrM = append(arrM, m)

			}

			chanOut := make(chan repository.ClosedConnect)
			go OneCv8.DropDoubleUsersDB(arrM, scdc, chanOut)

			for {
				cc, ok := <-chanOut
				if !ok {

					break
				}

				srv.Setting.ClosedConnects = append(srv.Setting.ClosedConnects, cc)
			}

			log.Println("chan closed")

		case <-ctx.Done():
			cancelFunc()
			return
		}
	}
}

func (srv *Server) serviceResetDataDisk(ctx context.Context, cancelFunc context.CancelFunc) {
	//Вернуть
	return

	//ticker := time.NewTicker(30 * time.Second)
	//
	//for {
	//	select {
	//	case <-ticker.C:
	//
	//		for k, v := range srv.ExchangeStorage {
	//
	//			downloadAt := v.DownloadAt
	//			subTime := time.Now().Sub(downloadAt)
	//
	//			if v.TimeTransferToDisk != 0 &&
	//				v.LenSoft() != 0 &&
	//				(subTime.Minutes() > v.TimeTransferToDisk || v.Size > 104857600) {
	//
	//				srv.Lock()
	//
	//				_ = v.TransferToDisk()
	//				srv.ExchangeStorage[k] = v
	//
	//				srv.Unlock()
	//
	//			}
	//
	//		}
	//
	//	case <-ctx.Done():
	//		cancelFunc()
	//		return
	//	}
	//}
}

func (srv *Server) findClaim(n string) (token.ClaimStore, bool) {
	for _, v := range srv.ClaimsStore {
		if v.Key == n {
			return v, true
		}
	}
	return token.ClaimStore{}, false
}
