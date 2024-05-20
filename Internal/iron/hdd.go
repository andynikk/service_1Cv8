package iron

import (
	"fmt"
	"github.com/StackExchange/wmi"
)

func GetDiskDrivers() (string, error) {
	type Win32_DiskDrive struct {
		Caption      string
		Name         string
		DeviceID     string
		Model        string
		Index        int
		Partitions   int
		Size         int
		PNPDeviceID  string
		Status       string
		SerialNumber string
		Manufacturer string
		MediaType    string
		Description  string
		SystemName   string
	}

	var dst []Win32_DiskDrive

	query := wmi.CreateQuery(&dst, "")
	err := wmi.Query(query, &dst)
	if err != nil {
		return "", err
	}

	var sn string
	for _, hdd := range dst {
		sn = fmt.Sprintf("%s-%s", sn, hdd.SerialNumber)
	}

	return sn, nil
}
