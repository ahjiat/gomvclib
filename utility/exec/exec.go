package exec

import (
	"os/exec"
	"bytes"
	"strings"
	"regexp"
)

func Wait(command string) (string, string, error) {
    var stdout, stderr bytes.Buffer
	var cmd *exec.Cmd
	re := regexp.MustCompile("[ \t]+")
	command = re.ReplaceAllString(command, " ")
	commandAry := strings.Split(command, " ")
	if len(commandAry) == 0 { return "", "", nil }
	name := commandAry[0] 
	if len(commandAry) == 1 {
		cmd = exec.Command(name)
	} else {
		cmd = exec.Command(name, commandAry[1:]...)
	}
	cmd.Stdout = &stdout
    cmd.Stderr = &stderr
	err := cmd.Run()
	return stdout.String(), stderr.String(), err
}

func NoWait(command string) {
	go Wait(command)
}