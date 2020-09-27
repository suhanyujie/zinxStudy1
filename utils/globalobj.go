package utils

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
	"strings"
	"zinx_study1/ziface"
)

// read config file and init server

type GlobalObj struct {
	// Server
	TcpServer ziface.IServer
	Host      string
	TcpPort   int
	Name      string

	// Zinx
	Version    string
	MaxConn    int
	MaxPkgSize uint32
	// work pool size
	WorkPoolSize uint32
	// 每个 work 协程的任务数量限制
	MaxTaskLenForOneWork uint32
}

var GlobalObject *GlobalObj

func init() {
	// 实例化 global object
	GlobalObject = &GlobalObj{
		Name:                 "Zinx default server",
		Version:              "0.3",
		Host:                 "0.0.0.0",
		TcpPort:              3001,
		MaxConn:              10,
		MaxPkgSize:           1024,
		WorkPoolSize:         2,
		MaxTaskLenForOneWork: 100,
	}
	// 从配置文件中加载配置
	GlobalObject.GetConfigFromFile()
}

func (this *GlobalObj) GetConfigFromFile() {
	configFile := "config/zinx.json"
	if len(os.Args) > 1 {
		cliParam1 := os.Args[1]
		// 启动路径中包含 `core/conf` 后缀，表示是单测环境，否则就是普通的执行环境
		if strings.HasSuffix(cliParam1, "-test.v") {
			configFile = "./../config/zinx.json"
		}
	}
	data, err := ioutil.ReadFile(configFile)
	if err != nil {
		log.Fatalf("get config error: %s\n", err)
	}
	err = json.Unmarshal(data, GlobalObject)
	if err != nil {
		log.Fatalf("config data unmarshal error: %s\n", err)
	}
}
