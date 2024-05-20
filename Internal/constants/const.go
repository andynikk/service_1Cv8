package constants

type Status int
type EventService int

const Divisor = 2

const (
	FormMain Status = iota
	FormAddDB
	FormBlock
	FormServer
	FormLogin
	FormSetting
	FormAddService
	FormServerDB
	FormSQLDB
	FormAddSQLSrv
	FormCopyFile
	FormResTlg
	FormTgMsg
)

const (
	EvStop EventService = iota
	EvDelete
	EvStart
	EvReboot
)

const (
	FormExchangeBasic Status = iota
	FormExchangeListToken
	FormExchangeKeyToken
	FormExchangeListQueues
	FormExchangeQueues
)

const SecretKey string = "s[pdlsda[pd" +
	"sdefds" +
	"fdsfewdsf wertre rt" +
	"egrt trgrt"

const (
	Blocking Status = iota
	Services
)

const (
	IntervalService = 0
	KillProcessKb   = 0
	Port            = "8080"
)

const HashKey = "aaadd"

const (
	Shag = 8388672 //524288000 //33554688 //8388672 //16777344 //4096
)

const (
	MASSAGE = "Планируется техническое обслуживание базы с %s до %s\n\n" +
		"Для администратора:\nДля того чтобы разрешить работу пользователей, воспользуйтесь консолью кластера " +
		"серверов или запустите \"1С:Предприятие\" с параметрами:\nENTERPRISE /S\"%s" +
		"\\%s\"/CРазрешитьРаботуПользователей /UC<код разрешения>"
)

const (
	PudgeSetting            = "./db/config/settings"
	PudgeDataBases          = "./db/config/databases"
	PudgePropertyDB         = "./db/config/property_db"
	PudgeServices           = "./db/config/services"
	PudgeBasesDoubleControl = "./db/config/bases_double_control"

	PudgeTokens = "./db/config/tokens"
	PudgeKey    = "./db/config/keys"

	Pudge = "./db/pudge"

	NumberPriorities   = 100
	TotalByteInPudge   = 104857600
	TotalCount         = 100000
	TimeTransferToDisk = 50
)

const (
	TypeEncryption = "sha512"

	TimeLivingCertificateYaer   = 10
	TimeLivingCertificateMounth = 0
	TimeLivingCe5rtificateDay   = 0
	TimeLiveToken               = 7
)

const (
	Soft TypeStorage = iota
	Hard
)

type TypeStorage int

func (st TypeStorage) String() string {
	return [...]string{"soft", "hard"}[st]
}
