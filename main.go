package main

import (
	"log"
	"net"
)

func main() {

	sock, err := net.ListenUDP("udp4", &net.UDPAddr{
		IP:   net.IPv4(0, 0, 0, 0),
		Port: 53,
	})

	if err != nil {
		log.Fatal(err)
		return
	}

	defer sock.Close()

	for {

		// 读取UDP数据包
		buf := make([]byte, 1024)
		n, addr, _ := sock.ReadFromUDP(buf)

		// DNS报文解析
		msg := UnpackMsg(buf)
		if msg.question[0].Name == "www.test.com" {
			msg.header.Flags.Qr = 1
		}

		log.Println(msg.header, " | ", msg.header.Flags, " | ", msg.question[0].Name)

		log.Println("from:", addr)
		log.Println("len:", n, "bytes:", buf[0:n])

		log.Println("str:", string(buf[:n]))

		ret, err := query(buf[:n])
		check_error(err)
		if err == nil {
			sock.WriteToUDP(ret, addr)
		}
		log.Println("return(", len(ret), "): ", ret)
	}

}

func query(msg []byte) ([]byte, error) {
	raddr, err := net.ResolveUDPAddr("udp", "172.19.240.19:53")
	if err != nil {
		return nil, err
	}

	conn, err := net.DialUDP("udp", nil, raddr)
	defer conn.Close()

	if err != nil {
		return nil, err
	}

	_, err = conn.Write(msg)
	if err != nil {
		return nil, err
	}

	ret := make([]byte, 4096)
	n, _, err := conn.ReadFromUDP(ret)
	if err != nil {
		return nil, err
	}

	return ret[0:n], nil
}

func check_error(err error) {
	if err != nil {
		log.Println(err)
	}
}

func debug(fmt ...interface{}) {
	log.Println(fmt...)
}
