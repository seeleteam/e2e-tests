package main

import (
	"encoding/json"
	"fmt"
	"net/mail"
	"net/smtp"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/scorredoira/email"
	"github.com/seeleteam/e2e-tests/bench"
	"github.com/seeleteam/e2e-tests/config"
	"github.com/seeleteam/e2e-tests/cover"
	"github.com/seeleteam/e2e-tests/store"
)

var (
	attachFile = []string{}
)

func main() {
	// now := time.Now()
	// next := now.Add(time.Hour * 24)
	// DoTest(now.Format("20060102"), next.Format("20060102"))
	yesterday := time.Now()
	for {
		now := time.Now()
		next := now.Add(time.Hour * 24)
		next = time.Date(next.Year(), next.Month(), next.Day(), config.StartHour, config.StartMin, config.StartSec, 0, next.Location())
		fmt.Println("now:", now)
		fmt.Println("next:", next)
		t := time.NewTimer(next.Sub(now))
		<-t.C
		t.Stop()
		weekday := next.Weekday()
		if weekday != time.Saturday && weekday != time.Sunday {
			fmt.Println("Go")
			DoTest(yesterday.Format("20060102"), next.Format("20060102"))
			yesterday = next
		}
	}
}

// DoTest seele test
func DoTest(yesterday, today string) {
	if updateresult := updateSeele(); updateresult != "" {
		fmt.Println("updateresult:", updateresult)
		return
	}

	workPath := filepath.Join(config.SeelePath, "/...")
	fmt.Printf("date:%s workPath:%s\n", today, workPath)
	// build
	buildresult := build(workPath)
	fmt.Println("build done")
	// cover
	coverResult, specified := cover.Run(workPath)
	coverbyte, err := json.Marshal(specified)
	if err != nil {
		fmt.Println("Marshal specified FAIL")
	}
	fmt.Println("cover done")
	// bench
	benchresult := bench.Run(config.SeelePath)
	fmt.Println("OutputCompressionReport result:\n", bench.OutputCompressionReport("bench_reports.zip"))
	fmt.Println("bench done")
	// save the result
	store.Save(today, buildresult, benchresult, coverbyte)

	message := ""
	if buildresult != "" || strings.Contains(coverResult, "FAIL") || strings.Contains(benchresult, "FAIL") {
		message += "😦 Appears to be a bug!\n\n"
	} else {
		attachFile = append(attachFile, config.CoverFileName+".html", "bench_reports.zip")
		message += "😁 Good day with no error~\n\n"
	}

	message += cover.PrintSpecifiedPkg(yesterday, specified)
	message += "\n\n============= Go build seele completed. ===============\n" + buildresult
	message += "\n\n============= Go cover seele completed. ===============\n" + coverResult
	message += "\n\n============= Go bench seele completed. ===============\n" + benchresult

	sendEmail(message, attachFile)
	filepath.Walk(".", func(path string, f os.FileInfo, err error) error {
		if strings.Contains(path, "main") || path == "." {
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

func sendEmail(message string, attachFile []string) {
	fmt.Println(message, attachFile)
	msg := email.NewMessage(config.Subject, message)
	msg.From, msg.To = mail.Address{Name: config.SenderName, Address: config.Sender}, strings.Split(config.Receivers, ";")
	for _, filePath := range attachFile {
		if err := msg.Attach(filePath); err != nil {
			fmt.Printf("failed to add attach file. path: %s, err: %s\n", filePath, err)
		}
	}

	hp := strings.Split(config.Host, ":")
	auth := smtp.PlainAuth("", config.Sender, config.Password, hp[0])

	if err := email.Send(config.Host, auth, msg); err != nil {
		fmt.Println("failed to send mail. err:", err)
	}
}
