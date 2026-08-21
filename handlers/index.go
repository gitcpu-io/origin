package handlers

import (
	"fmt"
	"os"

	"github.com/gitcpu-io/origin/configs"

	"github.com/gitcpu-io/zgo"
	"github.com/kataras/iris/v12"
)

/*
@Time : 2019-03-06 19:46
@Author : rubinus.chu
@File : index
@project: origin
*/

var commitOld = "latest"
var commit = "latest"

func init() {
	if os.Getenv("COMMIT") != "" {
		commit = os.Getenv("COMMIT")
	}
	if os.Getenv("COMMIT_OLD") != "" {
		commitOld = os.Getenv("COMMIT_OLD")
	}
}

type indexPage struct {
	Title     string
	Message   string
	Version   string
	HostName  string
	IP        string
	Commit    string
	CommitOld string
}

func IndexPage(ctx iris.Context) {
	project := configs.Conf.Project
	hostName, _ := os.Hostname()
	ctx.ViewData("", indexPage{
		Title:     project,
		Message:   fmt.Sprintf("%s welcome by zgo engine %s ...", project, zgo.Version),
		Version:   configs.Conf.Version,
		HostName:  hostName,
		IP:        zgo.Utils.GetIntranetIP(),
		Commit:    commit,
		CommitOld: commitOld,
	})
	err := ctx.View("index.html")
	if err != nil {
		zgo.Log.Error(err)
		return
	}
}

func FourZeroFourPage(ctx iris.Context) {
	project := configs.Conf.Project
	hostName, _ := os.Hostname()
	ctx.ViewData("", indexPage{
		Title:     project,
		Message:   fmt.Sprintf("%d -- %s by zgo engine %s ...", 404, project, zgo.Version),
		Version:   configs.Conf.Version,
		HostName:  hostName,
		IP:        zgo.Utils.GetIntranetIP(),
		Commit:    commit,
		CommitOld: commitOld,
	})
	err := ctx.View("404.html")
	if err != nil {
		zgo.Log.Error(err)
		return
	}
}

func FiveZeroZeroPage(ctx iris.Context) {
	project := configs.Conf.Project
	hostName, _ := os.Hostname()
	ctx.ViewData("", indexPage{
		Title:     project,
		Message:   fmt.Sprintf("%d -- %s by zgo engine %s ...", 500, project, zgo.Version),
		Version:   configs.Conf.Version,
		HostName:  hostName,
		IP:        zgo.Utils.GetIntranetIP(),
		Commit:    commit,
		CommitOld: commitOld,
	})
	err := ctx.View("500.html")
	if err != nil {
		zgo.Log.Error(err)
		return
	}
}

func Health(ctx iris.Context) {
	err := ctx.JSONP(map[string]string{
		"health": "true",
	})
	if err != nil {
		zgo.Log.Error(err)
		return
	}
}
