package OneCv8

import (
	"fmt"
	"os/exec"
)

type ConfigDB struct {
	Command string
	Server  string
	Port    string
	DB      string
	Key     string
	User    string
	Pwd     string
}

func OpenConfigDB(openConfigDB ConfigDB) error {
	//command := "C:\\Program Files\\1cv8\\8.3.20.2257\\bin\\1cv8.exe"
	port := openConfigDB.Port
	if port != "" {
		port = ":" + port
	}
	arg1 := fmt.Sprintf("DESIGNER /S %s%s\\%s /WA- /N%s /P%s /UC%s", openConfigDB.Server,
		port, openConfigDB.DB, openConfigDB.User, openConfigDB.Pwd, openConfigDB.Key)

	cmd := exec.Command(openConfigDB.Command, arg1)
	err := cmd.Start()

	return err
}
