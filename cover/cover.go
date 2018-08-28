package cover

import (
	"encoding/json"
	"fmt"
	"os/exec"
	"strings"

	"github.com/seeleteam/e2e-tests/config"
	"github.com/seeleteam/e2e-tests/store"
)

// Run coverPath
func Run(coverPath string) (all string, specified map[string]string) {
	coverFileName := config.CoverFileName
	specified = make(map[string]string)
	// go test github.com/seeleteam/go-seele/... -coverprofile=seele_cover
	coverbyte, err := exec.Command("go", "test", coverPath, "-coverprofile="+coverFileName).CombinedOutput()
	if err != nil {
		return fmt.Sprintf("cover FAIL: %s %s", err, string(coverbyte)), nil
	}

	// remove useless output
	outs, pkgs := strings.Split(string(coverbyte), "\n"), strings.Split(config.CoverPackage, ",")
	for _, out := range outs {
		// ? == 63
		if out == "" || out[0] == 63 {
			continue
		}

		for _, pkg := range pkgs {
			if strings.Contains(out, pkg) {
				specified[pkg] = out
			}
		}

		all += out + "\n"
	}

	// go tool cover -html=covprofile -o coverage.html
	if err := exec.Command("go", "tool", "cover", "-html="+coverFileName, "-o", coverFileName+".html").Run(); err != nil {
		return fmt.Sprintf("tool cover FAIL: %s", err), nil
	}

	return all, specified
}

// PrintSpecifiedPkg coverage compare
func PrintSpecifiedPkg(yestoday string, specified map[string]string) string {
	result := "\n============= Change in coverage of major packages compared to yesterday ===============\n\n"
	yestodaySpec := make(map[string]string)
	_, _, coverByte := store.Get(yestoday)
	if err := json.Unmarshal(coverByte, &yestodaySpec); err != nil {
		return ""
	}

	for k, v := range specified {
		out, ok := yestodaySpec[k]
		if !ok {
			result += v + "\n"
		} else {
			result += out + " --> " + v[strings.Index(v, "coverage"):] + "\n"
		}
	}

	return result
}
