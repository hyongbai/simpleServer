package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
)

type Server struct {
	ServerName string
	ServerIP   string
}

type Serverslice struct {
	Servers []Server
}

func main() {
	var s Serverslice
	//	str := `{"servers":[{"serverName":"Shanghai_VPN","serverIP":"127.0.0.1"},{"serverName":"Beijing_VPN","serverIP":"127.0.0.2"}]}`
	str, _ := ioutil.ReadFile("json.json")
	fmt.Println(string(str))
	json.Unmarshal(str, &s)
	fmt.Println(s)
}