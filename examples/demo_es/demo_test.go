package demo_es

import (
	"context"
	"fmt"
	"github.com/gitcpu-io/origin/configs"
	"github.com/gitcpu-io/zgo"
	"testing"
	"time"
)

const (
	new_write = "label_new"
)

func TestGet(t *testing.T) {
	configs.InitConfig("", "local", "", "", 0, 0)

	err := zgo.Engine(&zgo.Options{
		Env:      configs.Conf.Env,
		Loglevel: configs.Conf.Loglevel,
		Project:  configs.Conf.Project,
		Es: []string{
			new_write,
		},
	})

	if err != nil {
		panic(err)
	}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	args := map[string]interface{}{}
	index := "active_bj_house_sell"
	table := "spider"
	dsl := `{"query": {"match_all": {}}}`

	//sellR, _ := zgo.Es.New(new_write)

	res, err := zgo.Es.SearchDsl(ctx, index, table, dsl, args)
	if err != nil {
		panic(err)
	}
	fmt.Print(res)
	//result, err := sellR.Search(ctx, args)

	//fmt.Print(result, err)
}
