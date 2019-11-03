package manipulator

import (
	"gocker-project/executor"
	"log"
	"strings"
)

type Command struct {
	ServiceName string
	Name        string
	Path        string
	Exec        string
}

func (command *Command) Run() {
	args := strings.Fields(command.Exec)
	executor.Run(command.ServiceName, command.Path, args[0], args[1:]...)
	log.Output(0, "["+command.ServiceName+"] executing \""+command.Name+"\" command ("+command.Exec+")")
}
