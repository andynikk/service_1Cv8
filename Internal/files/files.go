package files

import (
	"Service_1Cv8/internal/constants"
	"bytes"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"strings"
)

type VC struct {
	b   []byte
	pos int64
}

type Files struct {
	Name        string
	Size        int64
	SizeStrings string
}

func ListFile(path string) ([]Files, error) {
	var listFiles []Files

	entries, err := os.ReadDir(path)
	if err != nil {
		return nil, err
	}

	for _, e := range entries {
		if e.IsDir() {
			continue
		}

		infoFile, _ := e.Info()
		file := Files{
			Name:        e.Name(),
			Size:        infoFile.Size(),
			SizeStrings: GroupSeparator(fmt.Sprintf("%d", infoFile.Size())),
		}

		listFiles = append(listFiles, file)
	}

	return listFiles, nil
}

func ReadFile(pathSource string, chanOut chan VC, chanInfo chan int64) {
	var pos int64 = 0

	file, err := os.Open(pathSource)
	if err != nil {
		close(chanOut)
		close(chanInfo)

		return
	}
	defer file.Close()

	fileInfo, err := file.Stat()
	balance := fileInfo.Size()

	for {

		vShag := min(constants.Shag, balance)

		mVc := new(VC)
		mVc.b = make([]byte, vShag)
		mVc.pos = pos

		_, err = file.Read(mVc.b)
		if err != nil {
			if err == io.EOF {
				break
			}
			log.Fatal(err)
		}

		chanOut <- *mVc

		pos = pos + vShag
		balance = balance - vShag

		chanInfo <- pos

		if balance <= 0 {
			break
		}
	}

	close(chanOut)
	close(chanInfo)
}

func GoWriteFile(newFile *os.File, chanOut chan VC) {

	for {
		vs, ok := <-chanOut
		if !ok {

			break
		}

		WriteFile(newFile, &vs)
	}

}

func WriteFile(newFile *os.File, vs *VC) {

	if _, err := newFile.WriteAt(vs.b, vs.pos); err != nil {
		log.Fatal(err)
	}

}

func GroupSeparator(source string) string {

	var strs []string

	runes := []rune(source)

	lenSource := len(source)
	numberTriples := lenSource/3 + 1

	sch := 0
	for i := 1; i <= numberTriples; i++ {

		beginning := len(runes) - sch - 3
		if beginning < 0 {
			beginning = 0
		}

		triple := string(runes[beginning : len(runes)-sch])
		if triple == "" {
			continue
		}

		strs = append(strs, triple)
		sch = sch + 3
	}

	strReturn := ""
	for i := 1; i <= len(strs); i++ {
		if strReturn != "" {
			strReturn = strReturn + " "
		}
		strReturn = strReturn + strs[len(strs)-i]
	}

	return strReturn
}

func DelFolder(n string) error {
	path := fmt.Sprintf("%s/%s", constants.Pudge, n)
	//if isOSWindows() {
	//	path = strings.Replace(path, "/", "\\", -1)
	//}

	err := os.RemoveAll(path)
	if err != nil {
		return err
	}

	return nil
}

func isOSWindows() bool {

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

func min(a, b int64) int64 {
	if a < b {
		return a
	}
	return b
}
