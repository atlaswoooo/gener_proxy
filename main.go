package main

import (
	accept "generproxy/handler/acceptreq"
	"generproxy/handler/monitor"
	"io"
	"log"
	"net"
	"os"
)

func init() {
	file, err := os.OpenFile("/home/qspace/log/generproxy.log", os.O_CREATE|os.O_RDWR|os.O_APPEND, os.ModePerm)
	if err != nil{
		log.Fatalln("Fail to open logger file :",err)
	}


	log.SetFlags(log.Lshortfile | log.Ldate | log.Lmicroseconds)
	log.SetOutput(io.MultiWriter(file,os.Stderr))
}

//server主函数
//server主函数

func main() {

	go monitor.Prome_test()

	l, err := net.Listen("tcp", ":80")
	log.Println("启动tcp服务，监听80端口")
	if err != nil {
		log.Panic(err)
	}

	for {
		client, err := l.Accept()
		if err != nil {
			log.Panic(err)
		}

		go accept.Handler(client)
	}
}

//

/*
func main() {

	//monitor.Prome_test()

	cert, err := tls.LoadX509KeyPair("server.pem", "server.key")
	if err != nil {
		log.Println(err)
		return
	}
	config := &tls.Config{Certificates: []tls.Certificate{cert}}
	l, err := tls.Listen("tcp", ":443", config)
	if err != nil {
		log.Println(err)
		return
	}

	for {
		client, err := l.Accept()
		if err != nil {
			log.Panic(err)
		}

		go accept.Handler(client)
	}
}
*/
