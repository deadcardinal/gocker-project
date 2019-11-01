package git

import (
	"gocker-project/executor"
	"log"
)

type Repository struct {
	ServiceName string
	Address     string
	Branch      string
	Path        string
}

func (repo *Repository) Clone() {
	executor.Run(repo.ServiceName, ".", "git", "clone", repo.Address, "-b", repo.Branch, repo.Path)
	log.Output(0, "["+repo.ServiceName+"] clone done ")
}

func (repo *Repository) Checkout() {
	executor.Run(repo.ServiceName, repo.Path, "git", "checkout", repo.Branch)
	log.Output(0, "["+repo.ServiceName+"] checkout  done ")
}

func (repo *Repository) Pull() {
	executor.Run(repo.ServiceName, repo.Path, "git", "pull")
	log.Output(0, "["+repo.ServiceName+"] pull done ")
}
