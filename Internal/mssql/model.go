package mssql

import (
	mssql "github.com/microsoft/go-mssqldb"
)

type ConnectSQLSetting struct {
	Server   string
	User     string
	Password string
	Database string
}

type DB struct {
	Mark          string
	ID            int
	Name          string
	RecoveryModel string
	State         string
}

type ConfigDB struct {
	FileName   string
	Creation   mssql.DateTime1
	Modified   mssql.DateTime1
	Attributes int
	DataSize   int64
	BinaryData []byte
	PartNo     int
}
