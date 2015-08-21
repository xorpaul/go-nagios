package nagios

import (
	"fmt"
	"github.com/kballard/go-shellquote"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"syscall"
	"time"
)

var (
	debug     bool
	buildtime string
)

type NagiosResult struct {
	exitCode  int
	text      string
	perfdata  string
	multiline []string
}

type ExecResult struct {
	returnCode int
	output     string
}

// Debugf is a helper function for debug logging if mainCfgSection["debug"] is set
func Debugf(s string) {
	if debug != false {
		fmt.Println("DEBUG " + fmt.Sprint(s))
	}
}

// nagiosExit uses the NagiosResult struct to output Nagios plugin compatible output and exit codes
func nagiosExit(nr NagiosResult) {
	text := nr.text
	exitCode := nr.exitCode
	switch {
	case nr.exitCode == 0:
		text = "OK: " + nr.text
		exitCode = nr.exitCode
	case nr.exitCode == 1:
		text = "WARNING: " + nr.text
		exitCode = nr.exitCode
	case nr.exitCode == 2:
		text = "CRITICAL: " + nr.text
		exitCode = nr.exitCode
	case nr.exitCode == 3:
		text = "UNKNOWN: " + nr.text
		exitCode = nr.exitCode
	default:
		text = "UNKNOWN: Exit code '" + string(nr.exitCode) + "'undefined :" + nr.text
		exitCode = 3
	}

	if len(nr.multiline) > 0 {
		multiline := ""
		for _, l := range nr.multiline {
			multiline = multiline + l + "\n"
		}
		fmt.Printf("%s|%s\n%s\n", text, nr.perfdata, multiline)
	} else {
		fmt.Printf("%s|%s\n", text, nr.perfdata)
	}
	os.Exit(exitCode)
}

func executeCommand(command string, timeout int, allowFail bool) ExecResult {
	Debugf("Executing " + command)
	parts := strings.SplitN(command, " ", 2)
	cmd := parts[0]
	cmdArgs := []string{}
	if len(parts) > 1 {
		args, err := shellquote.Split(parts[1])
		if err != nil {
			Debugf("executeCommand(): err: " + fmt.Sprint(err))
			os.Exit(1)
		} else {
			cmdArgs = args
		}
	}

	before := time.Now()
	out, err := exec.Command(cmd, cmdArgs...).CombinedOutput()
	duration := time.Since(before).Seconds()
	er := ExecResult{0, string(out)}
	if msg, ok := err.(*exec.ExitError); ok { // there is error code
		er.returnCode = msg.Sys().(syscall.WaitStatus).ExitStatus()
	}
	Debugf("Executing " + command + " took " + strconv.FormatFloat(duration, 'f', 5, 64) + "s")
	if err != nil && !allowFail {
		fmt.Println("executeCommand(): command failed: "+command, err)
		fmt.Println("executeCommand(): Output: " + string(out))
		os.Exit(1)
	}
	return er
}
