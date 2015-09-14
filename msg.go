package main

import (
	"strconv"
	"strings"
)

// http://www.rfc-editor.org/rfc/rfc1035.txt
//
// DNS报文格式：
//
// +---------------------------+
// | 报文头 |
// +---------------------------+
// | 问题　 | 向服务器提出的查询部分
// +---------------------------+
// | 回答　 | 服务器回复的资源记录
// +---------------------------+
// | 授权  | 权威的资源记录
// +---------------------------+
// | 格外的 | 格外的资源记录
// +---------------------------+
//
// 说明：查询包只有包头和问题部分，回复包是在查询包的基础上追加了回答、授权、额外资源部分，并且修改了包头的相关标识。
//
// 报文头（12字节）：
//                                 1  1  1  1  1  1
//   0  1  2  3  4  5  6  7  8  9  0  1  2  3  4  5
// +--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+
// |                      ID                       |
// +--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+
// |QR|   Opcode  |AA|TC|RD|RA|   Z    |   RCODE   |
// +--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+
// |                    QDCOUNT                    |
// +--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+
// |                    ANCOUNT                    |
// +--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+
// |                    NSCOUNT                    |
// +--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+
// |                    ARCOUNT                    |
// +--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+
//
// 查询部分结构：
//   0  1  2  3  4  5  6  7  8  9  0  1  2  3  4  5
// +--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+
// |                                               |
// /                     QNAME                     /
// /                                               /
// +--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+
// |                     QTYPE                     |
// +--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+
// |                     QCLASS                    |
// +--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+
//
// 应答资源结构（包括回答、授权、额外资源都使用同一结构）：
//                                 1  1  1  1  1  1
//   0  1  2  3  4  5  6  7  8  9  0  1  2  3  4  5
// +--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+
// |                                               |
// /                                               /
// /                      NAME                     /
// |                                               |
// +--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+
// |                      TYPE                     |
// +--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+
// |                     CLASS                     |
// +--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+
// |                      TTL                      |
// |                                               |
// +--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+
// |                   RDLENGTH                    |
// +--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+--|
// /                     RDATA                     /
// /                                               /
// +--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+
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
	Header   *Header
	Question []Question // 查询
	Answer   []Resource // 回答，服务器回复的资源记录
	Ns       []Resource // 授权，权威的资源记录
	Extra    []Resource // 额外的, 额外的资源记录
}

func (m *Msg) GetQuestion(n uint) Question {
	return m.Question[n]
}

// 设为回复
func (m *Msg) SetResponse() {
	m.Header.Flags.Qr = 1
}

func (m *Msg) AddAnswer(rs Resource) {
	m.Answer = append(m.Answer, rs)
	m.Header.Ancount++
}

func (m *Msg) AddNs(rs Resource) {
	m.Ns = append(m.Ns, rs)
}

func (m *Msg) AddExtra(rs Resource) {
	m.Extra = append(m.Extra, rs)
}

// 消息解析
func UnpackMsg(buf []byte) *Msg {
	pkt := NewPacket(buf)
	msg := new(Msg)
	msg.Header = unpackHeader(pkt)
	msg.Question = unpackQuestion(pkt)

	return msg
}

// 消息打包
func PackMsg(msg *Msg) []byte {
	pkt := NewPacket(make([]byte, 0, 1024))
	// 包头
	pkt.WriteBytes(packHeader(msg.Header))
	// 查询段
	for _, q := range msg.Question {
		pkt.WriteBytes(packQuestion(q))
	}
	for _, rs := range msg.Answer {
		pkt.WriteBytes(packResource(rs))
	}
	for _, rs := range msg.Ns {
		pkt.WriteBytes(packResource(rs))
	}
	for _, rs := range msg.Extra {
		pkt.WriteBytes(packResource(rs))
	}

	return pkt.Bytes()
}

// 解析包头（12字节）
func unpackHeader(pkt *Packet) *Header {
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
	flags.Zero = 0               // 第5-7位
	flags.Rcode = i & 0xF        // 最后4位

	return flags
}

// 解析查询段，只支持单个域名
func unpackQuestion(pkt *Packet) []Question {
	question := make([]Question, 1)
	name := make([]string, 0)
	for {
		b := pkt.ReadByte()
		if b == 0 {
			break
		}
		name = append(name, string(pkt.ReadBytes(uint(b))))
	}
	question[0].Name = strings.Join(name, ".")
	question[0].Type = pkt.ReadUint16()
	question[0].Class = pkt.ReadUint16()

	return question
}

// 打包包头
func packHeader(hd *Header) []byte {
	pkt := NewPacket(make([]byte, 0, 12))
	pkt.WriteUint16(hd.Id)
	flags := 0
	flags |= hd.Flags.Qr << 15 & 0x8000
	flags |= hd.Flags.Opcode << 11 & 0x7800
	flags |= hd.Flags.Aa << 10 & 0x400
	flags |= hd.Flags.Tc << 9 & 0x200
	flags |= hd.Flags.Rd << 8 & 0x100
	flags |= hd.Flags.Ra << 7 & 0x80
	flags |= hd.Flags.Zero << 4 & 0x70
	flags |= hd.Flags.Rcode & 0xF
	pkt.WriteUint16(uint16(flags))
	pkt.WriteUint16(hd.Qdcount)
	pkt.WriteUint16(hd.Ancount)
	pkt.WriteUint16(hd.Nscount)
	pkt.WriteUint16(hd.Arcount)

	return pkt.Bytes()
}

// 打包查询
func packQuestion(qs Question) []byte {
	pkt := NewPacket(make([]byte, 0, 256))
	pkt.WriteBytes(packName(qs.Name))
	pkt.WriteUint16(qs.Type)
	pkt.WriteUint16(qs.Class)

	return pkt.Bytes()
}

// 打包资源
func packResource(rs Resource) []byte {
	pkt := NewPacket(make([]byte, 0, 256))
	pkt.WriteBytes(packName(rs.Name))
	pkt.WriteUint16(rs.Type)
	pkt.WriteUint16(rs.Class)
	pkt.WriteUint(rs.TTL)
	pkt.WriteUint16(uint16(len(rs.Rdata)))
	pkt.WriteString(rs.Rdata)

	return pkt.Bytes()
}

// 打包域名
func packName(name string) []byte {
	parts := strings.Split(name, ".")
	buf := make([]byte, 0, 256)
	for _, v := range parts {
		buf = append(buf, byte(len(v)))
		buf = append(buf, []byte(v)...)
	}
	buf = append(buf, byte(0x0))

	return buf
}

func packIp(ip string) string {
	bs := make([]byte, 4)
	parts := strings.Split(ip, ".")
	for i := 0; i < 4; i++ {
		i32, _ := strconv.Atoi(parts[i])
		bs[i] = byte(i32)
	}

	return string(bs)
}

func NewA(domain string, ip string) Resource {
	var rs Resource
	rs.Name = domain
	rs.Type = TypeA
	rs.Class = ClassIN
	rs.TTL = 284
	rs.Rdlenth = 4
	rs.Rdata = packIp(ip)

	return rs
}
