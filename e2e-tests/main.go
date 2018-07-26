package main

import (
	"fmt"
	"io/ioutil"
	"net/mail"
	"net/smtp"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	"github.com/scorredoira/email"
	"github.com/seeleteam/e2e-tests/store"
)

// config.go
const (
	SeelePath     = "github.com/seeleteam/go-seele"
	CoverFileName = "seele_coverage_detail"

	Subject    = "Daily Test Report"
	Sender     = "wangfeifan@zsbatech.com"
	Password   = "Wff19940326..."
	SenderName = "Seele-e2e"
	// Receivers  = "wangfeifan@zsbatech.com"
	Receivers = "rdc@zsbatech.com"
	Host      = "smtp.exmail.qq.com:25"

	StartHour = 04
	StartMin  = 00
	StartSec  = 00
)

var (
	attachFile = []string{}
)

func main() {
	// now := time.Now()
	// next := now.Add(time.Hour * 24)
	// DoTest(now.Format("20060102"), next.Format("20060102"))
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
		DoTest(now.Format("20060102"), next.Format("20060102"))
	}
}

// DoTest seele test
func DoTest(yesterday, today string) {
	if updateresult := updateSeele(); updateresult != "" {
		fmt.Println("updateresult:", updateresult)
		return
	}

	workPath := filepath.Join(SeelePath, "/...")
	fmt.Printf("date:%s workPath:%s\n", today, workPath)

	buildresult := build(workPath)
	fmt.Println("build done")
	coverresult := cover(workPath)
	fmt.Println("cover done")
	benchresult := bench(workPath)
	fmt.Println("bench done")
	store.Save(today, buildresult, coverresult, benchresult)

	message := ""
	if buildresult != "" || strings.Contains(coverresult, "FAIL") || strings.Contains(benchresult, "FAIL") {
		message += "😦 Appears to be a bug!\n\n"
	} else {
		message += "😁 Good day with no error~\n\n"
	}
	message += "\n=============Go build seele started. ===============\n" + buildresult
	message += "=============Go build seele completed. ===============\n\n"

	message += "\n=============Go cover seele started. ===============\n" + coverresult
	message += "=============Go cover seele completed. ===============\n\n"

	message += "\n=============Go bench seele started. ===============\n" + benchresult
	message += "=============Go bench seele completed. ===============\n\n"

	sendEmail(message, attachFile)
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

func updateSeele() string {
	if updateout, err := exec.Command("git", "pull").CombinedOutput(); err != nil {
		return fmt.Sprintf("update err: %s %s", err, string(updateout))
	}
	return ""
}

func build(buildPath string) string {
	// go build github.com/seeleteam/go-seele/...
	buildout, err := exec.Command("go", "build", buildPath).CombinedOutput()
	if err != nil {
		return fmt.Sprintf("build err: %s %s", err, string(buildout))
	}

	return ""
}

func cover(coverPath string) string {
	// go test github.com/seeleteam/go-seele/... -coverprofile=seele_cover
	coverout, err := exec.Command("go", "test", coverPath, "-coverprofile="+CoverFileName).CombinedOutput()
	if err != nil {
		return fmt.Sprintf("cover err: %s %s", err, string(coverout))
	}

	// go tool cover -html=covprofile -o coverage.html
	if err := exec.Command("go", "tool", "cover", "-html="+CoverFileName, "-o", CoverFileName+".html").Run(); err != nil {
		return fmt.Sprintf("tool cover err: %s\n", err)
	}

	attachFile = append(attachFile, CoverFileName+".html")
	return string(coverout)
}

func bench(benchPath string) string {
	benchout, walkPath := "", filepath.Join(os.Getenv("GOPATH"), "src", SeelePath)
	filepath.Walk(walkPath, func(path string, f os.FileInfo, err error) error {
		if strings.Contains(path, ".git") || strings.Contains(path, "vendor") || strings.Contains(path, "crypto") || !f.IsDir() {
			return nil
		}

		// go test github.com/seeleteam/go-seele/... -bench=. -run Benchmark
		path = path[strings.Index(path, "src")+4:]
		out, err := exec.Command("go", "test", path, "-bench=.", "-run", "Benchmark").CombinedOutput()
		if err != nil {
			if strings.Contains(string(out), "no Go files") {
				return nil
			}
			benchout += fmt.Sprintf("path: %s\nbench err: %s %s", path, err, string(out))
			return nil
		}

		benchout += string(out)
		if strings.Contains(string(out), "no test files") {
			return nil
		}

		// go test github.com/seeleteam/go-seele/core -bench=. -cpuprofile core.prof -run Benchmark
		ostype, cpuName := runtime.GOOS, ""
		if ostype == "windows" {
			cpuName = path[strings.LastIndex(path, "\\")+1:]
		} else {
			cpuName = path[strings.LastIndex(path, "/")+1:]
		}
		cpuout, err := exec.Command("go", "test", path, "-bench=.", "-cpuprofile", cpuName+".prof", "-run", "Benchmark").CombinedOutput()
		if err != nil {
			fmt.Println(fmt.Errorf("cpuout err: %s %s", err, string(cpuout)))
			return nil
		}

		// go tool pprof -svg core.prof  core.svg
		profout, err := exec.Command("go", "tool", "pprof", "-svg", cpuName+".prof", ">", cpuName+".svg").CombinedOutput()
		if err != nil {
			fmt.Println(fmt.Errorf("profout err: %s %s", err, string(profout)))
			return nil
		}

		if err := ioutil.WriteFile(cpuName+"_cpu_detail.svg", profout, os.ModePerm); err != nil {
			fmt.Printf("writefile err: %s\n", err)
			return nil
		}
		attachFile = append(attachFile, cpuName+"_cpu_detail.svg")
		return nil
	})

	return benchout
}

func sendEmail(message string, attachFile []string) {
	fmt.Println(message, attachFile)
	msg := email.NewMessage(Subject, message)
	msg.From, msg.To = mail.Address{Name: SenderName, Address: Sender}, strings.Split(Receivers, ";")
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
