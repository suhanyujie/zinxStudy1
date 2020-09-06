package utils

import (
	"encoding/json"
	"io/ioutil"
	"log"
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
}

var GlobalObject *GlobalObj

func init() {
	// 实例化 global object
	GlobalObject = &GlobalObj{
		Name:       "Zinx default server",
		Version:    "0.3",
		Host:       "0.0.0.0",
		TcpPort:    3001,
		MaxConn:    10,
		MaxPkgSize: 1024,
	}
	// 从配置文件中加载配置
	GlobalObject.GetConfigFromFile()
}

func (this *GlobalObj) GetConfigFromFile() {
	data, err := ioutil.ReadFile("config/zinx.json")
	if err != nil {
		log.Fatalf("get config error: %s\n", err)
	}
	err = json.Unmarshal(data, &GlobalObject)
	if err != nil {
		log.Fatalf("config data unmarshal error: %s\n", err)
	}
}
