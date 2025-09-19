package executors

import (
	"os"
	"os/exec"
)

// runCommand executes a command and returns its output
func runCommand(name string, args ...string) (string, error) {
	cmd := exec.Command(name, args...)
	output, err := cmd.CombinedOutput()
	return string(output), err
}

// writeLog writes content to a log file
func writeLog(filename, content string) error {
	return os.WriteFile(filename, []byte(content), 0644)
}
