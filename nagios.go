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
	Debug     bool
	buildtime string
)

type NagiosResult struct {
	ExitCode  int
	Text      string
	Perfdata  string
	Multiline []string
}

type ExecResult struct {
	ReturnCode int
	Output     string
}

// Debugf is a helper function for debug logging if mainCfgSection["debug"] is set
func Debugf(s string) {
	if Debug != false {
		fmt.Println("DEBUG " + fmt.Sprint(s))
	}
}

// NagiosExit uses the NagiosResult struct to output Nagios plugin compatible output and exit codes
func NagiosExit(nr NagiosResult) {
	text := nr.Text
	exitCode := nr.ExitCode
	switch {
	case nr.ExitCode == 0:
		text = "OK: " + nr.Text
		exitCode = nr.ExitCode
	case nr.ExitCode == 1:
		text = "WARNING: " + nr.Text
		exitCode = nr.ExitCode
	case nr.ExitCode == 2:
		text = "CRITICAL: " + nr.Text
		exitCode = nr.ExitCode
	case nr.ExitCode == 3:
		text = "UNKNOWN: " + nr.Text
		exitCode = nr.ExitCode
	default:
		text = "UNKNOWN: Exit code '" + string(nr.ExitCode) + "'undefined :" + nr.Text
		exitCode = 3
	}

	if len(nr.Multiline) > 0 {
		multiline := ""
		for _, l := range nr.Multiline {
			multiline = multiline + l + "\n"
		}
		fmt.Printf("%s|%s\n%s\n", text, nr.Perfdata, multiline)
	} else {
		fmt.Printf("%s|%s\n", text, nr.Perfdata)
	}
	os.Exit(exitCode)
}

func ExecuteCommand(command string, timeout int, allowFail bool) ExecResult {
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
		er.ReturnCode = msg.Sys().(syscall.WaitStatus).ExitStatus()
	}
	Debugf("Executing " + command + " took " + strconv.FormatFloat(duration, 'f', 5, 64) + "s")
	if err != nil && !allowFail {
		fmt.Println("executeCommand(): command failed: "+command, err)
		fmt.Println("executeCommand(): Output: " + string(out))
		os.Exit(1)
	}
	return er
}
