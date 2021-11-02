package chart

import (
	"fmt"
	"gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/plumbing"
	"helm-try-2/pkg/config"
	"helm.sh/helm/v3/pkg/action"
	"helm.sh/helm/v3/pkg/chart"
	"helm.sh/helm/v3/pkg/chart/loader"
	"os"
	"path"
	//"helm.sh/helm/v3/pkg/cli"
)

type Chart struct {
	Name string
	Git
	Config *config.Config
	action.ChartPathOptions
}
type Git struct {
	Url string
	Tag string
}

func (c Chart) Get() (*chart.Chart, error) {
	dir := c.Name
	if dir == "" {
		if url := c.Git.Url; url != "" {
			pwd, err := os.Getwd()
			dir = path.Join(pwd, "git")
			_, err = os.Stat(dir)
			if !os.IsNotExist(err) {
				os.RemoveAll(dir)
			}
			os.Mkdir(dir, 0644)
			_, err = git.PlainClone(dir, false,
				&git.CloneOptions{
					URL:           url,
					ReferenceName: plumbing.NewTagReferenceName(c.Tag),
				})

			if err != nil {
				return nil, err
			}
		} else {
			return nil, fmt.Errorf("need to provide git's url")
		}
	}

	chartPath, err := c.ChartPathOptions.LocateChart(dir, c.Config.EnvSettings)
	if err != nil {
		return nil, err
	}
	fmt.Printf("CHART PATH: %s", chartPath)

	// Check chart dependencies to make sure all are present in /charts
	chartRequested, err := loader.Load(chartPath)

	return chartRequested, err

}
