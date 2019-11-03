package executor

import (
	"bufio"
	"log"
	"os"
	"os/exec"
	"strings"
)

func Run(service string, dir string, name string, args ...string) {
	cmd := exec.Command(name, args...)
	cmd.Dir = dir

	stderr, _ := cmd.StderrPipe()
	cmd.Start()

	scanner := bufio.NewScanner(stderr)
	scanner.Split(bufio.ScanLines)
	for scanner.Scan() {
		m := scanner.Text()
		log.Output(0, "["+service+"] "+m)
	}

	cmd.Wait()
}

func ExecInput(service string, dir string, input string) error {
	// Remove the newline character.
	input = strings.TrimSuffix(input, "\n")
	commandLine := strings.Split(input, " ")

	// Prepare the command to execute.
	cmd := exec.Command(commandLine[0], commandLine[1:]...)
	cmd.Dir = dir

	// Set the correct output device.
	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout

	// Execute the command and return the error.
	return cmd.Run()
}
