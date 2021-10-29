package release

import (
	"helm-try-2/pkg/config"
	"helm.sh/helm/v3/pkg/action"
)

type UninstallRel struct {
	config *config.Config
	name   string
}

func (r *UninstallRel) Run() error {
	unInstall := action.NewUninstall(r.config.Configuration)
	_, err := unInstall.Run(r.name)
	return err
}
