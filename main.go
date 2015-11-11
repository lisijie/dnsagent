package main

import (
	"flag"
	"github.com/lisijie/go-conf"
	"log"
	"net"
	"regexp"
	"strings"
)

var (
	config   *goconf.Config
	confFile = flag.String("conf", "./config.ini", "配置文件路径")
	dns      = flag.String("dns", "8.8.8.8:53", "DNS地址（本地查不到时向该服务器查询）")
	isDebug  = flag.Bool("debug", false, "是否启用调试模式")
)

func main() {
	flag.Parse()

	var err error
	config, err = goconf.NewConfig(*confFile)
	if err != nil {
		log.Fatal(err)
		return
	}

	sock, err := net.ListenUDP("udp4", &net.UDPAddr{
		IP:   net.IPv4(0, 0, 0, 0),
		Port: 53,
	})

	if err != nil {
		log.Fatal(err)
		return
	}

	defer sock.Close()

	log.Println("启动服务并监听 0.0.0.0:53 端口...")

	// 找出有包含通配符的，用正则替代
	regMap := make(map[string]string)
	for name, ip := range config.GetAll() {
		if strings.Contains(name, "*") {
			rv := regexp.QuoteMeta(name)
			rv = strings.Replace(rv, "\\*", "(.*)", -1)
			regMap[rv] = ip
		}
	}

	for {
		// 读取UDP数据包
		buf := make([]byte, 1024)
		n, addr, _ := sock.ReadFromUDP(buf)

		go func() {
			// DNS报文解析
			msg := UnpackMsg(buf[:n])
			if ip := config.GetString(msg.GetQuestion(0).Name); ip != "" {
				msg.SetResponse()
				msg.AddAnswer(NewA(msg.GetQuestion(0).Name, ip))
				ret := PackMsg(msg)
				sock.WriteToUDP(ret, addr)
				debug("[L]解析: ", msg.GetQuestion(0).Name)
			} else {
				for rexp, ip := range regMap {
					if ok, _ := regexp.MatchString(rexp, msg.GetQuestion(0).Name); ok {
						msg.SetResponse()
						msg.AddAnswer(NewA(msg.GetQuestion(0).Name, ip))
						ret := PackMsg(msg)
						sock.WriteToUDP(ret, addr)
						debug("[L2]解析: ", msg.GetQuestion(0).Name)
						return
					}
				}
				ret, err := query(buf[:n])
				check_error(err)
				if err == nil {
					sock.WriteToUDP(ret, addr)
				}
				debug("[R]解析: ", msg.GetQuestion(0).Name)
			}
		}()

	}

}

func query(msg []byte) ([]byte, error) {
	raddr, err := net.ResolveUDPAddr("udp", *dns)
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
	if *isDebug {
		log.Println(fmt...)
	}
}
