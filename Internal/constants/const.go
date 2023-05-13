package constants

type Status int

const Divisor = 2

const (
	FormMain Status = iota
	FormAddDB
	FormPropertyDB
	FormServer
	FormLogin
	FormSetting
	FormAddService
	FormServerDB
	FormSQLDB
	FormAddSQLSrv
	FormCopyFile
)

const (
	Blocking Status = iota
	Services
)

const (
	IntervalService = 15
	KillProcessKb   = 12000000
	Port            = "8080"
)

const (
	Shag = 8388672 //524288000 //33554688 //8388672 //16777344 //4096
)

const (
	MASSAGE = "Планируется техническое обслуживание базы с %s до %s\n\n" +
		"Для администратора:\nДля того чтобы разрешить работу пользователей, воспользуйтесь консолью кластера " +
		"серверов или запустите \"1С:Предприятие\" с параметрами:\nENTERPRISE /S\"%s" +
		"\\%s\"/CРазрешитьРаботуПользователей /UC<код разрешения>"
)

var HashKey = "asddesflplkpofef fck=dfdw?Wdwk dow98933"
