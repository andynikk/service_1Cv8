// Package postgresql: работа с базой данных
package postgresql

import (
	"context"
	"errors"

	"Service_1Cv8/internal/constants/errs"
	"Service_1Cv8/internal/environment"
	"Service_1Cv8/internal/postgresql/model"

	"github.com/jackc/pgx/v5/pgxpool"
)

type DBConnector struct {
	Pool *pgxpool.Pool
	Cfg  *environment.DBConfig
}

// NewDBConnector создание конекта с базой и установка свойств конфигурации БД
func NewDBConnector(dbCfg *environment.DBConfig) (*DBConnector, error) {

	if dbCfg.DatabaseDsn == "" {
		return nil, errors.New("пустой путь к базе")
	}

	ctx, cancelFunc := context.WithCancel(context.Background())
	pool, err := pgxpool.New(ctx, dbCfg.DatabaseDsn)
	if err != nil {
		cancelFunc = nil
		return nil, err
	}

	dbc := DBConnector{
		pool,
		dbCfg,
	}

	cancelFunc()
	return &dbc, nil
}

// NewAccount метод для создания нового экаунта из ДБ конектора
// вызывает методы объекта user.
// Проверяет есть ли такой пользователь.
// Если нет, то создает
func (dbc *DBConnector) NewAccount(user *model.User) error {
	//ctx := context.Background()
	//conn, err := dbc.Pool.Acquire(ctx)
	//if err != nil {
	//	return errs.ErrErrorServer
	//}
	//defer conn.Release()
	//
	//ctxVW := context.WithValue(ctx, model.KeyContext("data"), user)
	//pc := PgxpoolConn{conn}
	//
	//user.HashPassword = cryptography.HashSHA256(user.Password, dbc.Cfg.Key)
	//recordExists, err := pc.CheckExistence(ctxVW)
	//if err != nil {
	//	return errs.ErrErrorServer
	//}
	//
	//if recordExists {
	//	return errs.ErrLoginBusy
	//}
	//
	//if _, err = conn.Exec(ctx, constants.QueryInsertUserTemplate, user.Name, user.HashPassword); err != nil {
	//	return errs.ErrErrorServer
	//}

	return nil
}

// CheckAccount проверяет, существует ли пользователь в базе данных
func (dbc *DBConnector) CheckAccount(user *model.User) error {

	//ctx := context.Background()
	//conn, err := dbc.Pool.Acquire(ctx)
	//if err != nil {
	//	return errs.ErrErrorServer
	//}
	//defer conn.Release()
	//
	//user.HashPassword = cryptography.HashSHA256(user.Password, dbc.Cfg.Key)
	//ctxVW := context.WithValue(ctx, model.KeyContext("data"), user)
	//var pc = PgxpoolConn{Conn: conn}
	//recordExists, err := pc.CheckExistence(ctxVW)
	//if err != nil {
	//	return errs.InvalidFormat
	//}
	//if recordExists {
	//	return nil
	//
	//}
	//conn.Release()

	return errs.ErrInvalidLoginPassword
}

/////////////////////////////////////

// DelAccount удаляет пользователя по имени и хешированному паролю
func (dbc *DBConnector) DelAccount(user *model.User) error {
	//ctx := context.Background()
	//conn, err := dbc.Pool.Acquire(ctx)
	//if err != nil {
	//	return errs.ErrErrorServer
	//}
	//defer conn.Release()
	//
	//ctxVW := context.WithValue(ctx, model.KeyContext("data"), user)
	//pc := PgxpoolConn{conn}
	//
	//err = pc.Delete(ctxVW)
	//if err != nil {
	//	return errs.ErrErrorServer
	//}

	return nil
}

// Select выбирает объекты из базы данных
func (dbc *DBConnector) Select(ctx context.Context, t string) (model.Appender, error) {

	//user := ctx.Value(model.KeyContext("user"))
	//conn, err := dbc.Pool.Acquire(ctx)
	//if err != nil {
	//	return nil, errs.ErrErrorServer
	//}
	//defer conn.Release()
	//
	//strUser := user.(string)
	//na, err := model.NewAppender(t, strUser)
	//if err != nil {
	//	return nil, errs.ErrErrorServer
	//}
	//
	//ctxVW := context.WithValue(ctx, model.KeyContext("data"), na)
	//pc := PgxpoolConn{conn}
	//
	//a, err := pc.Select(ctxVW)
	//
	//if err != nil {
	//	return nil, errs.ErrErrorServer
	//}
	//

	a := model.Appender{}
	return a, nil
}

// Update добавляет/обновляет объекты базы данных
func (dbc *DBConnector) Update(u model.Updater) error {
	//ctx := context.Background()
	//conn, err := dbc.Pool.Acquire(ctx)
	//if err != nil {
	//	return errs.ErrErrorServer
	//}
	//defer conn.Release()
	//
	//ctxVW := context.WithValue(ctx, model.KeyContext("data"), u)
	//pc := PgxpoolConn{conn}
	//
	//recordExists, err := pc.CheckExistence(ctxVW)
	//
	//if err != nil {
	//	return errs.InvalidFormat
	//}
	//if recordExists {
	//
	//	if err = pc.Update(ctxVW); err != nil {
	//		return errs.InvalidFormat
	//	}
	//	return nil
	//}
	//
	//if err = pc.Insert(ctxVW); err != nil {
	//	return errs.InvalidFormat
	//}

	return nil
}

// Delete удаляет объекты из базы данных
func (dbc *DBConnector) Delete(u model.Updater) error {
	//ctx := context.Background()
	//conn, err := dbc.Pool.Acquire(ctx)
	//if err != nil {
	//	return errs.ErrErrorServer
	//}
	//defer conn.Release()
	//
	//ctxVW := context.WithValue(ctx, model.KeyContext("data"), u)
	//pc := PgxpoolConn{Conn: conn}
	//
	//err = pc.Delete(ctxVW)
	//if err != nil {
	//	return errs.ErrErrorServer
	//}
	return nil
}
