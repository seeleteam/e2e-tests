/**
*  @file
*  @copyright defined in go-seele/LICENSE
 */

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
	today := time.Now()
	yesterday := today.Add(-1 * 24 * time.Hour)
	day := today.Weekday()

	// no work on Saturday and Sunday
	if day != time.Saturday && day != time.Sunday {
		DoTest(yesterday.Format("20060102"), today.Format("20060102"))
	}
}

// DoTest seele test
func DoTest(yesterday, today string) {
	attachFile = attachFile[:0]
	workPath := filepath.Join(config.SeelePath, "/...")
	fmt.Printf("yesterday:%s, today:%s, workPath:%s\n", yesterday, today, workPath)
	dir, err := os.Getwd()
	if err != nil {
		panic(fmt.Sprintf("Failed to get working directory %s", err))
	}

	fmt.Println("currrent working dir", dir)
	// build
	// buildresult := build(workPath)
	// fmt.Println("build done")
	// cover
	coverResult, specified := cover.Run(workPath)
	coverbyte, err := json.Marshal(specified)
	if err != nil {
		fmt.Println("Marshal specified FAIL")
	}
	fmt.Println("cover done")
	// bench
	benchresult := bench.Run(config.SeelePath)
	// the google have some mistake about the pprof formats
	compressd := bench.OutputCompressionReport("bench_reports.zip")
	fmt.Println("OutputCompressionReport result:\n", compressd)
	if compressd == "" {
		attachFile = append(attachFile, "bench_reports.zip")
	}
	fmt.Println("bench done")
	// save the result
	// store.Save(today, buildresult, benchresult, coverbyte)
	store.Save(today, benchresult, coverbyte)

	message := ""
	// if buildresult != "" || strings.Contains(coverResult, "FAIL") || strings.Contains(benchresult, "FAIL") {
	if strings.Contains(coverResult, "FAIL") || strings.Contains(benchresult, "FAIL") {
		message += "😦 Appears to be a bug!\n\n"
	} else {
		attachFile = append(attachFile, config.CoverFileName+".html")
		message += "😁 Good day with no error~\n\n"
	}

	// message += cover.PrintSpecifiedPkg(yesterday, specified)
	// message += "\n\n============= Go build seele completed. ===============\n" + buildresult
	message += "\n\n============= Go cover seele completed. ===============\n" + coverResult
	message += "\n\n============= Go bench seele completed. ===============\n" + benchresult

	sendEmail(message, attachFile)
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
	msg.Cc = strings.Split(config.CC, ";")
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
