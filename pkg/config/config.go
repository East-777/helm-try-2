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
	repoConfig = "F:/GoLand/helm-try-2/testdata/repositories.yaml"
	repoCache  = "F:/GoLand/helm-try-2/testdata/repository"
)

func New(namespace string) (*Config, error) {

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
	settings.RegistryConfig = repoConfig
	settings.RepositoryCache = repoCache

	return &Config{
		Configuration: config,
		EnvSettings:   settings,
	}, nil
}

type Config struct {
	*action.Configuration
	*cli.EnvSettings
}
