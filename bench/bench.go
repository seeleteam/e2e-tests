package bench

import (
	"archive/zip"
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/seeleteam/e2e-tests/config"
)

// Run bench
func Run(benchPath string) string {
	benchout, walkPath := "", filepath.Join(os.Getenv("GOPATH"), "src", benchPath)
	filepath.Walk(walkPath, func(path string, f os.FileInfo, err error) error {
		if strings.Contains(path, ".git") || strings.Contains(path, "vendor") || strings.Contains(path, "crypto") || f == nil || !f.IsDir() {
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

		if !strings.Contains(string(out), "Benchmark") {
			return nil
		}
		benchout += string(out)

		ostype, cpuName := runtime.GOOS, ""
		if ostype == "windows" {
			cpuName = path[strings.LastIndex(path, "\\")+1:]
		} else {
			cpuName = path[strings.LastIndex(path, "/")+1:]
		}

		// go test github.com/seeleteam/go-seele/core -bench=. -cpuprofile core.prof -run Benchmark
		cpuout, err := exec.Command("go", "test", path, "-bench=.", "-cpuprofile", cpuName+".prof", "-run", "Benchmark").CombinedOutput()
		if err != nil {
			fmt.Println(fmt.Errorf("cpuout err: %s %s", err, string(cpuout)))
			return nil
		}

		// go tool pprof core.prof | top 10
		in := bytes.NewBuffer(nil)
		cmd := exec.Command("go", "tool", "pprof", cpuName+".prof")
		cmd.Stdin = in
		in.WriteString("top " + config.BenchTopN)
		topN, err := cmd.CombinedOutput()
		if err != nil {
			fmt.Println(fmt.Errorf("topN err: %s %s", err, string(topN)))
			return nil
		}

		benchout += string(topN)[strings.Index(string(topN), "(pprof)"):] + "\n\n"

		// go tool pprof -png core.prof > core.png
		profout, err := exec.Command("go", "tool", "pprof", "-png", cpuName+".prof").CombinedOutput()
		if err != nil {
			fmt.Println(fmt.Errorf("profout err: %s %s", err, string(profout)))
			return nil
		}
		body := string(profout)
		googlebugs := "Main binary filename not available."
		if strings.Contains(body, googlebugs) {
			body = body[strings.Index(body, googlebugs)+len(googlebugs)+1:]
		}

		if err := ioutil.WriteFile(cpuName+"_cpu_detail.png", []byte(body), os.ModePerm); err != nil {
			fmt.Printf("writefile err: %s\n", err)
			return nil
		}

		return nil
	})

	return benchout
}

// OutputCompressionReport output bench zip reports
func OutputCompressionReport(zipName string) string {
	// compressPkg := "./reports"
	// if err := os.MkdirAll(compressPkg, os.ModePerm); err != nil {
	// 	return fmt.Sprintf("MkdirAll compress package[%s] FAIL: %s", compressPkg, err)
	// }

	file, err := os.Create(zipName)
	if err != nil {
		return fmt.Sprintf("Create reports.zip FAIL: %s", err)
	}

	writer := zip.NewWriter(file)
	defer writer.Close()

	if err := filepath.Walk(".", func(path string, f os.FileInfo, err error) error {
		if f != nil && strings.Contains(f.Name(), ".png") {
			body, err := ioutil.ReadFile(path)
			if err != nil {
				return fmt.Errorf("ReadFile fileInfo[%s] FAIL: %s", path, err)
			}

			fw, err := writer.Create(path)
			if err != nil {
				return fmt.Errorf("Create fileInfo[%s] FAIL: %s", path, err)
			}

			if _, err := fw.Write(body); err != nil {
				return fmt.Errorf("Write body[%s] FAIL: %s", string(body), err)
			}
		}

		return nil
	}); err != nil {
		return err.Error()
	}

	return ""
}
