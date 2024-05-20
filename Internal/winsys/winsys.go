package winsys

import (
	"bytes"
	"errors"
	"fmt"
	"golang.org/x/sys/windows"
	"golang.org/x/sys/windows/svc/mgr"
	"golang.org/x/text/encoding/charmap"
	"log"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"syscall"
	"time"
	"unsafe"

	"Service_1Cv8/internal/files"
	"Service_1Cv8/internal/repository"
)

func RebootWindows() error {

	user32 := syscall.MustLoadDLL("user32")
	defer user32.Release()

	kernel32 := syscall.MustLoadDLL("kernel32")
	defer user32.Release()

	advapi32 := syscall.MustLoadDLL("advapi32")
	defer advapi32.Release()

	ExitWindowsEx := user32.MustFindProc("ExitWindowsEx")
	GetCurrentProcess := kernel32.MustFindProc("GetCurrentProcess")
	GetLastError := kernel32.MustFindProc("GetLastError")
	OpenProdcessToken := advapi32.MustFindProc("OpenProcessToken")
	LookupPrivilegeValue := advapi32.MustFindProc("LookupPrivilegeValueW")
	AdjustTokenPrivileges := advapi32.MustFindProc("AdjustTokenPrivileges")

	currentProcess, _, _ := GetCurrentProcess.Call()

	const tokenAdjustPrivileges = 0x0020
	const tokenQuery = 0x0008
	var hToken uintptr

	result, _, err := OpenProdcessToken.Call(currentProcess, tokenAdjustPrivileges|tokenQuery, uintptr(unsafe.Pointer(&hToken)))
	if result != 1 {
		return err
	}

	const SeShutdownName = "SeShutdownPrivilege"

	type Luid struct {
		lowPart  uint32 // DWORD
		highPart int32  // long
	}
	type LuidAndAttributes struct {
		luid       Luid   // LUID
		attributes uint32 // DWORD
	}

	type TokenPrivileges struct {
		privilegeCount uint32 // DWORD
		privileges     [1]LuidAndAttributes
	}

	var tkp TokenPrivileges

	result, _, err = LookupPrivilegeValue.Call(uintptr(0), uintptr(unsafe.Pointer(syscall.StringToUTF16Ptr(SeShutdownName))), uintptr(unsafe.Pointer(&(tkp.privileges[0].luid))))
	if result != 1 {
		return err
	}

	const SePrivilegeEnabled uint32 = 0x00000002

	tkp.privilegeCount = 1
	tkp.privileges[0].attributes = SePrivilegeEnabled

	result, _, err = AdjustTokenPrivileges.Call(hToken, 0, uintptr(unsafe.Pointer(&tkp)), 0, uintptr(0), 0)
	if result != 1 {
		return err
	}

	result, _, _ = GetLastError.Call()
	if result != 0 {
		return err
	}

	const ewxForceIfHung = 0x00000010
	const ewxReboot = 0x00000002
	const ewxShutdown = 0x00000001
	const shutdownReasonMajorSoftware = 0x00030000

	result, _, err = ExitWindowsEx.Call(ewxReboot|ewxForceIfHung, shutdownReasonMajorSoftware)
	if result != 1 {
		return err
	}

	return nil
}

func RebootRemoteWindows(workspace string) error {

	//arg1 := fmt.Sprintf("shutdown /m \\\\%s /r /f", workspace)
	//arg1 := "shutdown /m \\\\" + workspace + " /r /f"
	aa := "/m \\\\" + workspace
	cmd := exec.Command("shutdown",
		aa, "/r", "/f")

	//arg1 := fmt.Sprintf("/m \\\\%s /r /f", workspace)
	//cmd := exec.Command("shutdown", arg1)

	err := cmd.Start()
	//err := cmd.Run()
	return err

}

func StopService(computerName, serviceName string) error {

	// Подключаемся к удаленному компьютеру
	handle, err := windows.OpenSCManager(windows.StringToUTF16Ptr(computerName), nil, windows.SC_MANAGER_CONNECT)
	if err != nil {
		return errors.New(fmt.Sprintf("%s: %s", "1", err.Error()))
	}
	defer windows.CloseServiceHandle(handle)

	// Открываем нужную службу
	serviceHandle, err := windows.OpenService(handle, windows.StringToUTF16Ptr(serviceName), windows.SERVICE_STOP|windows.SERVICE_QUERY_STATUS)
	if err != nil {
		return errors.New(fmt.Sprintf("%s: %s", "2", err.Error()))
	}
	defer windows.CloseServiceHandle(serviceHandle)

	// Отправляем службе команду на остановку
	var serviceStatus windows.SERVICE_STATUS
	if err := windows.ControlService(serviceHandle, windows.SERVICE_CONTROL_STOP, &serviceStatus); err != nil {
		return errors.New(fmt.Sprintf("%s: %s", "3", err.Error()))
		return err
	}

	return nil
}

