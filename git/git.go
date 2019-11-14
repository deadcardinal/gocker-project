package git

import (
	"fmt"
	"gocker-project/env"
	"gocker-project/yaml"
	"log"
	"os"
	"sync"
)

type Git struct {
	Repos map[string]Repository
}

func (git *Git) CloneSpecifiedRepos(repos []Repository) {
	var wg sync.WaitGroup
	wg.Add(len(repos))

	for _, repo := range repos {
		if _, err := os.Stat(repo.Path); os.IsNotExist(err) {
			go func(repo Repository) {
				repo.Clone()
				wg.Done()
			}(repo)
		} else {
			go func(repo Repository) {
				repo.Checkout()
				repo.Pull()
				wg.Done()
			}(repo)
		}
	}

	wg.Wait()
}

func (git *Git) CheckoutSpecifiedRepos(repos []Repository, branch string) {
	var wg sync.WaitGroup
	wg.Add(len(repos))

	for _, repo := range repos {
		if _, err := os.Stat(repo.Path); !os.IsNotExist(err) {
			go func(repo Repository, branch string) {
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

func (git *Git) FilterRepos(services []string) []Repository {
	repos := make([]Repository, 0)

	for _, service := range services {
		if _, ok := git.Repos[service]; ok {
			repos = append(repos, git.Repos[service])
		}
	}

	return repos
}

func (git *Git) CloneAllRepos() {
	repos := make([]Repository, 0, len(git.Repos))
	for _, repo := range git.Repos {
		repos = append(repos, repo)
	}

	git.CloneSpecifiedRepos(repos)
}

func (git *Git) CheckoutAllRepos(branch string) {
	repos := make([]Repository, 0, len(git.Repos))
	for _, repo := range git.Repos {
		repos = append(repos, repo)
	}

	git.CheckoutSpecifiedRepos(repos, branch)
}

func NewFromConfig(config yaml.Config) Git {
	git := Git{
		Repos: map[string]Repository{},
	}

	appConfig := env.GetConfig()

	for service, serviceData := range config.Services {
		repo := Repository{}

		for label, labelValue := range serviceData.Labels {
			if label == "project.git" {
				repo.Address = labelValue
				repo.ServiceName = service

				if repo.Path == "" {
					repo.Path = fmt.Sprintf("%s/%s", appConfig.SourceDir, service)
				}
			}

			if label == "project.git.branch" {
				repo.Branch = labelValue
			}

			if label == "project.src" {
				repo.Path = fmt.Sprintf("%s/%s", labelValue, service)
			}
		}

		if (Repository{} != repo) {
			git.Repos[service] = repo
		}
	}

	return git
}
