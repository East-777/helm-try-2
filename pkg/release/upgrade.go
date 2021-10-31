package release

import (
	"fmt"
	"helm-try-2/pkg/config"
	"helm.sh/helm/v3/pkg/action"
	"helm.sh/helm/v3/pkg/chart/loader"
	"helm.sh/helm/v3/pkg/cli/values"
	"helm.sh/helm/v3/pkg/downloader"
	"helm.sh/helm/v3/pkg/getter"
	"helm.sh/helm/v3/pkg/release"
)

type UpgradeRel struct {
	config      *config.Config
	releaseName string
	chart       string
}

func (u UpgradeRel) Upgrade() (*release.Release, error) {
	client := action.NewUpgrade(u.config.Configuration)
	valueOpts := &values.Options{}

	client.Namespace = u.config.EnvSettings.Namespace()

	if client.Version == "" && client.Devel {
		client.Version = ">0.0.0-0"
	}

	chartPath, err := client.ChartPathOptions.LocateChart(u.releaseName, u.config.EnvSettings)
	if err != nil {
		return nil, err
	}

	p := getter.All(u.config.EnvSettings)
	vals, err := valueOpts.MergeValues(p)
	if err != nil {
		return nil, err
	}

	// Check chart dependencies to make sure all are present in /charts
	ch, err := loader.Load(chartPath)
	if err != nil {
		return nil, err
	}
	if req := ch.Metadata.Dependencies; req != nil {
		if err := action.CheckDependencies(ch, req); err != nil {
			if client.DependencyUpdate {
				man := &downloader.Manager{
					ChartPath:        chartPath,
					Keyring:          client.ChartPathOptions.Keyring,
					SkipUpdate:       false,
					Getters:          p,
					RepositoryConfig: u.config.EnvSettings.RepositoryConfig,
					RepositoryCache:  u.config.EnvSettings.RepositoryCache,
					Debug:            u.config.EnvSettings.Debug,
				}
				if err := man.Update(); err != nil {
					return nil, err
				}
				// Reload the chart with the updated Chart.lock file.
				if ch, err = loader.Load(chartPath); err != nil {
					return nil, fmt.Errorf("failed reloading chart after repo update")
				}
			} else {
				return nil, err
			}
		}
	}

	return client.Run(u.releaseName, ch, vals)
}
