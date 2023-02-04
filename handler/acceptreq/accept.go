package accept

import (
	"fmt"
	"generproxy/handler/monitor"
	"log"
	"net"
)

//过滤链接
//过滤链接
//返回值: true 需要过滤

//主handler异常处理
//有异常返回true
func abnormalHandle(client net.Conn) bool {



	var ret = false
	if client == nil {
		return true
	}

	//过滤本地局域网ip
	if filterReq(client.RemoteAddr().String()) {
		return true
	}

	return ret
}

//parseIP4header
//parseIP4header
func parseIP4header(buffer []byte) {
	//var versionAndIHL byte
	//var version byte
	//var IHL byte
	//var headerLength byte
	//var typeOfService uint16
	//var totalLength uint16
	//var identificationAndFlagsAndFragmentOffset uint32
	//var TTL byte
	//var protocolNum byte
	//var protocol string
	//var headerChecksum uint16
	//var sourceAddress string
	//var destinationAddress string
	//
	//
	////fmt.Printf("buffer:%+v  \n", buffer)
	////fmt.Printf("buffer[2]:%+v  \n", buffer[2])
	////fmt.Printf("buffer[3]:%+v  \n", buffer[3])
	////fmt.Printf(" buffer[2:4]:%+v  \n", buffer[2:4])
	//
	//versionAndIHL = buffer[0]
	//version = versionAndIHL >> 4
	//IHL = versionAndIHL & 0x0F
	//headerLength = IHL << 2
	//typeOfService = uint16(buffer[1] & 0xFF)
	//totalLength = binary.LittleEndian.Uint16(buffer[2:4]) & 0xFFFF
	//identificationAndFlagsAndFragmentOffset = binary.LittleEndian.Uint32(buffer[4:8])
	//TTL = buffer[8] & 0x0F
	//protocolNum = buffer[9] & 0x0F
	//if protocolNum == 6 {
	//	protocol = "UDP"
	//} else if protocolNum == 17 {
	//	protocol = "TCP"
	//}
	//headerChecksum = binary.LittleEndian.Uint16(buffer[10:12]) & 0xFFFF
	//sourceAddress = fmt.Sprintf("%d.%d.%d.%d", buffer[12], buffer[13], buffer[14], buffer[15])
	//destinationAddress = fmt.Sprintf("%d.%d.%d.%d", buffer[16], buffer[17], buffer[18], buffer[19])
	//
	//




}


//handle处理主函数
//handle处理主函数
func Handler(client net.Conn) {
	//异常处理
	defer client.Close()
	if abnormalHandle(client) {
		return
	}
	//上报监控
	monitor.RecordMetrics()

	//获取用户地理位置
	country, location := getUserGeolocation(client.RemoteAddr())
	log.Printf("客户端接入信息: %+v %+v %+v", client.RemoteAddr(), country, location)

	for{
		//从客户端获取数据method，url
		var buffer [10240]byte

		n, err := client.Read(buffer[:])
		//log.Printf("读取客户数据")
		if err != nil {
			log.Printf("从客户端获取数据 err: %v", err)
			return
		}
		go readClientData(buffer[:n],client)
	}

	
}

func readClientData(buffer []byte,client net.Conn)  {
	var address string

	var (
		destinationIp string
		destinationPort uint16
	)
	//如果是ipv4数据包
	if buffer[0] == 69{
		//报文头部为20字节
		ipv4Header := ParseIPV4Header(buffer[:20])
		log.Printf("ip buffer:%v\n",buffer[:20])

		destinationIp = ipv4Header.DestinationAddress


		if ipv4Header.Protocol == "TCP"{
			//tcp报文首部长度也为20字节
			tcpHeader := ParseTcpHeader(buffer[20:40])
			log.Printf("%v\n",tcpHeader.Flags)
			if tcpHeader.Flags == 0x002{
				log.Printf("是tcp握手第一次包")
				//直接封装tcp第二次握手包放回给client
				tcpResponce := &TCPHeader{}
				tcpResponce.SourcePort = tcpHeader.DestinationPort
				tcpResponce.DestinationPort = tcpHeader.SourcePort


			}

			destinationPort = tcpHeader.DestinationPort

			address = fmt.Sprintf("%s:%d",destinationIp,destinationPort)
			//log.Printf("客户端需要访问的tcp真实地址 : %v", address)

			//获得了请求的host和port，向服务端发起tcp连接
			server, err := net.Dial("tcp", address)
			if err != nil {
				log.Printf(" 向服务端发起tcp连接 err: %v", err)
				return
			}
			defer server.Close()

			server.Write(buffer[40:])

			readeResponse := func() {
				var data []byte
				n, err := server.Read(data)
				if err != nil{
					log.Printf("从目标服务器读响应错误：%v\n",err)
					return
				}
				log.Printf("目标服务器响应数据：%v\n",data)
				//写回客户端
				_, err = client.Write(data[:n])
				if err != nil {
					log.Printf("写回客户端失败")
					return
				}
			}

			readeResponse()

			//io.Copy(client,server)
			log.Printf("收到目标服务器响应，发送给客户端")


			//server.Write(buffer[40:])
			//
			////将客户端的请求转发至服务端，将服务端的响应转发给客户端。io.Copy为阻塞函数，文件描述符不关闭就不停止
			//go io.Copy(server, client)
			//io.Copy(client, server)


		}else if ipv4Header.Protocol == "UDP"{
			//udp首部8字节
			udpHeader := ParseUDPHeader(buffer[20:28])

			//log.Printf("udp数据报首部：%v\n",udpHeader)

			destinationPort = udpHeader.DestinationPort
			address = fmt.Sprintf("%s:%d",destinationIp,destinationPort)
			//log.Printf("客户端需要访问的udp真实地址 : %v", address)

			server, err := net.ListenUDP("udp", &net.UDPAddr{
				IP:   []byte(destinationIp),
				Port: int(destinationPort)})
			if err != nil {
				log.Printf(" 向服务端发起udp连接 err: %v", err)
				log.Println(err)
				return
			}
			defer server.Close()
		}
	}else{
		return
	}



	//indexByte := bytes.IndexByte(buffer[:], '\n')
	//if indexByte == -1{
	//	return
	//}
	//fmt.Sscanf(string(buffer[:indexByte]), "%s%s", &method, &URL)
	//
	//hostPortURL, err := url.Parse(URL)
	//if err != nil {
	//	log.Printf("从客户端数据读入method,url err: %v", err)
	//	return
	//} else {
	//	log.Printf("从客户端数据读入method,url : %v %v", method, URL)
	//}
	//
	//// 如果方法是CONNECT，则为https协议
	//if method == CONNECT_ {
	//	address = hostPortURL.Scheme + ":" + hostPortURL.Opaque
	//} else { //否则为http协议
	//	address = hostPortURL.Host
	//	// 如果host不带端口，则默认为80
	//	if strings.Index(hostPortURL.Host, ":") == -1 { //host不带端口， 默认80
	//		address = hostPortURL.Host + ":80"
	//	}
	//}



	////如果使用https协议，需先向客户端表示连接建立完毕
	//if method == CONNECT_ {
	//	fmt.Fprint(client, "HTTP/1.1 200 Connection established\r\n")
	//} else { //如果使用http协议，需将从客户端得到的http请求转发给服务端
	//	server.Write(buffer)
	//}


}