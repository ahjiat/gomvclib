package exec

import (
	"os/exec"
	"bytes"
)

func Wait(command string) (string, string, error) {
    var stdout, stderr bytes.Buffer
	var cmd *exec.Cmd
	cmd = exec.Command("/bin/bash", "-c", command)
	cmd.Stdout = &stdout
    cmd.Stderr = &stderr
	err := cmd.Run()
	return stdout.String(), stderr.String(), err
}
func NoWait(command string) error {
	var stdout, stderr bytes.Buffer
	var cmd *exec.Cmd
	cmd = exec.Command("/bin/bash", "-c", command)
	cmd.Stdout = &stdout
    cmd.Stderr = &stderr
	return cmd.Start()
}