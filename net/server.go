package net

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net"
)

// 种子结点
var SeedNodes=[]string{"localhost:3000"}

//
func StartServer(nodeID string){
	var port string
	_, _ = fmt.Scanf("PORT:%s", &port)
	nodeAddress:=fmt.Sprintf("localhost:%s",port)
	listen,err:=net.Listen("tcp",nodeAddress)
	if err!=nil{
		log.Panicf("listen ERROR:%v\n",err)
	}
	defer listen.Close()
	if nodeAddress!=SeedNodes[0]{
		// 不是主节点，发送请求，同步数据
	}
	for{
		conn,err:=listen.Accept()
		if err!=nil{
			log.Panicf("Listen connect ERROR:%v\n",err)
		}
		request,err:=ioutil.ReadAll(conn)
		if err!=nil{
			log.Panicf("Read message ERROR:%v\n",err)
		}
		fmt.Printf("Receive message:%v\n",request)

	}
}

func SendMessage(to,from string){
	fmt.Printf("Connect to server[%s]...\n",to)
	conn,err:=net.Dial("tcp",to)
	if err!=nil{
		log.Panicf("connect to server ERROR:%s\n",to)
	}
	_, err = io.Copy(conn, bytes.NewReader([]byte(from)))
	if err!=nil{
		log.Panicf("Receive message ERROR:%v\n",err)
	}
}