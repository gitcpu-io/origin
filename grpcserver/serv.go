package grpcserver

import (
	"context"
	"fmt"
	"github.com/rubinus/origin/config"
	"github.com/rubinus/origin/grpchandlers"
	"github.com/rubinus/origin/pb/helloworld"
	"github.com/rubinus/zgo"
)

/*
@Time : 2019-06-10 18:47
@Author : rubinus.chu
@File : serv
@project: origin
*/

func Start() {
	server, err := zgo.Grpc.Server(context.TODO())
	if err != nil {
		fmt.Println("-----grpc grpcserver is error :", err)
		return
	}

	pb_helloworld.RegisterHelloWorldServiceServer(server, &grpchandlers.HelloWorldServer{})

	//在这里继续添加你的rpc服务
	// todo add your rpc service

	fmt.Printf("Now listening GRPC Serv on: %s\n", config.Conf.RpcPort)
	_, err = zgo.Grpc.Run(context.TODO(), server, config.Conf.RpcPort)
	if err != nil {
		fmt.Println("=====grpc grpcserver is error :", err)
		return
	}
}