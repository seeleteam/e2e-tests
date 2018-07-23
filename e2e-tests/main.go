package main

import (
	"fmt"
	"io/ioutil"
	"net/mail"
	"net/smtp"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/scorredoira/email"
)

// config.go
const (
	SeelePath     = "github.com/seeleteam/go-seele"
	CoverFileName = "seele_coverage_detail"

	Subject    = "Daily Test Report"
	Sender     = "wangfeifan@zsbatech.com"
	Password   = "Wff19940326..."
	SenderName = "Seele-e2e"
	// Receivers  = "wangfeifan@zsbatech.com;wangfeifan@zsbatech.com"
	Receivers = "liuwensi@zsbatech.com;jiazhiwei@zsbatech.com;libingyi@zsbatech.com;qiaozhigang@zsbatech.com;qiubo@zsbatech.com;tianyu@zsbatech.com;wangfeifan@zsbatech.com;yanpeng@zsbatech.com;yanpengfei@zsbatech.com;zhoushengkun@zsbatech.com;yangmingyan@zsbatech.com;"
	Host      = "smtp.exmail.qq.com:25"

	StartHour = 04
	StartMin  = 00
	StartSec  = 00
)

var (
	attachFile = []string{CoverFileName + ".html"}
)

func main() {
	for {
		now := time.Now()
		next := now.Add(time.Hour * 24)
		next = time.Date(next.Year(), next.Month(), next.Day(), StartHour, StartMin, StartSec, 0, next.Location())
		fmt.Println("now:", now)
		fmt.Println("next:", next)
		t := time.NewTimer(next.Sub(now))
		<-t.C
		t.Stop()
		fmt.Println("Go")
		DoTest()
	}

}

// DoTest seele test
func DoTest() {
	message := "嘿嘿😁, 测一下行不行啊\n\n"

	workPath := filepath.Join(SeelePath, "/...")
	message += build(workPath)
	message += "=============Go build seele completed. ===============\n\n"

	message += cover(workPath)
	message += "=============Go cover seele completed. ===============\n\n"

	message += bench(workPath)
	message += "=============Go bench seele completed. ===============\n\n"

	sendEmail(message, attachFile)

	fmt.Println(message, attachFile)
	filepath.Walk(".", func(path string, f os.FileInfo, err error) error {
		if strings.Contains(path, "main.go") || path == "." {
			return nil
		}

		fmt.Println("remove path:", path)
		if err := os.Remove(path); err != nil {
			fmt.Println("remove err:", err)
		}
		return nil
	})
}

func build(buildPath string) string {
	// go build github.com/seeleteam/go-seele/...
	buildout, err := exec.Command("go", "build", buildPath).Output()
	if err != nil {
		return fmt.Sprintf("build err: %s\n%s", err, string(buildout))
	}

	return ""
}

func cover(coverPath string) string {
	// go test github.com/seeleteam/go-seele/... -coverprofile=seele_cover
	coverout, err := exec.Command("go", "test", coverPath, "-coverprofile="+CoverFileName).Output()
	if err != nil {
		return fmt.Sprintf("cover err: %s\n%s", err, string(coverout))
	}

	// go tool cover -html=covprofile -o coverage.html
	if _, err := exec.Command("go", "tool", "cover", "-html="+CoverFileName, "-o", CoverFileName+".html").Output(); err != nil {
		return fmt.Sprintf("tool cover err: %s\n", err)
	}

	return string(coverout)
}

func bench(benchPath string) string {
	// go test github.com/seeleteam/go-seele/... -bench=.
	benchout, err := exec.Command("go", "test", benchPath, "-bench=.").Output()
	if err != nil {
		return fmt.Sprintf("bench err: %s\n%s", err, string(benchout))
	}

	walkPath := filepath.Join(os.Getenv("GOPATH"), "src", SeelePath)
	filepath.Walk(walkPath, func(path string, f os.FileInfo, err error) error {
		if strings.Contains(path, ".git") || strings.Contains(path, "cmd") || strings.Contains(path, "vendor") || !f.IsDir() {
			return nil
		}

		// go test github.com/seeleteam/go-seele/core -bench=. -cpuprofile core.prof
		path = path[strings.Index(path, "src")+4:]
		cpuName := path[strings.LastIndex(path, "\\")+1:]
		cpuout, err := exec.Command("go", "test", path, "-bench=.", "-cpuprofile", cpuName+".prof").Output()
		if err != nil {
			fmt.Println(fmt.Errorf("cpuout err: %s\n%s", err, string(cpuout)))
			return nil
		}

		if strings.Contains(string(cpuout), "no test files") {
			return nil
		}

		// go tool pprof -svg core.prof  core.svg
		profout, err := exec.Command("go", "tool", "pprof", "-svg", cpuName+".prof", "", cpuName+".svg").Output()
		if err != nil {
			fmt.Println(fmt.Errorf("profout err: %s\n%s", err, string(profout)))
			return nil
		}

		// fmt.Println("profout:", string(profout))

		if err := ioutil.WriteFile(cpuName+"_cpu_detail.svg", profout, os.ModePerm); err != nil {
			fmt.Printf("writefile err: %s\n", err)
			return nil
		}
		attachFile = append(attachFile, cpuName+"_cpu_detail.svg")
		return nil
	})

	return string(benchout)
}

func sendEmail(message string, attachFile []string) {
	msg := email.NewMessage(Subject, message)
	msg.From = mail.Address{Name: SenderName, Address: Sender}
	msg.To = strings.Split(Receivers, ";")
	for _, filePath := range attachFile {
		if err := msg.Attach(filePath); err != nil {
			fmt.Printf("failed to add attach file. path: %s, err: %s\n", filePath, err)
		}
	}

	hp := strings.Split(Host, ":")
	auth := smtp.PlainAuth("", Sender, Password, hp[0])

	if err := email.Send(Host, auth, msg); err != nil {
		fmt.Println("failed to send mail. err:", err)
	}
}
