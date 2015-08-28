# go-nagios
Golang package for Nagios plugins

## Usage

```
import "github.com/xorpaul/go-nagios"
```

```
ExecuteCommand

VARIABLES

var (
    Debug bool
)


FUNCTIONS

func Debugf(s string)
    Debugf is a helper function for debug logging if Debug is true

func NagiosExit(nr NagiosResult)
    NagiosExit uses the NagiosResult struct to output Nagios plugin
    compatible output and exit codes


TYPES

type ExecResult struct {
    ReturnCode int
    Output     string
}


func ExecuteCommand(command string, timeout int, allowFail bool) ExecResult
    ExecuteCommand executes a shell command and provides the stdout and
    stderr output combined and the exit code. First arg is the command to
    execute as a simple string. Second arg is a timeout parameter as a int
    after which to kill the command and return. Third arg is a bool flag if
    the command is allowed to fail or not.



type NagiosResult struct {
    ExitCode  int
    Text      string
    Perfdata  string
    Multiline []string
}
```

## Example usage

```
// create NagiosResult with default values (UNKNOWN in this case)
nr := nagios.NagiosResult{ExitCode: 3, Text: "uncatched case", Perfdata: ""}

// execute a command with ExecuteCommand
er := nagios.ExecuteCommand("nfsstat -c -l", 1, false)
nr = doSomethingWithTheCommandOutputOrExitCode(er.Output)

nagios.NagiosExit(nr)
```


Look at https://github.com/xorpaul/check_nfs_client as an example.