func ClearServerCache(computerName string) error {
	path := fmt.Sprintf("\\\\%s\\c$\\Program Files\\1cv8\\srvinfo", computerName)
	entries, err := os.ReadDir(path)
	if err != nil {
		return err
	}
	for _, e := range entries {

		if !e.IsDir() || e.Name()[:4] != "reg_" {
			continue
		}

		podEntries, err := os.ReadDir(fmt.Sprintf("%s\\%s", path, e.Name()))
		if err != nil {
			return err
		}
		for _, pe := range podEntries {
			if !pe.IsDir() || len(pe.Name()) < 7 || pe.Name()[:7] != "snccntx" {
				continue
			}

			files, err := os.ReadDir(fmt.Sprintf("%s\\%s\\%s\\", path, e.Name(), pe.Name()))
			if err != nil {
				return err
			}
			for _, f := range files {
				if f.IsDir() {
					err = os.RemoveAll(fmt.Sprintf("%s\\%s\\%s\\%s", path, e.Name(), pe.Name(), f.Name()))
					continue
				}
				err = os.Remove(fmt.Sprintf("%s\\%s\\%s\\%s", path, e.Name(), pe.Name(), f.Name()))
				if err != nil {
					return err
				}
			}
		}
	}

	return nil
}

func StartService(IP, computerName, serviceName string) error {

	m, err := mgr.ConnectRemote(computerName)
	if err != nil {
		return err
	}
	defer m.Disconnect()

	s, err := m.OpenService(serviceName)
	if err != nil {
		return err
	}
	defer s.Close()

	err = s.Start(serviceName)
	if err != nil {
		return err
	}

	return nil
}

func IsOSWindows() bool {

	var stderr bytes.Buffer
	defer stderr.Reset()

	var out bytes.Buffer
	defer out.Reset()

	cmd := exec.Command("cmd", "ver")
	cmd.Stdin = strings.NewReader("some input")
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	err := cmd.Run()
	if err != nil {
		return false
	}
	myOS := out.String()
	return strings.Contains(myOS, "Microsoft Windows")
}

func Convert_CP866_To_unicode(b byte) rune {
	r := charmap.CodePage866.DecodeByte(b)
	return r
	//for key, val := range charset.Table_CP866 {
	//	if int32(key) == v {
	//		return []rune{val}, true
	//	}
	//}
}

func SubstitutionRune(r []rune) ([]rune, bool) {

	changed := false
	for k, v := range r {

		//unicode.Is(unicode.Cyrillic, v)
		//b, ok := charmap.CodePage866.EncodeRune(v)
		//for key, val := range charset.Table_CP866 {
		//	if int32(key) == v {
		//		return []rune{val}, true
		//	}
		//}

		if v >= 128 && v <= 159 {
			r[k] = v + 912
			changed = true
		} else if v >= 160 && v <= 175 {
			r[k] = v + 912
			changed = true
		} else if v >= 224 && v <= 239 {
			r[k] = v + 864
			changed = true
		} else if v == 240 {
			r[k] = 1025
			changed = true
		} else if v == 241 {
			r[k] = 1105
			changed = true
		}
	}

	return r, changed
}

func KillWinProc(hostname, nameProc string, sizeKill int) ([]repository.ClosedTask, error) {

	log.Println("Start kill win proc")
	if hostname == "" {
		hostname, _ = os.Hostname()
	}

	//cmd := exec.Command("tasklist.exe",
	//	"/fo", "csv", "/nh", "/s", "1c1.telematika.local", // "/U", "telematika\\andrey.mikhailov", "/P", "cNjgh005",
	//	"/fi", fmt.Sprintf("IMAGENAME eq %s", nameProc), "/fi", "MEMUSAGE ge 12000000")

	cmd := exec.Command("tasklist.exe",
		"/fo", "csv", "/nh", "/fi", fmt.Sprintf("IMAGENAME eq %s", nameProc), "/fi", "MEMUSAGE ge 12000000")
	cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}
	out, err := cmd.Output()
	if err != nil {
		log.Println(3, err)
		return nil, err
	}

	log.Println(fmt.Sprintf("Task items %d", len(out)))

	var cts []repository.ClosedTask

	strOut := string(out)
	bukva := ""
	stroka := ""
	for _, v := range strOut {
		if v == 65533 || v == 13 {
			continue
		}

		bukva = string(v)
		if bukva == "\n" {
			bukva = ""

			a := strings.FieldsFunc(stroka, Split)
			a[0] = strings.Replace(a[0], "\"", "", -1)

			if nameProc != a[0] {
				stroka = ""
				continue
			}

			strSize := a[len(a)-1]
			strSize = strings.Replace(strSize, "K", "", -1)
			strSize = strings.Replace(strSize, "\"", "", -1)
			strSize = strings.Replace(strSize, " ", "", -1)

			size, err := strconv.Atoi(strSize)
			if err != nil {
				stroka = ""
				log.Println(1, err)
				continue
			}

			log.Println(fmt.Sprintf("Size task %d", size))
			if size >= sizeKill {
				pid := strings.Replace(a[1], "\"", "", -1)
				//cmd = exec.Command("taskkill", "/s", hostname, "/Pid", pid, "/f", "/t")
				cmd = exec.Command("taskkill", "/Pid", pid, "/f", "/t")

				ct := repository.ClosedTask{
					Time: time.Now(),
					Host: hostname,
					Size: files.GroupSeparator(strSize),
				}

				cts = append(cts, ct)
				err = cmd.Start()
				if err != nil {
					stroka = ""
					log.Println(2, err)
					continue
				}
			}
			stroka = ""
		}

		stroka = stroka + bukva
	}

	return cts, nil
}

func Split(r rune) bool {
	return r == ','
}
