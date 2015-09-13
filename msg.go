package main

import "strings"

// 包头，12字节
//                                     1  1  1  1  1  1
//      0  1  2  3  4  5  6  7  8  9  0  1  2  3  4  5
//    +--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+
//    |                      ID                       |
//    +--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+
//    |QR|   Opcode  |AA|TC|RD|RA|   Z    |   RCODE   |
//    +--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+
//    |                    QDCOUNT                    |
//    +--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+
//    |                    ANCOUNT                    |
//    +--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+
//    |                    NSCOUNT                    |
//    +--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+
//    |                    ARCOUNT                    |
//    +--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+
//
type Header struct {
	Id      uint16       // ID
	Flags   *HeaderFlags // 标识信息，共16位
	Qdcount uint16       // 报文请求段中的问题记录数
	Ancount uint16       // 报文回答段中的回答记录数
	Nscount uint16       // 报文授权段中的授权记录数
	Arcount uint16       // 报文附加段中的附加记录数
}

// 头部标识位信息
type HeaderFlags struct {
	Qr     int // 长度1位，值0是请求，1是应答
	Opcode int // 长度4位，值0是标准查询，1是反向查询，2死服务器状态查询。
	Aa     int // 长度1位，授权应答(Authoritative Answer) – 这个比特位在应答的时候才有意义，指出给出应答的服务器是查询域名的授权解析服务器。
	Tc     int // 长度1位，截断(TrunCation) – 用来指出报文比允许的长度还要长，导致被截断。
	Rd     int // 长度1位，期望递归(Recursion Desired) – 这个比特位被请求设置，应答的时候使用的相同的值返回。如果设置了RD，就建议域名服务器进行递归解析，递归查询的支持是可选的。
	Ra     int // 长度1位，支持递归(Recursion Available) – 这个比特位在应答中设置或取消，用来代表服务器是否支持递归查询。
	Zero   int // 长度3位，保留值，值为0.
	Rcode  int // 长度4位，返回码，通常为0(没有差错)和3(名字差错)
}

// 查询段结构
type Question struct {
	Name  string // 可变长查询域名
	Type  uint16 // 16位，查询资源类型
	Class uint16 // 16位，查询类别
}

// 应答结构,资源记录格式(Resource record)
type Resource struct {
	Name    string // 资源记录包含的域名
	Type    uint16 // 2个字节表示资源记录的类型，指出RDATA数据的含义
	Class   uint16 // 2个字节表示RDATA的类
	TTL     uint   // 4字节无符号整数表示资源记录可以缓存的时间。0代表只能被传输，但是不能被缓存。
	Rdlenth uint16 // 2个字节无符号整数表示RDATA的长度
	Rdata   string // 不定长字符串来表示记录，格式根TYPE和CLASS有关。比如，TYPE是A，CLASS 是 IN，那么RDATA就是一个4个字节的ARPA网络地址。
}

// 消息结构
type Msg struct {
	header   *Header
	question []Question
	Answer   []Resource
	Ns       []Resource
	Extra    []Resource
}

func (m *Msg) GetHeader() *Header {
	return m.header
}

func (m *Msg) GetQuestion(n uint) Question {
	return m.question[n]
}

// 消息解析
func UnpackMsg(buf []byte) *Msg {
	pkt := NewPacket(buf)
	msg := new(Msg)
	msg.header = unpackHeader(pkt)
	msg.question = unpackQuestion(pkt)

	return msg
}

func PackMsg(msg *Msg) []byte {
	pkt := NewPacket(make([]byte, 0, 1024))

	return pkt.Bytes()
}

// 解析包头（12字节）
func unpackHeader(pkt *packet) *Header {
	hd := new(Header)
	hd.Id = pkt.ReadUint16()
	hd.Flags = unpackFlags(pkt.ReadUint16())
	hd.Qdcount = pkt.ReadUint16()
	hd.Ancount = pkt.ReadUint16()
	hd.Nscount = pkt.ReadUint16()
	hd.Arcount = pkt.ReadUint16()

	return hd
}

// 解析头部标识位信息
func unpackFlags(si uint16) *HeaderFlags {
	i := int(si)
	flags := new(HeaderFlags)
	flags.Qr = i >> 15           // 第16位
	flags.Opcode = i >> 11 & 0xF // 第15-12位
	flags.Aa = i >> 10 & 0x1     // 第11位
	flags.Tc = i >> 9 & 0x1      // 第10位
	flags.Rd = i >> 8 & 0x1      // 第9位
	flags.Ra = i >> 7 & 0x1      // 第8位
	flags.Zero = 0
	flags.Rcode = i & 0xF // 最后4位

	return flags
}

// 解析查询段，只支持单个域名
func unpackQuestion(pkt *packet) []Question {
	question := make([]Question, 1)
	name := make([]string, 0)
	for {
		b := pkt.ReadByte()
		if b == 0 {
			break
		}
		name = append(name, string(pkt.ReadBytes(uint(b))))
	}
	debug(name)
	question[0].Name = strings.Join(name, ".")
	question[0].Type = pkt.ReadUint16()
	question[0].Class = pkt.ReadUint16()

	return question
}
