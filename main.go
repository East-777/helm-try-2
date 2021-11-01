package main

import (
	"fmt"
	"helm.sh/helm/v3/pkg/action"
	"helm.sh/helm/v3/pkg/cli"
	"helm.sh/helm/v3/pkg/cli/values"
	"k8s.io/cli-runtime/pkg/genericclioptions"
	"log"
	"os"
	"path"
)

const (
	repoConfig = "testdata/repositories.yaml"
	repoCache  = "testdata/repository"
)

////
var settings = cli.New()
var actionConfig = new(action.Configuration)

func init() {
	d, _ := os.Getwd()
	kubeConfig := path.Join(d, "config")

	getter := &genericclioptions.ConfigFlags{
		KubeConfig: &kubeConfig,
	}
	//kc := kube.New(getter)
	//actionConfig.RESTClientGetter = getter
	//actionConfig.KubeClient = kc

	actionConfig.Init(getter, settings.Namespace(), "memory", log.Printf)

	settings.KubeConfig = kubeConfig
	settings.RepositoryConfig = path.Join(d, repoConfig)
	settings.RepositoryCache = path.Join(d, repoCache)

	//storage := storage.Init(driver.NewMemory())
	//storage.Create(&release.Release{
	//	Name: "myrelease",
	//	Info: &release.Info{Status: release.StatusDeployed},
	//	Chart: &chart.Chart{
	//		Metadata: &chart.Metadata{
	//			Name:    "Myrelease-Chart",
	//			Version: "1.2.3",
	//		},
	//	},
	//	Version: 1,
	//})
	//
	//actionConfig.Releases = storage
}

func main() {
	//d, _ := os.Getwd()
	//config, err := clientcmd.BuildConfigFromFlags("", path.Join(d, "config"))
	//if err != nil {
	//	fmt.Printf("%v", err)
	//	return
	//}
	//
	//clientset, err := kubernetes.NewForConfig(config)
	//if err != nil {
	//	fmt.Printf("%v", err)
	//	return
	//}
	//namespaceList, err := clientset.CoreV1().Namespaces().List(context.Background(), metav1.ListOptions{})
	//if err != nil {
	//	fmt.Printf("%v", err)
	//	return
	//}
	//
	//for _, namespace := range namespaceList.Items {
	//	fmt.Println(namespace.Name)
	//}

	//o := &repoAddOptions{}
	//o.name = "aliyuncs"
	//o.url = "https://apphub.aliyuncs.com"
	//o.repoFile = settings.RepositoryConfig
	//o.repoCache = settings.RepositoryCache
	//
	//o.run()

	client := action.NewInstall(actionConfig)
	valueOpts := &values.Options{}

	//name := "my-nginx"
	//chart := "aliyuncs/nginx"

	request := &InstallRequest{
		config:   actionConfig,
		settings: settings,
		name:     "my-nginx",
		chart:    "nginx",
	}

	//request.name = name
	//request.chart = chart
	release, err := request.runInstall(client, valueOpts)
	if err != nil {
		fmt.Printf("install err:%s", err)
	}

	fmt.Printf("release name:%s", release.Name)
	//store := actionConfig.Releases
	//store.Create(release)
	//
	//history, err := store.History(release.Name)
	//
	//fmt.Printf("History num:%s", len(history))
	//
	//unInstall := action.NewUninstall(actionConfig)
	//
	//unInstall.Run(release.Name)
	//
	//history2, _ := store.History(release.Name)
	//fmt.Printf("History2 num:%s", len(history2))
}

//func main() {
//
//	url := "https://github.com/East-777/helm-try-2.git"
//
//	referenceName := "v1.0.0"
//
//	//git.Clone()
//
//	repository, err := git.PlainClone("F:/GoLand/git", false,
//		&git.CloneOptions{
//			URL:           url,
//			ReferenceName: plumbing.NewTagReferenceName(referenceName),
//		})
//
//	if err != nil {
//		fmt.Printf("err:%s", err)
//	}
//
//	_, _ = repository.Config()
//
//}

//func main() {
//
//	config, err := config.New("default")
//
//	if err != nil {
//		fmt.Printf("err:%s", err)
//		return
//	}
//
//	install := &release.InstallRel{
//		Config: config,
//		Name:   "my-nginx",
//		Chart:  "aliyuncs/nginx",
//	}
//
//	r, err := install.Install()
//
//	if err != nil {
//		fmt.Printf("err:%s", err)
//		return
//	}
//
//	fmt.Printf("version:%s",r.Version)
//
//}
