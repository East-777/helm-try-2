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
	"os"
)

type InstallRel struct {
	config *config.Config
	name   string
	chart  string
}

func (r *InstallRel) Install() (*release.Release, error) {

	client := action.NewInstall(r.config.Configuration)
	valueOpts := &values.Options{}

	fmt.Printf("Original chart version: %s", client.Version)
	if client.Version == "" && client.Devel {
		fmt.Printf("setting version to >0.0.0-0")
		client.Version = ">0.0.0-0"
	}

	client.ReleaseName = r.name

	chartPath, err := client.ChartPathOptions.LocateChart(r.chart, r.config.EnvSettings)
	if err != nil {
		return nil, err
	}
	fmt.Printf("CHART PATH: %s", chartPath)

	p := getter.All(r.config.EnvSettings)
	vals, err := valueOpts.MergeValues(p)
	if err != nil {
		return nil, err
	}

	// Check chart dependencies to make sure all are present in /charts
	chartRequested, err := loader.Load(chartPath)
	if err != nil {
		return nil, err
	}

	if req := chartRequested.Metadata.Dependencies; req != nil {
		// If CheckDependencies returns an error, we have unfulfilled dependencies.
		// As of Helm 2.4.0, this is treated as a stopping condition:
		// https://github.com/helm/helm/issues/2209
		if err := action.CheckDependencies(chartRequested, req); err != nil {
			if client.DependencyUpdate {
				man := &downloader.Manager{
					Out:              os.Stdout,
					ChartPath:        chartPath,
					Keyring:          client.ChartPathOptions.Keyring,
					SkipUpdate:       false,
					Getters:          p,
					RepositoryConfig: r.config.EnvSettings.RepositoryConfig,
					RepositoryCache:  r.config.EnvSettings.RepositoryCache,
					Debug:            r.config.EnvSettings.Debug,
				}
				if err := man.Update(); err != nil {
					return nil, err
				}
				// Reload the chart with the updated Chart.lock file.
				if chartRequested, err = loader.Load(chartPath); err != nil {
					return nil, fmt.Errorf("failed reloading chart after repo update")
				}
			} else {
				return nil, err
			}
		}
	}

	client.Namespace = r.config.EnvSettings.Namespace()

	return client.Run(chartRequested, vals)
}
