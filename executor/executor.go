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
	input = strings.TrimSuffix(input, "\n")
	commandLine := strings.Split(input, " ")

	cmd := exec.Command(commandLine[0], commandLine[1:]...)
	cmd.Dir = dir

	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout

	return cmd.Run()
}
