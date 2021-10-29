package repo

import (
	"fmt"
	"github.com/pkg/errors"
	"helm-try-2/pkg/config"
	"helm.sh/helm/v3/pkg/helmpath"
	"helm.sh/helm/v3/pkg/repo"
	"os"
	"path/filepath"
)

type RemoveRepo struct {
	config *config.Config

	names []string
}

func (o *RemoveRepo) Remove() error {

	r, err := repo.LoadFile(o.config.EnvSettings.RegistryConfig)
	if isNotExist(err) || len(r.Repositories) == 0 {
		return errors.New("no repositories configured")
	}

	for _, name := range o.names {
		if !r.Remove(name) {
			return errors.Errorf("no repo named %q found", name)
		}
		if err := r.WriteFile(o.config.EnvSettings.RegistryConfig, 0644); err != nil {
			return err
		}

		if err := removeRepoCache(o.config.EnvSettings.RepositoryCache, name); err != nil {
			return err
		}
		fmt.Printf("%q has been removed from your repositories\n", name)
	}

	return nil
}

func isNotExist(err error) bool {
	return os.IsNotExist(errors.Cause(err))
}

func removeRepoCache(root, name string) error {
	idx := filepath.Join(root, helmpath.CacheChartsFile(name))
	if _, err := os.Stat(idx); err == nil {
		os.Remove(idx)
	}

	idx = filepath.Join(root, helmpath.CacheIndexFile(name))
	if _, err := os.Stat(idx); os.IsNotExist(err) {
		return nil
	} else if err != nil {
		return errors.Wrapf(err, "can't remove index file %s", idx)
	}
	return os.Remove(idx)
}
