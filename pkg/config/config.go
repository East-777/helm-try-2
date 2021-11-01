package config

import (
	"helm.sh/helm/v3/pkg/action"
	"helm.sh/helm/v3/pkg/cli"
	"k8s.io/cli-runtime/pkg/genericclioptions"
	"log"
	"os"
	"path"
)

const (
	repoConfig = "testdata/repositories.yaml"
	repoCache  = "testdata/repository"
)

func New(namespace string) (*Config, error) {
	if namespace == "" {
		namespace = "default"
	}

	settings := cli.New()
	config := &action.Configuration{}

	pwd, _ := os.Getwd()
	kubeConfig := path.Join(pwd, "config")

	getter := &genericclioptions.ConfigFlags{
		KubeConfig: &kubeConfig,
	}

	if err := config.Init(getter, namespace, "memory", log.Printf); err != nil {
		return nil, err
	}

	settings.KubeConfig = kubeConfig
	settings.RepositoryConfig = path.Join(pwd, repoConfig)
	settings.RepositoryCache = path.Join(pwd, repoCache)

	return &Config{
		Configuration: config,
		EnvSettings:   settings,
	}, nil
}

type Config struct {
	*action.Configuration
	*cli.EnvSettings
}
