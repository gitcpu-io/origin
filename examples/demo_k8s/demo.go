package demo_limiter

import (
	"context"
	"fmt"
	"github.com/gitcpu-io/origin/configs"
	"github.com/gitcpu-io/zgo"
  "os"
  "strings"
)

/*
@Time : 2022-01-19 14:23
@Author : rubinus.chu
@File : demo
@project: origin
*/

func CallK8s() {
  // 准备cpath
  var cpath string
  if cpath == "" {
    pwd, err := os.Getwd()
    if err == nil {
      arr := strings.Split(pwd,"/")
      cp := strings.Join(arr[:len(arr) - 2],"/")
      cpath = fmt.Sprintf("%s/%s", cp, "configs")
    }
  }
	configs.InitConfig(cpath, "local", "", "", 0, 0)
	err := zgo.Engine(&zgo.Options{
		CPath:    configs.Conf.CPath,
		Env:      configs.Conf.Env,
		Project:  configs.Conf.Project,
		Loglevel: configs.Conf.Loglevel,
	})

	if err != nil {
		panic(err)
	}

	//###############################################################开始使用k8s
	//1.通过配置选项Func
	cof := zgo.K8s.ConfigOption().GetFunc()

	//2.为参数赋值
	kco, err := zgo.K8s.ConfigOption().Build(
		cof.WithMasterUrl("https://kubernetes.docker.internal:6443"),
		cof.WithKubeConfig("/Users/rubinus/.kube/config"),
	)
	if err != nil {
		panic(err)
	}
	//2-1.为另一个赋值
	kco_1, err := zgo.K8s.ConfigOption().Build(
		cof.WithMasterUrl("https://install-prow-dns-809743c1.hcp.eastasia.azmk8s.io:443"),
		cof.WithKubeConfig("/Users/rubinus/app/kubeconfig/config-prow"),
	)

	if err != nil {
		zgo.Log.Infof("----%+v", err)
		return
	}

	//3.调用生成config
	config_0, err := zgo.K8s.Builder().BuildConfig(kco)
	if err != nil {
		return
	}

	//3-1.调用生成另外一个config
	config_1, err := zgo.K8s.Builder().BuildConfig(kco_1)
	if err != nil {
		return
	}

	//4. 打印config
	fmt.Println("config: ", zgo.K8s.GetContext(config_0.Host))
	fmt.Println("config-1: ", zgo.K8s.GetContext(config_1.Host))

	//5.生成clientset
	_, err = zgo.K8s.Builder().BuildClientSet(config_0.Host, config_0)
	if err != nil {
		return
	}

	//5-1生成另外一个clientset
	_, err = zgo.K8s.Builder().BuildClientSet(config_1.Host, config_1)
	if err != nil {
		return
	}

	//6.打印clientset
	fmt.Println("clientset: ", zgo.K8s.GetClientSet(config_0.Host))
	fmt.Println("clientset-1: ", zgo.K8s.GetClientSet(config_1.Host))

	//7. 打印version
	info, err := zgo.K8s.ServerVersion(config_0.Host)
	if err != nil {
		fmt.Println("==ServerVersion==err: ", err)
		return
	}
	zgo.Log.Infof("%s,%s", info.Platform, info.GitVersion)

	//7-1. 打印另一个version
	info_1, err := zgo.K8s.ServerVersion(config_1.Host)
	if err != nil {
		fmt.Println("err: ", err)
		return
	}
	fmt.Println("\n----")
	zgo.Log.Infof("%s,%s", info_1.Platform, info_1.GitVersion)

	//8.
	zgo.K8s.UseContext(config_0.Host) //使用指定的configs.host的context
	dlist, err := zgo.K8s.Deployment().List(context.TODO(), "default", "", "", -1, false)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Printf("总量：%v\n", len(dlist.Items))
	for _, item := range dlist.Items {
		fmt.Println(item.Namespace, item.Name)
	}
}
