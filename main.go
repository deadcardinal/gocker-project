package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"strings"
	"sync"

	"github.com/alexflint/go-arg"
	"gopkg.in/yaml.v3"
)

type Config struct {
	Services map[string]Service
}

type Service struct {
	Labels   map[string]string `yaml:"labels"`
	Repo     Repo
	Commands Command
}

type Manipulator struct {
	Commands map[string]Command
}

type Command struct {
	ServiceName string
	Name        string
	Path        string
	Exec        string
}

type Git struct {
	Repos map[string]Repo
}

type Repo struct {
	ServiceName string
	Address     string
	Branch      string
	Path        string
}

type UpdateCmd struct {
	Services []string `arg:"positional"`
}

type CheckoutCmd struct {
	Branch   string   `arg:"positional"`
	Services []string `arg:"positional"`
}

type CommandCmd struct {
	Command string `arg:"positional" arg:"required"`
}

func (manipulator *Manipulator) run(services []Service, name string) {

}

func (git *Git) CloneSpecifiedRepos(repos []Repo) {
	var wg sync.WaitGroup
	wg.Add(len(repos))

	for _, repo := range repos {
		if _, err := os.Stat(repo.Path); os.IsNotExist(err) {
			go func(repo Repo) {
				repo.Clone()
				wg.Done()
			}(repo)
		} else {
			go func(repo Repo) {
				repo.Checkout()
				repo.Pull()
				wg.Done()
			}(repo)
		}
	}

	wg.Wait()
}

func (git *Git) CheckoutSpecifiedRepos(repos []Repo, branch string) {
	var wg sync.WaitGroup
	wg.Add(len(repos))

	for _, repo := range repos {
		if _, err := os.Stat(repo.Path); !os.IsNotExist(err) {
			go func(repo Repo, branch string) {
				repo.Branch = branch
				repo.Checkout()
				repo.Pull()
				wg.Done()
			}(repo, branch)
		} else {
			log.Output(0, "["+repo.ServiceName+"] clone done ")
			wg.Done()
		}
	}

	wg.Wait()
}

func (git *Git) FilterRepos(services []string) []Repo {
	repos := make([]Repo, 0)

	for _, service := range services {
		if _, ok := git.Repos[service]; ok {
			repos = append(repos, git.Repos[service])
		}
	}

	return repos
}

func (repo *Repo) Clone() {
	run(repo.ServiceName, ".", "git", "clone", repo.Address, "-b", repo.Branch, repo.Path)
	log.Output(0, "["+repo.ServiceName+"] clone done ")
}

func (repo *Repo) Checkout() {
	run(repo.ServiceName, repo.Path, "git", "checkout", repo.Branch)
	log.Output(0, "["+repo.ServiceName+"] checkout  done ")
}

func (repo *Repo) Pull() {
	run(repo.ServiceName, repo.Path, "git", "pull")
	log.Output(0, "["+repo.ServiceName+"] pull done ")
}

func (git *Git) CloneAllRepos() {
	repos := make([]Repo, 0, len(git.Repos))
	for _, repo := range git.Repos {
		repos = append(repos, repo)
	}

	git.CloneSpecifiedRepos(repos)
}

func (git *Git) CheckoutAllRepos(branch string) {
	repos := make([]Repo, 0, len(git.Repos))
	for _, repo := range git.Repos {
		repos = append(repos, repo)
	}

	git.CheckoutSpecifiedRepos(repos, branch)
}

func main() {
	var args struct {
		File     string       `default:"docker-compose.yml" arg:"-f" help:"config file"`
		Update   *UpdateCmd   `arg:"subcommand:update"`
		Checkout *CheckoutCmd `arg:"subcommand:checkout"`
		Command  *CommandCmd  `arg:"subcommand:command"`
	}
	arg.MustParse(&args)

	var config = parseYaml(args.File)

	if args.Update != nil {
		runUpdate(config, args.Update.Services)
	}

	if args.Checkout != nil {
		runCheckout(config, args.Checkout.Services, args.Checkout.Branch)
	}

	if args.Command != nil {
		runCommand(config)
	}
}

func runUpdate(config Config, services []string) {
	var git = parseGit(config)
	var filteredRepos = git.FilterRepos(services)

	if len(services) > 0 {
		git.CloneSpecifiedRepos(filteredRepos)
	} else {
		git.CloneAllRepos()
	}
}

func runCheckout(config Config, services []string, branch string) {
	var git = parseGit(config)
	var filteredRepos = git.FilterRepos(services)

	if len(services) > 0 {
		git.CheckoutSpecifiedRepos(filteredRepos, branch)
	} else {
		git.CheckoutAllRepos(branch)
	}
}

func runCommand(config Config) {
	var manipulator = parseManipulator(config)
	manipulator.Commands = nil
}

func parseYaml(file string) Config {
	data, err := ioutil.ReadFile(file)

	log.Output(0, "[system] use file "+file)
	if err != nil {
		log.Fatal(err)
	}

	config := Config{}
	err = yaml.Unmarshal(data, &config)
	if err != nil {
		log.Fatal(err)
	}

	return config
}

func parseGit(config Config) Git {
	git := Git{
		Repos: map[string]Repo{},
	}

	for service, serviceData := range config.Services {
		repo := Repo{}

		for label, labelValue := range serviceData.Labels {
			if label == "project.git" {
				repo.Address = labelValue
				repo.ServiceName = service
				repo.Path = "src/" + service
			}

			if label == "project.git.branch" {
				repo.Branch = labelValue
			}
		}

		if (Repo{} != repo) {
			git.Repos[service] = repo
		}
	}

	return git
}

func parseManipulator(config Config) Manipulator {
	manipulator := Manipulator{
		Commands: map[string]Command{},
	}

	for service, serviceData := range config.Services {
		command := Command{}

		for label, labelValue := range serviceData.Labels {
			if label != "project.git" && label != "project.git.branch" && strings.HasPrefix(label, "project.") {
				firstDot := strings.Index(label, ".") + 1
				command.Name = label[firstDot:len(label)]
				command.Exec = labelValue
				command.ServiceName = service
				command.Path = "src/" + service
			}
		}

		if (Command{} != command) {
			manipulator.Commands[service] = command
		}
	}

	fmt.Printf("%v\n", manipulator.Commands)

	return manipulator
}

func run(service string, dir string, name string, args ...string) {
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
