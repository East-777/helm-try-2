package release

import (
	"fmt"
	"helm-try-2/pkg/config"
	"helm.sh/helm/v3/pkg/action"
	"strconv"
)

type RollbackRel struct {
	config      *config.Config
	releaseName string
	prevision   string
}

func (r *RollbackRel) Rollback() error {
	client := action.NewRollback(r.config.Configuration)

	ver, err := strconv.Atoi(r.prevision)

	if err != nil {

		return fmt.Errorf("could not convert revision to a number: %v", err)
	}

	client.Version = ver

	return client.Run(r.releaseName)
}
