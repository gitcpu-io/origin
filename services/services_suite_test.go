package services_test

import (
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/gitcpu-io/zgo"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

const (
	label_bj = "mongo_label_bj"
)

func TestServices(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Services Suite")
}

var _ = BeforeSuite(func() {
	// 准备cpath
	var cpath string
	if cpath == "" {
		pwd, err := os.Getwd()
		if err == nil {
			arr := strings.Split(pwd, "/")
			cp := strings.Join(arr[:len(arr)-2], "/")
			cpath = fmt.Sprintf("%s/%s", cp, "origin/configs")
		}
	}
	//调用zgo
	err := zgo.Engine(&zgo.Options{
		CPath:   cpath,
		Env:     "local",
		Project: "origin",
		Mongo: []string{
			label_bj,
			//label_sh,
		},
	})
	//测试时表示使用mgo，在origin中使用一次
	if err != nil {
		panic(err)
	}
})

var _ = AfterSuite(func() {
})
