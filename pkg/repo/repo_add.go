package repo

import (
	"context"
	"fmt"
	"github.com/gofrs/flock"
	"gopkg.in/yaml.v2"
	"helm-try-2/pkg/config"
	"helm.sh/helm/v3/pkg/getter"
	"helm.sh/helm/v3/pkg/repo"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"time"
)

type AddRepo struct {
	Config *config.Config

	Name      string
	Url       string
	caFile    string
	keyFile   string
	certFile  string
	username  string
	password  string
	RepoFile  string
	RepoCache string
}

func (o *AddRepo) Add() error {

	// Ensure the file directory exists as it is required for file locking
	err := os.MkdirAll(filepath.Dir(o.RepoFile), os.ModePerm)
	if err != nil && !os.IsExist(err) {
		return err
	}

	// Acquire a file lock for process synchronization
	repoFileExt := filepath.Ext(o.RepoFile)
	var lockPath string
	if len(repoFileExt) > 0 && len(repoFileExt) < len(o.RepoFile) {
		lockPath = strings.Replace(o.RepoFile, repoFileExt, ".lock", 1)
	} else {
		lockPath = o.RepoFile + ".lock"
	}
	fileLock := flock.New(lockPath)
	lockCtx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	locked, err := fileLock.TryLockContext(lockCtx, time.Second)
	if err == nil && locked {
		defer fileLock.Unlock()
	}
	if err != nil {
		return err
	}

	b, err := ioutil.ReadFile(o.RepoFile)
	if err != nil && !os.IsNotExist(err) {
		return err
	}

	var f repo.File
	if err := yaml.Unmarshal(b, &f); err != nil {
		return err
	}

	c := repo.Entry{
		Name:     o.Name,
		URL:      o.Url,
		Username: o.username,
		Password: o.password,
		CertFile: o.certFile,
		KeyFile:  o.keyFile,
		CAFile:   o.caFile,
	}

	if f.Has(o.Name) {
		return fmt.Errorf("repository %q already exists\n", o.Name)
	}

	r, err := repo.NewChartRepository(&c, getter.All(o.Config.EnvSettings))
	if err != nil {
		return err
	}

	if o.RepoCache != "" {
		r.CachePath = o.RepoCache
	}
	if _, err := r.DownloadIndexFile(); err != nil {
		return fmt.Errorf("err: %s ,looks like %q is not a valid chart repository or cannot be reached\n", err, o.Url)
	}

	f.Update(&c)

	if err := f.WriteFile(o.RepoFile, 0644); err != nil {
		return err
	}
	fmt.Printf("%q has been added to your repositories\n", o.Name)
	return nil
}
