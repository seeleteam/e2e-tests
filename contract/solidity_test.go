/**
*  @file
*  @copyright defined in go-seele/LICENSE
 */

package contract

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"strings"
	"testing"
)

func Test_SolidityFiles(t *testing.T) {
	// ensure the simulator tool exists.
	exeFile := `..\..\go-seele\cmd\vm\vm.exe`
	if _, err := os.Stat(exeFile); os.IsNotExist(err) {
		t.Fatalf("Cannot find the solidity simulator: %v", exeFile)
	}

	// ensure the config file exists.
	configFile := "solidity_test.config"
	config, err := ioutil.ReadFile(configFile)
	if err != nil {
		t.Fatalf("Cannot find the solidity test config file: %v", configFile)
	}

	// execute commands and validate outputs in config file.
	cmdPrefix := "vm "
	commentPrefix1 := "#"
	commentPrefix2 := "//"
	outputs := make(map[string]bool)
	for i, line := range strings.Split(string(config), "\n") {
		// ignore empty line and comment line
		if line = strings.TrimSpace(line); len(line) == 0 || strings.HasPrefix(line, commentPrefix1) || strings.HasPrefix(line, commentPrefix2) {
			continue
		}

		if strings.HasPrefix(line, cmdPrefix) {
			outputs = execute(t, exeFile, line[len(cmdPrefix):])
		} else if !outputs[line] {
			fmt.Println("outputs:", outputs)
			t.Fatalf("Failed to assert output, line = %v", i+1)
		}
	}
}

// execute executes the specified command and returns the outputs.
func execute(t *testing.T, exeFile, args string) map[string]bool {
	cmd := exec.Command(exeFile, strings.Split(args, " ")...)
	output, err := cmd.Output()
	if err != nil {
		fmt.Println("Failed to execute simulator:")
		fmt.Println("\tcommand: vm", args)
		fmt.Println("\terror  :", err.Error())
		t.Fatal()
	}

	lines := make(map[string]bool)
	for _, l := range strings.Split(string(output), "\n") {
		if l = strings.TrimSpace(l); len(l) == 0 {
			continue
		}

		/* if lines[l] {
			t.Fatalf("Found 2 same lines in simulator execution outputs: %v", l)
		} */

		lines[l] = true
	}

	return lines
}
