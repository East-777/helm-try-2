package main

import (
	"context"
	"fmt"
	"github.com/gofrs/flock"
	"golang.org/x/term"
	"helm.sh/helm/v3/pkg/getter"
	"helm.sh/helm/v3/pkg/repo"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"sigs.k8s.io/yaml"
	"strings"
	"time"
)

type repoAddOptions struct {
	name                 string
	url                  string
	username             string
	password             string
	passwordFromStdinOpt bool
	passCredentialsAll   bool
	forceUpdate          bool
	allowDeprecatedRepos bool

	certFile              string
	keyFile               string
	caFile                string
	insecureSkipTLSverify bool

	repoFile  string
	repoCache string

	// Deprecated, but cannot be removed until Helm 4
	deprecatedNoUpdate bool
}

func (o *repoAddOptions) run() error {
	//// Block deprecated repos
	//if !o.allowDeprecatedRepos {
	//	for oldURL, newURL := range deprecatedRepos {
	//		if strings.Contains(o.url, oldURL) {
	//			return fmt.Errorf("repo %q is no longer available; try %q instead", o.url, newURL)
	//		}
	//	}
	//}

	// Ensure the file directory exists as it is required for file locking
	err := os.MkdirAll(filepath.Dir(o.repoFile), os.ModePerm)
	if err != nil && !os.IsExist(err) {
		return err
	}

	// Acquire a file lock for process synchronization
	repoFileExt := filepath.Ext(o.repoFile)
	var lockPath string
	if len(repoFileExt) > 0 && len(repoFileExt) < len(o.repoFile) {
		lockPath = strings.Replace(o.repoFile, repoFileExt, ".lock", 1)
	} else {
		lockPath = o.repoFile + ".lock"
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

	b, err := ioutil.ReadFile(o.repoFile)
	if err != nil && !os.IsNotExist(err) {
		return err
	}

	var f repo.File
	if err := yaml.Unmarshal(b, &f); err != nil {
		return err
	}

	if o.username != "" && o.password == "" {
		if o.passwordFromStdinOpt {
			passwordFromStdin, err := io.ReadAll(os.Stdin)
			if err != nil {
				return err
			}
			password := strings.TrimSuffix(string(passwordFromStdin), "\n")
			password = strings.TrimSuffix(password, "\r")
			o.password = password
		} else {
			fd := int(os.Stdin.Fd())
			password, err := term.ReadPassword(fd)
			if err != nil {
				return err
			}
			o.password = string(password)
		}
	}

	c := repo.Entry{
		Name:                  o.name,
		URL:                   o.url,
		Username:              o.username,
		Password:              o.password,
		PassCredentialsAll:    o.passCredentialsAll,
		CertFile:              o.certFile,
		KeyFile:               o.keyFile,
		CAFile:                o.caFile,
		InsecureSkipTLSverify: o.insecureSkipTLSverify,
	}

	// If the repo exists do one of two things:
	// 1. If the configuration for the name is the same continue without error
	// 2. When the config is different require --force-update
	if !o.forceUpdate && f.Has(o.name) {
		existing := f.Get(o.name)
		if c != *existing {

			// The input coming in for the name is different from what is already
			// configured. Return an error.
			return fmt.Errorf("repository name (%s) already exists, please specify a different name", o.name)
		}

		// The add is idempotent so do nothing
		return nil
	}

	r, err := repo.NewChartRepository(&c, getter.All(settings))
	if err != nil {
		return err
	}

	if o.repoCache != "" {
		r.CachePath = o.repoCache
	}
	if _, err := r.DownloadIndexFile(); err != nil {
		return fmt.Errorf("Err: %s,looks like %q is not a valid chart repository or cannot be reached", err, o.url)
	}

	f.Update(&c)

	if err := f.WriteFile(o.repoFile, 0644); err != nil {
		return err
	}
	return nil
}
