module github.com/gitcpu-io/origin

go 1.16

require (
	github.com/0x19/goesl v0.0.0-20191107044804-3efcc2f41ccb
	github.com/fiorix/go-eventsocket v0.0.0-20180331081222-a4a0ee7bd315
	github.com/gitcpu-io/zgo v1.1.2
	github.com/gogo/protobuf v1.3.2
	github.com/golang/mock v1.4.4
	github.com/golang/protobuf v1.5.2
	github.com/kataras/iris/v12 v12.1.8
	github.com/onsi/ginkgo v1.14.0
	github.com/onsi/gomega v1.10.1
	github.com/op/go-logging v0.0.0-20160315200505-970db520ece7 // indirect
	github.com/spf13/pflag v1.0.5
	go.etcd.io/etcd/client/v3 v3.5.2
	golang.org/x/net v0.7.0
	google.golang.org/grpc v1.45.0
	google.golang.org/grpc/examples v0.0.0-20220329220628-b6873c006da7
	gopkg.in/ini.v1 v1.62.0 // indirect
)

//replace github.com/gitcpu-io/zgo => ../zgo
