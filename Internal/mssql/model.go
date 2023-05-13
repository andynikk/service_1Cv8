package mssql

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
