package executor

import (
	"bufio"
	"log"
	"os/exec"
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
